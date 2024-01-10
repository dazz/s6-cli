BINARY_NAME=s6-cli
S6_PATH := ./examples/s6-overlay/s6-rc.d
ARGS := -p $(S6_PATH)

.DEFAULT: help

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

.PHONY: dep
dep:
	@go mod download

## build: build binary file
.PHONY: build
build:
	@GOARCH=amd64 GOOS=linux go build -o /Users/dazz/bin/${BINARY_NAME} -v ./cmd/s6cli

## build: build binary file
.PHONY: build-darwin
build-darwin:
	@GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME} -v ./cmd/s6cli



## clean: clean binary file
.PHONY: clean
clean:
	@go clean
	@rm -f ${BINARY_NAME}

## run: run binary file with good args
.PHONY: run
run:
	@go run ./cmd/s6cli $(ARGS)

## test: run all tests
.PHONY: test
test:
	@go test ./...

## test-coverage: run all tests with coverage
.PHONY: test-coverage
test-coverage:
	@go test -coverprofile=coverage.out -v ./...
	@go tool cover -html=coverage.out

## nix: build binary file with nix
.PHONY: nix
nix:
	@nix-shell --show-trace

# ==================================================================================== #
# RUN COMMANDS OF CLI WITH DEFAULT ARGS
# ==================================================================================== #

## lint: lint s6-overlay directories and files
.PHONY: lint
lint:
	@go run ./cmd/s6cli $(ARGS) lint

## mermaid: generate mermaid graph
.PHONY: mermaid
mermaid:
	@go run ./cmd/s6cli $(ARGS) mermaid

.PHONY: create-oneshot
create-oneshot:
	go run ./cmd/s6cli $(ARGS) create o test

.PHONY: create-longrun
create-longrun:
	go run ./cmd/s6cli $(ARGS) create l test

.PHONY: create-bundle
create-bundle:
	go run ./cmd/s6cli $(ARGS) create b test

.PHONY: remove
remove:
	go run ./cmd/s6cli $(ARGS) remove test

.PHONY: re-create
re-create: remove create-oneshot

