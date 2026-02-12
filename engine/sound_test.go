package engine

import "testing"

func TestRandomSound_Any(t *testing.T) {
	rng := NewSeededRandomizer(240, 0)
	r := RandomSound(rng, "")
	if r.Sound == "" {
		t.Error("sound should not be empty")
	}
	if r.Category == "" {
		t.Error("category should not be empty")
	}
}

func TestRandomSound_SpecificCategory(t *testing.T) {
	rng := NewSeededRandomizer(241, 0)
	r := RandomSound(rng, "combat")
	if r.Category != "combat" {
		t.Errorf("expected combat category, got %q", r.Category)
	}
	if r.Sound == "" {
		t.Error("sound should not be empty")
	}
}

func TestRandomSound_UnknownCategoryFallback(t *testing.T) {
	rng := NewSeededRandomizer(242, 0)
	r := RandomSound(rng, "nonexistent")
	if r.Sound == "" {
		t.Error("should fall back to any category")
	}
}

func TestSoundTotalEntries(t *testing.T) {
	if len(allSoundsFlat) < 200 {
		t.Errorf("expected 200+ sounds, got %d", len(allSoundsFlat))
	}
}

func TestSoundCategoryNames(t *testing.T) {
	names := SoundCategoryNames()
	if len(names) != 10 {
		t.Errorf("expected 10 categories, got %d", len(names))
	}
}
