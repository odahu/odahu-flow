*** Settings ***
Library  SeleniumLibrary  timeout=10s
Resource        ${RES_DIR}/keywords.robot

*** Variables ***
${DASHBOARD.EMPTY_SPACE}            xpath:/html/body/div[2]/div[1]
${DASHBOARD.BENTO_MENU.TAB}         xpath:/html/body/div[2]/div[3]

# Text values
${DASHBOARD.ODAHU_TITLE}            ODAHU
${DASHBOARD.HEADING}                Getting Started
${DASHBOARD.URL}                    ${EDGE_URL}/dashboard
${DASHBOARD.BENTO_HEADING}          ODAHU Components

# ODAHU Components
${DASHBOARD.BENTO_MENU.DOCS.TEXT}              Docs
${DASHBOARD.BENTO_MENU.API_GATEWAY.TEXT}       API Gateway
${DASHBOARD.BENTO_MENU.MLFLOW.TEXT}            ML Metrics
${DASHBOARD.BENTO_MENU.SERVICE_CATALOG.TEXT}   Service Catalog
${DASHBOARD.BENTO_MENU.GRAFANA.TEXT}           Cluster Monitoring
${DASHBOARD.BENTO_MENU.JUPYTERHUB.TEXT}        JupyterHub
${DASHBOARD.BENTO_MENU.AIRFLOW.TEXT}           Airflow
${DASHBOARD.BENTO_MENU.KIBANA.TEXT}            Kibana
${DASHBOARD.BENTO_MENU.FEEDBACK_STORAGE.TEXT}  Feedback storage

# ODAHU Components links
${DASHBOARD.BENTO_MENU.DOCS.URL}                ${DOCS.URL}
${DASHBOARD.BENTO_MENU.API_GATEWAY.URL}         ${EDGE_URL}/swagger/index.html
${DASHBOARD.BENTO_MENU.MLFLOW.URL}              ${EDGE_URL}/mlflow/
${DASHBOARD.BENTO_MENU.SERVICE_CATALOG.URL}     ${EDGE_URL}/service-catalog/catalog/index.html
${DASHBOARD.BENTO_MENU.GRAFANA.URL}             ${EDGE_URL}/grafana/
${DASHBOARD.BENTO_MENU.JUPYTERHUB.URL}          ${EDGE_URL}/jupyterhub/hub/
${DASHBOARD.BENTO_MENU.AIRFLOW.URL}             ${EDGE_URL}/airflow/
${DASHBOARD.BENTO_MENU.KIBANA.URL}              ${EDGE_URL}/kibana/app/home

# Button locators
${DASHBOARD.ODAHU_COMPONENTS.BENTO_BUTTON}  xpath://*[@id="root"]/div/header/div/div[2]/button
${DASHBOARD.USER_INFO_BUTTON}       xpath://*[@id="root"]/div/header/div/div[3]/button
${DASHBOARD.INFO_BUTTON}            xpath://*[@id="root"]/div/header/div/div[4]/button
${DASHBOARD.LOG_OUT_BUTTON}         xpath:/html/body/div[3]/div[3]/div/button

# ODAHU Components button locators
#                             button description:/html/body/div[2]/div[3]/div[2]/div[1]/div/button/div/p
${DASHBOARD.BENTO_MENU.DOCS}               xpath:/html/body/div[2]/div[3]/div[2]/div[1]/div/button
${DASHBOARD.BENTO_MENU.API_GATEWAY}        xpath:/html/body/div[2]/div[3]/div[2]/div[2]/div/button
${DASHBOARD.BENTO_MENU.MLFLOW}             xpath:/html/body/div[2]/div[3]/div[2]/div[3]/div/button
${DASHBOARD.BENTO_MENU.SERVICE_CATALOG}    xpath:/html/body/div[2]/div[3]/div[2]/div[4]/div/button
${DASHBOARD.BENTO_MENU.GRAFANA}            xpath:/html/body/div[2]/div[3]/div[2]/div[5]/div/button
${DASHBOARD.BENTO_MENU.JUPYTERHUB}         xpath:/html/body/div[2]/div[3]/div[2]/div[6]/div/button
${DASHBOARD.BENTO_MENU.AIRFLOW}            xpath:/html/body/div[2]/div[3]/div[2]/div[7]/div/button
${DASHBOARD.BENTO_MENU.KIBANA}             xpath:/html/body/div[2]/div[3]/div[2]/div[8]/div/button
${DASHBOARD.BENTO_MENU.FEEDBACK_STORAGE}   xpath:/html/body/div[2]/div[3]/div[2]/div[9]/div/button

# Text locators
${DASHBOARD.ODAHU_COMPONENTS.BENTO_TEXT}  xpath:/html/body/div[2]/div[3]/div[1]/h2
${DASHBOARD.ODAHU_UI_TEXT}          xpath:/html/body/div[4]/div[3]/div/h6
${DASHBOARD.USER_INFO.USERNAME}     xpath:/html/body/div[3]/div[3]/div/h4
${DASHBOARD.USER_INFO.EMAIL}        xpath:/html/body/div[3]/div[3]/div/h6

# Docs urls
${DOCS.URL}               https://docs.odahu.org
${DOCS.QUICKSTART_LINK}   ${DOCS.URL}/tutorials_wine.html
${DOCS.CONNECTIONS_LINK}  ${DOCS.URL}/ref_connections.html
${DOCS.TRAINING_LINK}     ${DOCS.URL}/ref_trainings.html
${DOCS.PACKAGING_LINK}    ${DOCS.URL}/ref_packagers.html
${DOCS.DEPLOYMENT_LINK}   ${DOCS.URL}/ref_deployments.html

# Docs link locators
${DASHBOARD.QUICKSTART_LINK}    xpath://*[@id="root"]/div/main/div[2]/div/div[1]/div/div[2]/ul/a[1]
${DASHBOARD.CONNECTIONS_LINK}   xpath://*[@id="root"]/div/main/div[2]/div/div[1]/div/div[2]/ul/a[2]
${DASHBOARD.TRAINING_LINK}      xpath://*[@id="root"]/div/main/div[2]/div/div[1]/div/div[2]/ul/a[3]
${DASHBOARD.PACKAGING_LINK}     xpath://*[@id="root"]/div/main/div[2]/div/div[1]/div/div[2]/ul/a[4]
${DASHBOARD.DEPLOYMENT_LINK}    xpath://*[@id="root"]/div/main/div[2]/div/div[1]/div/div[2]/ul/a[5]

*** Keywords ***
Validate "Dashboard" page loaded
    Validate page  ${DASHBOARD.URL}  ${DASHBOARD.ODAHU_TITLE}  ${DASHBOARD.HEADING}

Click on Empty field on "Dashboard" page
    click element  ${DASHBOARD.EMPTY_SPACE}

Open User info Tab
    click button  ${DASHBOARD.USER_INFO_BUTTON}

Validate "Username" and "Email"
    Validate visible element and text  ${DASHBOARD.USER_INFO.USERNAME}  ${DASHBOARD.USER_INFO.USERNAME.TEXT}
    Validate visible element and text  ${DASHBOARD.USER_INFO.EMAIL}  ${DASHBOARD.USER_INFO.EMAIL.TEXT}

Click "LOG OUT" Button
    click button  ${DASHBOARD.LOG_OUT_BUTTON}

Click "Info" button
    click button  ${DASHBOARD.INFO_BUTTON}

Validate text of "Info" pop-up
    [Arguments]  ${odahu_ui_version}
    Validate visible element and text  ${DASHBOARD.ODAHU_UI_TEXT}  ${odahu_ui_version}

Open Bento Menu (ODAHU Components)
    click button  ${DASHBOARD.ODAHU_COMPONENTS.BENTO_BUTTON}

Validate Bento Menu
    wait until element contains  ${DASHBOARD.ODAHU_COMPONENTS.BENTO_TEXT}  ${DASHBOARD.BENTO_HEADING}
    wait until element is visible  ${DASHBOARD.BENTO_MENU.TAB}
    element should be visible  ${DASHBOARD.BENTO_MENU.TAB}

Validate Bento Menu closed
    wait until element is not visible  ${DASHBOARD.BENTO_MENU.TAB}
    element should not be visible  ${DASHBOARD.BENTO_MENU.TAB}

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
