*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_status_codes_400-401-403
${RES_DIR}              ${CURDIR}/resources


*** Settings ***
Documentation       tests for API status codes 400, 401, 403
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.Configuration
Library             odahuflow.robot.libraries.sdk_wrapper.Connection
Library             odahuflow.robot.libraries.sdk_wrapper.Toolchain
Library             odahuflow.robot.libraries.sdk_wrapper.Packager
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Library             odahuflow.robot.libraries.sdk_wrapper.ModelDeployment
Library             odahuflow.robot.libraries.sdk_wrapper.ModelRoute
Library             odahuflow.robot.libraries.sdk_wrapper.Model
Suite Setup         Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
Suite Teardown      Remove file  ${LOCAL_CONFIG}
Force Tags          api  sdk  negative  test
Test Timeout        1 minute

*** Keywords ***
Try Call API
    [Arguments]  ${error}  ${command}  @{options}
    Call API and get Error  ${error}  ${command}  @{options}

Try Call API and continue on Failure
    [Arguments]     ${error}  ${command}  @{options}
    ${result}       Call API and continue on Failure  ${command}  @{options}
    should contain  ${result}  ${error}

*** Test Cases ***
Status Code 400 - Bad Request
    [Template]  Try Call API
    # connection
#    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; unknown type: . Supported types: [s3 gcs azureblob git docker ecr] (status: 400)
#    ...  connection post  ${RES_DIR}/connection/invalid/no_type
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: ID is not valid (status: 400)
    ...  connection put   ${RES_DIR}/connection/invalid/conn_invalid_id.yaml
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: s3 type requires that keyID and keySecret parameters must be non-empty (status: 400)
    ...  connection post  ${RES_DIR}/connection/invalid/s3_no_required_parameters.yaml
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: s3 type requires that keyID and keySecret parameters must be non-empty (status: 400)
    ...  connection put   ${RES_DIR}/connection/invalid/s3_no_required_parameters.yaml
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; gcs type requires that keySecret parameter must be non-empty (status: 400)
    ...  connection post  ${RES_DIR}/connection/invalid/gcs_no_required_parameters
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; gcs type requires that keySecret parameter must be non-empty (status: 400)
    ...  connection put   ${RES_DIR}/connection/invalid/gcs_no_required_parameters
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; azureblob type requires that keySecret parameter containsHTTP endpoint with SAS Token (status: 400)
    ...  connection post  ${RES_DIR}/connection/invalid/azureblob_no_required_parameters.json
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; azureblob type requires that keySecret parameter containsHTTP endpoint with SAS Token (status: 400)
    ...  connection put   ${RES_DIR}/connection/invalid/azureblob_no_required_parameters.json
#    Error  connection post  ${RES_DIR}/connection/invalid/git_no_required_parameters.json
#    Error  connection put   ${RES_DIR}/connection/invalid/git_no_required_parameters.json
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty (status: 400)
    ...  connection post  ${RES_DIR}/connection/invalid/docker_no_required_parameters.yaml
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty (status: 400)
    ...  connection put   ${RES_DIR}/connection/invalid/docker_no_required_parameters.yaml
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; not valid uri for ecr type: docker-credential-ecr-login can only be used with Amazon Elastic Container Registry.; ecr type requires that keyID and keySecret parameters must be non-empty (status: 400)
    ...  connection post  ${RES_DIR}/connection/invalid/ecr_no_required_parameters.json
    WrongHttpStatusCode: Got error from server: Validation of connection is failed: the uri parameter is empty; not valid uri for ecr type: docker-credential-ecr-login can only be used with Amazon Elastic Container Registry.; ecr type requires that keyID and keySecret parameters must be non-empty (status: 400)
    ...  connection put   ${RES_DIR}/connection/invalid/ecr_no_required_parameters.json
    # toolchains
    WrongHttpStatusCode: Got error from server: Validation of toolchain integration is failed: entrypoint must be no empty; defaultImage must be no empty (status: 400)
    ...  toolchain post  ${RES_DIR}/toolchain/invalid/toolchain_no_required_parameters.json
    WrongHttpStatusCode: Got error from server: Validation of toolchain integration is failed: entrypoint must be no empty (status: 400)
    ...  toolchain put  ${RES_DIR}/toolchain/invalid/toolchain_no_required_parameters.yaml
    # packagers
    WrongHttpStatusCode: Got error from server: Validation of packaging integration is failed: entrypoint must be nonempty; default image must be nonempty (status: 400)
    ...  packager post  ${RES_DIR}/packager/invalid/cli_no_required_params.yaml
    WrongHttpStatusCode: Got error from server: Validation of packaging integration is failed: entrypoint must be nonempty; default image must be nonempty (status: 400)
    ...  packager put  ${RES_DIR}/packager/invalid/cli_no_required_params.yaml
    WrongHttpStatusCode: Got error from server: Validation of packaging integration is failed: entrypoint must be nonempty; default image must be nonempty (status: 400)
    ...  packager post  ${RES_DIR}/packager/invalid/rest_no_required_params.json
    WrongHttpStatusCode: Got error from server: Validation of packaging integration is failed: entrypoint must be nonempty; default image must be nonempty (status: 400)
    ...  packager put  ${RES_DIR}/packager/invalid/rest_no_required_params.json
#     # model training
#     Error  training post
#     Error  training put
#     # model packaging
#     Error  packaging post
#     Error  packaging put
#     # model deployment
#     Error  deployment post
#     Error  deployment put

## also create 401, 403 for different user types (data-scientist, viewer, admin)
#
#Status Code 401 - Unathorized
#    [Template]  Template Error Keyword
#    # config
#    config get
#    # connection
#    connection get
#    connection get id
#    connection get id decrypted
#    connection post
#    connection put
#    connection delete
#    # toolchains
#    toolchain get
#    toolchain get id
#    toolchain post
#    toolchain put
#    toolchain delete
#    # packagers
#    packager get
#    packager get id
#    packager post
#    packager put
#    packager delete
#    # training
#    training get
#    training get id
#    training get log
#    training post
#    training put
#    training delete
#    # packaging
#    packaging get
#    packaging get id
#    packaging get log
#    packaging post
#    packaging put
#    packaging delete
#    # deployment
#    deployment get
#    deployment get id
#    deployment post
#    deployment put
#    deployment delete
#    # route
#    route get
#    route get id
#    route post
#    route put
#    route delete
#    # model
#    model get
#    model post
#
#Status Code 403 - Forbidden
#    [Template]  Template Error Keyword
#    # config
#    config get
#    # connection
#    connection get
#    connection get id
#    connection get id decrypted
#    connection post
#    connection put
#    connection delete
#    # toolchains
#    toolchain get
#    toolchain get id
#    toolchain post
#    toolchain put
#    toolchain delete
#    # packagers
#    packager get
#    packager get id
#    packager post
#    packager put
#    packager delete
#    # training
#    training get
#    training get id
#    training get log
#    training post
#    training put
#    training delete
#    # packaging
#    packaging get
#    packaging get id
#    packaging get log
#    packaging post
#    packaging put
#    packaging delete
#    # deployment
#    deployment get
#    deployment get id
#    deployment post
#    deployment put
#    deployment delete
#    # route
#    route get
#    route get id
#    route post
#    route put
#    route delete
#    # model
#    model get
#    model post
#