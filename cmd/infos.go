package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/infos"
	"github.com/spf13/cobra"
	"github.com/teambenny/goetl"
)

// extractInfosCmd represents the command to extract chain information.
var extractInfosCmd = &cobra.Command{
	Use:   "infos [source]",
	Short: "Extract chain information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chainName, _ := cmd.Flags().GetString(flagChainName)
		src := args[0]

		read, err := infos.NewInfoReader(chainName, src, logger)
		if err != nil {
			return err
		}

		write, err := newCSVWriter(cmd, args)
		if err != nil {
			return err
		}

		pipeline := goetl.NewPipeline(read, write)
		pipeline.Name = "Delegators"

		err = <-pipeline.Run()
		return err
	},
}

func init() {
	extractCmd.AddCommand(extractInfosCmd)
}
