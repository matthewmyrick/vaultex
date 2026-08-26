package cmd

import (
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage vaults",
	Long:  "Create, delete, list vaults and manage their symlinks.",
}

func init() {
	vaultCmd.AddCommand(vaultCreateCmd)
	vaultCmd.AddCommand(vaultDeleteCmd)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultLinkCmd)
	vaultCmd.AddCommand(vaultUnlinkCmd)
}
