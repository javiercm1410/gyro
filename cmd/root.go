package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "rotate-me",
	Short:   "A CLI tool designed to rotate AWS Access Key and users credentials",
	Version: "1.0.0", // TODO: Import version from package.json
}

func init() {
	RootCmd.PersistentFlags().Int32P(
		"quantity",
		"n",
		50,
		"User number to be listed`",
	)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
