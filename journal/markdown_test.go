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

func TestRender_CharEntry(t *testing.T) {
	j := New("Test", "/tmp/test.md")
	j.AddEntry(Entry{Type: EntryNarrative, Label: "Elara", Markdown: "I search the room"})
	md := Render(j)
	// Without portrait dir, no image tags or bold name
	if strings.Contains(md, "![Elara]") {
		t.Error("should NOT contain image tag")
	}
	if strings.Contains(md, "**Elara:**") {
		t.Error("should NOT contain bold name")
	}
	if !strings.Contains(md, "I search the room") {
		t.Error("should contain the text")
	}
}

func TestStripCharMarkdown_RoundTrip(t *testing.T) {
	// Backward compat: old journals with image tags should strip correctly
	text := "![Elara](portraits/elara.png)\n\n**Elara:** I search the room"
	got := stripCharMarkdown(text, "Elara")
	if got != "I search the room" {
		t.Errorf("expected plain text, got %q", got)
	}
}

func TestStripCharMarkdown_PlainText(t *testing.T) {
	// Old entries without image tags should pass through unchanged
	got := stripCharMarkdown("I search the room", "Elara")
	if got != "I search the room" {
		t.Errorf("expected unchanged text, got %q", got)
	}
}
