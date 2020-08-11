*** Settings ***
Documentation       API of configuration
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Configuration
Suite Setup         Run Keywords
...                 Login to the api and edge
Force Tags          api  sdk  configuration
Test Timeout        5 minutes


*** Test Cases ***
Get configuration
    [Documentation]  create git connection and check that one exists
    ${result}                           Call API  config get
    should not be equal as strings      ${result}  ${EMPTY}
