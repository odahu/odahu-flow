ODAHUFLOW_OPERATOR_GENERATED_ENTITIES = ../../helms/odahu-flow-core/templates/operator/generated
LINTER_ADDITIONAL_ARGS =
KUBEBUILDER_CONTROLPLANE_START_TIMEOUT=600s

ODAHUFLOW_NAMESPACE=odahu-flow
ODAHUFLOW_SECRET_CONFIG=odahu-flow-operator-config
# Name of a sed binary
SED_BIN=sed

# Hardcoded vault token for local development
TEST_VAULT_ROOT_TOKEN=test_root_token

-include .env
.EXPORT_ALL_VARIABLES:

all: help

## help: Run unit tests
test: generate manifests
	gotestsum --junitfile operator-report.xml -- -coverprofile=operator-cover.out -covermode=count ./pkg/...

## build-operator: Build operator executable binary
build-operator: generate
	go build -o operator github.com/odahu/odahu-flow/packages/operator/cmd/operator

## build-api: Build api executable binary
build-api:
	go build -o api github.com/odahu/odahu-flow/packages/operator/cmd/api

## build-model-trainer: Build model trainer executable binary
build-model-trainer:
	go build -o trainer github.com/odahu/odahu-flow/packages/operator/cmd/trainer

## build-model-packager: Build model packager executable binary
build-model-packager:
	go build -o packager github.com/odahu/odahu-flow/packages/operator/cmd/packager

## build-model-packager: Build model packager executable binary
build-service-catalog:
	go build -o service-catalog github.com/odahu/odahu-flow/packages/operator/cmd/service_catalog

## build-all: Build all executable binyries
build-all: build-api build-model-trainer build-model-packager build-operator build-service-catalog

## run-operator: Start operator
run-operator:
	go run ./cmd/operator/main.go

## run-api: Start api server
run-api:
	go run ./cmd/operator/main.go

## install-crds: Install CRDs into a cluster
install-crds: manifests
	kubectl apply -f config/crds

## manifests: Generate manifests e.g. CRD, RBAC etc.
manifests:  swag
	rm -fv ./config/crds/*.yaml
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

	mkdir -p ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}
	@echo "removing old manifests from HELM's generated dictionary"

	rm -fv ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/*.yaml
	@echo "copying actual manifests to HELM's generated dictionary"
	cp -f ./config/crds/*.yaml ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}
	cp config/rbac/rbac_role.yaml ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml

	# All Odahuflow CRD's and operator role are generated resources.
	# We cannot add, for example, a custom label during the generation process.
	# That's why we add them using sed tool after generation.

	"${SED_BIN}" -i '/^  name:.*$$/a \ \ labels:' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml

	"${SED_BIN}" -i '/^  labels:$$/a \ \ \ \ {{- include "odahuflow.helm-labels" (dict "component" "operator" "root" .) | nindent 4 }}' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/*.yaml

	"${SED_BIN}" -i 's/manager-role/"{{ .Release.Name }}-operator"/g' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml
	"${SED_BIN}" -i '1i{{- if .Values.rbac }}' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml
	"${SED_BIN}" -i '$$a{{- end }}' ${ODAHUFLOW_OPERATOR_GENERATED_ENTITIES}/operator-rbac-role.yaml

## run-vault: Start vault in dev model for local development
run-value:
	vault server -dev -dev-root-token-id $TEST_VAULT_ROOT_TOKEN

## swag: Generate swagger files
swag:
	swag init -g pkg/webserver/server.go
	# Remove the line with generated timestamp
	"${SED_BIN}" -i '3d' docs/docs.go
	"${SED_BIN}" -i 's/connection\.//g' docs/*
	"${SED_BIN}" -i 's/routes\.//g' docs/*
	"${SED_BIN}" -i 's/deployment\.//g' docs/*
	"${SED_BIN}" -i 's/training\.//g' docs/*
	"${SED_BIN}" -i 's/packaging\.//g' docs/*
	"${SED_BIN}" -i 's/config\.//g' docs/*
	"${SED_BIN}" -i 's/user\.//g' docs/*
	"${SED_BIN}" -i 's/v1alpha1\.//g' docs/*

## apply-crds: Apply all odahuflow crds
apply-crds:
	kubectl apply -f config/crds/

## delete-crds: Delete all odahuflow crds
delete-crds:
	kubectl delete -f config/crds/

## recreate-crds: Recreate all odahuflow crds
recreate-crds: delete-crds apply-crds

## generate: Generate golang codes
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

## generate-all: Execute all code generators
generate-all: generate manifests swag

## export-config-from-k8s: Export operator config from k8s secret. It helps to run operator or api locally with properly set of config options
export-config-from-k8s:
	kubectl get -n ${ODAHUFLOW_NAMESPACE} secrets ${ODAHUFLOW_SECRET_CONFIG} -o=jsonpath="{$$.data['config\\.yaml']}" | base64 -d > config.yaml

## turn-off-odahu-flow-operator: The command scales down the odahuflow operator to the zero number of pods.
turn-off-odahu-flow-operator:
	kubectl scale --replicas 0 -n ${ODAHUFLOW_NAMESPACE} statefulset/odahu-flow-operator

## turn-on-odahu-flow-operator: The command scales up the odahuflow operator to the one number of pods.
turn-on-odahu-flow-operator:
	kubectl scale --replicas 1 -n ${ODAHUFLOW_NAMESPACE} statefulset/odahu-flow-operator

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
