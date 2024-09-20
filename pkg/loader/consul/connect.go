package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func (c *API) Connect(opts ...Option) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.client == nil {
		o := option{}
		o.apply(opts...)

		config := api.Config{
			HttpClient: o.HttpClient,
		}

		if err := c.connect(&config); err != nil {
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

	c.client = client.KV()

	return nil
}
