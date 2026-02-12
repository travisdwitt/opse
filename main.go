package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"opse/journal"
	"opse/ui"
)

func main() {
	if len(os.Args) > 1 {
		path := os.Args[1]
		loaded, err := journal.Load(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", path, err)
			os.Exit(1)
		}
		runApp(loaded)
		return
	}

	files, _ := journal.ListJournals(".")
	home := ui.NewHome(files)
	p := tea.NewProgram(home, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	h := result.(ui.HomeModel)
	switch h.Choice {
	case "quit", "":
		return
	case "new":
		name := sanitizeFilename(h.Title)
		path := filepath.Join(".", name+".md")
		j := journal.New(h.Title, path)
		runApp(j)
	case "open":
		loaded, err := journal.Load(h.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading %s: %v\n", h.Path, err)
			os.Exit(1)
		}
		runApp(loaded)
	}
}

func runApp(j *journal.Journal) {
	app := ui.NewApp(j)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func sanitizeFilename(title string) string {
	name := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			return r
		}
		if r == ' ' {
			return '_'
		}
		return -1
	}, title)
	if name == "" {
		name = "adventure"
	}
	return time.Now().Format("2006-01-02") + "_" + strings.ToLower(name)
}
