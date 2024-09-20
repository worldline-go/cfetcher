package vault

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

func (c *API) Connect(opts ...Option) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.client == nil {
		o := option{}
		o.apply(opts...)

		configDef := api.DefaultConfig()

		config := &api.Config{
			Address: configDef.Address,
		}

		if o.HttpClient != nil {
			config.HttpClient = o.HttpClient
		}

		if err := c.connect(config); err != nil {
			return err
		}

		if err := c.login(context.Background()); err != nil {
			return err
		}
	}

	return nil
}

func (c *API) connect(config *api.Config) error {
	// Get a new client
	client, err := api.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	c.client = client

	return nil
}

func (c *API) login(ctx context.Context) error {
	// A combination of a Role ID and Secret ID is required to log in to Vault with an AppRole.
	// First, let's get the role ID given to us by our Vault administrator.
	roleID := os.Getenv(EnvVaultRoleID)
	if roleID == "" {
		return fmt.Errorf("no role ID was provided in VAULT_ROLE_ID env var")
	}

	// check default path
	appRoleBasePath := c.AppRoleBasePath
	if appRoleBasePath == "" {
		appRoleBasePath = os.Getenv(EnvVaultApproleBasePath)
	}

	if appRoleBasePath == "" {
		appRoleBasePath = "auth/approle/login"
	}

	secret, err := c.client.Logical().WriteWithContext(ctx, appRoleBasePath, map[string]interface{}{
		"role_id":   roleID,
		"secret_id": os.Getenv(EnvVaultRoleSecret),
	})
	if err != nil {
		return fmt.Errorf("failed to login to vault: %w", err)
	}

	// Set the token
	c.client.SetToken(secret.Auth.ClientToken)

	return nil
}
