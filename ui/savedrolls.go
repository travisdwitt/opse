package ui

import (
	"fmt"
	"strings"

	"opse/engine"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type srState int

const (
	srBrowsing srState = iota
	srCreatingName
	srCreatingExpr
	srCreatingFolder
	srConfirmDelete
)

type srItem struct {
	isFolder  bool
	collapsed bool // only meaningful for folder items
	rollID    string
	label     string
	folder    string
}

type SavedRollsModel struct {
	state     srState
	config    *engine.SavedRollsConfig
	collapsed map[string]bool // folder name → collapsed
	items     []srItem
	cursor    int
	scrollOff int
	input     textinput.Model
	newName   string
	errMsg    string
}

func NewSavedRolls(cfg *engine.SavedRollsConfig) SavedRollsModel {
	ti := textinput.New()
	ti.CharLimit = 64
	m := SavedRollsModel{
		config:    cfg,
		collapsed: make(map[string]bool),
		input:     ti,
	}
	m.rebuildItems()
	return m
}

func (m *SavedRollsModel) rebuildItems() {
	m.items = nil
	byFolder := m.config.ByFolder()

	// Unsorted rolls first
	if rolls, ok := byFolder[""]; ok {
		for _, r := range rolls {
			m.items = append(m.items, srItem{
				rollID: r.ID,
				label:  fmt.Sprintf("%s  %s", r.Name, DimStyle.Render(r.Expression)),
			})
		}
	}

	// Folder groups
	for _, folder := range m.config.Folders {
		isCollapsed := m.collapsed[folder.Name]
		m.items = append(m.items, srItem{
			isFolder:  true,
			collapsed: isCollapsed,
			label:     folder.Name,
			folder:    folder.Name,
		})
		if !isCollapsed {
			for _, r := range byFolder[folder.Name] {
				m.items = append(m.items, srItem{
					rollID: r.ID,
					label:  fmt.Sprintf("  %s  %s", r.Name, DimStyle.Render(r.Expression)),
					folder: folder.Name,
				})
			}
		}
	}

	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *SavedRollsModel) SetConfig(cfg *engine.SavedRollsConfig) {
	m.config = cfg
	m.rebuildItems()
}

// Update returns a rollID to execute (non-empty when user selects a roll).
func (m *SavedRollsModel) Update(msg tea.Msg) (string, tea.Cmd) {
	switch m.state {
	case srBrowsing:
		return m.updateBrowsing(msg)
	case srCreatingName:
		return m.updateCreatingName(msg)
	case srCreatingExpr:
		return m.updateCreatingExpr(msg)
	case srCreatingFolder:
		return m.updateCreatingFolder(msg)
	case srConfirmDelete:
		return m.updateConfirmDelete(msg)
	}
	return "", nil
}

func (m *SavedRollsModel) updateBrowsing(msg tea.Msg) (string, tea.Cmd) {
	kmsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return "", nil
	}
	switch kmsg.String() {
	case "j", "down":
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "left", "h":
		if len(m.items) > 0 && m.items[m.cursor].isFolder {
			m.collapsed[m.items[m.cursor].folder] = true
			m.rebuildItems()
		} else if len(m.items) > 0 && m.items[m.cursor].folder != "" {
			// Collapse parent folder and move cursor to it
			m.collapsed[m.items[m.cursor].folder] = true
			folder := m.items[m.cursor].folder
			m.rebuildItems()
			for i, item := range m.items {
				if item.isFolder && item.folder == folder {
					m.cursor = i
					break
				}
			}
		}
	case "right", "l":
		if len(m.items) > 0 && m.items[m.cursor].isFolder {
			m.collapsed[m.items[m.cursor].folder] = false
			m.rebuildItems()
		}
	case "enter":
		if len(m.items) > 0 {
			item := m.items[m.cursor]
			if item.isFolder {
				// Toggle collapse on enter too
				m.collapsed[item.folder] = !m.collapsed[item.folder]
				m.rebuildItems()
			} else {
				return item.rollID, nil
			}
		}
	case "n":
		m.state = srCreatingName
		m.input.SetValue("")
		m.input.Placeholder = "Roll name (e.g. Attack)"
		m.input.Focus()
		m.errMsg = ""
		return "", m.input.Cursor.BlinkCmd()
	case "f":
		m.state = srCreatingFolder
		m.input.SetValue("")
		m.input.Placeholder = "Folder name"
		m.input.Focus()
		m.errMsg = ""
		return "", m.input.Cursor.BlinkCmd()
	case "d":
		if len(m.items) > 0 {
			m.state = srConfirmDelete
		}
	}
	return "", nil
}

func (m *SavedRollsModel) updateCreatingName(msg tea.Msg) (string, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "enter":
			name := strings.TrimSpace(m.input.Value())
			if name == "" {
				return "", nil
			}
			m.newName = name
			m.state = srCreatingExpr
			m.input.SetValue("")
			m.input.Placeholder = "Dice expression (e.g. 2d6+3)"
			m.errMsg = ""
			return "", nil
		case "esc":
			m.state = srBrowsing
			return "", nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return "", cmd
}

func (m *SavedRollsModel) updateCreatingExpr(msg tea.Msg) (string, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "enter":
			expr := strings.TrimSpace(m.input.Value())
			if expr == "" {
				return "", nil
			}
			if _, err := engine.ParseDice(expr); err != nil {
				m.errMsg = "Invalid dice expression"
				return "", nil
			}
			// Determine folder: if cursor is on a folder or a roll inside a folder, use that folder
			folder := ""
			if len(m.items) > 0 && m.cursor < len(m.items) {
				item := m.items[m.cursor]
				if item.isFolder || item.folder != "" {
					folder = item.folder
				}
			}
			m.config.Add(engine.SavedRoll{
				ID:         fmt.Sprintf("sr_%d", len(m.config.Rolls)+1),
				Name:       m.newName,
				Expression: expr,
				Folder:     folder,
			})
			engine.SaveSavedRolls(m.config)
			m.rebuildItems()
			m.state = srBrowsing
			m.errMsg = ""
			return "", nil
		case "esc":
			m.state = srBrowsing
			m.errMsg = ""
			return "", nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return "", cmd
}

func (m *SavedRollsModel) updateCreatingFolder(msg tea.Msg) (string, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "enter":
			name := strings.TrimSpace(m.input.Value())
			if name == "" {
				return "", nil
			}
			m.config.AddFolder(name)
			engine.SaveSavedRolls(m.config)
			m.rebuildItems()
			m.state = srBrowsing
			return "", nil
		case "esc":
			m.state = srBrowsing
			return "", nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return "", cmd
}

func (m *SavedRollsModel) updateConfirmDelete(msg tea.Msg) (string, tea.Cmd) {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "y":
			if len(m.items) > 0 {
				item := m.items[m.cursor]
				if item.isFolder {
					m.config.DeleteFolder(item.folder)
					delete(m.collapsed, item.folder)
				} else {
					m.config.Delete(item.rollID)
				}
				engine.SaveSavedRolls(m.config)
				m.rebuildItems()
			}
			m.state = srBrowsing
		case "n", "esc":
			m.state = srBrowsing
		}
	}
	return "", nil
}

func (m *SavedRollsModel) View(width, height int) string {
	boxW := width - 4
	if boxW < 40 {
		boxW = 40
	}
	boxH := height - 2
	if boxH < 10 {
		boxH = 10
	}

	contentH := boxH - 4 // padding(2) + header(1) + footer(1)

	header := ResultLabelStyle.Render("Saved Rolls")

	var body string
	switch m.state {
	case srBrowsing, srConfirmDelete:
		body = m.viewBrowse(contentH - 2) // room for confirm line(s)
		if m.state == srConfirmDelete && len(m.items) > 0 {
			item := m.items[m.cursor]
			warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
			if item.isFolder {
				count := m.config.FolderRollCount(item.folder)
				if count > 0 {
					body += "\n" + warnStyle.Render(
						fmt.Sprintf("Deleting this folder will delete all %d saved roll(s) in it. Continue? (y/n)", count))
				} else {
					body += "\n" + warnStyle.Render(
						fmt.Sprintf("Delete folder \"%s\"? (y/n)", item.folder))
				}
			} else {
				body += "\n" + warnStyle.Render(
					fmt.Sprintf("Delete \"%s\"? (y/n)", item.label))
			}
		}
	case srCreatingName:
		body = m.viewBrowse(contentH - 3)
		body += "\n" + ResultLabelStyle.Render("New Roll — Name:") + "\n" + m.input.View()
		if m.errMsg != "" {
			body += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.errMsg)
		}
	case srCreatingExpr:
		body = m.viewBrowse(contentH - 3)
		body += "\n" + ResultLabelStyle.Render(fmt.Sprintf("New Roll \"%s\" — Expression:", m.newName)) + "\n" + m.input.View()
		if m.errMsg != "" {
			body += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.errMsg)
		}
	case srCreatingFolder:
		body = m.viewBrowse(contentH - 3)
		body += "\n" + ResultLabelStyle.Render("New Folder — Name:") + "\n" + m.input.View()
	}

	var footerText string
	switch m.state {
	case srBrowsing:
		footerText = "j/k: navigate | ←/→: collapse/expand | Enter: roll | n: new | f: folder | d: delete | Esc: close"
	case srCreatingName, srCreatingExpr, srCreatingFolder:
		footerText = "Enter: confirm | Esc: cancel"
	case srConfirmDelete:
		footerText = "y: confirm | n: cancel"
	}
	footer := DimStyle.Render(footerText)

	content := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(boxW).
		Height(boxH).
		Render(content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

func (m *SavedRollsModel) viewBrowse(maxLines int) string {
	if len(m.items) == 0 {
		return DimStyle.Render("No saved rolls. Press 'n' to create one.")
	}

	var lines []string
	for i, item := range m.items {
		var line string
		if item.isFolder {
			arrow := "▾"
			if item.collapsed {
				arrow = "▸"
			}
			prefix := arrow + " 📁 " + item.label
			if i == m.cursor {
				line = SrSelectedStyle.Render(prefix)
			} else {
				line = CategoryStyle.Render(prefix)
			}
		} else if i == m.cursor {
			line = SrSelectedStyle.Render(item.label)
		} else {
			line = ItemStyle.Render("  " + item.label)
		}
		lines = append(lines, line)
	}

	// Scrolling
	if maxLines < 1 {
		maxLines = 1
	}
	if m.cursor < m.scrollOff {
		m.scrollOff = m.cursor
	}
	if m.cursor >= m.scrollOff+maxLines {
		m.scrollOff = m.cursor - maxLines + 1
	}

	end := m.scrollOff + maxLines
	if end > len(lines) {
		end = len(lines)
	}

	return strings.Join(lines[m.scrollOff:end], "\n")
}
