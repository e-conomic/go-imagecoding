package imagecoding

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	sample, err := os.ReadFile("testdata/gamer.png")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	img, _, _, _, err := Transform(sample, true, DefaultScale)
	if assert.NoError(t, err) {
		assert.Equal(t, 1240, img.Bounds().Dx())
		assert.Equal(t, 1317, img.Bounds().Dy())
	}
}

func TestTransformNoError(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			"world",
			"testdata/world-political.jpg",
		},
	}
	for _, tt := range tests {
		sample, err := os.ReadFile(tt.filename)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, err := Transform(sample, true, DefaultScale)
			assert.NoError(t, err)
		})
	}
}

func TestTransformEmpty(t *testing.T) {
	empty := []byte{}
	_, _, _, _, err := Transform(empty, true, DefaultScale)
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyInput, err)
	}
}

func BenchmarkPNGTransform(b *testing.B) {
	sample, err := os.ReadFile("testdata/gamer.png")
	if !assert.NoError(b, err) {
		b.FailNow()
	}
	for n := 0; n < b.N; n++ {
		_, _, _, _, err := Transform(sample, true, DefaultScale)
		if !assert.NoError(b, err) {
			b.FailNow()
		}
	}
}
