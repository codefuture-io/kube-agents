#  Copyright 2026 CodeFuture Authors
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

# Image URL to use all building/pushing image targets
VERSION ?= 0.1.0
IMG ?= codefuture/kube-agents

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILDDATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAG_OPTIONS = -ldflags "-X github.com/codefuture-io/kube-agents/version.Version=$(VERSION) \
                      -X github.com/codefuture-io/kube-agents/version.GitCommit=$(GIT_COMMIT) \
                      -X github.com/codefuture-io/kube-agents/version.BuildDate=$(BUILDDATE)"

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: mod-tidy
mod-tidy: ## Run go mod tidy.
	go mod tidy

##@ Build

.PHONY: build
build: fmt vet ## Build kube-agents binary.
	go build $(LDFLAG_OPTIONS) -o bin/kube-agents ./cmd/kube-agents/

.PHONY: run
run: fmt vet ## Run kube-agents from your host.
	go run ./cmd/kube-agents/

.PHONY: docker-build
docker-build: ## Build docker image.
	$(CONTAINER_TOOL) build -t ${IMG}:${VERSION} \
		--build-arg VERSION=$(VERSION) \
		--build-arg GITCOMMIT=$(GIT_COMMIT) \
		--build-arg BUILDDATE=$(BUILDDATE) .

.PHONY: docker-push
docker-push: ## Push docker image.
	$(CONTAINER_TOOL) push ${IMG}:${VERSION}

PLATFORMS ?= linux/amd64,linux/arm64

.PHONY: docker-buildx
docker-buildx: ## Build and push docker image for cross-platform support.
	- $(CONTAINER_TOOL) buildx create --name kube-agents
	$(CONTAINER_TOOL) buildx use kube-agents
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag $(IMG):$(VERSION) \
			--build-arg VERSION=$(VERSION) \
			--build-arg GITCOMMIT=$(GIT_COMMIT) \
			--build-arg BUILDDATE=$(BUILDDATE) .
	- $(CONTAINER_TOOL) buildx rm kube-agents

##@ Test

.PHONY: test
test: ## Run tests with race detection.
	go test -race ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage report.
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
