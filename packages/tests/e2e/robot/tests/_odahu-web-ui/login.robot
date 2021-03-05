*** Variables ***
${RES_DIR}          ${CURDIR}/resources

*** Settings ***
Documentation   testing login to ODAHU WebUI
Resource        ${RES_DIR}/common.robot
Test Setup      Begin Web test
Test Teardown   End Web test
Force Tags      web-ui  login

*** Test Cases ***
Should be able to "Log In" with valid credentials
    Login to ODAHU WebUI  ${UI_USERNAME}  ${UI_PASSWORD}

Should be able to "Log Out" and be redirected to keycloak
    Login to ODAHU WebUI  ${UI_USERNAME}  ${UI_PASSWORD}
    Log Out from ODAHU WebUI

Should see error when try to Log In with invalid credentials
    [Template]  Fail Login with invalid credentials
    invalid         ${UI_PASSWORD}
    ${UI_USERNAME}  some-password
    user-name       p@ssw0rd
    ${UI_USERNAME}  ${EMPTY}
    ${EMPTY}        ${UI_PASSWORD}
    ${EMPTY}        ${EMPTY}
