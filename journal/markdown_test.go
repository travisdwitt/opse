package journal

import (
	"strings"
	"testing"

	"opse/engine"
)

func TestRenderOracleYesNo(t *testing.T) {
	r := engine.OracleYesNoResult{
		Likelihood: "Even", AnswerRoll: 5, Answer: true,
		ModRoll: 6, Modifier: "and...",
	}
	md := RenderOracleYesNo(r)
	if !strings.Contains(md, "Oracle (Yes/No, Even)") {
		t.Error("should contain oracle label")
	}
	if !strings.Contains(md, "Yes, and...") {
		t.Errorf("should contain 'Yes, and...', got %q", md)
	}
}

func TestRenderDiceRoll_WithDropped(t *testing.T) {
	r := engine.DiceRollResult{
		Expression: engine.DiceExpression{Raw: "4d6kh3", Modifier: 0},
		Rolls:      []int{6, 4, 3, 2},
		Kept:       []bool{true, true, true, false},
		Subtotal:   13,
		Total:      13,
	}
	md := RenderDiceRoll(r)
	if !strings.Contains(md, "~~2~~") {
		t.Error("dropped die should have strikethrough")
	}
	if !strings.Contains(md, "[6]") {
		t.Error("kept die should be in brackets")
	}
}

func TestRenderDiceRoll_WithModifier(t *testing.T) {
	r := engine.DiceRollResult{
		Expression: engine.DiceExpression{Raw: "2d6+5", Modifier: 5},
		Rolls:      []int{3, 4},
		Kept:       []bool{true, true},
		Subtotal:   7,
		Total:      12,
	}
	md := RenderDiceRoll(r)
	if !strings.Contains(md, "+ 5") {
		t.Error("should show modifier")
	}
	if !strings.Contains(md, "**12**") {
		t.Error("should show total in bold")
	}
}

func TestRender_IncludesTitleAndDate(t *testing.T) {
	j := New("My Adventure", "/tmp/test.md")
	j.AddEntry(Entry{Type: EntryNarrative, Markdown: "Hello world"})
	md := Render(j)
	if !strings.HasPrefix(md, "# My Adventure") {
		t.Error("should start with title")
	}
	if !strings.Contains(md, "*Started:") {
		t.Error("should contain start date")
	}
	if !strings.Contains(md, "Hello world") {
		t.Error("should contain entry")
	}
}
