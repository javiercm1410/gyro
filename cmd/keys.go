package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/rotator/pkg/providers/aws"
	"github.com/javiercm1410/rotator/pkg/utils"
	"github.com/spf13/cobra"
)

type ListCommandOptions struct {
	Quantity int32
	Path     string
	User     string
	TimeZone string
	Format   string
	Age      int
	Expired  bool
}

var listCmd = &cobra.Command{
	Use:     "get-keys",
	Short:   "Get IAM Access Keys",
	Aliases: []string{"keys", "k"},
	Example: "rotate keys",
	Run: func(cmd *cobra.Command, args []string) {
		options := configureListFlags(cmd)
		var wrapper iam.UserWrapper
		wrapper.IamClient = iam.DeclareConfig()

		inputs := iam.GetUserAccessKeyInputs{
			MaxUsers: options.Quantity,
			TimeZone: options.TimeZone,
			Age:      options.Age,
			Expired:  options.Expired,
			UserName: options.User,
			Client:   wrapper,
		}

		userKeyData, err := iam.GetUserAccessKey(inputs)
		if err != nil {
			log.Errorf("Failed to get users. Here's why: %v\n", err)
			os.Exit(1)
		}

		utils.DisplayData(options.Format, options.Path, options.Age, userKeyData)
	},
}

func configureListFlags(cmd *cobra.Command) ListCommandOptions {
	quantity, _ := cmd.Flags().GetInt32("quantity")
	timeZone, _ := cmd.Flags().GetString("timezone")
	format, _ := cmd.Flags().GetString("format")
	userName, _ := cmd.Flags().GetString("username")
	path, _ := cmd.Flags().GetString("output-file")
	age, _ := cmd.Flags().GetInt("age")
	expired, _ := cmd.Flags().GetBool("expired")

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

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringP("timezone", "t", "America/Santo_Domingo", "Timezone for displaying dates")
	listCmd.PersistentFlags().StringP("format", "f", "json", "Output format (json, table, text)")
	listCmd.PersistentFlags().StringP("output-file", "o", "./output.json", "Save results to file")
	listCmd.PersistentFlags().StringP("username", "u", "", "Filter by specific IAM username")
	listCmd.PersistentFlags().IntP("age", "a", 90, "Consider keys stale after N days")
	listCmd.PersistentFlags().BoolP("expired-only", "x", false, "Show only expired keys")
}
