package engine

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var diceRegex = regexp.MustCompile(
	`(?i)^(\d*)d(\d+)(!)?(?:(kh|kl)(\d+))?([+-]\d+)?$`,
)

func ParseDice(input string) (DiceExpression, error) {
	clean := strings.ReplaceAll(strings.TrimSpace(input), " ", "")
	m := diceRegex.FindStringSubmatch(clean)
	if m == nil {
		return DiceExpression{}, fmt.Errorf("invalid dice expression: %q", input)
	}

	count := 1
	if m[1] != "" {
		count, _ = strconv.Atoi(m[1])
	}
	sides, _ := strconv.Atoi(m[2])
	explode := m[3] == "!"
	keepMode := strings.ToLower(m[4])
	keepN := 0
	if m[5] != "" {
		keepN, _ = strconv.Atoi(m[5])
	}
	modifier := 0
	if m[6] != "" {
		modifier, _ = strconv.Atoi(m[6])
	}

	if count < 1 || count > 999 {
		return DiceExpression{}, fmt.Errorf("dice count must be 1-999, got %d", count)
	}
	if sides < 2 || sides > 100 {
		return DiceExpression{}, fmt.Errorf("die sides must be 2-100, got %d", sides)
	}
	if keepN > count {
		keepN = count
	}
	if keepMode != "" && keepN == 0 {
		keepN = 1
	}

	return DiceExpression{
		Count: count, Sides: sides, Explode: explode,
		KeepMode: keepMode, KeepN: keepN, Modifier: modifier,
		Raw: input,
	}, nil
}

func RollDice(rng *Randomizer, expr DiceExpression) DiceRollResult {
	rolls := make([]int, 0, expr.Count)
	for range expr.Count {
		val := rng.RollDN(expr.Sides)
		if expr.Explode {
			total := val
			depth := 0
			for val == expr.Sides && depth < 100 {
				val = rng.RollDN(expr.Sides)
				total += val
				depth++
			}
			rolls = append(rolls, total)
		} else {
			rolls = append(rolls, val)
		}
	}

	kept := make([]bool, len(rolls))
	if expr.KeepMode == "" {
		for i := range kept {
			kept[i] = true
		}
	} else {
		type iv struct{ i, v int }
		indexed := make([]iv, len(rolls))
		for i, v := range rolls {
			indexed[i] = iv{i, v}
		}
		if expr.KeepMode == "kh" {
			sort.Slice(indexed, func(a, b int) bool { return indexed[a].v > indexed[b].v })
		} else {
			sort.Slice(indexed, func(a, b int) bool { return indexed[a].v < indexed[b].v })
		}
		for i := 0; i < expr.KeepN && i < len(indexed); i++ {
			kept[indexed[i].i] = true
		}
	}

	subtotal := 0
	for i, v := range rolls {
		if kept[i] {
			subtotal += v
		}
	}

	return DiceRollResult{
		Expression: expr, Rolls: rolls, Kept: kept,
		Subtotal: subtotal, Total: subtotal + expr.Modifier,
	}
}
