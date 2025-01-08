package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/infos"
	"github.com/spf13/cobra"
)

var extractInfosCmd = &cobra.Command{
	Use:   "infos [source] [dest]",
	Short: "Extract chain informations into a CSV file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chainName, _ := cmd.Flags().GetString(flagChainName)
		pipeline, err := infos.Pipeline(chainName, args[0], args[1], logger)
		if err != nil {
			return err
		}

		err = <-pipeline.Run()

		return err
	},
}

func init() {
	extractCmd.AddCommand(extractInfosCmd)
}
