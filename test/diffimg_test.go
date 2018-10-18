package diffimg_test

import (
	"github.com/nicolashahn/diffimg-go/pkg/diffimg"
	"image/color"
	"testing"
)

// Test helpers
// Templated tests for each function, note the lowercase `test` in the name

func testAbsDiffUint8(t *testing.T, x, y, expected uint8) {
	val := diffimg.AbsDiffUint8(x, y)
	if val != expected {
		t.Errorf("AbsDiffUint8(%v, %v): expected %v, got %v\n",
			x, y, 10, val)
	}
}

func testRgbaArrayUint8(t *testing.T, r, g, b, a uint8, expected [4]uint8) {
	col := color.RGBA{r, g, b, a}
	colArr := diffimg.RgbaArrayUint8(col)
	if colArr != expected {
		t.Errorf("RgbaArrayUint8(1,2,3,4): expected %v, got %v\n",
			expected, colArr)
	}
}

func testGetRatio(
	t *testing.T, im1file, im2file string, expected float64) {
	im1 := diffimg.LoadImage(im1file)
	im2 := diffimg.LoadImage(im2file)
	ratio := diffimg.GetRatio(im1, im2)
	if ratio != expected {
		t.Errorf("GetRatio(%v, %v): expected %v, got %v\n",
			im1file, im2file, expected, ratio)
	}
}

func testGetRatioFromImage(
	t *testing.T, im1file, im2file string, expected float64) {
	im1 := diffimg.LoadImage(im1file)
	im2 := diffimg.LoadImage(im2file)
	diffIm := diffimg.CreateDiffImage(im1, im2, false)
	ratio := diffimg.GetRatioFromImage(diffIm, false)
	if ratio != expected {
		t.Errorf("GetRatioFromImage(%v, %v): expected %v, got %v\n",
			im1file, im2file, expected, ratio)
	}
}

func testBothRatioMethods(
	t *testing.T, im1file, im2file string, expected float64) {
	testGetRatio(t, im1file, im2file, expected)
	testGetRatioFromImage(t, im1file, im2file, expected)
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
	testBothRatioMethods(t, "data/black.png", "data/white.png", 0.75)
	// Same image
	testBothRatioMethods(t, "data/im1.png", "data/im1.png", 0)
	// Image with non-homogenous alpha
	testBothRatioMethods(t, "data/mario-circle-node.png",
		"data/mario-circle-cs.png",
		0.002123925685759868)
}
