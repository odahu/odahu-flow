*** Settings ***
Library     SeleniumLibrary  timeout=10s
Resource    ${PAGE_OBJECTS}/keycloak.robot
Resource    ${PAGE_OBJECTS}/dashboard.robot
Variables   ../../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}

*** Variables ***
${PAGE_OBJECTS}  ${CURDIR}/PO

# SHOULD BE CUSTOMISABLE
${BROWSER}          chrome
${DASHBOARD.USER_INFO.USERNAME.TEXT}    anonymous
${DASHBOARD.USER_INFO.EMAIL.TEXT}       anonymous@email.org

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

Open and Validate links
    [Arguments]  ${link_locator}  ${page_url}
    [Teardown]   run keywords
    ...          close window
    ...          AND  switch window  ${handle}
    dashboard.Validate "Dashboard" page loaded
    click link   ${link_locator}
    ${handle}    switch window  locator=NEW
    location should be  ${page_url}

Open Menu with ODAHU Components
    dashboard.Validate "Dashboard" page loaded
    dashboard.Open Bento Menu (ODAHU Components)
    dashboard.Validate Bento Menu

Close Menu with ODAHU Components
    dashboard.Validate "Dashboard" page loaded
    dashboard.Click on Empty field on "Dashboard" page
    dashboard.Validate Bento Menu closed

Validate ODAHU Components are visible
    [Arguments]  ${button_locator}  ${button_description}  ${validation_url}=${EMPTY}
    Open Menu with ODAHU Components  # Setup for test
    dashboard.Validate ODAHU components visible and description present  ${button_locator}  ${button_description}
    run keyword if  '${validation_url}' != '${EMPTY}'
    ...     dashboard.Validate ODAHU components links  ${button_locator}  ${validation_url}
