package iam

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/charmbracelet/log"
)

type UserWrapper struct {
	IamClient *iam.Client
}

type UserData interface {
	isUserData()
}

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

// Implement the interface
func (u UserAccessKeyData) isUserData() {}

// DeclareConfig initializes the IAM client using the default AWS configuration.
func DeclareConfig() *iam.Client {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Warn("Couldn't load default configuration. Ensure AWS account setup.")
		log.Error(err)
		return nil
	}
	return iam.NewFromConfig(sdkConfig)
}

// ListUsers fetches a list of IAM users up to the specified maximum.
func (wrapper UserWrapper) ListUsers(maxUsers int32) ([]types.User, error) {
	var users []types.User

	input := &iam.ListUsersInput{
		MaxItems: aws.Int32(maxUsers),
	}

	for {
		result, err := wrapper.IamClient.ListUsers(context.TODO(), input)
		if err != nil {
			log.Errorf("Couldn't list users. Error: %v", err)
			return nil, err
		}

		users = append(users, result.Users...)

		if result.Marker != nil {
			input.Marker = result.Marker
		}
		// maxUsers != because 50 is the default value, so it should get all
		if !result.IsTruncated || maxUsers != 50 {
			break
		}
	}

	return users, nil
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
		return AccessKeyData{}, false, errors.New("Couldn't get last used data for access key")
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
