*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_toolchain
${RES_DIR}              ${CURDIR}/resources/toolchain
${MLFLOW}               mlflow-api-testing
${MLFLOW_NOT_EXIST}     mlflow-api-not-exist

*** Settings ***
Documentation       API of toolchains
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.Toolchain
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup All Resources
Suite Teardown      Run Keywords
...                 Cleanup All Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  toolchain
Test Timeout        5 minutes

*** Keywords ***
Cleanup All Resources
    Cleanup resource  toolchain-integration  ${MLFLOW}
    Cleanup resource  toolchain-integration  ${MLFLOW_NOT_EXIST}

*** Test Cases ***
Get list of toolchains
    [Documentation]  check that toolchains that would be created do not exist now
    Command response list should not contain id  toolchain  ${MLFLOW}

Create mlflow toolchain
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow_create.yaml
    ${check}                    Call API  toolchain get id  ${MLFLOW}
    Default Docker image should be equal  ${check}  created
    Default Entrypoint should be equal  ${check}  created

Update mlflow toolchain
    sleep                       1s
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow_update.json
    ${check}                    Call API  toolchain get id  ${MLFLOW}
    Default Docker image should be equal  ${check}  updated
    Default Entrypoint should be equal  ${check}  updated

Get updated list of toolchains
    Command response list should contain id  toolchain  ${MLFLOW}

Get mlflow toolchains by id
    ${result}                   Call API  toolchain get id  ${MLFLOW}
    ID should be equal          ${result}  ${MLFLOW}

Delete mlflow toolchain
    ${result}                   Call API  toolchain delete  ${MLFLOW}
    should be equal             ${result.get('message')}  ToolchainIntegration ${MLFLOW} was deleted

Check that toolchains do not exist
    Command response list should not contain id  toolchain  ${MLFLOW}

#############################
#    NEGATIVE TEST CASES    #
#############################
Try Create Toolchain that already exists
    [Tags]                      negative
    [Setup]                     Cleanup resource  toolchain-integration  ${MLFLOW}
    [Teardown]                  Cleanup resource  toolchain-integration  ${MLFLOW}
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow_update.json
    ${EntityAlreadyExists}      format string  ${409 Conflict Template}  ${MLFLOW}
    Call API and get Error      ${EntityAlreadyExists}  toolchain post  ${RES_DIR}/valid/mlflow_update.json

Try Update not existing Toolchain
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${404NotFound}  toolchain put  ${RES_DIR}/invalid/mlflow_update_not_exist.yaml

Try Update deleted Toolchain
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW}
    Call API and get Error      ${404NotFound}  toolchain put  ${RES_DIR}/valid/mlflow_create.yaml

Try Get id not existing Toolchain
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${404NotFound}  toolchain get id  ${MLFLOW_NOT_EXIST}

Try Get id deleted Toolchain
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW}
    Call API and get Error      ${404NotFound}  toolchain get id  ${MLFLOW}

Try Delete not existing Toolchain
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${404NotFound}  toolchain delete  ${MLFLOW_NOT_EXIST}

Try Delete deleted Toolchain
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW}
    Call API and get Error      ${404NotFound}  toolchain delete  ${MLFLOW}
