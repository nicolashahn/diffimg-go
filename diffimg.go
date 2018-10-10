package main

import (
  "image"
  "fmt"
  "os"
)
import _ "image/png"

// Helper functions

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

func Abs(x int) uint32 {
	if x < 0 {
		return uint32(-x)
	}
	return uint32(x)
}

// Main diffimg functions

func checkArgs() {
  if len(os.Args) != 3 {
    panic("must pass exactly two arguments (image1 and image2)")
  }
}

func loadImage(filepath string) image.Image {
  file, err := os.Open(filepath)
  checkErr(err)
  defer file.Close()
  im, _, err := image.Decode(file)
  checkErr(err)
  return im
}

func checkDimensions(im1, im2 image.Image) {
  if im1.Bounds() != im2.Bounds() {
    panic("image dimensions are different")
  }
}

// TODO Convert 2nd image color model to 1st
func checkColorModel(im1, im2 image.Image) {
  if im1.ColorModel() != im2.ColorModel() {
    panic("color models are different")
  }
}

// Return the channel information at (x,y) as a uint32 array {r g b a}
func rgbaArray(im image.Image, x, y int) [4]uint32 {
  px := im.At(x,y)
  r, g, b, a := px.RGBA()
  return [4]uint32{r, g, b, a}
}

// Absolute difference between the channel values of the pixels at the same
// coordinates (x,y) in im1, im2
// Example (real channel values max at 65535, not 255 though):
// RGBA for im1 at (x,y): (100, 100, 180, 255)
// RGBA for im2 at (x,y): (120, 100, 100, 255)
// abs(120-100) + abs(100-100) + abs(180-100) + abs(255-255)
// 20 + 0 + 80 + 0 = 100
// return 100
func pixelDiff(im1, im2 image.Image, x, y int) uint32 {
  rgba1 := rgbaArray(im1, x, y)
  rgba2 := rgbaArray(im2, x, y)
  var pixDiffVal uint32
  for i, _ := range rgba1 {
    chanDiff := Abs(int(rgba1[i]) - int(rgba2[i]))
    pixDiffVal += chanDiff
  }
  return pixDiffVal
}

// Calculate difference ratio between two Images
// Adds up all the differences in each pixel's channel values, and averages
// over all pixels
func DiffImages(im1, im2 image.Image) float64 {
  var sumPixDiffVals uint32
  width := im1.Bounds().Max.X
  height := im2.Bounds().Max.Y
  for y := 0; y < height; y++ {
    for x := 0; x < width; x++ {
      sumPixDiffVals += pixelDiff(im1, im2, x, y)
    }
  }
  // Sum of max channel values for all pixels in the image
  totalPixVals := (height * width) * (65535 * 4)
  return float64(sumPixDiffVals) / float64(totalPixVals)
}

func main () {
  checkArgs()
  im1 := loadImage(os.Args[1])
  im2 := loadImage(os.Args[2])
  checkDimensions(im1, im2) 
  checkColorModel(im1, im2) 
  ratio := DiffImages(im1, im2)
  percentage := ratio * 100
  fmt.Printf("Images differ by %v%%", percentage)
}
