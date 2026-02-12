package engine

func OracleYesNo(rng *Randomizer, likelihood string) OracleYesNoResult {
	answerRoll := rng.RollD6()
	modRoll := rng.RollD6()

	threshold := 4
	switch likelihood {
	case "Likely":
		threshold = 3
	case "Unlikely":
		threshold = 5
	}

	answer := answerRoll >= threshold

	var modifier string
	switch modRoll {
	case 1:
		modifier = "but..."
	case 6:
		modifier = "and..."
	}

	return OracleYesNoResult{
		Likelihood: likelihood,
		AnswerRoll: answerRoll,
		Answer:     answer,
		ModRoll:    modRoll,
		Modifier:   modifier,
	}
}

var oracleHowTable = [7]string{
	"",
	"Surprisingly lacking",
	"Less than expected",
	"About average",
	"About average",
	"More than expected",
	"Extraordinary",
}

func OracleHow(rng *Randomizer) OracleHowResult {
	roll := rng.RollD6()
	return OracleHowResult{Roll: roll, Result: oracleHowTable[roll]}
}
