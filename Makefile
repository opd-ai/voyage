# Voyage Makefile
# Build targets for native, WebAssembly, and mobile builds

# Version information (override via environment or ldflags)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go settings
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
GOMOD := $(GOCMD) mod

# Output directories
BUILD_DIR := build
WASM_DIR := $(BUILD_DIR)/wasm
ANDROID_DIR := $(BUILD_DIR)/android
IOS_DIR := $(BUILD_DIR)/ios

# Source files
MAIN_PKG := ./cmd/voyage

# Test settings
TEST_FLAGS := -tags headless -race

# Mobile settings
ANDROID_SDK ?= $(ANDROID_HOME)
MIN_ANDROID_API ?= 21
APP_ID := ai.opd.voyage
APP_NAME := Voyage

.PHONY: all build build-wasm build-android build-ios clean test lint vet tidy run help serve-wasm mobile-setup

# Default target
all: build

# Build native binary
build:
	@echo "Building native binary..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/voyage $(MAIN_PKG)
	@echo "Build complete: $(BUILD_DIR)/voyage"

# Build WebAssembly binary
build-wasm:
	@echo "Building WebAssembly binary..."
	@mkdir -p $(WASM_DIR)
	GOOS=js GOARCH=wasm $(GOBUILD) $(LDFLAGS) -o $(WASM_DIR)/voyage.wasm $(MAIN_PKG)
	@if [ -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)/; \
	elif [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WASM_DIR)/; \
	elif [ -f "web/wasm_exec.js" ]; then \
		cp web/wasm_exec.js $(WASM_DIR)/; \
	else \
		echo "Downloading wasm_exec.js for Go 1.24..."; \
		curl -sL "https://raw.githubusercontent.com/golang/go/go1.24.0/lib/wasm/wasm_exec.js" -o $(WASM_DIR)/wasm_exec.js; \
	fi
	@cp web/index.html $(WASM_DIR)/ 2>/dev/null || $(MAKE) create-wasm-html
	@echo "WASM build complete: $(WASM_DIR)/voyage.wasm"
	@echo "Serve with: make serve-wasm"

# Create the HTML file for WASM if it doesn't exist
create-wasm-html:
	@echo "Error: web/index.html not found. Please create it first."
	@exit 1

# Serve WASM build locally for testing
serve-wasm: build-wasm
	@echo "Starting local server at http://localhost:8080"
	@echo "Press Ctrl+C to stop"
	@cd $(WASM_DIR) && python3 -m http.server 8080

# Run tests
test:
	@echo "Running tests..."
	xvfb-run -a $(GOTEST) $(TEST_FLAGS) ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	xvfb-run -a $(GOTEST) $(TEST_FLAGS) -v ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) -tags headless ./...

# Run linting (if golangci-lint is available)
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, running go vet only"; \
		$(MAKE) vet; \
	fi

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Run the game
run: build
	./$(BUILD_DIR)/voyage

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f voyage
	@echo "Clean complete"

# Setup gomobile for mobile builds
mobile-setup:
	@echo "Setting up gomobile..."
	@if ! command -v gomobile > /dev/null; then \
		echo "Installing gomobile..."; \
		go install golang.org/x/mobile/cmd/gomobile@latest; \
		gomobile init; \
	else \
		echo "gomobile already installed"; \
	fi

# Build Android APK
build-android: mobile-setup
	@echo "Building Android APK..."
	@mkdir -p $(ANDROID_DIR)
	@if [ -z "$(ANDROID_SDK)" ]; then \
		echo "Error: ANDROID_HOME or ANDROID_SDK not set"; \
		echo "Please install Android SDK and set ANDROID_HOME"; \
		exit 1; \
	fi
	gomobile build -target=android/arm64,android/arm \
		-androidapi $(MIN_ANDROID_API) \
		-o $(ANDROID_DIR)/$(APP_NAME).apk \
		$(MAIN_PKG)
	@echo "Android APK built: $(ANDROID_DIR)/$(APP_NAME).apk"

# Build iOS app (requires macOS with Xcode)
build-ios: mobile-setup
	@echo "Building iOS app..."
	@mkdir -p $(IOS_DIR)
	@if [ "$$(uname)" != "Darwin" ]; then \
		echo "Error: iOS builds require macOS with Xcode"; \
		exit 1; \
	fi
	gomobile build -target=ios \
		-bundleid $(APP_ID) \
		-o $(IOS_DIR)/$(APP_NAME).app \
		$(MAIN_PKG)
	@echo "iOS app built: $(IOS_DIR)/$(APP_NAME).app"

# Build iOS framework for integration
build-ios-framework: mobile-setup
	@echo "Building iOS framework..."
	@mkdir -p $(IOS_DIR)
	@if [ "$$(uname)" != "Darwin" ]; then \
		echo "Error: iOS builds require macOS with Xcode"; \
		exit 1; \
	fi
	gomobile bind -target=ios \
		-bundleid $(APP_ID) \
		-o $(IOS_DIR)/Voyage.xcframework \
		$(MAIN_PKG)
	@echo "iOS framework built: $(IOS_DIR)/Voyage.xcframework"

# Build Android AAR library
build-android-aar: mobile-setup
	@echo "Building Android AAR..."
	@mkdir -p $(ANDROID_DIR)
	gomobile bind -target=android \
		-androidapi $(MIN_ANDROID_API) \
		-o $(ANDROID_DIR)/voyage.aar \
		$(MAIN_PKG)
	@echo "Android AAR built: $(ANDROID_DIR)/voyage.aar"

# Show help
help:
	@echo "Voyage Build Targets"
	@echo "===================="
	@echo ""
	@echo "Building:"
	@echo "  make build          - Build native binary"
	@echo "  make build-wasm     - Build WebAssembly binary for browsers"
	@echo "  make build-android  - Build Android APK (requires Android SDK)"
	@echo "  make build-ios      - Build iOS app (requires macOS + Xcode)"
	@echo ""
	@echo "Mobile Libraries:"
	@echo "  make build-android-aar - Build Android AAR library"
	@echo "  make build-ios-framework - Build iOS XCFramework"
	@echo "  make mobile-setup   - Install gomobile if not present"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run all tests"
	@echo "  make test-verbose - Run all tests with verbose output"
	@echo "  make vet          - Run go vet"
	@echo "  make lint         - Run linting"
	@echo ""
	@echo "Running:"
	@echo "  make run          - Build and run the game"
	@echo "  make serve-wasm   - Build WASM and serve locally on port 8080"
	@echo ""
	@echo "Maintenance:"
	@echo "  make tidy         - Tidy Go modules"
	@echo "  make clean        - Remove build artifacts"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  BUILD_TIME=$(BUILD_TIME)"
	@echo "  ANDROID_SDK=$(ANDROID_SDK)"
	@echo "  MIN_ANDROID_API=$(MIN_ANDROID_API)"
