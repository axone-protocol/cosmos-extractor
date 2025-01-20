package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	flagMkDirs = "mk-dirs"
)

// metaCmd represents the meta docs command to generate documentation.
var metaDocsCmd = &cobra.Command{
	Use:   "docs [dir]",
	Short: "Generate documentation to a directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := args[0]

		shouldMkDirs, err := cmd.Flags().GetBool(flagMkDirs)
		if err != nil {
			return err
		}
		if shouldMkDirs {
			if err := mkDirs(targetPath); err != nil {
				return err
			}
		}

		rootCmd.DisableAutoGenTag = true

		err = doc.GenMarkdownTree(rootCmd, targetPath)
		if err != nil {
			return err
		}

		return nil
	},
}

func mkDirs(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func init() {
	metaCmd.AddCommand(metaDocsCmd)
	metaDocsCmd.Flags().Bool(flagMkDirs, false, "create directories if they do not exist")
}
