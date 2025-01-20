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

// extractCmd represents the command to extract data from a chain. It is a parent command for all extract commands.
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract data from a chain (snapshot)",
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.PersistentFlags().StringP(flagChainName, "n", "cosmos", "name of the chain")
	extractCmd.PersistentFlags().StringP(flagOutput, "o", "", "output file (defaults to stdout)")
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
