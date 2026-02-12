package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

type LogViewModel struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func NewLogView() LogViewModel {
	return LogViewModel{}
}

func (l *LogViewModel) SetSize(width, height int) {
	if !l.ready {
		l.viewport = viewport.New(width, height)
		l.setWrappedContent()
		l.ready = true
	} else {
		l.viewport.Width = width
		l.viewport.Height = height
		l.setWrappedContent()
	}
}

func (l *LogViewModel) SetContent(content string) {
	l.content = content
	if l.ready {
		l.setWrappedContent()
	}
}

// setWrappedContent pre-wraps content to the viewport width so that the
// viewport's internal line count matches the actual visual line count.
// Without this, GotoBottom miscalculates when styled lines wrap.
func (l *LogViewModel) setWrappedContent() {
	if l.viewport.Width <= 0 {
		l.viewport.SetContent(l.content)
		return
	}
	lines := strings.Split(l.content, "\n")
	var wrapped []string
	for _, line := range lines {
		w := ansi.StringWidth(line)
		if w > l.viewport.Width {
			wlines := strings.Split(ansi.Wrap(line, l.viewport.Width, ""), "\n")
			wrapped = append(wrapped, wlines...)
		} else {
			wrapped = append(wrapped, line)
		}
	}
	l.viewport.SetContent(strings.Join(wrapped, "\n"))
}

func (l *LogViewModel) ScrollToBottom() {
	if l.ready {
		l.viewport.GotoBottom()
	}
}

func (l *LogViewModel) Update(msg tea.Msg) (*LogViewModel, tea.Cmd) {
	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}

func (l *LogViewModel) View(width int, focused bool) string {
	style := LogStyle
	if focused {
		style = LogStyleFocused
	}
	return style.Width(width).Render(l.viewport.View())
}
