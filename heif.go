// +build heif darwin

package imagecoding

import (
	"image"
	"image/color"
	"runtime"
)

import "github.com/strukturag/libheif/go/heif"

func ConfigHeif(data []byte) (image.Config, string, error) {
	if len(data) == 0 {
		return image.Config{}, string(Heif), EmptyInputError
	}
	ctx, err := heif.NewContext()
	if err != nil {
		return image.Config{}, string(Heif), err
	}
	err = ctx.ReadFromMemory(data)
	if err != nil {
		return image.Config{}, string(Heif), err
	}
	img, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return image.Config{}, string(Heif), err
	}
	cfg := image.Config{
		ColorModel: color.YCbCrModel,
		Width:      img.GetWidth(),
		Height:     img.GetHeight(),
	}
	runtime.KeepAlive(ctx)
	return cfg, string(Heif), nil
}

func DecodeHeif(data []byte) (image.Image, error) {
	if len(data) == 0 {
		return nil, EmptyInputError
	}
	ctx, err := heif.NewContext()
	if err != nil {
		return nil, err
	}
	err = ctx.ReadFromMemory(data)
	if err != nil {
		return nil, err
	}
	imgh, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return nil, err
	}
	img, err := imgh.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	runtime.KeepAlive(ctx)
	if err != nil {
		return nil, err
	}
	goimg, err := img.GetImage()
	if err != nil {
		return nil, err
	}
	return goimg, nil
}

func TransformHeif(data []byte, grayscale bool, scale ScaleFunc) (out image.Image, width, height int, scaleFactor float64, err error) {
	if len(data) == 0 {
		return nil, 0, 0, 0, EmptyInputError
	}
	ctx, err := heif.NewContext()
	if err != nil {
		return nil, 0, 0, 0, err
	}
	err = ctx.ReadFromMemory(data)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	imgh, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	width = imgh.GetWidth()
	height = imgh.GetHeight()

	// Calculate scaling factor
	scaledW, scaledH, scaleFactor := scale(width, height)

	var img *heif.Image
	img, err = imgh.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	runtime.KeepAlive(ctx)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Scale if required
	if scaleFactor > 1.1 || scaleFactor < 0.9 {
		img, err = img.ScaleImage(scaledW, scaledH)
		if err != nil {
			return nil, 0, 0, 0, err
		}
	} else {
		scaleFactor = 1
	}

	goimg, err := img.GetImage()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// libheif does not support conversion from YUV/RGB -> Gray Scale
	if grayscale {
		// Drop the channels we don't need by converting to image.Gray
		bounds := goimg.Bounds()
		imgGray := image.NewGray(bounds)
		for y := 0; y < bounds.Max.Y; y++ {
			for x := 0; x < bounds.Max.X; x++ {
				imgGray.Set(x, y, color.GrayModel.Convert(goimg.At(x, y)))
			}
		}
		goimg = imgGray
	}
	return goimg, width, height, scaleFactor, nil
}
