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
			log.Errorf("Couldn't list users: %v", err)
			return nil, err
		}
	}

	var (
		userKeyData []UserData
		mu          sync.Mutex
		wg          sync.WaitGroup
		errors      []error
	)

	wg.Add(len(usersData))

	for _, user := range usersData {
		go func(user types.User) {
			defer wg.Done()
			keyData, err := input.Client.ListAccessKeys(*user.UserName, input.TimeZone, input.Expired, input.Age)
			if err != nil {
				log.Errorf("Couldn't list access keys for user %s: %v", *user.UserName, err)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			mu.Lock()
			if input.Expired {
				if keyData.Keys != nil {
					userKeyData = append(userKeyData, keyData)
				}
			} else {
				userKeyData = append(userKeyData, keyData)
			}
			mu.Unlock()
		}(user)
	}

	wg.Wait()
	if len(errors) > 0 {
		return nil, errors[0] // Returning the first error as an example
	}

	if input.UserName != "" {
		sort.Slice(userKeyData, func(i, j int) bool {
			return userKeyData[j].(UserAccessKeyData).UserName > userKeyData[i].(UserAccessKeyData).UserName
		})
	}

	return userKeyData, nil
}

// func GetUserConsoleAccess(input GetUserAccessKeyInputs) ([]UserData, error) {
// 	var usersData []types.User
// 	var err error

// 	if input.UserName != "" {
// 		usersData = []types.User{{UserName: aws.String(input.UserName)}}
// 	} else {
// 		usersData, err = input.Client.ListUsers(input.MaxUsers)
// 		if err != nil {
// 			log.Errorf("Couldn't list users: %v", err)
// 			return nil, err
// 		}
// 	}

// 	return userKeyData, nil
// }
