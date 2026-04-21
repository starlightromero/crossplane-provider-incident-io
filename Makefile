# Project Configuration
PROJECT_NAME := crossplane-provider-incident-io
MODULE := github.com/avodah-inc/$(PROJECT_NAME)
VERSION ?= 0.1.0

# Registry
REGISTRY := ghcr.io/avodah-inc/$(PROJECT_NAME)
IMAGE_TAG := $(REGISTRY):$(VERSION)

# Build
PROVIDER_BINARY := bin/provider
GENERATOR_BINARY := bin/generator
PLATFORMS := linux/amd64,linux/arm64

# Tools
GOLANGCI_LINT_VERSION ?= v1.57.2
KUBECONFORM_VERSION ?= v0.6.4
CRANK ?= crank

# Directories
CRD_DIR := package/crds
XPKG_DIR := _output/xpkg
XPKG_FILE := $(XPKG_DIR)/$(PROJECT_NAME)-$(VERSION).xpkg

# Go
GOFLAGS ?=
GO := go

.PHONY: all generate build docker-build docker-push xpkg-build xpkg-push \
	lint test validate-crds trivy-fs trivy-image sonar clean help

all: generate build

##@ Code Generation

generate: ## Run Upjet code generation via cmd/generator/main.go
	$(GO) run $(MODULE)/cmd/generator
	@echo "Generating managed resource method sets..."
	python3 hack/generate-managed.py
	@echo "Generating deepcopy methods..."
	controller-gen object:headerFile=hack/boilerplate.go.txt paths="./apis/..."
	@echo "Generating CRD YAML..."
	controller-gen crd:allowDangerousTypes=true paths="./apis/..." output:crd:artifacts:config=package/crds

##@ Build

build: ## Compile the provider binary from cmd/provider/main.go
	$(GO) build -o $(PROVIDER_BINARY) -ldflags "-X main.Version=$(VERSION)" $(MODULE)/cmd/provider

##@ Container Image

docker-build: ## Build multi-arch container image (linux/amd64, linux/arm64) with distroless base
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag $(IMAGE_TAG) \
		--tag $(REGISTRY):latest \
		--build-arg VERSION=$(VERSION) \
		.

docker-push: ## Push container image to ghcr.io
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag $(IMAGE_TAG) \
		--tag $(REGISTRY):latest \
		--build-arg VERSION=$(VERSION) \
		--push \
		.

##@ Crossplane Package

xpkg-build: ## Build Crossplane provider package
	@mkdir -p $(XPKG_DIR)
	$(CRANK) xpkg build \
		--package-root package \
		--output $(XPKG_FILE)

xpkg-push: ## Push xpkg to ghcr.io
	$(CRANK) xpkg push \
		--package $(XPKG_FILE) \
		$(REGISTRY):$(VERSION)

##@ Quality

lint: ## Run linters
	golangci-lint run ./...

test: ## Run tests
	$(GO) test -race -coverprofile=coverage.out ./...

##@ Validation

validate-crds: ## Validate generated CRDs with kubeconform in strict mode
	kubeconform \
		-schema-location default \
		-schema-location "https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json" \
		-strict \
		-summary \
		$(CRD_DIR)/

##@ Security

trivy-fs: ## Run Trivy filesystem scan on source code
	trivy fs . --include-dev-deps --severity CRITICAL,HIGH,MEDIUM --exit-code 1

trivy-image: ## Run Trivy container image scan
	trivy image $(IMAGE_TAG) --severity CRITICAL,HIGH,MEDIUM --exit-code 1

sonar: ## Run SonarQube scan
	sonar-scanner

##@ Cleanup

clean: ## Remove build artifacts
	rm -rf bin/ $(XPKG_DIR) coverage.out

##@ Help

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
