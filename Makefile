.PHONY: build test clean dev web

build: web
	go build -o skm ./cmd/skm

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
