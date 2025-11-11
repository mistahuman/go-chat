BINARY := bin/go-chat
PKG := ./...

.PHONY: build run clean fmt lint test tidy

build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(dir $(BINARY))
	@go build -o $(BINARY) ./

run: build
	@$(BINARY)

clean:
	@echo "Removing build artifacts"
	@rm -rf bin

fmt:
	@go fmt $(PKG)

lint:
	@go vet $(PKG)

test:
	@go test $(PKG)

tidy:
	@go mod tidy
