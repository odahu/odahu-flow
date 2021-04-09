*** Settings ***
Documentation   tests for "Connections" page and CRUD for connection entities
Resource        ${RES_DIR}/common.robot
Suite Setup     Setup
Suite Teardown  Teardown
Test Setup      run keywords
...             Test Setup
...             AND  Go to "Connections" page
Test Teardown   Test Teardown
Force Tags      web-ui  connection

*** Variables ***
${RES_DIR}                          ${CURDIR}/resources
${UI.CONNECTIONS.GCS.URI}           gs://${FEEDBACK_BUCKET}/output
${UI.CONNECTIONS.DOCKER.WEB_LINK}   https://${ODAHU_UI_DOCKER_REPO}

${UI.CONNECTIONS.AZUREBLOB.SAS.MOCK}    d0phbHJYVXRuRkVNSS9LN01ERU5HL2JQeFJmaUNZRVhBTVBMRUtFWQ==
${UI.CONNECTIONS.S3.URI.MOCK}       s3://raw-data/model/input
${UI.CONNECTIONS.ECR.URI.MOCK}      5555555555.dkr.ecr.eu-central-1.amazonaws.com/odahuflow
${UI.CONNECTIONS.AWS.REGION.MOCK}       us-east-1
${UI.CONNECTIONS.AWS.KEYID.MOCK}        QUtJQUlPU0ZPRE5ON0VYQU1QTEU=
${UI.CONNECTIONS.AWS.KEYSECRET.MOCK}    d0phbHJYVXRuRkVNSS9LN01ERU5HL2JQeFJmaUNZRVhBTVBMRUtFWQ==

${COMMON_CONNECTION_NAME}   odahu-ui-connection
${GIT_TYPE}     git
${GCS_TYPE}     gcs
${DOCKER_TYPE}  docker
${AZUREBLOB_TYPE}  azureblob
${S3_TYPE}      s3
${ECR_TYPE}     ecr

*** Test Cases ***
Create "Git" Connection and validate that one exists
    Create "Git" Connection  id=${COMMON_CONNECTION_NAME}-${GIT_TYPE}  type=${GIT_TYPE}  web_ui_link=${ODAHU_UI_GIT_WEB_LINK}
    ...  description=Git repository with the Odahu-Flow examples for tests
    ...  git_ssh_url=${ODAHU_UI_GIT_URI}  reference=${ODAHU_UI_GIT_REFERENCE}  ssh_private_key=${ODAHU_UI_GIT_KEYSECRET}

Create "Docker" Connection and validate that one exists
    Create "Docker" Connection  id=${COMMON_CONNECTION_NAME}-${DOCKER_TYPE}  type=${DOCKER_TYPE}
    ...  web_ui_link=${UI.CONNECTIONS.DOCKER.WEB_LINK}  description=Tests GCR docker repository for model packaging
    ...  uri=${ODAHU_UI_DOCKER_REPO}  username=${ODAHU_UI_DOCKER_USERNAME}  password=${ODAHU_UI_DOCKER_PASSWORD}

Create "GCS" Connection and validate that one exists
    Create "GCS" Connection  id=${COMMON_CONNECTION_NAME}-${GCS_TYPE}  type=${GCS_TYPE}  description=Connection for storage of trained artifacts
    ...  project=${ODAHU_UI_GCS_PROJECT}  uri=${UI.CONNECTIONS.GCS.URI}  SA_secret=${ODAHU_UI_GCS_SA_SECRET}

Create "Azureblob" Connection and validate that one exists
    Create "Azureblob" Connection  id=${COMMON_CONNECTION_NAME}-${AZUREBLOB_TYPE}  type=${AZUREBLOB_TYPE}
    ...  web_ui_link=https//${FEEDBACK_BUCKET}  description=Tests GCR docker repository for model packaging
    ...  uri=${FEEDBACK_BUCKET}  SAS_Token=${UI.CONNECTIONS.AZUREBLOB.SAS.MOCK}

Create "S3" Connection and validate that one exists
    Create "S3" Connection  id=${COMMON_CONNECTION_NAME}-${S3_TYPE}  type=${S3_TYPE}
    ...  description=Tests S3 bucket storage for model trained artifacts
    ...  uri=${UI.CONNECTIONS.S3.URI.MOCK}  region=${UI.CONNECTIONS.AWS.REGION.MOCK}
    ...  access_key_id=${UI.CONNECTIONS.AWS.KEYID.MOCK}  access_key_secret=${UI.CONNECTIONS.AWS.KEYSECRET.MOCK}

Create "ECR" Connection and validate that one exists
    Create "ECR" Connection  id=${COMMON_CONNECTION_NAME}-${ECR_TYPE}  type=${ECR_TYPE}
    ...  description=Tests ECR docker repository for model packaging
    ...  uri=${UI.CONNECTIONS.ECR.URI.MOCK}  region=${UI.CONNECTIONS.AWS.REGION.MOCK}
    ...  access_key_id=${UI.CONNECTIONS.AWS.KEYID.MOCK}  access_key_secret=${UI.CONNECTIONS.AWS.KEYSECRET.MOCK}

Open View of connections and validate values
    [Template]  Open View of the connection
    ${COMMON_CONNECTION_NAME}-${GIT_TYPE}
    ${COMMON_CONNECTION_NAME}-${GCS_TYPE}
    ${COMMON_CONNECTION_NAME}-${DOCKER_TYPE}
    ${COMMON_CONNECTION_NAME}-${AZUREBLOB_TYPE}
    ${COMMON_CONNECTION_NAME}-${S3_TYPE}
    ${COMMON_CONNECTION_NAME}-${ECR_TYPE}
