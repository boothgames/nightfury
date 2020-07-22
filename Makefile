.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

APP=nightfury
SRC_PACKAGES=$(shell go list ./... | grep -v "vendor")
VERSION?=1.0
BUILD?=$(shell git describe --always --dirty 2> /dev/null)
GOLINT:=$(shell command -v golint 2> /dev/null)
APP_EXECUTABLE="./out/$(APP)"
RICHGO=$(shell command -v richgo 2> /dev/null)
GOMETA_LINT=$(shell command -v golangci-lint 2> /dev/null)
GOLANGCI_LINT_VERSION=v1.29.0
SHELL=bash -o pipefail
BUILD_ARGS="-s -w -X github.com/boothgames/nightfury/cmd.version=$(VERSION)-$(BUILD)"
GO111MODULE=on

ifeq ($(GOMETA_LINT),)
	GOMETA_LINT=$(shell command -v $(PWD)/bin/golangci-lint 2> /dev/null)
endif

ifeq ($(RICHGO),)
	GO_BINARY=go
else
	GO_BINARY=richgo
endif

ifeq ($(BUILD),)
	BUILD=dev
endif

ifdef CI_COMMIT_SHORT_SHA
	BUILD=$(CI_COMMIT_SHORT_SHA)
endif

all: setup build

ci: setup-common build-common

ensure-build-dir:
	mkdir -p out

build-deps: ## Install dependencies
	go get
	go mod tidy

update-deps: ## Update dependencies
	go get -u

compile: compile-app  ## Compile nightfury

run: compile  ## run nightfury with default arguments
	./out/nightfury server

compile-app: ensure-build-dir
	$(GO_BINARY) build -ldflags $(BUILD_ARGS) -o $(APP_EXECUTABLE) ./main.go

install: ## Install nightfury
	go install -ldflags $(BUILD_ARGS)

compile-linux: ensure-build-dir ## Compile nightfury for linux
	GOOS=linux GOARCH=amd64 $(GO_BINARY) build -ldflags $(BUILD_ARGS) -o $(APP_EXECUTABLE) ./main.go

build: setup build-deps fmt build-common ## Build the application

build-common: vet lint-all test compile

compress: compile ## Compress the binary
	upx $(APP_EXECUTABLE)

fmt:
	$(GO_BINARY) fmt $(SRC_PACKAGES)

vet:
	$(GO_BINARY) vet $(SRC_PACKAGES)

setup-common:
ifeq ($(GOMETA_LINT),)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION)
endif
ifeq ($(GOLINT),)
	$(GO_BINARY) get -u golang.org/x/lint/golint
endif

setup-richgo:
ifeq ($(RICHGO),)
	$(GO_BINARY) get -u github.com/kyoh86/richgo
endif

setup: setup-richgo setup-common ensure-build-dir ## Setup environment

lint-all: lint setup-common
	$(GOMETA_LINT) run

lint:
	./scripts/lint $(SRC_PACKAGES)

test-all: test

test: ensure-build-dir ## Run tests
	ENVIRONMENT=test $(GO_BINARY) test $(SRC_PACKAGES) -race -coverprofile ./out/coverage -short -v | grep -viE "start|no test files"

test-cover-html: ## Run tests with coverage
	mkdir -p ./out
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(SRC_PACKAGES),\
	ENVIRONMENT=test $(GO_BINARY) test -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	$(GO_BINARY) tool cover -html=coverage-all.out -o out/coverage.html

dev-docker-build: ## Build nightfury server docker image with local chartmuseum repo added to it
	docker build --build-arg BUILD_MODE="dev" -f docker/nightfury/Dockerfile -t local-nightfury .
	@echo "To run:"
	@echo "$ docker run --rm -it -p 5443:5443 --name local-nightfury local-nightfury:latest"

dev-docker-run:  ## Run nightfury server in local docker container
	docker run --rm -it -p 5443:5443 --name local-nightfury local-nightfury:latest
