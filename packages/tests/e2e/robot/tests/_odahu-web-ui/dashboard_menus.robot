*** Variables ***
${RES_DIR}          ${CURDIR}/resources
${UI_VERSION}       ODAHU version: ${ODAHU_WEB_UI_VERSION}

*** Settings ***
Documentation   testing Dashboard page and Menus (Sidebar and ODAHU Components)
Resource        ${RES_DIR}/common.robot
Suite Setup     Setup
Suite Teardown  Teardown
Test Setup      Test Setup
Test Teardown   Test Teardown
Force Tags      web-ui  dashboard

*** Test Cases ***
#       --------- DASHBOARD -----------
Check links to documention on "Dashboard" should opens and lead to docs.odahu.org
    [Template]  Open and Validate links
    ${DASHBOARD.QUICKSTART_LINK}    ${DOCS.QUICKSTART_LINK}
    ${DASHBOARD.CONNECTIONS_LINK}   ${DOCS.CONNECTIONS_LINK}
    ${DASHBOARD.TRAINING_LINK}      ${DOCS.TRAINING_LINK}
    ${DASHBOARD.PACKAGING_LINK}     ${DOCS.PACKAGING_LINK}
    ${DASHBOARD.DEPLOYMENT_LINK}    ${DOCS.DEPLOYMENT_LINK}

Check "Dashboard" charts are visible
    [Template]  Common.Validate that chart is visible
    ${DASHBOARD.CHART.CONNECTION}   Connections
    ${DASHBOARD.CHART.TRAINING}     Trainings
    ${DASHBOARD.CHART.PACKAGING}    Packaging
    ${DASHBOARD.CHART.DEPLOYMENT}   Deployment


#       --------- HEADER -----------
Check "ODAHU Components" menu opens
    Open Menu with ODAHU Components

Check "ODAHU Components" menu closes
    Open Menu with ODAHU Components
    Close Menu with ODAHU Components

Check "ODAHU Components" present and lead to right link
    [Template]  Validate ODAHU Components are visible and lead to the right link
    ${HEADER.BENTO_MENU.DOCS}                ${HEADER.BENTO_MENU.DOCS.TEXT}                ${HEADER.BENTO_MENU.DOCS.URL}
    ${HEADER.BENTO_MENU.API_GATEWAY}         ${HEADER.BENTO_MENU.API_GATEWAY.TEXT}         ${HEADER.BENTO_MENU.API_GATEWAY.URL}
    ${HEADER.BENTO_MENU.MLFLOW}              ${HEADER.BENTO_MENU.MLFLOW.TEXT}              ${HEADER.BENTO_MENU.MLFLOW.URL}
    ${HEADER.BENTO_MENU.SERVICE_CATALOG}     ${HEADER.BENTO_MENU.SERVICE_CATALOG.TEXT}     ${HEADER.BENTO_MENU.SERVICE_CATALOG.URL}
    ${HEADER.BENTO_MENU.GRAFANA}             ${HEADER.BENTO_MENU.GRAFANA.TEXT}             ${HEADER.BENTO_MENU.GRAFANA.URL}
    ${HEADER.BENTO_MENU.JUPYTERHUB}          ${HEADER.BENTO_MENU.JUPYTERHUB.TEXT}          ${HEADER.BENTO_MENU.JUPYTERHUB.URL}
    ${HEADER.BENTO_MENU.AIRFLOW}             ${HEADER.BENTO_MENU.AIRFLOW.TEXT}             ${HEADER.BENTO_MENU.AIRFLOW.URL}
    ${HEADER.BENTO_MENU.KIBANA}              ${HEADER.BENTO_MENU.KIBANA.TEXT}              ${HEADER.BENTO_MENU.KIBANA.URL}
    ${HEADER.BENTO_MENU.FEEDBACK_STORAGE}    ${HEADER.BENTO_MENU.FEEDBACK_STORAGE.TEXT}    # Cannot validate url,  tries to login

Validate "User Info" matches
    Validate "User Info" button and text fields match

Validate "ODAHU UI Version" matches
    Validate "Info" button present and ODAHU UI Version match  ${UI_VERSION}


#       --------- SIDEBAR -----------
Check "Sidebar" extends
    Extend "SideBar" and validate

Check "Sidebar" extends and shrinks
    Extend "SideBar" and validate
    Shrink "SideBar" and validate

Check "Sidebar" links change color when user switch between them
    [Tags]  test
    Extend "SideBar" and validate
    # switch to ODAHU page through link
    Check ODAHU page Icons changes the color when selected
    # switch to ODAHU page through entities
    Go to entity and validate ODAHU page icons  ${COMMON.TESTING_ICON_CONN_ENTITY}  ${SIDEBAR.LINK.CONNECTIONS.ICON}
    Go to entity and validate ODAHU page icons  ${COMMON.TESTING_ICON_TI_ENTITY}  ${SIDEBAR.LINK.TOOLCHAINS.ICON}
