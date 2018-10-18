// This package provides image diffing tools. You can create a difference image
// from two given images, or find their difference ratio/percentage.

package diffimg

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

//////////////////////
// Helper functions //
//////////////////////

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func LoadImage(filepath string) image.Image {
	file, err := os.Open(filepath)
	CheckErr(err)
	defer file.Close()
	im, _, err := image.Decode(file)
	CheckErr(err)
	return im
}

func GetWidthAndHeight(im image.Image) (int, int) {
	width := im.Bounds().Max.X
	height := im.Bounds().Max.Y
	return width, height
}

func CheckDimensions(im1, im2 image.Image) {
	if im1.Bounds() != im2.Bounds() {
		im1w, im1h := GetWidthAndHeight(im1)
		im2w, im2h := GetWidthAndHeight(im2)
		fmt.Fprintf(os.Stderr, "Image dimensions are different: %vx%v, %vx%v\n",
			im1w, im1h, im2w, im2h)
		os.Exit(1)
	}
}

func RgbaArrayUint8(col color.Color) [4]uint8 {
	// The values are stored as uint8, but returned by RGBA() as uint32 for
	// compatibility, we need to manually convert back to uint8
	// See https://golang.org/src/image/color/color.go
	r, g, b, a := col.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	return [4]uint8{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Abs(x - y)
func AbsDiffUint8(x, y uint8) uint8 {
	// int16 cast because result of subtraction could be negative
	diff := int16(x) - int16(y)
	if diff < 0 {
		return uint8(-diff)
	}
	return uint8(diff)
}

////////////////////////////////////////////////////////
// Get difference ratio without creating a diff image //
// (This is faster than creating a diff image first)  //
////////////////////////////////////////////////////////

// Absolute difference between the channel values of the pixels at the same
// coordinates (x,y) in im1, im2
// Example:
// RGBA for im1 at (x,y): (100, 100, 180, 255)
// RGBA for im2 at (x,y): (120, 100, 100, 255)
// abs(100-120) + abs(100-100) + abs(180-100) + abs(255-255)
// 20 + 0 + 80 + 0 = 100
// return 100
func sumPixelDiff(im1, im2 image.Image, x, y int) uint16 {
	rgba1 := RgbaArrayUint8(im1.At(x, y))
	rgba2 := RgbaArrayUint8(im2.At(x, y))
	var pixDiffVal uint16
	for i, _ := range rgba1 {
		chanDiff := uint16(AbsDiffUint8(rgba1[i], rgba2[i]))
		pixDiffVal += chanDiff
	}
	return pixDiffVal
}

// Calculate difference ratio between two Images
// Adds up all the differences in each pixel's channel values, and averages
// over all pixels
func GetRatio(im1, im2 image.Image) float64 {
	var sum uint64
	width, height := GetWidthAndHeight(im1)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sum += uint64(sumPixelDiff(im1, im2, x, y))
		}
	}
	// Sum of max channel values for all pixels in the image
	totalPixVals := (height * width) * (255 * 4)
	return float64(sum) / float64(totalPixVals)
}

/////////////////////////
// Generate diff image //
/////////////////////////

// NOTE for the following functions relating to generated diff images:
// The invertAlpha param flips the alpha value - without it, if two pixels have
// different RGB values but the same alpha, the resulting pixel will be
// invisible. For most images with an unused alpha channel (ie the image is
// fully opaque) this means the diff image will be fully transparent. In this
// case, invertAlpha should be set to true. When diffing graphics like logos or
// other images that make use of a transparent background, you may want to set
// invertAlpha to false to see where the difference in overlap is. The param
// does not affect the diff ratio, unless you use differing bool values in
// CreateDiffImage() and GetRatioFromImage().

// Return a color created from diffing each of the image's color values
// at (x,y).
func pixelDiff(im1, im2 image.Image, x, y int, invertAlpha bool) color.Color {
	rgba1 := RgbaArrayUint8(im1.At(x, y))
	rgba2 := RgbaArrayUint8(im2.At(x, y))
	var rgba3 [4]uint8
	for i, _ := range rgba1 {
		rgba3[i] = AbsDiffUint8(rgba1[i], rgba2[i])
	}
	r, g, b, a := rgba3[0], rgba3[1], rgba3[2], rgba3[3]
	if invertAlpha {
		a = 255 - a
	}
	newColor := color.RGBA{r, g, b, a}
	return newColor
}

// Create a new image made by diffing each color value (RGBA) at each pixel in
// im1 and im2
func CreateDiffImage(im1, im2 image.Image, invertAlpha bool) image.Image {
	width, height := GetWidthAndHeight(im1)
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	diffIm := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newPixel := pixelDiff(im1, im2, x, y, invertAlpha)
			diffIm.Set(x, y, newPixel)
		}
	}
	return diffIm
}

////////////////////////////////////////////
// Get difference ratio from a diff image //
////////////////////////////////////////////

// Sum the channel values at the given coordinates of the diff image
func sumDiffPixelValues(diffIm image.Image, x, y int, invertAlpha bool) uint64 {
	rgba := RgbaArrayUint8(diffIm.At(x, y))
	if invertAlpha {
		rgba[3] = 255 - rgba[3]
	}
	var sum uint64
	for _, v := range rgba {
		sum += uint64(v)
	}
	return sum
}

// Get the ratio by summing the diff image's pixel channel values
func GetRatioFromImage(diffIm image.Image, invertAlpha bool) float64 {
	width, height := GetWidthAndHeight(diffIm)
	var sum uint64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sum += sumDiffPixelValues(diffIm, x, y, invertAlpha)
		}
	}
	totalPixVals := (height * width) * (255 * 4)
	return float64(sum) / float64(totalPixVals)
}
