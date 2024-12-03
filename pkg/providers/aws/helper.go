package iam

import (
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/charmbracelet/log"
)

type GetUserAccessKeyInputs struct {
	MaxUsers int32
	TimeZone string
	UserName string
	Client   UserWrapper
	Age      int
	Expired  bool
}

func GetUserAccessKey(input GetUserAccessKeyInputs) ([]UserData, error) {
	var usersData []types.User
	var err error
	if input.UserName != "" {
		usersData = []types.User{{UserName: aws.String(input.UserName)}}
	} else {
		usersData, err = input.Client.ListUsers(input.MaxUsers)
		if err != nil {
			log.Errorf("Couldn't list users. Here's why: %v\n", err)
			return nil, err
		}
	}

	var userKeyData []UserData

	wg := sync.WaitGroup{}
	wg.Add(len(usersData))

	for _, user := range usersData {
		go func() {
			defer wg.Done()
			keyData, err := input.Client.ListAccessKeys(*user.UserName, input.TimeZone, input.Expired, input.Age)
			if err != nil {
				log.Errorf("Couldn't list users. Here's why: %v\n", err)
				return
			}
			if input.Expired {
				if keyData.Keys != nil {
					userKeyData = append(userKeyData, keyData)
				}
			} else {
				userKeyData = append(userKeyData, keyData)
			}
		}()
	}

	wg.Wait()

	// Sort the slice by age in descending order
	sort.Slice(userKeyData, func(i, j int) bool {
		return userKeyData[j].(UserAccessKeyData).UserName > userKeyData[i].(UserAccessKeyData).UserName
	})

	return userKeyData, nil
}
