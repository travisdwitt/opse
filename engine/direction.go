package engine

type dirEntry struct{ Name, Abbrev, Arrow string }

var dirs4 = []dirEntry{
	{"North", "N", "↑"}, {"East", "E", "→"},
	{"South", "S", "↓"}, {"West", "W", "←"},
}

var dirs8 = []dirEntry{
	{"North", "N", "↑"}, {"Northeast", "NE", "↗"},
	{"East", "E", "→"}, {"Southeast", "SE", "↘"},
	{"South", "S", "↓"}, {"Southwest", "SW", "↙"},
	{"West", "W", "←"}, {"Northwest", "NW", "↖"},
}

var dirs16 = []dirEntry{
	{"North", "N", "↑"}, {"North-Northeast", "NNE", "↑"},
	{"Northeast", "NE", "↗"}, {"East-Northeast", "ENE", "→"},
	{"East", "E", "→"}, {"East-Southeast", "ESE", "→"},
	{"Southeast", "SE", "↘"}, {"South-Southeast", "SSE", "↓"},
	{"South", "S", "↓"}, {"South-Southwest", "SSW", "↓"},
	{"Southwest", "SW", "↙"}, {"West-Southwest", "WSW", "←"},
	{"West", "W", "←"}, {"West-Northwest", "WNW", "←"},
	{"Northwest", "NW", "↖"}, {"North-Northwest", "NNW", "↑"},
}

func RandomDirection(rng *Randomizer, points int) DirectionResult {
	table := dirs8
	switch points {
	case 4:
		table = dirs4
	case 16:
		table = dirs16
	}
	d := table[rng.Intn(len(table))]
	return DirectionResult{Direction: d.Name, Abbrev: d.Abbrev, Arrow: d.Arrow}
}
