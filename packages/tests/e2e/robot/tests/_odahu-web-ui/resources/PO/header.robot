*** Settings ***
Documentation   Page object for header of ODAHU Web UI
Library         SeleniumLibrary  timeout=10s
Resource        ${PAGE_OBJECTS_RES}/keywords.robot

*** Variables ***
${HEADER.ODAHU_TITLE}            ODAHU
${HEADER.BENTO_HEADING}          ODAHU Components

${HEADER.CLOSE_BENTO_LOCATOR}    xpath:/html/body/div[2]/div[1]  # when bento menu open. it's another object on the page
${HEADER.BENTO_MENU.TAB}         xpath:/html/body/div[2]/div[3]

# Text locators
${HEADER.ODAHU_COMPONENTS.BENTO_TEXT}  xpath:/html/body/div[2]/div[3]/div[1]/h2
${HEADER.ODAHU_TITLE.TEXT}      xpath://*[@id="root"]/div/header/div/h6
${HEADER.ODAHU_UI_VERSION}      xpath:/html/body/div[4]/div[3]/div/h6
${HEADER.USER_INFO.USERNAME}    xpath:/html/body/div[3]/div[3]/div/h4
${HEADER.USER_INFO.EMAIL}       xpath:/html/body/div[3]/div[3]/div/h6

# ODAHU Components
${HEADER.BENTO_MENU.DOCS.TEXT}              Docs
${HEADER.BENTO_MENU.API_GATEWAY.TEXT}       API Gateway
${HEADER.BENTO_MENU.MLFLOW.TEXT}            ML Metrics
${HEADER.BENTO_MENU.SERVICE_CATALOG.TEXT}   Service Catalog
${HEADER.BENTO_MENU.GRAFANA.TEXT}           Cluster Monitoring
${HEADER.BENTO_MENU.JUPYTERHUB.TEXT}        JupyterHub
${HEADER.BENTO_MENU.AIRFLOW.TEXT}           Airflow
${HEADER.BENTO_MENU.KIBANA.TEXT}            Kibana
${HEADER.BENTO_MENU.FEEDBACK_STORAGE.TEXT}  Feedback storage

# ODAHU Components links
${HEADER.BENTO_MENU.DOCS.URL}                ${DOCS.URL}
${HEADER.BENTO_MENU.API_GATEWAY.URL}         ${EDGE_URL}/swagger/index.html
${HEADER.BENTO_MENU.MLFLOW.URL}              ${EDGE_URL}/mlflow/
${HEADER.BENTO_MENU.SERVICE_CATALOG.URL}     ${EDGE_URL}/service-catalog/catalog/index.html
${HEADER.BENTO_MENU.GRAFANA.URL}             ${EDGE_URL}/grafana/
${HEADER.BENTO_MENU.JUPYTERHUB.URL}          ${EDGE_URL}/jupyterhub/hub/
${HEADER.BENTO_MENU.AIRFLOW.URL}             ${EDGE_URL}/airflow/
${HEADER.BENTO_MENU.KIBANA.URL}              ${EDGE_URL}/kibana/app/home

# Button locators
${HEADER.ODAHU_COMPONENTS.BENTO_BUTTON}  xpath://*[@id="root"]/div/header/div/div[2]/button
${HEADER.USER_INFO_BUTTON}       xpath://*[@id="root"]/div/header/div/div[3]/button
${HEADER.INFO_BUTTON}            xpath://*[@id="root"]/div/header/div/div[4]/button
${HEADER.LOG_OUT_BUTTON}         xpath:/html/body/div[3]/div[3]/div/button

# ODAHU Components button locators
#                    button description xpath:/html/body/div[2]/div[3]/div[2]/div[1]/div/button/div/p
${HEADER.BENTO_MENU.DOCS}               xpath:/html/body/div[2]/div[3]/div[2]/div[1]/div/button
${HEADER.BENTO_MENU.API_GATEWAY}        xpath:/html/body/div[2]/div[3]/div[2]/div[2]/div/button
${HEADER.BENTO_MENU.MLFLOW}             xpath:/html/body/div[2]/div[3]/div[2]/div[3]/div/button
${HEADER.BENTO_MENU.SERVICE_CATALOG}    xpath:/html/body/div[2]/div[3]/div[2]/div[4]/div/button
${HEADER.BENTO_MENU.GRAFANA}            xpath:/html/body/div[2]/div[3]/div[2]/div[5]/div/button
${HEADER.BENTO_MENU.JUPYTERHUB}         xpath:/html/body/div[2]/div[3]/div[2]/div[6]/div/button
${HEADER.BENTO_MENU.AIRFLOW}            xpath:/html/body/div[2]/div[3]/div[2]/div[7]/div/button
${HEADER.BENTO_MENU.KIBANA}             xpath:/html/body/div[2]/div[3]/div[2]/div[8]/div/button
${HEADER.BENTO_MENU.FEEDBACK_STORAGE}   xpath:/html/body/div[2]/div[3]/div[2]/div[9]/div/button

*** Keywords ***
Validate Header Loaded
    Validate ODAHU Web UI page loaded  ${HEADER.ODAHU_TITLE.TEXT}  ${HEADER.ODAHU_COMPONENTS.BENTO_BUTTON}

Click on Empty field on "Dashboard" page when Bento Menu open
    click element  ${HEADER.CLOSE_BENTO_LOCATOR}

Open User info Tab
    click button  ${HEADER.USER_INFO_BUTTON}

Validate "Username" and "Email"
    Validate visible element and text  ${HEADER.USER_INFO.USERNAME}  ${COMMON.USER_INFO.USERNAME.TEXT}
    Validate visible element and text  ${HEADER.USER_INFO.EMAIL}  ${COMMON.USER_INFO.EMAIL.TEXT}

Click "LOG OUT" Button
    click button  ${HEADER.LOG_OUT_BUTTON}

Click "Info" button
    click button  ${HEADER.INFO_BUTTON}

Validate text of "Info" pop-up
    [Arguments]  ${odahu_ui_version}
    Validate visible element and text  ${HEADER.ODAHU_UI_VERSION}  ${odahu_ui_version}

Open Bento Menu (ODAHU Components)
    click button  ${HEADER.ODAHU_COMPONENTS.BENTO_BUTTON}

Validate Bento Menu
    wait until element contains  ${HEADER.ODAHU_COMPONENTS.BENTO_TEXT}  ${HEADER.BENTO_HEADING}
    wait until element is visible  ${HEADER.BENTO_MENU.TAB}
    element should be visible  ${HEADER.BENTO_MENU.TAB}

Validate Bento Menu closed
    wait until element is not visible  ${HEADER.BENTO_MENU.TAB}
    element should not be visible  ${HEADER.BENTO_MENU.TAB}

Validate ODAHU components visible and description present
    # button description  -> xpath:/html/body/div[2]/div[3]/div[2]/div[1]/div/button/div/p
    # button_locator      -> xpath:/html/body/div[2]/div[3]/div[2]/div[1]/div/button
    [Arguments]  ${button_locator}  ${button_description}
    page should contain button  ${button_locator}
    element should be visible  ${button_locator}
    Validate visible element and text  ${button_locator}/div/p  ${button_description}

Validate ODAHU components links
    [Arguments]  ${button_locator}  ${validation_url}
    [Teardown]   run keywords
    ...          close window
    ...          AND  switch window  ${handle}
    page should contain button  ${button_locator}
    element should be visible  ${button_locator}
    click button  ${button_locator}
    ${handle}     switch window  locator=NEW
    location should contain  ${validation_url}
