GO = go
ARCH = amd64
OS = linux
GO_ENV = GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0
BIN_NAME = s3downloader-$(OS)-$(ARCH)

.PHONY: build-linux

build-linux: main.go
	$(GO_ENV) $(GO) build -o bin/$(BIN_NAME) $<
	chmod +x bin/$(BIN_NAME)

.PHONY: release
release:
	./scripts/release.sh
