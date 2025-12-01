package cmd

import (
	"fmt"

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
			Expired:  options.Expired,
			UserName: options.User,
			Client:   wrapper,
		},
		Notify:     options.Notify,
		ExpireOnly: options.ExpireOnly,
	}, baseOptions
}

func askForConfirmation() bool {
	fmt.Println("Confirmation? (y/n)")
	var response string
	fmt.Scanln(&response)
	return response == "y"
}

var rotateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Rotate credentials for a specific IAM user",
	Run: func(cmd *cobra.Command, args []string) {
		inputs, baseOptions := initRotateCommand(cmd)

		userPasswordData := iam.GetLoginProfiles(inputs.GetWrapperInputs)

		utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, userPasswordData)

		if !askForConfirmation() {
			fmt.Println("Operation aborted.")
			return
		}
		fmt.Println("Operation confirmed.")

		iam.UserWrapper.RotateLoginProfiles(inputs.GetWrapperInputs.Client, userPasswordData)
	},
}

var rotateKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Rotate credentials for a specific IAM key",
	Run: func(cmd *cobra.Command, args []string) {
		inputs, baseOptions := initRotateCommand(cmd)

		userKeyData := iam.GetUserAccessKey(inputs.GetWrapperInputs)

		utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, userKeyData)

		if !askForConfirmation() {
			fmt.Println("Operation aborted.")
			return
		}
		fmt.Println("Operation confirmed.")

		// iam.RotateAccessKeys(inputs.GetWrapperInputs.Client, userKeyData)
	},
}

func init() {
	RootCmd.AddCommand(rotateCmd)
	rotateCmd.AddCommand(rotateUserCmd)
	rotateCmd.AddCommand(rotateKeyCmd)

	initializeBaseCommandFlags(rotateCmd)
}
