package engine

var pacingMoves = [7]string{
	"",
	"Foreshadow Trouble",
	"Reveal a New Detail",
	"An NPC Takes Action",
	"Advance a Threat",
	"Advance a Plot",
	"Add a RANDOM EVENT to the scene",
}

var failureMoves = [7]string{
	"",
	"Cause Harm",
	"Put Someone in a Spot",
	"Offer a Choice",
	"Advance a Threat",
	"Reveal an Unwelcome Truth",
	"Foreshadow Trouble",
}

func PacingMove(rng *Randomizer, deck *Deck) PacingMoveResult {
	roll := rng.RollD6()
	result := PacingMoveResult{Roll: roll, Result: pacingMoves[roll]}
	if roll == 6 {
		evt := RandomEvent(deck, rng)
		result.RandomEvent = &evt
	}
	return result
}

func FailureMove(rng *Randomizer) FailureMoveResult {
	roll := rng.RollD6()
	return FailureMoveResult{Roll: roll, Result: failureMoves[roll]}
}
