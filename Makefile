# 仓库根目录
REPODIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
# 构建输出目录
BUILDDIR := $(REPODIR)/dist

.PHONY: build
build:
	@echo "Cleaning up and downloading modules..."
	go mod tidy
	@echo "Building for Linux amd64..."
	mkdir -p $(BUILDDIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/app $(REPODIR)
	@echo "Build completed: $(BUILDDIR)/app"

run:
	cd $(REPODIR) && ./dist/app
