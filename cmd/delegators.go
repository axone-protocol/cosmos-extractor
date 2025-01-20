package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/delegators"
	"github.com/spf13/cobra"
	"github.com/teambenny/goetl"

	"cosmossdk.io/math"
)

const (
	flagHrp       = "hrp"
	flagMinShares = "min-shares"
	flagMaxShares = "max-shares"
)

// extractDelegatorsCmd represents the command to extract all delegators.
var extractDelegatorsCmd = &cobra.Command{
	Use:   "delegators [source]",
	Short: "Extract all delegators",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		read, err := newDelegatorsReader(cmd, args)
		if err != nil {
			return err
		}

		processors := []goetl.Processor{}
		processors = append(processors, read)

		hrps, err := cmd.Flags().GetStringSlice(flagHrp)
		if err != nil {
			return err
		}
		if len(hrps) != 0 {
			p, err := delegators.NewAddressEnhancer(hrps, logger)
			if err != nil {
				return err
			}
			processors = append(processors, p)
		}

		write, err := newCSVWriter(cmd, args)
		if err != nil {
			return err
		}
		processors = append(processors, write)

		pipeline := goetl.NewPipeline(processors...)
		pipeline.Name = "Delegators"

		err = <-pipeline.Run()
		return err
	},
}

func newDelegatorsReader(cmd *cobra.Command, args []string) (goetl.Processor, error) {
	chainName, _ := cmd.Flags().GetString(flagChainName)
	src := args[0]

	delegatorsReaderOpts := []delegators.ReaderOption{
		delegators.WithChainName(chainName),
		delegators.WithLogger(logger),
	}

	v, err := getShares(cmd, flagMinShares)
	if err != nil {
		return nil, err
	}
	if !v.IsNil() {
		delegatorsReaderOpts = append(delegatorsReaderOpts, delegators.WithMinSharesFilter(v))
	}

	v, err = getShares(cmd, flagMaxShares)
	if err != nil {
		return nil, err
	}

	if !v.IsNil() {
		delegatorsReaderOpts = append(delegatorsReaderOpts, delegators.WithMaxSharesFilter(v))
	}

	return delegators.NewDelegatorsReader(src, delegatorsReaderOpts...)
}

func getShares(cmd *cobra.Command, flag string) (math.LegacyDec, error) {
	shares, err := cmd.Flags().GetString(flag)
	if err != nil {
		return math.LegacyDec{}, err
	}
	if shares == "" {
		return math.LegacyDec{}, nil
	}
	return math.LegacyNewDecFromStr(shares)
}

func init() {
	extractCmd.AddCommand(extractDelegatorsCmd)

	extractDelegatorsCmd.Flags().StringSliceP(
		flagHrp,
		"p",
		[]string{},
		"one or more Human-Readable Parts (HRPs) to append delegator addresses in the given Bech32 formats (e.g., cosmos, osmo). "+
			"Can be used multiple times for different HRPs.")

	extractDelegatorsCmd.Flags().String(flagMinShares, "", "filter delegators with minimum shares (in native token)")
	extractDelegatorsCmd.Flags().String(flagMaxShares, "", "filter delegators with maximum shares (in native token)")
}
