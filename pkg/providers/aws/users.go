package iam

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/charmbracelet/log"
)

type UserLoginData struct {
	UserName     string
	LastUsedTime time.Time
	LoginProfile *types.LoginProfile
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

// GetLoginProfile fetches login profile info for a specific user.
func (wrapper UserWrapper) GetLoginProfile(user types.User, expired, debug bool, stale int) (UserLoginData, error) {
	input := &iam.GetLoginProfileInput{
		UserName: user.UserName,
	}

	userLoginProfile := UserLoginData{
		UserName: *user.UserName,
	}

	if user.PasswordLastUsed != nil && !user.PasswordLastUsed.IsZero() {
		userLoginProfile.LastUsedTime = *user.PasswordLastUsed
	}

	result, err := wrapper.IamClient.GetLoginProfile(context.TODO(), input)
	if err != nil {
		if debug {
			log.Infof("Couldn't list login profile for user %s. Error: %v", *user.UserName, err)
		}
		return UserLoginData{}, err

	}

	userLoginProfile.LoginProfile = result.LoginProfile

	if expired {
		if time.Since(*result.LoginProfile.CreateDate).Hours() > float64(stale*24) {
			return userLoginProfile, nil
		} else {
			return UserLoginData{}, nil
		}
	}

	return userLoginProfile, nil
}

func GetLoginProfiles(input GetWrapperInputs) []UserData {
	var usersData []types.User
	var err error

	if input.UserName != "" {
		inputGetUser := &iam.GetUserInput{
			UserName: &input.UserName,
		}

		selectedUser, err := input.Client.IamClient.GetUser(context.TODO(), inputGetUser)
		if err != nil {
			log.Fatalf("Failed to get users: %v", err)
		}

		if !selectedUser.User.PasswordLastUsed.IsZero() {
			usersData = []types.User{{
				UserName:         aws.String(input.UserName),
				PasswordLastUsed: selectedUser.User.PasswordLastUsed,
			}}
		} else {
			log.Error("User haven't accessed AWS console yet")
		}

	} else {
		usersData, err = input.Client.ListUsers(input.MaxUsers)
		if err != nil {
			log.Fatalf("Failed to get users: %v", err)
		}
	}

	var (
		userLoginProfiles []UserData
		mu                sync.Mutex
		wg                sync.WaitGroup
		// errors            []error
	)

	wg.Add(len(usersData))

	for _, user := range usersData {
		go func(user types.User) {
			defer wg.Done()
			userLogin, err := input.Client.GetLoginProfile(user, input.Expired, false, input.Age)
			if err != nil {
				// mu.Lock()

				// We do this on purpose to avoid logs for profiles without login: 404 error
				// errors = append(errors, err)
				// mu.Unlock()
				return
			}

			mu.Lock()
			userLoginProfiles = append(userLoginProfiles, userLogin)

			mu.Unlock()
		}(user)
	}

	wg.Wait()

	sort.Slice(userLoginProfiles, func(i, j int) bool {
		return userLoginProfiles[j].(UserLoginData).UserName > userLoginProfiles[i].(UserLoginData).UserName
	})

	// if len(errors) > 0 {
	// 	// CHange error message
	// 	return userLoginProfiles, errors[0] // Returning the first error as an example
	// }

	return userLoginProfiles
}
