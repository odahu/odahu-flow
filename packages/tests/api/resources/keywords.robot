*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_cli
${MD_SIMPLE_MODEL}     simple-model-cli

*** Settings ***
Documentation       API keywords
Resource            ../../e2e/robot/resources/keywords.robot
Resource            ../../e2e/robot/resources/variables.robot
Variables           ../e2e/robot/load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Library             odahuflow.robot.libraries.sdkWrapperForApi.Connection
Suite Setup         Run Keywords  Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                               Login to the api and edge

*** Keywords ***
CallAPI
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
    [Return]        ${res}


*** Test Cases ***
testing of API for connection
    CallAPI               connection get
    Strict Call API       connection get
