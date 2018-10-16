// This package provides image diffing tools. You can create a difference image
// from two given images, or find their difference ratio/percentage.

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

// This command line flag is used in multiple functions, so it's simpler to
// keep as a global variable
var invertAlphaPtr *bool

//////////////////////
// Helper functions //
//////////////////////

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getWidthAndHeight(im image.Image) (int, int) {
	width := im.Bounds().Max.X
	height := im.Bounds().Max.Y
	return width, height
}

func parseArgs() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "Require exactly two args: filename1, filename2")
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, "Image dimensions are different:",
			im1.Bounds(), im2.Bounds())
		os.Exit(1)
	}
}

// TODO Convert 2nd image color model to 1st
func checkColorModel(im1, im2 image.Image) {
	if im1.ColorModel() != im2.ColorModel() {
		fmt.Fprintln(os.Stderr, "Color models are different:",
			im1.ColorModel(), im2.ColorModel())
		os.Exit(1)
	}
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

// Abs(x - y)
func absDiffUint8(x, y uint8) uint8 {
	int16x := int16(x)
	int16y := int16(y)
	diff := int16x - int16y
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
	rgba1 := rgbaArrayUint8(im1.At(x, y))
	rgba2 := rgbaArrayUint8(im2.At(x, y))
	var pixDiffVal uint16
	for i, _ := range rgba1 {
		chanDiff := uint16(absDiffUint8(rgba1[i], rgba2[i]))
		pixDiffVal += chanDiff
	}
	return pixDiffVal
}

// Calculate difference ratio between two Images
// Adds up all the differences in each pixel's channel values, and averages
// over all pixels
func GetRatio(im1, im2 image.Image) float64 {
	var sum uint64
	width, height := getWidthAndHeight(im1)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sum += uint64(sumPixelDiff(im1, im2, x, y))
		}
	}
	// Sum of max channel values for all pixels in the image
	totalPixVals := (height * width) * (255 * 4)
	return float64(sum) / float64(totalPixVals)
}

////////////////////////////////////////////
// Generate diff image, get ratio from it //
////////////////////////////////////////////

// Return a color created from diffing each of the image's color values
// at (x,y)
func pixelDiff(im1, im2 image.Image, x, y int) color.Color {
	rgba1 := rgbaArrayUint8(im1.At(x, y))
	rgba2 := rgbaArrayUint8(im2.At(x, y))
	var rgba3 [4]uint8
	for i, _ := range rgba1 {
		rgba3[i] = absDiffUint8(rgba1[i], rgba2[i])
	}
	r, g, b, a := rgba3[0], rgba3[1], rgba3[2], rgba3[3]
	if *invertAlphaPtr {
		a = 255 - a
	}
	newColor := color.RGBA{r, g, b, a}
	return newColor
}

// Sum the channel values at the given coordinates of the diff image
func sumDiffPixelValues(diffIm image.Image, x, y int) uint64 {
	rgba := rgbaArrayUint8(diffIm.At(x, y))
	if *invertAlphaPtr {
		rgba[3] = 255 - rgba[3]
	}
	var sum uint64
	for _, v := range rgba {
		sum += uint64(v)
	}
	return sum
}

// Create a new image made by diffing each color value (RGBA) at each pixel in
// im1 and im2
func CreateDiffImage(im1, im2 image.Image) image.Image {
	width, height := getWidthAndHeight(im1)
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	diffIm := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newPixel := pixelDiff(im1, im2, x, y)
			diffIm.Set(x, y, newPixel)
		}
	}
	return diffIm
}

// Get the ratio by summing the diff image's pixel channel values
func GetRatioFromImage(diffIm image.Image) float64 {
	width, height := getWidthAndHeight(diffIm)
	var sum uint64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sum += sumDiffPixelValues(diffIm, x, y)
		}
	}
	totalPixVals := (height * width) * (255 * 4)
	return float64(sum) / float64(totalPixVals)
}

/////////////////////////
// Main (CLI behavior) //
/////////////////////////

func main() {

	// Command line flags
	createDiffImPtr := flag.Bool("generate", false, "Generate a diff image file")
	returnRatioPtr := flag.Bool("ratio", false,
		"Output a ratio (0-1.0) instead of the percentage sentence")
	invertAlphaPtr = flag.Bool("invertalpha", false,
		"Invert the alpha channel for the generated diff image")

	parseArgs()

	im1 := loadImage(flag.Args()[0])
	im2 := loadImage(flag.Args()[1])

	// Ensure images are compatible for diffing
	checkDimensions(im1, im2)
	checkColorModel(im1, im2)

	var ratio float64
	if *createDiffImPtr {
		diffIm := CreateDiffImage(im1, im2)
		ratio = GetRatioFromImage(diffIm)
		newFile, _ := os.Create("diff.png")
		png.Encode(newFile, diffIm)
	} else {
		// Just getting the ratio without creating a diffIm is faster
		ratio = GetRatio(im1, im2)
	}

	if *returnRatioPtr {
		fmt.Println(ratio)
	} else {
		percentage := ratio * 100
		fmt.Printf("Images differ by %v%%\n", percentage)
	}
}
