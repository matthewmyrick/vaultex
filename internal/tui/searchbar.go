package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/matthewmyrick/vaultex/internal/search"
)

type SearchModel struct {
	input      textinput.Model
	index      *search.Index
	results    []search.Result
	cursor     int
	maxVisible int
	selected   *search.Result
	quitting   bool
	width      int
}

func NewSearchBar(index *search.Index, initialQuery string) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search markdown files..."
	ti.Focus()
	ti.PromptStyle = PromptStyle
	ti.Prompt = "  🔍 "
	ti.CharLimit = 256

	if initialQuery != "" {
		ti.SetValue(initialQuery)
	}

	m := SearchModel{
		input:      ti,
		index:      index,
		maxVisible: 10,
		width:      80,
	}

	m.results = search.Search(index, initialQuery)

	return m
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if len(m.results) > 0 && m.cursor < len(m.results) {
				r := m.results[m.cursor]
				m.selected = &r
			}
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "ctrl+n":
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	prevValue := m.input.Value()
	m.input, cmd = m.input.Update(msg)

	if m.input.Value() != prevValue {
		m.results = search.Search(m.index, m.input.Value())
		m.cursor = 0
	}

	return m, cmd
}

func (m SearchModel) View() string {
	if m.selected != nil || m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.input.View())
	b.WriteString("\n")

	count := len(m.results)
	b.WriteString(DimStyle.Render(fmt.Sprintf("  %d result(s)", count)))
	b.WriteString("\n\n")

	b.WriteString(RenderResultList(m.results, m.cursor, m.maxVisible, m.width))
	b.WriteString("\n")
	b.WriteString(DimStyle.Render("  ↑/↓ navigate • enter open • esc quit"))
	b.WriteString("\n")

	return b.String()
}

func (m SearchModel) Selected() *search.Result {
	return m.selected
}

func RunSearchBar(index *search.Index, initialQuery string) (*search.Result, error) {
	model := NewSearchBar(index, initialQuery)
	p := tea.NewProgram(model)
	result, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}

	sm := result.(SearchModel)
	return sm.Selected(), nil
}
