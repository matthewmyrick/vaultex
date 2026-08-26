package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type VaultItem struct {
	Name        string
	Description string
	LinkCount   int
}

type PickerModel struct {
	items    []VaultItem
	filtered []VaultItem
	input    textinput.Model
	cursor   int
	choice   string
	quitting bool
}

func NewPicker(items []VaultItem) PickerModel {
	ti := textinput.New()
	ti.Placeholder = "Search vaults..."
	ti.Focus()
	ti.PromptStyle = PromptStyle
	ti.Prompt = "  > "

	return PickerModel{
		items:    items,
		filtered: items,
		input:    ti,
	}
}

func (m PickerModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				m.choice = m.filtered[m.cursor].Name
			}
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	query := strings.ToLower(m.input.Value())
	if query == "" {
		m.filtered = m.items
	} else {
		m.filtered = filterVaults(m.items, query)
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}

	return m, cmd
}

func (m PickerModel) View() string {
	if m.choice != "" || m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(TitleStyle.Render("  Select a vault"))
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	maxVisible := 10
	if len(m.filtered) < maxVisible {
		maxVisible = len(m.filtered)
	}

	if len(m.filtered) == 0 {
		b.WriteString(DimStyle.Render("    No matching vaults"))
		b.WriteString("\n")
	}

	for i := 0; i < maxVisible; i++ {
		item := m.filtered[i]
		cursor := "  "
		style := NormalStyle
		if i == m.cursor {
			cursor = CursorStyle.Render("> ")
			style = SelectedStyle
		}

		line := fmt.Sprintf("  %s%s", cursor, style.Render(item.Name))
		if item.Description != "" {
			line += DimStyle.Render(fmt.Sprintf("  %s", item.Description))
		}
		line += DimStyle.Render(fmt.Sprintf("  (%d links)", item.LinkCount))

		b.WriteString(line)
		b.WriteString("\n")
	}

	if len(m.filtered) > 10 {
		b.WriteString(DimStyle.Render(fmt.Sprintf("    ... and %d more", len(m.filtered)-10)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(DimStyle.Render("  ↑/↓ navigate • enter select • esc quit"))
	b.WriteString("\n")

	return b.String()
}

func (m PickerModel) Choice() string {
	return m.choice
}

func RunPicker(items []VaultItem) (string, error) {
	model := NewPicker(items)
	p := tea.NewProgram(model)
	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("picker error: %w", err)
	}

	picker := result.(PickerModel)
	return picker.Choice(), nil
}

func filterVaults(items []VaultItem, query string) []VaultItem {
	var filtered []VaultItem
	for _, item := range items {
		name := strings.ToLower(item.Name)
		desc := strings.ToLower(item.Description)
		if strings.Contains(name, query) || strings.Contains(desc, query) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
