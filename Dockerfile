FROM golang:1.15-buster as builder

ENV GOLANGCI_VERSION "1.30.0"
ENV GOLANGCI_SHASUM "c8e8fc5753e74d2eb489ad428dbce219eb9907799a57c02bcd8b80b4b98c60d4"

WORKDIR /app

# Install libheif from Debian Bullseye, only required for optional libheif
# Install libjpeg-turbo from Debian Experimental
RUN \
  echo "deb http://deb.debian.org/debian bullseye main" | tee -a /etc/apt/sources.list \
  && echo "deb http://deb.debian.org/debian experimental main" | tee -a /etc/apt/sources.list \
  && apt-get update \
  && apt-get install -t bullseye -y --no-install-recommends gcc-8-base libgcc-8-dev\
# Start libheif
  && apt-get install -t bullseye -y -o APT::Immediate-Configure=0 --no-install-recommends libheif-dev \
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
