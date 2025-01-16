package delegators

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/axone-protocol/cosmos-extractor/pkg/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"
	"github.com/teambenny/goetl/etlutil"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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

	lastSeenAddr := ""
	keepers.Bank.IterateAllBalances(ctx, func(addr sdk.AccAddress, _ sdk.Coin) (stop bool) {
		if addr.String() == lastSeenAddr {
			return false
		}
		lastSeenAddr = addr.String()

		for _, val := range validators {
			valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
			etlutil.KillPipelineIfErr(err, killChan)

			delegation, err := keepers.Staking.GetDelegation(ctx, addr, valAddr)
			if err != nil {
				if errors.Is(err, stakingtypes.ErrNoDelegation) {
					continue
				}

				r.logger.Error(err.Error())
				killChan <- err
				return true
			}

			payload := Delegation{
				ChainName:           r.chainName,
				DelegatorNativeAddr: delegation.DelegatorAddress,
				DelegatorCosmosAddr: convertAndEncodeMust("cosmos", delegation.DelegatorAddress),
				DelegatorAxoneAddr:  convertAndEncodeMust("axone", delegation.DelegatorAddress),
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

func convertAndEncodeMust(hrp string, bech string) string {
	_, bytes, err := bech32.DecodeAndConvert(bech)
	if err != nil {
		panic(err)
	}

	encoded, err := bech32.ConvertAndEncode(hrp, bytes)
	if err != nil {
		panic(err)
	}

	return encoded
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
