package imagecoding

import (
	"bufio"
	"bytes"
	"errors"
	"image"

	"github.com/harukasan/go-libwebp/webp"
)

func EncodeWebP(buf *bytes.Buffer, img image.Image) ([]byte, error) {
	// Efficiency level between 0 (fastest, lowest compression) and 9 (slower, best compression)
	// At some point maybe we want to be able to configure this, for now, gotta go fast!
	config, configErr := webp.ConfigLosslessPreset(0)

	// Graph is the closest match to what we're doing
	config.SetImageHint(webp.HintGraph)

	// Set a bunch of lossy compression settings, or it won't run
	// These are based on the "text" presets
	config.SetSegments(2)
	config.SetSNSStrength(0)
	config.SetPass(1)
	config.SetFilterType(webp.StrongFilter)
	config.SetFilterStrength(0)
	config.SetPreprocessing(webp.PreprocessingNone)
	config.SetLossless(true)

	if configErr != nil {
		return nil, configErr
	}

	pageWriter := bufio.NewWriter(buf)
	var err error

	switch v := img.(type) {
	case *image.Gray:
		err = webp.EncodeGray(pageWriter, v, config)
	case *webp.RGBImage:
		err = webp.EncodeRGBA(pageWriter, img, config)
	case *image.RGBA:
		err = webp.EncodeRGBA(pageWriter, img, config)
	default:
		return nil, errors.New("unsupported image type")
	}
	if err != nil {
		return nil, err
	}
	pageWriter.Flush()
	return buf.Bytes(), nil
}
