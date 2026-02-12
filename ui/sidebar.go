package ui

import (
	"fmt"
	"strings"
)

type SidebarItem struct {
	Label    string
	Shortcut string
	Action   string
}

type SidebarCategory struct {
	Name  string
	Items []SidebarItem
}

type SidebarModel struct {
	Categories []SidebarCategory
	cursor     int
	height     int
	scrollOff  int
}

func NewSidebar() SidebarModel {
	return SidebarModel{
		Categories: []SidebarCategory{
			{Name: "ORACLE", Items: []SidebarItem{
				{"Yes/No (Likely)", "1", "oracle_likely"},
				{"Yes/No (Even)", "2", "oracle_even"},
				{"Yes/No (Unlikely)", "3", "oracle_unlikely"},
				{"How", "4", "oracle_how"},
			}},
			{Name: "FOCUS", Items: []SidebarItem{
				{"Action", "5", "focus_action"},
				{"Detail", "6", "focus_detail"},
				{"Topic", "7", "focus_topic"},
			}},
			{Name: "SCENE", Items: []SidebarItem{
				{"Set the Scene", "8", "set_scene"},
				{"Random Event", "9", "random_event"},
			}},
			{Name: "GM MOVES", Items: []SidebarItem{
				{"Pacing Move", "0", "pacing_move"},
				{"Failure Move", "-", "failure_move"},
			}},
			{Name: "GENERATORS", Items: []SidebarItem{
				{"Generic", "=", "generic"},
				{"Plot Hook", "", "plot_hook"},
				{"NPC", "", "npc"},
				{"Dungeon Theme", "", "dungeon_theme"},
				{"Dungeon Room", "", "dungeon_room"},
				{"Hex", "", "hex"},
			}},
			{Name: "TOOLS", Items: []SidebarItem{
				{"Dice Roller", "/", "dice_roller"},
				{"Coin Flip", "", "coin_flip"},
				{"Card Draw", "", "card_draw"},
				{"Direction", "", "direction"},
				{"Weather", "", "weather"},
				{"Color", "", "color"},
				{"Sound", "", "sound"},
			}},
		},
	}
}

func (s *SidebarModel) flatItems() []SidebarItem {
	var items []SidebarItem
	for _, cat := range s.Categories {
		items = append(items, cat.Items...)
	}
	return items
}

func (s *SidebarModel) totalItems() int {
	n := 0
	for _, cat := range s.Categories {
		n += len(cat.Items)
	}
	return n
}

func (s *SidebarModel) MoveUp() {
	if s.cursor > 0 {
		s.cursor--
	}
}

func (s *SidebarModel) MoveDown() {
	if s.cursor < s.totalItems()-1 {
		s.cursor++
	}
}

func (s *SidebarModel) Selected() SidebarItem {
	return s.flatItems()[s.cursor]
}



func (s *SidebarModel) View(width, height int, focused bool) string {
	s.height = height
	var b strings.Builder
	idx := 0

	// Calculate visible lines
	lines := make([]string, 0)
	for _, cat := range s.Categories {
		lines = append(lines, CategoryStyle.Render(cat.Name))
		for _, item := range cat.Items {
			shortcut := "   "
			if item.Shortcut != "" {
				shortcut = ShortcutStyle.Render(fmt.Sprintf("[%s]", item.Shortcut))
			}
			label := item.Label
			if idx == s.cursor {
				label = ItemSelectedStyle.Render("▸ " + label)
			} else {
				label = ItemStyle.Render("  " + label)
			}
			lines = append(lines, fmt.Sprintf(" %s %s", shortcut, label))
			idx++
		}
		lines = append(lines, "")
	}

	// Scroll handling: height includes padding (2) but not border
	contentHeight := height - 2
	if contentHeight < 1 {
		contentHeight = 1
	}
	if s.scrollOff > len(lines)-contentHeight {
		s.scrollOff = len(lines) - contentHeight
	}
	if s.scrollOff < 0 {
		s.scrollOff = 0
	}

	end := s.scrollOff + contentHeight
	if end > len(lines) {
		end = len(lines)
	}

	for _, line := range lines[s.scrollOff:end] {
		b.WriteString(line)
		b.WriteString("\n")
	}

	style := SidebarStyle
	if focused {
		style = SidebarFocusedStyle
	}
	return style.Width(width).Height(height).Render(b.String())
}
