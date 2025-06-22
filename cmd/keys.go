package cmd

import (
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
		inputs, options := configureListCommand(cmd)

		userKeyData := iam.GetUserAccessKey(inputs)

		utils.DisplayData(options.Format, options.Path, options.Age, userKeyData)
	},
}

func init() {
	RootCmd.AddCommand(keysCmd)

	initializeBaseCommandFlags(keysCmd)
}
