package engine

func RandomEvent(deck *Deck, rng *Randomizer) RandomEventResult {
	eventFn := func() *RandomEventResult {
		evt := RandomEvent(deck, rng)
		return &evt
	}

	actionDraw := deck.Draw(eventFn)
	action := CardTableResult{
		Draw:      actionDraw,
		TableName: "Action Focus",
		Entry:     actionFocusTable[actionDraw.Card.Rank],
	}

	topicDraw := deck.Draw(eventFn)
	topic := CardTableResult{
		Draw:      topicDraw,
		TableName: "Topic Focus",
		Entry:     topicFocusTable[topicDraw.Card.Rank],
	}

	return RandomEventResult{Action: action, Topic: topic}
}
