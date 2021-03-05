*** Variables ***
${RES_DIR}          ${CURDIR}/resources
${PAGE_OBJECTS}     ${RES_DIR}/PO

${UI_VERSION}           1.5.0-b1614899374909
${ODAHU_UI_VERSION}  ODAHU version: ${UI_VERSION}

*** Settings ***
Documentation   testing Dashboard page and Menus (Sidebar and ODAHU Components)
Resource        ${RES_DIR}/common.robot
Suite Setup     Setup
Suite Teardown  Teardown
Test Teardown   Test Setup
Force Tags      web-ui  dashboard

*** Test Cases ***
Validate ODAHU UI Version matches
    Validate "Info" button present and ODAHU UI Version match  ${ODAHU_UI_VERSION}

Validate User Info matches
    Validate "User Info" button and text fields match
