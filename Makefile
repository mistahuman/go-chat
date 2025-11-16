GO ?= go
DOCKER ?= docker
DOCKER_COMPOSE ?= docker compose
BINARY := bin/chatserver
PKG := ./...
IMAGE ?= go-chat:local
DOC_PORT ?= 6060

.PHONY: build run clean fmt lint test tidy doc docker-build docker-run compose-up compose-down compose-logs

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
	@echo "Starting documentation server on http://localhost:$(DOC_PORT)"
	@$(GO) run golang.org/x/tools/cmd/godoc@latest -http=:$(DOC_PORT)

docker-build:
	@$(DOCKER) build -t $(IMAGE) .

docker-run: docker-build
	@$(DOCKER) run --rm -p 8080:8080 $(IMAGE)

compose-up:
	@$(DOCKER_COMPOSE) up --build

compose-down:
	@$(DOCKER_COMPOSE) down

compose-logs:
	@$(DOCKER_COMPOSE) logs -f
