package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommandMsg struct {
	Command string
	Args    []string
}

type NarrativeMsg struct {
	Text string
}

type InputModel struct {
	textarea     textarea.Model
	autocomplete AutocompleteModel
	focused      bool
}

func NewInput() InputModel {
	ta := textarea.New()
	ta.Placeholder = "Write your story... (or /roll 2d6, /flip, /weather)"
	ta.CharLimit = 10000
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.Focus()
	return InputModel{textarea: ta, focused: true}
}

func (m *InputModel) Focus()        { m.textarea.Focus(); m.focused = true }
func (m *InputModel) Blur()         { m.textarea.Blur(); m.focused = false }
func (m *InputModel) Focused() bool { return m.focused }

func (m *InputModel) Update(msg tea.Msg) (*InputModel, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok && m.autocomplete.visible {
		switch kmsg.String() {
		case "tab":
			completed := m.autocomplete.Complete()
			if completed != "" {
				m.textarea.SetValue(completed)
				m.textarea.CursorEnd()
				m.autocomplete.Hide()
			}
			return m, nil
		case "up":
			m.autocomplete.MoveUp()
			return m, nil
		case "down":
			m.autocomplete.MoveDown()
			return m, nil
		case "esc":
			m.autocomplete.Hide()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	m.autocomplete.Update(m.textarea.Value())
	return m, cmd
}

func (m *InputModel) Submit() tea.Msg {
	text := strings.TrimSpace(m.textarea.Value())
	m.textarea.Reset()
	m.autocomplete.Hide()
	if text == "" {
		return nil
	}

	if strings.HasPrefix(text, "/") {
		parts := strings.Fields(text)
		cmd := strings.ToLower(strings.TrimPrefix(parts[0], "/"))
		known := map[string]bool{
			"roll": true, "r": true, "flip": true, "f": true,
			"draw": true, "card": true, "shuffle": true,
			"dir": true, "direction": true,
			"weather": true, "w": true,
			"color": true, "sound": true,
			"scene": true, "char": true,
		}
		if known[cmd] {
			return CommandMsg{Command: cmd, Args: parts[1:]}
		}
	}

	return NarrativeMsg{Text: text}
}

func (m *InputModel) View(width int, focused bool) string {
	style := InputStyle
	if focused {
		style = InputStyleFocused
	}
	inputBox := style.Width(width).Render(m.textarea.View())

	if m.autocomplete.visible {
		acView := m.autocomplete.View()
		return lipgloss.JoinVertical(lipgloss.Left, acView, inputBox)
	}
	return inputBox
}
