*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow
${RESULT_DIR}               ${CURDIR}/training_train_results

${INPUT_FILE}               ${RES_DIR}/request.json
${DEFAULT_RESULT_DIR}       ~/.odahuflow/local_training/training_output

${LOCAL_CONFIG}             odahuflow/local_training


*** Settings ***
Documentation       local trainings & packagings with spec on cluster
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${DEFAULT_RESULT_DIR}
Suite Teardown      Run Keywords
...                 Remove Directory  ${RESULT_DIR}  recursive=True  AND
...                 Remove Directory  ${DEFAULT_RESULT_DIR}  recursive=True  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          cli  local  training
Test Timeout        180 minutes

*** Keywords ***
Run Training with local spec
    [Arguments]  ${train options}  ${artifact path}
        ${result}  StrictShell  odahuflowctl --verbose local train run ${train options}

        # fetch the training artifact name from stdout
        Create File  ${RESULT_DIR}/train_result.txt  ${result.stdout}
        ${artifact_name}    StrictShell  tail -n 1 ${RESULT_DIR}/train_result.txt | awk '{ print $2 }'
        Remove File  ${RESULT_DIR}/train_result.txt
        ${full artifact path}  set variable  ${artifact path}/${artifact_name.stdout}

        # check the training artifact validity
        ${response}  StrictShell  odahuflowctl --verbose gppi -m ${full artifact path} predict ${INPUT_FILE} ${RESULT_DIR}
        ${result_path}  StrictShell  echo "${response.stdout}" | tail -n 1 | awk '{ print $3 }'

        ${response}   Get File  ${result_path.stdout}
        Should be equal as Strings  ${response}  ${WINE_MODEL_RESULT}

Try Run Training with local spec
    [Arguments]  ${error}  ${train options}
        ${result}  FailedShell  odahuflowctl --verbose local train run ${train options}
        ${result}  Catenate  ${result.stdout}  ${result.stderr}
        should contain  ${result}  ${error}

Run Packaging with api server spec
    [Teardown]  Shell  docker rm -f ${container_id.stdout}
    [Arguments]  ${command}
        ${pack_result}  StrictShell  odahuflowctl --verbose ${command}

        Create File  ${RESULT_DIR}/pack_result.txt  ${pack_result.stdout}
        ${image_name}  StrictShell  tail -n 1 ${RESULT_DIR}/pack_result.txt | awk '{ print $4 }'
        Remove File  ${RESULT_DIR}/pack_result.txt

        StrictShell  docker images --all
        ${container_id}  StrictShell  docker run -d --rm -p 5001:5000 ${image_name.stdout}

        Sleep  5 sec
        StrictShell  docker container list -as -f id=${container_id.stdout}

        ${MODEL_HOST}    Get local model host
        ${result_model}  StrictShell  odahuflowctl --verbose model invoke --url ${MODEL_HOST}:5001 --json-file ${RES_DIR}/request.json
        Should be equal as Strings  ${result_model.stdout}  ${WINE_MODEL_RESULT}

Try Run Packaging with api server spec
    [Arguments]  ${error}  ${command}
        ${result}  FailedShell  odahuflowctl --verbose ${command}
        should contain  ${result.stdout}  ${error}

*** Test Cases ***
Try Run and Fail Packaging with invalid credentials
    [Tags]   negative
    [Setup]  StrictShell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    [Template]  Try Run Packaging with api server spec
    ${INVALID_CREDENTIALS_ERROR}    local pack --url ${API_URL} --token "invalid" run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image
    ${MISSED_CREDENTIALS_ERROR}     local pack --url ${API_URL} --token "${EMPTY}" run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image
    ${INVALID_URL_ERROR}            local pack --url "invalid" --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image
    ${INVALID_URL_ERROR}            local pack --url "${EMPTY}" --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image

Try Run and Fail invalid Training
    [Tags]   negative
    [Setup]  Login to the api and edge
    [Teardown]  Shell  odahuflowctl logout
    [Template]  Try Run Training with local spec
    # missing required option
    Error: Missing option '--train-id' / '--id'.
    ...  -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}
    # not valid value for option
    # for file & dir options
    Error: [Errno 21] Is a directory: '${ARTIFACT_DIR}/dir'
    ...  --id "wine-dir-artifact-template" --manifest-file "${ARTIFACT_DIR}/dir" --output ${RESULT_DIR}
    Error: ${ARTIFACT_DIR}/file/training.yaml is not a directory
    ...  --id "local id file with spaces" -d "${ARTIFACT_DIR}/file/training.yaml" --output-dir ${RESULT_DIR}
    Error: Resource file '${ARTIFACT_DIR}/file/not-existing.yaml' not found
    ...  --id "local id file with spaces" -f "${ARTIFACT_DIR}/file/not-existing.yaml" --manifest-dir "${ARTIFACT_DIR}/not-existing" --output-dir ${RESULT_DIR}
    # no training either locally or on the server
    Error: Got error from server: entity "not-existing-training" is not found (status: 404)
    ...  --train-id not-existing-training

Run Valid Training with local spec
    [Template]  Run Training with local spec
    # id	file/dir	output
    --id wine-dir-artifact-template -d "${ARTIFACT_DIR}/dir" --manifest-file ${ARTIFACT_DIR}/file/training.yaml --output-dir ${RESULT_DIR}  ${RESULT_DIR}
    --train-id wine-e2e-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"  ${DEFAULT_RESULT_DIR}
    --id "local id file with spaces" --manifest-file "${ARTIFACT_DIR}/file/training.yaml" --manifest-file "${ARTIFACT_DIR}/dir/training_cluster.json" --output ${RESULT_DIR}  ${RESULT_DIR}
    --train-id train-artifact-hardcoded --manifest-dir "${ARTIFACT_DIR}/dir"  ${DEFAULT_RESULT_DIR}

Run Valid Packaging with api server spec
    [Setup]     Run Keywords
    ...         Login to the api and edge  AND
    ...         StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Teardown]  Shell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Template]  Run Packaging with api server spec
    # id	file/dir	artifact path	artifact name	package-targets
    local pack run -f ${ARTIFACT_DIR}/dir/packaging --id pack-dir --artifact-path ${RESULT_DIR} --artifact-name wine-dir-1.0 --no-disable-package-targets
    local packaging --url ${API_URL} --token ${AUTH_TOKEN} run --id pack-dir --artifact-name simple-model --no-disable-package-targets --disable-target docker-push

List trainings in default output dir
    ${list_result}  StrictShell  odahuflowctl --verbose local train list
    Should contain  ${list_result.stdout}  Training artifacts:
    Should contain  ${list_result.stdout}  simple-model
    Should contain  ${list_result.stdout}  wine-name-1
    ${line number}  Split To Lines  ${list_result.stdout}
    ${line number}  Get length   ${line number}
    Should be equal as integers  ${line number}  3

Cleanup training artifacts from default output dir
    StrictShell  odahuflowctl --verbose local train cleanup-artifacts
    ${list_result}  StrictShell  odahuflowctl --verbose local train list
    Should be Equal  ${list_result.stdout}  Artifacts not found
