package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/gyro/pkg/providers/aws"
	"github.com/spf13/cobra"
)

var Version string = "dev"

var RootCmd = &cobra.Command{
	Use:     "gyro",
	Short:   "A CLI tool designed to rotate AWS Access Key and user credentials",
	Version: Version,
}

type BaseCommandOptions struct {
	Quantity int32
	Path     string
	User     string
	TimeZone string
	Format   string
	Age      int
	Expired  bool
}

type RotateCommandOptions struct {
	BaseCommandOptions
	ExpireOnly bool
	DryRun     bool
	Notify     bool
	// AutoApprove bool
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

func configureListFlags(cmd *cobra.Command) BaseCommandOptions {
	quantity, _ := cmd.Flags().GetInt32("quantity")
	timeZone, _ := cmd.Flags().GetString("timezone")
	format, _ := cmd.Flags().GetString("format")
	userName, _ := cmd.Flags().GetString("username")
	path, _ := cmd.Flags().GetString("output-file")
	age, _ := cmd.Flags().GetInt("age")
	expired, _ := cmd.Flags().GetBool("expired-only")

	return BaseCommandOptions{
		Quantity: quantity,
		User:     userName,
		TimeZone: timeZone,
		Format:   format,
		Path:     path,
		Age:      age,
		Expired:  expired,
	}
}

func configureListCommand(cmd *cobra.Command) (iam.GetWrapperInputs, BaseCommandOptions) {
	options := configureListFlags(cmd)

	wrapper := iam.UserWrapper{
		IamClient: iam.DeclareConfig(),
	}

	inputs := iam.GetWrapperInputs{
		MaxUsers: options.Quantity,
		TimeZone: options.TimeZone,
		Age:      options.Age,
		Expired:  options.Expired,
		UserName: options.User,
		Client:   wrapper,
	}
	return inputs, options
}

func configureRotateCommand(cmd *cobra.Command) (RotateCommandOptions, BaseCommandOptions) {
	listOptions := configureListFlags(cmd)
	expireOnly, _ := cmd.Flags().GetBool("expire-only")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	notify, _ := cmd.Flags().GetBool("notify")

	return RotateCommandOptions{
		BaseCommandOptions: listOptions,
		ExpireOnly:         expireOnly,
		DryRun:             dryRun,
		Notify:             notify,
	}, listOptions
}

func initializeBaseCommandFlags(cmd *cobra.Command) {
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
