*** Settings ***
Documentation   Page Object for "Connections" page
Library         SeleniumLibrary  timeout=10s
Library         String
Resource        ${PAGE_OBJECTS_RES}/keywords.robot
Resource        ${PAGE_OBJECTS_RES}/entities.robot
Resource        ../../../../resources/keywords.robot

*** Variables ***
# Text values
${CONNECTIONS.HEADING}          Connections
${CONNECTIONS.URL}              ${EDGE_URL}/connections

# Alert messages
${CONNECTIONS.CONNECTION_CREATED}  Success\nConnection was created


@{CONNECTIONS.TABLE_HEADERS}    ID  Type  URI  Description  WEB UI  Created at  Updated at

${CONNECTIONS.ALERT.SUCCESS}                    css:#root > div > div.MuiSnackbar-root.MuiSnackbar-anchorOriginBottomLeft > div
${CONNECTIONS.ALERT.SUCCESS.CREATED_MESSAGE}    css:#root > div > div.MuiSnackbar-root.MuiSnackbar-anchorOriginBottomLeft > div > div.MuiAlert-message

# ---- Locators for "New" button ----
${CONNECTIONS.METADATA.NEXT_BUTTON}         xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[1]/div/div/div/div/div[2]/div/button
# ${CONNECTIONS.SPECIFICATION.BACK_BUTTON}  xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[3]/div/div/div/div/div[2]/div/button[1]
${CONNECTIONS.SPECIFICATION.NEXT_BUTTON}    xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[3]/div/div/div/div/div[2]/div/button[2]
# ${CONNECTIONS.REVIEW.BACK_BUTTON}         xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/div/button[1]
${CONNECTIONS.REVIEW.SUBMIT_BUTTON}         xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/div/button[2]

# -- Metadata --
${CONNECTIONS.ID}               name=id
${CONNECTIONS.TYPE.LISTBOX}     id=mui-component-select-spec.type
${CONNECTIONS.TYPE.LIST}        xpath://*[@id="menu-spec.type"]/div[3]/ul
${CONNECTIONS.TYPE}             name=spec.type
${CONNECTIONS.WEBLINK}          name=spec.webUILink
${CONNECTIONS.DESCRIPTION}      name=spec.description
# -- Specification --
${CONNECTIONS.SPECIFICATION.STEP_CONTENT}   xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[3]/div
${CONNECTIONS.URI}              name=spec.uri
${CONNECTIONS.GIT.REFERENCE}    name=spec.reference
# ${CONNECTIONS.REGION.LIST}      id=mui-autocomplete-37240  # for s3 connection type
${CONNECTIONS.REGION}           name=spec.region
${CONNECTIONS.KEYID}            name=spec.keyID
${CONNECTIONS.KEYSECRET}        name=spec.keySecret
${CONNECTIONS.USERNAME}         name=spec.username
${CONNECTIONS.PASSWORD}         name=spec.password
# # -- Review --
${CONNECTIONS.REVIEW.STEP_CONTENT}  xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div
# ${CONNECTIONS.REVIEW.ID}                    xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/ul/li[1]/div/p
# ${CONNECTIONS.REVIEW.TYPE}                  xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/ul/li[2]/div/p
# ${CONNECTIONS.REVIEW.GIT.WEB_UI_LINK}       xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/ul/li[3]
# ${CONNECTIONS.REVIEW.GIT.URI}               xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/ul/li[4]
# ${CONNECTIONS.REVIEW.GIT.REFERRENCE}        xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/ul/li[5]
# ${CONNECTIONS.REVIEW.GIT.SSH_PRIVATE_KEY}   xpath://*[@id="root"]/div/main/div[2]/div/div/form/div/div[5]/div/div/div/div/ul/li[6]

*** Keywords ***
Validate "Connections" page loaded
    Entities.Validate page with entities    ${CONNECTIONS.URL}  ${CONNECTIONS.HEADING}

Validate "Connections" page table loaded
    Entities.Validate Table loaded  @{CONNECTIONS.TABLE_HEADERS}

Fill "ID" field with value
    [Arguments]  ${id}
    ${random}   String.Generate Random String  33  [LOWER][NUMBERS]
    input text  ${CONNECTIONS.ID}  ${id}-${random}

Select Type from the list
    [Arguments]  ${type}
    ${get_items}    get webelements  ${CONNECTIONS.TYPE.LIST}/li
    FOR  ${element}  IN  @{get_items}
        ${text}  get text  ${element}
        log  ${text}
        Exit For Loop IF    "${text}" == "${type}"
    END
    click element  ${element}

Wait until "Type" list is visible
    wait until element is visible  ${CONNECTIONS.TYPE.LIST}

Select Type
    [Arguments]  ${type}
    click element  ${CONNECTIONS.TYPE.LISTBOX}
    Wait until "Type" list is visible
    Select Type from the list  ${type}
    capture page screenshot

Fill Web UI link field
    [Arguments]  ${web_ui_link}
    input text  ${CONNECTIONS.WEBLINK}  ${web_ui_link}

Add description for connection
    [Arguments]  ${description}
    input text  ${CONNECTIONS.DESCRIPTION}  ${description}

Click "Next" button on "Metadata" Step
    click button  ${CONNECTIONS.METADATA.NEXT_BUTTON}

Validate "Specification" Step Content loaded
    element should be visible  ${CONNECTIONS.SPECIFICATION.STEP_CONTENT}

### --  Connections' Specifications  -- ###
Insert URI
    [Arguments]  ${uri}
    input text  ${CONNECTIONS.URI}  ${uri}

# -- GIT Specification --
Insert Git SSH URL
    [Arguments]  ${git_ssh_url}
    input text  ${CONNECTIONS.URI}  ${git_ssh_url}

Insert Reference
    [Arguments]  ${reference}
    input text  ${CONNECTIONS.GIT.REFERENCE}  ${reference}

Insert SSH private key
    [Arguments]  ${ssh_private_key}
    input password  ${CONNECTIONS.KEYSECRET}  ${ssh_private_key}

# -- GCS Specification --
Insert Project
    [Arguments]  ${project}
    input text  ${CONNECTIONS.REGION}  ${project}

Insert Service account secret
    [Arguments]  ${SA_secret}
    ${SA_secret_base64}  Shell  base64 <<< ${SA_secret}
    input password  ${CONNECTIONS.KEYSECRET}  ${SA_secret_base64.stdout}

# -- Docker Specification --
Insert Username
    [Arguments]  ${username}
    input password  ${CONNECTIONS.USERNAME}  ${username}

Insert Password
    [Arguments]  ${password}
    ${password_base64}  Shell  base64 <<< ${password}
    input password  ${CONNECTIONS.PASSWORD}  ${password_base64.stdout}

# -- Azureblob Specification --
Insert SAS Token
    [Arguments]  ${SAS_Token}
    ${SAS_Token_base64}  Shell  base64 <<< ${SAS_Token}
    input password  ${CONNECTIONS.KEYSECRET}  ${SAS_Token_base64.stdout}

# -- S3, ECR Specification --
Insert Access Key ID
    [Arguments]  ${access_key_id}
    ${access_key_id_base64}  Shell  base64 <<< ${access_key_id}
    input password  ${CONNECTIONS.KEYID}  ${access_key_id_base64.stdout}

Insert Access Key Secret
    [Arguments]  ${access_key_secret}
    ${access_key_secret_base64}  Shell  base64 <<< ${access_key_secret}
    input password  ${CONNECTIONS.KEYSECRET}  ${access_key_secret_base64.stdout}

Insert S3 Region
    [Arguments]  ${region}
    input text  ${CONNECTIONS.REGION}  ${region}

Insert ECR Region
    [Arguments]  ${region}
    input text  ${CONNECTIONS.REGION}  ${region}

# ===========================

Click "Next" button on "Specification" Step
    click button  ${CONNECTIONS.SPECIFICATION.NEXT_BUTTON}

Validate "Review" Step Content loaded
    page should contain element  ${CONNECTIONS.REVIEW.STEP_CONTENT}
    mouse over  ${CONNECTIONS.REVIEW.STEP_CONTENT}
    element should be visible  ${CONNECTIONS.REVIEW.STEP_CONTENT}

Click "Submit" button
    mouse over  ${CONNECTIONS.REVIEW.SUBMIT_BUTTON}
    capture page screenshot
    click button  ${CONNECTIONS.REVIEW.SUBMIT_BUTTON}

Validate Alert pops up and says connection created
    Validate visible element  ${CONNECTIONS.ALERT.SUCCESS}
    capture page screenshot
    element text should be  ${CONNECTIONS.ALERT.SUCCESS.CREATED_MESSAGE}  ${CONNECTIONS.CONNECTION_CREATED}

# Validate that connection created with specified values
#     [Arguments]

Fill Metadata during Connection creation
    [Arguments]  ${id}  ${type}  ${web_ui_link}=${EMPTY}  ${description}=${EMPTY}
    Fill "ID" field with value  ${id}
    Select Type  ${type}
    run keyword if  '${web_ui_link}' != '${EMPTY}'  Fill Web UI link field  ${web_ui_link}
    run keyword if  '${description}' != '${EMPTY}'  Add description for connection  ${description}
    Click "Next" button on "Metadata" Step

Fill "Git" Specification during Connection creation
    [Arguments]  ${git_ssh_url}  ${ssh_private_key}  ${reference}=${EMPTY}
    Validate "Specification" Step Content loaded
    Insert Git SSH URL  ${git_ssh_url}
    run keyword if  '${reference}' != '${EMPTY}'  Insert Reference  ${reference}
    Insert SSH private key  ${ssh_private_key}
    Click "Next" button on "Specification" Step

Fill "Docker" Specification during Connection creation
    [Arguments]  ${uri}  ${username}  ${password}
    Validate "Specification" Step Content loaded
    Insert URI  ${uri}
    Insert Username  ${username}
    Insert Password  ${password}
    Click "Next" button on "Specification" Step


Fill "GCS" Specification during Connection creation
    [Arguments]  ${project}  ${uri}  ${SA_secret}
    Validate "Specification" Step Content loaded
    Insert Project  ${project}
    Insert URI  ${uri}
    Insert Service account secret  ${SA_secret}
    Click "Next" button on "Specification" Step

Fill "Azureblob" Specification during Connection creation
    [Arguments]  ${uri}  ${SAS_Token}
    Validate "Specification" Step Content loaded
    Insert URI  ${uri}
    Insert SAS Token  ${SAS_Token}
    Click "Next" button on "Specification" Step

Fill "S3" Specification during Connection creation
    [Arguments]  ${uri}  ${region}  ${access_key_id}  ${access_key_secret}
    Validate "Specification" Step Content loaded
    Insert URI  ${uri}
    Insert S3 Region  ${region}
    Insert Access Key ID  ${access_key_id}
    Insert Access Key Secret  ${access_key_secret}
    Click "Next" button on "Specification" Step

Fill "ECR" Specification during Connection creation
    [Arguments]  ${uri}  ${region}  ${access_key_id}  ${access_key_secret}
    Validate "Specification" Step Content loaded
    Insert URI  ${uri}
    Insert ECR Region  ${region}
    Insert Access Key ID  ${access_key_id}
    Insert Access Key Secret  ${access_key_secret}
    Click "Next" button on "Specification" Step

Review Step during Connection creation
    # [Arguments]  ${id}  ${type}           ${git_ssh_url}  ${ssh_private_key}
    # ...          ${web_ui_link}=${EMPTY}  ${reference}=${EMPTY}
    Validate "Review" Step Content loaded
    capture page screenshot
    Click "Submit" button
