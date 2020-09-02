*** Variables ***
${LOCAL_CONFIG}     odahuflow/api_packager
${RES_DIR}          ${CURDIR}/resources/packager
${DOCKER_CLI}       docker-cli-api-testing
${DOCKER_REST}      docker-rest-api-testing
${DOCKER_INVALID}   docker-rest-api-not-exist

*** Settings ***
Documentation       API of packagers
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Packager
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup All Resources
Suite Teardown      Run Keywords
...                 Cleanup All Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  packager
Test Timeout        5 minutes

*** Keywords ***
Cleanup All Resources
    Cleanup resource  packaging-integration  ${DOCKER_CLI}
    Cleanup resource  packaging-integration  ${DOCKER_REST}
    Cleanup resource  packaging-integration  ${DOCKER_INVALID}

*** Test Cases ***
Get list of packagers
    [Documentation]  check that packagers that would be created do not exist now
    Command response list should not contain id  packager  ${DOCKER_CLI}  ${DOCKER_REST}

Create Docker CLI packager
    Call API                    packager post  ${RES_DIR}/valid/docker_cli_create.yaml
    ${check}                    Call API  packager get id  ${DOCKER_CLI}
    Default Docker image should be equal  ${check}  created
    Default Entrypoint should be equal  ${check}  created

Create Docker REST packager
    Call API                    packager post  ${RES_DIR}/valid/docker_rest_create.json
    ${check}                    Call API  packager get id  ${DOCKER_REST}
    Default Docker image should be equal  ${check}  created
    Default Entrypoint should be equal  ${check}  created

Update Docker CLI packager
    sleep                       1s
    Call API                    packager put  ${RES_DIR}/valid/docker_cli_update.json
    ${check}                    Call API  packager get id  ${DOCKER_CLI}
    Default Docker image should be equal  ${check}  updated
    Default Entrypoint should be equal  ${check}  updated

Update Docker REST packager
    Call API                    packager put  ${RES_DIR}/valid/docker_rest_update.yaml
    ${check}                    Call API  packager get id  ${DOCKER_REST}
    Default Docker image should be equal  ${check}  updated
    Default Entrypoint should be equal  ${check}  updated

Get updated list of packagers
    Command response list should contain id  packager  ${DOCKER_CLI}  ${DOCKER_REST}

Get Docker CLI and REST packagers by id
    ${result}                   Call API  packager get id  ${DOCKER_CLI}
    ID should be equal          ${result}  ${DOCKER_CLI}
    ${result}                   Call API  packager get id  ${DOCKER_REST}
    ID should be equal          ${result}  ${DOCKER_REST}

Delete Docker CLI packager
    ${result}                   Call API   packager delete  ${DOCKER_CLI}
    should be equal             ${result.get('message')}  Packaging integration ${DOCKER_CLI} was deleted

Delete Docker REST packager
    ${result}                   Call API  packager delete  ${DOCKER_REST}
    should be equal             ${result.get('message')}  Packaging integration ${DOCKER_REST} was deleted

Check that packagers do not exist
    Command response list should not contain id  packager  ${DOCKER_CLI}  ${DOCKER_REST}

#############################
#    NEGATIVE TEST CASES    #
#############################
Try Create Packager that already exists
    [Tags]                      negative
    [Teardown]                  Cleanup resource  packaging-integration  ${DOCKER_CLI}
    Call API                    packager post  ${RES_DIR}/valid/docker_cli_create.yaml
    ${EntityAlreadyExists}      Format EntityAlreadyExists  ${DOCKER_CLI}
    Call API and get Error      ${EntityAlreadyExists}  packager post  ${RES_DIR}/valid/docker_cli_create.yaml

Try Update not existing and deleted Packager
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DOCKER_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  packager put  ${RES_DIR}/invalid/docker_rest_update.not_exist.json
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DOCKER_CLI}
    Call API and get Error      ${WrongHttpStatusCode}  packager put  ${RES_DIR}/valid/docker_cli_update.json

Try Get id not existing and deleted Packager
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DOCKER_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  packager get id  ${DOCKER_INVALID}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DOCKER_REST}
    Call API and get Error      ${WrongHttpStatusCode}  packager get id  ${DOCKER_REST}

Try Delete not existing and deleted Packager
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DOCKER_INVALID}
    Call API and get Error      ${WrongHttpStatusCode}  packager delete  ${DOCKER_INVALID}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DOCKER_CLI}
    Call API and get Error      ${WrongHttpStatusCode}  packager delete  ${DOCKER_CLI}
