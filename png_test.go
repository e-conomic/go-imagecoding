package imagecoding

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/Nr90/imgsim"
	"github.com/stretchr/testify/assert"
)

func TestEncodePng(t *testing.T) {
	var buf bytes.Buffer
	imgbytes, err := EncodePng(&buf, img)
	assert.NoError(t, err)

	output, err := png.Decode(bytes.NewReader(imgbytes))
	assert.NoError(t, err)
	outputHash := imgsim.AverageHash(output)
	assert.Equal(t, ref, outputHash)
}

func BenchmarkPNG(b *testing.B) {
	var err error
	var buf bytes.Buffer
	buf.Grow(10 << 20)
	for n := 0; n < b.N; n++ {
		_, err = EncodePng(&buf, img)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}

func BenchmarkPNGComplex(b *testing.B) {
	sample, err := os.ReadFile("testdata/world-political.jpg")
	if !assert.NoError(b, err) {
		b.FailNow()
	}
	img, _, _, _, err := TransformJpeg(sample, true, DefaultScale)
	if !assert.NoError(b, err) {
		b.FailNow()
	}
	var buf bytes.Buffer
	buf.Grow(10 << 20)
	for n := 0; n < b.N; n++ {
		_, err = EncodePng(&buf, img)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
