package imagecoding

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"testing"

	"github.com/Nr90/imgsim"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

type Circle struct {
	X, Y, R float64
}

func (c *Circle) brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	}
	return 255
}

var (
	img = makeTestImage()
	ref = imgsim.AverageHash(img)
)

func makeTestImage() *image.Gray {
	var w, h int = int(math.Floor(2480)), int(math.Floor(3508))
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := math.Min(hw, hh) / 4
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), math.Min(hw, hh) / 2}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), math.Min(hw, hh) / 2}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), math.Min(hw, hh) / 2}

	m := image.NewGray(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.GrayModel.Convert(
				color.RGBA{
					cr.brightness(float64(x), float64(y)),
					cg.brightness(float64(x), float64(y)),
					cb.brightness(float64(x), float64(y)),
					255,
				},
			)
			m.Set(x, y, c)
		}
	}
	return m
}

func readImageHash(filename string, t *testing.T, diffhash bool) imgsim.Hash {
	imgbytes, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Error reading file %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(imgbytes))
	if err != nil {
		t.Fatalf("Error reading file %v", err)
	}
	if diffhash {
		return imgsim.DifferenceHash(img)
	}
	return imgsim.AverageHash(img)
}
