S6_PATH := ./data/s6-overlay/s6-rc.d
ARGS := -p $(S6_PATH)

.PHONY: build
build:
	go build -o s6-cli -v ./cmd/s6-cli

.PHONY: run
run:
	go run ./cmd/s6-cli $(ARGS)

.PHONY: lint
lint:
	go run ./cmd/s6-cli $(ARGS) lint

.PHONY: mermaid
mermaid:
	go run ./cmd/s6-cli $(ARGS) mermaid