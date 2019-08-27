package imagecoding

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/h2non/filetype"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

// Transform scales, colormaps and orients an image according to input param
func Transform(data []byte, grayscale bool, scale ScaleFunc) (out image.Image, width, height int, scaleFactor float64, err error) {
	// Look at the magic bytes to determine the file type
	kind, err := filetype.Match(data)
	if err != nil {
		return nil, 0, 0, 0, errors.New("could not determine file type")
	}
	format := ImgFormat(kind.Extension)

	var img image.Image
	imagefile := bytes.NewReader(data)

	switch format {
	case Webp:
		img, err = webp.Decode(imagefile)
	case Png:
		img, err = png.Decode(imagefile)
	case Jpeg:
		// Early return for JPEG fast path
		return TransformJpeg(data, grayscale, scale)
	case Tiff:
		orient := GetOrientation(imagefile)
		img, err = tiff.Decode(imagefile)
		img = FixOrientation(img, orient)
	case Gif:
		img, err = gif.Decode(imagefile)
	case Bmp:
		img, err = bmp.Decode(imagefile)
	case Heif:
		return TransformHeif(data, grayscale, scale)
	default:
		err = image.ErrFormat
	}
	if err != nil {
		return nil, 0, 0, 0, err
	}
	width = img.Bounds().Dx()
	height = img.Bounds().Dy()

	// Scale the image
	imgWidth, imgHeight, scaleFactor := scale(img.Bounds().Dx(), img.Bounds().Dy())
	if scaleFactor > 1.1 || scaleFactor < 0.9 {
		img = imaging.Resize(img, imgWidth, imgHeight, imaging.CatmullRom)
	} else {
		scaleFactor = 1
	}

	if grayscale {
		// Drop the channels we don't need by converting to image.Gray
		bounds := img.Bounds()
		imgGray := image.NewGray(bounds)
		for y := 0; y < bounds.Max.Y; y++ {
			for x := 0; x < bounds.Max.X; x++ {
				oldPixel := img.At(x, y)
				pixel := color.GrayModel.Convert(oldPixel)
				imgGray.Set(x, y, pixel)
			}
		}
		img = imgGray
	}

	return img, width, height, scaleFactor, nil
}
