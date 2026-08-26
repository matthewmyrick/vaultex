package tui

import (
	"fmt"
	"strings"

	"github.com/matthewmyrick/vaultex/internal/search"
)

func RenderResult(r search.Result, selected bool, width int) string {
	var b strings.Builder

	nameStyle := NormalStyle
	if selected {
		nameStyle = SelectedStyle
	}

	title := r.Entry.Title
	if len(title) > width-10 {
		title = title[:width-13] + "..."
	}

	b.WriteString(nameStyle.Render(title))
	b.WriteString("  ")
	b.WriteString(DimStyle.Render(r.Entry.Dir+"/"+r.Entry.Filename+".md"))
	b.WriteString("\n")

	preview := r.Entry.Preview
	maxPreview := width - 6
	if maxPreview > 0 && len(preview) > maxPreview {
		preview = preview[:maxPreview] + "..."
	}
	b.WriteString(DimStyle.Render(preview))

	return b.String()
}

func RenderResultList(results []search.Result, cursor, maxVisible, width int) string {
	if len(results) == 0 {
		return DimStyle.Render("  No results found")
	}

	offset := 0
	if cursor >= maxVisible {
		offset = cursor - maxVisible + 1
	}

	end := offset + maxVisible
	if end > len(results) {
		end = len(results)
	}

	var b strings.Builder

	if offset > 0 {
		b.WriteString(DimStyle.Render(fmt.Sprintf("  ↑ %d more above", offset)))
		b.WriteString("\n")
	}

	for i := offset; i < end; i++ {
		selected := i == cursor
		prefix := "  "
		if selected {
			prefix = CursorStyle.Render("> ")
		}

		b.WriteString("  " + prefix)
		b.WriteString(RenderResult(results[i], selected, width-4))
		b.WriteString("\n")
	}

	remaining := len(results) - end
	if remaining > 0 {
		b.WriteString(DimStyle.Render(fmt.Sprintf("  ↓ %d more below", remaining)))
		b.WriteString("\n")
	}

	return b.String()
}
