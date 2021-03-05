*** Variables ***
${PAGE_OBJECTS}  ${CURDIR}/PO

# should be customisable
${BROWSER}          chrome
${CLUSTER_URL}      https://cluster.odahu.org
${UI_USERNAME}      username
${UI_PASSWORD}      password
${DASHBOARD.USER_INFO.USERNAME.TEXT}    dashboard_user
${DASHBOARD.USER_INFO.EMAIL.TEXT}       dashboard_user@email.org

*** Settings ***
Library     SeleniumLibrary  timeout=10s
Resource    ${PAGE_OBJECTS}/keycloak.robot
Resource    ${PAGE_OBJECTS}/dashboard.robot


*** Keywords ***
#       --------- COMMON -----------
Begin Web Test
    open browser  ${CLUSTER_URL}  ${BROWSER}
    maximize browser window

End Web Test
    close all browsers

Setup
    Begin Web Test
    Login to ODAHU WebUI  ${UI_USERNAME}  ${UI_PASSWORD}

Teardown
    End Web Test

Test Setup
    reload page

#       --------- LOGIN -----------
Login to ODAHU WebUI
    [Arguments]     ${username}  ${password}
    Keycloak.Validate "Log In" page loaded
    Keycloak.Log In  ${username}  ${password}
    Dashboard.Validate "Dashboard" page loaded

Log Out from ODAHU WebUI
    Dashboard.Open User info Tab
    Dashboard.Click "LOG OUT" Button
    Keycloak.Validate "Log In" page loaded

Fail Login with invalid credentials
    [Arguments]     ${username}  ${password}
    Keycloak.Validate "Log In" page
    Keycloak.Log In  ${username}  ${password}
    Keycloak."Username or email" field should contain  ${username}
    Keycloak.Validate "Invalid username or password" alert shows up
    Keycloak.Validate "Log In" page after trying to login

#       --------- DASHBOARD -----------
Validate "Info" button present and ODAHU UI Version match
    [Arguments]  ${odahu_ui_version}
    dashboard.Validate "Dashboard" page loaded
    dashboard.Click "Info" button
    dashboard.Validate text of "Info" pop-up  ${odahu_ui_version}

Validate "User Info" button and text fields match
    dashboard.Validate "Dashboard" page loaded
    dashboard.Open User info Tab
    dashboard.Validate "Username" and "Email"
