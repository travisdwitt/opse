package engine

import "testing"

func TestFlipCoins_Single(t *testing.T) {
	rng := NewSeededRandomizer(210, 0)
	r := FlipCoins(rng, 1)
	if len(r.Flips) != 1 {
		t.Errorf("expected 1 flip, got %d", len(r.Flips))
	}
	if r.Heads+r.Tails != 1 {
		t.Error("heads + tails should equal 1")
	}
}

func TestFlipCoins_Multi(t *testing.T) {
	rng := NewSeededRandomizer(211, 0)
	r := FlipCoins(rng, 10)
	if len(r.Flips) != 10 {
		t.Errorf("expected 10 flips, got %d", len(r.Flips))
	}
	if r.Heads+r.Tails != 10 {
		t.Errorf("heads(%d) + tails(%d) should equal 10", r.Heads, r.Tails)
	}
}

func TestFlipCoins_ZeroDefaultsToOne(t *testing.T) {
	rng := NewSeededRandomizer(212, 0)
	r := FlipCoins(rng, 0)
	if len(r.Flips) != 1 {
		t.Errorf("expected 1 flip for count=0, got %d", len(r.Flips))
	}
}
