package cmd

import (
	"encoding/json"
	"strings"

	"github.com/axone-protocol/cosmos-extractor/internal/version"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	flagLong   = "long"
	flagFormat = "format"
)

// versionCmd represents the command to interactively print the application binary version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the application binary version information",
	RunE: func(cmd *cobra.Command, _ []string) error {
		verInfo := version.NewInfo()

		if long, _ := cmd.Flags().GetBool(flagLong); !long {
			cmd.Println(verInfo.Version)
			return nil
		}

		var (
			bz  []byte
			err error
		)

		format, _ := cmd.Flags().GetString(flagFormat)
		switch strings.ToLower(format) {
		case "json":
			bz, err = json.Marshal(verInfo)

		default:
			bz, err = yaml.Marshal(&verInfo)
		}

		if err != nil {
			return err
		}

		cmd.Println(string(bz))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().Bool(flagLong, false, "print long version information")
	versionCmd.Flags().StringP(flagFormat, "f", "text", "output format (text|json)")
}
