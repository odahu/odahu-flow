*** Variables ***
${RES_DIR}                      ${CURDIR}/resources
${ARTIFACT_DIR}                 ${RES_DIR}/artifacts/odahuflow

${LOCAL_CONFIG}                 odahuflow/E2E_local
${MODEL_RESULT}                 {"prediction": [6.3881577909662886, 4.675934265196686], "columns": ["quality"]}

${LOCAL_DOCKER_CONTAINER}       E2E_local_model
${CLUSTER_DOCKER_CONTAINER}     E2E_spec_on_cluster_model

${LOCAL_MODEL_OUTPUT_DIR}       ${CURDIR}/${LOCAL_DOCKER_CONTAINER}
${CLUSTER_MODEL_OUTPUT_DIR}     ${CURDIR}/${CLUSTER_DOCKER_CONTAINER}

*** Settings ***
Documentation       OdahuFlow's API operational check for operations on ModelTraining resources
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Suite Setup         Set Environment Variable    ODAHUFLOW_CONFIG    ${LOCAL_CONFIG}
# Suite Teardown      Remove File     ${LOCAL_CONFIG}
Force Tags          cli  local  e2e
# Test Timeout        30 minutes

*** Keywords ***
Run Training with local spec
    [Arguments]  ${options}
        StrictShell  odahuflowctl --verbose local train run ${options}

Run Packaging with api server spec
    [Arguments]  ${auth data}  ${options}
        StrictShell  odahuflowctl --verbose local pack ${auth data} run ${options}

*** Test Cases ***
# Run E2E local model
#     [Setup]                     StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${LOCAL_MODEL_OUTPUT_DIR}
#     [Teardown]                  Run Keywords
#     ...                         Remove Directory  ${LOCAL_MODEL_OUTPUT_DIR}  recursive=True  AND
#     ...                         Shell  docker stop -t 3 "${LOCAL_DOCKER_CONTAINER}"
#     # training
#     ${result_train}             StrictShell     odahuflowctl --verbose local train run --train-id wine-e2e-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"
#     # check that training artifact exists and take artifact name for packaging
#     ${result_list}              StrictShell  odahuflowctl --verbose local train list
#     ${artifact_name_cmd}        StrictShell  echo "${result_list.stdout}" | tail -n 1 | awk '{ print $2 }'
#     ${artifact_name_dir}        list directory  ${LOCAL_MODEL_OUTPUT_DIR}
#     Should Be Equal As Strings  ${artifact_name_cmd.stdout}  @{artifact_name_dir}
#     # packaing
#     ${pack_result}              StrictShell  odahuflowctl --verbose local pack run --pack-id pack-dir -d "${ARTIFACT_DIR}/dir" --artifact-name ${artifact_name_cmd.stdout}
#     ${image_name}               StrictShell  echo "${pack_result.stdout}" | tail -n 1 | awk '{ print $4 }'
#     # deployment
#     Run   docker run --name "${LOCAL_DOCKER_CONTAINER}" -d --rm -p 5000:5000 ${image_name.stdout}
#
#     Sleep  5 sec
#     Shell     docker container list -as -f name=${LOCAL_DOCKER_CONTAINER}
#     # model invoke
#     ${result_model}              StrictShell  odahuflowctl --verbose model invoke --url http://0:5000 --json-file ${RES_DIR}/request.json
#     Should be equal as Strings  ${result_model.stdout}  ${MODEL_RESULT}

Run E2E spec on cluster model
    [Setup]         Run Keywords
    ...             Login to the api and edge  AND
    ...             Shell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${CLUSTER_MODEL_OUTPUT_DIR}  AND
    ...             StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/dir/e2e.training.yaml
    [Teardown]      Run Keywords
    ...             Remove Directory  ${CLUSTER_MODEL_OUTPUT_DIR}  recursive=True  AND
    ...             Shell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/dir/e2e.training.yaml  AND
    ...             StrictShell  docker stop -t 3 "${CLUSTER_DOCKER_CONTAINER}"

    ${result_train}             StrictShell  odahuflowctl --verbose local train run --train-id e2e-artifact-hardcoded -d "${ARTIFACT_DIR}/file"
    ${artifact_name_dir}        list directory  ${CLUSTER_MODEL_OUTPUT_DIR}
    ${pack_result}              StrictShell  odahuflowctl --verbose local pack run --pack-id e2e-pack-file-image -d "${ARTIFACT_DIR}/dir" --artifact-name my-training

    Create File  ${RES_DIR}/pack_result.txt  ${pack_result.stdout}
    ${image_name}    Shell  (tail -n 1 ${RES_DIR}/pack_result.txt | awk '{ print $4 }'
    Remove File  ${RES_DIR}/pack_result.txt

    Run   docker run --name "${CLUSTER_DOCKER_CONTAINER}" -d --rm -p 5000:5000 ${image_name.stdout}

    Sleep  5 sec
    Shell     docker container list -as -f name=${CLUSTER_DOCKER_CONTAINER}

    ${result_model}              StrictShell  odahuflowctl --verbose model invoke --url http://0:5000 --json-file ${RES_DIR}/request.json
    Should be equal as Strings  ${result_model.stdout}  ${MODEL_RESULT}
