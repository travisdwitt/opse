package ui

import (
	"fmt"
	"strings"
	"time"

	"opse/engine"
)

func FormatEntryHeader(ts time.Time, source string) string {
	if ts.IsZero() {
		return ""
	}
	return DimStyle.Render(fmt.Sprintf("%s  %s", ts.Format("15:04"), source))
}

func RenderCardForTUI(c engine.Card) string {
	symbol := engine.SuitSymbols[c.Suit]
	style := SuitWhiteStyle
	if c.Suit == engine.Diamonds || c.Suit == engine.Hearts {
		style = SuitRedStyle
	}
	return fmt.Sprintf("%s%s", c.Rank.String(), style.Render(symbol))
}

func cardDomain(c engine.Card) string {
	return DimStyle.Render(c.Domain())
}

func RenderOracleYesNoTUI(r engine.OracleYesNoResult) string {
	answer := "No"
	if r.Answer {
		answer = "Yes"
	}
	result := answer
	if r.Modifier != "" {
		result += ", " + r.Modifier
	}
	modStr := "(none)"
	if r.Modifier != "" {
		modStr = r.Modifier
	}
	body := fmt.Sprintf(
		" Answer: %d → %s\n Modifier: %d → %s\n\n Result: %s",
		r.AnswerRoll, answer, r.ModRoll, modStr,
		ResultLabelStyle.Render(result),
	)
	title := fmt.Sprintf("Oracle: Yes/No (%s)", r.Likelihood)
	return ResultBlockStyle.Render(ResultLabelStyle.Render(title) + "\n" + body)
}

func RenderOracleHowTUI(r engine.OracleHowResult) string {
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Oracle: How") + "\n " + r.Result)
}

func RenderCardTableTUI(r engine.CardTableResult) string {
	card := RenderCardForTUI(r.Draw.Card)
	domain := cardDomain(r.Draw.Card)
	body := fmt.Sprintf(" %s %s — %s", card, r.Entry, domain)
	return ResultBlockStyle.Render(ResultLabelStyle.Render(r.TableName) + "\n" + body)
}

func RenderRandomEventTUI(r engine.RandomEventResult) string {
	body := fmt.Sprintf(
		" What happens: %s %s — %s\n Involving:    %s %s — %s",
		RenderCardForTUI(r.Action.Draw.Card), r.Action.Entry, cardDomain(r.Action.Draw.Card),
		RenderCardForTUI(r.Topic.Draw.Card), r.Topic.Entry, cardDomain(r.Topic.Draw.Card),
	)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Random Event") + "\n" + body)
}

func RenderSetTheSceneTUI(r engine.SetTheSceneResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, " Complication: %s", r.Complication.Result)
	if r.Altered && r.AlteredScene != nil {
		fmt.Fprintf(&b, "\n Altered Scene: %s", r.AlteredScene.Result)
		if r.AlteredScene.Cascade != nil {
			switch c := r.AlteredScene.Cascade.(type) {
			case *engine.SceneComplicationResult:
				fmt.Fprintf(&b, "\n   → %s", c.Result)
			case *engine.PacingMoveResult:
				fmt.Fprintf(&b, "\n   → Pacing Move: %s", c.Result)
				if c.RandomEvent != nil {
					fmt.Fprintf(&b, "\n     → Random Event: %s / %s",
						c.RandomEvent.Action.Entry, c.RandomEvent.Topic.Entry)
				}
			case *engine.RandomEventResult:
				fmt.Fprintf(&b, "\n   → Random Event: %s / %s",
					c.Action.Entry, c.Topic.Entry)
			}
		}
	} else {
		fmt.Fprintf(&b, "\n %s", DimStyle.Render("(not altered)"))
	}
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Set the Scene") + "\n" + b.String())
}

func RenderPacingMoveTUI(r engine.PacingMoveResult) string {
	body := " " + r.Result
	if r.RandomEvent != nil {
		body += fmt.Sprintf("\n\n %s %s — %s\n %s %s — %s",
			RenderCardForTUI(r.RandomEvent.Action.Draw.Card),
			r.RandomEvent.Action.Entry, cardDomain(r.RandomEvent.Action.Draw.Card),
			RenderCardForTUI(r.RandomEvent.Topic.Draw.Card),
			r.RandomEvent.Topic.Entry, cardDomain(r.RandomEvent.Topic.Draw.Card),
		)
	}
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Pacing Move") + "\n" + body)
}

func RenderFailureMoveTUI(r engine.FailureMoveResult) string {
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Failure Move") + "\n " + r.Result)
}

func RenderGenericTUI(r engine.GenericGeneratorResult) string {
	body := fmt.Sprintf(
		" What it does:    %s %s — %s\n How it looks:    %s %s — %s\n How significant: %s",
		RenderCardForTUI(r.Action.Draw.Card), r.Action.Entry, cardDomain(r.Action.Draw.Card),
		RenderCardForTUI(r.Detail.Draw.Card), r.Detail.Entry, cardDomain(r.Detail.Draw.Card),
		r.Significance.Result,
	)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Generic Generator") + "\n" + body)
}

func RenderPlotHookTUI(r engine.PlotHookResult) string {
	body := fmt.Sprintf(
		" Objective:  %s\n Adversary:  %s\n Reward:     %s",
		r.Objective, r.Adversary, r.Reward,
	)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Plot Hook") + "\n" + body)
}

func RenderNPCTUI(r engine.NPCResult) string {
	body := fmt.Sprintf(
		" Identity: %s %s — %s\n Goal:     %s %s — %s\n Feature:  %s — %s %s — %s\n Attitude: %s\n Topic:    %s %s — %s",
		RenderCardForTUI(r.Identity.Draw.Card), r.Identity.Entry, cardDomain(r.Identity.Draw.Card),
		RenderCardForTUI(r.Goal.Draw.Card), r.Goal.Entry, cardDomain(r.Goal.Draw.Card),
		r.Feature, RenderCardForTUI(r.FeatureDetail.Draw.Card), r.FeatureDetail.Entry, cardDomain(r.FeatureDetail.Draw.Card),
		r.Attitude.Result,
		RenderCardForTUI(r.Topic.Draw.Card), r.Topic.Entry, cardDomain(r.Topic.Draw.Card),
	)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("NPC") + "\n" + body)
}

func RenderDungeonThemeTUI(r engine.DungeonThemeResult) string {
	body := fmt.Sprintf(
		" How it looks: %s %s — %s\n How it's used: %s %s — %s",
		RenderCardForTUI(r.Looks.Draw.Card), r.Looks.Entry, cardDomain(r.Looks.Draw.Card),
		RenderCardForTUI(r.Used.Draw.Card), r.Used.Entry, cardDomain(r.Used.Draw.Card),
	)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Dungeon Theme") + "\n" + body)
}

func RenderDungeonRoomTUI(r engine.DungeonRoomResult) string {
	body := fmt.Sprintf(
		" Location:  %s\n Encounter: %s\n Object:    %s\n Exits:     %s",
		r.Location, r.Encounter, r.Object, r.Exits,
	)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Dungeon Room") + "\n" + body)
}

func RenderHexTUI(r engine.HexResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, " Terrain:  %s\n Contents: %s", r.Terrain, r.Contents)
	if r.Feature != "" {
		fmt.Fprintf(&b, "\n Feature:  %s", r.Feature)
	}
	fmt.Fprintf(&b, "\n Event:    %s", r.Event)
	if r.RandomEvent != nil {
		fmt.Fprintf(&b, "\n   %s %s / %s %s",
			RenderCardForTUI(r.RandomEvent.Action.Draw.Card), r.RandomEvent.Action.Entry,
			RenderCardForTUI(r.RandomEvent.Topic.Draw.Card), r.RandomEvent.Topic.Entry)
	}
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Hex") + "\n" + b.String())
}

func RenderDiceRollTUI(r engine.DiceRollResult) string {
	parts := make([]string, len(r.Rolls))
	for i, v := range r.Rolls {
		if r.Kept[i] {
			parts[i] = fmt.Sprintf("[%d]", v)
		} else {
			parts[i] = DimStyle.Render(fmt.Sprintf(" %d ", v))
		}
	}
	mod := ""
	if r.Expression.Modifier > 0 {
		mod = fmt.Sprintf(" + %d", r.Expression.Modifier)
	}
	if r.Expression.Modifier < 0 {
		mod = fmt.Sprintf(" - %d", -r.Expression.Modifier)
	}
	body := fmt.Sprintf(" Rolls: %s%s = %s",
		strings.Join(parts, " "), mod,
		ResultLabelStyle.Render(fmt.Sprintf("%d", r.Total)))
	title := fmt.Sprintf("Dice: %s", r.Expression.Raw)
	return ResultBlockStyle.Render(ResultLabelStyle.Render(title) + "\n" + body)
}

func RenderCoinFlipTUI(r engine.CoinFlipResult) string {
	if len(r.Flips) == 1 {
		face := "Tails"
		if r.Flips[0] {
			face = "Heads"
		}
		return ResultBlockStyle.Render(ResultLabelStyle.Render("Coin Flip") + "\n " + face)
	}
	faces := make([]string, len(r.Flips))
	for i, f := range r.Flips {
		if f {
			faces[i] = "H"
		} else {
			faces[i] = "T"
		}
	}
	body := fmt.Sprintf(" %s (%dH / %dT)", strings.Join(faces, ", "), r.Heads, r.Tails)
	return ResultBlockStyle.Render(ResultLabelStyle.Render(fmt.Sprintf("Coin Flip (%d)", len(r.Flips))) + "\n" + body)
}

func RenderCardDrawTUI(r engine.CardDrawResult) string {
	cards := make([]string, len(r.Cards))
	for i, c := range r.Cards {
		cards[i] = RenderCardForTUI(c)
	}
	body := fmt.Sprintf(" %s\n %s", strings.Join(cards, "  "),
		DimStyle.Render(fmt.Sprintf("(%d remaining)", r.Remaining)))
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Card Draw") + "\n" + body)
}

func RenderDirectionTUI(r engine.DirectionResult) string {
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Direction") + "\n " + r.Arrow + " " + r.Direction)
}

func RenderWeatherTUI(r engine.WeatherResult) string {
	body := fmt.Sprintf(" Condition:   %s\n Temperature: %s\n Wind:        %s",
		r.Condition, r.Temperature, r.Wind)
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Weather") + "\n" + body)
}

func RenderColorTUI(r engine.ColorResult) string {
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Random Color") + "\n " + r.Color)
}

func RenderSoundTUI(r engine.SoundResult) string {
	body := fmt.Sprintf(" *%s*  %s", r.Sound, DimStyle.Render("("+r.Category+")"))
	return ResultBlockStyle.Render(ResultLabelStyle.Render("Random Sound") + "\n" + body)
}
