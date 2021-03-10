*** Variables ***
${RES_DIR}          ${CURDIR}/resources

${UI_VERSION}  ODAHU version: ${ODAHU_WEB_UI_VERSION}

*** Settings ***
Documentation   testing Dashboard page and Menus (Sidebar and ODAHU Components)
Resource        ${RES_DIR}/common.robot
Suite Setup     Setup
Suite Teardown  Teardown
Test Teardown   Test Teardown
Force Tags      web-ui  dashboard

*** Test Cases ***
Validate ODAHU UI Version matches
    Validate "Info" button present and ODAHU UI Version match  ${UI_VERSION}

Validate User Info matches
    Validate "User Info" button and text fields match

Check links to documention on "Dashboard" should opens and lead to docs.odahu.org
    [Template]  Open and Validate links
    ${DASHBOARD.QUICKSTART_LINK}    ${DOCS.QUICKSTART_LINK}
    ${DASHBOARD.CONNECTIONS_LINK}   ${DOCS.CONNECTIONS_LINK}
    ${DASHBOARD.TRAINING_LINK}      ${DOCS.TRAINING_LINK}
    ${DASHBOARD.PACKAGING_LINK}     ${DOCS.PACKAGING_LINK}
    ${DASHBOARD.DEPLOYMENT_LINK}    ${DOCS.DEPLOYMENT_LINK}

Check ODAHU Components menu opens
    Open Menu with ODAHU Components

Check ODAHU Components menu closes
    Open Menu with ODAHU Components
    Close Menu with ODAHU Components

Check ODAHU Components present and lead to right link
    [Template]  Validate ODAHU Components are visible
    ${DASHBOARD.BENTO_MENU.DOCS}                ${DASHBOARD.BENTO_MENU.DOCS.TEXT}                ${DASHBOARD.BENTO_MENU.DOCS.URL}
    ${DASHBOARD.BENTO_MENU.API_GATEWAY}         ${DASHBOARD.BENTO_MENU.API_GATEWAY.TEXT}         ${DASHBOARD.BENTO_MENU.API_GATEWAY.URL}
    ${DASHBOARD.BENTO_MENU.MLFLOW}              ${DASHBOARD.BENTO_MENU.MLFLOW.TEXT}              ${DASHBOARD.BENTO_MENU.MLFLOW.URL}
    ${DASHBOARD.BENTO_MENU.SERVICE_CATALOG}     ${DASHBOARD.BENTO_MENU.SERVICE_CATALOG.TEXT}     ${DASHBOARD.BENTO_MENU.SERVICE_CATALOG.URL}
    ${DASHBOARD.BENTO_MENU.GRAFANA}             ${DASHBOARD.BENTO_MENU.GRAFANA.TEXT}             ${DASHBOARD.BENTO_MENU.GRAFANA.URL}
    ${DASHBOARD.BENTO_MENU.JUPYTERHUB}          ${DASHBOARD.BENTO_MENU.JUPYTERHUB.TEXT}          ${DASHBOARD.BENTO_MENU.JUPYTERHUB.URL}
    ${DASHBOARD.BENTO_MENU.AIRFLOW}             ${DASHBOARD.BENTO_MENU.AIRFLOW.TEXT}             ${DASHBOARD.BENTO_MENU.AIRFLOW.URL}
    ${DASHBOARD.BENTO_MENU.KIBANA}              ${DASHBOARD.BENTO_MENU.KIBANA.TEXT}              ${DASHBOARD.BENTO_MENU.KIBANA.URL}
    ${DASHBOARD.BENTO_MENU.FEEDBACK_STORAGE}    ${DASHBOARD.BENTO_MENU.FEEDBACK_STORAGE.TEXT}    # Cannot validate url,  tries to login
