package ui

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func makeTestImage(w, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestResizeNearestNeighbor(t *testing.T) {
	src := makeTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	dst := ResizeNearestNeighbor(src, 10, 10)
	bounds := dst.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("expected 10x10, got %dx%d", bounds.Dx(), bounds.Dy())
	}
	r, g, b, _ := dst.At(5, 5).RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("expected red pixel, got (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestResizeBilinear(t *testing.T) {
	src := makeTestImage(100, 100, color.RGBA{128, 64, 200, 255})
	dst := ResizeBilinear(src, 10, 10)
	bounds := dst.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("expected 10x10, got %dx%d", bounds.Dx(), bounds.Dy())
	}
	// Solid color image should stay approximately the same color
	r, g, b, _ := dst.At(5, 5).RGBA()
	if diff(r>>8, 128) > 2 || diff(g>>8, 64) > 2 || diff(b>>8, 200) > 2 {
		t.Errorf("expected ~(128,64,200), got (%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func diff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestRenderPortraitArt(t *testing.T) {
	img := makeTestImage(portraitPixW, portraitPixH, color.RGBA{0, 128, 255, 255})
	art := RenderPortraitArt(img)
	lines := strings.Split(art, "\n")
	if len(lines) != portraitCellH {
		t.Errorf("expected %d lines, got %d", portraitCellH, len(lines))
	}
	hasBraille := false
	for _, r := range art {
		if r >= 0x2800 && r <= 0x28FF {
			hasBraille = true
			break
		}
	}
	if !hasBraille {
		t.Error("expected braille character in output")
	}
	if !strings.Contains(art, "\x1b[38;2;") {
		t.Error("expected true-color foreground escape")
	}
	if !strings.Contains(art, "\x1b[48;2;") {
		t.Error("expected true-color background escape")
	}
}

func TestRenderEmptyPortraitBox(t *testing.T) {
	box := RenderEmptyPortraitBox()
	if box == "" {
		t.Error("expected non-empty box")
	}
	if !strings.Contains(box, "╭") && !strings.Contains(box, "┌") {
		t.Error("expected rounded border character")
	}
}

func TestPortraitTotalWidth(t *testing.T) {
	w := PortraitTotalWidth()
	if w != portraitCellW+2 {
		t.Errorf("expected %d, got %d", portraitCellW+2, w)
	}
}

func TestMedianOf8(t *testing.T) {
	v := [8]float64{1, 8, 3, 6, 2, 7, 4, 5}
	got := medianOf8(v)
	want := 4.5 // (4+5)/2
	if got != want {
		t.Errorf("expected %f, got %f", want, got)
	}
}
