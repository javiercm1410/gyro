package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "gyro",
	Short:   "A CLI tool designed to rotate AWS Access Key and user credentials",
	Version: getVersion(),
}

func init() {
	RootCmd.PersistentFlags().Int32P(
		"quantity",
		"n",
		50,
		"Number of users to be listed",
	)

	// Add validation example:
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		quantity, _ := cmd.Flags().GetInt32("quantity")
		if quantity < 1 {
			return fmt.Errorf("quantity must be a positive number, got %d", quantity)
		}
		return nil
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

func getVersion() string {
	// Placeholder for dynamic version retrieval logic
	return "1.0.0"
}
