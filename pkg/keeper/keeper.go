package keeper

import (
	"fmt"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Keepers struct {
	dbPath string
	codec  *codec.ProtoCodec
	keys   map[string]*storetypes.KVStoreKey

	Logger log.Logger
	Store  storetypes.CommitMultiStore

	Account authkeeper.AccountKeeper
	Bank    bankkeeper.BaseKeeper
	Staking *stakingkeeper.Keeper
}

// OpenStore opens an existing store at the given path and returns the keepers for the store.
func OpenStore(dbPath string, logger log.Logger) (*Keepers, error) {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
	)

	codec, err := newCodec()
	if err != nil {
		return nil, err
	}

	authKeeper := newAuthKeeper(codec, keys)
	bankKeeper := newBankKeeper(codec, keys, authKeeper, logger)
	stakingKeeper := newStakingKeeper(codec, keys, authKeeper, bankKeeper)

	store, err := newCommitMultiStore(dbPath, keys, logger)
	if err != nil {
		return nil, err
	}

	return &Keepers{
		dbPath:  dbPath,
		codec:   codec,
		keys:    keys,
		Logger:  logger,
		Store:   store,
		Account: authKeeper,
		Bank:    bankKeeper,
		Staking: stakingKeeper,
	}, nil
}

var modulePermissions = map[string][]string{
	authtypes.FeeCollectorName:     nil,
	stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
}

func newCodec() (*codec.ProtoCodec, error) {
	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterInterface("/cosmos.auth.v1beta1.BaseAccount", (*sdk.AccountI)(nil))

	return codec.NewProtoCodec(interfaceRegistry), nil
}

func newAuthKeeper(
	codec *codec.ProtoCodec, keys map[string]*storetypes.KVStoreKey,
) authkeeper.AccountKeeper {
	return authkeeper.NewAccountKeeper(
		codec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		modulePermissions,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

func newBankKeeper(
	codec *codec.ProtoCodec, keys map[string]*storetypes.KVStoreKey, authKeeper authkeeper.AccountKeeper, logger log.Logger,
) bankkeeper.BaseKeeper {
	return bankkeeper.NewBaseKeeper(
		codec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		authKeeper,
		map[string]bool{},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)
}

func newStakingKeeper(
	codec *codec.ProtoCodec, keys map[string]*storetypes.KVStoreKey, authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.BaseKeeper,
) *stakingkeeper.Keeper {
	return stakingkeeper.NewKeeper(
		codec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		authKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
}

func newCommitMultiStore(
	dbPath string, keys map[string]*storetypes.KVStoreKey, logger log.Logger,
) (storetypes.CommitMultiStore, error) {
	logger.Debug("Opening store", "path", dbPath)

	db, err := dbm.NewDB("application", dbm.GoLevelDBBackend, dbPath)
	if err != nil {
		return nil, err
	}

	ms := store.NewCommitMultiStore(db, logger, metrics.NewNoOpMetrics())
	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, nil)
	}

	if err := ms.LoadLatestVersion(); err != nil {
		return nil, err
	}

	commitID := ms.LastCommitID()

	logger.Debug("Store loaded", "version", commitID.Version, "hash", fmt.Sprintf("%x", commitID.Hash))

	return ms, nil
}
