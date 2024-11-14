package op

import (
	"context"
	"fmt"

	parser "github.com/javiercm1410/rotator/pkg/utils"

	"github.com/1password/onepassword-sdk-go"
)

type OnePasswordClient struct {
	Ctx    context.Context
	Client *onepassword.Client
}

type OnepasswordOptions struct {
	Vault string
	Items []string
	Path  []string
}

func NewClient(token string) (*OnePasswordClient, error) {
	if token == "" {
		return nil, fmt.Errorf("unauthorized: OnePassword token not present")
	}

	ctx := context.TODO()
	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo("My 1Password Integration", "v1.0.0"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OnePassword client: %v", err)
	}

	return &OnePasswordClient{
		Ctx:    ctx,
		Client: client,
	}, nil
}

func (client *OnePasswordClient) GenerateEnvFile(options OnepasswordOptions) error {
	// var vaultItem onepassword.Item
	for _, item := range options.Items {
		vaultItem, err := client.Client.Items.Get(client.Ctx, options.Vault, item)
		if err != nil {
			return fmt.Errorf("failed to get vault item: %v", err)
		}

		envData := make(parser.EnvVarObject)
		for _, field := range vaultItem.Fields {
			envData[field.Title] = field.Value
		}

		for _, path := range options.Path {
			if err := parser.GenerateEnvFile(envData, path); err != nil {
				return fmt.Errorf("failed to generate env file at %s: %v", path, err)
			}
		}
	}

	return nil
}
