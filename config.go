// Package imagecoding provides go bindings for en-/de-coding images using
// image processing c libraries found in common systems
package imagecoding

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/png"

	"github.com/h2non/filetype"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
)

// DecodeConfig is like image.DecodeConfig but supports additional formats
// For JPEGs it uses jpeg-turbos internal function for compatibility
func DecodeConfig(content []byte) (image.Config, string, error) {
	// Look at the magic bytes to determine the file type
	kind, err := filetype.Match(content)
	if err != nil {
		return image.Config{}, "", image.ErrFormat
	}

	switch ImgFormat(kind.Extension) {
	// If this is an image, resize it
	case Jpeg:
		return ConfigJpeg(content)
	case Heif:
		return ConfigHeif(content)
	default:
		c, fmt, err := image.DecodeConfig(bytes.NewReader(content))
		return c, fmt, err
	}
}
