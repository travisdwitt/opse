package engine

import "testing"

func TestGenericGenerator(t *testing.T) {
	rng := NewSeededRandomizer(100, 0)
	deck := NewDeck(rng)
	r := GenericGenerator(deck, rng)
	if r.Action.Entry == "" {
		t.Error("action empty")
	}
	if r.Detail.Entry == "" {
		t.Error("detail empty")
	}
	if r.Significance.Result == "" {
		t.Error("significance empty")
	}
}

func TestPlotHook(t *testing.T) {
	rng := NewSeededRandomizer(101, 0)
	for range 100 {
		r := PlotHook(rng)
		if r.Objective == "" {
			t.Error("objective empty")
		}
		if r.Adversary == "" {
			t.Error("adversary empty")
		}
		if r.Reward == "" {
			t.Error("reward empty")
		}
	}
}

func TestNPCGenerator(t *testing.T) {
	rng := NewSeededRandomizer(102, 0)
	deck := NewDeck(rng)
	r := NPCGenerator(deck, rng)
	if r.Identity.Entry == "" {
		t.Error("identity empty")
	}
	if r.Goal.Entry == "" {
		t.Error("goal empty")
	}
	if r.Feature == "" {
		t.Error("feature empty")
	}
	if r.FeatureDetail.Entry == "" {
		t.Error("feature detail empty")
	}
	if r.Attitude.Result == "" {
		t.Error("attitude empty")
	}
	if r.Topic.Entry == "" {
		t.Error("topic empty")
	}
	// Verify suit domains are preserved
	if r.Identity.Draw.Card.Domain() == "" {
		t.Error("identity should have suit domain")
	}
	if r.Goal.Draw.Card.Domain() == "" {
		t.Error("goal should have suit domain")
	}
}
