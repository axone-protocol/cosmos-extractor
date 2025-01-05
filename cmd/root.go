package cmd

import (
	"os"

	pkglogger "github.com/axone-protocol/cosmos-extractor/pkg/logger"
	"github.com/spf13/cobra"

	"cosmossdk.io/log"
)

var logger = log.NewNopLogger()

const (
	flagLogLevel = "log_level"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "cosmos-extractor",
	Short: "A template fo Golang projects",
	Long:  "A template fo Golang projects",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		logLevel := cmd.Flag(flagLogLevel).Value.String()
		filterFn, err := log.ParseLogLevel(logLevel)
		if err != nil {
			return nil
		}
		logger = log.NewLogger(os.Stderr, log.FilterOption(filterFn))
		pkglogger.InstallETLLogger(logger)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String(flagLogLevel, "info", "The logging level (trace|debug|info|warn|error|fatal|panic)")
}
