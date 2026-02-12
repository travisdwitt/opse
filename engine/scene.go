package engine

var sceneComplications = [7]string{
	"",
	"Hostile forces oppose you",
	"An obstacle blocks your way",
	"Wouldn't it suck if...",
	"An NPC acts suddenly",
	"All is not as it seems",
	"Things actually go as planned",
}

var alteredScenes = [7]string{
	"",
	"A major detail of the scene is enhanced or somehow worse",
	"The environment is different",
	"Unexpected NPCs are present",
	"Add a SCENE COMPLICATION",
	"Add a PACING MOVE",
	"Add a RANDOM EVENT",
}

func SceneComplication(rng *Randomizer) SceneComplicationResult {
	roll := rng.RollD6()
	return SceneComplicationResult{Roll: roll, Result: sceneComplications[roll]}
}

func SetTheScene(rng *Randomizer, deck *Deck) SetTheSceneResult {
	comp := SceneComplication(rng)
	altRoll := rng.RollD6()
	result := SetTheSceneResult{
		Complication: comp,
		AlteredRoll:  altRoll,
		Altered:      altRoll >= 5,
	}

	if !result.Altered {
		return result
	}

	altSceneRoll := rng.RollD6()
	alt := AlteredSceneResult{
		Roll:   altSceneRoll,
		Result: alteredScenes[altSceneRoll],
	}

	switch altSceneRoll {
	case 4:
		cascade := SceneComplication(rng)
		alt.Cascade = &cascade
	case 5:
		cascade := PacingMove(rng, deck)
		alt.Cascade = &cascade
	case 6:
		cascade := RandomEvent(deck, rng)
		alt.Cascade = &cascade
	}

	result.AlteredScene = &alt
	return result
}
