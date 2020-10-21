*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow
${RESULT_DIR}               ${CURDIR}/training_train_results

${INPUT_FILE}               ${RES_DIR}/request.json
${DEFAULT_OUTPUT_DIR}       ~/.odahuflow/local_training/training_output

${MODEL_RESULT}             {"prediction": [6.3881577909662886, 4.675934265196686], "columns": ["quality"]}
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
...                 StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${DEFAULT_OUTPUT_DIR}
# Suite Teardown    Run Keywords
...                 Remove Directory  ${RESULT_DIR}  recursive=True  AND
...                 Remove Directory  ${DEFAULT_OUTPUT_DIR}  recursive=True  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          cli  local  training
# Test Timeout        90 minutes

*** Keywords ***
Run Training with local spec
    [Teardown]  Remove File  ${RES_DIR}/results.json
    [Arguments]  ${train options}  ${artifact path}
        ${result}  StrictShell  odahuflowctl --verbose local train run ${train options}

        # fetch the training artifact name from stdout
        Create File  ${RES_DIR}/train_result.txt  ${result.stdout}
        ${artifact_name}    Shell  (tail -n 1 ${RES_DIR}/train_result.txt | awk '{ print $2 }')
        Remove File  ${RES_DIR}/train_result.txt
        ${full artifact path}  set variable  ${artifact path}/${artifact_name.stdout}

        # check the training artifact validity
        ${response}  StrictShell  odahuflowctl --verbose gppi -m ${full artifact path} predict ${INPUT_FILE} ${RESULT_DIR}
        ${result_path}  StrictShell  echo "${response.stdout}" | tail -n 1 | awk '{ print $3 }'

        ${response}   Get File  ${result_path.stdout}
        Should be equal as Strings  ${response}  ${MODEL_RESULT}

Run Packaging with api server spec
    [Arguments]  ${command}
        StrictShell  odahuflowctl --verbose ${command}

*** Test Cases ***
Run Valid Training with local spec
    [Template]  Run Training with local spec
    # id	file/dir	output
    --id wine-dir-artifact-template -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}  ${RESULT_DIR}
    --train-id wine-e2e-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"  ${DEFAULT_OUTPUT_DIR}
    --id wine-id-file --manifest-file "${ARTIFACT_DIR}/file/training.yaml" --output ${RESULT_DIR}  ${RESULT_DIR}
    --train-id train-artifact-hardcoded --manifest-dir "${ARTIFACT_DIR}/dir"  ${DEFAULT_OUTPUT_DIR}

Run Valid Packaging with api server spec
    [Setup]     Run Keywords
    ...         Login to the api and edge  AND
    ...         StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Teardown]  StrictShell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Template]  Run Packaging with api server spec
    # id	file/dir	artifact path	artifact name	package-targets
    local pack run -f ${ARTIFACT_DIR}/dir/packaging --id pack-dir --output ${RESULT_DIR}/wine-dir-1.0 --artifact-name wine-dir-1.0
    local pack --url ${API_URL} --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image

List trainings in default output dir
    ${list_result}  StrictShell  odahuflowctl --verbose local train list
    Should contain  ${list_result.stdout}  Training artifacts:
    Should contain  ${list_result.stdout}  my-training
    Should contain  ${list_result.stdout}  wine-name-1
    ${line number}  Split To Lines  ${list_result.stdout}
    ${line number}  Get length   ${line number}
    Should be equal as integers  ${line number}  3

Cleanup training artifacts from default output dir
    StrictShell  odahuflowctl --verbose local train cleanup-artifacts
    ${list_result}  StrictShell  odahuflowctl --verbose local train list
    Should be Equal  ${list_result.stdout}  Artifacts not found
