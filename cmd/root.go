package cmd

import (
	"fmt"
	"os"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "vaultex",
	Short: "A CLI for managing and searching markdown document vaults",
	Long: `Vaultex is a CLI tool for organizing markdown documents into vaults.
It provides symlink-based vault management, an interactive shell with
fuzzy search, and live search across your markdown files.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return nil
		}
		return initConfig()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default: ./vaultex.yaml)")

	rootCmd.AddCommand(vaultCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(viewCmd)
}

func initConfig() error {
	var path string
	var err error

	if cfgFile != "" {
		path = cfgFile
	} else {
		path, err = config.ResolveConfigPath()
		if err != nil {
			return fmt.Errorf("failed to resolve config path: %w", err)
		}
	}

	config.SetConfigPath(path)

	cfg, err = config.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	vaultsDir := config.GetVaultsDir()
	if err := os.MkdirAll(vaultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create vaults directory: %w", err)
	}

	return nil
}

func saveConfig() error {
	return cfg.Save(config.GetConfigPath())
}
