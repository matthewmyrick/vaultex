package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/spf13/cobra"
)

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vaults",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		names := cfg.VaultNames()
		if len(names) == 0 {
			fmt.Println("No vaults configured. Create one with: vaultex vault create <name>")
			return nil
		}

		vaultsDir := config.GetVaultsDir()

		for _, name := range names {
			vc, _ := cfg.GetVault(name)
			fmt.Printf("\n  %s", name)
			if vc.Description != "" {
				fmt.Printf("  -  %s", vc.Description)
			}
			fmt.Println()

			if len(vc.Links) == 0 {
				fmt.Println("    (no links)")
				continue
			}

			for _, link := range vc.Links {
				symPath := filepath.Join(vaultsDir, name, link.Name)
				status := "ok"
				if _, err := os.Stat(symPath); err != nil {
					status = "broken"
				}
				fmt.Printf("    %s -> %s [%s]\n", link.Name, link.Target, status)
			}
		}
		fmt.Println()
		return nil
	},
}
