package engine

import "testing"

func TestOracleYesNo_Likely(t *testing.T) {
	rng := NewSeededRandomizer(42, 0)
	for range 100 {
		r := OracleYesNo(rng, "Likely")
		if r.AnswerRoll >= 3 && !r.Answer {
			t.Errorf("Likely: roll %d should be Yes", r.AnswerRoll)
		}
		if r.AnswerRoll < 3 && r.Answer {
			t.Errorf("Likely: roll %d should be No", r.AnswerRoll)
		}
	}
}

func TestOracleYesNo_Even(t *testing.T) {
	rng := NewSeededRandomizer(43, 0)
	for range 100 {
		r := OracleYesNo(rng, "Even")
		if r.AnswerRoll >= 4 && !r.Answer {
			t.Errorf("Even: roll %d should be Yes", r.AnswerRoll)
		}
		if r.AnswerRoll < 4 && r.Answer {
			t.Errorf("Even: roll %d should be No", r.AnswerRoll)
		}
	}
}

func TestOracleYesNo_Unlikely(t *testing.T) {
	rng := NewSeededRandomizer(44, 0)
	for range 100 {
		r := OracleYesNo(rng, "Unlikely")
		if r.AnswerRoll >= 5 && !r.Answer {
			t.Errorf("Unlikely: roll %d should be Yes", r.AnswerRoll)
		}
		if r.AnswerRoll < 5 && r.Answer {
			t.Errorf("Unlikely: roll %d should be No", r.AnswerRoll)
		}
	}
}

func TestOracleYesNo_Modifier(t *testing.T) {
	rng := NewSeededRandomizer(45, 0)
	sawBut, sawAnd, sawNone := false, false, false
	for range 1000 {
		r := OracleYesNo(rng, "Even")
		switch r.ModRoll {
		case 1:
			if r.Modifier != "but..." {
				t.Errorf("mod roll 1 should be 'but...', got %q", r.Modifier)
			}
			sawBut = true
		case 6:
			if r.Modifier != "and..." {
				t.Errorf("mod roll 6 should be 'and...', got %q", r.Modifier)
			}
			sawAnd = true
		default:
			if r.Modifier != "" {
				t.Errorf("mod roll %d should be empty, got %q", r.ModRoll, r.Modifier)
			}
			sawNone = true
		}
	}
	if !sawBut || !sawAnd || !sawNone {
		t.Error("didn't see all modifier outcomes")
	}
}

func TestOracleHow_AllValues(t *testing.T) {
	rng := NewSeededRandomizer(50, 0)
	seen := make(map[string]bool)
	for range 1000 {
		r := OracleHow(rng)
		if r.Result == "" {
			t.Error("OracleHow returned empty result")
		}
		seen[r.Result] = true
	}
	expected := []string{
		"Surprisingly lacking", "Less than expected",
		"About average", "More than expected", "Extraordinary",
	}
	for _, e := range expected {
		if !seen[e] {
			t.Errorf("never saw %q", e)
		}
	}
}
