SHELL := /bin/bash

PROJECTNAME := $(shell basename "$(PWD)")
PYLINT_FOLDER=target/pylint
PYDOCSTYLE_FOLDER=target/pydocstyle
PROJECTS_PYLINT=sdk cli tests
PROJECTS_PYCODESTYLE="sdk cli"
BUILD_PARAMS=
ODAHUFLOW_VERSION=0.11.0
CREDENTIAL_SECRETS=.secrets.yaml
SANDBOX_PYTHON_TOOLCHAIN_IMAGE=
ROBOT_FILES=**/*.robot
ROBOT_THREADS=6
ROBOT_OPTIONS=-e disable
E2E_PYTHON_TAGS=
COMMIT_ID=
TEMP_DIRECTORY=
BUILD_TAG=latest
TAG=
# Example of DOCKER_REGISTRY: nexus.domain.com:443/
DOCKER_REGISTRY=
HELM_ADDITIONAL_PARAMS=
# Specify gcp auth keys
GOOGLE_APPLICATION_CREDENTIALS=
MOCKS_DIR=target/mocks
SWAGGER_FILE=packages/operator/docs/swagger.yaml
PYTHON_MODEL_DIR=packages/sdk/odahuflow/sdk/models
SWAGGER_CODEGEN_BIN=java -jar swagger-codegen-cli.jar
PIP_BIN=pip3

HIERA_KEYS_DIR=
ODAHUFLOW_PROFILES_DIR=

EXPORT_HIERA_DOCKER_IMAGE := odahuflow/terraform:${BUILD_TAG}
SECRET_DIR := $(CURDIR)/.secrets
CLUSTER_PROFILE := ${SECRET_DIR}/cluster_profile.json

-include .env

.EXPORT_ALL_VARIABLES:

.PHONY: install-all install-cli install-sdk

all: help

check-tag:
	@if [ "${TAG}" == "" ]; then \
	    echo "TAG is not defined, please define the TAG variable" ; exit 1 ;\
	fi
	@if [ "${DOCKER_REGISTRY}" == "" ]; then \
	    echo "DOCKER_REGISTRY is not defined, please define the DOCKER_REGISTRY variable" ; exit 1 ;\
	fi

## install-all: Install all python packages
install-all: install-sdk install-cli install-robot

## install-sdk: Install sdk python package
install-sdk:
	cd packages/sdk && \
		"${PIP_BIN}" install ${BUILD_PARAMS} -e .

## install-cli: Install cli python package
install-cli:
	cd packages/cli && \
		"${PIP_BIN}" install ${BUILD_PARAMS} -e .

## install-robot: Install robot tests
install-robot:
	cd packages/robot && \
		"${PIP_BIN}" install ${BUILD_PARAMS} -e .

## docker-build-odahu-flow-cli: Build image with odahuflow cli
docker-build-odahu-flow-cli:
	docker build -t odahu/odahu-flow-cli:${BUILD_TAG} -f containers/odahu-flow-cli/Dockerfile .

## docker-build-robot-tests: Build pipeline agent docker image
docker-build-robot-tests:
	docker build -t odahu/odahu-flow-robot-tests:${BUILD_TAG} -f containers/robot-tests/Dockerfile .

## docker-build-api: Build api docker image
docker-build-api:
	docker build --target api -t odahu/odahu-flow-api:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-controller: Build controller docker image
docker-build-controller:
	docker build --target controller -t odahu/odahu-flow-controller:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-model-trainer: Build model builder docker image
docker-build-model-trainer:
	docker build --target model-trainer -t odahu/odahu-flow-model-trainer:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-model-packager: Build model packager docker image
docker-build-model-packager:
	docker build --target model-packager -t odahu/odahu-flow-model-packager:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-operator: Build operator docker image
docker-build-operator:
	docker build --target operator -t odahu/odahu-flow-operator:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-service-catalog: Build service catalog docker image
docker-build-service-catalog:
	docker build --target service-catalog -t odahu/odahu-flow-service-catalog:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-feedback-collector: Build feedback collector image
docker-build-feedback-collector:
	docker build --target collector -t odahu/odahu-flow-feedback-collector:${BUILD_TAG} -f containers/feedback/Dockerfile .

## docker-build-feedback-rq-catcher: Build feedback rq-catcher image
docker-build-feedback-rq-catcher:
	docker build --target rq-catcher -t odahu/odahu-flow-feedback-rq-catcher:${BUILD_TAG} -f containers/feedback/Dockerfile .

## docker-build-all: Build all docker images
docker-build-all:  docker-build-robot-tests docker-build-api docker-build-controller  docker-build-model-trainer  docker-build-operator  docker-build-feedback-collector docker-build-service-catalog

## docker-push-robot-tests: Push pipeline agent docker image
docker-push-robot-tests:
	docker tag odahu/odahu-flow-robot-tests:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-robot-tests:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahu-flow-robot-tests:${TAG}

## docker-push-api: Push api docker image
docker-push-api:  check-tag
	docker tag odahu/odahu-flow-api:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-api:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahu-flow-api:${TAG}

## docker-push-api: Push controller docker image
docker-push-controller:  check-tag
	docker tag odahu/odahu-flow-controller:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-controller:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahu-flow-controller:${TAG}

## docker-push-model-packager: Push model packager docker image
docker-push-model-packager:  check-tag
	docker tag odahu/odahu-flow-model-packager:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-model-packager:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahu-flow-model-packager:${TAG}

## docker-push-model-trainer: Push model trainer docker image
docker-push-model-trainer:  check-tag
	docker tag odahu/odahu-flow-model-trainer:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-model-trainer:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahu-flow-model-trainer:${TAG}

## docker-push-operator: Push operator docker image
docker-push-operator:  check-tag
	docker tag odahu/odahu-flow-operator:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-operator:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahu-flow-operator:${TAG}

## docker-push-service-catalog: Push service catalog docker image
docker-push-service-catalog:  check-tag
	docker tag odahu/odahu-flow-service-catalog:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahu-flow-service-catalog:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/service-catalog:${TAG}

## docker-push-feedback-collector: Push feedback collector docker image
docker-push-feedback-collector:  check-tag
	docker tag odahu/odahu-flow-feedback-collector:${BUILD_TAG} ${DOCKER_REGISTRY}odahu/odahu-flow-feedback-collector:${TAG}
	docker push ${DOCKER_REGISTRY}odahu/odahu-flow-feedback-collector:${TAG}

## docker-push-all: Push all docker images
docker-push-all:  docker-push-robot-tests docker-push-api  docker-push-model-trainer  docker-push-operator  docker-push-feedback-collector docker-push-service-catalog

## helm-install: Install the odahuflow helm chart from source code
helm-install: helm-delete
	helm install helms/odahuflow --atomic --wait --timeout 320 --namespace odahuflow --name odahuflow --debug ${HELM_ADDITIONAL_PARAMS}

## helm-delete: Delete the odahuflow helm release
helm-delete:
	helm delete --purge odahuflow || true

## install-python-linter: Install python test dependencies
install-python-linter:
	pip install pipenv pylint
	cd packages/cli && pipenv install --system --three --dev

## python-lint: Lints python source code
python-lint: python-format
	scripts/lint.sh

## generate-python-client: Generate python models
generate-python-client:
	mkdir -p ${MOCKS_DIR}
	rm -rf ${MOCKS_DIR}/python
	$(SWAGGER_CODEGEN_BIN) generate \
		-i ${SWAGGER_FILE} \
		-l python-flask \
		-o ${MOCKS_DIR}/python \
		--model-package sdk.models \
		-c scripts/swagger/python.conf.json
	# bug in swagger generator
	mv ${MOCKS_DIR}/python/odahuflow/sdk.models/* ${MOCKS_DIR}/python/odahuflow/sdk/models/

	# replace the util script location
	sed -i 's/from odahuflow import util/from odahuflow.sdk.models import util/g' ${MOCKS_DIR}/python/odahuflow/sdk/models/*
	mv ${MOCKS_DIR}/python/odahuflow/util.py ${MOCKS_DIR}/python/odahuflow/sdk/models/

	rm -rf ${PYTHON_MODEL_DIR}
	mkdir -p ${PYTHON_MODEL_DIR}
	cp -r ${MOCKS_DIR}/python/odahuflow/sdk/models/* ${PYTHON_MODEL_DIR}
	git add ${PYTHON_MODEL_DIR}

## install-python-tests: Install python test dependencies
install-python-tests:
	pip install pipenv
	cd packages/cli && pipenv install --system --three --dev

## python-unittests: Run pythoon unit tests
python-unittests:
	DEBUG=true VERBOSE=true pytest -s --cov --cov-report term-missing\
	          packages/cli packages/sdk

## setup-e2e-robot: Prepare a test data for the e2e robot tests
setup-e2e-robot: export VERBOSE=true
setup-e2e-robot:
	odahu-flow-authenticate-test-user ${CLUSTER_PROFILE}

	./packages/tests/stuff/training_stuff.sh setup

## cleanup-e2e-robot: Delete a test data after the e2e robot tests
cleanup-e2e-robot: export VERBOSE=true
cleanup-e2e-robot:
	odahu-flow-authenticate-test-user ${CLUSTER_PROFILE}

	./packages/tests/stuff/training_stuff.sh cleanup

## e2e-robot: Run e2e robot tests
e2e-robot:
	pabot --verbose --processes ${ROBOT_THREADS} \
	      -v CLUSTER_PROFILE:${CLUSTER_PROFILE} \
	      --listener odahuflow.robot.process_reporter \
	      --outputdir target packages/tests/e2e/robot/tests/${ROBOT_FILES}

## update-python-deps: Update all python dependecies in the Pipfiles
update-python-deps:
	scripts/update_python_deps.sh

## install-vulnerabilities-checker: Install the vulnerabilities checker
install-vulnerabilities-checker:
	./scripts/install-git-secrets-hook.sh install_binaries

## check-vulnerabilities: Сheck vulnerabilities in the source code
check-vulnerabilities:
	./scripts/install-git-secrets-hook.sh install_hooks
	git secrets --scan -r


## install-python-formatter: Install python formatter
install-python-formatter:
	pip install black==20.8b1

## python-format: Format python packages
python-format:
	black packages/sdk
	black packages/cli
	black packages/robot

## configure-git: Configure repo git config
configure-git:
	git config blame.ignoreRevsFile .git-blame-ignore-revs

## help: Show the help message
help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sort | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo
