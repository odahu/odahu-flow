*** Variables ***
${RES_DIR}          ${CURDIR}/resources/toolchain
${MLFLOW}           mlflow-api-testing
${MLFLOW_GPU}       mlflow-gpu-api-testing

*** Settings ***
Documentation       API of toolchains
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Toolchain
Suite Setup         Run Keywords
...                 Login to the api and edge
Force Tags          api  sdk  toolchain


*** Test Cases ***
Get list of toolchains
    [Documentation]  check that toolchains that would be created do not exist now
    Command response list should not contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}

Create mlflow toolchain
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow_create.yaml

Create mlflow-gpu toolchain
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow-gpu_create.json

Update mlflow toolchain
    sleep                       1s
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow_update.json

Update mlflow-gpu toolchain
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow-gpu_update.yaml

Get updated list of toolchains
    Command response list should contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}

Get mlflow and mlflow-gpu toolchains by id
    Call API                    toolchain get id  ${MLFLOW}
    Call API                    toolchain get id  ${MLFLOW_GPU}

Delete mlflow toolchain
    Call API                    toolchain delete  ${MLFLOW}

Delete mlflow-gpu toolchain
    Call API                    toolchain delete  ${MLFLOW_GPU}

Check that toolchains do not exist
    Command response list should not contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}
