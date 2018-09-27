SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export GO111MODULE := on
export GOBIN       := $(PWD)/bin
export PATH        := $(PWD)/bin:$(PATH)

setup:
	curl -sfL https://git.io/vp6lP | sh

test:
	go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN)

cover: test
	go tool cover -html=coverage.out

fmt:
	find . -name '*.go' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint:
	gometalinter --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...

ci: build lint test

build:
	go build -o org-stats .

.DEFAULT_GOAL := build
