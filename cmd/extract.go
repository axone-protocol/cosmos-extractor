package cmd

import (
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract data from a chain (snapshot)",
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
