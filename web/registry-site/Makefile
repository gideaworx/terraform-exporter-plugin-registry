mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

yarn := $(mkfile_dir)/.yarn/releases/yarn-3.5.0.cjs

.PHONY: build

build: 
	$(yarn) install --immutable
	$(yarn) build
