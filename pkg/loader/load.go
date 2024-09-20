package loader

import (
	"github.com/worldline-go/cfetcher/pkg/loader/consul"
	"github.com/worldline-go/cfetcher/pkg/loader/vault"
)

type Loaders struct {
	Vault  vault.API  `cfg:"vault"`
	Consul consul.API `cfg:"consul"`
}
