package engine

import "testing"

func TestGenerateRandomPortrait(t *testing.T) {
	rng := NewRandomizer()
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		p := GenerateRandomPortrait(rng)
		if p.FaceShape >= faceShapeCount {
			t.Errorf("face shape %d out of range", p.FaceShape)
		}
		if p.EyeStyle >= eyeStyleCount {
			t.Errorf("eye style %d out of range", p.EyeStyle)
		}
		if p.HairStyle >= hairCount {
			t.Errorf("hair style %d out of range", p.HairStyle)
		}
		if p.Accessory >= accCount {
			t.Errorf("accessory %d out of range", p.Accessory)
		}
		// Track unique structural combos
		key := DescribePortrait(p)
		seen[key] = true
	}
	// With 50 random draws from 36k+ combos, we should see variety
	if len(seen) < 5 {
		t.Errorf("expected variety, only saw %d unique descriptions", len(seen))
	}
}

func TestRenderPortraitImage(t *testing.T) {
	rng := NewRandomizer()
	for i := 0; i < 20; i++ {
		p := GenerateRandomPortrait(rng)
		img := RenderPortraitImage(p)
		bounds := img.Bounds()
		if bounds.Dx() != pgSize || bounds.Dy() != pgSize {
			t.Errorf("expected %dx%d, got %dx%d", pgSize, pgSize, bounds.Dx(), bounds.Dy())
		}
	}
}

func TestRenderPortraitImageAllTraitCombos(t *testing.T) {
	// Test boundary values for each trait
	base := PortraitParams{
		BG: [3]uint8{40, 40, 60}, Skin: [3]uint8{200, 160, 130},
		HairColor: [3]uint8{100, 60, 30}, EyeColor: [3]uint8{60, 100, 60},
		AccColor: [3]uint8{140, 140, 150},
	}

	for face := uint8(0); face < faceShapeCount; face++ {
		for hair := uint8(0); hair < hairCount; hair++ {
			for acc := uint8(0); acc < accCount; acc++ {
				p := base
				p.FaceShape = face
				p.HairStyle = hair
				p.Accessory = acc
				img := RenderPortraitImage(p)
				if img.Bounds().Dx() != pgSize {
					t.Errorf("face=%d hair=%d acc=%d: wrong size", face, hair, acc)
				}
			}
		}
	}
}

func TestDescribePortrait(t *testing.T) {
	p := PortraitParams{FaceShape: 0, HairStyle: 2}
	got := DescribePortrait(p)
	if got != "Round, Medium hair" {
		t.Errorf("expected 'Round, Medium hair', got %q", got)
	}
}
