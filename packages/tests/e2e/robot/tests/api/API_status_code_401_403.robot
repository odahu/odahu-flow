*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_401-403
${RES_DIR}              ${CURDIR}/resources
${invalid_token}        not-valid-token


*** Settings ***
Documentation       tests for API status codes 401, 403
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
# Suite Teardown      Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  negative
Test Timeout        1 minute

*** Test Cases ***
Status Code 401 - Unathorized
    [Template]  Try Call API - Unathorized
    # config
    config get
    # connection
    connection get
    connection get id  ${VCS_CONNECTION}
    connection get id decrypted  ${VCS_CONNECTION}
    connection post  ${RES_DIR}/connection/valid/docker_connection_create.json
    connection put  ${RES_DIR}/connection/valid/git_connection_update.yaml
    connection delete  not-exist
    # toolchains
    toolchain get
    toolchain get id  ${TOOLCHAIN_INTEGRATION}
    toolchain post  ${RES_DIR}/toolchain/valid/mlflow_create.yaml
    toolchain put  ${RES_DIR}/toolchain/valid/mlflow_update.json
    toolchain delete  not-exist
    # packagers
    packager get
    packager get id  ${PI_REST}
    packager post  ${RES_DIR}/packager/valid/docker_rest_create.json
    packager put  ${RES_DIR}/packager/valid/docker_rest_update.yaml
    packager delete  not-exist
    # training
    training get
    training get id
    training get log
    training post
    training put
    training delete
    # packaging
    packaging get
    packaging get id
    packaging get log
    packaging post
    packaging put
    packaging delete
    # deployment
    deployment get
    deployment get id
    deployment post
    deployment put
    deployment delete
    # route
    route get
    route get id
    route post
    route put
    route delete
    # model
    model get
    model post

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