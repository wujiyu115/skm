.PHONY: build test clean dev web install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build: web
	go build -ldflags "-X github.com/ejoy/skm/internal/cli.Version=$(VERSION)" -o skm ./cmd/skm

install: web
	go install -ldflags "-X github.com/ejoy/skm/internal/cli.Version=$(VERSION)" ./cmd/skm

web:
	cd web && npm ci && npm run build
	rm -rf internal/server/dist
	cp -r web/dist internal/server/dist

test:
	go test ./... -v

dev:
	SKM_DEV=1 go run ./cmd/skm serve &
	cd web && npx vite

clean:
	rm -f skm
	rm -rf internal/server/dist web/dist
