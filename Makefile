.PHONY: all clean build

all: clean build

clean:
	rm -f build/*

deps:
	glide install

build: build-linux-amd64

build-linux-amd64:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/go-geoip2-httpd -ldflags "-s"
