package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/gyro/pkg/providers/aws"
	"github.com/javiercm1410/gyro/pkg/utils"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:     "get-keys",
	Short:   "Get IAM Access Keys",
	Aliases: []string{"keys", "k"},
	Example: "gyro keys",
	Run: func(cmd *cobra.Command, args []string) {
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

		userKeyData, err := iam.GetUserAccessKey(inputs)
		if err != nil {
			log.Error("Failed to get users", "error", err)
			os.Exit(1)
		}

		utils.DisplayData(options.Format, options.Path, options.Age, userKeyData)
	},
}

func init() {
	RootCmd.AddCommand(keysCmd)

	initializeListCommandFlags(keysCmd)
}
