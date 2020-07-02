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
    ${check}                    Call API  toolchain get id  ${MLFLOW}
    should be equal             ${check.spec.default_image}  created
    should be equal             ${check.spec.entrypoint}  created

Create mlflow-gpu toolchain
    Call API                    toolchain post  ${RES_DIR}/valid/mlflow-gpu_create.json
    ${check}                    Call API  toolchain get id  ${MLFLOW_GPU}
    should be equal             ${check.spec.default_image}  created
    should be equal             ${check.spec.entrypoint}  created

Update mlflow toolchain
    sleep                       1s
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow_update.json
    ${check}                    Call API  toolchain get id  ${MLFLOW}
    should be equal             ${check.spec.default_image}  updated
    should be equal             ${check.spec.entrypoint}  updated

Update mlflow-gpu toolchain
    Call API                    toolchain put  ${RES_DIR}/valid/mlflow-gpu_update.yaml
    ${check}                    Call API  toolchain get id  ${MLFLOW_GPU}
    should be equal             ${check.spec.default_image}  updated
    should be equal             ${check.spec.entrypoint}  updated

Get updated list of toolchains
    Command response list should contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}

Get mlflow and mlflow-gpu toolchains by id
    ${result}                   Call API  toolchain get id  ${MLFLOW}
    should be equal             ${result.id}  ${MLFLOW}
    ${result}                   Call API  toolchain get id  ${MLFLOW_GPU}
    should be equal             ${result.id}  ${MLFLOW_GPU}

Delete mlflow toolchain
    ${result}                   Call API  toolchain delete  ${MLFLOW}
    should be equal             ${result.get('message')}  ToolchainIntegration ${MLFLOW} was deleted

Delete mlflow-gpu toolchain
    ${result}                   Call API  toolchain delete  ${MLFLOW_GPU}
    should be equal             ${result.get('message')}  ToolchainIntegration ${MLFLOW_GPU} was deleted

Check that toolchains do not exist
    Command response list should not contain id  toolchain  ${MLFLOW}  ${MLFLOW_GPU}
