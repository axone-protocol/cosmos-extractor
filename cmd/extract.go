package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagChainName = "chain-name"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract data from a chain (snapshot)",
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.PersistentFlags().StringP(flagChainName, "n", "cosmos", "Name of the chain")
}
