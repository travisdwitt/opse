package journal

import (
	"os"
	"time"
)

type EntryType string

const (
	EntryNarrative EntryType = "narrative"
	EntryGenerator EntryType = "generator"
	EntryOracle    EntryType = "oracle"
	EntryTool      EntryType = "tool"
	EntryScene     EntryType = "scene"
)

type Entry struct {
	Timestamp time.Time
	Type      EntryType
	Label     string
	Markdown  string
}

type Journal struct {
	Title     string
	CreatedAt time.Time
	Entries   []Entry
	FilePath  string
	dirty     bool
}

func New(title, filePath string) *Journal {
	return &Journal{
		Title:     title,
		CreatedAt: time.Now(),
		FilePath:  filePath,
	}
}

func (j *Journal) AddEntry(e Entry) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	j.Entries = append(j.Entries, e)
	j.dirty = true
}

func (j *Journal) Save() error {
	if !j.dirty {
		return nil
	}
	md := Render(j)
	err := os.WriteFile(j.FilePath, []byte(md), 0644)
	if err == nil {
		j.dirty = false
	}
	return err
}

func (j *Journal) IsDirty() bool {
	return j.dirty
}
