// Package diffimg provides image diffing tools. You can create a difference
// image from two given images, or quantify their difference as a ratio/%.
package diffimg

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

const maxChannelVal = 255

//////////////////////
// Helper functions //
//////////////////////

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func getWidthAndHeight(im image.Image) (int, int) {
	width := im.Bounds().Max.X
	height := im.Bounds().Max.Y
	return width, height
}

func rgbaArrayUint8(col color.Color) [4]uint8 {
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

// abs(x - y)
func absDiffUint8(x, y uint8) uint8 {
	if x > y {
		return x - y
	}
	return y - x
}

// LoadImage opens a file and tries to decode it as an image.
func LoadImage(filepath string) image.Image {
	file, err := os.Open(filepath)
	checkErr(err)
	defer file.Close()
	im, _, err := image.Decode(file)
	checkErr(err)
	return im
}

// CheckDimensions ensures that the images have the same width and height.
func CheckDimensions(im1, im2 image.Image) {
	if im1.Bounds() != im2.Bounds() {
		im1w, im1h := getWidthAndHeight(im1)
		im2w, im2h := getWidthAndHeight(im2)
		fmt.Fprintf(os.Stderr, "Image dimensions are different: %vx%v, %vx%v\n",
			im1w, im1h, im2w, im2h)
		os.Exit(1)
	}
}

////////////////////////////////////////////////////////
// Get difference ratio without creating a diff image //
// (This is faster than creating a diff image first)  //
////////////////////////////////////////////////////////

// NOTE for the following functions relating to generated diff images:
// The ignoreAlpha param ignores the alpha channel when doing the diff ratio
// calculation and generating a diff image (it sets alpha to max for all
// pixels) - without it, if two pixels have different RGB values but the same
// alpha, the resulting pixel will be invisible. For most images with an
// "unused" alpha channel (ie the image is fully opaque) this means the diff
// image will be fully transparent. In this case, ignoreAlpha should be set to
// true. When diffing graphics like logos or other images that make use of a
// transparent background, you may want to set ignoreAlpha to false to see
// where the difference in overlap is. You will get a different diff ratio if
// you set ignoreAlpha to true because the calculation is done with 3 channels
// instead of 4.

// sumPixelDiff gets the absolute difference between the channel values of the
// pixels at the same coordinates (x,y) in im1, im2
// Example:
// RGBA for im1 at (x,y): (100, 100, 180, maxChannelVal)
// RGBA for im2 at (x,y): (120, 100, 100, maxChannelVal)
// abs(100-120) + abs(100-100) + abs(180-100) + abs(maxChannelVal-maxChannelVal)
// 20 + 0 + 80 + 0 = 100
// return 100
func sumPixelDiff(im1, im2 image.Image, x, y int, ignoreAlpha bool) uint16 {
	rgba1 := rgbaArrayUint8(im1.At(x, y))
	rgba2 := rgbaArrayUint8(im2.At(x, y))
	if ignoreAlpha {
		rgba1[3] = 0
		rgba2[3] = 0
	}
	var pixDiffVal uint16
	for i := range rgba1 {
		chanDiff := uint16(absDiffUint8(rgba1[i], rgba2[i]))
		pixDiffVal += chanDiff
	}
	return pixDiffVal
}

func toRGBA(im image.Image) *image.RGBA {
	if rgba, ok := im.(*image.RGBA); ok {
		return rgba
	}
	rect := im.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, im, rect.Min, draw.Src)
	return rgba
}

// GetRatio calculates difference ratio between two Images
// Adds up all the differences in each pixel's channel values, and averages
// over all pixels.
func GetRatio(im1, im2 image.Image, ignoreAlpha bool) float64 {
	rgba1, rgba2 := toRGBA(im1), toRGBA(im2)

	var sum uint64
	for i := range rgba1.Pix {
		if ignoreAlpha && i%3 == 0 {
			continue
		}
		sum += uint64(absDiffUint8(rgba1.Pix[i], rgba2.Pix[i]))
	}

	max := float64(len(rgba1.Pix)) * maxChannelVal
	if ignoreAlpha {
		max *= 0.75 // only RGB, not A
	}

	return float64(sum) / max
}

/////////////////////////
// Generate diff image //
/////////////////////////

// pixelDiff returns a color created from diffing each of the image's color
// values at (x,y).
func pixelDiff(im1, im2 image.Image, x, y int, ignoreAlpha bool) color.Color {
	rgba1 := rgbaArrayUint8(im1.At(x, y))
	rgba2 := rgbaArrayUint8(im2.At(x, y))
	var rgba3 [4]uint8
	for i := range rgba1 {
		rgba3[i] = absDiffUint8(rgba1[i], rgba2[i])
	}
	r, g, b, a := rgba3[0], rgba3[1], rgba3[2], rgba3[3]
	if ignoreAlpha {
		a = maxChannelVal
	}
	newColor := color.RGBA{r, g, b, a}
	return newColor
}

// CreateDiffImage creates a new image made by diffing each color value (RGBA)
// at each pixel in im1 and im2.
func CreateDiffImage(im1, im2 image.Image, ignoreAlpha bool) image.Image {
	width, height := getWidthAndHeight(im1)
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	diffIm := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newPixel := pixelDiff(im1, im2, x, y, ignoreAlpha)
			diffIm.Set(x, y, newPixel)
		}
	}
	return diffIm
}

////////////////////////////////////////////
// Get difference ratio from a diff image //
////////////////////////////////////////////

// sumDiffPixelValues sums the channel values at the given coordinates of the
// diff image.
func sumDiffPixelValues(diffIm image.Image, x, y int, ignoreAlpha bool) uint64 {
	rgba := rgbaArrayUint8(diffIm.At(x, y))
	if ignoreAlpha {
		rgba[3] = 0
	}
	var sum uint64
	for _, v := range rgba {
		sum += uint64(v)
	}
	return sum
}

// GetRatioFromImage gets the ratio by summing the diff image's pixel channel
// values.
func GetRatioFromImage(diffIm image.Image, ignoreAlpha bool) float64 {
	width, height := getWidthAndHeight(diffIm)
	var sum uint64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sum += sumDiffPixelValues(diffIm, x, y, ignoreAlpha)
		}
	}
	var numChannels = 4
	if ignoreAlpha {
		numChannels = 3
	}
	totalPixVals := (height * width) * (maxChannelVal * numChannels)
	return float64(sum) / float64(totalPixVals)
}
