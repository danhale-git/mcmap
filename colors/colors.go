package colors

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

func MapColors(mapPath string) (map[string]color.Color, error) {
	dir, err := os.Stat(mapPath)
	if err != nil {
		return nil, err
	}

	if !dir.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", mapPath)
	}

	entries, err := os.ReadDir(mapPath)
	if err != nil {
		return nil, err
	}

	colors := make(map[string]color.Color)

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".png") {
			continue
		}

		f, err := os.Open(filepath.Join(mapPath, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("opening texture file: %w", err)
		}

		c, err := GetColor(f)
		if err != nil {
			return nil, fmt.Errorf("getting average color for %s: %w", e.Name(), err)
		}

		// TODO: Probably want to get the top texture for everything

		key := strings.Split(e.Name(), "_")[0]
		if strings.Contains(key, ".") {
			key = strings.Split(key, ".")[0]
		}
		colors[key] = c
	}

	return colors, nil
}

func GetColor(f *os.File) (color.Color, error) {
	m, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	bounds := m.Bounds()

	var r, g, b, a uint32

	// Sum all color values
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pr, pg, pb, pa := m.At(x, y).RGBA()
			r += pr
			g += pg
			b += pb
			a += pa
		}
	}

	// Calculate average color
	c := color.RGBA{
		R: uint8(r / uint32(bounds.Max.Y)),
		G: uint8(g / uint32(bounds.Max.Y)),
		B: uint8(b / uint32(bounds.Max.Y)),
		A: 255, //uint8(a / uint32(bounds.Max.Y)), // not bothering with alpha for now
	}

	return c, nil
}
