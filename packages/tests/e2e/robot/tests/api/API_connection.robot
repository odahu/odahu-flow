*** Variables ***
${RES_DIR}                          ${CURDIR}/resources/connection
${GIT_VALID}                        git-connection-valid
${DOCKER_VALID}                     docker-ci-connection-valid
${GIT_INVALID}                      git-connection-invalid
${CONN_SYNTAX_ERROR}                docker-ci-syntax-error

*** Settings ***
Documentation                       API of conections
Resource                            ../../resources/variables.robot
Resource                            ../../resources/keywords.robot
Resource                            ./resources/keywords.robot
Variables                           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library                             String
Library                             odahuflow.robot.libraries.sdk_wrapper
Library                             odahuflow.robot.libraries.sdk_wrapper.Connection
Suite Setup                         Run Keywords
...                                 Login to the api and edge
Force Tags                          api  sdk  connection
Test Timeout                        5 minutes

*** Keywords ***
Format WrongHttpStatusCode
    [Arguments]                     ${entity name}
    ${error output}                 format string  WrongHttpStatusCode: Got error from server: entity "${entity name}" is not found (status: 404)
    [return]                        ${error output}

Format EntityAlreadyExists
    [Arguments]                     ${entity name}
    ${error output}                 format string  EntityAlreadyExists: Got error from server: entity "${entity name}" already exists (status: 409)
    [return]                        ${error output}

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

Try Create Connection that already exists
    Call API                    connection post  ${RES_DIR}/valid/docker_connection_create.json
    ${EntityAlreadyExists}      Format EntityAlreadyExists  ${DOCKER_VALID}
    Call API and get Error      ${EntityAlreadyExists}  connection post  ${RES_DIR}/valid/docker_connection_create.json

Try Update not existing and deleted connection
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection put  ${RES_DIR}/invalid/git_connection_update.not_exist.yaml
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_VALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection put  ${RES_DIR}/invalid/git_connection_update.deleted.yaml

Try Get id not existing and deleted connection
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection get id  ${GIT_INVALID}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_VALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection get id  ${GIT_VALID}

Try Get id decrypted not existing and deleted connection
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection get id decrypted  ${GIT_INVALID}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_VALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection get id decrypted  ${GIT_VALID}

Try Delete not existing and deleted connection
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection delete  ${GIT_INVALID}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${GIT_VALID}
    Call API and get Error      ${WrongHttpStatusCode}  connection delete  ${GIT_VALID}
