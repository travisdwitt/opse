package engine

import "testing"

func TestHexCrawl_ContentsFeature(t *testing.T) {
	rng := NewSeededRandomizer(120, 0)
	foundFeature := false
	for range 500 {
		r := HexCrawl(rng, NewDeck(rng))
		if r.ContentsRoll == 6 {
			foundFeature = true
			if r.Feature == "" {
				t.Error("contents 6 should populate feature")
			}
		}
		if r.ContentsRoll != 6 && r.Feature != "" {
			t.Error("contents != 6 should not have feature")
		}
	}
	if !foundFeature {
		t.Error("never rolled contents 6")
	}
}

func TestHexCrawl_EventRandomEvent(t *testing.T) {
	rng := NewSeededRandomizer(121, 0)
	foundEvent := false
	for range 500 {
		r := HexCrawl(rng, NewDeck(rng))
		if r.EventRoll >= 5 {
			foundEvent = true
			if r.RandomEvent == nil {
				t.Error("event 5+ should have RandomEvent")
			}
		}
		if r.EventRoll < 5 && r.RandomEvent != nil {
			t.Error("event <5 should not have RandomEvent")
		}
	}
	if !foundEvent {
		t.Error("never rolled event 5+")
	}
}
