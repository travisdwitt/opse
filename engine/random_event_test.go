package engine

import "testing"

func TestRandomEvent_ValidResults(t *testing.T) {
	rng := NewSeededRandomizer(90, 0)
	deck := NewDeck(rng)
	for range 50 {
		r := RandomEvent(deck, rng)
		if r.Action.TableName != "Action Focus" {
			t.Errorf("action table name = %q", r.Action.TableName)
		}
		if r.Topic.TableName != "Topic Focus" {
			t.Errorf("topic table name = %q", r.Topic.TableName)
		}
		if r.Action.Entry == "" {
			t.Error("action entry empty")
		}
		if r.Topic.Entry == "" {
			t.Error("topic entry empty")
		}
	}
}
