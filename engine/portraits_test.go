package engine

import "testing"

func TestSavedPortraitsConfig_Add(t *testing.T) {
	cfg := &SavedPortraitsConfig{}
	cfg.Add(SavedPortrait{Name: "Elara", Params: PortraitParams{FaceShape: 1, HairStyle: 3}})
	if len(cfg.Portraits) != 1 {
		t.Fatalf("expected 1 portrait, got %d", len(cfg.Portraits))
	}
	if cfg.Portraits[0].Name != "Elara" {
		t.Error("expected name Elara")
	}
}

func TestSavedPortraitsConfig_Delete(t *testing.T) {
	cfg := &SavedPortraitsConfig{}
	cfg.Add(SavedPortrait{Name: "Elara"})
	cfg.Add(SavedPortrait{Name: "Thorin"})
	cfg.Delete("elara") // case-insensitive
	if len(cfg.Portraits) != 1 {
		t.Fatalf("expected 1 portrait after delete, got %d", len(cfg.Portraits))
	}
	if cfg.Portraits[0].Name != "Thorin" {
		t.Errorf("expected Thorin to remain, got %s", cfg.Portraits[0].Name)
	}
}

func TestSavedPortraitsConfig_DeleteNonExistent(t *testing.T) {
	cfg := &SavedPortraitsConfig{}
	cfg.Add(SavedPortrait{Name: "Elara"})
	cfg.Delete("nobody")
	if len(cfg.Portraits) != 1 {
		t.Error("delete of non-existent should not remove anything")
	}
}

func TestSavedPortraitsConfig_FindByName(t *testing.T) {
	cfg := &SavedPortraitsConfig{}
	cfg.Add(SavedPortrait{Name: "Elara", Params: PortraitParams{FaceShape: 1}})
	cfg.Add(SavedPortrait{Name: "Thorin", Params: PortraitParams{FaceShape: 2}})

	// Exact match
	p := cfg.FindByName("Elara")
	if p == nil || p.Params.FaceShape != 1 {
		t.Error("expected to find Elara with face shape 1")
	}

	// Case-insensitive
	p = cfg.FindByName("THORIN")
	if p == nil || p.Params.FaceShape != 2 {
		t.Error("expected to find Thorin case-insensitively")
	}

	// Not found
	if cfg.FindByName("Nobody") != nil {
		t.Error("expected nil for unknown name")
	}
}

func TestLoadSavedPortraits_NonExistent(t *testing.T) {
	cfg := &SavedPortraitsConfig{}
	if len(cfg.Portraits) != 0 {
		t.Error("new config should have empty portraits")
	}
}
