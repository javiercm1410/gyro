package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/gyro/pkg/providers/aws"
	"github.com/javiercm1410/gyro/pkg/utils"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:     "get-users",
	Short:   "Get IAM users",
	Aliases: []string{"users", "u"},
	Example: "gyro users",
	Run: func(cmd *cobra.Command, args []string) {
		options := configureListFlags(cmd)

		wrapper := iam.UserWrapper{
			IamClient: iam.DeclareConfig(),
		}

		inputs := iam.GetUserAccessKeyInputs{
			MaxUsers: options.Quantity,
			TimeZone: options.TimeZone,
			Age:      options.Age,
			Expired:  options.Expired,
			UserName: options.User,
			Client:   wrapper,
		}

		userPasswordData, err := wrapper.ListUsers(inputs.MaxUsers)
		if err != nil {
			log.Error("Failed to get users", "error", err)
			os.Exit(1)
		}

		utils.DisplayData(options.Format, options.Path, options.Age, []iam.UserData{userPasswordData})
	},
}

func init() {
	RootCmd.AddCommand(usersCmd)

	initializeListCommandFlags(usersCmd)
}
