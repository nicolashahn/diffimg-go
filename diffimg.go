package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)


// Helper functions


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func abs(x int) uint32 {
	if x < 0 {
		return uint32(-x)
	}
	return uint32(x)
}

func getWidthAndHeight(im image.Image) (int, int) {
	width := im.Bounds().Max.X
	height := im.Bounds().Max.Y
	return width, height
}

func checkArgs() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "must pass exactly two arguments (image1 and image2)")
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
		fmt.Fprintln(os.Stderr, "image dimensions are different:",
			im1.Bounds(), im2.Bounds())
		os.Exit(1)
	}
}

// TODO Convert 2nd image color model to 1st
func checkColorModel(im1, im2 image.Image) {
	if im1.ColorModel() != im2.ColorModel() {
		panic("color models are different")
	}
}


// Functions for actually getting the ratio, generating diff image


func rgbaArrayUint16(col color.Color) [4]uint16 {
	r, g, b, a := col.RGBA()
	return [4]uint16{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// return a color created from diffing each of the image's color values
// at (x,y)
func pixelDiff(im1, im2 image.Image, x, y int) color.Color {
	rgba1 := rgbaArrayUint16(im1.At(x,y))
	rgba2 := rgbaArrayUint16(im2.At(x,y))
	var rgba3 [4]uint16
	for i, _ := range rgba1 {
		rgba3[i] = uint16(abs(int(rgba1[i]) - int(rgba2[i])))
	}
	r, g, b, a := rgba3[0], rgba3[1], rgba3[2], rgba3[3]
	// invert alpha channel - if both images have full alpha, the resulting new
	// pixel would be completely transparent
	a = 65535 - a
	newColor := color.RGBA64{r,g,b,a}
	return newColor
}

// absolute difference between the channel values of the pixels at the same
// coordinates (x,y) in im1, im2
// Example (real channel values max at 65535, not 255 though):
// RGBA for im1 at (x,y): (100, 100, 180, 255)
// RGBA for im2 at (x,y): (120, 100, 100, 255)
// abs(100-120) + abs(100-100) + abs(180-100) + abs(255-255)
// 20 + 0 + 80 + 0 = 100
// return 100
func sumPixelDiff(im1, im2 image.Image, x, y int) uint64 {
	rgba1 := rgbaArrayUint16(im1.At(x,y))
	rgba2 := rgbaArrayUint16(im2.At(x,y))
	var pixDiffVal uint64
	for i, _ := range rgba1 {
		chanDiff := uint64(abs(int(rgba1[i]) - int(rgba2[i])))
		pixDiffVal += chanDiff
	}
	return pixDiffVal
}

func sumDiffPixelValues(diffIm image.Image, x, y int) uint64 {
	rgba := rgbaArrayUint16(diffIm.At(x,y))
	// diffIm's alpha channel was inverted
	rgba[3] = 65535 - rgba[3]
	var sum uint64
	for _,v := range rgba {
		sum += uint64(v)
	}
	return sum
}

// Calculate difference ratio between two Images
// Adds up all the differences in each pixel's channel values, and averages
// over all pixels
func GetRatio(im1, im2 image.Image) float64 {
	var sumPixDiffVals uint64
	width, height := getWidthAndHeight(im1)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sumPixDiffVals += sumPixelDiff(im1, im2, x, y)
		}
	}
	// Sum of max channel values for all pixels in the image
	totalPixVals := (height * width) * (65535 * 4)
	return float64(sumPixDiffVals) / float64(totalPixVals)
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
	totalPixVals := (height * width) * (65535 * 4)
	return float64(sum) / float64(totalPixVals)
}

func main() {

	createDiffIm := true

	checkArgs()
	im1 := loadImage(os.Args[1])
	im2 := loadImage(os.Args[2])
	
	// Ensure images are compatible
	checkDimensions(im1, im2)
	checkColorModel(im1, im2)

	var ratio float64
	if createDiffIm {
		diffIm := CreateDiffImage(im1, im2)
		ratio = GetRatioFromImage(diffIm)
		newFile, _ := os.Create("diffimg.png")
		png.Encode(newFile, diffIm)
	} else {
		// Just getting the ratio without creating a diffIm is faster
		ratio = GetRatio(im1, im2)
	}

	percentage := ratio * 100
	fmt.Printf("Images differ by %v%%", percentage)
}
