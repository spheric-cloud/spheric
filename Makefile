# Image URL to use all building/pushing image targets
CONTROLLER_IMG ?= controller:latest
APISERVER_IMG ?= apiserver:latest
SPHERELET_IMG ?= spherelet:latest

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.30.0

# Docker image name for the mkdocs based local development setup
IMAGE=spheric/documentation

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	# controller-manager
	$(CONTROLLER_GEN) rbac:roleName=manager-role webhook paths="./internal/controllers/...;./api/..." output:rbac:artifacts:config=config/controller/rbac

.PHONY: generate
generate: models-schema deepcopy-gen client-gen lister-gen informer-gen defaulter-gen conversion-gen openapi-gen applyconfiguration-gen
	MODELS_SCHEMA=$(MODELS_SCHEMA) \
	DEEPCOPY_GEN=$(DEEPCOPY_GEN) \
	CLIENT_GEN=$(CLIENT_GEN) \
	LISTER_GEN=$(LISTER_GEN) \
	INFORMER_GEN=$(INFORMER_GEN) \
	DEFAULTER_GEN=$(DEFAULTER_GEN) \
	CONVERSION_GEN=$(CONVERSION_GEN) \
	OPENAPI_GEN=$(OPENAPI_GEN) \
	APPLYCONFIGURATION_GEN=$(APPLYCONFIGURATION_GEN) \
	./hack/update-codegen.sh

.PHONY: proto
proto: goimports protoc-gen-go protoc-gen-go-grpc
	PROTOC_GEN_GO=$(PROTOC_GEN_GO) \
	PROTOC_GEN_GO_GRPC=$(PROTOC_GEN_GO_GRPC) \
	./hack/update-proto.sh
	$(GOIMPORTS) -w ./iri-api

.PHONY: fmt
fmt: goimports ## Run goimports against code.
	$(GOIMPORTS) -w .

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: golangci-lint ## Run golangci-lint on the code.
	$(GOLANGCI_LINT) run ./...

.PHONY: clean
clean: ## Clean any artifacts that can be regenerated.
	rm -rf client-go/applyconfigurations
	rm -rf client-go/informers
	rm -rf client-go/listers
	rm -rf client-go/spheric
	rm -rf client-go/openapi

.PHONY: add-license
add-license: addlicense ## Add license headers to all go files.
	find . -name '*.go' -exec $(ADDLICENSE) -f hack/license-header.txt {} +

.PHONY: check-license
check-license: addlicense ## Check that every file has a license header present.
	find . -name '*.go' -exec $(ADDLICENSE) -check -c 'Spheric authors' {} +

.PHONY: check
check: generate manifests add-license fmt lint test # Generate manifests, code, lint, add licenses, test

.PHONY: docs
docs: gen-crd-api-reference-docs ## Run go generate to generate API reference documentation.
	$(GEN_CRD_API_REFERENCE_DOCS) -api-dir ./api/core/v1alpha1 -config ./hack/api-reference/config.json -template-dir ./hack/api-reference/template -out-file ./docs/api-reference/core.md

.PHONY: start-docs
start-docs: ## Start the local mkdocs based development environment.
	docker build -t $(IMAGE) -f docs/Dockerfile . --load
	docker run -p 8000:8000 -v `pwd`/:/docs $(IMAGE)

.PHONY: clean-docs
clean-docs: ## Remove all local mkdocs Docker images (cleanup).
	docker container prune --force --filter "label=project=spheric_documentation"

.PHONY: test
test: manifests generate fmt vet test-only ## Run tests.

.PHONY: test-only
test-only: envtest ## Run *only* the tests - no generation, linting etc.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -coverprofile cover.out

.PHONY: extract-openapi
extract-openapi: envtest openapi-extractor
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" $(OPENAPI_EXTRACTOR) \
		--apiserver-package="spheric.cloud/spheric/cmd/apiserver" \
		--apiserver-build-opts=mod \
		--apiservices="./config/apiserver/apiservice/bases" \
		--attach-control-plane-output \
		--output="./gen"

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	go build -o bin/manager ./cmd/spheric-controller-manager

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./cmd/spheric-controller-manager

.PHONY: docker-build
docker-build: \
	docker-build-spheric-apiserver docker-build-spheric-controller-manager \
	docker-build-spherelet ## Build docker image with the manager.

.PHONY: docker-build-spheric-apiserver
docker-build-spheric-apiserver: ## Build apiserver.
	docker build --target apiserver -t ${APISERVER_IMG} .

.PHONY: docker-build-spheric-controller-manager
docker-build-spheric-controller-manager: ## Build controller-manager.
	docker build --target controller-manager -t ${CONTROLLER_IMG} .

.PHONY: docker-build-spherelet
docker-build-spherelet: ## Build spherelet image.
	docker build --target spherelet -t ${SPHERELET_IMG} .

.PHONY: docker-build-irictl-machine
docker-build-irictl-machine: ## Build irictl-machine image.
	docker build --target irictl-machine -t ${IRICTL_MACHINE_IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${CONTROLLER_IMG}
	docker push ${APISERVER_IMG}

##@ Deployment

.PHONY: install
install: manifests kustomize ## Install API server & API services into the K8s cluster specified in ~/.kube/config. This requires APISERVER_IMG to be available for the cluster.
	cd config/apiserver/server && $(KUSTOMIZE) edit set image apiserver=${APISERVER_IMG}
	kubectl apply -k config/apiserver/default

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall API server & API services from the K8s cluster specified in ~/.kube/config.
	kubectl delete -k config/apiserver/default

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/controller/manager && $(KUSTOMIZE) edit set image controller=${CONTROLLER_IMG}
	kubectl apply -k config/controller/default

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	kubectl delete -k config/controller/default

##@ Kind Deployment plumbing

.PHONY: kind-build-apiserver
kind-build-apiserver: ## Build the apiserver for usage in kind.
	docker build --target apiserver -t apiserver .

.PHONY: kind-build-controller
kind-build-controller: ## Build the controller for usage in kind.
	docker build --target manager -t controller .

.PHONY: kind-build
kind-build: kind-build-apiserver kind-build-controller ## Build the apiserver and controller for usage in kind.

.PHONY: kind-load-apiserver
kind-load-apiserver: ## Load the apiserver image into the kind cluster.
	kind load docker-image apiserver

.PHONY: kind-load-controller
kind-load-controller: ## Load the controller image into the kind cluster.
	kind load docker-image controller

.PHONY: kind-load-spherelet
kind-load-spherelet:
	kind load docker-image ${SPHERELET_IMG}

.PHONY: kind-load
kind-load: kind-load-apiserver kind-load-controller ## Load the apiserver and controller in kind.

.PHONY: kind-restart-apiserver
kind-restart-apiserver: ## Restart the apiserver in kind. Useless if the manifests are not in place (deployed e.g. via kind-apply / kind-deploy).
	kubectl -n spheric-system delete rs -l control-plane=apiserver

.PHONY: kind-restart-controller
kind-restart-controller: ## Restart the controller in kind. Useless if the manifests are not in place (deployed e.g. via kind-apply / kind-deploy).
	kubectl -n spheric-system delete rs -l control-plane=controller-manager

.PHONY: kind-restart
kind-restart: kind-restart-apiserver kind-restart-controller ## Restart the apiserver and controller in kind. Restart is useless if the manifests are not in place (deployed e.g. via kind-apply / kind-deploy).

.PHONY: kind-build-load-restart-controller
kind-build-load-restart-controller: kind-build-controller kind-load-controller kind-restart-controller ## Build, load and restart the controller in kind. Restart is useless if the manifests are not in place (deployed e.g. via kind-apply / kind-deploy).

.PHONY: kind-build-load-restart-apiserver
kind-build-load-restart-apiserver: kind-build-apiserver kind-load-apiserver kind-restart-apiserver ## Build, load and restart the apiserver in kind. Restart is useless if the manifests are not in place (deployed e.g. via kind-apply / kind-deploy).

.PHONY: kind-build-load-restart
kind-build-load-restart: kind-build-load-restart-apiserver kind-build-load-restart-controller ## Build load and restart the apiserver and controller in kind. Restart is useless if the manifests are not in place (deployed e.g. via kind-apply / kind-deploy).

.PHONY: kind-apply-apiserver
kind-apply-apiserver: manifests kustomize ## Applies the apiserver manifests in kind. Caution, without loading the images, the pods won't come up. Use kind-install / kind-deploy for a deployment including loading the images.
	kubectl apply -k config/apiserver/kind

.PHONY: kind-install
kind-install: kind-build-load-restart-apiserver kind-apply-apiserver ## Build and load and apply apiserver in kind. Restarts apiserver if it was present.

.PHONY: kind-uninstall
kind-uninstall: manifests kustomize ## Uninstall API server & API services from the K8s cluster specified in ~/.kube/config.
	kubectl delete -k config/apiserver/kind

.PHONY: kind-apply
kind-apply: ## Apply the config in kind. Caution: Without loading the images, the pods won't come up. Use kind-deploy for a deployment including loading the images.
	kubectl apply -k config/kind

.PHONY: kind-delete
kind-delete: ## Delete the config from kind.
	kubectl delete -k config/kind

.PHONY: kind-deploy
kind-deploy: kind-build-load-restart kind-apply ## Build and load apiserver and controller into the kind cluster, then apply the config. Restarts apiserver / controller if they were present.

##@ Tools

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
OPENAPI_EXTRACTOR ?= $(LOCALBIN)/openapi-extractor
DEEPCOPY_GEN ?= $(LOCALBIN)/deepcopy-gen
CLIENT_GEN ?= $(LOCALBIN)/client-gen
LISTER_GEN ?= $(LOCALBIN)/lister-gen
INFORMER_GEN ?= $(LOCALBIN)/informer-gen
DEFAULTER_GEN ?= $(LOCALBIN)/defaulter-gen
CONVERSION_GEN ?= $(LOCALBIN)/conversion-gen
OPENAPI_GEN ?= $(LOCALBIN)/openapi-gen
APPLYCONFIGURATION_GEN ?= $(LOCALBIN)/applyconfiguration-gen
GEN_CRD_API_REFERENCE_DOCS ?= $(LOCALBIN)/gen-crd-api-reference-docs
ADDLICENSE ?= $(LOCALBIN)/addlicense
PROTOC_GEN_GO ?= $(LOCALBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC ?= $(LOCALBIN)/protoc-gen-go-grpc
MODELS_SCHEMA ?= $(LOCALBIN)/models-schema
GOIMPORTS ?= $(LOCALBIN)/goimports
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

## Tool Versions
KUSTOMIZE_VERSION ?= v5.4.3
CODE_GENERATOR_VERSION ?= v0.31.0
CONTROLLER_TOOLS_VERSION ?= v0.16.1
OPENAPI_GEN_VERSION ?= f7e401e7b4c2199f15e2cf9e37a2faa2209f286a # Unfortunately, no tagged releases - watch https://github.com/kubernetes/kube-openapi/issues/383 for changes.
GEN_CRD_API_REFERENCE_DOCS_VERSION ?= v0.3.0
ADDLICENSE_VERSION ?= v1.1.1
PROTOC_GEN_GO_VERSION ?= v1.34.2
PROTOC_GEN_GO_GRPC_VERSION ?= v1.5.1
GOIMPORTS_VERSION ?= v0.24.0
GOLANGCI_LINT_VERSION ?= v2.1.6

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: deepcopy-gen
deepcopy-gen: $(DEEPCOPY_GEN) ## Download deepcopy-gen locally if necessary.
$(DEEPCOPY_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/deepcopy-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GENERATOR_VERSION)

.PHONY: client-gen
client-gen: $(CLIENT_GEN) ## Download client-gen locally if necessary.
$(CLIENT_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/client-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/client-gen@$(CODE_GENERATOR_VERSION)

.PHONY: lister-gen
lister-gen: $(LISTER_GEN) ## Download lister-gen locally if necessary.
$(LISTER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/lister-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/lister-gen@$(CODE_GENERATOR_VERSION)

.PHONY: informer-gen
informer-gen: $(INFORMER_GEN) ## Download informer-gen locally if necessary.
$(INFORMER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/informer-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/informer-gen@$(CODE_GENERATOR_VERSION)

.PHONY: defaulter-gen
defaulter-gen: $(DEFAULTER_GEN) ## Download defaulter-gen locally if necessary.
$(DEFAULTER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/defaulter-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/defaulter-gen@$(CODE_GENERATOR_VERSION)

.PHONY: conversion-gen
conversion-gen: $(CONVERSION_GEN) ## Download conversion-gen locally if necessary.
$(CONVERSION_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/conversion-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/conversion-gen@$(CODE_GENERATOR_VERSION)

.PHONY: openapi-gen
openapi-gen: $(OPENAPI_GEN) ## Download openapi-gen locally if necessary.
$(OPENAPI_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/openapi-gen || GOBIN=$(LOCALBIN) go install k8s.io/kube-openapi/cmd/openapi-gen@$(OPENAPI_GEN_VERSION)

.PHONY: applyconfiguration-gen
applyconfiguration-gen: $(APPLYCONFIGURATION_GEN) ## Download applyconfiguration-gen locally if necessary.
$(APPLYCONFIGURATION_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/applyconfiguration-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/applyconfiguration-gen@$(CODE_GENERATOR_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: openapi-extractor
openapi-extractor: $(OPENAPI_EXTRACTOR) ## Download openapi-extractor locally if necessary.
$(OPENAPI_EXTRACTOR): $(LOCALBIN)
	test -s $(LOCALBIN)/openapi-extractor || GOBIN=$(LOCALBIN) go install github.com/ironcore-dev/openapi-extractor/cmd/openapi-extractor@latest

.PHONY: gen-crd-api-reference-docs
gen-crd-api-reference-docs: $(GEN_CRD_API_REFERENCE_DOCS) ## Download gen-crd-api-reference-docs locally if necessary.
$(GEN_CRD_API_REFERENCE_DOCS): $(LOCALBIN)
	test -s $(LOCALBIN)/gen-crd-api-reference-docs || GOBIN=$(LOCALBIN) go install github.com/ahmetb/gen-crd-api-reference-docs@$(GEN_CRD_API_REFERENCE_DOCS_VERSION)

.PHONY: addlicense
addlicense: $(ADDLICENSE) ## Download addlicense locally if necessary.
$(ADDLICENSE): $(LOCALBIN)
	test -s $(LOCALBIN)/addlicense || GOBIN=$(LOCALBIN) go install github.com/google/addlicense@$(ADDLICENSE_VERSION)

.PHONY: protoc-gen-go
protoc-gen-go: $(PROTOC_GEN_GO) ## Download protoc-gen-go locally if necessary.
$(PROTOC_GEN_GO): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-go || GOBIN=$(LOCALBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc: $(PROTOC_GEN_GO_GRPC) ## Download protoc-gen-go-grpc locally if necessary.
$(PROTOC_GEN_GO_GRPC): $(LOCALBIN)
	test -s $(LOCALBIN)/protoc-gen-go-grpc || GOBIN=$(LOCALBIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)

.PHONY: models-schema
models-schema: $(MODELS_SCHEMA) ## Install models-schema locally if necessary.
$(MODELS_SCHEMA): $(LOCALBIN)
	test -s $(LOCALBIN)/models-schema || GOBIN=$(LOCALBIN) go install spheric.cloud/spheric/models-schema

.PHONY: goimports
goimports: $(GOIMPORTS) ## Download goimports locally if necessary.
$(GOIMPORTS): $(LOCALBIN)
	test -s $(LOCALBIN)/goimports || GOBIN=$(LOCALBIN) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	test -s $(LOCALBIN)/golangci-lint || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
