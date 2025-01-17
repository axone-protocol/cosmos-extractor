package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/delegators"
	"github.com/spf13/cobra"
	"github.com/teambenny/goetl"
)

const (
	flagHrp = "hrp"
)

var extractDelegatorsCmd = &cobra.Command{
	Use:   "delegators [source]",
	Short: "Extract all delegators",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chainName, _ := cmd.Flags().GetString(flagChainName)
		src := args[0]

		processors := []goetl.Processor{}

		read, err := delegators.NewDelegatorsReader(chainName, src, logger)
		if err != nil {
			return err
		}
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

func init() {
	extractCmd.AddCommand(extractDelegatorsCmd)

	extractDelegatorsCmd.Flags().StringSliceP(
		flagHrp,
		"p",
		[]string{},
		"One or more Human-Readable Parts (HRPs) to append delegator addresses in the given Bech32 formats (e.g., cosmos, osmo). "+
			"Can be used multiple times for different HRPs.")
}
