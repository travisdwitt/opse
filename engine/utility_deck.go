package engine

type UtilityDeck struct {
	cards         []Card
	includeJokers bool
	rng           *Randomizer
}

func NewUtilityDeck(rng *Randomizer, jokers bool) *UtilityDeck {
	d := &UtilityDeck{rng: rng, includeJokers: jokers}
	d.Shuffle()
	return d
}

func (d *UtilityDeck) Shuffle() {
	d.cards = make([]Card, 0, 54)
	for _, s := range []Suit{Clubs, Diamonds, Spades, Hearts} {
		for r := RankTwo; r <= RankAce; r++ {
			d.cards = append(d.cards, Card{Rank: r, Suit: s})
		}
	}
	if d.includeJokers {
		d.cards = append(d.cards, Card{Rank: RankJoker}, Card{Rank: RankJoker})
	}
	for i := len(d.cards) - 1; i > 0; i-- {
		j := d.rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *UtilityDeck) Draw(count int) CardDrawResult {
	if count < 1 {
		count = 1
	}
	if count > len(d.cards) {
		count = len(d.cards)
	}
	drawn := make([]Card, count)
	copy(drawn, d.cards[:count])
	d.cards = d.cards[count:]
	return CardDrawResult{Cards: drawn, Remaining: len(d.cards)}
}

func (d *UtilityDeck) Remaining() int { return len(d.cards) }
