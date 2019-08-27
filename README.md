# Go Image Coding
[![Go Report Card](https://goreportcard.com/badge/github.com/e-conomic/go-imagecoding)](https://goreportcard.com/report/github.com/e-conomic/go-imagecoding)
[![go-doc](https://godoc.org/github.com/e-conomic/go-imagecoding?status.svg)](https://godoc.org/github.com/e-conomic/go-imagecoding)

Go (bindings for) image en-/de-coding using libraries found on debian and other
common systems. Go comes with image processing built-in, written in go, but
often they can't quite match the performance of established image processing
libraries such as jpeg-turbo etc.

While general purpose, this project is opinionated and tailor to our needs.

### Install (debian)

See Dockerfile

### Install (OS X)

```bash
brew install jpeg-turbo libpng libheif webp pkg-config
```

### HEIF/HEIC

This package optionally supports heif, to include heif; add `-tags heif` to your gobuild. It's enabled by default on darwin (macOS).