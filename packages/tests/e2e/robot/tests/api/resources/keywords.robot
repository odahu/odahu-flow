*** Settings ***
Documentation       API keywords
Resource            ../../../resources/keywords.robot
Resource            ../../../resources/variables.robot
Variables           ../../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Library             odahuflow.robot.libraries.sdk_wrapper.Connection

*** Keywords ***
Call API
    [Arguments]         ${keyword}  @{arguments}
    ${result}=          Run Keyword  ${keyword}    @{arguments}
    Log                 ${result}
    [Return]            ${result}

Strict Call API
    [Arguments]     ${command}
    ${result}=      Call API  ${command}
                    # Should Be Equal  ${res.rc}  ${0}
    [Return]        ${result}

Fail Call API
    [Arguments]     ${command}
    ${result}=      Call API  ${command}
                    # Should Not Be Equal  ${res.rc}  ${0}
    [Return]        ${result}
