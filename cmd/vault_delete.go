package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/matthewmyrick/vaultex/internal/config"
	"github.com/matthewmyrick/vaultex/internal/vault"
	"github.com/spf13/cobra"
)

var vaultDeleteForce bool

var vaultDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if _, err := cfg.GetVault(name); err != nil {
			return err
		}

		if !vaultDeleteForce {
			fmt.Printf("Delete vault %q? This removes all symlinks (not source files). [y/N] ", name)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		vaultsDir := config.GetVaultsDir()
		if err := vault.Delete(vaultsDir, name); err != nil {
			return fmt.Errorf("failed to delete vault directory: %w", err)
		}

		if err := cfg.RemoveVault(name); err != nil {
			return err
		}

		if err := saveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("Deleted vault %q\n", name)
		return nil
	},
}

func init() {
	vaultDeleteCmd.Flags().BoolVarP(&vaultDeleteForce, "force", "f", false, "skip confirmation prompt")
}
