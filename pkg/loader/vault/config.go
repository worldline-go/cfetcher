package vault

import (
	"sync"

	"github.com/hashicorp/vault/api"
)

var (
	EnvVaultRoleID          = "VAULT_ROLE_ID"
	EnvVaultApproleBasePath = "VAULT_APPROLE_BASE_PATH"
	EnvVaultRoleSecret      = "VAULT_ROLE_SECRET"
)

type API struct {
	AppRoleBasePath string `cfg:"app_role_base_path"`
	client          *api.Client
	m               sync.Mutex
}
