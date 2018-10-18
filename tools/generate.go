// Script I used for quickly generating test images. Not really meant for
// public consumption but may be useful if you're looking to quickly generate
// an image of a certain size with all RGBA values the same.

package main

import "os"
import "image"
import "image/png"
import "image/color"

func main() {
	w, h := 2,2
	upLeft := image.Point{0,0}
	lowRight := image.Point{w,h}
	im := image.NewRGBA(image.Rectangle{upLeft,lowRight})
	for y := 0; y < h; y++ {
		for x := 0; x < h; x++ {
			pix := color.NRGBA64{0,0,65535,40000}
			im.Set(x,y,pix)
		}
	}
	newFile, _ := os.Create("images/im2.png")
	png.Encode(newFile, im)
}
