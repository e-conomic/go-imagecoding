FROM golang:1.13-buster as builder

ENV GOLANGCI_VERSION "1.18.0"
ENV GOLANGCI_SHASUM "0ef2c502035d5f12d6d3a30a7c4469cfcae4dd3828d15fbbfb799c8331cd51c4"

WORKDIR /app

# Install libheif from Debian Bullseye, only required for optional libheif
# Install libjpeg-turbo from Debian Experimental
RUN \
  echo "deb http://deb.debian.org/debian bullseye main" | tee -a /etc/apt/sources.list \
  && echo "deb http://deb.debian.org/debian experimental main" | tee -a /etc/apt/sources.list \
  && apt-get update \
# Start libheif
  && apt-get install -t bullseye -y --no-install-recommends libheif-dev \
# Start turbojpeg
  && apt-get install -t experimental -y --no-install-recommends libturbojpeg0-dev \
# Install dep packages
  && apt-get install -t buster -y --no-install-recommends libwebp-dev libpng-dev autoconf libtool make nasm pkg-config libgomp1 \
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
