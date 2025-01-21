package delegators

import (
	"context"
	"io"

	"github.com/axone-protocol/cosmos-extractor/pkg/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/samber/lo"
	"github.com/teambenny/goetl"
	"github.com/teambenny/goetl/etldata"

	"cosmossdk.io/log"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type ReaderOption func(*delegatorsReader) error

func WithChainName(chainName string) ReaderOption {
	return func(r *delegatorsReader) error {
		r.chainName = chainName
		return nil
	}
}

func WithLogger(logger log.Logger) ReaderOption {
	return func(r *delegatorsReader) error {
		r.logger = logger
		return nil
	}
}

func WithMinSharesFilter(minShares math.LegacyDec) ReaderOption {
	return func(r *delegatorsReader) error {
		r.minSharesFilter = minShares
		return nil
	}
}

func WithMaxSharesFilter(maxShares math.LegacyDec) ReaderOption {
	return func(r *delegatorsReader) error {
		r.maxSharesFilter = maxShares
		return nil
	}
}

type delegatorsReader struct {
	chainName       string
	src             string
	logger          log.Logger
	closer          io.Closer
	minSharesFilter math.LegacyDec
	maxSharesFilter math.LegacyDec
}

// NewDelegatorsReader returns a new Reader that reads delegators from a blockchain data stores.
func NewDelegatorsReader(src string, options ...ReaderOption) (goetl.Processor, error) {
	r := &delegatorsReader{
		chainName: "mystery",
		src:       src,
		logger:    log.NewNopLogger(),
	}

	for _, option := range options {
		if err := option(r); err != nil {
			return nil, err
		}
	}

	return r, nil
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

	err = keepers.Account.Accounts.Walk(ctx, nil, func(addr sdk.AccAddress, _ sdk.AccountI) (stop bool, err error) {
		delegations := lo.FlatMap([]sdk.AccAddress{addr},
			extractDelegations(ctx, r.logger, keepers.Staking, killChan),
		)
		shares := lo.Reduce(delegations, computeShares(), math.LegacyZeroDec())

		if (!r.maxSharesFilter.IsNil() && shares.GT(r.maxSharesFilter)) ||
			(!r.minSharesFilter.IsNil() && shares.LT(r.minSharesFilter)) {
			return false, nil
		}

		for _, delegation := range delegations {
			json, err := delegatorToPayload(r, delegation)
			if err != nil {
				return true, err
			}

			outputChan <- json
		}

		return false, nil
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

func extractDelegations(
	ctx context.Context, logger log.Logger, stakingKeeper *stakingkeeper.Keeper, killChan chan error,
) func(delegator sdk.AccAddress, _ int) []stakingtypes.Delegation {
	return func(delegator sdk.AccAddress, _ int) []stakingtypes.Delegation {
		delegation, err := stakingKeeper.GetAllDelegatorDelegations(ctx, delegator)
		if err != nil {
			logger.Error(err.Error())
			killChan <- err
			return nil
		}
		return delegation
	}
}

func computeShares() func(acc math.LegacyDec, delegation stakingtypes.Delegation, _ int) math.LegacyDec {
	return func(acc math.LegacyDec, delegation stakingtypes.Delegation, _ int) math.LegacyDec {
		return acc.Add(delegation.Shares)
	}
}

func delegatorToPayload(r *delegatorsReader, delegation stakingtypes.Delegation) (etldata.JSON, error) {
	payload := Delegation{
		ChainName:           r.chainName,
		DelegatorNativeAddr: delegation.DelegatorAddress,
		ValidatorAddr:       delegation.ValidatorAddress,
		Shares:              delegation.Shares.String(),
	}

	json, err := etldata.NewJSON(payload)
	return json, err
}
