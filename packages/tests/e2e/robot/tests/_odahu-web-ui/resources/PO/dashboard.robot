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

*** Keywords ***
Validate "Dashboard" page loaded
    Validate page  ${DASHBOARD.URL}  ${HEADER.ODAHU_TITLE}  ${DASHBOARD.HEADING}
