package iam

import (
	"context"
	"errors"
	"fmt"
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
	MatchesCriteria bool
	IsExpired       bool
}

type AccessKeyRotationResult struct {
	UserName        string
	AccessKeyId     string
	SecretAccessKey string
}

type UserAccessKeyData struct {
	UserName string
	Keys     []AccessKeyData
}

// ListAccessKeys fetches access keys for a specific user.
func (wrapper UserWrapper) ListAccessKeys(userName, timeZone string, expired bool, stale int) (UserAccessKeyData, error) {
	var keys []AccessKeyData
	hasMatch := false

	// this should be on presentation only
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

	if len(result.AccessKeyMetadata) == 0 {
		return UserAccessKeyData{}, nil
	}

	for _, key := range result.AccessKeyMetadata {
		keyData, err := wrapper.getAccessKeyDetails(key, loc, expired, stale)
		if err != nil {
			log.Errorf("Couldn't fetch access key details for user %s. Error: %v", userName, err)
			continue
		}
		keys = append(keys, keyData)
		if keyData.MatchesCriteria {
			hasMatch = true
		}
	}

	if expired && !hasMatch {
		return UserAccessKeyData{}, nil
	}

	return UserAccessKeyData{
		UserName: userName,
		Keys:     keys,
	}, nil
}

// getAccessKeyDetails fetches and processes details for a single access key.
func (wrapper UserWrapper) getAccessKeyDetails(key types.AccessKeyMetadata, loc *time.Location, expired bool, stale int) (AccessKeyData, error) {
	accessKeyInput := &iam.GetAccessKeyLastUsedInput{
		AccessKeyId: aws.String(*key.AccessKeyId),
	}

	lastUsed, err := wrapper.IamClient.GetAccessKeyLastUsed(context.TODO(), accessKeyInput)
	if err != nil {
		return AccessKeyData{}, errors.New("couldn't get last used data for access key")
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

	keyData.MatchesCriteria = true
	if time.Since(keyData.CreateDate).Hours() > float64(stale*24) {
		keyData.IsExpired = true
	}

	if expired {
		if !keyData.IsExpired {
			keyData.MatchesCriteria = false
		}
	}

	return keyData, nil
}

// RotateAccessKeys rotates the access keys for the provided users.
func (wrapper UserWrapper) RotateAccessKeys(keys []UserData, skipConfirmation bool) []UserData {
	var results []UserData
	for _, keyData := range keys {
		user, ok := keyData.(UserAccessKeyData)
		if !ok {
			log.Warnf("Skipping invalid user data type: %T", keyData)
			continue
		}

		// Check for expired active keys and prompt for deactivation
		for _, key := range user.Keys {
			if key.IsExpired && key.KeyStatus == types.StatusTypeActive {
				shouldDeactivate := skipConfirmation
				if !skipConfirmation {
					fmt.Printf("User %s has an expired active access key (%s). Do you want to deactivate it? (y/n): ", user.UserName, *key.Id)
					var response string
					fmt.Scanln(&response)
					if response == "y" {
						shouldDeactivate = true
					}
				}

				if shouldDeactivate {
					updateInput := &iam.UpdateAccessKeyInput{
						UserName:    aws.String(user.UserName),
						AccessKeyId: key.Id,
						Status:      types.StatusTypeInactive,
					}
					_, err := wrapper.IamClient.UpdateAccessKey(context.TODO(), updateInput)
					if err != nil {
						log.Errorf("Failed to deactivate access key %s: %v", *key.Id, err)
					} else {
						log.Infof("Successfully deactivated access key %s", *key.Id)
					}
				}
			}
		}

		if len(user.Keys) >= 2 {
			// Find oldest key
			oldestKey := user.Keys[0]
			for _, k := range user.Keys {
				if k.CreateDate.Before(oldestKey.CreateDate) {
					oldestKey = k
				}
			}

			if !skipConfirmation {
				fmt.Printf("User %s has 2 access keys. Do you want to delete the oldest key (%s created on %s)? (y/n): ", user.UserName, *oldestKey.Id, oldestKey.CreateDate)
				var response string
				fmt.Scanln(&response)
				if response != "y" {
					log.Warnf("Skipping rotation for user %s as they have 2 keys", user.UserName)
					continue
				}
			}

			// Delete key
			deleteInput := &iam.DeleteAccessKeyInput{
				UserName:    aws.String(user.UserName),
				AccessKeyId: oldestKey.Id,
			}
			_, err := wrapper.IamClient.DeleteAccessKey(context.TODO(), deleteInput)
			if err != nil {
				log.Errorf("Failed to delete access key %s for user %s: %v", *oldestKey.Id, user.UserName, err)
				continue
			}
			log.Infof("Successfully deleted access key %s for user %s", *oldestKey.Id, user.UserName)
		}

		// Create new key
		createInput := &iam.CreateAccessKeyInput{
			UserName: aws.String(user.UserName),
		}
		createOutput, err := wrapper.IamClient.CreateAccessKey(context.TODO(), createInput)
		if err != nil {
			log.Errorf("Failed to create access key for user %s: %v", user.UserName, err)
			continue
		}

		log.Infof("Successfully rotated access key for user: %s", user.UserName)
		log.Infof("Access Key ID: %s", *createOutput.AccessKey.AccessKeyId)
		log.Infof("Secret Access Key: %s", *createOutput.AccessKey.SecretAccessKey)

		results = append(results, AccessKeyRotationResult{
			UserName:        user.UserName,
			AccessKeyId:     *createOutput.AccessKey.AccessKeyId,
			SecretAccessKey: *createOutput.AccessKey.SecretAccessKey,
		})
	}

	return results
}

func GetUserAccessKey(input GetWrapperInputs) []UserData {
	var usersData []types.User
	var err error

	if input.UserName != "" {
		usersData = []types.User{{UserName: aws.String(input.UserName)}}
	} else {
		usersData, err = input.Client.ListUsers(input.MaxUsers)
		if err != nil {
			log.Fatalf("Failed to get users: %v", err)
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
			if keyData.Keys != nil {
				if input.Expired {
					userKeyData = append(userKeyData, keyData)
				} else {
					userKeyData = append(userKeyData, keyData)
				}
			}
			mu.Unlock()
		}(user)
	}

	wg.Wait()
	if len(errors) > 0 {
		for _, err := range errors {
			log.Errorf("Failed to get user: %v", err)
		}
	}

	//This work if commented, we can leave like that but let's see why
	// if input.UserName != "" {
	sort.Slice(userKeyData, func(i, j int) bool {
		return userKeyData[j].(UserAccessKeyData).UserName > userKeyData[i].(UserAccessKeyData).UserName
	})
	// }

	return userKeyData
}
