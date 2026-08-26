package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/matthewmyrick/vaultex/internal/vault"
	"github.com/spf13/cobra"
)

var vaultUnlinkCmd = &cobra.Command{
	Use:   "unlink <vault> <link-name>",
	Short: "Remove a symlink from a vault",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName := args[0]
		linkName := args[1]

		if _, err := cfg.GetVault(vaultName); err != nil {
			return err
		}

		vaultsDir := config.GetVaultsDir()
		vaultPath := filepath.Join(vaultsDir, vaultName)

		if err := vault.RemoveLink(vaultPath, linkName); err != nil {
			return err
		}

		if err := cfg.RemoveLink(vaultName, linkName); err != nil {
			return err
		}

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Unlinked %q from vault %q\n", linkName, vaultName)
		return nil
	},
}
