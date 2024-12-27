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

type ListCommandOptions struct {
	Quantity int32
	Path     string
	User     string
	TimeZone string
	Format   string
	Age      int
	Expired  bool
}

func init() {
	RootCmd.PersistentFlags().Int32P(
		"quantity",
		"n",
		50,
		"Number of users to be listed",
	)

	RootCmd.PersistentFlags().BoolP(
		"debug",
		"d",
		false,
		"Show detailed output",
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

func configureListFlags(cmd *cobra.Command) ListCommandOptions {
	quantity, _ := cmd.Flags().GetInt32("quantity")
	timeZone, _ := cmd.Flags().GetString("timezone")
	format, _ := cmd.Flags().GetString("format")
	userName, _ := cmd.Flags().GetString("username")
	path, _ := cmd.Flags().GetString("output-file")
	age, _ := cmd.Flags().GetInt("age")
	expired, _ := cmd.Flags().GetBool("expired-only")

	return ListCommandOptions{
		Quantity: quantity,
		User:     userName,
		TimeZone: timeZone,
		Format:   format,
		Path:     path,
		Age:      age,
		Expired:  expired,
	}
}

func initializeListCommandFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("timezone", "t", "America/Santo_Domingo", "Timezone for displaying dates")
	cmd.PersistentFlags().StringP("format", "f", "table", "Output format (json, table, file)")
	cmd.PersistentFlags().StringP("output-file", "o", "./output.json", "Save results to file")
	cmd.PersistentFlags().StringP("username", "u", "", "Filter by specific IAM username")
	cmd.PersistentFlags().IntP("age", "a", 90, "Consider keys stale after N days")
	cmd.PersistentFlags().BoolP("expired-only", "x", false, "Show only expired keys")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		age, _ := cmd.Flags().GetInt("age")
		if age < 1 {
			return fmt.Errorf("age must be greater than 0, got %d", age)
		}

		format, _ := cmd.Flags().GetString("format")
		validFormats := map[string]bool{"json": true, "table": true, "text": true}
		if !validFormats[format] {
			return fmt.Errorf("invalid format '%s'. Valid options are: json, table, text", format)
		}

		timeZone, _ := cmd.Flags().GetString("timezone")
		if timeZone == "" {
			return fmt.Errorf("timezone cannot be empty")
		}

		return nil
	}
}
