package engine

import "testing"

func TestSeededRandomizerDeterministic(t *testing.T) {
	r1 := NewSeededRandomizer(42, 0)
	r2 := NewSeededRandomizer(42, 0)
	for i := 0; i < 100; i++ {
		if r1.Intn(1000) != r2.Intn(1000) {
			t.Fatalf("seeded randomizers diverged at iteration %d", i)
		}
	}
}

func TestRollD6Range(t *testing.T) {
	rng := NewSeededRandomizer(1, 2)
	for i := 0; i < 1000; i++ {
		v := rng.RollD6()
		if v < 1 || v > 6 {
			t.Fatalf("RollD6 returned %d, want 1-6", v)
		}
	}
}

func TestRollDNRange(t *testing.T) {
	rng := NewSeededRandomizer(3, 4)
	for _, n := range []int{2, 4, 6, 8, 10, 12, 20, 100} {
		for i := 0; i < 500; i++ {
			v := rng.RollDN(n)
			if v < 1 || v > n {
				t.Fatalf("RollDN(%d) returned %d", n, v)
			}
		}
	}
}

func TestCoinFlipBool(t *testing.T) {
	rng := NewSeededRandomizer(5, 6)
	heads, tails := 0, 0
	for i := 0; i < 1000; i++ {
		if rng.CoinFlip() {
			heads++
		} else {
			tails++
		}
	}
	if heads == 0 || tails == 0 {
		t.Fatal("CoinFlip never produced one of the outcomes")
	}
}
