*** Settings ***
Library  SeleniumLibrary  timeout=10s
Resource        ${RES_DIR}/keywords.robot

*** Variables ***
${DASHBOARD.ODAHU_TITLE}          ODAHU
${DASHBOARD.HEADING}    Getting Started
${DASHBOARD.URL}        ${CLUSTER_URL}/dashboard

${DASHBOARD.LOG_OUT_BUTTON}         xpath:/html/body/div[3]/div[3]/div/button
${DASHBOARD.USER_INFO_BUTTON}       xpath://*[@id="root"]/div/header/div/div[3]/button
${DASHBOARD.INFO_BUTTON}            xpath://*[@id="root"]/div/header/div/div[4]/button
${DASHBOARD.ODAHU_UI_TEXT}          xpath:/html/body/div[4]/div[3]/div/h6
${DASHBOARD.USER_INFO.USERNAME}     xpath:/html/body/div[3]/div[3]/div/h4
${DASHBOARD.USER_INFO.EMAIL}        xpath:/html/body/div[3]/div[3]/div/h6

*** Keywords ***
Validate "Dashboard" page loaded
    Validate page  ${DASHBOARD.URL}  ${DASHBOARD.ODAHU_TITLE}  ${DASHBOARD.HEADING}

Open User info Tab
    Validate button and click  ${DASHBOARD.USER_INFO_BUTTON}

Validate "Username" and "Email"
    Validate visible element and text  ${DASHBOARD.USER_INFO.USERNAME}  ${DASHBOARD.USER_INFO.USERNAME.TEXT}
    Validate visible element and text  ${DASHBOARD.USER_INFO.EMAIL}  ${DASHBOARD.USER_INFO.EMAIL.TEXT}

Click "LOG OUT" Button
    Validate button and click  ${DASHBOARD.LOG_OUT_BUTTON}

Click "Info" button
    Validate button and click  ${DASHBOARD.INFO_BUTTON}

Validate text of "Info" pop-up
    [Arguments]  ${odahu_ui_version}
    Validate visible element and text  ${DASHBOARD.ODAHU_UI_TEXT}  ${odahu_ui_version}
