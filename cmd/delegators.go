package cmd

import (
	"github.com/axone-protocol/cosmos-extractor/pkg/delegators"
	"github.com/spf13/cobra"
	"github.com/teambenny/goetl"
)

var extractDelegatorsCmd = &cobra.Command{
	Use:   "delegators [source]",
	Short: "Extract all delegators",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chainName, _ := cmd.Flags().GetString(flagChainName)
		src := args[0]

		read, err := delegators.NewDelegatorsReader(chainName, src, logger)
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
	extractCmd.AddCommand(extractDelegatorsCmd)
}
