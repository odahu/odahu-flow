# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

ODAHUFLOW_OPERATOR_GENERATED_ENTITIES = ../../helms/odahu-flow-core/templates/operator/generated
LINTER_ADDITIONAL_ARGS =
KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=600s

ODAHUFLOW_NAMESPACE=odahu-flow
ODAHUFLOW_SECRET_CONFIG=odahu-flow-operator-config
# Name of a sed binary
SED_BIN=sed

GO_BIN_DATA_EXTRA_ARGS=


# Hardcoded vault token for local development
TEST_VAULT_ROOT_TOKEN=test_root_token

-include .env
.EXPORT_ALL_VARIABLES:

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: help

## help: Run unit tests
test: generate manifests
	gotestsum --junitfile operator-report.xml -- -coverprofile=coverage.txt ./...

## build-operator: Build operator executable binary
build-operator: generate
	go build -o operator ./main.go

## build-webserver: Build api executable binary
build-webserver:
	go build -o webserver ./cmd/api

## build-controller: Build odahu system controller
build-controller:
	go build -o controller ./cmd/controller

## build-model-trainer: Build model trainer executable binary
build-model-trainer:
	go build -o trainer ./cmd/trainer

## build-model-packager: Build model packager executable binary
build-model-packager:
	go build -o packager ./cmd/packager

## build-model-packager: Build model packager executable binary
build-service-catalog:
	go build -o service-catalog ./cmd/service_catalog

## build-all: Build all executable binyries
build-all: build-webserver build-controller build-model-trainer build-model-packager build-operator build-service-catalog

## run-operator: Start operator
run-operator:
	go run ./cmd/operator/main.go

## run-api: Start api server
run-api:
	go run ./cmd/operator/main.go

## Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."


## Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	rm -fv ./config/crds/*.yaml

	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crds

	mkdir -p ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}
	@echo "removing old manifests from HELM's generated dictionary"

	rm -fv ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/*.yaml
	@echo "copying actual manifests to HELM's generated dictionary"
	cp -f ./config/crds/*.yaml ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}
	cp config/rbac/role.yaml ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml

	# All Odahuflow CRD's and operator role are generated resources.
	# We cannot add, for example, a custom label during the generation process.
	# That's why we add them using sed tool after generation.

	"${SED_BIN}" -i '/^  name:.*$$/a \ \ labels:' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml

	"${SED_BIN}" -i '/^  labels:$$/a \ \ \ \ {{- include "odahuflow.helm-labels" (dict "component" "operator" "root" .) | nindent 4 }}' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/*.yaml

	"${SED_BIN}" -i 's/manager-role/"{{ .Release.Name }}-operator"/g' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml
	"${SED_BIN}" -i '1i{{- if .Values.rbac }}' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml
	"${SED_BIN}" -i '$$a{{- end }}' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml

## swag: Generate swagger files
swag:
	swag init -g pkg/apiserver/server.go --exclude pkg/servicecatalog

	# Remove the line with generated timestamp
	"${SED_BIN}" -i '3d' docs/docs.go
	"${SED_BIN}" -i 's/connection\.//g' docs/*
	"${SED_BIN}" -i 's/routes\.//g' docs/*
	"${SED_BIN}" -i 's/deployment\.//g' docs/*
	"${SED_BIN}" -i 's/training\.//g' docs/*
	"${SED_BIN}" -i 's/packaging\.//g' docs/*
	"${SED_BIN}" -i 's/config\.//g' docs/*
	"${SED_BIN}" -i 's/user\.//g' docs/*
	"${SED_BIN}" -i 's/httputil\.//g' docs/*
	"${SED_BIN}" -i 's/event\.//g' docs/*
	"${SED_BIN}" -i 's/v1alpha1\.//g' docs/*
	"${SED_BIN}" -i 's/v1alpha1\.//g' docs/*

	"${SED_BIN}" -i 's/model\.//g' pkg/servicecatalog/docs/*
	"${SED_BIN}" -i 's/httputil\.//g' pkg/servicecatalog/docs/*


## generate-all: Execute all code generators
generate-all: generate manifests swag


## generate-resources: Generate go-bindata
generate-resources:
	go-bindata -pkg postgres -prefix pkg/database/migrations/postgres/sources -o pkg/database/migrations/postgres/bindata.go ${GO_BIN_DATA_EXTRA_ARGS} pkg/database/migrations/postgres/sources
	git add pkg/database/migrations/postgres/bindata.go
	go-bindata -ignore .*_test.rego -pkg bindata -prefix pkg/deployment/assets -o pkg/deployment/bindata/bindata.go ${GO_BIN_DATA_EXTRA_ARGS} pkg/deployment/assets
	git add pkg/deployment/bindata/bindata.go


apply-crds:
	kubectl apply -f config/crds/

## delete-crds: Delete all odahuflow crds
delete-crds:
	kubectl delete -f config/crds/

## recreate-crds: Recreate all odahuflow crds
recreate-crds: delete-crds apply-crds

## lint: Start golangci-lint linter
lint:
	golangci-lint run ${LINTER_ADDITIONAL_ARGS}

## static-content: Generate static content for API server. Read more about installation and usage: https://github.com/rakyll/statik
static-content:
	rm -rf pkg/static
	statik -src static -dest pkg -p static
	git add pkg/static

## help: Show the help message
help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo


controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
