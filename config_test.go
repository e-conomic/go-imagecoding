package imagecoding

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		width    int
		height   int
	}{
		{"rose", "testdata/rose.jpg", 227, 149},
		{"jpeg-warning-invalid-sos", "testdata/samsung-invalid-sos.jpg", 440, 500},
		{"f1", "testdata/f1-exif.jpg", 40, 80},
		{"f2", "testdata/f2-exif.jpg", 40, 80},
	}
	for _, tt := range tests {
		imgbytes, err := os.ReadFile(tt.filename)
		assert.NoError(t, err)
		exHeight := tt.height
		exWidth := tt.width
		t.Run(tt.name, func(t *testing.T) {
			config, _, err := DecodeConfig(imgbytes)
			if assert.NoError(t, err) {
				assert.Equal(t, exHeight, config.Height)
				assert.Equal(t, exWidth, config.Width)
			}
		})
	}
}
