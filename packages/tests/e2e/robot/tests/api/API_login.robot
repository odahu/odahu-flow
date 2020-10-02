*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_login

${invalid_url}          https://invalid-url
${invalid_token}        _not-valid-token
${invalid_id}           12invalid-id
${invalid_secret}       234-invalid-client-secret
${invalid_issuer}       https://-invalid-issuer

*** Settings ***
Documentation       API for login
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.Login
Suite Setup         Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
Suite Teardown      Remove file  ${LOCAL_CONFIG}
Force Tags          api  sdk  security  login
Test Timeout        1 minute

*** Keywords ***
Verify login
    [Arguments]             &{keyword arguments}
    Call API  config get    &{keyword arguments}

Try login
    [Arguments]             ${error}  &{keyword arguments}
    Call API and get Error  ${error}  config get  &{keyword arguments}

*** Test Cases ***
Verify login with valid credentials
    [Template]  Verify login
    base_url=${API_URL}  token=${AUTH_TOKEN}
    base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}

Try login with invalid credentials
    [Template]  Try login
    ${IncorrectToken}  base_url=${API_URL}  token=${EMPTY}
    ${IncorrectToken}  base_url=${API_URL}  token=${invalid_token}

    ${IncorrectToken}  base_url=${API_URL}  client_id=${EMPTY}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}
    ${IncorrectToken}  base_url=${API_URL}  client_id=${invalid_id}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}
    ${IncorrectToken}  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${EMPTY}  issuer_url=${ISSUER}
    ${IncorrectToken}  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${invalid_secret}  issuer_url=${ISSUER}
    ${IncorrectToken}  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${EMPTY}
    ${IncorrectToken}  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${invalid_issuer}

Try login with empty url and valid token
    ${error}                    format string  ${APIConnectionException}  base url=${EMPTY}
    Call API and get Error      ${error}  config get  base_url=${EMPTY}  token=${AUTH_TOKEN}

Try login with invalid url and valid token
    [Timeout]                   20 minutes
    ${error}                    format string  ${APIConnectionException}  base url=${invalid_url}
    Call API and get Error      ${error}  config get  base_url=${invalid_url}  token=${AUTH_TOKEN}
