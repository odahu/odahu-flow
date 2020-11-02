*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow
${RESULT_DIR}               ${CURDIR}/training_train_results

${INPUT_FILE}               ${RES_DIR}/request.json
${DEFAULT_RESULT_DIR}       ~/.odahuflow/local_training/training_output

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
...                 StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${DEFAULT_RESULT_DIR}
# Suite Teardown    Run Keywords
# ...                 Remove Directory  ${RESULT_DIR}  recursive=True  AND
# ...                 Remove Directory  ${DEFAULT_RESULT_DIR}  recursive=True  AND
# ...                 Remove File  ${LOCAL_CONFIG}
Force Tags          cli  local  training
# Test Timeout        90 minutes

*** Keywords ***
Run Training with local spec
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

Try Run Training with local spec
    [Arguments]  ${output}  ${error}  ${train options}
        ${result}  FailedShell  odahuflowctl --verbose local train run ${train options}
        Run Keyword And Continue On Failure  log many  ${result}  ${result.stdout}  ${result.stderr}
        should contain  ${result.${output}}  ${error}

Run Packaging with api server spec
    [Arguments]  ${command}
        ${pack_result}  StrictShell  odahuflowctl --verbose ${command}

        Create File  ${RES_DIR}/pack_result.txt  ${pack_result.stdout}
        ${image_name}    Shell  tail -n 1 ${RES_DIR}/pack_result.txt | awk '{ print $4 }'
        Remove File  ${RES_DIR}/pack_result.txt

        StrictShell  docker images --all
        ${container_id}  StrictShell  docker run -d --rm -p 5001:5000 ${image_name.stdout}

        Sleep  5 sec
        Shell  docker container list -as -f id=${container_id.stdout}

        ${result_model}             StrictShell  odahuflowctl --verbose model invoke --url http://0:5001 --json-file ${RES_DIR}/request.json
        Should be equal as Strings  ${result_model.stdout}  ${MODEL_RESULT}

Try Run Packaging with api server spec
    [Arguments]  ${error}  ${command}
        ${result}  FailedShell  odahuflowctl --verbose ${command}
        should contain  ${result.stdout}  ${error}

*** Test Cases ***
Run Valid Training with local spec
    [Template]  Run Training with local spec
    # id	file/dir	output
    --id wine-dir-artifact-template -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}  ${RESULT_DIR}
    --train-id wine-e2e-default-template -f "${ARTIFACT_DIR}/file/training.default.artifact.template.json"  ${DEFAULT_RESULT_DIR}
    --id "wine id file" --manifest-file "${ARTIFACT_DIR}/file/training.yaml" --output ${RESULT_DIR}  ${RESULT_DIR}
    --train-id train-artifact-hardcoded --manifest-dir "${ARTIFACT_DIR}/dir"  ${DEFAULT_RESULT_DIR}

Run Valid Packaging with api server spec
    [Setup]     Run Keywords
    ...         Login to the api and edge  AND
    ...         StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Teardown]  Shell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Template]  Run Packaging with api server spec
    # id	file/dir	artifact path	artifact name	package-targets
    local pack run -f ${ARTIFACT_DIR}/dir/packaging --id pack-dir --artifact-path ${RESULT_DIR}/wine-dir-1.0 --artifact-name wine-dir-1.0
    local pack --url ${API_URL} --token ${AUTH_TOKEN} run --id simple-model -a simple-model

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

# negative tests
Try Run invalid Training with local spec
    [Template]  Try Run Training with local spec
    # missing required option
    Error  -d "${ARTIFACT_DIR}/dir" --output-dir ${RESULT_DIR}
    Error  --id "wine id file" --output-dir ${RESULT_DIR}
    # incompatible options
    Error  --id "wine id file" -d "${ARTIFACT_DIR}/dir" --manifest-file "${ARTIFACT_DIR}/file/training.yaml"
    Error  --train-id wine-e2e-default-template --id wine-e2e -f ${RESULT_DIR} --manifest-dir "${ARTIFACT_DIR}/file/training.yaml"
    Error  --id "wine_id_file" -d "${ARTIFACT_DIR}/dir" --manifest-dir "${ARTIFACT_DIR}/file" --output ${RESULT_DIR}
    # not valid value for option
    # for file & dir options
    Error  --id "wine-dir-artifact-template" --manifest-file "${ARTIFACT_DIR}/dir" --output ${RESULT_DIR}
    Error  --id "wine id file" -d "${ARTIFACT_DIR}/file/training.yaml" --output-dir ${RESULT_DIR}
    # no training either locally or on the server
    Error  --train-id not-existing-training

Try Run invalid Packaging with api server spec
    [Setup]  Shell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    [Template]  Try Run Packaging with api server spec
    # invalid credentials
    Error  local pack --url ${API_URL} --token "invalid" run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image
    Error  local pack --url ${API_URL} --token ${EMPTY} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image
    Error  local pack --url "invalid" --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image
    Error  local pack --url ${EMPTY} --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-file-image