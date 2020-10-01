*** Variables ***
${LOCAL_CONFIG}        odahuflow/api_connection
${RES_DIR}             ${CURDIR}/resources/connection
${GIT_VALID}           git-connection-valid
${DOCKER_VALID}        docker-ci-connection-valid
${GIT_NOT_EXIST}       git-connection-not-exist

*** Settings ***
Documentation          API of conections
Resource               ../../resources/variables.robot
Resource               ../../resources/keywords.robot
Resource               ./resources/keywords.robot
Variables              ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library                String
Library                odahuflow.robot.libraries.sdk_wrapper.Connection
Suite Setup            Run Keywords
...                    Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                    Login to the api and edge  AND
...                    Cleanup All Resources
Suite Teardown         Run Keywords
...                    Cleanup All Resources  AND
...                    Remove File  ${LOCAL_CONFIG}
Force Tags             api  sdk  connection
Test Timeout           5 minutes

*** Keywords ***
Cleanup All Resources
    Cleanup resource  connection  ${GIT_VALID}
    Cleanup resource  connection  ${DOCKER_VALID}
    Cleanup resource  connection  ${GIT_NOT_EXIST}

*** Test Cases ***
Create GIT connection
    [Documentation]  create git connection and check that one exists
    Call API                    connection post  ${RES_DIR}/valid/git_connection_create.yaml
    Command response list should contain id  connection  ${GIT_VALID}

Create Docker connection
    [Documentation]  create docker connection and check that one exists
    Call API                    connection post  ${RES_DIR}/valid/docker_connection_create.json
    Command response list should contain id  connection  ${DOCKER_VALID}

Update GIT connection
    [Documentation]  updating GIT connection
    sleep                       1s  # waiting to have, for sure, different time between createdAt and updatedAt
    Call API                    connection put  ${RES_DIR}/valid/git_connection_update.yaml
    ${result}                   Call API  connection get id  ${GIT_VALID}
    should be equal             ${result.spec.description}  link to the Git repo odahu-flow-examples
    CreatedAt and UpdatedAt times should not be equal  ${result}

Get GIT connection
    [Documentation]  getting GIT connection
    ${result}                   Call API  connection get id  ${GIT_VALID}
    keySecret connection should be equal  ${result}  ${CONN_SECRET_MASK}

Update Docker connection
    [Documentation]  update Docker connection
    Call API                    connection put  ${RES_DIR}/valid/docker_connection_update.json
    ${result}                   Call API  connection get id  ${DOCKER_VALID}
    ID should be equal          ${result}  ${DOCKER_VALID}
    should be equal             ${result.spec.description}  updated
    CreatedAt and UpdatedAt times should not be equal  ${result}

Get Decrypted Docker connection
    [Documentation]  get decrypted Docker connection
    ${result}                   Call API  connection get id decrypted  ${DOCKER_VALID}
    Password connection should not be equal    ${result}  ${CONN_SECRET_MASK}

Delete GIT connection
    Call API                    connection delete  ${GIT_VALID}
    Command response list should not contain id  connection  ${GIT_VALID}

Delete Docker connection
    Call API                    connection delete  ${DOCKER_VALID}
    Command response list should not contain id  connection  ${DOCKER_VALID}

#############################
#    NEGATIVE TEST CASES    #
#############################
Try Create Connection that already exists
    [Tags]                      negative
    [Setup]                     cleanup resource  connection  ${DOCKER_VALID}
    [Teardown]                  cleanup resource  connection  ${DOCKER_VALID}
    Call API                    connection post  ${RES_DIR}/valid/docker_connection_create.json
    ${EntityAlreadyExists}      format string  ${409 Conflict Template}  ${DOCKER_VALID}
    Call API and get Error      ${EntityAlreadyExists}  connection post  ${RES_DIR}/valid/docker_connection_create.json

Try Update not existing Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  connection put  ${RES_DIR}/invalid/git_connection_update.not_exist.yaml

Try Update deleted Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_VALID}
    Call API and get Error      ${404NotFound}  connection put  ${RES_DIR}/valid/git_connection_update.yaml

Try Get id not existing Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  connection get id  ${GIT_NOT_EXIST}

Try Get id deleted Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_VALID}
    Call API and get Error      ${404NotFound}  connection get id  ${GIT_VALID}

Try Get id decrypted not existing Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  connection get id decrypted  ${GIT_NOT_EXIST}

Try Get id decrypted deleted Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_VALID}
    Call API and get Error      ${404NotFound}  connection get id decrypted  ${GIT_VALID}

Try Delete not existing Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  connection delete  ${GIT_NOT_EXIST}

Try Delete deleted Connection
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${GIT_VALID}
    Call API and get Error      ${404NotFound}  connection delete  ${GIT_VALID}
