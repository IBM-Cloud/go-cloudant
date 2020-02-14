GOPACKAGES=$(shell glide novendor)
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOLANGCI_LINT_EXISTS:=$(shell golangci-lint --version 2>/dev/null)

.PHONY: all
all: deps fmt vet test

.PHONY: deps
deps:
	glide install

# Format the code
.PHONY: fmt
fmt:
ifdef GOLANGCI_LINT_EXISTS
	golangci-lint run --disable-all --enable=gofmt --fix
else
	@echo "golangci-lint is not installed"
endif

.PHONY: lint
lint:
ifdef GOLANGCI_LINT_EXISTS
	golangci-lint run
else
	@echo "golangci-lint is not installed"
endif

.PHONY: test
test:
	go test -v -race ${GOPACKAGES}

.PHONY: vet
vet:
	go vet ${GOPACKAGES}

.PHONY: golangci-update
golangci-update:
	./tools/update_golangci.sh
