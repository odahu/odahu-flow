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

Try Run Training
    [Arguments]  ${error}  ${train options}
        ${result}  FailedShell  odahuflowctl --verbose ${train options}
        ${result}  Catenate  ${result.stdout}  ${result.stderr}
        should contain  ${result}  ${error}

*** Test Cases ***
Try Run and Fail Training with invalid credentials
    [Tags]   negative
    [Setup]  StrictShell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    [Template]  Try Run Training
    ${INVALID_CREDENTIALS_ERROR}    local training --url "${API_URL}" --token "invalid" run -f ${ARTIFACT_DIR}/file/training.yaml --id not-exist
    ${MISSED_CREDENTIALS_ERROR}     local training --url "${API_URL}" --token "${EMPTY}" run -f ${ARTIFACT_DIR}/file/training.yaml --id not-exist
    ${INVALID_URL_ERROR}            local training --url "invalid" --token "${AUTH_TOKEN}" run -f ${ARTIFACT_DIR}/file/training.yaml --id not-exist
    ${INVALID_URL_ERROR}            local training --url "${EMPTY}" --token "${AUTH_TOKEN}" run -f ${ARTIFACT_DIR}/file/training.yaml --id not-exist

Try Run and Fail invalid Training
    [Tags]   negative
    [Setup]  Login to the api and edge
    [Teardown]  Shell  odahuflowctl logout
    [Template]  Try Run Training
    # missing required option
    Error: Missing option '--train-id' / '--id'.
    ...  run -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}
    # not valid value for option
    # for file & dir options
    Error: [Errno 21] Is a directory: '${ARTIFACT_DIR}/dir'
    ...  run --id "local-dir-artifact-template" --manifest-file "${ARTIFACT_DIR}/dir" --output ${RESULT_DIR}
    Error: ${ARTIFACT_DIR}/file/training.yaml is not a directory
    ...  run --id "local id file with spaces" -d "${ARTIFACT_DIR}/file/training.yaml" --output-dir ${RESULT_DIR}
    Error: Resource file '${ARTIFACT_DIR}/file/not-existing.yaml' not found
    ...  run --id "local id file with spaces" -f "${ARTIFACT_DIR}/file/not-existing.yaml" --manifest-dir "${ARTIFACT_DIR}/not-existing" --output-dir ${RESULT_DIR}
    # no training either locally or on the server
    Error: Got error from server: entity "not-existing-training" is not found (status: 404)
    ...  run --train-id not-existing-training

Run Valid Training with local spec
    [Template]  Run Training with local spec
    # id	file/dir	output
    --id local-dir-artifact-template -d "${ARTIFACT_DIR}/dir" --manifest-file ${ARTIFACT_DIR}/file/training.yaml --output-dir ${RESULT_DIR}  ${RESULT_DIR}
    --train-id local-host-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"  ${DEFAULT_RESULT_DIR}
    --id "local id file with spaces" --manifest-file "${ARTIFACT_DIR}/file/training.yaml" --manifest-file "${ARTIFACT_DIR}/dir/training_cluster.json" --output ${RESULT_DIR}  ${RESULT_DIR}
    --train-id local-dir-cluster-artifact-hardcoded --manifest-dir "${ARTIFACT_DIR}/dir"  ${DEFAULT_RESULT_DIR}

Run Valid Packaging with local spec
    [Setup]     StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Teardown]  Shell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Template]  Run Packaging with local spec
    # id	file/dir	artifact path	artifact name	package-targets
    run --pack-id local-dir-spec-targets -d ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_RESULT_DIR} --disable-package-targets
    run --pack-id local-dir-spec-targets --manifest-dir ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_RESULT_DIR} -a wine-local-1.0

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
