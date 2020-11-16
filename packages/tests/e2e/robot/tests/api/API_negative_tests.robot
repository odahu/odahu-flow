*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_status_codes_400-401-403

*** Settings ***
Documentation       tests for API status codes 400, 401, 403
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.Login
Suite Setup         Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
Suite Teardown      Remove file  ${LOCAL_CONFIG}
Force Tags          api  sdk  negative
Test Timeout        1 minute

*** Keywords ***
Template Error Keyword
    [Arguments]  ${command}  @{options}
    Call API and get Error  ${expected_error}  ${command}  @{options}

*** Test Cases ***
Status Code 400 - Bad Request
    [Template]  Template Error Keyword
    # connection
    connection post
    connection put
    connection post
    connection put
    connection post
    connection put
    connection post
    connection put
    connection post
    connection put
    connection post
    connection put
    # toolchains
    toolchain post
    toolchain put
    # packagers
    packager post
    packager put
    # model training
    training post
    training put
    # model packaging
    packaging post
    packaging put
    # model deployment
    deployment post
    deployment put

# also create 401, 403 for different user types (data-scientist, viewer, admin)

Status Code 401 - Unathorized
    [Template]  Template Error Keyword
    # config
    config get
    # connection
    connection get
    connection get id
    connection get id decrypted
    connection post
    connection put
    connection delete
    # toolchains
    toolchain get
    toolchain get id
    toolchain post
    toolchain put
    toolchain delete
    # packagers
    packager get
    packager get id
    packager post
    packager put
    packager delete
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

Status Code 403 - Forbidden
    [Template]  Template Error Keyword
    # config
    config get
    # connection
    connection get
    connection get id
    connection get id decrypted
    connection post
    connection put
    connection delete
    # toolchains
    toolchain get
    toolchain get id
    toolchain post
    toolchain put
    toolchain delete
    # packagers
    packager get
    packager get id
    packager post
    packager put
    packager delete
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
