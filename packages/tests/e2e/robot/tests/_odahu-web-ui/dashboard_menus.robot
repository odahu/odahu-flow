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
    ${HEADER.BENTO_MENU.DOCS}                ${HEADER.BENTO_MENU.DOCS.TEXT}                ${HEADER.BENTO_MENU.DOCS.URL}
    ${HEADER.BENTO_MENU.API_GATEWAY}         ${HEADER.BENTO_MENU.API_GATEWAY.TEXT}         ${HEADER.BENTO_MENU.API_GATEWAY.URL}
    ${HEADER.BENTO_MENU.MLFLOW}              ${HEADER.BENTO_MENU.MLFLOW.TEXT}              ${HEADER.BENTO_MENU.MLFLOW.URL}
    ${HEADER.BENTO_MENU.SERVICE_CATALOG}     ${HEADER.BENTO_MENU.SERVICE_CATALOG.TEXT}     ${HEADER.BENTO_MENU.SERVICE_CATALOG.URL}
    ${HEADER.BENTO_MENU.GRAFANA}             ${HEADER.BENTO_MENU.GRAFANA.TEXT}             ${HEADER.BENTO_MENU.GRAFANA.URL}
    ${HEADER.BENTO_MENU.JUPYTERHUB}          ${HEADER.BENTO_MENU.JUPYTERHUB.TEXT}          ${HEADER.BENTO_MENU.JUPYTERHUB.URL}
    ${HEADER.BENTO_MENU.AIRFLOW}             ${HEADER.BENTO_MENU.AIRFLOW.TEXT}             ${HEADER.BENTO_MENU.AIRFLOW.URL}
    ${HEADER.BENTO_MENU.KIBANA}              ${HEADER.BENTO_MENU.KIBANA.TEXT}              ${HEADER.BENTO_MENU.KIBANA.URL}
    ${HEADER.BENTO_MENU.FEEDBACK_STORAGE}    ${HEADER.BENTO_MENU.FEEDBACK_STORAGE.TEXT}    # Cannot validate url,  tries to login

# Check Sidebar extends
#
#
# Check Sidebar extends and shrinks