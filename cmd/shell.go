package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/matthewmyrick/vaultex/internal/shell"
	"github.com/matthewmyrick/vaultex/internal/tui"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open an interactive shell in a vault",
	Long:  "Select a vault with fuzzy search, then drop into an interactive subshell scoped to that vault.",
	RunE: func(cmd *cobra.Command, args []string) error {
		names := cfg.VaultNames()
		if len(names) == 0 {
			return fmt.Errorf("no vaults configured - create one with: vaultex vault create <name>")
		}

		var vaultName string

		if len(names) == 1 {
			vaultName = names[0]
		} else {
			items := make([]tui.VaultItem, 0, len(names))
			for _, name := range names {
				vc, _ := cfg.GetVault(name)
				items = append(items, tui.VaultItem{
					Name:        name,
					Description: vc.Description,
					LinkCount:   len(vc.Links),
				})
			}

			choice, err := tui.RunPicker(items)
			if err != nil {
				return fmt.Errorf("vault picker failed: %w", err)
			}
			if choice == "" {
				return nil
			}
			vaultName = choice
		}

		vaultsDir := config.GetVaultsDir()
		vaultPath := filepath.Join(vaultsDir, vaultName)

		fmt.Printf("Entering vault %q...\n", vaultName)
		return shell.LaunchShell(vaultName, vaultPath)
	},
}
