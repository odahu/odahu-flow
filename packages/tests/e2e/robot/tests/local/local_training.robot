*** Variables ***
${RESULT_DIR}               ${CURDIR}/local_train_results
${DEFAULT_RESULT_DIR}       ~/.odahuflow/local_training/training_output

${LOCAL_CONFIG}             odahuflow/local_training

*** Settings ***
Documentation       Local run of trainings with specs on cluster and host
...                 and packagings with specs on host, accent on trainings.
...                 Validated that training and packaging can be run without (logout from) cluster
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Resource            ./variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Suite Setup         Run Setup Only Once  Local Setup
Suite Teardown      Run Teardown Only Once  Local Cleanup
Force Tags          training
Test Timeout        120 minutes

*** Keywords ***
Try Run Local Training
    [Arguments]  ${error}  ${train options}
        ${result}  FailedShell  odahuflowctl --verbose local training ${train options}
        ${result}  Catenate  ${result.stdout}  ${result.stderr}
        should contain  ${result}  ${error}

Run Local Packaging
    [Arguments]  ${options}
        Run Packaging  5000  ${options}

*** Test Cases ***
Try Run and Fail Training with invalid credentials
    [Tags]   negative
    [Template]  Try Run Local Training
    [Setup]  StrictShell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    ${INVALID_CREDENTIALS_ERROR}    --url "${API_URL}" --token "invalid" run -f ${ARTIFACT_DIR}/file/training.json --id not-exist
    ${MISSED_CREDENTIALS_ERROR}     --url "${API_URL}" --token "${EMPTY}" run -f ${ARTIFACT_DIR}/file/training.json --id not-exist
    ${INVALID_URL_ERROR}            --url "invalid" --token "${AUTH_TOKEN}" run -f ${ARTIFACT_DIR}/file/training.json --id not-exist
    ${INVALID_URL_ERROR}            --url "${EMPTY}" --token "${AUTH_TOKEN}" run -f ${ARTIFACT_DIR}/file/training.json --id not-exist

Try Run and Fail invalid Training
    [Tags]   negative
    [Template]  Try Run Local Training
    # missing required option
    Error: Missing option '--train-id' / '--id'.
    ...  run -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}
    # not valid value for option
    # for file & dir options
    Error: [Errno 21] Is a directory: '${ARTIFACT_DIR}/dir'
    ...  run --id "local-dir-artifact-template" --manifest-file "${ARTIFACT_DIR}/dir" --output ${RESULT_DIR}
    Error: ${ARTIFACT_DIR}/file/training.json is not a directory
    ...  run --id "local id file with spaces" -d "${ARTIFACT_DIR}/file/training.json" --output-dir ${RESULT_DIR}
    Error: Resource file '${ARTIFACT_DIR}/file/not-existing.yaml' not found
    ...  run --id "local id file with spaces" -f "${ARTIFACT_DIR}/file/not-existing.yaml" --manifest-dir "${ARTIFACT_DIR}/not-existing" --output-dir ${RESULT_DIR}
    # no training either locally or on the server
    Error: Got error from server: entity "not-existing-training" is not found (status: 404)
    ...  run --train-id not-existing-training

Run Valid Training with local specs, logout from cluster
    [Template]  Run Local Training
    [Setup]  StrictShell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    # id	file/dir	output
    run --id local-dir-artifact-template -d "${ARTIFACT_DIR}/dir" --manifest-file ${ARTIFACT_DIR}/file/training.json --output-dir ${RESULT_DIR}
    run --train-id local-host-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"
    run --id "local id file with spaces" --manifest-file "${ARTIFACT_DIR}/file/training.json" --manifest-file "${ARTIFACT_DIR}/dir/training_cluster.yaml" --output ${RESULT_DIR}

Run Valid Training with specs on cluster
    [Template]  Run Local Training
    # id	file/dir	output
    run -f ${ARTIFACT_DIR}/dir/packaging.yaml --id local-dir-cluster-artifact-template --output ${DEFAULT_RESULT_DIR}
    --url ${API_URL} --token "${AUTH_TOKEN}" run --train-id local-dir-cluster-artifact-hardcoded

Run Valid Packaging with local spec, logout from cluster
    [Template]  Run Local Packaging
    [Setup]  StrictShell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    # id	file/dir	artifact path	artifact name	package-targets
    run --pack-id local-dir-spec-targets -d ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_RESULT_DIR} --no-disable-package-targets --disable-target docker-push
    run --pack-id local-dir-spec-targets --manifest-dir ${ARTIFACT_DIR}/dir --artifact-path ${RESULT_DIR} -a wine-local-1 --no-disable-package-targets --disable-target docker-push

List trainings in default output dir
    ${list_result}  StrictShell  odahuflowctl --verbose local train list
    Should contain  ${list_result.stdout}  Training artifacts:
    Should contain  ${list_result.stdout}  simple-model
    Should contain  ${list_result.stdout}  wine-local-1
    Should contain  ${list_result.stdout}  wine-cluster-1
    ${lines}  Split To Lines  ${list_result.stdout}
    ${line number}  Get length   ${lines}
    Should be equal as integers  ${line number}  4

Cleanup training artifacts from default output dir
    StrictShell  odahuflowctl --verbose local train cleanup-artifacts
    ${list_result}  StrictShell  odahuflowctl --verbose local train list
    Should be Equal  ${list_result.stdout}  Artifacts not found
