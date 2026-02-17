package ui

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/muesli/termenv"
)

const (
	// Portrait cell dimensions. Terminal cells are ~2:1 h:w, so W×(W/2) looks square.
	// 12×6 cells = 24×24 braille pixel resolution.
	portraitCellW = 12
	portraitCellH = 6

	// Braille pixels: each cell encodes a 2×4 pixel grid.
	portraitPixW = portraitCellW * 2 // 24 pixels wide
	portraitPixH = portraitCellH * 4 // 24 pixels tall
)

// Braille dot layout within a 2×4 pixel block:
//
//	(0,0)=dot1  (1,0)=dot4
//	(0,1)=dot2  (1,1)=dot5
//	(0,2)=dot3  (1,2)=dot6
//	(0,3)=dot7  (1,3)=dot8
var brailleBit = [2][4]rune{
	{0x01, 0x02, 0x04, 0x40}, // x=0: dots 1,2,3,7
	{0x08, 0x10, 0x20, 0x80}, // x=1: dots 4,5,6,8
}

// PortraitTotalWidth returns the total width of a bordered portrait box.
func PortraitTotalWidth() int {
	return portraitCellW + 2 // +2 for rounded border left+right
}

// SupportsPortraits checks if the terminal supports true color for braille art.
func SupportsPortraits() bool {
	return termenv.ColorProfile() == termenv.TrueColor
}

// ResizeNearestNeighbor scales src to w x h using nearest-neighbor interpolation.
func ResizeNearestNeighbor(src image.Image, w, h int) *image.RGBA {
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		srcY := bounds.Min.Y + y*srcH/h
		for x := 0; x < w; x++ {
			srcX := bounds.Min.X + x*srcW/w
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

// ResizeBilinear scales src to w x h using bilinear interpolation for smoother results.
func ResizeBilinear(src image.Image, w, h int) *image.RGBA {
	bounds := src.Bounds()
	srcW := float64(bounds.Dx())
	srcH := float64(bounds.Dy())
	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		sy := (float64(y)+0.5)*srcH/float64(h) - 0.5
		y0 := int(math.Floor(sy))
		y1 := y0 + 1
		fy := sy - float64(y0)
		y0 = clampInt(y0, 0, bounds.Dy()-1) + bounds.Min.Y
		y1 = clampInt(y1, 0, bounds.Dy()-1) + bounds.Min.Y

		for x := 0; x < w; x++ {
			sx := (float64(x)+0.5)*srcW/float64(w) - 0.5
			x0 := int(math.Floor(sx))
			x1 := x0 + 1
			fx := sx - float64(x0)
			x0 = clampInt(x0, 0, bounds.Dx()-1) + bounds.Min.X
			x1 = clampInt(x1, 0, bounds.Dx()-1) + bounds.Min.X

			r00, g00, b00, _ := src.At(x0, y0).RGBA()
			r01, g01, b01, _ := src.At(x1, y0).RGBA()
			r10, g10, b10, _ := src.At(x0, y1).RGBA()
			r11, g11, b11, _ := src.At(x1, y1).RGBA()

			r := bilinear(float64(r00), float64(r01), float64(r10), float64(r11), fx, fy)
			g := bilinear(float64(g00), float64(g01), float64(g10), float64(g11), fx, fy)
			b := bilinear(float64(b00), float64(b01), float64(b10), float64(b11), fx, fy)

			dst.Set(x, y, color.RGBA{uint8(r / 256), uint8(g / 256), uint8(b / 256), 255})
		}
	}
	return dst
}

func bilinear(v00, v01, v10, v11, fx, fy float64) float64 {
	return v00*(1-fx)*(1-fy) + v01*fx*(1-fy) + v10*(1-fx)*fy + v11*fx*fy
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// luminance returns perceptual brightness of an RGB color (0.0–1.0).
func luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 65535.0
}

// avgColor returns the average RGB of the given colors.
func avgColor(colors []color.Color) (uint8, uint8, uint8) {
	if len(colors) == 0 {
		return 0, 0, 0
	}
	var rSum, gSum, bSum uint64
	for _, c := range colors {
		r, g, b, _ := c.RGBA()
		rSum += uint64(r)
		gSum += uint64(g)
		bSum += uint64(b)
	}
	n := uint64(len(colors))
	return uint8(rSum / n >> 8), uint8(gSum / n >> 8), uint8(bSum / n >> 8)
}

// RenderPortraitArt converts an image to braille art using true-color ANSI escapes
// and bilinear downsampling. Each cell represents a 2×4 pixel block split by luminance
// into foreground (braille dots) and background, giving 48×48 effective resolution.
func RenderPortraitArt(img image.Image) string {
	resized := ResizeBilinear(img, portraitPixW, portraitPixH)
	var b strings.Builder
	for cellY := 0; cellY < portraitCellH; cellY++ {
		if cellY > 0 {
			b.WriteByte('\n')
		}
		for cellX := 0; cellX < portraitCellW; cellX++ {
			var pixels [2][4]color.Color
			var lums [8]float64
			i := 0
			for x := 0; x < 2; x++ {
				for y := 0; y < 4; y++ {
					px := resized.At(cellX*2+x, cellY*4+y)
					pixels[x][y] = px
					lums[i] = luminance(px)
					i++
				}
			}

			median := medianOf8(lums)

			var pattern rune
			var fgPixels, bgPixels []color.Color
			for x := 0; x < 2; x++ {
				for y := 0; y < 4; y++ {
					if luminance(pixels[x][y]) >= median {
						pattern |= brailleBit[x][y]
						fgPixels = append(fgPixels, pixels[x][y])
					} else {
						bgPixels = append(bgPixels, pixels[x][y])
					}
				}
			}

			fgR, fgG, fgB := avgColor(fgPixels)
			bgR, bgG, bgB := avgColor(bgPixels)
			if len(bgPixels) == 0 {
				bgR, bgG, bgB = fgR, fgG, fgB
			}

			fmt.Fprintf(&b, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%c",
				fgR, fgG, fgB, bgR, bgG, bgB, 0x2800+pattern)
		}
		b.WriteString("\x1b[0m")
	}
	return b.String()
}

// medianOf8 returns the median of exactly 8 values (average of 4th and 5th).
func medianOf8(v [8]float64) float64 {
	sorted := v
	for i := 1; i < 8; i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && sorted[j] > key {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}
	return (sorted[3] + sorted[4]) / 2
}

// RenderEmptyPortraitBox returns a bordered box of the same size as a portrait,
// filled with spaces.
func RenderEmptyPortraitBox() string {
	row := strings.Repeat(" ", portraitCellW)
	rows := make([]string, portraitCellH)
	for i := range rows {
		rows[i] = row
	}
	return PortraitBorderStyle.Render(strings.Join(rows, "\n"))
}
