package consul

import (
	"sync"

	"github.com/hashicorp/consul/api"
)

type API struct {
	client *api.KV
	m      sync.Mutex
}
