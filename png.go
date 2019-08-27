package imagecoding

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"unsafe"
)

// -Wno-incompatible-pointer-types: ignore a warning from cgo https://github.com/golang/go/issues/19832

// #cgo pkg-config: libpng
// #cgo CFLAGS: -Wno-incompatible-pointer-types
// #include <stdlib.h>
// #include <png.h>
// int cpng_image_data_size(const png_image *image) {
//	    return PNG_IMAGE_DATA_SIZE(*image);
// }
import "C"

// EncodePng will encode an image.Gray to PNG bytes, using libpng's simplified API for performance
func EncodePng(buf *bytes.Buffer, img image.Image) ([]byte, error) {
	var pix []uint8
	var stride int
	var format C.png_uint_32

	switch v := img.(type) {
	case *image.Gray:
		pix = v.Pix
		stride = v.Stride
		format = C.PNG_FORMAT_GRAY
	case *RGBImage:
		pix = v.Pix
		stride = v.Stride
		format = C.PNG_FORMAT_RGB
	case *image.RGBA:
		pix = v.Pix
		stride = v.Stride
		format = C.PNG_FORMAT_RGBA
	default:
		return nil, errors.New("unsupported image type")
	}

	var op C.png_controlp
	info := &C.png_image{
		opaque:           op,
		version:          C.PNG_IMAGE_VERSION,
		width:            C.uint(img.Bounds().Dx()),
		height:           C.uint(img.Bounds().Dy()),
		format:           format,
		flags:            C.PNG_IMAGE_FLAG_FAST,
		colormap_entries: 0,
	}
	defer C.png_image_free(info)

	// Get memory size
	var result C.int
	var imageBytes C.png_alloc_size_t

	// Make a buffer of sufficient size
	maxSize := C.cpng_image_data_size(info)

	// Prepare a buffer
	buf.Reset()
	buf.Grow(int(maxSize))
	imageBytes = C.png_alloc_size_t(maxSize)
	buf.WriteByte(0)

	// Write the actual PNG
	result = C.png_image_write_to_memory(
		info,
		unsafe.Pointer(&(buf.Bytes())[0]),
		&imageBytes,
		0,
		unsafe.Pointer(&pix[0]),
		C.int(stride),
		nil,
	)
	if result == 0 {
		return nil, fmt.Errorf("libpng threw an error (%x) %q", info.warning_or_error, C.GoString(&info.message[:][0]))
	}

	return buf.Bytes()[:int(imageBytes)], nil
}
