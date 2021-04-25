*** Variables ***
${RES_DIR}              ${CURDIR}/resources/wine
${LOCAL_CONFIG}         odahuflow/config_e2e_mlflow_wine
${WINE_ID}              test-e2e-wine

*** Settings ***
Documentation       Check wine model
Test Timeout        120 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         E2E test setup  ${WINE_ID}
Suite Teardown      E2E test teardown  ${WINE_ID}
Force Tags          e2e  wine  cli

*** Test Cases ***
Wine model
    Download file  mlflow/sklearn/wine/odahuflow/training.odahuflow.yaml  ${RES_DIR}/training.odahuflow.yaml
    Download file  mlflow/sklearn/wine/odahuflow/packaging.odahuflow.yaml  ${RES_DIR}/packaging.odahuflow.yaml
    Download file  mlflow/sklearn/wine/odahuflow/deployment.odahuflow.yaml  ${RES_DIR}/deployment.odahuflow.yaml
    Download file  mlflow/sklearn/wine/odahuflow/request.json  ${RES_DIR}/request.json

    Run example model  ${WINE_ID}  ${RES_DIR}
