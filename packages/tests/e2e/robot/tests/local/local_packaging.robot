*** Variables ***
${RES_DIR}                  ${CURDIR}/resources
${ARTIFACT_DIR}             ${RES_DIR}/artifacts/odahuflow
${RESULT_DIR}               ${CURDIR}/packaging_train_results

${INPUT_FILE}               ${RES_DIR}/request.json
${DEFAULT_RESULT_DIR}       ~/.odahuflow/local_packaging/training_output

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
...                 Login to the api and edge  AND
...                 StrictShell  odahuflowctl --verbose config set LOCAL_MODEL_OUTPUT_DIR ${DEFAULT_RESULT_DIR}
Suite Teardown      Run Keywords
...                 Remove Directory  ${RESULT_DIR}  recursive=True  AND
...                 Remove Directory  ${DEFAULT_RESULT_DIR}  recursive=True  AND
...                 Remove File  ${LOCAL_CONFIG}  AND
...                 StrictShell  odahuflowctl logout
Force Tags          cli  local  packaging
Test Timeout        180 minutes

*** Keywords ***
Run Training with server spec
    [Arguments]  ${command}  ${artifact path}
        ${result}  StrictShell  odahuflowctl --verbose ${command}

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

Run Packaging with local spec
    [Teardown]  Shell  docker rm -f ${container_id.stdout}
    [Arguments]  ${options}
        ${pack_result}  StrictShell  odahuflowctl --verbose local packaging run ${options}

        Create File  ${RESULT_DIR}/pack_result.txt  ${pack_result.stdout}
        ${image_name}    StrictShell  tail -n 1 ${RESULT_DIR}/pack_result.txt | awk '{ print $4 }'
        Remove File  ${RESULT_DIR}/pack_result.txt

        StrictShell  docker images --all
        ${container_id}  StrictShell  docker run -d --rm -p 5002:5000 ${image_name.stdout}

        Sleep  5 sec
        StrictShell  docker container list -as -f id=${container_id.stdout}

        ${MODEL_HOST}    Get local model host
        ${result_model}  StrictShell  odahuflowctl --verbose model invoke --url ${MODEL_HOST}:5002 --json-file ${RES_DIR}/request.json
        Should be equal as Strings  ${result_model.stdout}  ${WINE_MODEL_RESULT}

Try Run Training with server spec
    [Arguments]  ${error}  ${command}
        ${result}  FailedShell  odahuflowctl --verbose ${command}
        should contain  ${result.stdout}  ${error}

Try Run Packaging with local spec
    [Arguments]  ${error}  ${options}
        ${result}  FailedShell  odahuflowctl --verbose local packaging ${options}
        ${result}  Catenate  ${result.stdout}  ${result.stderr}
        should contain  ${result}  ${error}

*** Test Cases ***
Run Valid Training with server spec
    [Setup]     StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/dir/training_cluster.json
    [Teardown]  StrictShell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/dir/training_cluster.json
    [Template]  Run Training with server spec
    # auth data     id      file/dir        output
    local train run -f ${ARTIFACT_DIR}/dir/packaging --id local-dir-cluster-artifact-template --output ${RESULT_DIR}  ${RESULT_DIR}
    local training --url ${API_URL} --token "${AUTH_TOKEN}" run --id local-dir-cluster-artifact-hardcoded  ${DEFAULT_RESULT_DIR}

Try Run and Fail Packaging with invalid credentials
    [Tags]   negative
    [Setup]  StrictShell  odahuflowctl logout
    [Teardown]  Login to the api and edge
    [Template]  Try Run Packaging with server spec
    ${INVALID_CREDENTIALS_ERROR}    local pack --url ${API_URL} --token "invalid" run -f ${ARTIFACT_DIR}/file/packaging.yaml --id not-exist
    ${MISSED_CREDENTIALS_ERROR}     local pack --url ${API_URL} --token "${EMPTY}" run -f ${ARTIFACT_DIR}/file/packaging.yaml --id not-exist
    ${INVALID_URL_ERROR}            local pack --url "invalid" --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/packaging.yaml --id not-exist
    ${INVALID_URL_ERROR}            local pack --url "${EMPTY}" --token ${AUTH_TOKEN} run -f ${ARTIFACT_DIR}/file/packaging.yaml --id not-exist

Try Run and Fail invalid Packaging
    [Tags]  negative
    [Template]  Try Run Packaging with local spec
    [Setup]     StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Teardown]  Shell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    # missing required option
    Error: Missing option '--pack-id' / '--id'.
    ...  --manifest-file ${ARTIFACT_DIR}/dir --artifact-path ${RESULT_DIR}/wine-name-1
    # not valid value for option
    # for file & dir options
    Error: [Errno 21] Is a directory: '${ARTIFACT_DIR}/dir'
    ...  --pack-id local-dir-spec-targets --manifest-file ${ARTIFACT_DIR}/dir --artifact-path ${RESULT_DIR}/wine-name-1
    Error: ${ARTIFACT_DIR}/file/packaging.yaml is not a directory
    ...  --id local-file-image-template -d ${ARTIFACT_DIR}/file/packaging.yaml -a ${RESULT_DIR}/wine-name-1
    Error: Resource file '${ARTIFACT_DIR}/file/not-existing.yaml' not found
    ...  --id local-file-image-template -f ${ARTIFACT_DIR}/file/not-existing.yaml -a ${RESULT_DIR}/wine-name-1
    Error: [Errno 2] No such file or directory: '${RESULT_DIR}/not-existing/mp.json'
    ...  --id local-file-image-template -f ${ARTIFACT_DIR}/file/packaging.yaml -a ${RESULT_DIR}/not-existing
    # no training either locally or on the server
    Error: Got error from server: entity "not-existing-packaging" is not found (status: 404)
    ...  --id not-existing-packaging --manifest-file ${ARTIFACT_DIR}/file/packaging.yaml -a simple-model
    # manifest on cluster but disabled target docker-pull (image not pulled locally)
    Exception: unauthorized: You don't have the needed permissions to perform this operation, and you may have invalid credentials.
    ...  --id local-cluster -a wine-name-1 --artifact-path ${RESULT_DIR} --disable-package-targets
    Exception: unauthorized: You don't have the needed permissions to perform this operation, and you may have invalid credentials.
    ...  --id local-cluster -a wine-name-1 --artifact-path ${RESULT_DIR}
    Exception: unauthorized: You don't have the needed permissions to perform this operation, and you may have invalid credentials.
    ...  --id local-cluster --no-disable-package-targets --disable-target docker-pull

Run Valid Packaging with local & cluster specs
    [Setup]     StrictShell  odahuflowctl --verbose bulk apply ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Teardown]  Shell  odahuflowctl --verbose bulk delete ${ARTIFACT_DIR}/file/packaging_cluster.yaml
    [Template]  Run Packaging with local spec
    # id	file/dir	artifact path	artifact name	package-targets
    # local
    run --pack-id local-file-image-template -f ${ARTIFACT_DIR}/file/packaging.yaml --artifact-path ${RESULT_DIR} --artifact-name wine-name-1
    run --id local-dir-spec-targets --manifest-dir ${ARTIFACT_DIR}/dir --disable-package-targets
    run --pack-id local-dir-spec-targets -d ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_RESULT_DIR} --disable-package-targets
    run --pack-id local-dir-spec-targets --manifest-dir ${ARTIFACT_DIR}/dir --artifact-path ${DEFAULT_RESULT_DIR}
    run --id local-file-image-template --manifest-file ${ARTIFACT_DIR}/file/packaging.yaml -a simple-model --disable-package-targets
    # cluster
    run --id local-dir-spec-targets -d ${ARTIFACT_DIR}/dir --no-disable-package-targets
    run --id local-cluster-spec-targets -f ${ARTIFACT_DIR}/file/packaging.yaml -a ${RESULT_DIR}/wine-name-1 --no-disable-package-targets  # watch for this
    run --id local-cluster-spec-targets --no-disable-package-targets --disable-target docker-push --disable-target not-existing
    run -f ${ARTIFACT_DIR}/dir/packaging --id local-cluster --artifact-path ${RESULT_DIR} --artifact-name wine-dir-1.0 --no-disable-package-targets
    --url ${API_URL} --token ${AUTH_TOKEN} run --id local-cluster-spec-targets --artifact-name simple-model --no-disable-package-targets --disable-target docker-push
