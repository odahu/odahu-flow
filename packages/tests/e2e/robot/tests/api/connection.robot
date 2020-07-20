*** Variables ***
${RES_DIR}             ${CURDIR}/resources/connection
${GIT_CONN}            git-connection-valid
${DOCKER_CONN}         docker-ci-connection-valid

*** Settings ***
Documentation       API of conections
Resource            ../../resources/variables.robot
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Connection
Suite Setup         Run Keywords
...                 Login to the api and edge
Force Tags          api  sdk  connection


*** Test Cases ***
Create GIT connection
    [Documentation]  create git connection and check that one exists
    Call API                    connection post  ${RES_DIR}/valid/git_connection_create.yaml
    Command response list should contain id  connection  ${GIT_CONN}

Create Docker connection
    [Documentation]  create docker connection and check that one exists
    Call API                    connection post  ${RES_DIR}/valid/docker_connection_create.json
    Command response list should contain id  connection  ${DOCKER_CONN}

Update GIT connection
    [Documentation]  updating GIT connection
    sleep                       1s  # waiting to have, for sure, different time between createdAt and updatedAt
    Call API                    connection put  ${RES_DIR}/valid/git_connection_update.yaml
    ${result}                   Call API  connection get id  ${GIT_CONN}
    ${result_status}            Log Status  ${result}
    should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Get GIT connection
    [Documentation]  getting GIT connection
    ${result}                   Call API  connection get id  ${GIT_CONN}
    keySecret connection should be equal  ${result}  ${CONN_SECRET_MASK}

Update Docker connection
    [Documentation]  update Docker connection
    Call API                    connection put  ${RES_DIR}/valid/docker_connection_create.json
    ${result}                   Call API  connection get id  ${DOCKER_CONN}
    should be equal             ${result.id}  ${DOCKER_CONN}

Get Decrypted Docker connection
    [Documentation]  get decrypted Docker connection
    ${result}                   Call API  connection get id decrypted  ${DOCKER_CONN}
    Password connection should not be equal    ${result}  ${CONN_SECRET_MASK}

Delete GIT connection
    Call API                    connection delete  ${GIT_CONN}
    Command response list should not contain id  connection  ${GIT_CONN}

Delete Docker connection
    Call API                    connection delete  ${DOCKER_CONN}
    Command response list should not contain id  connection  ${DOCKER_CONN}
