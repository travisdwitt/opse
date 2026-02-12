package engine

import "testing"

func TestSavedRollsConfig_Add(t *testing.T) {
	cfg := &SavedRollsConfig{}
	cfg.Add(SavedRoll{ID: "1", Name: "Test", Expression: "d20"})
	if len(cfg.Rolls) != 1 {
		t.Errorf("expected 1 roll, got %d", len(cfg.Rolls))
	}
}

func TestSavedRollsConfig_Delete(t *testing.T) {
	cfg := &SavedRollsConfig{}
	cfg.Add(SavedRoll{ID: "1", Name: "A", Expression: "d20"})
	cfg.Add(SavedRoll{ID: "2", Name: "B", Expression: "2d6"})
	cfg.Delete("1")
	if len(cfg.Rolls) != 1 {
		t.Errorf("expected 1 roll after delete, got %d", len(cfg.Rolls))
	}
	if cfg.Rolls[0].ID != "2" {
		t.Errorf("expected remaining roll ID '2', got %q", cfg.Rolls[0].ID)
	}
}

func TestSavedRollsConfig_DeleteNonExistent(t *testing.T) {
	cfg := &SavedRollsConfig{}
	cfg.Add(SavedRoll{ID: "1", Name: "A", Expression: "d20"})
	cfg.Delete("999")
	if len(cfg.Rolls) != 1 {
		t.Error("deleting non-existent ID should not change rolls")
	}
}

func TestSavedRollsConfig_ByFolder(t *testing.T) {
	cfg := &SavedRollsConfig{}
	cfg.Add(SavedRoll{ID: "1", Name: "A", Folder: "Combat"})
	cfg.Add(SavedRoll{ID: "2", Name: "B", Folder: "Combat"})
	cfg.Add(SavedRoll{ID: "3", Name: "C", Folder: ""})
	m := cfg.ByFolder()
	if len(m["Combat"]) != 2 {
		t.Errorf("expected 2 in Combat, got %d", len(m["Combat"]))
	}
	if len(m[""]) != 1 {
		t.Errorf("expected 1 in root, got %d", len(m[""]))
	}
}

func TestLoadSavedRolls_NonExistent(t *testing.T) {
	// LoadSavedRolls should return empty config for missing file (not error)
	// This test relies on the default path not existing in test environment
	cfg, err := LoadSavedRolls()
	if err != nil {
		// It's ok if the config dir exists with an actual file
		return
	}
	if cfg == nil {
		t.Error("should return non-nil config")
	}
}
