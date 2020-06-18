*** Variables ***
${RES_DIR}             ${CURDIR}/resources/connections

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
Force Tags          api  testing


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
    Should Be Equal             ${result_id}  git-connection-valid

# Create Docker connection
#     [Documentation]  create docker connection and check that one exists
#     Call API                    connection post  ${RES_DIR}/valid/docker_connection_create.json
#     ${result}                   Call API  connection get
#     length should be            ${result}  2
#     list should contain value   ${result}  docker-ci-connectio-valid

Update GIT connection
    [Documentation]  updating GIT connection
    sleep                       1s  # waiting to have, for sure, different time between createdAt and updatedAt
    Call API                    connection put  ${RES_DIR}/valid/git_connection_update.yaml
    ${result}                   Call API  connection get id  git-connection-valid
    ${result_id}                Log id  ${result}
    Should Be Equal             ${result_id}  git-connection-valid
    ${result_status}            Log Status  ${result}
    Should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Put and Get Decrypted Docker connection
    [Documentation]  update and get decrypted Docker connection
    Call API                    connection put  ${RES_DIR}/valid/docker_connection_create.json
    ${result}                   Call API  connection get id decrypted  docker-ci-connectio-valid
    Log                         ${result.spec.password}
    should not be equal         ${result.spec.password}  *****

Delete GIT connection
    Call API                    connection delete  git-connection-valid
    ${result}                   Call API  connection get
    length should be            ${result}  1

Delete Docker connection
    Call API                    connection delete  docker-ci-connectio-valid
    ${result}                   Call API  connection get
    length should be            ${result}  0


