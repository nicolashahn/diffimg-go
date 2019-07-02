package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nicolashahn/diffimg-go/diffimg"
	"github.com/nicolashahn/diffimg-go/imgutil"
)

func main() {
	// Parse flags
	diffOutput := flag.String("filename", "", "Generate a diff image and save it at this filename.\n(should have a .png extension)\nIf not passed, only a ratio will be returned.")
	ignoreAlpha := flag.Bool("ignorealpha", false, "Ignore the alpha channel when doing the ratio calculation, or if \ngenerating an image, invert the alpha channel for the generated image.")

	flag.Parse()

	// Check that both images have been supplied
	apath, bpath := flag.Arg(0), flag.Arg(1)
	if apath == "" || bpath == "" {
		fmt.Fprintln(os.Stderr, "requires two images as arguments")
		flag.Usage()
		os.Exit(1)
	}

	// Load Images
	a, aerr := imgutil.Load(apath)
	b, berr := imgutil.Load(bpath)
	if aerr != nil || berr != nil {
		if aerr != nil {
			fmt.Fprintf(os.Stderr, "failed to load %q: %v\n", apath, aerr)
		}
		if berr != nil {
			fmt.Fprintf(os.Stderr, "failed to load %q: %v\n", bpath, berr)
		}
		os.Exit(1)
	}

	// Ensure images are compatible for diffing
	diffimg.CheckDimensions(a, b)

	var ratio float64
	if *diffOutput != "" {
		diffIm := diffimg.CreateDiffImage(a, b, *ignoreAlpha)
		ratio = diffimg.GetRatioFromImage(diffIm, *ignoreAlpha)

		err := imgutil.WritePNG(*diffOutput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %q: %v\n", *diffoutput, berr)
			os.Exit(1)
		}
	} else {
		// Just getting the ratio without creating a diffIm is faster
		ratio = diffimg.GetRatio(a, b, *ignoreAlpha)
	}

	fmt.Println(ratio)
}
