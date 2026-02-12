package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type homeState int

const (
	homeMenu homeState = iota
	homeNaming
	homeBrowsing
	homeConfirmDelete
)

type HomeModel struct {
	state      homeState
	cursor     int
	files      []string
	fileCursor int
	textInput  textinput.Model
	width      int
	height     int

	Choice string
	Title  string
	Path   string
}

func NewHome(files []string) HomeModel {
	ti := textinput.New()
	ti.Placeholder = "My Epic Quest"
	ti.CharLimit = 100
	ti.Width = 30

	return HomeModel{
		state:     homeMenu,
		files:     files,
		textInput: ti,
	}
}

func (m HomeModel) Init() tea.Cmd { return nil }

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case homeMenu:
			return m.updateMenu(msg)
		case homeNaming:
			return m.updateNaming(msg)
		case homeBrowsing:
			return m.updateBrowsing(msg)
		case homeConfirmDelete:
			return m.updateConfirmDelete(msg)
		}
	}
	return m, nil
}

func (m HomeModel) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+q", "q":
		m.Choice = "quit"
		return m, tea.Quit
	case "j", "down":
		if m.cursor < 2 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "enter":
		switch m.cursor {
		case 0:
			m.state = homeNaming
			m.textInput.Focus()
			return m, m.textInput.Cursor.BlinkCmd()
		case 1:
			if len(m.files) > 0 {
				m.state = homeBrowsing
			}
		case 2:
			m.Choice = "quit"
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m HomeModel) updateNaming(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		title := strings.TrimSpace(m.textInput.Value())
		if title == "" {
			title = "New Adventure"
		}
		m.Choice = "new"
		m.Title = title
		return m, tea.Quit
	case "esc":
		m.state = homeMenu
		m.textInput.Reset()
		return m, nil
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m HomeModel) updateBrowsing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.fileCursor < len(m.files)-1 {
			m.fileCursor++
		}
	case "k", "up":
		if m.fileCursor > 0 {
			m.fileCursor--
		}
	case "enter":
		m.Choice = "open"
		m.Path = m.files[m.fileCursor]
		return m, tea.Quit
	case "d":
		if len(m.files) > 0 {
			m.state = homeConfirmDelete
		}
	case "esc":
		m.state = homeMenu
		return m, nil
	case "ctrl+q":
		m.Choice = "quit"
		return m, tea.Quit
	}
	return m, nil
}

func (m HomeModel) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		path := m.files[m.fileCursor]
		os.Remove(path)
		m.files = append(m.files[:m.fileCursor], m.files[m.fileCursor+1:]...)
		if len(m.files) == 0 {
			m.state = homeMenu
		} else {
			if m.fileCursor >= len(m.files) {
				m.fileCursor = len(m.files) - 1
			}
			m.state = homeBrowsing
		}
	case "n", "esc":
		m.state = homeBrowsing
	}
	return m, nil
}

func (m HomeModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string
	switch m.state {
	case homeMenu:
		content = m.viewMenu()
	case homeNaming:
		content = m.viewNaming()
	case homeBrowsing, homeConfirmDelete:
		content = m.viewBrowsing()
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m HomeModel) viewMenu() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3"))

	title := titleStyle.Render("ONE PAGE SOLO ENGINE")
	subtitle := DimStyle.Render("A minimalist toolkit for GM-less RPG adventures.")

	items := []string{"New Adventure", "Open Adventure", "Quit"}
	var menuLines []string
	for i, item := range items {
		if i == 1 && len(m.files) == 0 {
			menuLines = append(menuLines, DimStyle.Render("    "+item+" (none found)"))
			continue
		}
		if i == m.cursor {
			menuLines = append(menuLines, ItemSelectedStyle.Render("  ▸ "+item))
		} else {
			menuLines = append(menuLines, ItemStyle.Render("    "+item))
		}
	}
	menu := strings.Join(menuLines, "\n")

	help := DimStyle.Render("[j/k] Navigate  [Enter] Select  [q] Quit")

	body := fmt.Sprintf("%s\n\n%s\n\n\n%s\n\n\n%s", title, subtitle, menu, help)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(2, 4).
		Render(body)
}

func (m HomeModel) viewNaming() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3"))

	title := titleStyle.Render("NEW ADVENTURE")
	prompt := "Enter a title for your adventure:"

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("252")).
		Padding(0, 1).
		Render(m.textInput.View())

	help := DimStyle.Render("[Enter] Create  [Esc] Back")

	body := fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", title, prompt, inputBox, help)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(2, 4).
		Render(body)
}

func (m HomeModel) viewBrowsing() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3"))

	title := titleStyle.Render("OPEN ADVENTURE")

	maxVisible := 10
	start := 0
	if m.fileCursor >= maxVisible {
		start = m.fileCursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.files) {
		end = len(m.files)
	}

	var fileLines []string
	for i := start; i < end; i++ {
		if i == m.fileCursor {
			fileLines = append(fileLines, ItemSelectedStyle.Render("  ▸ "+m.files[i]))
		} else {
			fileLines = append(fileLines, ItemStyle.Render("    "+m.files[i]))
		}
	}
	fileList := strings.Join(fileLines, "\n")

	counter := DimStyle.Render(fmt.Sprintf("(%d of %d)", m.fileCursor+1, len(m.files)))

	var footer string
	if m.state == homeConfirmDelete {
		warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		footer = warnStyle.Render(fmt.Sprintf("Delete \"%s\"? This cannot be undone. (y/n)", m.files[m.fileCursor]))
	} else {
		footer = DimStyle.Render("[j/k] Navigate  [Enter] Open  [d] Delete  [Esc] Back")
	}

	body := fmt.Sprintf("%s  %s\n\n%s\n\n%s", title, counter, fileList, footer)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(2, 4).
		Render(body)
}
