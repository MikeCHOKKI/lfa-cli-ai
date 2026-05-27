BINARY=bin/lfa
VERSION ?= $(shell grep 'Version = ' cmd/version.go | sed 's/.*"\(.*\)"/\1/')
LDFLAGS=-ldflags "-s -w -X github.com/lfa-cli/lfa-cli-ai/cmd.Version=$(VERSION)"
EXT ?=
BUILD_CMD=go build $(LDFLAGS)

.PHONY: all build test lint coverage cross clean run version release bump-patch bump-minor

all: lint test build

build:
	$(BUILD_CMD) -o $(BINARY)

build-matrix:
	$(BUILD_CMD) -o $(BINARY)-$(GOOS)-$(GOARCH)$(EXT)

test:
	go test ./... -v

lint:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@v0.6.0 ./...

coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

cross: clean
	GOOS=linux GOARCH=amd64 $(BUILD_CMD) -o $(BINARY)-linux-amd64
	GOOS=linux GOARCH=arm64 $(BUILD_CMD) -o $(BINARY)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD) -o $(BINARY)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(BUILD_CMD) -o $(BINARY)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(BUILD_CMD) -o $(BINARY)-windows-amd64.exe
	GOOS=windows GOARCH=arm64 $(BUILD_CMD) -o $(BINARY)-windows-arm64.exe

release: cross
	cd bin && sha256sum lfa-* > checksums.txt
	@echo "Release v$(VERSION) ready:"
	@ls -lh bin/

version:
	@echo "$(VERSION)"

bump-patch:
	$(eval NEW := $(shell echo $(VERSION) | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3+1}'))
	sed -i 's/Version = "$(VERSION)"/Version = "$(NEW)"/' cmd/version.go
	@echo "v$(VERSION) -> v$(NEW)"

bump-minor:
	$(eval NEW := $(shell echo $(VERSION) | awk -F. '{printf "%d.%d.0", $$1, $$2+1}'))
	sed -i 's/Version = "$(VERSION)"/Version = "$(NEW)"/' cmd/version.go
	@echo "v$(VERSION) -> v$(NEW)"

clean:
	rm -rf bin/ coverage.out coverage.html

run: build
	./$(BINARY)
