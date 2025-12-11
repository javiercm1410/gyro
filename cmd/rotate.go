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
		} else {
			log.Fatalf("Invalid argument '%s'. Valid options are: 'users' and 'keys'", args[0])
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
		Notify:           options.Notify,
		ExpireOnly:       options.ExpireOnly,
		SkipConfirmation: options.SkipConfirmation,
		SkipCurrentUser:  options.SkipCurrentUser,
	}, baseOptions
}

func askForConfirmation() bool {
	fmt.Println("Confirmation? (y/n)")
	var response string
	fmt.Scanln(&response)
	return response == "y"
}

var rotateUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Rotate credentials for a specific IAM user",
	Run: func(cmd *cobra.Command, args []string) {
		inputs, baseOptions := initRotateCommand(cmd)

		userPasswordData := iam.GetLoginProfiles(inputs.GetWrapperInputs)

		// if inputs.SkipCurrentUser {
		// 	iam.RemoveCurrentUser(userPasswordData)
		// }

		utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, userPasswordData)

		if len(userPasswordData) > 0 {
			if !inputs.SkipConfirmation && !askForConfirmation() {
				fmt.Println("Operation aborted.")
				return
			}
			fmt.Println("Operation confirmed.")

			userResults := iam.UserWrapper.RotateLoginProfiles(inputs.GetWrapperInputs.Client, userPasswordData)
			utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, userResults)
		}
	},
}

var rotateKeyCmd = &cobra.Command{
	Use:     "key",
	Aliases: []string{"keys"},
	Short:   "Rotate credentials for a specific IAM key",
	Run: func(cmd *cobra.Command, args []string) {
		inputs, baseOptions := initRotateCommand(cmd)

		userKeyData := iam.GetUserAccessKey(inputs.GetWrapperInputs)

		// if inputs.SkipCurrentUser {
		// 	iam.RemoveCurrentUser(userKeyData)
		// }

		utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, userKeyData)

		if len(userKeyData) > 0 {
			if !inputs.SkipConfirmation && !askForConfirmation() {
				fmt.Println("Operation aborted.")
				return
			}
			fmt.Println("Operation confirmed.")

			keyResults := iam.UserWrapper.RotateAccessKeys(inputs.GetWrapperInputs.Client, userKeyData, inputs.SkipConfirmation)
			utils.DisplayData(baseOptions.Format, baseOptions.Path, baseOptions.Age, keyResults)
		}
	},
}

func init() {
	RootCmd.AddCommand(rotateCmd)
	rotateCmd.AddCommand(rotateUserCmd)
	rotateCmd.AddCommand(rotateKeyCmd)

	initializeBaseCommandFlags(rotateCmd)
}
