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
CLUSTER_NAME=
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
SWAGGER_FILE=odahuFlow/operator/docs/swagger.yaml
PYTHON_MODEL_DIR=odahuFlow/sdk/odahuflow/sdk/models
SWAGGER_CODEGEN_BIN=java -jar swagger-codegen-cli.jar

HIERA_KEYS_DIR=
ODAHUFLOW_PROFILES_DIR=

CLUSTER_NAME=
CLOUD_PROVIDER=

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
	cd odahuFlow/sdk && \
		rm -rf build dist *.egg-info && \
		pip3 install ${BUILD_PARAMS} -e . && \
		python setup.py sdist && \
    	python setup.py bdist_wheel

## install-cli: Install cli python package
install-cli:
	cd odahuFlow/cli && \
		rm -rf build dist *.egg-info && \
		pip3 install ${BUILD_PARAMS} -e . && \
		python setup.py sdist && \
    	python setup.py bdist_wheel

## install-robot: Install robot tests
install-robot:
	cd odahuFlow/robot && \
		rm -rf build dist *.egg-info && \
		pip3 install ${BUILD_PARAMS} -e . && \
		python setup.py sdist && \
    	python setup.py bdist_wheel

## docker-build-pipeline-agent: Build pipeline agent docker image
docker-build-pipeline-agent:
	docker build -t odahu/odahuflow-pipeline-agent:${BUILD_TAG} -f containers/pipeline-agent/Dockerfile .

## docker-build-edi: Build edi docker image
docker-build-edi:
	docker build --target edi -t odahu/odahuflow-edi:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-model-trainer: Build model builder docker image
docker-build-model-trainer:
	docker build --target model-trainer -t odahu/odahuflow-model-trainer:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-model-packager: Build model packager docker image
docker-build-model-packager:
	docker build --target model-packager -t odahu/odahuflow-model-packager:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-operator: Build operator docker image
docker-build-operator:
	docker build --target operator -t odahu/odahuflow-operator:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-service-catalog: Build service catalog docker image
docker-build-service-catalog:
	docker build --target service-catalog -t odahu/odahuflow-service-catalog:${BUILD_TAG} -f containers/operator/Dockerfile .

## docker-build-feedback-aggregator: Build feedback aggregator image
docker-build-feedback-aggregator:
	docker build --target aggregator -t odahu/odahuflow-feedback-aggregator:${BUILD_TAG} -f containers/feedback-aggregator/Dockerfile .

## docker-build-feedback-collector: Build feedback collector image
docker-build-feedback-collector:
	docker build --target collector -t odahu/odahuflow-feedback-aggregator:${BUILD_TAG} -f containers/feedback-aggregator/Dockerfile .

## docker-build-fluentd: Build fluentd image
docker-build-fluentd:
	docker build -t odahu/fluentd:${BUILD_TAG} -f containers/fluentd/Dockerfile .

## docker-build-all: Build all docker images
docker-build-all:  docker-build-pipeline-agent docker-build-edi  docker-build-model-trainer  docker-build-operator  docker-build-feedback-aggregator docker-build-fluentd docker-build-service-catalog

## docker-push-pipeline-agent: Push pipeline agent docker image
docker-push-pipeline-agent:
	docker tag odahu/odahuflow-pipeline-agent:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahuflow-pipeline-agent:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahuflow-pipeline-agent:${TAG}

## docker-push-edi: Push edi docker image
docker-push-edi:  check-tag
	docker tag odahu/odahuflow-edi:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahuflow-edi:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahuflow-edi:${TAG}

## docker-push-model-trainer: Push model builder docker image
docker-push-model-trainer:  check-tag
	docker tag odahu/odahuflow-model-trainer:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/model-trainer:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahuflow-model-trainer:${TAG}

## docker-push-operator: Push operator docker image
docker-push-operator:  check-tag
	docker tag odahu/odahuflow-operator:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahuflow-operator:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/odahuflow-operator:${TAG}

## docker-push-service-catalog: Push service catalog docker image
docker-push-service-catalog:  check-tag
	docker tag odahu/odahuflow-service-catalog:${BUILD_TAG} ${DOCKER_REGISTRY}/odahu/odahuflow-service-catalog:${TAG}
	docker push ${DOCKER_REGISTRY}/odahu/service-catalog:${TAG}

## docker-push-feedback-aggregator: Push feedback aggregator docker image
docker-push-feedback-aggregator:  check-tag
	docker tag odahu/odahuflow-feedback-aggregator:${BUILD_TAG} ${DOCKER_REGISTRY}odahu/odahuflow-feedback-aggregator:${TAG}
	docker push ${DOCKER_REGISTRY}odahu/odahuflow-feedback-aggregator:${TAG}

## docker-push-fluentd: Push fluentd docker image
docker-push-fluentd:  check-tag
	docker tag odahu/fluentd:${BUILD_TAG} ${DOCKER_REGISTRY}odahu/fluentd:${TAG}
	docker push ${DOCKER_REGISTRY}odahu/fluentd:${TAG}

## docker-push-all: Push all docker images
docker-push-all:  docker-push-pipeline-agent docker-push-edi  docker-push-model-trainer  docker-push-operator  docker-push-feedback-aggregator docker-push-fluentd docker-push-service-catalog

## helm-install: Install the odahuflow helm chart from source code
helm-install: helm-delete
	helm install helms/odahuflow --atomic --wait --timeout 320 --namespace odahuflow --name odahuflow --debug ${HELM_ADDITIONAL_PARAMS}

## helm-delete: Delete the odahuflow helm release
helm-delete:
	helm delete --purge odahuflow || true

## install-python-linter: Install python test dependencies
install-python-linter:
	pip install pipenv
	cd containers/pipeline-agent && pipenv install --system --three --dev

## python-lint: Lints python source code
python-lint:
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
	cd odahuFlow/cli && pipenv install --system --three --dev

## python-unittests: Run pythoon unit tests
python-unittests:
	DEBUG=true VERBOSE=true pytest \
	          odahuFlow/cli

## setup-e2e-robot: Prepare a test data for the e2e robot tests
setup-e2e-robot:
	odahuflow-authenticate-test-user ${CLUSTER_PROFILE}

	./odahuFlow/tests/stuff/training_stuff.sh setup

## cleanup-e2e-robot: Delete a test data after the e2e robot tests
cleanup-e2e-robot:
	odahuflow-authenticate-test-user ${CLUSTER_PROFILE}

	./odahuFlow/tests/stuff/training_stuff.sh cleanup

## e2e-robot: Run e2e robot tests
e2e-robot:
	pabot --verbose --processes ${ROBOT_THREADS} \
	      -v CLUSTER_PROFILE:${CLUSTER_PROFILE} \
	      --listener odahuflow.robot.process_reporter \
	      --outputdir target odahuFlow/tests/e2e/robot/tests/${ROBOT_FILES}

## update-python-deps: Update all python dependecies in the Pipfiles
update-python-deps:
	scripts/update_python_deps.sh

## install-vulnerabilities-checker: Install the vulnerabilities-checker
install-vulnerabilities-checker:
	./scripts/install-git-secrets-hook.sh install_binaries

## check-vulnerabilities: Сheck vulnerabilities in the source code
check-vulnerabilities:
	./scripts/install-git-secrets-hook.sh install_hooks
	git secrets --scan -r

## help: Show the help message
help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sort | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo
