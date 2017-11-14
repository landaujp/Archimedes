# メタ情報
NAME     := archimedes
VERSION  := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLGAS  := -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'

# 必要なツール類をセットアップする
## Setup
setup:
	go get github.com/Masterminds/glide
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Songmu/make2help/cmd/make2help

# テストを実行する
## Run tests
test: deps
	go test $$(glide novendor)

# glideを使って依存パッケージをインストールする
## Install dependencies
deps: setup
	glide install

## Update dependencies
update: setup
	glide udpate

## Lint
lint: setup
	go vet $$(glide novendor)
	for pkg in $$(glide novendor -x); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

## Format source codes
fmt: setup
	goimports -w $$(glide nv -x)

## build binaries ex. make bin/archimedes
bin/%: cmd/%/main.go deps
	go generate $< && go build -ldflags "$(LDFLGAS)" -o $@ $< cmd/$*/bindata.go

## build binaries for Mac
darwin-build: cmd/last/main.go cmd/depth/main.go deps
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go generate cmd/last/main.go  &&  go build -ldflags "$(LDFLGAS)" -o bin/darwin-amd64/last  cmd/last/main.go  cmd/last/bindata.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go generate cmd/depth/main.go &&  go build -ldflags "$(LDFLGAS)" -o bin/darwin-amd64/depth cmd/depth/main.go cmd/depth/bindata.go

## build binaries for Linux
linux-build: cmd/last/main.go cmd/depth/main.go deps
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go generate cmd/last/main.go  &&  go build -ldflags "$(LDFLGAS)" -o bin/linux-amd64/last  cmd/last/main.go  cmd/last/bindata.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go generate cmd/depth/main.go &&  go build -ldflags "$(LDFLGAS)" -o bin/linux-amd64/depth cmd/depth/main.go cmd/depth/bindata.go

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps udpate test lint help
