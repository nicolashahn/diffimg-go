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
		fmt.Fprintln(os.Stderr, "Require exactly two args: filename1, filename2")
		os.Exit(1)
	}
}

func main() {

	// Command line flags
	createDiffImPtr := flag.Bool("generate", false, "Generate a diff image file")
	returnRatioPtr := flag.Bool("ratio", false,
		"Output a ratio (0-1.0) instead of the percentage sentence")
	ignoreAlphaPtr := flag.Bool("ignorealpha", false,
		"Invert the alpha channel for the generated diff image")

	parseArgs()

	im1 := diffimg.LoadImage(flag.Args()[0])
	im2 := diffimg.LoadImage(flag.Args()[1])

	// Ensure images are compatible for diffing
	diffimg.CheckDimensions(im1, im2)

	var ratio float64
	if *createDiffImPtr {
		diffIm := diffimg.CreateDiffImage(im1, im2, *ignoreAlphaPtr)
		ratio = diffimg.GetRatioFromImage(diffIm, *ignoreAlphaPtr)
		newFile, _ := os.Create("diff.png")
		png.Encode(newFile, diffIm)
	} else {
		// Just getting the ratio without creating a diffIm is faster
		ratio = diffimg.GetRatio(im1, im2, *ignoreAlphaPtr)
	}

	if *returnRatioPtr {
		fmt.Println(ratio)
	} else {
		percentage := ratio * 100
		fmt.Printf("Images differ by %v%%\n", percentage)
	}
}
