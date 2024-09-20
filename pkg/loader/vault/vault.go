package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path"
	"strings"
)

// Load specific key from vault.
func (c *API) Load(ctx context.Context, prefix, key string) ([]byte, error) {
	if c.client == nil {
		return nil, fmt.Errorf("vault client is not initialized")
	}

	// Get the key
	secret, err := c.client.KVv2(prefix).Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	// Decode the value
	v, err := c.codec(secret.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	return v, nil
}

// ListFolder lists all keys in a folder.
//   - prefix: like finops
//   - folder: with trailing slash.
func (c *API) ListFolder(ctx context.Context, prefix, folder string) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("vault client is not initialized")
	}

	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	folderList := path.Join(prefix, "metadata", folder)

	// Get the key
	secret, err := c.client.Logical().ListWithContext(ctx, folderList)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	keys, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to get keys from folder")
	}

	values := make([]string, 0, len(keys))
	for _, k := range keys {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("failed to key not a string, got %T", k)
		}

		// check has some space chars
		if strings.Contains(key, " ") {
			slog.Warn("key [" + key + "] contains space chars, skipping")

			continue
		}

		values = append(values, folder+key)
	}

	return values, nil
}

func (c *API) codec(v map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func (c *API) Travel(ctx context.Context, prefix, folder string, fn func(key string, value []byte) error) error {
	if folder != "" && !strings.HasSuffix(folder, "/") {
		// this is a key
		value, err := c.Load(ctx, prefix, folder)
		if err != nil {
			return err
		}

		if err := fn(folder, value); err != nil {
			return err
		}

		return nil
	}

	keys, err := c.ListFolder(ctx, prefix, folder)
	if err != nil {
		return err
	}

	for _, key := range keys {
		// check key has trailing slash
		if strings.HasSuffix(key, "/") {
			if err := c.Travel(ctx, prefix, key, fn); err != nil {
				return err
			}

			continue
		}

		value, err := c.Load(ctx, prefix, key)
		if err != nil {
			return err
		}

		if err := fn(key, value); err != nil {
			return err
		}
	}

	return nil
}
