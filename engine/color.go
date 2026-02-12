package engine

var colors = []string{
	"crimson", "scarlet", "ruby", "vermillion", "burgundy", "cherry",
	"blood red", "rose", "coral", "rust", "amber", "copper",
	"tangerine", "burnt orange", "peach", "apricot", "terracotta",
	"sienna", "gold", "saffron", "honey", "lemon", "canary",
	"mustard", "champagne", "flaxen", "emerald", "jade", "olive",
	"forest green", "mint", "sage", "viridian", "chartreuse", "moss",
	"pine", "cobalt", "sapphire", "navy", "cerulean", "azure",
	"teal", "indigo", "slate blue", "powder blue", "midnight blue",
	"violet", "amethyst", "plum", "lavender", "mauve", "orchid",
	"royal purple", "lilac", "magenta", "ivory", "bone", "pearl",
	"ash", "charcoal", "obsidian", "silver", "steel", "slate",
	"smoke", "mahogany", "chestnut", "umber", "ochre", "sandy",
	"taupe", "walnut", "driftwood",
}

func RandomColor(rng *Randomizer) ColorResult {
	return ColorResult{Color: colors[rng.Intn(len(colors))]}
}
