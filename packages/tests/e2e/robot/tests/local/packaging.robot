*** Variables ***
*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${RESULT_DIR}               ${CURDIR}/packaging_train_results
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow

${INPUT_FILE}               ${RES_DIR}/request.json
${OUTPUT_DIR}               ${RES_DIR}
${DEFAULT_OUTPUT_DIR}       ~/.odahuflow/local_packaging/training_output

${MODEL_RESULT}             {"prediction": [6.3881577909662886, 4.675934265196686], "columns": ["quality"]}

${LOCAL_CONFIG}             odahuflow/local_packaging

*** Settings ***
Documentation       trainings with spec on cluster & local packagings
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${DEFAULT_OUTPUT_DIR}
# Suite Teardown      Remove File  ${LOCAL_CONFIG}
Force Tags          cli  local  packaging
# Test Timeout        90 minutes

*** Keywords ***
Run Training with api server spec
    [Arguments]  ${command}
        StrictShell  odahuflowctl --verbose ${command}

Run Packaging with local spec
    [Arguments]  ${options}
        StrictShell  odahuflowctl --verbose local packaging run ${options}

*** Test Cases ***
Run Valid Training with api server spec
    [Setup]     Run Keywords
    ...         Login to the api and edge  AND
    ...         StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/dir/training_cluster.json
    [Teardown]  StrictShell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/dir/training_cluster.json
    [Template]  Run Training with api server spec
    # auth data     id      file/dir        output
    local training run -f ${ARTIFACT_DIR}/dir/packaging --id wine-packaging --output ${RESULT_DIR}
    local training --url ${API_URL} --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/training.yaml --id pack-artifact-hardcoded

Run Valid Packaging with local spec
    [Template]  Run Packaging with local spec
    # id	file/dir	artifact path	artifact name	package-targets
    --id pack-dir -d ${ARTIFACT_DIR}/dir --no-disable-package-targets
    --pack-id pack-file-image -f ${ARTIFACT_DIR}/file/packaging.yaml --artifact-path ${RESULT_DIR}/wine-name-1 --artifact-name wine-name-1
    --id pack-dir --manifest-dir ${ARTIFACT_DIR}/dir --disable-package-targets
    --pack-id pack-dir -d ${ARTIFACT_DIR}/dir --artifact-path ${RESULT_DIR}/wine-name-1 --disable-package-targets
    --id pack-file-image -f ${ARTIFACT_DIR}/file/packaging.yaml -a ${RESULT_DIR}/wine-name-1 --no-disable-package-targets
    --pack-id pack-dir --manifest-dir ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_OUTPUT_DIR}/my-training
    --id pack-file-image --manifest-file ${ARTIFACT_DIR}/file/packaging.yaml -a wine-name-1 --disable-package-targets
