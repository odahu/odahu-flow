*** Variables ***
${RES_DIR}          ${CURDIR}/resources
${UI_VERSION}       ODAHU version: ${ODAHU_WEB_UI_VERSION}

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

Check Dashboard charts are visible
    [Template]  Common.Validate that chart is visible
    ${DASHBOARD.CHART.CONNECTION}   Connections
    ${DASHBOARD.CHART.TRAINING}     Trainings
    ${DASHBOARD.CHART.PACKAGING}    Packaging
    ${DASHBOARD.CHART.DEPLOYMENT}   Deployment
