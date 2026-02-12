package engine

import "testing"

func TestPacingMove_Roll6HasEvent(t *testing.T) {
	rng := NewSeededRandomizer(80, 0)
	found := false
	for range 500 {
		deck := NewDeck(rng)
		r := PacingMove(rng, deck)
		if r.Roll == 6 {
			found = true
			if r.RandomEvent == nil {
				t.Error("pacing move roll 6 should have RandomEvent")
			}
			if r.RandomEvent.Action.Entry == "" || r.RandomEvent.Topic.Entry == "" {
				t.Error("random event entries should not be empty")
			}
		}
		if r.Roll != 6 && r.RandomEvent != nil {
			t.Errorf("pacing move roll %d should not have RandomEvent", r.Roll)
		}
	}
	if !found {
		t.Error("never rolled a 6 for pacing move")
	}
}

func TestFailureMove_AllEntries(t *testing.T) {
	rng := NewSeededRandomizer(81, 0)
	seen := make(map[string]bool)
	for range 1000 {
		r := FailureMove(rng)
		if r.Result == "" {
			t.Error("empty result")
		}
		seen[r.Result] = true
	}
	expected := []string{
		"Cause Harm", "Put Someone in a Spot", "Offer a Choice",
		"Advance a Threat", "Reveal an Unwelcome Truth", "Foreshadow Trouble",
	}
	for _, e := range expected {
		if !seen[e] {
			t.Errorf("never saw %q", e)
		}
	}
}
