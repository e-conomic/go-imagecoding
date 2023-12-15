package imagecoding

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
	"testing"

	"github.com/Nr90/imgsim"
)

func TestJpegExifReference(t *testing.T) {
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
			img, err := jpeg.Decode(bytes.NewReader(jpegbytes))
			if err != nil {
				t.Errorf("Error during jpeg decode %v", err)
			}

			orient := GetOrientation(bytes.NewReader(jpegbytes))
			fixedimg := FixOrientation(img, orient)
			hash := imgsim.AverageHash(fixedimg)
			if hash != refhash {
				t.Errorf("Output does not match reference %v", hash)
			}
		})
	}
}
