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

	if *diffOutput != "" {
		diff, ratio, err := diffimg.RGBA(a, b, *ignoreAlpha)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to diff: %v\n", err)
			os.Exit(1)
		}

		err = imgutil.WritePNG(*diffOutput, diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %q: %v\n", *diffOutput, berr)
			os.Exit(1)
		}

		fmt.Println(ratio)
	} else {
		ratio, err := diffimg.RatioRGBA(a, b, *ignoreAlpha)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to diff: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(ratio)
	}
}
