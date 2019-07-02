// Script I used for quickly generating test images. Not really meant for
// public consumption but may be useful if you're looking to quickly generate
// an image of a certain size with all RGBA values the same.

package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/nicolashahn/diffimg-go/imgutil"
)

func main() {
	size := image.Point{2, 2}
	m := image.NewRGBA(image.Rectangle{image.ZP, size})
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			pix := color.NRGBA64{0, 0, 65535, 40000}
			m.Set(x, y, pix)
		}
	}

	outfile := "images/im2.png"
	err := imgutil.WritePNG(outfile, m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %q: %v\n", outfile, err)
	}
}
