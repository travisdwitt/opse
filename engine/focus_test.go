package engine

import "testing"

func TestActionFocusAllRanks(t *testing.T) {
	for r := RankTwo; r <= RankAce; r++ {
		if _, ok := actionFocusTable[r]; !ok {
			t.Errorf("actionFocusTable missing rank %d", r)
		}
	}
}

func TestDetailFocusAllRanks(t *testing.T) {
	for r := RankTwo; r <= RankAce; r++ {
		if _, ok := detailFocusTable[r]; !ok {
			t.Errorf("detailFocusTable missing rank %d", r)
		}
	}
}

func TestTopicFocusAllRanks(t *testing.T) {
	for r := RankTwo; r <= RankAce; r++ {
		if _, ok := topicFocusTable[r]; !ok {
			t.Errorf("topicFocusTable missing rank %d", r)
		}
	}
}

func TestActionFocusResult(t *testing.T) {
	rng := NewSeededRandomizer(60, 0)
	deck := NewDeck(rng)
	r := ActionFocus(deck)
	if r.TableName != "Action Focus" {
		t.Errorf("expected table name 'Action Focus', got %q", r.TableName)
	}
	if r.Entry == "" {
		t.Error("ActionFocus returned empty entry")
	}
}

func TestDetailFocusResult(t *testing.T) {
	rng := NewSeededRandomizer(61, 0)
	deck := NewDeck(rng)
	r := DetailFocus(deck)
	if r.TableName != "Detail Focus" {
		t.Errorf("expected table name 'Detail Focus', got %q", r.TableName)
	}
	if r.Entry == "" {
		t.Error("DetailFocus returned empty entry")
	}
}

func TestTopicFocusResult(t *testing.T) {
	rng := NewSeededRandomizer(62, 0)
	deck := NewDeck(rng)
	r := TopicFocus(deck)
	if r.TableName != "Topic Focus" {
		t.Errorf("expected table name 'Topic Focus', got %q", r.TableName)
	}
	if r.Entry == "" {
		t.Error("TopicFocus returned empty entry")
	}
}
