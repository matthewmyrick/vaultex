package viewer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/matthewmyrick/vaultex/internal/shell"
	"golang.org/x/term"
)

func RenderMarkdown(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create renderer: %w", err)
	}

	rendered, err := renderer.Render(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	return rendered, nil
}

func ViewInPager(content string) error {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}

	cmd := exec.Command(pager, "-R")
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func OpenInEditor(filePath string) error {
	editor := shell.DetectEditor()

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ViewFile(filePath string, raw bool) error {
	if raw {
		return OpenInEditor(filePath)
	}

	rendered, err := RenderMarkdown(filePath)
	if err != nil {
		return err
	}

	return ViewInPager(rendered)
}

