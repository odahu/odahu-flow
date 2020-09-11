*** Variables ***
${LOCAL_CONFIG}     odahuflow/api_authn

*** Settings ***
Documentation       API for login
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Configuration
Suite Setup         Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
Suite Teardown      Remove file  ${LOCAL_CONFIG}
Force Tags          api  sdk  security  authn

*** Test Cases ***
Verify login with token
    ${result}           Call API  config get  base_url=${API_URL}  token=${AUTH_TOKEN}

Verify login with client id, client secret and issuer
    Log                 ${API_URL}
    Log                 ${SA_CLIENT_ID}
    Log                 ${SA_CLIENT_SECRET}
    Log                 ${ISSUER}
    ${result}           Call API  config get  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}

Try login with not valid token
    ${result}           Call API  config get  base_url=${API_URL}  token="not-valid-token"
