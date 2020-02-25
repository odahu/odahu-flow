*** Variables ***
${LOCAL_CONFIG}         odahuflow/config_common_config
${TEST_VALUE}           test

*** Settings ***
Documentation       Login cli config command
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Force Tags          cli  config  common
Suite Setup         Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
Suite Teardown      Remove file  ${LOCAL_CONFIG}

*** Test Cases ***
All parameters
    ${res}=  Shell  odahuflowctl --verbose config all
    Should be equal  ${res.rc}  ${0}
    should contain  ${res.stdout}  API_URL
    should contain  ${res.stdout}  API_TOKEN

Config path
    ${res}=  Shell  odahuflowctl --verbose config path
    Should be equal  ${res.rc}  ${0}
    should contain  ${res.stdout}  ${LOCAL_CONFIG}

Set config value
    ${res}=  Shell  odahuflowctl --verbose config get ODAHUFLOWCTL_OAUTH_AUTH_URL
    Should be equal  ${res.rc}  ${0}
    should not contain  ${res.stdout}  ${TEST_VALUE}

    ${res}=  Shell  odahuflowctl --verbose config set ODAHUFLOWCTL_OAUTH_AUTH_URL test
    Should be equal  ${res.rc}  ${0}

    ${res}=  Shell  odahuflowctl --verbose config get ODAHUFLOWCTL_OAUTH_AUTH_URL
    Should be equal  ${res.rc}  ${0}
    should contain  ${res.stdout}  ${TEST_VALUE}
