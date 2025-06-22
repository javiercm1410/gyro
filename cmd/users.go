package cmd

import (
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
		inputs, options := configureListCommand(cmd)

		userPasswordData := iam.GetLoginProfiles(inputs)

		utils.DisplayData(options.Format, options.Path, inputs.Age, userPasswordData)
	},
}

func init() {
	RootCmd.AddCommand(usersCmd)

	initializeBaseCommandFlags(usersCmd)
}
