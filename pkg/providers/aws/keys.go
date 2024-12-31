package iam

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/charmbracelet/log"
)

type AccessKeyData struct {
	Id              *string
	CreateDate      time.Time
	KeyStatus       types.StatusType
	LastUsedTime    time.Time
	LastUsedService string
}

type UserAccessKeyData struct {
	UserName string
	Keys     []AccessKeyData
}

// ListAccessKeys fetches access keys for a specific user.
func (wrapper UserWrapper) ListAccessKeys(userName, timeZone string, expired bool, stale int) (UserAccessKeyData, error) {
	var keys []AccessKeyData
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Errorf("Couldn't list Load Time Zone %s. Error: %v", timeZone, err)
		return UserAccessKeyData{}, err
	}

	input := &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}

	result, err := wrapper.IamClient.ListAccessKeys(context.TODO(), input)
	if err != nil {
		log.Errorf("Couldn't list access keys for user %s. Error: %v", userName, err)
		return UserAccessKeyData{}, err
	}

	for _, key := range result.AccessKeyMetadata {
		keyData, addToList, err := wrapper.getAccessKeyDetails(key, loc, expired, stale)
		if err != nil {
			log.Errorf("Couldn't fetch access key details for user %s. Error: %v", userName, err)
			continue
		}
		if addToList {
			keys = append(keys, keyData)
		}
	}

	return UserAccessKeyData{
		UserName: userName,
		Keys:     keys,
	}, nil
}

// getAccessKeyDetails fetches and processes details for a single access key.
func (wrapper UserWrapper) getAccessKeyDetails(key types.AccessKeyMetadata, loc *time.Location, expired bool, stale int) (AccessKeyData, bool, error) {
	accessKeyInput := &iam.GetAccessKeyLastUsedInput{
		AccessKeyId: aws.String(*key.AccessKeyId),
	}

	lastUsed, err := wrapper.IamClient.GetAccessKeyLastUsed(context.TODO(), accessKeyInput)
	if err != nil {
		return AccessKeyData{}, false, errors.New("couldn't get last used data for access key")
	}

	keyData := AccessKeyData{
		Id:         key.AccessKeyId,
		CreateDate: key.CreateDate.In(loc),
		KeyStatus:  key.Status,
	}

	if lastUsed.AccessKeyLastUsed.LastUsedDate != nil {
		keyData.LastUsedTime = lastUsed.AccessKeyLastUsed.LastUsedDate.In(loc)
		keyData.LastUsedService = *lastUsed.AccessKeyLastUsed.ServiceName
	} else {
		keyData.LastUsedService = "n/a"
	}

	if expired {
		if time.Since(keyData.CreateDate).Hours() > float64(stale*24) {
			return keyData, true, nil
		} else {
			return keyData, false, nil
		}
	}

	return keyData, true, nil
}

func GetUserAccessKey(input GetWrapperInputs) ([]UserData, error) {
	var usersData []types.User
	var err error

	if input.UserName != "" {
		usersData = []types.User{{UserName: aws.String(input.UserName)}}
	} else {
		usersData, err = input.Client.ListUsers(input.MaxUsers)
		if err != nil {
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

	//This work if commented, we can leave like that but let's see why
	// if input.UserName != "" {
	sort.Slice(userKeyData, func(i, j int) bool {
		return userKeyData[j].(UserAccessKeyData).UserName > userKeyData[i].(UserAccessKeyData).UserName
	})
	// }

	return userKeyData, nil
}
