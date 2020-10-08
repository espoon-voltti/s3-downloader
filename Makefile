GO			= go
ARCH 		= amd64
OS 			= linux
GO_ENV 		= GOOS=$(OS) GOARCH=$(ARCH) GO111MODULE=on CGO_ENABLED=0
BIN_NAME 	= s3downloader-$(OS)-$(ARCH)
BUILD_PATH 	= ./bin/$(BIN_NAME)

.PHONY: tidy phony release

fmt:
	@$(GO_ENV) $(GO) fmt .

tidy:
	@$(GO_ENV) $(GO) mod tidy

build: main.go
	$(GO_ENV) $(GO) build -o $(BUILD_PATH) $<
	chmod +x $(BUILD_PATH)

release:
	./scripts/release.sh

.DEFAULT_GOAL := build
