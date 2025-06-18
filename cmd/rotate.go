package cmd

import (
	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/gyro/pkg/providers/aws"
	"github.com/spf13/cobra"
)

var rotateCmd = &cobra.Command{
	Use:     "rotate",
	Short:   "Rotate IAM Access Keys",
	Example: "gyro rotate [users|keys] [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("No arguments provided. Valid options are: 'users' and 'keys'")
		}
	},
}

func initRotateCommand(cmd *cobra.Command) iam.RotateWrapperInputs {
	options := configureRotateFlags(cmd)

	wrapper := iam.UserWrapper{
		IamClient: iam.DeclareConfig(),
	}

	return iam.RotateWrapperInputs{
		GetWrapperInputs: iam.GetWrapperInputs{
			MaxUsers: options.Quantity,
			TimeZone: options.TimeZone,
			Age:      options.Age,
			Expired:  options.Expired,
			UserName: options.User,
			Client:   wrapper,
		},
		DryRun:     options.DryRun,
		Notify:     options.Notify,
		ExpireOnly: options.ExpireOnly,
	}
}

var rotateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Rotate credentials for a specific IAM user",
	Run: func(cmd *cobra.Command, args []string) {
		inputs := initRotateCommand(cmd)

		iam.GetLoginProfiles(inputs.GetWrapperInputs)

	},
}

var rotateKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Rotate credentials for a specific IAM key",
	Run: func(cmd *cobra.Command, args []string) {
		inputs := initRotateCommand(cmd)

		iam.GetUserAccessKey(inputs.GetWrapperInputs)

	},
}

func init() {
	RootCmd.AddCommand(rotateCmd)
	rotateCmd.AddCommand(rotateUserCmd)
	rotateCmd.AddCommand(rotateKeyCmd)

	initializeListCommandFlags(rotateCmd)
}
