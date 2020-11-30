*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_status_codes_400-401-403
${RES_DIR}              ${CURDIR}/resources


*** Settings ***
Documentation       tests for API status codes 400, 401, 403
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
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
    [Arguments]  ${format_string}  ${error}  ${command}  @{options}
    ${error}        format string   ${format_string}  ${error}
    Call API and get Error  ${error}  ${command}  @{options}

Try Call API and continue on Failure
    [Arguments]     ${format_string}  ${error}  ${command}  @{options}
    ${error}        format string   ${format_string}  ${error}
    ${result}       Call API and continue on Failure  ${command}  @{options}
    should contain  ${result}  ${error}

*** Test Cases ***
Status Code 400 - Bad Request
    [Template]  Try Call API
    # connection
#    ${400 BadRequest Template}  ${FailedConn} ${unknown type}
#    ...  connection post  ${RES_DIR}/connection/invalid/no_type
    ${400 BadRequest Template}  ${FailedConn} ${invalid_id}
    ...  connection put   ${RES_DIR}/connection/invalid/conn_invalid_id.yaml
    ${400 BadRequest Template}  ${FailedConn} ${s3_empty_keyID_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/s3_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${s3_empty_keyID_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/s3_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${gcs_empty_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/gcs_no_required_parameters
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${gcs_empty_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/gcs_no_required_parameters
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${azureblob_req_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/azureblob_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${azureblob_req_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/azureblob_no_required_parameters.json
#    ${400 BadRequest Template}  Error
#     ...  connection post  ${RES_DIR}/connection/invalid/git_no_required_parameters.json
#    ${400 BadRequest Template}  Error
#     ...  connection put   ${RES_DIR}/connection/invalid/git_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}
    ...  connection post  ${RES_DIR}/connection/invalid/docker_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}
    ...  connection put   ${RES_DIR}/connection/invalid/docker_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${ecr_invalid_uri}; ${ecr_empty_keyID_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/ecr_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${ecr_invalid_uri}; ${ecr_empty_keyID_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/ecr_no_required_parameters.json
    # toolchains
    ${400 BadRequest Template}  ${FailedTI} ${TI_empty_entrypoint}; ${TI_empty_defaultImage}
    ...  toolchain post  ${RES_DIR}/toolchain/invalid/toolchain_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedTI} ${TI_empty_entrypoint}
    ...  toolchain put  ${RES_DIR}/toolchain/invalid/toolchain_no_required_parameters.yaml
    # packagers
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager post  ${RES_DIR}/packager/invalid/cli_no_required_params.yaml
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager put  ${RES_DIR}/packager/invalid/cli_no_required_params.yaml
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager post  ${RES_DIR}/packager/invalid/rest_no_required_params.json
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager put  ${RES_DIR}/packager/invalid/rest_no_required_params.json
    # model training
    ${400 BadRequest Template}  ${FailedTrain} ${empty_model_name}; ${empty_model_version}; ${empty_VCS}; ${empty_toolchain}
    ...  training post  ${RES_DIR}/training_packaging/invalid/training_no_required_params.yaml
    ${400 BadRequest Template}  ${FailedTrain} ${empty_model_name}; ${empty_model_version}; ${empty_VCS}; ${empty_toolchain}
    ...  training put  ${RES_DIR}/training_packaging/invalid/training_no_required_params.yaml
    # model packaging
    ${400 BadRequest Template}  ${FailedPack} entity "" is not found; ${empty_artifactName}; ${empty_integrationName}
    ...  packaging post  ${RES_DIR}/training_packaging/invalid/packaging_no_required_params.json
    ${400 BadRequest Template}  ${FailedPack} entity "" is not found; ${empty_artifactName}; ${empty_integrationName}
    ...  packaging put  ${RES_DIR}/training_packaging/invalid/packaging_no_required_params.json
    # model deployment
    ${400 BadRequest Template}  ${max_smaller_min_replicas}; ${empty_image}
    ...  deployment post  ${RES_DIR}/deploy_route_model/invalid/deployment_empty_required_params.json
    ${400 BadRequest Template}  ${max_smaller_min_replicas}; ${empty_image}
    ...  deployment put  ${RES_DIR}/deploy_route_model/invalid/deployment_empty_required_params.json
    ${400 BadRequest Template}  ${positive_livenessProbe}; ${positive_readinessProbe}; ${max_smaller_min_replicas}; ${min_num_of_max_replicas}; ${min_num_of_min_replicas}
    ...  deployment post  ${RES_DIR}/deploy_route_model/invalid/deployment_validation_checks.yaml
    ${400 BadRequest Template}  ${positive_livenessProbe}; ${positive_readinessProbe}; ${max_smaller_min_replicas}; ${min_num_of_max_replicas}; ${min_num_of_min_replicas}
    ...  deployment put  ${RES_DIR}/deploy_route_model/invalid/deployment_validation_checks.yaml

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
#
## also create 403 for different user types (data-scientist, viewer, admin)
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