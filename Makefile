BINARY=bin/lfa
LDFLAGS=-ldflags "-s -w -X github.com/lfa-cli/lfa-cli-ai/cmd.Version=0.1.0"
BUILD_CMD=go build $(LDFLAGS) -o $(BINARY)

.PHONY: build test lint coverage cross clean run

build:
	$(BUILD_CMD)

test:
	go test ./... -v

lint:
	staticcheck ./...

coverage:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

cross:
	GOOS=linux GOARCH=amd64 $(BUILD_CMD)-linux-amd64
	GOOS=linux GOARCH=arm64 $(BUILD_CMD)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(BUILD_CMD)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(BUILD_CMD)-windows-amd64.exe
	GOOS=windows GOARCH=arm64 $(BUILD_CMD)-windows-arm64.exe

clean:
	rm -rf bin/ coverage.out coverage.html

run: build
	./$(BINARY)
