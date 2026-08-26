package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/matthewmyrick/vaultex/internal/search"
	"github.com/matthewmyrick/vaultex/internal/shell"
	"github.com/matthewmyrick/vaultex/internal/tui"
	"github.com/matthewmyrick/vaultex/internal/vault"
	"github.com/matthewmyrick/vaultex/internal/viewer"
	"github.com/spf13/cobra"
)

var searchVault string

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search markdown files in a vault",
	Long:  "Live search across content, titles, filenames, and directories within a vault.",
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultName, err := resolveVault()
		if err != nil {
			return err
		}
		if vaultName == "" {
			return nil
		}

		vaultsDir := config.GetVaultsDir()
		vaultPath := filepath.Join(vaultsDir, vaultName)

		files, err := vault.ScanMarkdown(vaultPath)
		if err != nil {
			return fmt.Errorf("failed to scan vault: %w", err)
		}

		if len(files) == 0 {
			fmt.Println("No markdown files found in vault.")
			return nil
		}

		index, err := search.BuildIndex(files)
		if err != nil {
			return fmt.Errorf("failed to build search index: %w", err)
		}

		initialQuery := strings.Join(args, " ")

		result, err := tui.RunSearchBar(index, initialQuery)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if result != nil {
			return viewer.OpenInEditor(result.Entry.Path)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchVault, "vault", "v", "", "vault to search (default: auto-detect from env)")
}

func resolveVault() (string, error) {
	if searchVault != "" {
		if _, err := cfg.GetVault(searchVault); err != nil {
			return "", err
		}
		return searchVault, nil
	}

	name, _, active := shell.GetActiveVault()
	if active {
		return name, nil
	}

	names := cfg.VaultNames()
	if len(names) == 0 {
		return "", fmt.Errorf("no vaults configured - create one with: vaultex vault create <name>")
	}

	if len(names) == 1 {
		return names[0], nil
	}

	items := make([]tui.VaultItem, 0, len(names))
	for _, n := range names {
		vc, _ := cfg.GetVault(n)
		items = append(items, tui.VaultItem{
			Name:        n,
			Description: vc.Description,
			LinkCount:   len(vc.Links),
		})
	}

	return tui.RunPicker(items)
}
