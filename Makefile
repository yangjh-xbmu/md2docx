.PHONY: build test lint clean release sync-styles

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X github.com/yangjh-xbmu/md2docx/internal/cli.Version=$(VERSION)"

# Sync embedded styles before build
sync-styles:
	@rm -rf cmd/md2docx/styles
	@cp -r styles cmd/md2docx/styles

build: sync-styles
	go build $(LDFLAGS) -o md2docx ./cmd/md2docx/

test: sync-styles
	go test -race ./...

lint:
	go vet ./...

clean:
	rm -f md2docx
	rm -rf dist/

release: sync-styles
	@mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/md2docx-darwin-arm64 ./cmd/md2docx/
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/md2docx-darwin-amd64 ./cmd/md2docx/
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/md2docx-linux-amd64 ./cmd/md2docx/
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/md2docx-windows-amd64.exe ./cmd/md2docx/
	@cd dist && shasum -a 256 md2docx-* > sha256sums.txt
	@echo "Release binaries and checksums in dist/"
