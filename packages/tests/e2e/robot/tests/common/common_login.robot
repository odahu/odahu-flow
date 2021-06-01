*** Variables ***
${LOCAL_CONFIG}         odahuflow/config_common_login

*** Settings ***
Documentation       Login cli command
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             Collections
Resource            ../../resources/keywords.robot
Force Tags          cli  common  security
Suite Setup         Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
Suite Teardown      Remove file  ${LOCAL_CONFIG}

*** Test Cases ***
Verifying of a valid token
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"
    Should be equal  ${res.rc}  ${0}
    should contain  ${res.stdout}  Success! Credentials have been saved

Verifying of a not valid token
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --token "not-valid-token"
    Should not be equal  ${res.rc}  ${0}
    should contain  ${res.stderr}  Credentials are not correct

Verifying misconfused cli parameters is not permitted
    [Documentation]  Intent of user how to login should be clear. Or using client credentials or token
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --token "not-valid-token" --client-id "client-id"
    Should not be equal  ${res.rc}  ${0}
    should contain  ${res.stderr}  You should use either --token or --client-id/--client-secret to login

Verifying of a valid client credentials flow
    [Documentation]  User should be able to login via ctl as service account using client_credentials oauth2 flow
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --client-id ${SA_CLIENT_ID} --client-secret ${SA_CLIENT_SECRET} --issuer ${ISSUER}
    Should be equal  ${res.rc}  ${0}
    should contain  ${res.stdout}  Success! Credentials have been saved

User try to make client_credentials login without passing issuer param
    [Documentation]
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --client-id ${SA_CLIENT_ID} --client-secret ${SA_CLIENT_SECRET}
    Should not be equal  ${res.rc}  ${0}
    should contain  ${res.stderr}  You must provide --issuer parameter to do client_credentials login

User try to make client_credentials login without passing client_secret param
    [Documentation]
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --client-id ${SA_CLIENT_ID} --issuer ${ISSUER}
    Should not be equal  ${res.rc}  ${0}
    should contain  ${res.stderr}  You must pass both client_id and client_secret to client_credentials login
