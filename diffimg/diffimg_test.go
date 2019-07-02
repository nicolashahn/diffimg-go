package diffimg

import (
	"fmt"
	"image/color"
	"path/filepath"
	"testing"

	"github.com/nicolashahn/diffimg-go/imgutil"
)

func Test(t *testing.T) {
	type Test struct {
		A, B        string
		IgnoreAlpha bool
		Expected    float64
	}

	var tests = []Test{
		{"black.png", "white.png", false, 0.75},
		{"im1.png", "im1.png", false, 0},
		{"mario-circle-node.png", "mario-circle-cs.png", false, 0.002123925685759868},

		{"black.png", "white.png", true, 1.0},
		{"im1.png", "im1.png", true, 0},
		{"mario-circle-node.png", "mario-circle-cs.png", true, 0.0017478156325230589},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v-%v: IgnoreAlpha:%v", test.A, test.B, test.IgnoreAlpha), func(t *testing.T) {
			a, err := imgutil.Load(filepath.Join("../testdata", test.A))
			if err != nil {
				t.Fatal(err)
			}
			b, err := imgutil.Load(filepath.Join("../testdata", test.B))
			if err != nil {
				t.Fatal(err)
			}

			ratio := GetRatio(a, b, test.IgnoreAlpha)
			if ratio != test.Expected {
				t.Errorf("GetRatio: got ratio %v, expected %v", ratio, test.Expected)
			}

			diff := CreateDiffImage(a, b, test.IgnoreAlpha)
			ratio = GetRatioFromImage(diff, test.IgnoreAlpha)
			if ratio != test.Expected {
				t.Errorf("CreateDiffImage: got ratio %v, expected %v", ratio, test.Expected)
			}
		})
	}
}

func TestRgbaArrayUint8(t *testing.T) {
	type Test struct {
		Color    color.RGBA
		Expected [4]uint8
	}

	var tests = []Test{
		{color.RGBA{1, 2, 3, 4}, [4]uint8{1, 2, 3, 4}},
		{color.RGBA{255, 255, 255, 255}, [4]uint8{255, 255, 255, 255}},
		{color.RGBA{0, 0, 0, 0}, [4]uint8{0, 0, 0, 0}},
	}

	for _, test := range tests {
		got := rgbaArrayUint8(test.Color)
		if got != test.Expected {
			t.Errorf("rgbaArrayUint8(%v): expected %v, got %v", test.Color, test.Expected, got)
		}
	}
}

func TestAbsDiffUint8(t *testing.T) {
	type Test struct {
		A, B, Result uint8
	}

	var tests = []Test{
		{30, 40, 10},
		{255, 1, 254},
		{127, 128, 1},
		{0, 0, 0},
	}

	for _, test := range tests {
		diffAB := absDiffUint8(test.A, test.B)
		if diffAB != test.Result {
			t.Errorf("absDiffUint8(%v, %v): expected %v, got %v", test.A, test.B, test.Result, diffAB)
		}

		diffBA := absDiffUint8(test.B, test.A)
		if diffBA != test.Result {
			t.Errorf("absDiffUint8(%v, %v): expected %v, got %v", test.B, test.A, test.Result, diffBA)
		}
	}
}
