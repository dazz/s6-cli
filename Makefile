S6_PATH := ./examples/s6-overlay/s6-rc.d
ARGS := -p $(S6_PATH)

.PHONY: build
build:
	@go build -o s6-cli -v ./cmd/s6cli

.PHONY: run
run:
	go run ./cmd/s6cli $(ARGS)

.PHONY: lint
lint:
	@go run ./cmd/s6cli $(ARGS) lint

.PHONY: mermaid
mermaid:
	@go run ./cmd/s6cli $(ARGS) mermaid

.PHONY: test
test:
	@go test -v ./...

nix:
	@nix-shell --show-trace