package imagecoding

import (
	"bytes"
	"golang.org/x/image/webp"
	"testing"

	"github.com/Nr90/imgsim"
)

func TestEncodeWebp(t *testing.T) {
	var buf bytes.Buffer
	imgbytes, err := EncodeWebP(&buf, img)
	if err != nil {
		t.Errorf("Error during libwebp encode %v", err)
	}

	output, err := webp.Decode(bytes.NewReader(imgbytes))
	if err != nil {
		t.Errorf("Error during webp decode %v", err)
	}
	outputHash := imgsim.AverageHash(output)
	if outputHash != ref {
		t.Errorf("libwebp output (%v) does not match reference (%v)", outputHash, ref)
	}
}

func BenchmarkWebp(b *testing.B) {
	var err error
	var buf bytes.Buffer
	buf.Grow(10 << 20)
	for n := 0; n < b.N; n++ {
		_, err = EncodeWebP(&buf, img)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
