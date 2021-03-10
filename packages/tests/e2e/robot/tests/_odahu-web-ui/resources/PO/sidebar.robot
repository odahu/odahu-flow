*** Settings ***
Documentation   Page Object for "sidebar" part of UI
Library         SeleniumLibrary  timeout=10s
Resource        ${RES_DIR}/keywords.robot

*** Test Cases ***
Test title
    [Tags]    DEBUG
    Provided precondition
    When action
    Then check expectations

*** Keywords ***
Validate SideBar exists and visible