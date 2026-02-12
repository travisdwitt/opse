package engine

import "testing"

func TestParseDice_Valid(t *testing.T) {
	tests := []struct {
		input        string
		count, sides int
		mod          int
		keep         string
		keepN        int
	}{
		{"d20", 1, 20, 0, "", 0},
		{"2d6", 2, 6, 0, "", 0},
		{"4d6kh3", 4, 6, 0, "kh", 3},
		{"2d8+5", 2, 8, 5, "", 0},
		{"d100-10", 1, 100, -10, "", 0},
		{"3d8+5", 3, 8, 5, "", 0},
		{"2d6!", 2, 6, 0, "", 0},
		{"2d20kl1", 2, 20, 0, "kl", 1},
	}
	for _, tt := range tests {
		expr, err := ParseDice(tt.input)
		if err != nil {
			t.Errorf("ParseDice(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if expr.Count != tt.count {
			t.Errorf("%q: count %d, want %d", tt.input, expr.Count, tt.count)
		}
		if expr.Sides != tt.sides {
			t.Errorf("%q: sides %d, want %d", tt.input, expr.Sides, tt.sides)
		}
		if expr.Modifier != tt.mod {
			t.Errorf("%q: mod %d, want %d", tt.input, expr.Modifier, tt.mod)
		}
		if expr.KeepMode != tt.keep {
			t.Errorf("%q: keep %q, want %q", tt.input, expr.KeepMode, tt.keep)
		}
		if expr.KeepN != tt.keepN {
			t.Errorf("%q: keepN %d, want %d", tt.input, expr.KeepN, tt.keepN)
		}
	}
}

func TestParseDice_Invalid(t *testing.T) {
	invalid := []string{"garbage", "0d6", "1d1", "1000d6", "d0", ""}
	for _, input := range invalid {
		_, err := ParseDice(input)
		if err == nil {
			t.Errorf("ParseDice(%q) should error", input)
		}
	}
}

func TestRollDice_CorrectCount(t *testing.T) {
	rng := NewSeededRandomizer(200, 0)
	expr, _ := ParseDice("5d6")
	r := RollDice(rng, expr)
	if len(r.Rolls) != 5 {
		t.Errorf("expected 5 rolls, got %d", len(r.Rolls))
	}
	for _, v := range r.Rolls {
		if v < 1 || v > 6 {
			t.Errorf("roll %d out of range", v)
		}
	}
}

func TestRollDice_KeepHighest(t *testing.T) {
	rng := NewSeededRandomizer(201, 0)
	expr, _ := ParseDice("4d6kh3")
	r := RollDice(rng, expr)
	keptCount := 0
	for _, k := range r.Kept {
		if k {
			keptCount++
		}
	}
	if keptCount != 3 {
		t.Errorf("expected 3 kept, got %d", keptCount)
	}
	// Verify subtotal equals sum of kept dice
	sum := 0
	for i, v := range r.Rolls {
		if r.Kept[i] {
			sum += v
		}
	}
	if sum != r.Subtotal {
		t.Errorf("subtotal %d != sum of kept %d", r.Subtotal, sum)
	}
}

func TestRollDice_Modifier(t *testing.T) {
	rng := NewSeededRandomizer(202, 0)
	expr, _ := ParseDice("1d6+5")
	r := RollDice(rng, expr)
	if r.Total != r.Subtotal+5 {
		t.Errorf("total %d != subtotal %d + 5", r.Total, r.Subtotal)
	}

	expr2, _ := ParseDice("1d6-3")
	r2 := RollDice(rng, expr2)
	if r2.Total != r2.Subtotal-3 {
		t.Errorf("total %d != subtotal %d - 3", r2.Total, r2.Subtotal)
	}
}

func TestRollDice_Exploding(t *testing.T) {
	rng := NewSeededRandomizer(203, 0)
	expr, _ := ParseDice("10d6!")
	r := RollDice(rng, expr)
	if len(r.Rolls) != 10 {
		t.Errorf("expected 10 rolls, got %d", len(r.Rolls))
	}
	// Exploding dice can produce values > sides
	hasHigh := false
	for _, v := range r.Rolls {
		if v > 6 {
			hasHigh = true
		}
	}
	_ = hasHigh // just ensure no panic
}
