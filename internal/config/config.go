package config

import (
	"log/slog"
	"os"

	"github.com/worldline-go/cfetcher/pkg/loader"
	"github.com/worldline-go/cfetcher/pkg/utils/client"
	"github.com/worldline-go/cfetcher/pkg/utils/file"
)

var EnvConfigFile = "CONFIG_FILE"

type Config struct {
	Loaders loader.Loaders `cfg:"loaders"`

	AuthService client.Provider `cfg:"auth_service"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	configFile := os.Getenv(EnvConfigFile)
	if configFile == "" {
		return cfg, nil
	}

	slog.Info("loading config from " + configFile)

	if err := file.Load(configFile, cfg); err != nil {
		return nil, err
	}

	slog.Info("config loaded", slog.Any("config", cfg))

	return cfg, nil
}
