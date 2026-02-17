package ui

import (
	"fmt"
	"strings"

	"opse/engine"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pbState int

const (
	pbBrowsing pbState = iota
	pbNaming
	pbConfirmDelete
)

type pbItem struct {
	isHeader bool
	isSaved  bool
	label    string
	params   engine.PortraitParams
	name     string // for saved portraits
}

// PortraitBrowserModel is a full-screen modal for generating and saving portraits.
type PortraitBrowserModel struct {
	state     pbState
	config    *engine.SavedPortraitsConfig
	rng       *engine.Randomizer
	generated []engine.PortraitParams
	items     []pbItem
	cursor    int
	scrollOff int
	input     textinput.Model
	errMsg    string
}

const generatedCount = 4

// NewPortraitBrowser creates a portrait browser with initial generated portraits.
func NewPortraitBrowser(cfg *engine.SavedPortraitsConfig, rng *engine.Randomizer) PortraitBrowserModel {
	ti := textinput.New()
	ti.CharLimit = 64
	m := PortraitBrowserModel{
		config: cfg,
		rng:    rng,
		input:  ti,
	}
	m.regenerate()
	return m
}

func (m *PortraitBrowserModel) regenerate() {
	m.generated = make([]engine.PortraitParams, generatedCount)
	for i := range m.generated {
		m.generated[i] = engine.GenerateRandomPortrait(m.rng)
	}
	m.rebuildItems()
}

func (m *PortraitBrowserModel) rebuildItems() {
	m.items = nil

	// Generated section
	m.items = append(m.items, pbItem{isHeader: true, label: "GENERATED"})
	for i, p := range m.generated {
		m.items = append(m.items, pbItem{
			label:  fmt.Sprintf("#%d  %s", i+1, engine.DescribePortrait(p)),
			params: p,
		})
	}

	// Saved section
	if len(m.config.Portraits) > 0 {
		m.items = append(m.items, pbItem{isHeader: true, label: "SAVED"})
		for _, sp := range m.config.Portraits {
			m.items = append(m.items, pbItem{
				isSaved: true,
				label:   sp.Name,
				params:  sp.Params,
				name:    sp.Name,
			})
		}
	}

	// Ensure cursor is on a valid non-header item
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.skipToNonHeader(1)
}

func (m *PortraitBrowserModel) skipToNonHeader(dir int) {
	for m.cursor >= 0 && m.cursor < len(m.items) && m.items[m.cursor].isHeader {
		m.cursor += dir
	}
	if m.cursor < 0 {
		m.cursor = 0
		m.skipToNonHeader(1)
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
		m.skipToNonHeader(-1)
	}
}

// SetConfig updates the saved portraits config and rebuilds the item list.
func (m *PortraitBrowserModel) SetConfig(cfg *engine.SavedPortraitsConfig) {
	m.config = cfg
	m.rebuildItems()
}

// Update handles key messages and returns a tea.Cmd.
func (m *PortraitBrowserModel) Update(msg tea.Msg) tea.Cmd {
	switch m.state {
	case pbBrowsing:
		return m.updateBrowsing(msg)
	case pbNaming:
		return m.updateNaming(msg)
	case pbConfirmDelete:
		return m.updateConfirmDelete(msg)
	}
	return nil
}

func (m *PortraitBrowserModel) updateBrowsing(msg tea.Msg) tea.Cmd {
	kmsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch kmsg.String() {
	case "j", "down":
		if m.cursor < len(m.items)-1 {
			m.cursor++
			m.skipToNonHeader(1)
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
			m.skipToNonHeader(-1)
		}
	case "g":
		m.regenerate()
	case "enter":
		if len(m.items) > 0 && !m.items[m.cursor].isHeader && !m.items[m.cursor].isSaved {
			m.state = pbNaming
			m.input.SetValue("")
			m.input.Placeholder = "Character name"
			m.input.Focus()
			m.errMsg = ""
			return m.input.Cursor.BlinkCmd()
		}
	case "d":
		if len(m.items) > 0 && m.items[m.cursor].isSaved {
			m.state = pbConfirmDelete
		}
	}
	return nil
}

func (m *PortraitBrowserModel) updateNaming(msg tea.Msg) tea.Cmd {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "enter":
			name := strings.TrimSpace(m.input.Value())
			if name == "" {
				return nil
			}
			if m.config.FindByName(name) != nil {
				m.errMsg = "Name already taken"
				return nil
			}
			m.config.Add(engine.SavedPortrait{
				Name:   name,
				Params: m.items[m.cursor].params,
			})
			engine.SaveSavedPortraits(m.config)
			m.rebuildItems()
			m.state = pbBrowsing
			m.errMsg = ""
			return nil
		case "esc":
			m.state = pbBrowsing
			m.errMsg = ""
			return nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return cmd
}

func (m *PortraitBrowserModel) updateConfirmDelete(msg tea.Msg) tea.Cmd {
	if kmsg, ok := msg.(tea.KeyMsg); ok {
		switch kmsg.String() {
		case "y":
			if len(m.items) > 0 && m.items[m.cursor].isSaved {
				m.config.Delete(m.items[m.cursor].name)
				engine.SaveSavedPortraits(m.config)
				m.rebuildItems()
			}
			m.state = pbBrowsing
		case "n", "esc":
			m.state = pbBrowsing
		}
	}
	return nil
}

// View renders the portrait browser as a centered overlay.
func (m *PortraitBrowserModel) View(width, height int) string {
	boxW := width - 4
	if boxW < 50 {
		boxW = 50
	}
	boxH := height - 2
	if boxH < 16 {
		boxH = 16
	}

	contentH := boxH - 4 // padding(2) + header(1) + footer(1)

	header := ResultLabelStyle.Render("Portraits")

	// Left side: preview of selected portrait
	var preview string
	if m.cursor >= 0 && m.cursor < len(m.items) && !m.items[m.cursor].isHeader {
		img := engine.RenderPortraitImage(m.items[m.cursor].params)
		preview = PortraitBorderStyle.Render(RenderPortraitArt(img))
	} else {
		preview = RenderEmptyPortraitBox()
	}

	// Right side: scrollable list
	listWidth := boxW - PortraitTotalWidth() - 6 // -6 for padding + gap
	if listWidth < 20 {
		listWidth = 20
	}
	listHeight := contentH - 2

	list := m.viewList(listHeight)

	// Compose left + right
	body := lipgloss.JoinHorizontal(lipgloss.Top, preview, "  ", list)

	// Naming / confirm overlay
	switch m.state {
	case pbNaming:
		body += "\n" + ResultLabelStyle.Render("Save as:") + " " + m.input.View()
		if m.errMsg != "" {
			body += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.errMsg)
		}
	case pbConfirmDelete:
		if m.cursor >= 0 && m.cursor < len(m.items) {
			warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
			body += "\n" + warnStyle.Render(
				fmt.Sprintf("Delete portrait \"%s\"? (y/n)", m.items[m.cursor].name))
		}
	}

	// Footer
	var footerText string
	switch m.state {
	case pbBrowsing:
		footerText = "j/k: navigate | Enter: save | g: regenerate | d: delete | Esc: close"
	case pbNaming:
		footerText = "Enter: confirm | Esc: cancel"
	case pbConfirmDelete:
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

func (m *PortraitBrowserModel) viewList(maxLines int) string {
	if len(m.items) == 0 {
		return DimStyle.Render("Press 'g' to generate portraits.")
	}

	var lines []string
	for i, item := range m.items {
		var line string
		if item.isHeader {
			line = CategoryStyle.Render(item.label)
		} else if i == m.cursor {
			prefix := "▸ " + item.label
			if item.isSaved {
				prefix += "  " + DimStyle.Render(engine.DescribePortrait(item.params))
			}
			line = SrSelectedStyle.Render(prefix)
		} else {
			prefix := "  " + item.label
			if item.isSaved {
				prefix += "  " + DimStyle.Render(engine.DescribePortrait(item.params))
			}
			line = ItemStyle.Render(prefix)
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
