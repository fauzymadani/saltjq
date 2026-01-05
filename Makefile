# Makefile for saltjq (minimal targets)

BINARY := saltjq
BUILD_CMD := go build -o $(BINARY) ./cmd/saltjq

.PHONY: all build fmt test vet run install clean help

all: build

# build the CLI binary
build:
	$(BUILD_CMD)

# format Go sources
fmt:
	gofmt -w .

# run unit tests
test:
	go test ./...

# static checks
vet:
	go vet ./...

# build and run local binary
run: build
	./$(BINARY)

# install to $GOBIN
install:
	go install ./cmd/saltjq

# clean the built binary
clean:
	-rm -f $(BINARY)

# show help
help:
	@echo "Makefile targets:"
	@echo "  make build   - build the saltjq binary"
	@echo "  make fmt     - format code with gofmt"
	@echo "  make test    - run go test ./..."
	@echo "  make vet     - run go vet ./..."
	@echo "  make run     - build and run the binary"
	@echo "  make install - install binary to GOBIN"
	@echo "  make clean   - remove built binary"

