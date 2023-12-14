package imagecoding

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/Nr90/imgsim"
	"github.com/stretchr/testify/assert"
)

// Test all possible exif orientation types
func TestJpegExif(t *testing.T) {
	// For control, read a normal image
	refhash := readImageHash("testdata/f1-exif.jpg", t, false)
	// Check all 8 orientations
	for f := 1; f <= 8; f++ {
		filename := fmt.Sprintf("testdata/f%d-exif.jpg", f)
		t.Run(fmt.Sprintf("f%d-exif.jpg", f), func(t *testing.T) {
			jpegbytes, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Error reading file %v", err)
			}
			img, _, _, scaleFactor, err := TransformJpeg(jpegbytes, true, DefaultScale)
			if err != nil {
				t.Errorf("Error during jpeg process %v", err)
			}
			if scaleFactor < 0 {
				t.Errorf("Unexpected scale factor %f", scaleFactor)
			}

			hash := imgsim.AverageHash(img)
			if hash != refhash {
				t.Errorf("Output does not match reference %v", hash)
			}
		})
	}
}

func TestJpegConformance(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{
			"rose",
			"testdata/rose.jpg",
		},
		{
			"jpeg-warning-invalid-sos",
			"testdata/samsung-invalid-sos.jpg",
		},
		{
			"world",
			"testdata/world-political.jpg",
		},
	}
	for _, tt := range tests {
		jpegbytes, err := os.ReadFile(tt.filename)
		if err != nil {
			t.Fatalf("Error reading file %v", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			// Run against libturbo
			_, _, _, turboScale, err := TransformJpeg(jpegbytes, true, DefaultScale)
			assert.NoError(t, err)

			// For native
			img, err := jpeg.Decode(bytes.NewReader(jpegbytes))
			assert.NoError(t, err)
			orient := GetOrientation(bytes.NewReader(jpegbytes))
			fixedimg := FixOrientation(img, orient)
			var refpng bytes.Buffer
			err = png.Encode(&refpng, fixedimg)
			assert.NoError(t, err)

			// Run against native implementation
			_, _, _, imagingScale, err := Transform(refpng.Bytes(), true, DefaultScale)
			assert.NoError(t, err)

			assert.InEpsilon(t, imagingScale, turboScale, 0.2, "Image scaling differ by more than 20%%")
		})
	}
}

func TestEncodeJpeg(t *testing.T) {
	var buf bytes.Buffer
	imgbytes, err := EncodeJpeg(&buf, img, 100)
	if err != nil {
		t.Errorf("Error during libjpeg-turbo encode %v", err)
	}

	output, err := jpeg.Decode(bytes.NewReader(imgbytes))
	if err != nil {
		t.Errorf("Error during jpeg decode %v", err)
	}
	outputHash := imgsim.AverageHash(output)
	if outputHash != ref {
		t.Errorf("libjpeg-turbo output (%v) does not match reference (%v)", outputHash, ref)
	}
}

func BenchmarkJPEG(b *testing.B) {
	var err error
	var buf bytes.Buffer
	buf.Grow(10 << 20)
	for n := 0; n < b.N; n++ {
		_, err = EncodeJpeg(&buf, img, 75)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}

func BenchmarkJPEGTransform(b *testing.B) {
	sample, err := os.ReadFile("testdata/world-political.jpg")
	if !assert.NoError(b, err) {
		b.FailNow()
	}
	for n := 0; n < b.N; n++ {
		img, _, _, _, err := TransformJpeg(sample, true, DefaultScale)
		if !assert.NoError(b, err) {
			b.FailNow()
		}
		if !assert.Equal(b, 1891, img.Bounds().Dx()) {
			b.FailNow()
		}
		if !assert.Equal(b, 1081, img.Bounds().Dy()) {
			b.FailNow()
		}
	}
}
