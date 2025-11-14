GO ?= go
BINARY := bin/chatserver
PKG := ./...

.PHONY: build run clean fmt lint test tidy doc

build: $(BINARY)

$(BINARY):
	@echo "Building $@"
	@mkdir -p $(dir $@)
	@$(GO) build -o $@ ./cmd/chatserver

run: build
	@$(BINARY)

clean:
	@echo "Removing build artifacts"
	@rm -rf $(dir $(BINARY))

fmt:
	@$(GO) fmt $(PKG)

lint:
	@$(GO) vet $(PKG)

test:
	@$(GO) test $(PKG)

tidy:
	@$(GO) mod tidy

doc:
	@$(GO) doc ./...
