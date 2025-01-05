package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/delegators"
	"github.com/spf13/cobra"
)

const (
	flagChainName = "chain-name"
)

var extractDelegatorsCmd = &cobra.Command{
	Use:   "delegators [source] [dest]",
	Short: "Extract all delegators into CSV files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chainName, _ := cmd.Flags().GetString(flagChainName)
		pipeline, err := delegators.Pipeline(chainName, args[0], args[1], logger)
		if err != nil {
			return err
		}

		err = <-pipeline.Run()

		return err
	},
}

func init() {
	extractCmd.AddCommand(extractDelegatorsCmd)

	extractDelegatorsCmd.Flags().StringP(flagChainName, "n", "cosmos", "Name of the chain")
}
