package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "rotate-me",
	Short:   "A CLI tool designed to rotate AWS Access Key",
	Version: "1.0.0", // TODO: Import version from package.json
}

func init() {
	RootCmd.PersistentFlags().IntP(
		"last-login",
		"l",
		0,
		"Search users that haven't logged since N days",
	)

	RootCmd.PersistentFlags().Int32P(
		"quantity",
		"n",
		50,
		"Test flag to pass value",
	)

	RootCmd.PersistentFlags().StringP(
		"user",
		"u",
		"fulano",
		"Test flag to pass value",
	)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
