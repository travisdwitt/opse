package engine

import "sort"

var soundCategories = map[string][]string{
	"impact": {
		"thud", "crash", "bang", "slam", "thump", "clang", "clatter",
		"crack", "smash", "crunch", "thwack", "wallop", "whack",
		"clunk", "bonk", "bam", "wham", "bash", "knock", "pound",
		"smack", "clink", "clash", "rattle", "jolt", "bump",
		"crumble", "shatter", "splinter", "collapse",
	},
	"voice": {
		"shout", "scream", "whisper", "murmur", "groan", "moan",
		"cry", "wail", "shriek", "gasp", "sigh", "laugh", "cackle",
		"sob", "whimper", "yell", "bark", "growl", "hiss", "chant",
		"sing", "hum", "mutter", "stammer", "babble", "giggle",
		"snicker", "roar", "bellow", "howl",
	},
	"animal": {
		"howl", "growl", "snarl", "bark", "screech", "squawk",
		"chirp", "hiss", "buzz", "croak", "ribbit", "caw", "hoot",
		"neigh", "bray", "bleat", "snort", "purr", "chitter",
		"squeal", "trumpet", "warble", "trill", "coo", "twitter",
	},
	"nature": {
		"thunder", "rumble", "gust", "howl", "patter", "downpour",
		"drip", "splash", "rush", "roar", "crackle", "rustle",
		"whoosh", "whistle", "trickle", "cascade", "gurgle", "surge",
		"hail", "lash", "whip", "sway", "groan", "snap", "crash",
	},
	"mechanical": {
		"clank", "grind", "screech", "squeak", "ratchet", "click",
		"tick", "whir", "hum", "drone", "clang", "ping", "ding",
		"chime", "gong", "scrape", "rasp", "creak", "jangle",
		"chain", "bolt", "latch", "gear", "spring", "lever",
	},
	"fire": {
		"crackle", "roar", "hiss", "pop", "sizzle", "whoosh",
		"blaze", "sputter", "flicker", "snap", "boom", "fizzle",
		"spark", "ignite", "smolder",
	},
	"water": {
		"splash", "drip", "gurgle", "bubble", "trickle", "pour",
		"gush", "slosh", "lap", "spray", "cascade", "surge",
		"plop", "splatter", "swirl", "ripple", "churn", "foam",
		"rush", "deluge",
	},
	"movement": {
		"whoosh", "swoosh", "flutter", "flap", "rustle", "shuffle",
		"stomp", "trudge", "scurry", "dash", "skitter", "slither",
		"creep", "thud", "patter", "gallop", "clop", "scramble",
		"tumble", "slide",
	},
	"eerie": {
		"whisper", "moan", "wail", "echo", "drone", "hum",
		"chime", "toll", "keen", "rasp", "scratch", "tap",
		"creak", "groan", "rattle", "thrum", "pulse", "vibrate",
		"resonance", "silence",
	},
	"combat": {
		"clang", "clash", "ring", "twang", "thwack", "slash",
		"stab", "pierce", "parry", "deflect", "shatter", "crack",
		"snap", "whip", "slice", "cleave", "chop", "hack",
		"thrust", "swipe", "volley", "ricochet", "whistle", "zing",
		"notch",
	},
}

type flatSound struct{ sound, category string }

var allSoundsFlat []flatSound

func init() {
	for cat, sounds := range soundCategories {
		for _, s := range sounds {
			allSoundsFlat = append(allSoundsFlat, flatSound{s, cat})
		}
	}
}

func SoundCategoryNames() []string {
	names := make([]string, 0, len(soundCategories))
	for k := range soundCategories {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func RandomSound(rng *Randomizer, category string) SoundResult {
	if category != "" {
		if sounds, ok := soundCategories[category]; ok {
			return SoundResult{Sound: sounds[rng.Intn(len(sounds))], Category: category}
		}
	}
	entry := allSoundsFlat[rng.Intn(len(allSoundsFlat))]
	return SoundResult{Sound: entry.sound, Category: entry.category}
}
