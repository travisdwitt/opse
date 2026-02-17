package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit          key.Binding
	Tab           key.Binding
	Enter         key.Binding
	Escape        key.Binding
	Up            key.Binding
	Down          key.Binding
	PageUp        key.Binding
	PageDown      key.Binding
	Help          key.Binding
	NewAdventure  key.Binding
	OpenAdventure key.Binding
	Save          key.Binding
	SavedRolls    key.Binding
	Portraits     key.Binding

	OracleLikely   key.Binding
	OracleEven     key.Binding
	OracleUnlikely key.Binding
	OracleHow      key.Binding
	ActionFocus    key.Binding
	DetailFocus    key.Binding
	TopicFocus     key.Binding
	SetScene       key.Binding
	RandomEvent    key.Binding
	PacingMove     key.Binding
	FailureMove    key.Binding
	Generic        key.Binding
}

var DefaultKeys = KeyMap{
	Quit:          key.NewBinding(key.WithKeys("ctrl+q"), key.WithHelp("ctrl+q", "quit")),
	Tab:           key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch panel")),
	Enter:         key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit/run")),
	Escape:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back to input")),
	Up:            key.NewBinding(key.WithKeys("up", "k")),
	Down:          key.NewBinding(key.WithKeys("down", "j")),
	PageUp:        key.NewBinding(key.WithKeys("pgup")),
	PageDown:      key.NewBinding(key.WithKeys("pgdown")),
	Help:          key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	NewAdventure:  key.NewBinding(key.WithKeys("ctrl+n")),
	OpenAdventure: key.NewBinding(key.WithKeys("ctrl+o")),
	Save:          key.NewBinding(key.WithKeys("ctrl+s")),
	SavedRolls:    key.NewBinding(key.WithKeys("ctrl+r")),
	Portraits:     key.NewBinding(key.WithKeys("ctrl+p")),

	OracleLikely:   key.NewBinding(key.WithKeys("1")),
	OracleEven:     key.NewBinding(key.WithKeys("2")),
	OracleUnlikely: key.NewBinding(key.WithKeys("3")),
	OracleHow:      key.NewBinding(key.WithKeys("4")),
	ActionFocus:    key.NewBinding(key.WithKeys("5")),
	DetailFocus:    key.NewBinding(key.WithKeys("6")),
	TopicFocus:     key.NewBinding(key.WithKeys("7")),
	SetScene:       key.NewBinding(key.WithKeys("8")),
	RandomEvent:    key.NewBinding(key.WithKeys("9")),
	PacingMove:     key.NewBinding(key.WithKeys("0")),
	FailureMove:    key.NewBinding(key.WithKeys("-")),
	Generic:        key.NewBinding(key.WithKeys("=")),
}
