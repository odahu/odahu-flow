*** Variables ***
${RES_DIR}              ${CURDIR}/resources/tensorflow
${LOCAL_CONFIG}         odahuflow/config_e2e_mlflow_tensorflow
${TENSORFLOW_ID}        test-e2e-tensorflow

*** Settings ***
Documentation       Check tensorflow model
Test Timeout        120 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup example resources  ${TENSORFLOW_ID}
Suite Teardown      Run Keywords
...                 Cleanup example resources  ${TENSORFLOW_ID}  AND
...                 Remove file  ${LOCAL_CONFIG}
Force Tags          e2e  tensorflow  api

*** Test Cases ***
Tensorflow example model
    Download file  mlflow/tensorflow/example/odahuflow/training.odahuflow.yaml  ${RES_DIR}/training.odahuflow.yaml
    Download file  mlflow/tensorflow/example/odahuflow/packaging.odahuflow.yaml  ${RES_DIR}/packaging.odahuflow.yaml
    Download file  mlflow/tensorflow/example/odahuflow/deployment.odahuflow.yaml  ${RES_DIR}/deployment.odahuflow.yaml
    Download file  mlflow/tensorflow/example/odahuflow/request.json  ${RES_DIR}/request.json

    Run example model  ${TENSORFLOW_ID}  ${RES_DIR}