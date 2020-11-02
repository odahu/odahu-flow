*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${RESULT_DIR}               ${CURDIR}/packaging_train_results
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow

${INPUT_FILE}               ${RES_DIR}/request.json
${DEFAULT_RESULT_DIR}       ~/.odahuflow/local_packaging/training_output

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
...                 StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${DEFAULT_RESULT_DIR}
Suite Teardown      Run Keywords
...                 Remove Directory  ${RESULT_DIR}  recursive=True  AND
...                 Remove Directory  ${DEFAULT_RESULT_DIR}  recursive=True  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          cli  local  packaging
# Test Timeout        90 minutes

*** Keywords ***
Run Training with api server spec
    [Arguments]  ${command}  ${artifact path}
        ${result}  StrictShell  odahuflowctl --verbose ${command}

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

Run Packaging with local spec
    [Teardown]  Shell  docker stop -t 3 ${container_id.stdout}
    [Arguments]  ${options}
        ${pack_result}  StrictShell  odahuflowctl --verbose local packaging run ${options}

        Create File  ${RES_DIR}/pack_result.txt  ${pack_result.stdout}
        ${image_name}    Shell  tail -n 1 ${RES_DIR}/pack_result.txt | awk '{ print $4 }'
        Remove File  ${RES_DIR}/pack_result.txt

        StrictShell  docker images --all
        ${container_id}  StrictShell  docker run -d --rm -p 5002:5000 ${image_name.stdout}

        Sleep  5 sec
        Shell  docker container list -as -f id=${container_id.stdout}

        ${result_model}             StrictShell  odahuflowctl --verbose model invoke --url http://0:5002 --json-file ${RES_DIR}/request.json
        Should be equal as Strings  ${result_model.stdout}  ${MODEL_RESULT}

Try Run Training with api server spec
    [Arguments]  ${error}  ${command}
        ${result}  FailedShell  odahuflowctl --verbose ${command}
        should contain  ${result.stdout}  ${error}

Try Run Packaging with local spec
    [Arguments]  ${error}  ${options}
        ${result}  FailedShell  odahuflowctl --verbose local packaging run ${options}
        should contain  ${result.stdout}  ${error}

*** Test Cases ***
Run Valid Training with api server spec
    [Setup]     Run Keywords
    ...         Login to the api and edge  AND
    ...         StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/dir/training_cluster.json
    [Teardown]  StrictShell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/dir/training_cluster.json
    [Template]  Run Training with api server spec
    # auth data     id      file/dir        output
    local training run -f ${ARTIFACT_DIR}/dir/packaging --id wine-packaging --output ${RESULT_DIR}  ${RESULT_DIR}
    local training --url ${API_URL} --token "${AUTH_TOKEN}" run --id train-artifact-hardcoded  ${DEFAULT_RESULT_DIR}

Run Valid Packaging with local spec
    [Template]  Run Packaging with local spec
    # id	file/dir	artifact path	artifact name	package-targets
    --id pack-dir -d ${ARTIFACT_DIR}/dir --no-disable-package-targets
    --pack-id pack-file-image -f ${ARTIFACT_DIR}/file/packaging.yaml --artifact-path ${RESULT_DIR}/wine-name-1 --artifact-name wine-name-1
    --id pack-dir --manifest-dir ${ARTIFACT_DIR}/dir --disable-package-targets
    --pack-id pack-dir -d ${ARTIFACT_DIR}/dir --artifact-path ${RESULT_DIR}/wine-name-1 --disable-package-targets
    --id pack-file-image -f ${ARTIFACT_DIR}/file/packaging.yaml -a ${RESULT_DIR}/wine-name-1 --no-disable-package-targets
    --pack-id pack-dir --manifest-dir ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_RESULT_DIR}/simple-model
    --id pack-file-image --manifest-file ${ARTIFACT_DIR}/file/packaging.yaml -a simple-model --disable-package-targets

# negative tests
Try Run invalid Training with api server spec
    [Setup]  Shell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    [Template]  Try Run Training with api server spec
    # invalid credentials
    Error  local training --url "${API_URL}" --token "invalid" run -f ${ARTIFACT_DIR}/file/training.yaml --id train-artifact-hardcoded
    Error  local training --url "${API_URL}" --token "${EMPTY}" run -f ${ARTIFACT_DIR}/file/training.yaml --id train-artifact-hardcoded
    Error  local training --url "invalid" --token "${AUTH_TOKEN}" run -f ${ARTIFACT_DIR}/file/training.yaml --id train-artifact-hardcoded
    Error  local training --url "${EMPTY}" --token "${AUTH_TOKEN}" run -f ${ARTIFACT_DIR}/file/training.yaml --id train-artifact-hardcoded

Try Run invalid Packaging with local spec
    [Template]  Try Run Packaging with local spec
    # missing required option
    Error  -d ${ARTIFACT_DIR}/dir
    Error  --pack-id pack-file-image --artifact-path ${RESULT_DIR}/wine-name-1 --artifact-name wine-name-1
    # incompatible options
    Error  --id pack-dir -f ${ARTIFACT_DIR}/file/packaging.yaml --manifest-dir ${ARTIFACT_DIR}/dir
    Error  --id pack-dir -d ${ARTIFACT_DIR}/file --manifest-dir ${ARTIFACT_DIR}/dir --disable-package-targets --no-disable-package-targets
    Error  --id pack-dir -d ${ARTIFACT_DIR}/dir --disable-package-targets --no-disable-package-targets
    # not valid value for option
    # for file & dir options
    Error  --pack-id pack-dir --manifest-file ${ARTIFACT_DIR}/dir --artifact-path ${RESULT_DIR}/wine-name-1 --disable-package-targets
    Error  --id pack-file-image -d ${ARTIFACT_DIR}/file/packaging.yaml -a ${RESULT_DIR}/wine-name-1 --no-disable-package-targets
    # no training either locally or on the server
    Error  --id not-existing-packaging --manifest-file ${ARTIFACT_DIR}/file/packaging.yaml -a simple-model --disable-package-targets
