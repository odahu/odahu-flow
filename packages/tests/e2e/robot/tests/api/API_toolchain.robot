*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_toolchain
${RES_DIR}              ${CURDIR}/resources/toolchain
${MLFLOW}               mlflow-api-testing
${MLFLOW_GPU}           mlflow-gpu-api-testing
${MLFLOW_NOT_EXIST}     mlflow-gpu-api-not-exist

*** Settings ***
Documentation       API of toolchains
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
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
    Cleanup resource  toolchain-integration  ${MLFLOW_GPU}
    Cleanup resource  toolchain-integration  ${MLFLOW_NOT_EXIST}

*** Test Cases ***
Get list of toolchains
    [Documentation]  check that toolchains that would be created do not exist now
    Command response list should not contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}

Create mlflow toolchain
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow_create.yaml
    ${check}                    Call API  toolchain get id  ${MLFLOW}
    Default Docker image should be equal  ${check}  created
    Default Entrypoint should be equal  ${check}  created

Create mlflow-gpu toolchain
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow-gpu_create.json
    ${check}                    Call API  toolchain get id  ${MLFLOW_GPU}
    Default Docker image should be equal  ${check}  created
    Default Entrypoint should be equal  ${check}  created

Update mlflow toolchain
    sleep                       1s
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow_update.json
    ${check}                    Call API  toolchain get id  ${MLFLOW}
    Default Docker image should be equal  ${check}  updated
    Default Entrypoint should be equal  ${check}  updated

Update mlflow-gpu toolchain
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow-gpu_update.yaml
    ${check}                    Call API  toolchain get id  ${MLFLOW_GPU}
    Default Docker image should be equal  ${check}  updated
    Default Entrypoint should be equal  ${check}  updated

Get updated list of toolchains
    Command response list should contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}

Get mlflow and mlflow-gpu toolchains by id
    ${result}                   Call API  toolchain get id  ${MLFLOW}
    ID should be equal          ${result}  ${MLFLOW}
    ${result}                   Call API  toolchain get id  ${MLFLOW_GPU}
    ID should be equal          ${result}  ${MLFLOW_GPU}

Delete mlflow toolchain
    ${result}                   Call API  toolchain delete  ${MLFLOW}
    should be equal             ${result.get('message')}  ToolchainIntegration ${MLFLOW} was deleted

Delete mlflow-gpu toolchain
    ${result}                   Call API  toolchain delete  ${MLFLOW_GPU}
    should be equal             ${result.get('message')}  ToolchainIntegration ${MLFLOW_GPU} was deleted

Check that toolchains do not exist
    Command response list should not contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}

#############################
#    NEGATIVE TEST CASES    #
#############################
Try Create Toolchain that already exists
    [Tags]                      negative
    [Teardown]                  Cleanup resource  toolchain-integration  ${MLFLOW}
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow_update.json
    ${EntityAlreadyExists}      Format EntityAlreadyExists  ${MLFLOW}
    Call API and get Error      ${EntityAlreadyExists}  toolchain post  ${RES_DIR}/valid/mlflow_update.json

Try Update not existing and deleted Toolchain
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  toolchain put  ${RES_DIR}/invalid/mlflow_gpu_update_not_exist.yaml
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${MLFLOW_GPU}
    Call API and get Error      ${WrongHttpStatusCode}  toolchain put  ${RES_DIR}/valid/mlflow-gpu_create.json

Try Get id not existing and deleted Toolchain
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  toolchain get id  ${MLFLOW_NOT_EXIST}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${MLFLOW}
    Call API and get Error      ${WrongHttpStatusCode}  toolchain get id  ${MLFLOW}

Try Delete not existing and deleted Toolchain
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  toolchain delete  ${MLFLOW_NOT_EXIST}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${MLFLOW_GPU}
    Call API and get Error      ${WrongHttpStatusCode}  toolchain delete  ${MLFLOW_GPU}
