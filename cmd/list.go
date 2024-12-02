package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/rotator/pkg/providers/aws"
	"github.com/javiercm1410/rotator/pkg/utils"
	"github.com/spf13/cobra"
)

type ListCommandOptions struct {
	LastLogin int
	Quantity  int32
	Path      string
	User      string
	TimeZone  string
	Output    string
	Stale     int
	Expired   bool
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
			MaxUsers: options.Quantity,
			TimeZone: options.TimeZone,
			Stale:    options.Stale,
			Expired:  options.Expired,
			UserName: options.User,
			Client:   wrapper,
		}

		userKeyData, err := iam.GetUserAccessKey(inputs)
		if err != nil {
			log.Errorf("Failed to get users. Here's why: %v\n", err)
			os.Exit(1)
		}

		utils.DisplayData(options.Output, options.Path, options.Stale, userKeyData)
	},
}

func configureListFlags(cmd *cobra.Command) ListCommandOptions {
	lastLogin, _ := cmd.Flags().GetInt("last-login")
	quantity, _ := cmd.Flags().GetInt32("quantity")
	timeZone, _ := cmd.Flags().GetString("timezone")
	output, _ := cmd.Flags().GetString("output")
	userName, _ := cmd.Flags().GetString("user")
	path, _ := cmd.Flags().GetString("path")
	stale, _ := cmd.Flags().GetInt("stale")
	expired, _ := cmd.Flags().GetBool("expired")

	return ListCommandOptions{
		LastLogin: lastLogin,
		Quantity:  quantity,
		User:      userName,
		TimeZone:  timeZone,
		Output:    output,
		Path:      path,
		Stale:     stale,
		Expired:   expired,
	}
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringP("timezone", "t", "America/Santo_Domingo", "Select the timezone to display the dates")
	listCmd.PersistentFlags().StringP("output", "o", "json", "Select the output format (json, table, text)")
	listCmd.PersistentFlags().StringP("path", "p", "./output.json", "Path to save the output file")
	listCmd.PersistentFlags().StringP("user", "u", "", "Select the user to get the access keys")
	listCmd.PersistentFlags().IntP("stale", "s", 90, "Select the stale days to get the expired access keys")
	listCmd.PersistentFlags().BoolP("expired", "x", false, "Only get the expired access keys")
}

// Check status
// Work on list users command
