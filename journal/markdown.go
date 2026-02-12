package journal

import (
	"fmt"
	"strings"

	"opse/engine"
)

func Render(j *Journal) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", j.Title)
	fmt.Fprintf(&b, "*Started: %s*\n\n", j.CreatedAt.Format("2006-01-02"))
	b.WriteString("---\n\n")

	for _, e := range j.Entries {
		if !e.Timestamp.IsZero() {
			source := "Engine"
			if e.Type == EntryNarrative {
				if e.Label != "" {
					source = e.Label
				} else {
					source = "User"
				}
			}
			fmt.Fprintf(&b, "*%s — %s*\n\n", e.Timestamp.Format("15:04"), source)
		}
		b.WriteString(e.Markdown)
		b.WriteString("\n\n")
	}
	return b.String()
}

func suitShort(c engine.Card) string {
	d := engine.SuitDomains[c.Suit]
	idx := strings.Index(d, " ")
	if idx < 0 {
		return d
	}
	return d[:idx]
}

func RenderOracleYesNo(r engine.OracleYesNoResult) string {
	answer := "No"
	if r.Answer {
		answer = "Yes"
	}
	result := answer
	if r.Modifier != "" {
		result += ", " + r.Modifier
	}
	return fmt.Sprintf("> **Oracle (Yes/No, %s):** %s", r.Likelihood, result)
}

func RenderCardTable(r engine.CardTableResult) string {
	return fmt.Sprintf("> **%s:** %s %s *(%s)*",
		r.TableName, r.Draw.Card.String(), r.Entry, suitShort(r.Draw.Card))
}

func RenderRandomEvent(r engine.RandomEventResult) string {
	return fmt.Sprintf(`> **Random Event**
> - **What happens:** %s %s *(%s)*
> - **Involving:** %s %s *(%s)*`,
		r.Action.Draw.Card.String(), r.Action.Entry, suitShort(r.Action.Draw.Card),
		r.Topic.Draw.Card.String(), r.Topic.Entry, suitShort(r.Topic.Draw.Card))
}

func RenderSetTheScene(r engine.SetTheSceneResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "> **Set the Scene**\n> - **Complication:** %s", r.Complication.Result)
	if r.Altered && r.AlteredScene != nil {
		fmt.Fprintf(&b, "\n> - **Altered Scene:** %s", r.AlteredScene.Result)
		if r.AlteredScene.Cascade != nil {
			switch c := r.AlteredScene.Cascade.(type) {
			case *engine.SceneComplicationResult:
				fmt.Fprintf(&b, "\n>   - **Cascade:** %s", c.Result)
			case *engine.PacingMoveResult:
				fmt.Fprintf(&b, "\n>   - **Cascade (Pacing Move):** %s", c.Result)
				if c.RandomEvent != nil {
					fmt.Fprintf(&b, "\n>     - **Random Event:** %s / %s",
						c.RandomEvent.Action.Entry, c.RandomEvent.Topic.Entry)
				}
			case *engine.RandomEventResult:
				fmt.Fprintf(&b, "\n>   - **Cascade (Random Event):** %s / %s",
					c.Action.Entry, c.Topic.Entry)
			}
		}
	}
	return b.String()
}

func RenderPacingMove(r engine.PacingMoveResult) string {
	s := fmt.Sprintf("> **Pacing Move:** %s", r.Result)
	if r.RandomEvent != nil {
		s += fmt.Sprintf("\n> - **Random Event:** %s %s / %s %s",
			r.RandomEvent.Action.Draw.Card.String(), r.RandomEvent.Action.Entry,
			r.RandomEvent.Topic.Draw.Card.String(), r.RandomEvent.Topic.Entry)
	}
	return s
}

func RenderGeneric(r engine.GenericGeneratorResult) string {
	return fmt.Sprintf(`> **Generic Generator**
> - **What it does:** %s %s *(%s)*
> - **How it looks:** %s %s *(%s)*
> - **How significant:** %s`,
		r.Action.Draw.Card.String(), r.Action.Entry, suitShort(r.Action.Draw.Card),
		r.Detail.Draw.Card.String(), r.Detail.Entry, suitShort(r.Detail.Draw.Card),
		r.Significance.Result)
}

func RenderPlotHook(r engine.PlotHookResult) string {
	return fmt.Sprintf(`> **Plot Hook**
> - **Objective:** %s
> - **Adversary:** %s
> - **Reward:** %s`,
		r.Objective, r.Adversary, r.Reward)
}

func RenderNPC(r engine.NPCResult) string {
	return fmt.Sprintf(`> **NPC**
> - **Identity:** %s *(%s)*
> - **Goal:** %s *(%s)*
> - **Feature:** %s — %s *(%s)*
> - **Attitude:** %s
> - **Topic:** %s *(%s)*`,
		r.Identity.Entry, suitShort(r.Identity.Draw.Card),
		r.Goal.Entry, suitShort(r.Goal.Draw.Card),
		r.Feature, r.FeatureDetail.Entry, suitShort(r.FeatureDetail.Draw.Card),
		r.Attitude.Result,
		r.Topic.Entry, suitShort(r.Topic.Draw.Card))
}

func RenderDungeonRoom(r engine.DungeonRoomResult) string {
	return fmt.Sprintf(`> **Dungeon Room**
> - **Location:** %s
> - **Encounter:** %s
> - **Object:** %s
> - **Exits:** %s`,
		r.Location, r.Encounter, r.Object, r.Exits)
}

func RenderHex(r engine.HexResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, `> **Hex**
> - **Terrain:** %s
> - **Contents:** %s`, r.Terrain, r.Contents)
	if r.Feature != "" {
		fmt.Fprintf(&b, "\n> - **Feature:** %s", r.Feature)
	}
	fmt.Fprintf(&b, "\n> - **Event:** %s", r.Event)
	if r.RandomEvent != nil {
		fmt.Fprintf(&b, "\n>   - **Random Event:** %s / %s",
			r.RandomEvent.Action.Entry, r.RandomEvent.Topic.Entry)
	}
	return b.String()
}

func RenderDiceRoll(r engine.DiceRollResult) string {
	parts := make([]string, len(r.Rolls))
	for i, v := range r.Rolls {
		if r.Kept[i] {
			parts[i] = fmt.Sprintf("[%d]", v)
		} else {
			parts[i] = fmt.Sprintf("~~%d~~", v)
		}
	}
	mod := ""
	if r.Expression.Modifier > 0 {
		mod = fmt.Sprintf(" + %d", r.Expression.Modifier)
	} else if r.Expression.Modifier < 0 {
		mod = fmt.Sprintf(" - %d", -r.Expression.Modifier)
	}
	return fmt.Sprintf("> **Dice (%s):** %s%s = **%d**",
		r.Expression.Raw, strings.Join(parts, " "), mod, r.Total)
}

func RenderCoinFlip(r engine.CoinFlipResult) string {
	if len(r.Flips) == 1 {
		face := "Tails"
		if r.Flips[0] {
			face = "Heads"
		}
		return fmt.Sprintf("> **Coin Flip:** %s", face)
	}
	faces := make([]string, len(r.Flips))
	for i, f := range r.Flips {
		if f {
			faces[i] = "H"
		} else {
			faces[i] = "T"
		}
	}
	return fmt.Sprintf("> **Coin Flip (%d):** %s (%dH / %dT)",
		len(r.Flips), strings.Join(faces, ", "), r.Heads, r.Tails)
}

func RenderCardDraw(r engine.CardDrawResult) string {
	cards := make([]string, len(r.Cards))
	for i, c := range r.Cards {
		cards[i] = c.String()
	}
	return fmt.Sprintf("> **Card Draw:** %s *(%d remaining)*",
		strings.Join(cards, " "), r.Remaining)
}

func RenderDirection(r engine.DirectionResult) string {
	return fmt.Sprintf("> **Direction:** %s %s", r.Arrow, r.Direction)
}

func RenderWeather(r engine.WeatherResult) string {
	return fmt.Sprintf("> **Weather:** %s, %s, %s", r.Condition, r.Temperature, r.Wind)
}

func RenderColor(r engine.ColorResult) string {
	return fmt.Sprintf("> **Color:** %s", r.Color)
}

func RenderSound(r engine.SoundResult) string {
	return fmt.Sprintf("> **Sound:** *%s* (%s)", r.Sound, r.Category)
}

func RenderDungeonTheme(r engine.DungeonThemeResult) string {
	return fmt.Sprintf(`> **Dungeon Theme**
> - **How it looks:** %s %s *(%s)*
> - **How it's used:** %s %s *(%s)*`,
		r.Looks.Draw.Card.String(), r.Looks.Entry, suitShort(r.Looks.Draw.Card),
		r.Used.Draw.Card.String(), r.Used.Entry, suitShort(r.Used.Draw.Card))
}
