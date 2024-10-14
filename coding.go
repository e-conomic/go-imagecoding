// Package imagecoding provides go bindings for en-/de-coding images using
// image processing c libraries found in common systems
package imagecoding

import (
	"errors"
	"image"
	"image/color"
	"math"
)

type ScaleFunc func(pageWidth, pageHeight int) (imgWidth, imgHeight int, scaleFactor float64)

const (
	A4Short = 210 / 25.4 * 150 // 210 mm / 25.4 mm/inch * 150 ppi ≈ 1204 pixels
	A4Long  = 297 / 25.4 * 150 // 297 mm / 25.4 mm/inch * 150 ppi ≈ 1754 pixels
)

// Calculcate at what scale to use for OCR optimized pages
// We prefer maximum what would be the equivalent for a A4 page at 150 ppi
func DefaultScale(pageWidth, pageHeight int) (imgWidth, imgHeight int, scaleFactor float64) {
	w := float64(pageWidth)
	h := float64(pageHeight)

	// The maximum size we will render, capped at A4 paper equivalent
	// Expressed as short & long for orientation support
	reqShort := math.Min(w, h)
	reqLong := math.Max(w, h)
	maxShort := math.Min(reqShort, A4Short)
	maxLong := math.Min(reqLong, A4Long)

	// Calculate the scale factor
	shortRatio := maxShort / math.Min(w, h)
	longRatio := maxLong / math.Max(w, h)
	scaleFactor = math.Min(shortRatio, longRatio)

	// Round to integers
	imgWidth = int(math.Round(w * scaleFactor))
	imgHeight = int(math.Round(h * scaleFactor))

	return imgWidth, imgHeight, scaleFactor
}

type ImgFormat string

const (
	Bmp  ImgFormat = "bmp"
	Gif  ImgFormat = "gif"
	Png  ImgFormat = "png"
	Jpeg ImgFormat = "jpg"
	Tiff ImgFormat = "tif"
	Webp ImgFormat = "webp"
	Heif ImgFormat = "heif"
)

var ErrEmptyInput = errors.New("empty input data")

// RGBImage represent image data which has RGB colors.
// RGBImage is compatible with image.RGBA, but does not have alpha channel to reduce using memory.
type RGBImage struct {
	// Pix holds the image's stream, in R, G, B order.
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// NewRGBImage allocates and returns RGB image
func NewRGBImage(r image.Rectangle) *RGBImage {
	w, h := r.Dx(), r.Dy()
	return &RGBImage{Pix: make([]uint8, 3*w*h), Stride: 3 * w, Rect: r}
}

// ColorModel returns RGB color model.
func (p *RGBImage) ColorModel() color.Model {
	return RGBModel
}

// Bounds implements image.Image.At
func (p *RGBImage) Bounds() image.Rectangle {
	return p.Rect
}

// At implements image.Image.At
func (p *RGBImage) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

// RGBAAt returns the color of the pixel at (x, y) as RGBA.
func (p *RGBImage) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
	return color.RGBA{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2], 0xFF}
}

// RGBModel is RGB color model instance
var RGBModel = color.ModelFunc(rgbModel)

func rgbModel(c color.Color) color.Color {
	if _, ok := c.(RGB); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

// RGB color
type RGB struct {
	R, G, B uint8
}

// RGBA implements Color.RGBA
func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(0xFFFF)
	return
}

// Make sure RGBImage implements image.Image.
// See https://golang.org/doc/effective_go.html#blank_implements.
var _ image.Image = new(RGBImage)
