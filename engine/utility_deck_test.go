package engine

import "testing"

func TestUtilityDeck_NoJokers(t *testing.T) {
	rng := NewSeededRandomizer(220, 0)
	d := NewUtilityDeck(rng, false)
	if d.Remaining() != 52 {
		t.Errorf("expected 52, got %d", d.Remaining())
	}
}

func TestUtilityDeck_WithJokers(t *testing.T) {
	rng := NewSeededRandomizer(221, 0)
	d := NewUtilityDeck(rng, true)
	if d.Remaining() != 54 {
		t.Errorf("expected 54, got %d", d.Remaining())
	}
}

func TestUtilityDeck_DrawReducesRemaining(t *testing.T) {
	rng := NewSeededRandomizer(222, 0)
	d := NewUtilityDeck(rng, false)
	r := d.Draw(5)
	if len(r.Cards) != 5 {
		t.Errorf("expected 5 cards, got %d", len(r.Cards))
	}
	if d.Remaining() != 47 {
		t.Errorf("expected 47 remaining, got %d", d.Remaining())
	}
}

func TestUtilityDeck_DrawMoreThanRemaining(t *testing.T) {
	rng := NewSeededRandomizer(223, 0)
	d := NewUtilityDeck(rng, false)
	for range 10 {
		d.Draw(5)
	}
	// 52 - 50 = 2 remaining
	r := d.Draw(5)
	if len(r.Cards) != 2 {
		t.Errorf("expected 2 cards, got %d", len(r.Cards))
	}
}
