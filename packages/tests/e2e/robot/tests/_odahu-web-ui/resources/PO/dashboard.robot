*** Settings ***
Documentation   Page Object for main part of "Dashboard" page
Library         SeleniumLibrary  timeout=10s
Resource        ${RES_DIR}/keywords.robot

*** Variables ***
# Text values
${DASHBOARD.HEADING}                Getting Started
${DASHBOARD.URL}                    ${EDGE_URL}/dashboard

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

# Dashboard chart frame locators
${DASHBOARD.CHART.CONNECTION}   xpath://*[@id="root"]/div/main/div[2]/div/div[2]/div
${DASHBOARD.CHART.TRAINING}     xpath://*[@id="root"]/div/main/div[2]/div/div[3]/div
${DASHBOARD.CHART.PACKAGING}    xpath://*[@id="root"]/div/main/div[2]/div/div[4]/div
${DASHBOARD.CHART.DEPLOYMENT}   xpath://*[@id="root"]/div/main/div[2]/div/div[5]/div

*** Keywords ***
Validate "Dashboard" page loaded
    Validate page  ${DASHBOARD.URL}  ${DASHBOARD.HEADING}

Validate that chart is visible
    # chart frame:       //*[@id="root"]/div/main/div[2]/div/div[4]/div
    # chart canvas:      //*[@id="root"]/div/main/div[2]/div/div[4]/div/div[2]/canvas
    # chart description: //*[@id="root"]/div/main/div[2]/div/div[4]/div/div[1]/div[2]/span/span
    [Arguments]  ${chart_locator}  ${chart_description}
    element should be visible  ${chart_locator}
    element should be visible  ${chart_locator}/div[2]/canvas
    Validate visible element and text  ${chart_locator}/div[1]/div[2]/span/span  ${chart_description}