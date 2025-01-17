package delegators

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/axone-protocol/cosmos-extractor/pkg/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/samber/lo"
	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"

	"cosmossdk.io/collections"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type delegatorsReader struct {
	chainName string
	src       string
	logger    log.Logger
	closer    io.Closer
}

// NewDelegatorsReader returns a new Reader that reads delegators from a blockchain data stores.
func NewDelegatorsReader(chainName, src string, logger log.Logger) (goetl.Processor, error) {
	return &delegatorsReader{
		chainName: chainName,
		src:       src,
		logger:    logger,
	}, nil
}

func (r *delegatorsReader) ProcessData(_ etldata.Payload, outputChan chan etldata.Payload, killChan chan error) {
	keepers, err := keeper.OpenStore(r.src, r.logger)
	if err != nil {
		r.logger.Error(err.Error())
		killChan <- err
		return
	}
	r.closer = keepers

	ctx := sdk.NewContext(keepers.Store, cmtproto.Header{}, false, keepers.Logger)

	validators, err := keepers.Staking.GetAllValidators(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		killChan <- err
		return
	}

	prefix, err := guessPrefixFromValoper(validators[0].OperatorAddress)
	if err != nil {
		r.logger.Error(err.Error())
		killChan <- err
		return
	}

	configureSdk(prefix)

	err = iterateAllAddresses(ctx, keepers.Bank, func(addr sdk.AccAddress) (stop bool) {
		delegations := lo.RejectMap(validators,
			toDelegations(ctx, addr, r.logger, keepers.Staking, killChan))

		for _, delegation := range delegations {
			payload := Delegation{
				ChainName:           r.chainName,
				DelegatorNativeAddr: delegation.DelegatorAddress,
				ValidatorAddr:       delegation.ValidatorAddress,
				Shares:              delegation.Shares.String(),
			}

			json, err := etldata.NewJSON(payload)
			if err != nil {
				r.logger.Error(err.Error())
				killChan <- err
				return true
			}

			outputChan <- json
		}

		return false
	})
	if err != nil {
		r.logger.Error(err.Error())
		killChan <- err
	}
}

func (r *delegatorsReader) Finish(_ chan etldata.Payload, killChan chan error) {
	if r.closer != nil {
		err := r.closer.Close()
		if err != nil {
			r.logger.Error(err.Error())
			killChan <- err
		}
	}
}

func (r *delegatorsReader) String() string {
	return "DelegatorsReader"
}

// IterateAllAddresses iterates over all the accounts that are provided to a callback.
// If true is returned from the callback, iteration is halted.
func iterateAllAddresses(ctx context.Context, bankKeeper bankkeeper.BaseKeeper, cb func(sdk.AccAddress) bool) error {
	lastSeenAddr := ""
	err := bankKeeper.Balances.Walk(ctx, nil, func(key collections.Pair[sdk.AccAddress, string], _ math.Int) (stop bool, err error) {
		addr := key.K1()
		if addr.String() == lastSeenAddr {
			return false, nil
		}
		lastSeenAddr = addr.String()

		return cb(addr), nil
	})

	return err
}

func toDelegations(
	ctx context.Context, address sdk.AccAddress, logger log.Logger, stakingKeeper *stakingkeeper.Keeper, killChan chan error,
) func(item stakingtypes.Validator, index int) (stakingtypes.Delegation, bool) {
	return func(item stakingtypes.Validator, _ int) (stakingtypes.Delegation, bool) {
		valAddr, err := sdk.ValAddressFromBech32(item.OperatorAddress)
		if err != nil {
			logger.Error(err.Error())
			killChan <- err
			return stakingtypes.Delegation{}, true
		}

		delegation, err := stakingKeeper.GetDelegation(ctx, address, valAddr)
		if err != nil {
			if errors.Is(err, stakingtypes.ErrNoDelegation) {
				return stakingtypes.Delegation{}, true
			}

			logger.Error(err.Error())
			killChan <- err
			return stakingtypes.Delegation{}, true
		}
		return delegation, false
	}
}

func guessPrefixFromValoper(valoper string) (string, error) {
	if idx := strings.Index(valoper, "valoper"); idx != -1 {
		return valoper[:idx], nil
	}
	return "", fmt.Errorf("valoper not found in operator address: %s", valoper)
}

func configureSdk(prefix string) {
	config := sdk.GetConfig()
	if config.GetBech32AccountAddrPrefix() != prefix {
		config.SetBech32PrefixForValidator(
			fmt.Sprintf("%svaloper", prefix),
			fmt.Sprintf("%svaloperpub", prefix),
		)
	}
}
