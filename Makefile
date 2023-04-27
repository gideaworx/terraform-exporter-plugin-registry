mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

BUILD_DIR := $(mkfile_dir)/build
WEB_DIR := $(BUILD_DIR)/web

GO := $(shell command -v go)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(WEB_DIR):
	mkdir -p $(WEB_DIR)

$(BUILD_DIR)/plugin-registry: $(BUILD_DIR)
	$(GO) build -trimpath -o $(BUILD_DIR)/plugin-registry ./backend/cmd/registry-mgmt/main.go

.PHONY: build-web build clean

build: $(BUILD_DIR)/plugin-registry

build-web: $(WEB_DIR) $(BUILD_DIR)/plugin-registry
	$(BUILD_DIR)/plugin-registry build-site -o $(BUILD_DIR)/web

clean:
	rm -fr $(BUILD_DIR)
