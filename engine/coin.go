package engine

func FlipCoins(rng *Randomizer, count int) CoinFlipResult {
	if count < 1 {
		count = 1
	}
	flips := make([]bool, count)
	heads, tails := 0, 0
	for i := range flips {
		flips[i] = rng.CoinFlip()
		if flips[i] {
			heads++
		} else {
			tails++
		}
	}
	return CoinFlipResult{Flips: flips, Heads: heads, Tails: tails}
}
