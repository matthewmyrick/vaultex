package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	MatchStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	PromptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212"))
)
