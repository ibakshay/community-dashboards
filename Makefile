TOOLS_BIN_DIR ?= $(shell pwd)/tmp/bin

export PATH := $(TOOLS_BIN_DIR):$(PATH)

GOLANGCILINTER_BINARY=$(TOOLS_BIN_DIR)/golangci-lint
MDOX_BINARY=$(TOOLS_BIN_DIR)/mdox
MDOX_VALIDATE_CONFIG?=.mdox.validate.yaml

TOOLING=$(MDOX_BINARY) $(GOLANGCILINTER_BINARY)

MD_FILES_TO_FORMAT=$(shell ls *.md)

GOCMD=go
GOMAIN=main.go
GOBUILD=$(GOCMD) build
GOOS?=$(shell go env GOOS)
ENVVARS=GOOS=$(GOOS) CGO_ENABLED=0

# Dashboard build configuration with defaults
OUTPUT_DIR_OPERATOR ?= ./dist/dashboards/operator
OUTPUT_DIR_PERSES ?= ./dist/dashboards/perses
OUTPUT_FORMAT_PERSES ?= json
PROJECT ?= default
DATASOURCE ?= prometheus-datasource

.PHONY: demo
start-demo:
	@echo "Setting up demo environment"

	@cd ./examples && docker-compose up -d

.PHONY: clean-demo
clean-demo:
	@echo "Cleaning up demo environment"

	@cd ./examples && docker-compose down -v

.PHONY: build-dashboards
build-dashboards:
	@echo "Building dashboards"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir="./examples/dashboards/operator" --output="operator" --project="perses-dev" --datasource="prometheus-datasource"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir="./examples/dashboards/perses" --output="yaml" --project="perses-dev" --datasource="prometheus-datasource"

# Adding a new target for building and testing dashboards locally with configurable flags
.PHONY: build-dashboards-local
build-dashboards-local:
	@echo "Building dashboards for local testing"
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir=$(OUTPUT_DIR_OPERATOR) --output="operator" --project=$(PROJECT) --datasource=$(DATASOURCE)
	@$(ENVVARS) $(GOCMD) run $(GOMAIN) --output-dir=$(OUTPUT_DIR_PERSES) --output=$(OUTPUT_FORMAT_PERSES) --project=$(PROJECT) --datasource=$(DATASOURCE)

.PHONY: deps
deps:
	$(ENVVARS) $(GOCMD) mod download

.PHONY: fmt
fmt:
	$(ENVVARS) $(GOCMD) fmt -x ./...

.PHONY: vet
vet:
	$(ENVVARS) $(GOCMD) vet ./...

.PHONY: check-golang
check-golang: $(GOLANGCILINTER_BINARY)
	$(GOLANGCILINTER_BINARY) run

.PHONY: fix-golang
fix-golang: $(GOLANGCILINTER_BINARY)
	$(GOLANGCILINTER_BINARY) run --fix

.PHONY: docs
docs: $(MDOX_BINARY)
	@echo ">> formatting and local/remote link check"
	$(MDOX_BINARY) fmt --soft-wraps -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) $(MD_FILES_TO_FORMAT)

.PHONY: check-docs
check-docs: $(MDOX_BINARY)
	@echo ">> checking formatting and local/remote links"
	$(MDOX_BINARY) fmt --soft-wraps --check -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) $(MD_FILES_TO_FORMAT)

.PHONY: tidy
tidy:
	go mod tidy -v
	cd scripts && go mod tidy -v -modfile=go.mod -compat=1.18

all: fmt vet deps check-golang check-docs

$(TOOLS_BIN_DIR):
	mkdir -p $(TOOLS_BIN_DIR)

$(TOOLING): $(TOOLS_BIN_DIR)
	@echo Installing tools from scripts/tools.go
	@cat scripts/tools.go | grep _ | awk -F'"' '{print $$2}' | GOBIN=$(TOOLS_BIN_DIR) xargs -tI % go install -mod=readonly -modfile=scripts/go.mod %
