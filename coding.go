// Package imagecoding provides go bindings for en-/de-coding images using
// image processing c libraries found in common systems
package imagecoding

import (
	"errors"
	"math"

	"github.com/pixiv/go-libwebp/webp"
)

// RGBImage is a good idea, so let's borrow it and make it our own
type RGBImage = webp.RGBImage
type RGB = webp.RGB

var RGBModel = webp.RGBModel

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

var EmptyInputError = errors.New("empty input data")
