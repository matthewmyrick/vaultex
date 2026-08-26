package search

import (
	"os"
	"regexp"
	"strings"

	"github.com/matthewmyrick/vaultex/internal/vault"
)

type SearchEntry struct {
	Title    string
	Filename string
	Dir      string
	Path     string
	Content  string
	Preview  string
}

type Index struct {
	Entries []SearchEntry
}

func BuildIndex(files []vault.MarkdownFile) (*Index, error) {
	entries := make([]SearchEntry, 0, len(files))

	for _, f := range files {
		content, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}

		text := string(content)
		title := extractTitle(text)
		if title == "" {
			title = f.Name
		}

		entries = append(entries, SearchEntry{
			Title:    title,
			Filename: f.Name,
			Dir:      f.Dir,
			Path:     f.Path,
			Content:  stripMarkdown(text),
			Preview:  generatePreview(text, 200),
		})
	}

	return &Index{Entries: entries}, nil
}

var headingRegex = regexp.MustCompile(`(?m)^#\s+(.+)$`)

func extractTitle(content string) string {
	match := headingRegex.FindStringSubmatch(content)
	if len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func generatePreview(content string, maxLen int) string {
	lines := strings.Split(content, "\n")
	var preview strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "---") {
			continue
		}

		if preview.Len() > 0 {
			preview.WriteString(" ")
		}
		preview.WriteString(line)

		if preview.Len() >= maxLen {
			break
		}
	}

	result := preview.String()
	if len(result) > maxLen {
		result = result[:maxLen] + "..."
	}
	return result
}

var markdownSyntax = regexp.MustCompile(`[#*_\[\]()` + "`" + `>~|]`)

func stripMarkdown(content string) string {
	result := markdownSyntax.ReplaceAllString(content, " ")
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")
	return strings.TrimSpace(result)
}
