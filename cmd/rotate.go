package cmd

import (
	"github.com/charmbracelet/log"
	iam "github.com/javiercm1410/gyro/pkg/providers/aws"
	"github.com/javiercm1410/gyro/pkg/utils"
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

func initRotateCommand(cmd *cobra.Command) (iam.RotateWrapperInputs, BaseCommandOptions) {
	options, baseOptions := configureRotateCommand(cmd)

	wrapper := iam.UserWrapper{
		IamClient: iam.DeclareConfig(),
	}

	return iam.RotateWrapperInputs{
		GetWrapperInputs: iam.GetWrapperInputs{
			MaxUsers: options.Quantity,
			TimeZone: options.TimeZone,
			Age:      options.Age,
			Expired:  true,
			UserName: options.User,
			Client:   wrapper,
		},
		DryRun:     options.DryRun,
		Notify:     options.Notify,
		ExpireOnly: options.ExpireOnly,
	}, baseOptions
}

var rotateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Rotate credentials for a specific IAM user",
	Run: func(cmd *cobra.Command, args []string) {
		inputs, baseOptions := initRotateCommand(cmd)

		userPasswordData := iam.GetLoginProfiles(inputs.GetWrapperInputs)

		utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, userPasswordData)

	},
}

var rotateKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Rotate credentials for a specific IAM key",
	Run: func(cmd *cobra.Command, args []string) {
		// inputs, _ := initRotateCommand(cmd)

		// iam.GetUserAccessKey(inputs)

	},
}

func init() {
	RootCmd.AddCommand(rotateCmd)
	rotateCmd.AddCommand(rotateUserCmd)
	rotateCmd.AddCommand(rotateKeyCmd)

	initializeBaseCommandFlags(rotateCmd)
}
