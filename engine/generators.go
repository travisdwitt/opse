package engine

func GenericGenerator(deck *Deck, rng *Randomizer) GenericGeneratorResult {
	return GenericGeneratorResult{
		Action:       ActionFocus(deck),
		Detail:       DetailFocus(deck),
		Significance: OracleHow(rng),
	}
}

var objectives = [7]string{
	"",
	"Eliminate a threat",
	"Learn the truth",
	"Recover something valuable",
	"Escort or deliver to safety",
	"Restore something broken",
	"Save an ally in peril",
}

var adversaries = [7]string{
	"",
	"A powerful organization",
	"Outlaws",
	"Guardians",
	"Local inhabitants",
	"Enemy horde or force",
	"A new or recurring villain",
}

var rewards = [7]string{
	"",
	"Money or valuables",
	"Money or valuables",
	"Knowledge and secrets",
	"Support of an ally",
	"Advance a plot arc",
	"A unique item of power",
}

func PlotHook(rng *Randomizer) PlotHookResult {
	o := rng.RollD6()
	a := rng.RollD6()
	r := rng.RollD6()
	return PlotHookResult{
		ObjectiveRoll: o, Objective: objectives[o],
		AdversaryRoll: a, Adversary: adversaries[a],
		RewardRoll: r, Reward: rewards[r],
	}
}

var identityTable = map[Rank]string{
	RankTwo: "Outlaw", RankThree: "Drifter", RankFour: "Tradesman",
	RankFive: "Commoner", RankSix: "Soldier", RankSeven: "Merchant",
	RankEight: "Specialist", RankNine: "Entertainer", RankTen: "Adherent",
	RankJack: "Leader", RankQueen: "Mystic", RankKing: "Adventurer",
	RankAce: "Lord",
}

var goalTable = map[Rank]string{
	RankTwo: "Obtain", RankThree: "Learn", RankFour: "Harm",
	RankFive: "Restore", RankSix: "Find", RankSeven: "Travel",
	RankEight: "Protect", RankNine: "Enrich Self", RankTen: "Avenge",
	RankJack: "Fulfill Duty", RankQueen: "Escape", RankKing: "Create",
	RankAce: "Serve",
}

var notableFeatures = [7]string{
	"",
	"Unremarkable",
	"Notable nature",
	"Obvious physical trait",
	"Quirk or mannerism",
	"Unusual equipment",
	"Unexpected age or origin",
}

func NPCGenerator(deck *Deck, rng *Randomizer) NPCResult {
	idDraw := deck.Draw(nil)
	identity := CardTableResult{
		Draw: idDraw, TableName: "Identity",
		Entry: identityTable[idDraw.Card.Rank],
	}

	goalDraw := deck.Draw(nil)
	goal := CardTableResult{
		Draw: goalDraw, TableName: "Goal",
		Entry: goalTable[goalDraw.Card.Rank],
	}

	featRoll := rng.RollD6()
	featDetail := DetailFocus(deck)
	attitude := OracleHow(rng)
	topic := TopicFocus(deck)

	return NPCResult{
		Identity:      identity,
		Goal:          goal,
		FeatureRoll:   featRoll,
		Feature:       notableFeatures[featRoll],
		FeatureDetail: featDetail,
		Attitude:      attitude,
		Topic:         topic,
	}
}
