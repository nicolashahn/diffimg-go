// Package diffimg provides image diffing tools. You can create a difference
// image from two given images, or quantify their difference as a ratio/%.
package diffimg

import (
	"fmt"
	"image"
	"image/color"
)

// RatioRGBA calculates absolute difference ratio in RGBA colorspace and the difference image.
func RatioRGBA(a, b image.Image, ignoreAlpha bool) (float64, error) {
	if a.Bounds().Size() != b.Bounds().Size() {
		return 0, fmt.Errorf("images are different size %v and %v", a.Bounds().Size(), b.Bounds().Size())
	}

	var totalColor uint64
	var totalAlpha uint64

	size := a.Bounds().Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			argba := rgba(a.At(x, y).RGBA())
			brgba := rgba(b.At(x, y).RGBA())
			drgba := absdiff(argba, brgba)

			totalColor += uint64(drgba.R) + uint64(drgba.G) + uint64(drgba.B)
			totalAlpha += uint64(drgba.A)
		}
	}

	if ignoreAlpha {
		return float64(totalColor) / float64(size.X*size.Y*3*255), nil
	}
	return float64(totalAlpha+totalColor) / float64(size.X*size.Y*4*255), nil
}

// RGBA calculates absolute difference ratio in RGBA colorspace and the difference image.
func RGBA(a, b image.Image, ignoreAlpha bool) (image.Image, float64, error) {
	if a.Bounds().Size() != b.Bounds().Size() {
		return nil, 0, fmt.Errorf("images are different size %v and %v", a.Bounds().Size(), b.Bounds().Size())
	}

	m := image.NewRGBA(a.Bounds())
	var totalColor uint64
	var totalAlpha uint64

	size := a.Bounds().Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			argba := rgba(a.At(x, y).RGBA())
			brgba := rgba(b.At(x, y).RGBA())
			drgba := absdiff(argba, brgba)

			totalColor += uint64(drgba.R) + uint64(drgba.G) + uint64(drgba.B)
			totalAlpha += uint64(drgba.A)

			if ignoreAlpha {
				drgba.A = 0xFF
			}
			m.SetRGBA(x, y, drgba)
		}
	}

	if ignoreAlpha {
		return m, float64(totalColor) / float64(size.X*size.Y*3*255), nil
	}
	return m, float64(totalAlpha+totalColor) / float64(size.X*size.Y*4*255), nil
}

// rgba converts uint16 to color.RGBA
func rgba(r, g, b, a uint32) color.RGBA {
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// absdiff calculates absolute difference for colors
func absdiff(a, b color.RGBA) color.RGBA {
	return color.RGBA{
		R: absdiff8(a.R, b.R),
		G: absdiff8(a.G, b.G),
		B: absdiff8(a.B, b.B),
		A: absdiff8(a.A, b.A),
	}
}

// absdiff8 calculates absolute diff for uint8
func absdiff8(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
