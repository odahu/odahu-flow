*** Variables ***
${RES_DIR}             ${CURDIR}/resources/connections
${GIT_CONN}            git-connection-valid
${DOCKER_CONN}         docker-ci-connectio-valid

*** Settings ***
Documentation       API of conections
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Connection
Suite Setup         Run Keywords
...                 Login to the api and edge
#...                Cleanup All Resources (to be created)
# Suite Teardown    Run Keywords
#...                Cleanup All Resources (to be created)
# Test Setup
# Test Teardown
Force Tags          api  connection


*** Test Cases ***
Get empty list of connections
    ${result}                   Call API  connection get
    should be empty             ${result}

Create GIT connection
    [Documentation]  create git connection and check that one exists
    Call API                    connection post  ${RES_DIR}/valid/git_connection_create.yaml
    ${result}                   Call API  connection get
    length should be            ${result}  1
    ${result_id}                Log id  @{result}
    should be equal             ${result_id}  ${GIT_CONN}

Create Docker connection
    [Documentation]  create docker connection and check that one exists
    Call API                    connection post  ${RES_DIR}/valid/docker_connection_create.json
    ${result}                   Call API  connection get
    length should be            ${result}  2
    Command response list should contain id  connection  ${DOCKER_CONN}

Update GIT connection
    [Documentation]  updating GIT connection
    sleep                       1s  # waiting to have, for sure, different time between createdAt and updatedAt
    Call API                    connection put  ${RES_DIR}/valid/git_connection_update.yaml
    ${result}                   Call API  connection get id  ${GIT_CONN}
    ${result_id}                Log id  ${result}
    should be equal             ${result_id}  ${GIT_CONN}
    keySecret connection should be equal  ${result}  *****
    ${result_status}            Log Status  ${result}
    should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Get GIT connection
    [Documentation]  getting GIT connection
    Call API                    connection get id  ${GIT_CONN}


Put Docker connection
    [Documentation]  update Docker connection
    Call API                    connection put  ${RES_DIR}/valid/docker_connection_create.json
    ${result}                   Call API  connection get id decrypted  ${DOCKER_CONN}
    should be equal             ${result.id}  ${DOCKER_CONN}

Get Decrypted Docker connection
    [Documentation]  get decrypted Docker connection
    ${result}                   Call API  connection get id decrypted  ${DOCKER_CONN}
    Password connection should not be equal    ${result}  *****

Delete GIT connection
    Call API                    connection delete  ${GIT_CONN}
    ${result}                   Call API  connection get
    length should be            ${result}  1

Delete Docker connection
    Call API                    connection delete  ${DOCKER_CONN}
    ${result}                   Call API  connection get
    length should be            ${result}  0


