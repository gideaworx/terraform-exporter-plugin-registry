GO := $(shell command -v go)

.PHONY: build-index generate-index

build-index:
	$(GO) build -o build-index ./cmd/main.go

generate-index: build-index
	./build-index
