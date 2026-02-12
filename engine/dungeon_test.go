package engine

import "testing"

func TestDungeonTheme(t *testing.T) {
	rng := NewSeededRandomizer(110, 0)
	deck := NewDeck(rng)
	r := DungeonTheme(deck)
	if r.Looks.Entry == "" {
		t.Error("looks empty")
	}
	if r.Used.Entry == "" {
		t.Error("used empty")
	}
}

func TestDungeonRoom_GroupedRanges(t *testing.T) {
	rng := NewSeededRandomizer(111, 0)
	for range 500 {
		r := DungeonRoom(rng)
		if r.Location == "" {
			t.Error("location empty")
		}
		if r.Encounter == "" {
			t.Error("encounter empty")
		}
		if r.Object == "" {
			t.Error("object empty")
		}
		if r.Exits == "" {
			t.Error("exits empty")
		}
	}
}

func TestDungeonEncounterGrouped(t *testing.T) {
	if dungeonEncounter(1) != "None" || dungeonEncounter(2) != "None" {
		t.Error("1-2 should be None")
	}
	if dungeonEncounter(3) != "Hostile enemies" || dungeonEncounter(4) != "Hostile enemies" {
		t.Error("3-4 should be Hostile enemies")
	}
	if dungeonEncounter(5) != "An obstacle blocks the way" {
		t.Error("5 should be obstacle")
	}
	if dungeonEncounter(6) != "Unique NPC or adversary" {
		t.Error("6 should be unique NPC")
	}
}
