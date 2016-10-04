GOPACKAGES=$(shell glide novendor)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all
all: deps fmt vet test

.PHONY: deps
deps:
	glide install

.PHONY: fmt
fmt:
	@if [ -n "$$(gofmt -l ${GOFILES})" ]; then echo 'Please run gofmt -l -w on your code.' && exit 1; fi

.PHONY: test
test:
	go test -v -race ${GOPACKAGES}

.PHONY: vet
vet:
	go vet ${GOPACKAGES}
