package cmd

import (
	"fmt"
	"regexp"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/matthewmyrick/vaultex/internal/vault"
	"github.com/spf13/cobra"
)

var vaultCreateDesc string

var vaultCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if !isValidVaultName(name) {
			return fmt.Errorf("invalid vault name %q: use only lowercase letters, numbers, and hyphens", name)
		}

		if err := cfg.AddVault(name, vaultCreateDesc); err != nil {
			return err
		}

		vaultsDir := config.GetVaultsDir()
		if err := vault.Create(vaultsDir, name); err != nil {
			cfg.RemoveVault(name)
			return fmt.Errorf("failed to create vault directory: %w", err)
		}

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Created vault %q\n", name)
		return nil
	},
}

func init() {
	vaultCreateCmd.Flags().StringVarP(&vaultCreateDesc, "description", "d", "", "vault description")
}

var validVaultName = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)

func isValidVaultName(name string) bool {
	return validVaultName.MatchString(name)
}
