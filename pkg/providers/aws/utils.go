package iam

import (
	"github.com/javiercm1410/rotator/pkg/utils"

	"github.com/charmbracelet/log"
)

type GetUserAccessKeyInputs struct {
	MaxUsers   int32
	TimeZone   string
	OutputType string
	Path       string
	Client     UserWrapper
}

func GetUserAccessKey(input GetUserAccessKeyInputs) error {
	usersData, err := input.Client.ListUsers(input.MaxUsers)
	if err != nil {
		log.Errorf("Couldn't list users. Here's why: %v\n", err)
		return err
	}

	var userKeyData []UserAccessKeyData
	for _, user := range usersData {
		keyData, err := input.Client.ListAccessKeys(*user.UserName, "America/Santo_Domingo")
		if err != nil {
			log.Errorf("Couldn't list users. Here's why: %v\n", err)
			return err
		}
		userKeyData = append(userKeyData, keyData)
	}

	utils.DisplayData(input.OutputType, input.OutputType, userKeyData)
	return nil
}
