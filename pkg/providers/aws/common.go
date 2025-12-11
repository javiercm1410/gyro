package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/charmbracelet/log"
)

type UserWrapper struct {
	IamClient *iam.Client
}

type UserData interface {
}

type GetWrapperInputs struct {
	MaxUsers int32
	TimeZone string
	UserName string
	Client   UserWrapper
	Age      int
	Expired  bool
}

type RotateWrapperInputs struct {
	GetWrapperInputs
	DryRun           bool
	Notify           bool
	ExpireOnly       bool
	SkipConfirmation bool
	SkipCurrentUser  bool
}

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
