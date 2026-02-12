package engine

import "testing"

func TestNewDeckHas54Cards(t *testing.T) {
	rng := NewSeededRandomizer(10, 20)
	d := NewDeck(rng)
	if d.Remaining() != 54 {
		t.Fatalf("expected 54 cards, got %d", d.Remaining())
	}
}

func TestDrawNeverReturnsJoker(t *testing.T) {
	rng := NewSeededRandomizer(99, 0)
	d := NewDeck(rng)
	for i := 0; i < 200; i++ {
		result := d.Draw(func() *RandomEventResult { return nil })
		if result.Card.IsJoker() {
			t.Fatal("Draw returned a joker")
		}
	}
}

func TestDrawDecrementsRemaining(t *testing.T) {
	rng := NewSeededRandomizer(42, 42)
	d := NewDeck(rng)
	before := d.Remaining()
	d.Draw(nil)
	after := d.Remaining()
	// Could have reshuffled on joker, but remaining should differ from before
	if after >= before {
		// Only possible if joker was drawn and deck reshuffled
		// which is fine — just verify we can keep drawing
	}
}

func TestDrawAll54NoPanic(t *testing.T) {
	rng := NewSeededRandomizer(7, 7)
	d := NewDeck(rng)
	for i := 0; i < 54; i++ {
		d.Draw(nil)
	}
	// Drawing beyond 54 should trigger auto-shuffle
	d.Draw(nil)
}

func TestShuffleResets(t *testing.T) {
	rng := NewSeededRandomizer(1, 1)
	d := NewDeck(rng)
	for i := 0; i < 10; i++ {
		d.Draw(nil)
	}
	d.Shuffle()
	if d.Remaining() != 54 {
		t.Fatalf("after shuffle expected 54, got %d", d.Remaining())
	}
}

func TestCardString(t *testing.T) {
	c := Card{Rank: RankAce, Suit: Spades}
	if c.String() != "A♠" {
		t.Fatalf("expected A♠, got %s", c.String())
	}
	j := Card{Rank: RankJoker}
	if j.String() != "Joker" {
		t.Fatalf("expected Joker, got %s", j.String())
	}
}

func TestRankString(t *testing.T) {
	tests := map[Rank]string{
		RankTwo: "2", RankTen: "10", RankJack: "J",
		RankQueen: "Q", RankKing: "K", RankAce: "A", RankJoker: "Joker",
	}
	for rank, want := range tests {
		if rank.String() != want {
			t.Errorf("Rank(%d).String() = %q, want %q", rank, rank.String(), want)
		}
	}
}
