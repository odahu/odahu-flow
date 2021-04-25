*** Variables ***
${RES_DIR}              ${CURDIR}/resources/flower_classifier
${LOCAL_CONFIG}         odahuflow/config_e2e_mlflow_flower_classifier
${FLOWER_CLASSIFIER}              test-e2e-flower-classifier

*** Settings ***
Documentation       Check flower classifier model
Test Timeout        120 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         E2E test setup  ${FLOWER_CLASSIFIER}
Suite Teardown      E2E test teardown  ${FLOWER_CLASSIFIER}
Force Tags          e2e  flower-classifier  cli

*** Test Cases ***
Flower classifier model
    Pass Execution If  not ${IS_GPU_ENABLED}  GPU node pools is not enabled on the cluster

    Download file  mlflow/tensorflow/flower_classifier/odahuflow/training.gpu.odahuflow.yaml  ${RES_DIR}/training.odahuflow.yaml
    Download file  mlflow/tensorflow/flower_classifier/odahuflow/packaging.odahuflow.yaml  ${RES_DIR}/packaging.odahuflow.yaml
    Download file  mlflow/tensorflow/flower_classifier/odahuflow/deployment.odahuflow.yaml  ${RES_DIR}/deployment.odahuflow.yaml
    Download file  mlflow/tensorflow/flower_classifier/odahuflow/request.json  ${RES_DIR}/request.json

    Run example model  ${FLOWER_CLASSIFIER}  ${RES_DIR}
