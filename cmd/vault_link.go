package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/matthewmyrick/vaultex/internal/vault"
	"github.com/spf13/cobra"
)

var vaultLinkName string

var vaultLinkCmd = &cobra.Command{
	Use:   "link <vault> <path>",
	Short: "Symlink a directory into a vault",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName := args[0]
		targetPath := args[1]

		if _, err := cfg.GetVault(vaultName); err != nil {
			return err
		}

		absTarget, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}

		linkName := vaultLinkName
		if linkName == "" {
			linkName = filepath.Base(absTarget)
		}

		vaultsDir := config.GetVaultsDir()
		vaultPath := filepath.Join(vaultsDir, vaultName)

		if err := vault.AddLink(vaultPath, absTarget, linkName); err != nil {
			return err
		}

		if err := cfg.AddLink(vaultName, linkName, absTarget); err != nil {
			vault.RemoveLink(vaultPath, linkName)
			return err
		}

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Linked %q -> %s in vault %q\n", linkName, absTarget, vaultName)
		return nil
	},
}

func init() {
	vaultLinkCmd.Flags().StringVarP(&vaultLinkName, "name", "n", "", "name for the symlink (default: directory basename)")
}
