package cmd

import (
	"github.com/axone-protocol/wallet-extractor/pkg/delegators"
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

		return delegators.Extract(chainName, args[0], args[1])
	},
}

func init() {
	extractCmd.AddCommand(extractDelegatorsCmd)

	extractDelegatorsCmd.Flags().StringP(flagChainName, "n", "cosmos", "Name of the chain")
}
