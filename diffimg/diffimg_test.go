package diffimg

import (
	"fmt"
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

			ratio, err := RatioRGBA(a, b, test.IgnoreAlpha)
			if err != nil {
				t.Fatalf("RatioRGBA failed: %v", err)
			}
			if ratio != test.Expected {
				t.Errorf("GetRatio: got ratio %v, expected %v", ratio, test.Expected)
			}

			_, ratio, err = RGBA(a, b, test.IgnoreAlpha)
			if err != nil {
				t.Fatalf("RGBA failed: %v", err)
			}
		})
	}
}
