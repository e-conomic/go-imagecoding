// +build heif darwin

package imagecoding

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeifTransform(t *testing.T) {
	sample, err := ioutil.ReadFile("testdata/world-political.heic")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	img, _, _, _, err := TransformHeif(sample, true, DefaultScale)
	if assert.NoError(t, err) {
		assert.Equal(t, 1754, img.Bounds().Dx())
		assert.Equal(t, 1002, img.Bounds().Dy())
	}
}

func BenchmarkHeifTransform(b *testing.B) {
	sample, err := ioutil.ReadFile("testdata/world-political.heic")
	if !assert.NoError(b, err) {
		b.FailNow()
	}
	for n := 0; n < b.N; n++ {
		_, _, _, _, err := TransformHeif(sample, true, DefaultScale)
		if !assert.NoError(b, err) {
			b.FailNow()
		}
	}
}
