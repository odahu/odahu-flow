*** Variables ***
${LOCAL_CONFIG}     odahuflow/api_configuration

*** Settings ***
Documentation       API of configuration
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Configuration
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge
Suite Teardown      Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  configuration
Test Timeout        5 minutes


*** Test Cases ***
Get configuration
    ${result}                           Call API  config get
    should not be equal as strings      ${result}  ${EMPTY}
