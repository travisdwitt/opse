package engine

import "testing"

func TestSceneComplicationValidResult(t *testing.T) {
	rng := NewSeededRandomizer(70, 0)
	for range 100 {
		r := SceneComplication(rng)
		if r.Roll < 1 || r.Roll > 6 {
			t.Errorf("roll %d out of range", r.Roll)
		}
		if r.Result == "" {
			t.Error("empty result")
		}
	}
}

func TestSetTheScene_NonAltered(t *testing.T) {
	rng := NewSeededRandomizer(70, 0)
	found := false
	for range 100 {
		r := SetTheScene(rng, NewDeck(rng))
		if !r.Altered {
			found = true
			if r.AlteredScene != nil {
				t.Error("non-altered scene should have nil AlteredScene")
			}
		}
	}
	if !found {
		t.Error("never got a non-altered scene")
	}
}

func TestSetTheScene_Altered(t *testing.T) {
	rng := NewSeededRandomizer(71, 0)
	found := false
	for range 200 {
		r := SetTheScene(rng, NewDeck(rng))
		if r.Altered {
			found = true
			if r.AlteredScene == nil {
				t.Error("altered scene should have non-nil AlteredScene")
			}
		}
	}
	if !found {
		t.Error("never got an altered scene in 200 tries")
	}
}

func TestSetTheScene_CascadeToComplication(t *testing.T) {
	rng := NewSeededRandomizer(72, 0)
	for range 500 {
		r := SetTheScene(rng, NewDeck(rng))
		if r.Altered && r.AlteredScene != nil && r.AlteredScene.Roll == 4 {
			if r.AlteredScene.Cascade == nil {
				t.Error("altered scene roll 4 should cascade to complication")
			}
			if _, ok := r.AlteredScene.Cascade.(*SceneComplicationResult); !ok {
				t.Error("altered scene roll 4 cascade should be SceneComplicationResult")
			}
			return
		}
	}
}

func TestSetTheScene_CascadeToPacingMove(t *testing.T) {
	rng := NewSeededRandomizer(73, 0)
	for range 500 {
		r := SetTheScene(rng, NewDeck(rng))
		if r.Altered && r.AlteredScene != nil && r.AlteredScene.Roll == 5 {
			if r.AlteredScene.Cascade == nil {
				t.Error("altered scene roll 5 should cascade to pacing move")
			}
			if _, ok := r.AlteredScene.Cascade.(*PacingMoveResult); !ok {
				t.Error("altered scene roll 5 cascade should be PacingMoveResult")
			}
			return
		}
	}
}
