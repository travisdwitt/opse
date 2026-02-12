package engine

import "math/rand/v2"

type Randomizer struct {
	rng *rand.Rand
}

func NewRandomizer() *Randomizer {
	return &Randomizer{rng: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))}
}

func NewSeededRandomizer(seed1, seed2 uint64) *Randomizer {
	return &Randomizer{rng: rand.New(rand.NewPCG(seed1, seed2))}
}

func (r *Randomizer) Intn(n int) int    { return r.rng.IntN(n) }
func (r *Randomizer) RollD6() int       { return r.Intn(6) + 1 }
func (r *Randomizer) RollD4() int       { return r.Intn(4) + 1 }
func (r *Randomizer) RollD12() int      { return r.Intn(12) + 1 }
func (r *Randomizer) RollDN(n int) int  { return r.Intn(n) + 1 }
func (r *Randomizer) CoinFlip() bool    { return r.Intn(2) == 0 }
