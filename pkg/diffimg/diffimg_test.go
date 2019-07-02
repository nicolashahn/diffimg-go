package diffimg

import (
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

// Test helpers
// Templated tests for each function, note the lowercase `test` in the name

func testAbsDiffUint8(t *testing.T, x, y, expected uint8) {
	val := absDiffUint8(x, y)
	if val != expected {
		t.Errorf("AbsDiffUint8(%v, %v): expected %v, got %v\n",
			x, y, 10, val)
	}
}

func testRgbaArrayUint8(t *testing.T, r, g, b, a uint8, expected [4]uint8) {
	col := color.RGBA{r, g, b, a}
	colArr := rgbaArrayUint8(col)
	if colArr != expected {
		t.Errorf("RgbaArrayUint8(1,2,3,4): expected %v, got %v\n",
			expected, colArr)
	}
}

func testGetRatio(
	t *testing.T, im1file, im2file string, ignoreAlpha bool, expected float64) {
	im1 := LoadImage(im1file)
	im2 := LoadImage(im2file)
	ratio := GetRatio(im1, im2, ignoreAlpha)
	if ratio != expected {
		t.Errorf("GetRatio(%v, %v): expected %v, got %v\n",
			im1file, im2file, expected, ratio)
	}
}

func testGetRatioFromImage(
	t *testing.T, im1file, im2file string, ignoreAlpha bool, expected float64) {
	im1 := LoadImage(im1file)
	im2 := LoadImage(im2file)
	diffIm := CreateDiffImage(im1, im2, ignoreAlpha)
	ratio := GetRatioFromImage(diffIm, ignoreAlpha)
	if ratio != expected {
		t.Errorf("GetRatioFromImage(%v, %v): expected %v, got %v\n",
			im1file, im2file, expected, ratio)
	}
}

func testBothRatioMethods(
	t *testing.T, im1file, im2file string, ignoreAlpha bool, expected float64) {
	testGetRatio(t, im1file, im2file, ignoreAlpha, expected)
	testGetRatioFromImage(t, im1file, im2file, ignoreAlpha, expected)
}

// Actual tests

func TestRgbaArrayUint8(t *testing.T) {
	testRgbaArrayUint8(t, 1, 2, 3, 4, [4]uint8{1, 2, 3, 4})
	testRgbaArrayUint8(t, 255, 255, 255, 255, [4]uint8{255, 255, 255, 255})
	testRgbaArrayUint8(t, 0, 0, 0, 0, [4]uint8{0, 0, 0, 0})
}

func TestAbsDiffUint8(t *testing.T) {
	testAbsDiffUint8(t, 30, 40, 10)
	testAbsDiffUint8(t, 255, 1, 254)
	testAbsDiffUint8(t, 1, 255, 254)
	testAbsDiffUint8(t, 0, 0, 0)
}

func TestGetRatio(t *testing.T) {
	// Pure black vs pure white image, both opaque
	testBothRatioMethods(t, "../../testdata/black.png", "../../testdata/white.png", false, 0.75)
	// Same image
	testBothRatioMethods(t, "../../testdata/im1.png", "../../testdata/im1.png", false, 0)
	// Image with non-homogenous alpha
	testBothRatioMethods(t, "../../testdata/mario-circle-node.png",
		"../../testdata/mario-circle-cs.png",
		false,
		0.002123925685759868)
}

func TestGetRatioIgnoreAlpha(t *testing.T) {
	testBothRatioMethods(t, "../../testdata/black.png", "../../testdata/white.png", true, 1.0)
	testBothRatioMethods(t, "../../testdata/im1.png", "../../testdata/im1.png", true, 0)
	testBothRatioMethods(t, "../../testdata/mario-circle-node.png",
		"../../testdata/mario-circle-cs.png",
		true,
		0.0017478156325230589)
}
