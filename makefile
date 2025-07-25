VERSION=$(shell git describe --tags --always)

.PHONY: build
# 一次性编译 amd64 和 arm64 的 Linux 二进制
build:
	mkdir -p bin/ && \
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/cmm-amd64 ./... && \
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/cmm-arm64 ./...