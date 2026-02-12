package engine

var actionFocusTable = map[Rank]string{
	RankTwo: "Seek", RankThree: "Oppose", RankFour: "Communicate",
	RankFive: "Move", RankSix: "Harm", RankSeven: "Create",
	RankEight: "Reveal", RankNine: "Command", RankTen: "Take",
	RankJack: "Protect", RankQueen: "Assist", RankKing: "Transform",
	RankAce: "Deceive",
}

var detailFocusTable = map[Rank]string{
	RankTwo: "Small", RankThree: "Large", RankFour: "Old",
	RankFive: "New", RankSix: "Mundane", RankSeven: "Simple",
	RankEight: "Complex", RankNine: "Unsavory", RankTen: "Specialized",
	RankJack: "Unexpected", RankQueen: "Exotic", RankKing: "Dignified",
	RankAce: "Unique",
}

var topicFocusTable = map[Rank]string{
	RankTwo: "Current Need", RankThree: "Allies", RankFour: "Community",
	RankFive: "History", RankSix: "Future Plans", RankSeven: "Enemies",
	RankEight: "Knowledge", RankNine: "Rumors", RankTen: "A Plot Arc",
	RankJack: "Recent Events", RankQueen: "Equipment", RankKing: "A Faction",
	RankAce: "The PCs",
}

func drawFocus(deck *Deck, table map[Rank]string, tableName string) CardTableResult {
	draw := deck.Draw(nil)
	return CardTableResult{
		Draw:      draw,
		TableName: tableName,
		Entry:     table[draw.Card.Rank],
	}
}

func ActionFocus(deck *Deck) CardTableResult {
	return drawFocus(deck, actionFocusTable, "Action Focus")
}

func DetailFocus(deck *Deck) CardTableResult {
	return drawFocus(deck, detailFocusTable, "Detail Focus")
}

func TopicFocus(deck *Deck) CardTableResult {
	return drawFocus(deck, topicFocusTable, "Topic Focus")
}
