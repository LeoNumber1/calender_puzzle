
.PHONY: build container build-local build-linux build-windows

# Project output directory.
OUTPUT_DIR := ./bin

# Build directory.
BUILD_DIR := ./build

NAME := puzzle
IMAGE_NAME := puzzle
CUR_PWD := $(shell pwd)

AUTHOR := $(shell git log --pretty=format:"%an"|head -n 1)
VERSION := $(shell git rev-list HEAD | head -1)
BUILD_DATE := $(shell date +%Y%m%d%H%M%S)

# Track code version with Docker Label.
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

export GO111MODULE=on

build: build-local build-linux build-windows

build-local:
	go build -v -o $(OUTPUT_DIR)/$(NAME)-mac

build-linux:
	GOOS=linux go build -v -o $(OUTPUT_DIR)/$(NAME)-linux

build-windows:
	GOOS=windows go build -v -o $(OUTPUT_DIR)/$(NAME)-win.exe

container:
	@docker build -t $(IMAGE_NAME):$(VERSION)-$(AUTHOR)-$(BUILD_DATE)                \
    	  --label $(DOCKER_LABELS)                                             \
    	  -f $(BUILD_DIR)/Dockerfile .;

clean:
	rm -rf $(OUTPUT_DIR)