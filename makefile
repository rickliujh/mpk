MODULE = github.com/rickliujh/mpk
APP_NAME := mpk
# Directory where the binary will be placed
BUILD_DIR := ./build
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*") # Find all Go files except in vendor

# Default target
all: clean build

# Ensure build directory exists, then build the Go binary
build: $(BUILD_DIR)/$(APP_NAME)

$(BUILD_DIR)/$(APP_NAME): $(GO_FILES)
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -v -o $(BUILD_DIR)/$(APP_NAME)
	@echo "Binary built at $(BUILD_DIR)/$(APP_NAME)"

# Clean the build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Tidy Go modules (optional)
tidy:
	@go mod tidy

# Format the Go code (optional)
fmt:
	@go fmt ./...

# Run the application (optional)
run:
	@go run .

# Phony targets to avoid name collisions
.PHONY: all build clean tidy fmt run
