mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

BUILD_DIR := $(mkfile_dir)/build

GO := $(shell command -v go)

.PHONY: build-dependencies build-registry-index build-web build clean

build: build-dependencies
	$(GO) build -trimpath -o $(BUILD_DIR)/plugin-registry ./main.go

build-dependencies: build-registry-index

build-registry-index: build-web
	BUILD_PATH=$(BUILD_DIR)/web $(MAKE) -C registry generate-index

build-web:
	BUILD_PATH=$(BUILD_DIR)/web $(MAKE) -C web/registry-site

clean:
	rm -fr $(BUILD_DIR)
