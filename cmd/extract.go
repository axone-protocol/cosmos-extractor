package cmd

import (
	"os"

	"github.com/axone-protocol/cosmos-extractor/pkg/csv"
	"github.com/spf13/cobra"
	"github.com/teambenny/goetl"
)

const (
	flagChainName = "chain-name"
	flagOutput    = "output"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract data from a chain (snapshot)",
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.PersistentFlags().StringP(flagChainName, "n", "cosmos", "Name of the chain")
	extractCmd.PersistentFlags().StringP(flagOutput, "o", "", "Output file (defaults to stdout)")
}

func newCSVWriter(cmd *cobra.Command, _ []string) (goetl.Processor, error) {
	options := []csv.Option{
		csv.WithWriterHeader(),
	}
	output, _ := cmd.Flags().GetString(flagOutput)
	if output != "" {
		options = append(options, csv.WithFile(output))
	} else {
		options = append(options, csv.WithWriter(os.Stdout))
	}

	return csv.NewCSVWriter(options...)
}
