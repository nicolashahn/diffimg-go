package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/nicolashahn/diffimg-go/diffimg"
)

func main() {
	diffOutput := flag.String("filename", "", "Generate a diff image and save it at this filename.\n(should have a .png extension)\nIf not passed, only a ratio will be returned.")
	ignoreAlpha := flag.Bool("ignorealpha", false, "Ignore the alpha channel when doing the ratio calculation, or if \ngenerating an image, invert the alpha channel for the generated image.")

	flag.Parse()

	apath, bpath := flag.Arg(0), flag.Arg(1)
	if apath == "" || bpath == "" {
		fmt.Fprintln(os.Stderr, "requires two images as arguments")
		flag.Usage()
		os.Exit(1)
	}

	a, aerr := LoadImage(apath)
	b, berr := LoadImage(bpath)
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
		newFile, _ := os.Create(*diffOutput)
		png.Encode(newFile, diffIm)
	} else {
		// Just getting the ratio without creating a diffIm is faster
		ratio = diffimg.GetRatio(a, b, *ignoreAlpha)
	}

	fmt.Println(ratio)
}

// LoadImage opens a file and tries to decode it as an image.
func LoadImage(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	m, _, err := image.Decode(file)
	return m, err
}
