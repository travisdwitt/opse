package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"opse/engine"
	"opse/journal"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FocusArea int

const (
	FocusInput FocusArea = iota
	FocusSidebar
	FocusLog
)

type AppModel struct {
	sidebar         SidebarModel
	logview         LogViewModel
	input           InputModel
	help            HelpModel
	savedRollsModal SavedRollsModel
	focus           FocusArea
	journal         *journal.Journal
	rng             *engine.Randomizer
	deck            *engine.Deck
	utilityDeck     *engine.UtilityDeck
	savedRolls      *engine.SavedRollsConfig
	keys            KeyMap
	width           int
	height          int
	showHelp        bool
	showSavedRolls  bool
	showSaveConfirm bool
	statusMsg       string
	statusExpiry    time.Time
}

func NewApp(j *journal.Journal) AppModel {
	rng := engine.NewRandomizer()
	savedRolls, _ := engine.LoadSavedRolls()
	m := AppModel{
		sidebar:         NewSidebar(),
		logview:         NewLogView(),
		input:           NewInput(),
		help:            NewHelp(),
		savedRollsModal: NewSavedRolls(savedRolls),
		focus:           FocusInput,
		journal:         j,
		rng:             rng,
		deck:            engine.NewDeck(rng),
		utilityDeck:     engine.NewUtilityDeck(rng, false),
		savedRolls:      savedRolls,
		keys:            DefaultKeys,
	}
	m.loadExistingEntries()
	return m
}

func (m AppModel) Init() tea.Cmd { return nil }

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case tea.KeyMsg:
		if m.showSaveConfirm {
			switch msg.String() {
			case "y":
				m.journal.Save()
				m.showSaveConfirm = false
				m.statusMsg = "Saved!"
				m.statusExpiry = time.Now().Add(3 * time.Second)
			case "n", "esc":
				m.showSaveConfirm = false
			}
			return m, nil
		}
		if m.showSavedRolls {
			if key.Matches(msg, m.keys.Escape) && m.savedRollsModal.state == srBrowsing {
				m.showSavedRolls = false
				return m, nil
			}
			rollID, cmd := m.savedRollsModal.Update(msg)
			if rollID != "" {
				m.showSavedRolls = false
				m.runSavedRoll(rollID)
			}
			return m, cmd
		}
		if m.showHelp {
			if key.Matches(msg, m.keys.Escape) || key.Matches(msg, m.keys.Help) {
				m.showHelp = false
				return m, nil
			}
			m.help.Update(msg)
			return m, nil
		}
		if key.Matches(msg, m.keys.Quit) {
			m.journal.Save()
			return m, tea.Quit
		}
		if key.Matches(msg, m.keys.Tab) {
			if m.focus == FocusInput && m.input.autocomplete.visible {
				_, cmd := m.input.Update(msg)
				return m, cmd
			}
			m.cycleFocus()
			return m, nil
		}
		if key.Matches(msg, m.keys.Escape) {
			m.setFocus(FocusInput)
			return m, nil
		}
		if key.Matches(msg, m.keys.Save) {
			if _, err := os.Stat(m.journal.FilePath); err == nil {
				m.showSaveConfirm = true
			} else {
				m.journal.Save()
				m.statusMsg = "Saved!"
				m.statusExpiry = time.Now().Add(3 * time.Second)
			}
			return m, nil
		}
		if key.Matches(msg, m.keys.SavedRolls) {
			m.showSavedRolls = true
			m.savedRollsModal.SetConfig(m.savedRolls)
			return m, nil
		}
		if key.Matches(msg, m.keys.Help) && m.focus != FocusInput {
			m.showHelp = true
			m.help.Reset()
			return m, nil
		}

		if m.focus != FocusInput {
			if action := m.matchShortcut(msg); action != "" {
				m.runAction(action)
				return m, nil
			}
		}

		return m.routeToFocused(msg)

	case NarrativeMsg:
		m.addNarrative(msg.Text)
		return m, nil

	case CommandMsg:
		m.runCommand(msg)
		return m, nil
	}

	return m, nil
}

func (m *AppModel) cycleFocus() {
	switch m.focus {
	case FocusInput:
		m.setFocus(FocusSidebar)
	case FocusSidebar:
		m.setFocus(FocusLog)
	case FocusLog:
		m.setFocus(FocusInput)
	}
}

func (m *AppModel) setFocus(f FocusArea) {
	m.focus = f
	if f == FocusInput {
		m.input.Focus()
	} else {
		m.input.Blur()
	}
}

func (m *AppModel) routeToFocused(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.focus {
	case FocusInput:
		if key.Matches(msg, m.keys.Enter) {
			if submitMsg := m.input.Submit(); submitMsg != nil {
				return m.Update(submitMsg)
			}
			return m, nil
		}
		_, cmd := m.input.Update(msg)
		return m, cmd
	case FocusSidebar:
		if key.Matches(msg, m.keys.Up) {
			m.sidebar.MoveUp()
		} else if key.Matches(msg, m.keys.Down) {
			m.sidebar.MoveDown()
		} else if key.Matches(msg, m.keys.Enter) {
			m.runAction(m.sidebar.Selected().Action)
		}
		return m, nil
	case FocusLog:
		_, cmd := m.logview.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *AppModel) matchShortcut(msg tea.KeyMsg) string {
	shortcuts := []struct {
		binding key.Binding
		action  string
	}{
		{m.keys.OracleLikely, "oracle_likely"},
		{m.keys.OracleEven, "oracle_even"},
		{m.keys.OracleUnlikely, "oracle_unlikely"},
		{m.keys.OracleHow, "oracle_how"},
		{m.keys.ActionFocus, "focus_action"},
		{m.keys.DetailFocus, "focus_detail"},
		{m.keys.TopicFocus, "focus_topic"},
		{m.keys.SetScene, "set_scene"},
		{m.keys.RandomEvent, "random_event"},
		{m.keys.PacingMove, "pacing_move"},
		{m.keys.FailureMove, "failure_move"},
		{m.keys.Generic, "generic"},
	}
	for _, s := range shortcuts {
		if key.Matches(msg, s.binding) {
			return s.action
		}
	}
	return ""
}

func (m *AppModel) runAction(action string) {
	now := time.Now()
	var label, md, tuiStr string
	var entryType journal.EntryType

	switch action {
	case "oracle_likely":
		r := engine.OracleYesNo(m.rng, "Likely")
		label = "Oracle (Yes/No, Likely)"
		md = journal.RenderOracleYesNo(r)
		tuiStr = RenderOracleYesNoTUI(r)
		entryType = journal.EntryOracle
	case "oracle_even":
		r := engine.OracleYesNo(m.rng, "Even")
		label = "Oracle (Yes/No, Even)"
		md = journal.RenderOracleYesNo(r)
		tuiStr = RenderOracleYesNoTUI(r)
		entryType = journal.EntryOracle
	case "oracle_unlikely":
		r := engine.OracleYesNo(m.rng, "Unlikely")
		label = "Oracle (Yes/No, Unlikely)"
		md = journal.RenderOracleYesNo(r)
		tuiStr = RenderOracleYesNoTUI(r)
		entryType = journal.EntryOracle
	case "oracle_how":
		r := engine.OracleHow(m.rng)
		label = "Oracle (How)"
		md = fmt.Sprintf("> **Oracle (How):** %s", r.Result)
		tuiStr = RenderOracleHowTUI(r)
		entryType = journal.EntryOracle
	case "focus_action":
		r := engine.ActionFocus(m.deck)
		label, md = r.TableName, journal.RenderCardTable(r)
		tuiStr = RenderCardTableTUI(r)
		entryType = journal.EntryGenerator
	case "focus_detail":
		r := engine.DetailFocus(m.deck)
		label, md = r.TableName, journal.RenderCardTable(r)
		tuiStr = RenderCardTableTUI(r)
		entryType = journal.EntryGenerator
	case "focus_topic":
		r := engine.TopicFocus(m.deck)
		label, md = r.TableName, journal.RenderCardTable(r)
		tuiStr = RenderCardTableTUI(r)
		entryType = journal.EntryGenerator
	case "random_event":
		r := engine.RandomEvent(m.deck, m.rng)
		label = "Random Event"
		md = journal.RenderRandomEvent(r)
		tuiStr = RenderRandomEventTUI(r)
		entryType = journal.EntryGenerator
	case "set_scene":
		r := engine.SetTheScene(m.rng, m.deck)
		label = "Set the Scene"
		md = journal.RenderSetTheScene(r)
		tuiStr = RenderSetTheSceneTUI(r)
		entryType = journal.EntryScene
	case "pacing_move":
		r := engine.PacingMove(m.rng, m.deck)
		label = "Pacing Move"
		md = journal.RenderPacingMove(r)
		tuiStr = RenderPacingMoveTUI(r)
		entryType = journal.EntryGenerator
	case "failure_move":
		r := engine.FailureMove(m.rng)
		label = "Failure Move"
		md = fmt.Sprintf("> **Failure Move:** %s", r.Result)
		tuiStr = RenderFailureMoveTUI(r)
		entryType = journal.EntryGenerator
	case "generic":
		r := engine.GenericGenerator(m.deck, m.rng)
		label = "Generic Generator"
		md = journal.RenderGeneric(r)
		tuiStr = RenderGenericTUI(r)
		entryType = journal.EntryGenerator
	case "plot_hook":
		r := engine.PlotHook(m.rng)
		label = "Plot Hook"
		md = journal.RenderPlotHook(r)
		tuiStr = RenderPlotHookTUI(r)
		entryType = journal.EntryGenerator
	case "npc":
		r := engine.NPCGenerator(m.deck, m.rng)
		label = "NPC"
		md = journal.RenderNPC(r)
		tuiStr = RenderNPCTUI(r)
		entryType = journal.EntryGenerator
	case "dungeon_theme":
		r := engine.DungeonTheme(m.deck)
		label = "Dungeon Theme"
		md = journal.RenderDungeonTheme(r)
		tuiStr = RenderDungeonThemeTUI(r)
		entryType = journal.EntryGenerator
	case "dungeon_room":
		r := engine.DungeonRoom(m.rng)
		label = "Dungeon Room"
		md = journal.RenderDungeonRoom(r)
		tuiStr = RenderDungeonRoomTUI(r)
		entryType = journal.EntryGenerator
	case "hex":
		r := engine.HexCrawl(m.rng, m.deck)
		label = "Hex"
		md = journal.RenderHex(r)
		tuiStr = RenderHexTUI(r)
		entryType = journal.EntryGenerator
	case "coin_flip":
		r := engine.FlipCoins(m.rng, 1)
		label = "Coin Flip"
		md = journal.RenderCoinFlip(r)
		tuiStr = RenderCoinFlipTUI(r)
		entryType = journal.EntryTool
	case "card_draw":
		r := m.utilityDeck.Draw(1)
		label = "Card Draw"
		md = journal.RenderCardDraw(r)
		tuiStr = RenderCardDrawTUI(r)
		entryType = journal.EntryTool
	case "direction":
		r := engine.RandomDirection(m.rng, 8)
		label = "Direction"
		md = journal.RenderDirection(r)
		tuiStr = RenderDirectionTUI(r)
		entryType = journal.EntryTool
	case "weather":
		r := engine.RandomWeather(m.rng)
		label = "Weather"
		md = journal.RenderWeather(r)
		tuiStr = RenderWeatherTUI(r)
		entryType = journal.EntryTool
	case "color":
		r := engine.RandomColor(m.rng)
		label = "Color"
		md = journal.RenderColor(r)
		tuiStr = RenderColorTUI(r)
		entryType = journal.EntryTool
	case "sound":
		r := engine.RandomSound(m.rng, "")
		label = "Sound"
		md = journal.RenderSound(r)
		tuiStr = RenderSoundTUI(r)
		entryType = journal.EntryTool
	case "dice_roller":
		// When selected from sidebar, do nothing — dice roller needs expression from input
		return
	default:
		return
	}

	m.journal.AddEntry(journal.Entry{
		Timestamp: now, Type: entryType, Label: label, Markdown: md,
	})
	m.refreshLog(tuiStr, now, "Engine")
	m.journal.Save()
}

func (m *AppModel) runSavedRoll(id string) {
	for _, r := range m.savedRolls.Rolls {
		if r.ID == id {
			expr, err := engine.ParseDice(r.Expression)
			if err != nil {
				return
			}
			now := time.Now()
			result := engine.RollDice(m.rng, expr)
			md := journal.RenderDiceRoll(result)
			tuiStr := RenderDiceRollTUI(result)
			m.journal.AddEntry(journal.Entry{
				Timestamp: now, Type: journal.EntryTool, Label: r.Name, Markdown: md,
			})
			m.refreshLog(tuiStr, now, "Engine")
			m.journal.Save()
			return
		}
	}
}

func (m *AppModel) runCommand(cmd CommandMsg) {
	now := time.Now()
	switch cmd.Command {
	case "roll", "r":
		if len(cmd.Args) == 0 {
			return
		}
		expr, err := engine.ParseDice(cmd.Args[0])
		if err != nil {
			return
		}
		r := engine.RollDice(m.rng, expr)
		md := journal.RenderDiceRoll(r)
		tuiStr := RenderDiceRollTUI(r)
		m.journal.AddEntry(journal.Entry{
			Timestamp: now, Type: journal.EntryTool, Label: "Dice", Markdown: md,
		})
		m.refreshLog(tuiStr, now, "Engine")

	case "flip", "f":
		count := 1
		if len(cmd.Args) > 0 {
			fmt.Sscanf(cmd.Args[0], "%d", &count)
		}
		r := engine.FlipCoins(m.rng, count)
		md := journal.RenderCoinFlip(r)
		tuiStr := RenderCoinFlipTUI(r)
		m.journal.AddEntry(journal.Entry{
			Timestamp: now, Type: journal.EntryTool, Label: "Coin Flip", Markdown: md,
		})
		m.refreshLog(tuiStr, now, "Engine")

	case "draw", "card":
		count := 1
		if len(cmd.Args) > 0 {
			fmt.Sscanf(cmd.Args[0], "%d", &count)
		}
		r := m.utilityDeck.Draw(count)
		md := journal.RenderCardDraw(r)
		tuiStr := RenderCardDrawTUI(r)
		m.journal.AddEntry(journal.Entry{
			Timestamp: now, Type: journal.EntryTool, Label: "Card Draw", Markdown: md,
		})
		m.refreshLog(tuiStr, now, "Engine")

	case "shuffle":
		m.utilityDeck.Shuffle()

	case "dir", "direction":
		points := 8
		if len(cmd.Args) > 0 {
			fmt.Sscanf(cmd.Args[0], "%d", &points)
		}
		r := engine.RandomDirection(m.rng, points)
		md := journal.RenderDirection(r)
		tuiStr := RenderDirectionTUI(r)
		m.journal.AddEntry(journal.Entry{
			Timestamp: now, Type: journal.EntryTool, Label: "Direction", Markdown: md,
		})
		m.refreshLog(tuiStr, now, "Engine")

	case "weather", "w":
		m.runAction("weather")
		return

	case "color":
		m.runAction("color")
		return

	case "sound":
		cat := ""
		if len(cmd.Args) > 0 {
			cat = cmd.Args[0]
		}
		r := engine.RandomSound(m.rng, cat)
		md := journal.RenderSound(r)
		tuiStr := RenderSoundTUI(r)
		m.journal.AddEntry(journal.Entry{
			Timestamp: now, Type: journal.EntryTool, Label: "Sound", Markdown: md,
		})
		m.refreshLog(tuiStr, now, "Engine")

	case "scene":
		m.runAction("set_scene")
		return

	case "char":
		if len(cmd.Args) < 2 {
			return
		}
		raw := strings.Join(cmd.Args, " ")
		var name, text string
		if strings.HasPrefix(raw, "\"") {
			end := strings.Index(raw[1:], "\"")
			if end < 0 {
				return
			}
			name = raw[1 : end+1]
			text = strings.TrimSpace(raw[end+2:])
		} else {
			name = cmd.Args[0]
			text = strings.Join(cmd.Args[1:], " ")
		}
		if name == "" || text == "" {
			return
		}
		m.journal.AddEntry(journal.Entry{
			Timestamp: now, Type: journal.EntryNarrative, Label: name, Markdown: text,
		})
		m.refreshLog(text, now, name)
		return

	}
	m.journal.Save()
}

func (m *AppModel) addNarrative(text string) {
	now := time.Now()
	m.journal.AddEntry(journal.Entry{
		Timestamp: now, Type: journal.EntryNarrative, Markdown: text,
	})
	m.refreshLog(text, now, "User")
	m.journal.Save()
}

func (m *AppModel) refreshLog(newTUIEntry string, ts time.Time, source string) {
	current := m.logview.content
	if current != "" {
		current += "\n\n"
	}
	header := FormatEntryHeader(ts, source)
	if header != "" {
		current += header + "\n"
	}
	current += newTUIEntry
	m.logview.SetContent(current)
	m.logview.ScrollToBottom()
}

func (m *AppModel) loadExistingEntries() {
	if len(m.journal.Entries) == 0 {
		return
	}
	var parts []string
	for _, e := range m.journal.Entries {
		source := "Engine"
		if e.Type == journal.EntryNarrative {
			if e.Label != "" {
				source = e.Label
			} else {
				source = "User"
			}
		}
		header := FormatEntryHeader(e.Timestamp, source)
		var entry string
		if header != "" {
			entry = header + "\n"
		}
		if e.Type == journal.EntryNarrative {
			entry += e.Markdown
		} else {
			entry += renderLoadedBlockquote(e.Markdown)
		}
		parts = append(parts, entry)
	}
	m.logview.SetContent(strings.Join(parts, "\n\n"))
}

func renderLoadedBlockquote(md string) string {
	lines := strings.Split(md, "\n")
	var stripped []string
	for _, line := range lines {
		stripped = append(stripped, strings.TrimPrefix(line, "> "))
	}
	text := strings.Join(stripped, "\n")
	// Remove markdown bold formatting for cleaner TUI display
	text = strings.ReplaceAll(text, "**", "")
	return ResultBlockStyle.Render(text)
}

func (m *AppModel) updateLayout() {
	// lipgloss Width/Height include padding but NOT border.
	// Sidebar style: Border(Rounded) + Padding(1,1)
	//   border: 2 horizontal, 2 vertical (added outside Width/Height)
	//   padding: 2 horizontal, 2 vertical (inside Width/Height)
	// Log/Input style: Border(Rounded) + Padding(0,1)
	//   border: 2 horizontal, 2 vertical (added outside Width/Height)
	//   padding: 2 horizontal, 0 vertical (inside Width/Height)
	const (
		sidebarStyleW = 28 // Width param → total rendered = 30
		sidebarBorderW = 2
		sidebarBorderH = 2
		sidebarPadH   = 2 // vertical padding inside Height
		mainBorderW   = 2
		mainBorderH   = 2
		mainPadW      = 2 // horizontal padding inside Width
		inputTextH    = 3
	)

	sidebarTotalW := sidebarStyleW + sidebarBorderW // 30
	mainTotalW := m.width - sidebarTotalW
	mainStyleW := mainTotalW - mainBorderW // Width param for log/input
	mainTextW := mainStyleW - mainPadW     // text area inside padding
	if mainTextW < 10 {
		mainTextW = 10
	}

	bodyH := m.height - 2                  // title (1) + help bar (1)
	inputTotalH := inputTextH + mainBorderH // 5
	logTotalH := bodyH - inputTotalH
	logViewportH := logTotalH - mainBorderH // viewport content height
	if logViewportH < 3 {
		logViewportH = 3
	}

	m.logview.SetSize(mainTextW, logViewportH)
	m.input.textarea.SetWidth(mainTextW)
	m.input.textarea.SetHeight(inputTextH)
}

func (m AppModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.showSavedRolls {
		return m.savedRollsModal.View(m.width, m.height)
	}
	if m.showHelp {
		return m.help.View(m.width, m.height)
	}

	const (
		sidebarStyleW  = 28
		sidebarBorderW = 2
		sidebarBorderH = 2
		mainBorderW    = 2
	)

	sidebarTotalW := sidebarStyleW + sidebarBorderW
	mainTotalW := m.width - sidebarTotalW
	mainStyleW := mainTotalW - mainBorderW
	bodyH := m.height - 2
	sidebarH := bodyH - sidebarBorderH

	title := TitleStyle.Render(fmt.Sprintf("OPSE — %s", m.journal.Title))

	// If autocomplete is visible, shrink log to make room for popup
	if m.input.autocomplete.visible {
		acHeight := len(m.input.autocomplete.suggestions) + 2 // items + border
		adjusted := m.logview.viewport.Height - acHeight
		if adjusted < 3 {
			adjusted = 3
		}
		m.logview.SetSize(m.logview.viewport.Width, adjusted)
	}

	sidebar := m.sidebar.View(sidebarStyleW, sidebarH, m.focus == FocusSidebar)
	logView := m.logview.View(mainStyleW, m.focus == FocusLog)
	inputView := m.input.View(mainStyleW, m.focus == FocusInput)
	mainPanel := lipgloss.JoinVertical(lipgloss.Left, logView, inputView)

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainPanel)

	var helpBar string
	if m.showSaveConfirm {
		warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
		helpBar = warnStyle.Render(fmt.Sprintf("Overwrite \"%s\"? (y/n)", m.journal.FilePath))
	} else if m.statusMsg != "" && time.Now().Before(m.statusExpiry) {
		helpBar = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render(m.statusMsg)
	} else {
		m.statusMsg = ""
		helpBar = HelpStyle.Render("Tab: switch | 1-9: generators | /: commands | Ctrl+R: saved rolls | ?: help | Ctrl+Q: quit")
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, body, helpBar)
}
