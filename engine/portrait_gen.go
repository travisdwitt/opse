package engine

import (
	"fmt"
	"image"
	"image/color"
)

const pgSize = 48

// PortraitParams holds the serializable parameters for a procedural portrait.
// Each trait index selects a variant; colors are fully random.
type PortraitParams struct {
	FaceShape  uint8    `json:"face"`
	EyeStyle   uint8    `json:"eyes"`
	NoseStyle  uint8    `json:"nose"`
	MouthStyle uint8    `json:"mouth"`
	HairStyle  uint8    `json:"hair"`
	Accessory  uint8    `json:"acc"`
	BodyStyle  uint8    `json:"body"`
	BG         [3]uint8 `json:"bg"`
	Skin       [3]uint8 `json:"skin"`
	HairColor  [3]uint8 `json:"hair_color"`
	EyeColor   [3]uint8 `json:"eye_color"`
	AccColor   [3]uint8 `json:"acc_color"`
}

const (
	faceShapeCount = 4 // round, oval, square, long
	eyeStyleCount  = 5 // round, narrow, large, dots, asymmetric
	noseStyleCount = 4 // dot, triangle, wide, long
	mouthCount     = 4 // smile, line, open, wide
	hairCount      = 9 // bald, short, medium, long, mohawk, spiky, afro, side-swept, buzz
	accCount       = 8 // none, glasses, hat, headband, beard, scarf, collar, hood
	bodyCount      = 3 // narrow, medium, broad
)

// GenerateRandomPortrait creates a fully random portrait.
func GenerateRandomPortrait(rng *Randomizer) PortraitParams {
	return PortraitParams{
		FaceShape:  uint8(rng.Intn(faceShapeCount)),
		EyeStyle:   uint8(rng.Intn(eyeStyleCount)),
		NoseStyle:  uint8(rng.Intn(noseStyleCount)),
		MouthStyle: uint8(rng.Intn(mouthCount)),
		HairStyle:  uint8(rng.Intn(hairCount)),
		Accessory:  uint8(rng.Intn(accCount)),
		BodyStyle:  uint8(rng.Intn(bodyCount)),
		BG:         randDarkColor(rng),
		Skin:       randSkinColor(rng),
		HairColor:  randBrightColor(rng),
		EyeColor:   randBrightColor(rng),
		AccColor:   randBrightColor(rng),
	}
}

func randDarkColor(rng *Randomizer) [3]uint8 {
	return [3]uint8{
		uint8(20 + rng.Intn(60)),
		uint8(20 + rng.Intn(60)),
		uint8(20 + rng.Intn(60)),
	}
}

func randSkinColor(rng *Randomizer) [3]uint8 {
	// Wide range: realistic warm tones to unrealistic colors
	return [3]uint8{
		uint8(80 + rng.Intn(176)),
		uint8(50 + rng.Intn(170)),
		uint8(30 + rng.Intn(170)),
	}
}

func randBrightColor(rng *Randomizer) [3]uint8 {
	return [3]uint8{
		uint8(30 + rng.Intn(210)),
		uint8(30 + rng.Intn(210)),
		uint8(30 + rng.Intn(210)),
	}
}

// DescribePortrait returns a short human-readable label for a portrait.
func DescribePortrait(p PortraitParams) string {
	faces := []string{"Round", "Oval", "Square", "Long"}
	hairs := []string{"Bald", "Short hair", "Medium hair", "Long hair",
		"Mohawk", "Spiky", "Afro", "Side-swept", "Buzz cut"}
	face := "Face"
	if int(p.FaceShape) < len(faces) {
		face = faces[p.FaceShape]
	}
	hair := ""
	if int(p.HairStyle) < len(hairs) {
		hair = hairs[p.HairStyle]
	}
	return fmt.Sprintf("%s, %s", face, hair)
}

// RenderPortraitImage draws a 48x48 portrait from the given parameters.
func RenderPortraitImage(params PortraitParams) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, pgSize, pgSize))
	bg := pgRgb(params.BG[0], params.BG[1], params.BG[2])
	skin := pgRgb(params.Skin[0], params.Skin[1], params.Skin[2])
	hairC := pgRgb(params.HairColor[0], params.HairColor[1], params.HairColor[2])
	eyeC := pgRgb(params.EyeColor[0], params.EyeColor[1], params.EyeColor[2])
	accC := pgRgb(params.AccColor[0], params.AccColor[1], params.AccColor[2])

	pgFillRect(img, 0, 0, pgSize, pgSize, bg)
	pgDrawBody(img, skin, accC, params.BodyStyle)
	pgDrawNeck(img, skin)
	pgDrawFace(img, skin, params.FaceShape)
	pgDrawHair(img, hairC, params.HairStyle, params.FaceShape)
	pgDrawEyes(img, eyeC, skin, params.EyeStyle)
	pgDrawNose(img, skin, params.NoseStyle)
	pgDrawMouth(img, skin, params.MouthStyle)
	pgDrawAccessory(img, accC, hairC, skin, params.Accessory, params.FaceShape)
	return img
}

// --- Drawing primitives ---

func pgRgb(r, g, b uint8) color.RGBA { return color.RGBA{r, g, b, 255} }

func pgDarker(c color.RGBA, amt uint8) color.RGBA {
	sub := func(v, a uint8) uint8 {
		if v < a {
			return 0
		}
		return v - a
	}
	return color.RGBA{sub(c.R, amt), sub(c.G, amt), sub(c.B, amt), 255}
}

func pgLighter(c color.RGBA, amt uint8) color.RGBA {
	add := func(v, a uint8) uint8 {
		if int(v)+int(a) > 255 {
			return 255
		}
		return v + a
	}
	return color.RGBA{add(c.R, amt), add(c.G, amt), add(c.B, amt), 255}
}

func pgSet(img *image.RGBA, x, y int, c color.RGBA) {
	if x >= 0 && x < pgSize && y >= 0 && y < pgSize {
		img.Set(x, y, c)
	}
}

func pgFillRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			pgSet(img, x+dx, y+dy, c)
		}
	}
}

func pgFillCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx, dy := float64(x-cx), float64(y-cy)
			if dx*dx+dy*dy <= float64(r*r) {
				pgSet(img, x, y, c)
			}
		}
	}
}

func pgFillOval(img *image.RGBA, cx, cy, rx, ry int, c color.RGBA) {
	for y := cy - ry; y <= cy+ry; y++ {
		for x := cx - rx; x <= cx+rx; x++ {
			dx := float64(x-cx) / float64(max(rx, 1))
			dy := float64(y-cy) / float64(max(ry, 1))
			if dx*dx+dy*dy <= 1.0 {
				pgSet(img, x, y, c)
			}
		}
	}
}

func pgHLine(img *image.RGBA, x1, x2, y int, c color.RGBA) {
	for x := x1; x <= x2; x++ {
		pgSet(img, x, y, c)
	}
}

func pgFillTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, c color.RGBA) {
	minX := min(x1, min(x2, x3))
	maxX := max(x1, max(x2, x3))
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			d1 := (x-x2)*(y1-y2) - (x1-x2)*(y-y2)
			d2 := (x-x3)*(y2-y3) - (x2-x3)*(y-y3)
			d3 := (x-x1)*(y3-y1) - (x3-x1)*(y-y1)
			hasNeg := d1 < 0 || d2 < 0 || d3 < 0
			hasPos := d1 > 0 || d2 > 0 || d3 > 0
			if !(hasNeg && hasPos) {
				pgSet(img, x, y, c)
			}
		}
	}
}

// --- Feature drawing functions ---

const cx = pgSize / 2 // horizontal center = 24

func pgDrawNeck(img *image.RGBA, skin color.RGBA) {
	pgFillRect(img, cx-4, 32, 8, 8, skin)
}

func pgDrawBody(img *image.RGBA, skin, accC color.RGBA, style uint8) {
	switch style {
	case 0: // narrow
		pgFillOval(img, cx, 44, 10, 6, accC)
	case 1: // medium
		pgFillOval(img, cx, 44, 14, 7, accC)
	default: // broad
		pgFillOval(img, cx, 44, 18, 8, accC)
	}
}

func pgDrawFace(img *image.RGBA, skin color.RGBA, shape uint8) {
	switch shape {
	case 0: // round
		pgFillCircle(img, cx, 22, 12, skin)
	case 1: // oval
		pgFillOval(img, cx, 22, 10, 13, skin)
	case 2: // square
		pgFillRect(img, cx-10, 11, 20, 22, skin)
		// round the corners slightly
		pgFillCircle(img, cx-8, 13, 3, skin)
		pgFillCircle(img, cx+8, 13, 3, skin)
		pgFillCircle(img, cx-8, 30, 3, skin)
		pgFillCircle(img, cx+8, 30, 3, skin)
	default: // long
		pgFillOval(img, cx, 22, 9, 14, skin)
	}
}

func pgDrawEyes(img *image.RGBA, eyeC, skin color.RGBA, style uint8) {
	white := pgRgb(240, 240, 240)
	pupil := pgRgb(20, 20, 20)
	sep := 5
	ey := 20

	switch style {
	case 0: // round
		pgFillCircle(img, cx-sep, ey, 3, white)
		pgFillCircle(img, cx+sep, ey, 3, white)
		pgFillCircle(img, cx-sep, ey, 2, eyeC)
		pgFillCircle(img, cx+sep, ey, 2, eyeC)
		pgFillCircle(img, cx-sep, ey, 1, pupil)
		pgFillCircle(img, cx+sep, ey, 1, pupil)
	case 1: // narrow
		pgFillRect(img, cx-sep-2, ey, 4, 2, white)
		pgFillRect(img, cx+sep-2, ey, 4, 2, white)
		pgSet(img, cx-sep, ey, eyeC)
		pgSet(img, cx+sep, ey, eyeC)
	case 2: // large
		pgFillOval(img, cx-sep, ey, 3, 2, white)
		pgFillOval(img, cx+sep, ey, 3, 2, white)
		pgFillCircle(img, cx-sep, ey, 2, eyeC)
		pgFillCircle(img, cx+sep, ey, 2, eyeC)
		pgFillCircle(img, cx-sep, ey, 1, pupil)
		pgFillCircle(img, cx+sep, ey, 1, pupil)
		// highlight
		pgSet(img, cx-sep+1, ey-1, white)
		pgSet(img, cx+sep+1, ey-1, white)
	case 3: // dots
		pgFillCircle(img, cx-sep, ey, 1, pupil)
		pgFillCircle(img, cx+sep, ey, 1, pupil)
	default: // asymmetric
		pgFillCircle(img, cx-sep, ey, 3, white)
		pgFillCircle(img, cx+sep, ey, 2, white)
		pgFillCircle(img, cx-sep, ey, 2, eyeC)
		pgFillCircle(img, cx+sep, ey, 1, eyeC)
		pgFillCircle(img, cx-sep, ey, 1, pupil)
		pgSet(img, cx+sep, ey, pupil)
	}
}

func pgDrawNose(img *image.RGBA, skin color.RGBA, style uint8) {
	d := pgDarker(skin, 25)
	ny := 25
	switch style {
	case 0: // dot
		pgFillCircle(img, cx, ny, 1, d)
	case 1: // triangle
		pgFillTriangle(img, cx, ny-2, cx-2, ny+2, cx+2, ny+2, d)
	case 2: // wide
		pgFillOval(img, cx, ny, 3, 2, d)
	default: // long
		pgFillRect(img, cx-1, ny-2, 2, 5, d)
	}
}

func pgDrawMouth(img *image.RGBA, skin color.RGBA, style uint8) {
	d := pgDarker(skin, 35)
	my := 30
	switch style {
	case 0: // smile
		for x := -3; x <= 3; x++ {
			curve := x * x / 6
			pgSet(img, cx+x, my+curve, d)
			pgSet(img, cx+x, my+curve+1, d)
		}
	case 1: // line
		pgHLine(img, cx-3, cx+3, my, d)
	case 2: // open
		pgFillOval(img, cx, my, 2, 2, d)
	default: // wide grin
		pgHLine(img, cx-5, cx+5, my, d)
		pgSet(img, cx-5, my-1, d)
		pgSet(img, cx+5, my-1, d)
	}
}

func pgDrawHair(img *image.RGBA, hairC color.RGBA, style uint8, faceShape uint8) {
	// hairline y depends on face shape
	topY := 10
	switch faceShape {
	case 1: // oval
		topY = 9
	case 2: // square
		topY = 11
	case 3: // long
		topY = 8
	}

	switch style {
	case 0: // bald — no hair
	case 1: // short
		pgFillOval(img, cx, topY, 12, 5, hairC)
	case 2: // medium
		pgFillOval(img, cx, topY, 12, 6, hairC)
		pgFillRect(img, cx-12, topY, 4, 16, hairC)
		pgFillRect(img, cx+8, topY, 4, 16, hairC)
	case 3: // long
		pgFillOval(img, cx, topY, 12, 6, hairC)
		pgFillRect(img, cx-12, topY, 4, 28, hairC)
		pgFillRect(img, cx+8, topY, 4, 28, hairC)
	case 4: // mohawk
		pgFillRect(img, cx-2, topY-6, 4, 10, hairC)
		pgFillRect(img, cx-3, topY-4, 6, 4, hairC)
	case 5: // spiky
		for _, xOff := range []int{-8, -4, 0, 4, 8} {
			pgFillTriangle(img, cx+xOff-2, topY+2, cx+xOff, topY-5, cx+xOff+2, topY+2, hairC)
		}
	case 6: // afro
		pgFillCircle(img, cx, topY-2, 14, hairC)
	case 7: // side-swept
		pgFillOval(img, cx-3, topY, 13, 6, hairC)
		pgFillRect(img, cx-14, topY, 5, 12, hairC)
	default: // buzz cut
		pgFillOval(img, cx, topY+1, 11, 4, hairC)
	}
}

func pgDrawAccessory(img *image.RGBA, accC, hairC, skin color.RGBA, acc uint8, faceShape uint8) {
	switch acc {
	case 0: // none
	case 1: // glasses
		frame := pgDarker(accC, 30)
		pgFillCircle(img, cx-5, 20, 3, frame)
		pgFillCircle(img, cx+5, 20, 3, frame)
		// clear the lens interior
		pgFillCircle(img, cx-5, 20, 2, pgLighter(accC, 40))
		pgFillCircle(img, cx+5, 20, 2, pgLighter(accC, 40))
		// bridge
		pgHLine(img, cx-2, cx+2, 19, frame)
	case 2: // hat
		pgFillRect(img, cx-10, 5, 20, 5, accC)
		pgFillRect(img, cx-14, 9, 28, 3, accC)
	case 3: // headband
		pgHLine(img, cx-11, cx+11, 14, accC)
		pgHLine(img, cx-11, cx+11, 15, accC)
	case 4: // beard
		pgFillOval(img, cx, 32, 7, 6, hairC)
		pgFillRect(img, cx-2, 33, 1, 5, pgDarker(hairC, 15))
		pgFillRect(img, cx+1, 33, 1, 5, pgDarker(hairC, 15))
	case 5: // scarf
		pgFillRect(img, cx-10, 34, 20, 5, accC)
		pgFillRect(img, cx-11, 35, 2, 8, accC)
	case 6: // collar/tie
		pgFillTriangle(img, cx, 36, cx-4, 46, cx+4, 46, accC)
		pgFillRect(img, cx-1, 34, 2, 4, accC)
	default: // hood
		pgFillOval(img, cx, 16, 16, 14, accC)
		// cut out face area
		pgDrawFace(img, skin, faceShape)
	}
}
