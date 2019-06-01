package main

import (
	"flag"
	"fmt"
	"github.com/nicolashahn/diffimg-go/pkg/diffimg"
	"image/png"
	"os"
)

func parseArgs() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "Requires two positional args: filename1, filename2")
		os.Exit(1)
	}
}

func main() {

	// Command line flags
	createDiffImPtr := flag.String("filename", "",
		`Generate a diff image and save it at this filename.
(should have a .png extension)
If not passed, only a ratio will be returned.`)
	ignoreAlphaPtr := flag.Bool("ignorealpha", false,
		`Ignore the alpha channel when doing the ratio calculation, or if 
generating an image, invert the alpha channel for the generated image.`)

	parseArgs()

	im1 := diffimg.LoadImage(flag.Args()[0])
	im2 := diffimg.LoadImage(flag.Args()[1])

	// Ensure images are compatible for diffing
	diffimg.CheckDimensions(im1, im2)

	var ratio float64
	if *createDiffImPtr != "" {
		diffIm := diffimg.CreateDiffImage(im1, im2, *ignoreAlphaPtr)
		ratio = diffimg.GetRatioFromImage(diffIm, *ignoreAlphaPtr)
		newFile, _ := os.Create(*createDiffImPtr)
		png.Encode(newFile, diffIm)
	} else {
		// Just getting the ratio without creating a diffIm is faster
		ratio = diffimg.GetRatio(im1, im2, *ignoreAlphaPtr)
	}

	fmt.Println(ratio)
}
