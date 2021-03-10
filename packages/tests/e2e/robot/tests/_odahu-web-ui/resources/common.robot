*** Settings ***
Library     SeleniumLibrary  timeout=10s
Resource    ${PAGE_OBJECTS}/keycloak.robot
Resource    ${PAGE_OBJECTS}/dashboard.robot
Resource    ${PAGE_OBJECTS}/header.robot
Variables   ../../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}

*** Variables ***
${PAGE_OBJECTS}  ${CURDIR}/PO

# SHOULD BE CUSTOMISABLE
${BROWSER}          chrome
${COMMON.USER_INFO.USERNAME.TEXT}    anonymous
${COMMON.USER_INFO.EMAIL.TEXT}       anonymous@email.org

*** Keywords ***
#       --------- COMMON -----------
Begin Web Test
    open browser  ${EDGE_URL}  ${BROWSER}
    maximize browser window

End Web Test
    close browser

Setup
    Begin Web Test
    Login to ODAHU WebUI  ${ODAHU_WEB_UI_USERNAME}  ${ODAHU_WEB_UI_PASSWORD}

Teardown
    End Web Test

Test Teardown
    reload page

#       --------- LOGIN -----------
Login to ODAHU WebUI
    [Arguments]     ${username}  ${password}
    Keycloak.Validate "Log In" page loaded
    Keycloak.Log In  ${username}  ${password}
    Dashboard.Validate "Dashboard" page loaded

Log Out from ODAHU WebUI
    Header.Open User info Tab
    Header.Click "LOG OUT" Button
    Keycloak.Validate "Log In" page loaded

Fail Login with invalid credentials
    [Arguments]     ${username}  ${password}
    Keycloak.Validate "Log In" page
    Keycloak.Log In  ${username}  ${password}
    Keycloak."Username or email" field should contain  ${username}
    Keycloak.Validate "Invalid username or password" alert shows up
    Keycloak.Validate "Log In" page after trying to login

#       --------- DASHBOARD AND MENUS -----------
Validate "Info" button present and ODAHU UI Version match
    [Arguments]  ${odahu_ui_version}
    Dashboard.Validate "Dashboard" page loaded
    Header.Click "Info" button
    Header.Validate text of "Info" pop-up  ${odahu_ui_version}

Validate "User Info" button and text fields match
    Dashboard.Validate "Dashboard" page loaded
    Header.Open User info Tab
    Header.Validate "Username" and "Email"

Open and Validate links
    [Arguments]  ${link_locator}  ${page_url}
    [Teardown]   run keywords
    ...          close window
    ...          AND  switch window  ${handle}
    Dashboard.Validate "Dashboard" page loaded
    click link   ${link_locator}
    ${handle}    switch window  locator=NEW
    location should be  ${page_url}

Open Menu with ODAHU Components
    Dashboard.Validate "Dashboard" page loaded
    Header.Open Bento Menu (ODAHU Components)
    Header.Validate Bento Menu

Close Menu with ODAHU Components
    Dashboard.Validate "Dashboard" page loaded
    Header.Click on Empty field on "Dashboard" page
    Header.Validate Bento Menu closed

Validate ODAHU Components are visible
    [Arguments]  ${button_locator}  ${button_description}  ${validation_url}=${EMPTY}
    Open Menu with ODAHU Components  # Setup for test
    Header.Validate ODAHU components visible and description present  ${button_locator}  ${button_description}
    run keyword if  '${validation_url}' != '${EMPTY}'
    ...     Header.Validate ODAHU components links  ${button_locator}  ${validation_url}

Validate that chart is visible
    [Arguments]  ${chart_locator}  ${chart_description}
    Dashboard.Validate "Dashboard" page loaded
    Dashboard.Validate that chart is visible  ${chart_locator}  ${chart_description}