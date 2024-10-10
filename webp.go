package imagecoding

import (
	"bufio"
	"bytes"
	"image"

	"github.com/kolesa-team/go-webp/encoder"
	newwebp "github.com/kolesa-team/go-webp/webp"
)

func EncodeWebP(buf *bytes.Buffer, img image.Image) ([]byte, error) {
	options, err := encoder.NewLosslessEncoderOptions(encoder.PresetDefault, 0)
	if err != nil {
		return nil, err
	}

	pageWriter := bufio.NewWriter(buf)

	err = newwebp.Encode(pageWriter, img, options)
	if err != nil {
		return nil, err
	}
	pageWriter.Flush()
	return buf.Bytes(), nil
}
