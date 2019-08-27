// +build !heif,!darwin

package imagecoding

import (
	"image"
)

func ConfigHeif(data []byte) (image.Config, string, error) {
	return image.Config{}, "", image.ErrFormat
}

func TransformHeif(data []byte, grayscale bool, scale ScaleFunc) (out image.Image, width, height int, scaleFactor float64, err error) {
	return nil, 0, 0, 0, image.ErrFormat
}
