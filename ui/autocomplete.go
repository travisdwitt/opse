package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var allCommands = []string{
	"roll", "r", "flip", "f", "draw", "card", "shuffle",
	"dir", "direction", "weather", "w", "color", "sound", "scene",
	"char",
}

type AutocompleteModel struct {
	visible     bool
	suggestions []string
	cursor      int
	query       string
}

// fuzzyMatch returns true if all characters in pattern appear in s in order.
func fuzzyMatch(s, pattern string) bool {
	pi := 0
	for i := 0; i < len(s) && pi < len(pattern); i++ {
		if s[i] == pattern[pi] {
			pi++
		}
	}
	return pi == len(pattern)
}

func (a *AutocompleteModel) Update(text string) {
	if !strings.HasPrefix(text, "/") || strings.Contains(text, " ") {
		a.visible = false
		a.suggestions = nil
		a.cursor = 0
		a.query = ""
		return
	}

	query := strings.ToLower(text[1:])
	if query == a.query && a.visible {
		return
	}
	a.query = query

	var matches []string
	if query == "" {
		// Bare "/" — show all commands
		matches = append(matches, allCommands...)
	} else {
		for _, cmd := range allCommands {
			if fuzzyMatch(cmd, query) {
				matches = append(matches, cmd)
			}
		}
	}

	// Hide if no matches, or if the only match is an exact match
	if len(matches) == 0 || (len(matches) == 1 && matches[0] == query) {
		a.visible = false
		a.suggestions = nil
		a.cursor = 0
		return
	}

	a.suggestions = matches
	a.visible = true
	if a.cursor >= len(matches) {
		a.cursor = 0
	}
}

func (a *AutocompleteModel) MoveUp() {
	if a.cursor > 0 {
		a.cursor--
	}
}

func (a *AutocompleteModel) MoveDown() {
	if a.cursor < len(a.suggestions)-1 {
		a.cursor++
	}
}

func (a *AutocompleteModel) Complete() string {
	if !a.visible || len(a.suggestions) == 0 {
		return ""
	}
	return "/" + a.suggestions[a.cursor] + " "
}

func (a *AutocompleteModel) Hide() {
	a.visible = false
	a.suggestions = nil
	a.cursor = 0
	a.query = ""
}

var (
	acItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	acSelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true).
			Padding(0, 1)

	acBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 0)
)

func (a *AutocompleteModel) View() string {
	if !a.visible || len(a.suggestions) == 0 {
		return ""
	}
	var items []string
	for i, s := range a.suggestions {
		label := "/" + s
		if i == a.cursor {
			items = append(items, acSelectedStyle.Render(label))
		} else {
			items = append(items, acItemStyle.Render(label))
		}
	}
	list := lipgloss.JoinVertical(lipgloss.Left, items...)
	return acBoxStyle.Render(list)
}
