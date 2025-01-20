package cmd

import (
	"github.com/spf13/cobra"
)

// metaCmd represents the meta command, which is a parent command for all meta commands.
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Meta commands for the tool",
}

func init() {
	rootCmd.AddCommand(metaCmd)
}
