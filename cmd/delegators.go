package cmd

import (
	"github.com/axone-protocol/wallet-extractor/pkg/delegators"
	"github.com/spf13/cobra"
)

var extractDelegatorsCmd = &cobra.Command{
	Use:   "delegators [source] [dest]",
	Short: "Extract all delegators",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		return delegators.Extract(args[0], args[1])
	},
}

func init() {
	extractCmd.AddCommand(extractDelegatorsCmd)
}
