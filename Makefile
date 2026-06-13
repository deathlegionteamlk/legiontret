.PHONY: build build-all clean test install run

BINARY=legiontret
VERSION?=1.0.0
LDFLAGS=-ldflags "-s -w -X github.com/deathlegionteam/legiontret/internal/config.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/legiontret

build-all:
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/legiontret-linux-amd64 ./cmd/legiontret
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/legiontret-linux-arm64 ./cmd/legiontret
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/legiontret-darwin-amd64 ./cmd/legiontret
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/legiontret-darwin-arm64 ./cmd/legiontret
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/legiontret-windows-amd64.exe ./cmd/legiontret
	@echo "All builds complete!"
	@ls -lh dist/

clean:
	rm -f $(BINARY)
	rm -rf dist/

test:
	go vet ./...
	go test ./... || true

install: build
	cp $(BINARY) /usr/local/bin/

run: build
	./$(BINARY) run gemma3

docker:
	docker build -t legiontret .
	docker compose up -d

checksums:
	cd dist && sha256sum * > checksums-sha256.txt
