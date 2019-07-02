package imgutil

import (
	"image"
	_ "image/jpeg" // to support loading jpeg-s
	"image/png"
	"os"
)

// Load opens a file and tries to decode it as an image.
func Load(filepath string) (image.Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	m, _, err := image.Decode(file)
	return m, err
}

// WritePNG writes image as a PNG to filepath.
func WritePNG(filepath string, m image.Image) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, m)
}
