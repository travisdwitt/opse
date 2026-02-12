package journal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewJournal(t *testing.T) {
	j := New("Test Adventure", "/tmp/test.md")
	if j.Title != "Test Adventure" {
		t.Errorf("expected title 'Test Adventure', got %q", j.Title)
	}
	if j.IsDirty() {
		t.Error("new journal should not be dirty")
	}
}

func TestAddEntryMarksDirty(t *testing.T) {
	j := New("Test", "/tmp/test.md")
	j.AddEntry(Entry{Type: EntryNarrative, Markdown: "Hello"})
	if !j.IsDirty() {
		t.Error("journal should be dirty after AddEntry")
	}
}

func TestSaveClearsDirty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	j := New("Test", path)
	j.AddEntry(Entry{Type: EntryNarrative, Markdown: "Hello"})
	if err := j.Save(); err != nil {
		t.Fatal(err)
	}
	if j.IsDirty() {
		t.Error("journal should not be dirty after Save")
	}
}

func TestSaveSkipsWhenClean(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	j := New("Test", path)
	// Save without adding anything
	if err := j.Save(); err != nil {
		t.Fatal(err)
	}
	// File should not exist since nothing was dirty
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should not exist when journal was never dirty")
	}
}

func TestRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "adventure.md")

	j := New("Test Adventure", path)
	j.AddEntry(Entry{Type: EntryNarrative, Markdown: "The hero entered the cave."})
	j.AddEntry(Entry{Type: EntryOracle, Markdown: "> **Oracle (Yes/No, Even):** Yes"})
	if err := j.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Title != "Test Adventure" {
		t.Errorf("loaded title = %q, want 'Test Adventure'", loaded.Title)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("loaded %d entries, want 2", len(loaded.Entries))
	}
	if loaded.Entries[0].Type != EntryNarrative {
		t.Errorf("entry[0] type = %q, want %q", loaded.Entries[0].Type, EntryNarrative)
	}
	if loaded.Entries[0].Markdown != "The hero entered the cave." {
		t.Errorf("entry[0] markdown = %q, want 'The hero entered the cave.'", loaded.Entries[0].Markdown)
	}
	if loaded.Entries[1].Type != EntryOracle {
		t.Errorf("entry[1] type = %q, want %q", loaded.Entries[1].Type, EntryOracle)
	}
}

func TestRoundTripCharacterVoice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "adventure.md")

	j := New("Voice Test", path)
	j.AddEntry(Entry{Type: EntryNarrative, Label: "Elara", Markdown: "I search the room."})
	j.AddEntry(Entry{Type: EntryNarrative, Markdown: "The room is dark and cold."})
	if err := j.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("loaded %d entries, want 2", len(loaded.Entries))
	}
	if loaded.Entries[0].Markdown != "I search the room." {
		t.Errorf("entry[0] markdown = %q", loaded.Entries[0].Markdown)
	}
	if loaded.Entries[1].Markdown != "The room is dark and cold." {
		t.Errorf("entry[1] markdown = %q", loaded.Entries[1].Markdown)
	}
}
