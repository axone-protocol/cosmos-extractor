package cmd

import (
	"github.com/spf13/cobra"
)

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Meta commands for the tool",
}

func init() {
	rootCmd.AddCommand(metaCmd)
}
