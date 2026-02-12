package engine

func hexTerrain(roll int) string {
	switch {
	case roll <= 2:
		return "Same as current hex"
	case roll <= 4:
		return "Common terrain"
	case roll == 5:
		return "Uncommon terrain"
	default:
		return "Rare terrain"
	}
}

func hexContents(roll int) string {
	if roll == 6 {
		return "Roll a FEATURE"
	}
	return "Nothing notable"
}

var hexFeatures = [7]string{
	"",
	"Notable structure",
	"Dangerous hazard",
	"A settlement",
	"Strange natural feature",
	"New region (set new terrain types)",
	"DUNGEON CRAWLER entrance",
}

func hexEvent(roll int) string {
	if roll >= 5 {
		return "RANDOM EVENT then SET THE SCENE"
	}
	return "None"
}

func HexCrawl(rng *Randomizer, deck *Deck) HexResult {
	tr := rng.RollD6()
	cr := rng.RollD6()
	er := rng.RollD6()

	result := HexResult{
		TerrainRoll: tr, Terrain: hexTerrain(tr),
		ContentsRoll: cr, Contents: hexContents(cr),
		EventRoll: er, Event: hexEvent(er),
	}

	if cr == 6 {
		fr := rng.RollD6()
		result.FeatureRoll = fr
		result.Feature = hexFeatures[fr]
	}

	if er >= 5 {
		evt := RandomEvent(deck, rng)
		result.RandomEvent = &evt
	}

	return result
}
