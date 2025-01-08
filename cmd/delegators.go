package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/delegators"
	"github.com/spf13/cobra"
)

var extractDelegatorsCmd = &cobra.Command{
	Use:   "delegators [source] [dest]",
	Short: "Extract all delegators into a CSV file",
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
}
