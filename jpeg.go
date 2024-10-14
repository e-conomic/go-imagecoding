package imagecoding

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"unsafe"

	"go.uber.org/zap"
)

// #cgo pkg-config: libturbojpeg
// #include <stdlib.h>
// #include <turbojpeg.h>
// int goTjCompress2(tjhandle handle, const unsigned char *srcBuf,
//                   int width, int pitch, int height, int pixelFormat,
//                   unsigned char *jpegBuf, unsigned long *jpegSize,
//                   int jpegSubsamp, int jpegQual, int flags) {
//      return tjCompress2(handle, srcBuf, width, pitch, height,
//                         pixelFormat, &jpegBuf, jpegSize,
//                         jpegSubsamp, jpegQual, flags | TJFLAG_NOREALLOC);
//}
import "C"

type TurboJpegOperation C.int

func EncodeJpeg(buf *bytes.Buffer, img image.Image, quality int) ([]byte, error) {
	var pix []uint8
	var format, stride, jpegSubsamp, cWidth, cHeight, flags, jpegQual, res C.int

	cWidth = C.int(img.Bounds().Dx())
	cHeight = C.int(img.Bounds().Dy())
	switch v := img.(type) {
	case *image.Gray:
		pix = v.Pix
		stride = C.int(v.Stride)
		format = C.TJPF_GRAY
		jpegSubsamp = C.TJSAMP_GRAY
	case *RGBImage:
		pix = v.Pix
		stride = C.int(v.Stride)
		format = C.TJPF_RGB
		jpegSubsamp = C.TJSAMP_420
	case *image.RGBA:
		pix = v.Pix
		stride = C.int(v.Stride)
		format = C.TJPF_RGBX
		jpegSubsamp = C.TJSAMP_420
	default:
		return nil, errors.New("unsupported image type")
	}

	tjHandle := C.tjInitCompress()
	if tjHandle == nil {
		return nil, fmt.Errorf("could not init libjpeg-turbo: %v", C.GoString(C.tjGetErrorStr2(tjHandle)))
	}
	defer C.tjDestroy(tjHandle)

	jpegQual = C.int(quality)
	flags = C.TJFLAG_NOREALLOC
	imageBytes := C.tjBufSize(cWidth, cHeight, jpegSubsamp)

	// Prepare a buffer
	buf.Reset()
	buf.Grow(int(imageBytes))
	buf.WriteByte(0)

	res = C.goTjCompress2(
		tjHandle,
		(*C.uchar)(unsafe.Pointer(&pix[0])),
		cWidth, stride, cHeight,
		format,
		(*C.uchar)(unsafe.Pointer(&(buf.Bytes())[0])),
		&imageBytes,
		jpegSubsamp, jpegQual, flags,
	)
	if res != 0 {
		if C.tjGetErrorCode(tjHandle) == C.TJERR_WARNING {
			zap.L().Warn(
				"jpeg compress warning",
				zap.String("jpgerror", C.GoString(C.tjGetErrorStr2(tjHandle))),
			)
		} else {
			return nil, fmt.Errorf("could not compress jpeg: %v", C.GoString(C.tjGetErrorStr2(tjHandle)))
		}
	}

	return buf.Bytes()[:int(imageBytes)], nil
}

func getTransformOperation(orient Orientation) []TurboJpegOperation {
	switch orient {
	case TopLeft:
		return []TurboJpegOperation{C.TJXOP_NONE}
	case TopRight:
		return []TurboJpegOperation{C.TJXOP_HFLIP}
	case BottomRight:
		return []TurboJpegOperation{C.TJXOP_ROT180}
	case BottomLeft:
		return []TurboJpegOperation{C.TJXOP_VFLIP}
	case LeftTop:
		return []TurboJpegOperation{C.TJXOP_TRANSPOSE}
	case RightTop:
		return []TurboJpegOperation{C.TJXOP_ROT90}
	case RightBottom:
		return []TurboJpegOperation{C.TJXOP_TRANSVERSE}
	case LeftBottom:
		return []TurboJpegOperation{C.TJXOP_ROT270}
	default:
		fmt.Printf("Unexpected orientation %v", orient)
		return []TurboJpegOperation{C.TJXOP_NONE}
	}
}

// ReOrientJpeg will transform a JPEG into a top left (normal) orientation
// It returns a buffer with JPEG encoding
func ReOrientJpeg(file []byte, orient Orientation) ([]byte, error) {
	if len(file) == 0 {
		return nil, ErrEmptyInput
	}
	// Init a transform
	tjHandle := C.tjInitTransform()
	if tjHandle == nil {
		return nil, fmt.Errorf("could not init libjpeg-turbo: %v", C.GoString(C.tjGetErrorStr2(tjHandle)))
	}
	defer C.tjDestroy(tjHandle)

	operations := getTransformOperation(orient)

	var destBuf *C.uchar
	var destSize C.ulong
	for _, operation := range operations {
		transform := C.tjtransform{
			op:      C.int(operation),
			options: C.TJXOPT_GRAY,
		}
		res := C.tjTransform(tjHandle,
			(*C.uchar)(unsafe.Pointer(&file[0])), C.ulong(len(file)),
			1,
			&destBuf,
			&destSize,
			&transform,
			0, // nolint:gocritic,staticcheck
		)
		if res != 0 {
			if C.tjGetErrorCode(tjHandle) == C.TJERR_WARNING {
				zap.L().Warn(
					"jpeg transform warning",
					zap.String("jpgerror", C.GoString(C.tjGetErrorStr())),
				)
			} else {
				return nil, fmt.Errorf("could not transform jpeg: %v", C.GoString(C.tjGetErrorStr2(tjHandle)))
			}
		}
	}
	result := C.GoBytes(unsafe.Pointer(destBuf), C.int(destSize))
	C.tjFree(destBuf)
	return result, nil
}

func ConfigJpeg(data []byte) (image.Config, string, error) {
	if len(data) == 0 {
		return image.Config{}, string(Jpeg), ErrEmptyInput
	}

	// Init Turbo-JPEG Decompression
	tjHandle := C.tjInitDecompress()
	if tjHandle == nil {
		return image.Config{}, string(Jpeg), fmt.Errorf("could not init libjpeg-turbo: %v", C.GoString(C.tjGetErrorStr2(tjHandle)))
	}
	defer C.tjDestroy(tjHandle)

	var cWidth, cHeight, jpegsubsamp, jpegcolorspace C.int
	res := C.tjDecompressHeader3(
		tjHandle,
		(*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(len(data)),
		&cWidth,
		&cHeight,
		&jpegsubsamp,
		&jpegcolorspace,
	)
	if res != 0 {
		return image.Config{}, string(Jpeg), fmt.Errorf("could not read JPEG header: %v", C.GoString(C.tjGetErrorStr2(tjHandle)))
	}

	var model color.Model
	switch jpegcolorspace {
	case C.TJCS_RGB:
		model = RGBModel
	case C.TJCS_YCbCr:
		model = color.YCbCrModel
	case C.TJCS_GRAY:
		model = color.GrayModel
	case C.TJCS_CMYK, C.TJCS_YCCK:
		model = color.CMYKModel
	default:
		return image.Config{}, string(Jpeg), fmt.Errorf("unknown jpeg color space: %d", jpegcolorspace)
	}

	return image.Config{
		ColorModel: model,
		Width:      int(cWidth),
		Height:     int(cHeight),
	}, string(Jpeg), nil
}

// TransformJpeg will scale and colormap an input JPEG file to an image.Gray or RGBImage
// This will use libjpeg-turbo to do it as efficiently as possible, utilizing DCT factors for fast scaling
func TransformJpeg(data []byte, grayscale bool, scale ScaleFunc) (out image.Image, width, height int, scaleFactor float64, err error) {
	if len(data) == 0 {
		return nil, 0, 0, -1, ErrEmptyInput
	}

	// Init Turbo-JPEG Decompression
	tjHandle := C.tjInitDecompress()
	if tjHandle == nil {
		return nil, 0, 0, -1, fmt.Errorf("could not init libjpeg-turbo: %v", C.GoString(C.tjGetErrorStr()))
	}
	defer C.tjDestroy(tjHandle)

	// Detect orientation
	orientation := GetOrientation(bytes.NewReader(data))
	if orientation != TopLeft {
		data, err = ReOrientJpeg(data, orientation)
		if err != nil {
			return nil, 0, 0, -1, err
		}
	}

	// Read the size of the jpeg
	conf, _, err := ConfigJpeg(data)
	if err != nil {
		return nil, 0, 0, -1, err
	}
	width = conf.Width
	height = conf.Height

	// Calculate our preferred scaling factor
	_, _, prefScaleFactor := scale(width, height)

	// Find the closest match for a DCT scaling
	var cNumScaleFactor C.int
	cScaleFactors := C.tjGetScalingFactors(&cNumScaleFactor)
	if cScaleFactors == nil {
		return nil, 0, 0, -1, errors.New("could not get libjpeg-turbo scale factors")
	}
	scaleFactors := uintptr(unsafe.Pointer(cScaleFactors))

	// Find the closest JPEG DCT scale factor
	selectedScaleFactorDiff := math.MaxFloat64
	var selectedScaleFactor uintptr
	var jpegScaleFactor float64
	for i := 0; i < int(cNumScaleFactor); i++ {
		offset := uintptr(C.sizeof_tjscalingfactor * i)
		sf := (*C.tjscalingfactor)(unsafe.Pointer(scaleFactors + offset))
		jpegScaleFactor = float64(sf.num) / float64(sf.denom)
		diff := math.Abs(prefScaleFactor - jpegScaleFactor)
		if diff < selectedScaleFactorDiff {
			selectedScaleFactorDiff = diff
			selectedScaleFactor = offset
		}
	}

	// Calculate the final image size
	sf := (*C.tjscalingfactor)(unsafe.Pointer(scaleFactors + selectedScaleFactor))
	scaledW := int(math.RoundToEven(float64(C.int(width)*sf.num+sf.denom-1) / float64(sf.denom)))
	scaledH := int(math.RoundToEven(float64(C.int(height)*sf.num+sf.denom-1) / float64(sf.denom)))
	scaleFactor = float64(sf.num) / float64(sf.denom)

	// Calculate the image stride and pitch
	var pixelFormat C.int
	var pitch int
	if grayscale {
		pixelFormat = C.TJPF_GRAY
		pitch = scaledW * 1 // C.tjPixelSize[C.TJPF_GRAY]
	} else {
		pixelFormat = C.TJPF_RGB
		pitch = scaledW * 3 // C.tjPixelSize[C.TJPF_RGB]
	}
	pixsize := scaledH * pitch
	buf := make([]uint8, pixsize)

	res := C.tjDecompress2(
		tjHandle,
		(*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(len(data)),
		(*C.uchar)(unsafe.Pointer(&buf[0])),
		C.int(scaledW),
		C.int(pitch),
		C.int(scaledH),
		pixelFormat,
		0,
	)
	if res != 0 {
		if C.tjGetErrorCode(tjHandle) == C.TJERR_WARNING {
			zap.L().Warn(
				"jpeg decompress warning",
				zap.String("jpgerror", C.GoString(C.tjGetErrorStr())),
			)
		} else {
			return nil, 0, 0, -1, fmt.Errorf("could not decompress jpeg: %v", C.GoString(C.tjGetErrorStr()))
		}
	}
	var img image.Image
	if grayscale {
		img = &image.Gray{
			Pix:    buf,
			Stride: pitch,
			Rect:   image.Rect(0, 0, scaledW, scaledH),
		}
	} else {
		img = &RGBImage{
			Pix:    buf,
			Stride: pitch,
			Rect:   image.Rect(0, 0, scaledW, scaledH),
		}
	}
	return img, width, height, scaleFactor, nil
}
