mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

BUILD_DIR := $(mkfile_dir)/build

GO := $(shell command -v go)
TMPINDEX := $(shell mktemp)

.PHONY: build-dependencies build-registry-index build-web build clean run

build: build-dependencies
	$(GO) build -trimpath -o $(BUILD_DIR)/plugin-registry ./main.go

build-dependencies: build-web

build-registry-index:
	mkdir -p $(BUILD_DIR)/web
	BUILD_PATH=$(BUILD_DIR)/web $(MAKE) -C registry generate-index

build-web: build-registry-index
	cp $(BUILD_DIR)/web/index.yaml $(TMPINDEX)
	BUILD_PATH=$(BUILD_DIR)/web $(MAKE) -C web/registry-site build
	mv $(TMPINDEX) $(BUILD_DIR)/web/index.yaml

clean:
	rm -fr $(BUILD_DIR)

run: clean build
	$(BUILD_DIR)/plugin-registry serve
