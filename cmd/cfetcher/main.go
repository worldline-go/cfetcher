package main

import (
	"context"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/rakunlabs/into"
	"github.com/rakunlabs/logi"
	"github.com/spf13/cobra"

	"github.com/worldline-go/cfetcher/internal/config"
	"github.com/worldline-go/cfetcher/pkg/loader/consul"
	"github.com/worldline-go/cfetcher/pkg/loader/vault"
	"github.com/worldline-go/cfetcher/pkg/utils/client"
	"github.com/worldline-go/cfetcher/pkg/utils/file"
)

var (
	version = "v0.0.0"
	commit  = "-"
	date    = "-"
)

var values = struct {
	VaultSave   string
	VaultPath   string
	VaultPrefix string

	ConsulSet    string
	ConsulSave   string
	ConsulPath   string
	ConsulPrefix string
}{
	VaultSave:   "",
	VaultPath:   "",
	VaultPrefix: "finops",

	ConsulSave:   "",
	ConsulPath:   "",
	ConsulPrefix: "finops",
}

var rootCmd = &cobra.Command{
	Use:           "cfetcher",
	Short:         "config fetcher",
	Long:          "load, save and change configuration in various places",
	SilenceUsage:  true,
	SilenceErrors: true,
	Example:       "cfetcher --vault-save=./out/finops_vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	rootCmd.Flags().StringVar(&values.VaultPath, "vault-path", values.VaultPath, "vault path to load, multiple comma or space separated (CONFIG_VAULT_PATH)")
	rootCmd.Flags().StringVar(&values.VaultPrefix, "vault-prefix", values.VaultPrefix, "vault prefix to load")
	rootCmd.Flags().StringVar(&values.VaultSave, "vault-save", values.VaultSave, "vault save to folder")

	rootCmd.Flags().StringVar(&values.ConsulSet, "consul-set", values.ConsulSet, "consul path to set")
	rootCmd.Flags().StringVar(&values.ConsulPath, "consul-path", values.ConsulPath, "consul path to load, multiple comma or space separated (CONFIG_CONSUL_PATH)")
	rootCmd.Flags().StringVar(&values.ConsulPrefix, "consul-prefix", values.ConsulPrefix, "consul prefix to load")
	rootCmd.Flags().StringVar(&values.ConsulSave, "consul-save", values.ConsulSave, "consul save to folder")
}

func main() {
	into.Init(
		runCommand,
		into.WithLogger(logi.InitializeLog(logi.WithCaller(false))),
		into.WithMsgf("cfetcher [%s]", version),
		into.WithStartFn(nil),
		into.WithStopFn(nil),
	)
}

func runCommand(ctx context.Context) error {
	rootCmd.Version = version
	rootCmd.Long += "\nversion: " + version + " commit: " + commit + " buildDate:" + date

	return rootCmd.ExecuteContext(ctx)
}

func run(ctx context.Context) error {
	slog.Info("cfetcher [" + version + "] commit: " + commit + "buildDate: " + date)

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if values.ConsulSet != "" {
		// get httpClient for custom authentication
		httpClient, err := client.NewHTTPClient(ctx, cfg.AuthService)
		if err != nil {
			return err
		}

		slog.Info("consul set from folder " + values.ConsulSet)

		if err := cfg.Loaders.Consul.Connect(consul.WithClient(httpClient)); err != nil {
			return err
		}

		if err := cfg.Loaders.Consul.Set(ctx, values.ConsulSet); err != nil {
			return err
		}
	}

	if values.ConsulSave != "" {
		// get httpClient for custom authentication
		httpClient, err := client.NewHTTPClient(ctx, cfg.AuthService)
		if err != nil {
			return err
		}

		slog.Info("consul save to folder " + values.ConsulSave)

		if err := cfg.Loaders.Consul.Connect(consul.WithClient(httpClient)); err != nil {
			return err
		}

		if values.ConsulPath == "" {
			values.ConsulPath = os.Getenv("CONFIG_CONSUL_PATH")
		}

		vPaths := strings.Fields(strings.ReplaceAll(values.ConsulPath, ",", " "))
		if len(vPaths) == 0 {
			// add empty to load all
			vPaths = append(vPaths, "")
		}

		for _, vPath := range vPaths {
			if err := cfg.Loaders.Consul.Travel(ctx, values.ConsulPrefix, vPath, func(key string, value []byte) error {
				// remove prefix in key
				key = strings.TrimPrefix(key, values.ConsulPrefix+"/")
				slog.Info("consul config saving [" + key + "]")

				filePath := path.Join(values.ConsulSave, key) + ".yaml"
				f, err := file.OpenFileWrite(filePath)
				if err != nil {
					return err
				}

				defer f.Close()

				if _, err := f.Write(value); err != nil {
					return err
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	if values.VaultSave != "" {
		slog.Info("vault save to folder " + values.VaultSave)

		// get httpClient for custom authentication
		httpClient, err := client.NewHTTPClient(ctx, client.Provider{})
		if err != nil {
			return err
		}

		if err := cfg.Loaders.Vault.Connect(vault.WithClient(httpClient)); err != nil {
			return err
		}

		if values.VaultPath == "" {
			values.VaultPath = os.Getenv("CONFIG_VAULT_PATH")
		}

		vPaths := strings.Fields(strings.ReplaceAll(values.VaultPath, ",", " "))
		if len(vPaths) == 0 {
			// add empty to load all
			vPaths = append(vPaths, "")
		}

		for _, vPath := range vPaths {
			if err := cfg.Loaders.Vault.Travel(ctx, values.VaultPrefix, vPath, func(key string, value []byte) error {
				slog.Info("vault config saving [" + key + "]")

				filePath := path.Join(values.VaultSave, key) + ".json"
				f, err := file.OpenFileWrite(filePath)
				if err != nil {
					return err
				}

				defer f.Close()

				if _, err := f.Write(value); err != nil {
					return err
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
