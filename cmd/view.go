package cmd

import (
	"github.com/matthewmyrick/vaultex/internal/viewer"
	"github.com/spf13/cobra"
)

var viewRaw bool

var viewCmd = &cobra.Command{
	Use:   "view <file>",
	Short: "View a markdown file",
	Long:  "Render a markdown file with syntax highlighting, or open it directly in your editor.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return viewer.ViewFile(args[0], viewRaw)
	},
}

func init() {
	viewCmd.Flags().BoolVar(&viewRaw, "raw", false, "open directly in $EDITOR instead of rendering")
}
