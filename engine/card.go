package engine

import "fmt"

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Spades
	Hearts
)

var SuitSymbols = map[Suit]string{
	Clubs: "♣", Diamonds: "♦", Spades: "♠", Hearts: "♥",
}

var SuitNames = map[Suit]string{
	Clubs: "Clubs", Diamonds: "Diamonds", Spades: "Spades", Hearts: "Hearts",
}

var SuitDomains = map[Suit]string{
	Clubs:    "Physical (appearance, existence)",
	Diamonds: "Technical (mental, operation)",
	Spades:   "Mystical (meaning, capability)",
	Hearts:   "Social (personal, connection)",
}

type Rank int

const (
	RankTwo   Rank = 2 + iota
	RankThree      // 3
	RankFour       // 4
	RankFive       // 5
	RankSix        // 6
	RankSeven      // 7
	RankEight      // 8
	RankNine       // 9
	RankTen        // 10
	RankJack       // 11
	RankQueen      // 12
	RankKing       // 13
	RankAce        // 14
	RankJoker = 99
)

func (r Rank) String() string {
	switch {
	case r >= 2 && r <= 10:
		return fmt.Sprintf("%d", int(r))
	case r == RankJack:
		return "J"
	case r == RankQueen:
		return "Q"
	case r == RankKing:
		return "K"
	case r == RankAce:
		return "A"
	case r == RankJoker:
		return "Joker"
	default:
		return "?"
	}
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) IsJoker() bool { return c.Rank == RankJoker }

func (c Card) String() string {
	if c.IsJoker() {
		return "Joker"
	}
	return c.Rank.String() + SuitSymbols[c.Suit]
}

func (c Card) Domain() string {
	return SuitDomains[c.Suit]
}

type DrawResult struct {
	Card       Card
	JokerDrawn bool
	JokerEvent *RandomEventResult
}

type Deck struct {
	cards []Card
	rng   *Randomizer
}

func NewDeck(rng *Randomizer) *Deck {
	d := &Deck{rng: rng}
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	d.cards = make([]Card, 0, 54)
	suits := []Suit{Clubs, Diamonds, Spades, Hearts}
	for _, s := range suits {
		for r := RankTwo; r <= RankAce; r++ {
			d.cards = append(d.cards, Card{Rank: r, Suit: s})
		}
	}
	d.cards = append(d.cards, Card{Rank: RankJoker}, Card{Rank: RankJoker})
	for i := len(d.cards) - 1; i > 0; i-- {
		j := d.rng.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

func (d *Deck) Draw(eventFn func() *RandomEventResult) DrawResult {
	if len(d.cards) == 0 {
		d.Shuffle()
	}
	card := d.cards[0]
	d.cards = d.cards[1:]

	if card.IsJoker() {
		d.Shuffle()
		var evt *RandomEventResult
		if eventFn != nil {
			evt = eventFn()
		}
		result := d.Draw(eventFn)
		result.JokerDrawn = true
		result.JokerEvent = evt
		return result
	}
	return DrawResult{Card: card}
}

func (d *Deck) Remaining() int { return len(d.cards) }
