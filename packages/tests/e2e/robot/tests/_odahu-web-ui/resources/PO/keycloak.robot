*** Settings ***
Documentation   Page Object for KEYCLOAK page
Library         SeleniumLibrary  timeout=10s

*** Variables ***
${KEYCLOAK.LOG_IN_TITLE}          Log in to ODAHU Cluster
${KEYCLOAK.LOG_IN_HEADING_TEXT}   Username or email
${KEYCLOAK.LOG_IN_HEADING}        xpath://*[@id="kc-form-login"]/div[1]/label
${KEYCLOAK.COMMON_URL}                 https://keycloak.cicd.odahu.org/auth/realms/legion/
${KEYCLOAK.OAUTH_OIDC_TOKEN_ENDPOINT}  https://keycloak.cicd.odahu.org/auth/realms/legion/protocol/openid-connect/auth
${KEYCLOAK.LOGIN_ACTIONS}              https://keycloak.cicd.odahu.org/auth/realms/legion/login-actions/authenticate?execution=
${KEYCLOAK.USERNAME_TEXTAREA}  id=username
${KEYCLOAK.PASSWORD_TEXTAREA}  id=password
${KEYCLOAK.LOG_IN_BUTTON}      name=login

${KEYCLOAK.ALERT_ERROR}                 xpath://*[@id="kc-content-wrapper"]/div[1]
${KEYCLOAK.INVALID_LOGIN_ERROR_TEXT}    Invalid username or password.

*** Keywords ***
Fill "Username or email"
    [Arguments]     ${username}
    input password  ${KEYCLOAK.USERNAME_TEXTAREA}  ${username}

"Username or email" field should contain
    [Arguments]     ${username}
    textfield value should be  ${KEYCLOAK.USERNAME_TEXTAREA}  ${username}

Fill "Password"
    [Arguments]     ${password}
    input password  ${KEYCLOAK.PASSWORD_TEXTAREA}  ${password}

Click "Log In" Button
    click button  ${KEYCLOAK.LOG_IN_BUTTON}

Log In
    [Arguments]     ${username}  ${password}
    Fill "Username or email"  ${username}
    Fill "Password"          ${password}
    Click "Log In" Button

Validate "Log In" page loaded
    wait until location contains    ${KEYCLOAK.OAUTH_OIDC_TOKEN_ENDPOINT}
    title should be                 ${KEYCLOAK.LOG_IN_TITLE}
    wait until page contains        ${KEYCLOAK.LOG_IN_HEADING_TEXT}

Validate "Log In" page
    wait until location contains    ${KEYCLOAK.COMMON_URL}
    title should be                 ${KEYCLOAK.LOG_IN_TITLE}
    wait until page contains        ${KEYCLOAK.LOG_IN_HEADING_TEXT}

Validate "Invalid username or password" alert shows up
    page should contain element  ${KEYCLOAK.ALERT_ERROR}
    element should be visible  ${KEYCLOAK.ALERT_ERROR}
    element text should be  ${KEYCLOAK.ALERT_ERROR}  ${KEYCLOAK.INVALID_LOGIN_ERROR_TEXT}

Validate "Log In" page after trying to login
    wait until location contains    ${KEYCLOAK.LOGIN_ACTIONS}
    title should be                 ${KEYCLOAK.LOG_IN_TITLE}
    wait until page contains        ${KEYCLOAK.LOG_IN_HEADING_TEXT}
    element should be visible       ${KEYCLOAK.LOG_IN_HEADING}

