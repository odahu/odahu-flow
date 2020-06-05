*** Settings ***
Force Tags          testing
Documentation       API keywords
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Library             odahuflow.robot.libraries.sdk_wrapper.Connection
# Suite Setup         Run Keywords  Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
# ...                               Login to the api and edge

Resource            variables.robot
Variables           ../load_variables_from_profiles.py
Library             String
Library             OperatingSystem
Library             Collections
Library             DateTime
Library             odahuflow.robot.libraries.k8s.K8s  ${ODAHUFLOW_NAMESPACE}
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.process.Process

*** Keywords ***
Call API
    [Arguments]         ${keyword}
    ${result}=          ${keyword}
    Log                 stdout = ${result.stdout}
    Log                 stderr = ${result.stderr}
    [Return]            ${result}

Strict Call API
    [Arguments]     ${command}
    ${res}=         CallAPI  ${command}
                    Should Be Equal  ${res.rc}  ${0}
    [Return]        ${res}

Fail Call API
    [Arguments]     ${command}
    ${res}=         CallAPI  ${command}
                    Should Not Be Equal  ${res.rc}  ${0}
    [Return]        ${res
