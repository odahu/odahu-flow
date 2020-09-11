*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_authn

${invalid_url}          -invalid-url
${invalid_token}        _not-valid-token
${invalid_id}           12invalid-id
${invalid_secret}       234-invalid-client-secret
${invalid_issuer}       https://-invalid-issuer

${APIConnectionException}   APIConnectionException: Can not reach {base url}
${IncorrectToken}           IncorrectAuthorizationToken: Refresh token is not correct.\nPlease login again

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

*** Keywords ***
Format APIConnectionException
    [Arguments]         ${base url}
    ${error output}     format string  ${APIConnectionException}  base url=${base url}
    [return]            ${error output}

Verify login
    [Arguments]             &{keyword arguments}
    Call API  config get    &{keyword arguments}

Try login
    [Arguments]             ${error code}  &{keyword arguments}
    Call API and get Error  ${error code}  config get  &{keyword arguments}

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

Try login with empty url
    ${error}        Format APIConnectionException   base url=${EMPTY}
    Call API and get Error      ${error}  config get  base_url=${EMPTY}  token=${AUTH_TOKEN}

Try login with invalid url and valid token
    ${error}        Format APIConnectionException   base url=${invalid_url}
    Call API and get Error      ${error}  config get  base_url=${invalid_url}  token=${AUTH_TOKEN}

Try login with valid credentials but empty client id
    Call API and get Error      ${IncorrectToken}  config get  base_url=${API_URL}  client_id=${EMPTY}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}

Try login with valid credentials but invalid client id
    Call API and get Error      ${IncorrectToken}  config get  base_url=${API_URL}  client_id=${invalid_id}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${ISSUER}

Try login with valid credentials but empty client secret
    Call API and get Error      ${IncorrectToken}  config get  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${EMPTY}  issuer_url=${ISSUER}

Try login with valid credentials but invalid client secret
    Call API and get Error      ${IncorrectToken}  config get  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${invalid_secret}  issuer_url=${ISSUER}

Try login with valid credentials but empty issuer
    Call API and get Error      ${IncorrectToken}  config get  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${EMPTY}

Try login with valid credentials but invalid issuer
    Call API and get Error      ${IncorrectToken}  config get  base_url=${API_URL}  client_id=${SA_CLIENT_ID}  client_secret=${SA_CLIENT_SECRET}  issuer_url=${invalid_issuer}
