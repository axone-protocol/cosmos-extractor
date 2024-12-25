package delegators

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/axone-protocol/wallet-extractor/pkg/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/gocarina/gocsv"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func Extract(chainName, src, dst string) error {
	logger := log.NewLogger(os.Stderr)

	logger.Info("Extracting wallets", "source", src, "destination", dst)

	keepers, err := keeper.OpenStore(src, logger)
	if err != nil {
		return err
	}

	ctx := sdk.NewContext(keepers.Store, cmtproto.Header{}, false, keepers.Logger)

	validators, err := keepers.Staking.GetAllValidators(ctx)
	if err != nil {
		panic(err)
	}

	keepers.Logger.Info("Analyzing validators", "count", len(validators))


	return extractDelegators(ctx, chainName, keepers, validators, dst)
}

func extractDelegators(
	ctx sdk.Context, chainName string, keepers *keeper.Keepers, validators []stakingtypes.Validator, destination string) error {
	file, err := os.OpenFile(path.Join(destination, "delegations.csv"), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	prefix, err := guessPrefixFromValoper(validators[0].OperatorAddress)
	if err != nil {
		return err
	}

	config := sdk.GetConfig()

	keepers.Bank.IterateAllBalances(ctx, func(addr sdk.AccAddress, _ sdk.Coin) (stop bool) {
		for _, val := range validators {
			if config.GetBech32AccountAddrPrefix() != prefix {
				config.SetBech32PrefixForValidator(
					fmt.Sprintf("%svaloper", prefix),
					fmt.Sprintf("%svaloperpub", prefix),
				)
			}

			valAddr, err := sdk.ValAddressFromBech32(val.OperatorAddress)
			if err != nil {
				panic(err)
			}
			delegation, err := keepers.Staking.GetDelegation(ctx, addr, valAddr)
			if err != nil {
				continue
			}

			record := Delegations{
				ChainName:           chainName,
				DelegatorNativeAddr: delegation.DelegatorAddress,
				DelegatorCosmosAddr: convertAndEncodeMust("cosmos", delegation.DelegatorAddress),
				DelegatorAxoneAddr:  convertAndEncodeMust("axone", delegation.DelegatorAddress),
				ValidatorAddr:       delegation.ValidatorAddress,
				Shares:              delegation.Shares.String(),
			}

			v, err := gocsv.MarshalStringWithoutHeaders(&[]Delegations{record})
			if err != nil {
				panic(err)
			}

			_, err = writer.WriteString(v)
			if err != nil {
				panic(err)
			}
		}

		return false
	})

	return nil
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
