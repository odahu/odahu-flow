*** Variables ***
${RES_DIR}             ${CURDIR}/resources/packager
${DOCKER_CLI}          docker-cli-api-testing
${DOCKER_REST}         docker-rest-api-testing

*** Settings ***
Documentation       API of conections
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.Packager
Suite Setup         Run Keywords
...                 Login to the api and edge
Force Tags          api  testing  packager


*** Test Cases ***
Get list of packagers
    [Documentation]  check that packagers that would be created do not exist now
    Command response list should not contain id  packager  ${DOCKER_CLI}  ${DOCKER_REST}

Create Docker CLI packager
    Call API                    packager post  ${RES_DIR}/valid/docker_cli_create.yaml

Create Docker REST packager
    Call API                    packager post  ${RES_DIR}/valid/docker_rest_create.json

Update Docker CLI packager
    Call API                    packager put  ${RES_DIR}/valid/docker_cli_update.json

Update Docker REST packager
    Call API                    packager put  ${RES_DIR}/valid/docker_rest_update.yaml

Get updated list of packagers
    Command response list should contain id  packager  ${DOCKER_CLI}  ${DOCKER_REST}

Get Docker CLI and REST packagers by id
    Call API                    packager get id  ${DOCKER_CLI}
    Call API                    packager get id  ${DOCKER_REST}

Delete Docker CLI packager
    Call API                    packager delete  ${DOCKER_CLI}

Delete Docker REST packager
    Call API                    packager delete  ${DOCKER_REST}

Check that packagers do not exist
    Command response list should not contain id  packager  ${DOCKER_CLI}  ${DOCKER_REST}
