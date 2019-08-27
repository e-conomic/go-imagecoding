package imagecoding

import (
	"image"
	"io"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"

	"go.uber.org/zap"
)

// Enum representation for Exif Orientation
type Orientation uint8

const (
	TopLeft     Orientation = 1
	TopRight    Orientation = 2
	BottomRight Orientation = 3
	BottomLeft  Orientation = 4
	LeftTop     Orientation = 5
	RightTop    Orientation = 6
	RightBottom Orientation = 7
	LeftBottom  Orientation = 8
)

// GetOrientation returns the image orientation from EXIF data
// https://www.daveperrett.com/articles/2012/07/28/exif-orientation-handling-is-a-ghetto/
//
// TopLeft     1: 0 degrees – the correct orientation, no adjustment is required.
// TopRight    2: 0 degrees, mirrored – image has been flipped back-to-front.
// BottomRight 3: 180 degrees – image is upside down.
// BottomLeft  4: 180 degrees, mirrored – image is upside down and flipped back-to-front.
// LeftTop     5: 90 degrees – image is on its side.
// RightTop    6: 90 degrees, mirrored – image is on its side and flipped back-to-front.
// RightBottom 7: 270 degrees – image is on its far side.
// LeftBottom  8: 270 degrees, mirrored – image is on its far side and flipped back-to-front.
func GetOrientation(reader io.Reader) Orientation {
	x, err := exif.Decode(reader)
	if err != nil {
		zap.L().Debug("exif decode error", zap.String("error", err.Error()))
		return TopLeft
	}
	if x != nil {
		orient, err := x.Get(exif.Orientation)
		if err != nil {
			zap.L().Debug("exif decode error", zap.String("error", err.Error()))
			return TopLeft
		}
		if orient != nil {
			intOrient, err := strconv.ParseUint(orient.String(), 10, 8)
			if err != nil {
				return TopLeft
			}
			return Orientation(intOrient)
		}
	}
	return TopLeft
}

// FixOrientation uses the imaging library to correct for orientation
func FixOrientation(img image.Image, orient Orientation) image.Image {
	switch orient {
	case TopLeft:
		return img
	case TopRight:
		return imaging.FlipH(img)
	case BottomRight:
		return imaging.Rotate180(img)
	case BottomLeft:
		return imaging.Rotate180(imaging.FlipH(img))
	case LeftTop:
		return imaging.Rotate270(imaging.FlipV(img))
	case RightTop:
		return imaging.Rotate270(img)
	case RightBottom:
		return imaging.Rotate90(imaging.FlipV(img))
	case LeftBottom:
		return imaging.Rotate90(img)
	default:
		return img
	}
}
