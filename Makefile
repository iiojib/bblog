.PHONY: all bookmarklet-js client-js bookmarklet platform clean build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64

BUILD_DIR := build
BOOKMARKLET_JS := $(BUILD_DIR)/bookmarklet.js
CLIENT_JS := $(BUILD_DIR)/client.js
BOOKMARKLET_TXT := $(BUILD_DIR)/bookmarklet
GO_BIN := bblog

all: platform bookmarklet client-js

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

bookmarklet-js: $(BUILD_DIR)
	esbuild --minify --bundle --outfile=$(BOOKMARKLET_JS) bookmarklet.ts

client-js: $(BUILD_DIR)
	esbuild --minify --bundle --outfile=$(CLIENT_JS) inject.ts

bookmarklet: bookmarklet-js build_bookmarklet.sh
	bash ./build_bookmarklet.sh $(BOOKMARKLET_JS) $(BOOKMARKLET_TXT)

platform: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/darwin_amd64/$(GO_BIN) .

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/darwin_arm64/$(GO_BIN) .

build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/linux_amd64/$(GO_BIN) .

build-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/linux_arm64/$(GO_BIN) .

clean:
	rm -rf $(BUILD_DIR)
