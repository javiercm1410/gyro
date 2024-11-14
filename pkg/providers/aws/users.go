package iam

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
	"github.com/charmbracelet/log"
)

type UserWrapper struct {
	IamClient *iam.Client
}

type AccessKeyData struct {
	Id              *string
	CreateDate      *time.Time
	KeyStatus       types.StatusType
	LastUsedTime    time.Time
	LastUsedService string
}

type UserAccessKeyData struct {
	UserName string
	Keys     []AccessKeyData
}

// type UserData struct {
// 	UserName             string
// 	LastConsoleLoginDate string
// 	AccessKey            string
// 	Active               string
// 	LastCredentialUsed   string
// }

func DeclareConfig() *iam.Client {
	// var IamClient *iam.Client
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Warn("Couldn't load default configuration. Have you set up your AWS account?")
		log.Error(err)
		return nil
	}
	IamClient := iam.NewFromConfig(sdkConfig)
	return IamClient
}

func (wrapper UserWrapper) ListUsers(maxUsers int32) ([]types.User, error) {
	var users []types.User

	input := &iam.ListUsersInput{
		MaxItems: aws.Int32(maxUsers),
	}

	for {
		result, err := wrapper.IamClient.ListUsers(context.TODO(), input)
		if err != nil {
			log.Errorf("Couldn't list users. Here's why: %v\n", err)
			return nil, err
		}

		users = append(users, result.Users...)

		if result.Marker != nil {
			input.Marker = result.Marker
		}
		if !result.IsTruncated || maxUsers != 50 {
			break
		}
	}

	return users, nil
}

func (wrapper UserWrapper) GetUser(userName string) (*types.User, error) {
	var user *types.User
	result, err := wrapper.IamClient.GetUser(context.TODO(), &iam.GetUserInput{
		UserName: aws.String(userName),
	})
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NoSuchEntityException:
				log.Error("User %v does not exist.\n", userName)
				err = nil
			default:
				log.Error("Couldn't get user %v. Here's why: %v\n", userName, err)
			}
		}
	} else {
		user = result.User
	}
	return user, err
}

func (wrapper UserWrapper) ListAccessKeys(userName, timeZone string) (UserAccessKeyData, error) {
	var keys []AccessKeyData

	input := &iam.ListAccessKeysInput{
		UserName: aws.String(userName),
	}

	result, err := wrapper.IamClient.ListAccessKeys(context.TODO(), input)
	if err != nil {
		log.Error("Couldn't list access keys for user %v. Here's why: %v\n", userName, err)
	}

	if len(result.AccessKeyMetadata) > 0 {
		for _, key := range result.AccessKeyMetadata {
			accessKeyInput := &iam.GetAccessKeyLastUsedInput{
				AccessKeyId: aws.String(*key.AccessKeyId),
			}
			lastUsed, err := wrapper.IamClient.GetAccessKeyLastUsed(context.TODO(), accessKeyInput)
			if err != nil {
				log.Error("Couldn't get access keys last login for user %v. Here's why: %v\n", userName, err)
			}

			if lastUsed.AccessKeyLastUsed.LastUsedDate == nil {
				keys = append(keys, AccessKeyData{
					Id:              key.AccessKeyId,
					LastUsedService: "n/a",
				})
				continue
			}

			loc, _ := time.LoadLocation(timeZone)
			keys = append(keys, AccessKeyData{
				Id:              key.AccessKeyId,
				CreateDate:      key.CreateDate,
				KeyStatus:       key.Status,
				LastUsedTime:    lastUsed.AccessKeyLastUsed.LastUsedDate.In(loc),
				LastUsedService: *lastUsed.AccessKeyLastUsed.ServiceName,
			})
		}
	}

	userKeyData := UserAccessKeyData{
		UserName: userName,
		Keys:     keys,
	}

	return userKeyData, err
}
