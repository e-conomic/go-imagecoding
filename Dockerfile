FROM golang:1.12-buster as builder

ENV GOLANGCI_VERSION "1.16.0"
ENV GOLANGCI_SHASUM "5343fc3ffcbb9910925f4047ec3c9f2e9623dd56a72a17ac76fb2886abc0976b"

WORKDIR /app

# Install libheif from Debian Bullseye, only required for optional libheif
# Install libjpeg-turbo from Debian Experimental
RUN apt-get update \
# Start libheif
  && wget -q http://deb.debian.org/debian/pool/main/libh/libheif/libheif-dev_1.5.0-1+b1_amd64.deb \
  && wget -q http://deb.debian.org/debian/pool/main/libh/libheif/libheif1_1.5.0-1+b1_amd64.deb \
  && wget -q http://deb.debian.org/debian/pool/main/x/x265/libx265-176_3.1.1-2_amd64.deb \
  && wget -q http://deb.debian.org/debian/pool/main/g/gcc-9/libstdc++6_9.2.1-4_amd64.deb \
  && wget -q http://deb.debian.org/debian/pool/main/g/gcc-9/gcc-9-base_9.2.1-4_amd64.deb \
# Start turbojpeg
  && wget -q http://deb.debian.org/debian/pool/main/libj/libjpeg-turbo/libturbojpeg0_2.0.2-1~exp2_amd64.deb \
  && wget -q http://deb.debian.org/debian/pool/main/libj/libjpeg-turbo/libturbojpeg0-dev_2.0.2-1~exp2_amd64.deb \
# Install dep packages
  && apt-get install -y ./*.deb \
  && apt-get install -y --no-install-recommends libwebp-dev libpng-dev autoconf libtool make nasm pkg-config libgomp1 \
  && apt-get clean

# Install GolangCI
RUN wget -q https://github.com/golangci/golangci-lint/releases/download/v$GOLANGCI_VERSION/golangci-lint-$GOLANGCI_VERSION-linux-amd64.tar.gz \
    && echo -n "$GOLANGCI_SHASUM  golangci-lint-$GOLANGCI_VERSION-linux-amd64.tar.gz" | shasum -c - \
    && tar xzf golangci-lint-$GOLANGCI_VERSION-linux-amd64.tar.gz \
    && rm golangci-lint-$GOLANGCI_VERSION-linux-amd64.tar.gz

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -tags heif ./...

RUN golangci-lint-$GOLANGCI_VERSION-linux-amd64/golangci-lint run
RUN go test -tags heif -v -race -cover -bench=. -benchmem ./...
