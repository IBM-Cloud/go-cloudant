GO111MODULE := on
export
GOPACKAGES=$(shell go list ./... | grep -v /vendor/)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all
all: deps fmt vet test

.PHONY: updatedeps
updatedeps:
	go get -u=patch ./...
	go mod tidy

.PHONY: fmt
fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

.PHONY: test
test:
	go test -race -covermode=atomic -coverprofile=cover.out ./...
	# go test -v -race ${GOPACKAGES}

.PHONY: vet
vet:
	go vet ${GOPACKAGES}
