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
Template Keyword
    [Arguments]  ${command}  @{options}
    ${command}  @{options}

*** Test Cases ***
Status Code 400
    [Template]  Template Keyword
    config get

Status Code 401
    [Template]  Template Keyword
    config get

Status Code 403
    [Template]  Template Keyword
    config get
