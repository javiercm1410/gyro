package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/rotator/pkg/providers/aws"
	"github.com/spf13/cobra"
)

type ListCommandOptions struct {
	LastLogin int
	Quantity  int32
	User      string
	TimeZone  string
	Output    string
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get IAM Access Keys",
	Aliases: []string{"users", "accounts"},
	Example: "rotate-me list",
	Run: func(cmd *cobra.Command, args []string) {
		options := configureListFlags(cmd)
		var wrapper iam.UserWrapper
		wrapper.IamClient = iam.DeclareConfig()

		inputs := iam.GetUserAccessKeyInputs{
			MaxUsers:   options.Quantity,
			TimeZone:   options.TimeZone,
			OutputType: options.Output,
			Client:     wrapper,
		}

		err := iam.GetUserAccessKey(inputs)
		if err != nil {
			log.Errorf("Failed to get users. Here's why: %v\n", err)
			os.Exit(1)
		}
	},
}

func configureListFlags(cmd *cobra.Command) ListCommandOptions {
	lastLogin, _ := cmd.Flags().GetInt("last-login")
	quantity, _ := cmd.Flags().GetInt32("quantity")
	userName, _ := cmd.Flags().GetString("user")
	timeZone, _ := cmd.Flags().GetString("timezone")
	output, _ := cmd.Flags().GetString("output")

	return ListCommandOptions{
		LastLogin: lastLogin,
		Quantity:  quantity,
		User:      userName,
		TimeZone:  timeZone,
		Output:    output,
	}
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("access-key", "k", false, "Vault ID")
	listCmd.PersistentFlags().BoolP("passwords", "w", false, "Vault ID")
	listCmd.PersistentFlags().BoolP("all", "a", false, "get both access-key and password")
	listCmd.PersistentFlags().StringP("timezone", "t", "America/Santo_Domingo", "get both access-key and password")
	listCmd.PersistentFlags().StringP("output", "o", "json", "json, table, text, default table")
	listCmd.PersistentFlags().StringP("path", "p", "json", "Json, table, text, default table")
}
