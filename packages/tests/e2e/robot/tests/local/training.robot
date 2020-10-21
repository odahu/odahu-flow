*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow
${RESULT_DIR}               ${CURDIR}/training_train_results

${INPUT_FILE}               ${RES_DIR}/request.json
${OUTPUT_DIR}               ${RES_DIR}
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
# Suite Teardown      Remove File  ${LOCAL_CONFIG}
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
        ${response}  StrictShell  odahuflowctl --verbose gppi -m ${full artifact path} predict ${INPUT_FILE} ${OUTPUT_DIR}
        ${result_path}  StrictShell  echo "${response.stdout}" | tail -n 1 | awk '{ print $3 }'
        log many     ${result_path}  ${result_path.stdout}
        ${response}   Get File  ${result_path.stdout}
        Should be equal as Strings  ${response}  ${MODEL_RESULT}

Run Packaging with api server spec
    [Arguments]  ${auth data}  ${id}  ${file}  ${options}
        Login to the api and edge
        StrictShell  odahuflowctl --verbose pack create -f ${ARTIFACT_DIR}/file/training.yaml
        StrictShell  odahuflowctl --verbose logout

        StrictShell  odahuflowctl --verbose local pack ${auth data} run ${options}

Run Packaging with local spec
    [Arguments]  ${options}
        StrictShell  odahuflowctl --verbose local packaging run ${options}

*** Test Cases ***
Run Valid Training with local spec
    [Template]  Run Training with local spec
    # id	file/dir	output
    --id wine-dir-artifact-template -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}  ${RESULT_DIR}
    --train-id wine-e2e-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"  ${DEFAULT_OUTPUT_DIR}
    --id wine-id-file --manifest-file "${ARTIFACT_DIR}/file/training.yaml" --output ${RESULT_DIR}  ${RESULT_DIR}
    --train-id wine-artifact-hardcoded --manifest-dir "${ARTIFACT_DIR}/dir"  ${DEFAULT_OUTPUT_DIR}

Run Valid Packaging with api server spec

    [Template]  Run Packaging with api server spec
    # id	file/dir	artifact path	artifact name	package-targets

List trainings in default output dir
    ${list result}  StrictShell  odahuflowctl --verbose local train list
    Should contain  ${list result}

Cleanup training artifacts
    StrictShell  odahuflowctl --verbose local train cleanup-artifacts
    ${list result}  StrictShell  odahuflowctl --verbose local train list
    Should be Equal  ${list result}  Artifacts not found
