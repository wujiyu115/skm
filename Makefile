.PHONY: build test clean

build:
	go build -o skm ./cmd/skm

test:
	go test ./... -v

clean:
	rm -f skm
