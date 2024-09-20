package consul

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
)

func (c *API) Set(ctx context.Context, folderPath string) error {
	if c.client == nil {
		return fmt.Errorf("consul client not initialized")
	}

	// walk in folder
	return walkDir(folderPath, func(path string) error {
		// read file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("path [%s] cannot read; %w", path, err)
		}

		path = strings.TrimPrefix(path, folderPath)
		path = strings.TrimSuffix(path, ".yaml")

		slog.Info("consul config set [" + path + "]")

		// set key
		_, err = c.client.Put(&api.KVPair{
			Key:   path,
			Value: data,
		}, nil)
		if err != nil {
			return fmt.Errorf("path [%s] cannot set; %w", path, err)
		}

		return nil
	})
}

func walkDir(dir string, fn func(path string) error) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err := walkDir(path, fn); err != nil {
				return err
			}
		} else {
			if err := fn(path); err != nil {
				return err
			}
		}
	}

	return nil
}

// Load specific key from consul.
func (c *API) Load(ctx context.Context, prefix, key string) ([]byte, error) {
	if c.client == nil {
		return nil, fmt.Errorf("consul client not initialized")
	}

	pair, _, err := c.client.Get(path.Join(prefix, key), nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, fmt.Errorf("key [%s] not found", key)
	}

	return pair.Value, nil
}

// ListFolder lists all keys in a folder.
//   - prefix: like finops
//   - folder: with trailing slash.
func (c *API) ListFolder(ctx context.Context, prefix, folder string) ([]string, error) {
	if c.client == nil {
		return nil, fmt.Errorf("consul client is not initialized")
	}

	v, _, err := c.client.Keys(path.Join(prefix, folder), "", nil)

	return v, err
}

func (c *API) Travel(ctx context.Context, prefix, folder string, fn func(key string, value []byte) error) error {
	if c.client == nil {
		return fmt.Errorf("consul client is not initialized")
	}

	v, _, err := c.client.List(prefix, nil)
	if err != nil {
		return err
	}

	for _, item := range v {
		if err := fn(item.Key, item.Value); err != nil {
			return err
		}
	}

	return nil
}
