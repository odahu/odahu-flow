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
    [Teardown]  run keywords
    ...         StrictShell  odahuflowctl logout  AND
    ...         reload config
    base_url=${API_URL}  token=${AUTH_TOKEN}
    base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}

Try login with invalid credentials
    [Template]  Try login
    ${MissedToken}              base_url=${API_URL}  token=${EMPTY}
    ${IncorrectToken}           base_url=${API_URL}  token=${invalid_token}

    ${MissedToken}              base_url=${API_URL}  client_id=${EMPTY}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}
    ${IncorrectCredentials}     base_url=${API_URL}  client_id=${invalid_id}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}
    ${MissedToken}              base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${EMPTY}  issuer_url=${ISSUER}
    ${IncorrectCredentials}     base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${invalid_secret}  issuer_url=${ISSUER}
    ${MissedToken}              base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${EMPTY}
    ${IncorrectCredentials}     base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${invalid_issuer}

Try login with empty url and valid token
    ${default_url}              StrictShell  odahuflowctl config get API_URL | grep default | awk '{print $2}' | xargs
    ${error}                    format string  ${APIConnectionException}  base url=${default_url.stdout}
    Call API and get Error      ${error}  config get  base_url=${EMPTY}  token=${AUTH_TOKEN}

Try login with invalid url and valid token
    [Timeout]                   10 minutes
    ${error}                    format string  ${APIConnectionException}  base url=${invalid_url}
    Call API and get Error      ${error}  config get  base_url=${invalid_url}  token=${AUTH_TOKEN}
