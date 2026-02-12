package engine

import "testing"

func TestDirection_4Point(t *testing.T) {
	rng := NewSeededRandomizer(230, 0)
	valid := map[string]bool{"North": true, "East": true, "South": true, "West": true}
	for range 100 {
		r := RandomDirection(rng, 4)
		if !valid[r.Direction] {
			t.Errorf("unexpected 4-point direction: %s", r.Direction)
		}
		if r.Arrow == "" {
			t.Error("arrow should not be empty")
		}
	}
}

func TestDirection_8Point(t *testing.T) {
	rng := NewSeededRandomizer(231, 0)
	seen := make(map[string]bool)
	for range 1000 {
		r := RandomDirection(rng, 8)
		seen[r.Direction] = true
	}
	if len(seen) != 8 {
		t.Errorf("expected 8 unique directions, got %d", len(seen))
	}
}

func TestDirection_16Point(t *testing.T) {
	rng := NewSeededRandomizer(232, 0)
	seen := make(map[string]bool)
	for range 5000 {
		r := RandomDirection(rng, 16)
		seen[r.Direction] = true
	}
	if len(seen) != 16 {
		t.Errorf("expected 16 unique directions, got %d", len(seen))
	}
}
