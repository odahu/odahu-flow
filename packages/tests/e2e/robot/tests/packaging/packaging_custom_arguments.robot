*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${PACK_ID}              test-custom-arguments-pack
${LOCAL_CONFIG}         odahuflow/config_packaging_custom_args

*** Settings ***
Documentation       Check packaging arguments and target
Test Timeout        60 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup resources
Suite Teardown      Run Keywords
...                 Cleanup resources  AND
...                 Remove file  ${LOCAL_CONFIG}
Test Teardown       Cleanup resources
Force Tags          packaging

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose pack delete --id ${PACK_ID} --ignore-not-found

*** Test Cases ***
Custom image name for docker-rest packager
    StrictShell  odahuflowctl --verbose pack create -f ${RES_DIR}/arguments/image_name_rest.yaml
    ${res}=  StrictShell  odahuflowctl pack get --id ${PACK_ID} -o 'jsonpath=$[0].status.results[0].value'

    Should contain  ${res.stdout}  simple-model:1.2

Custom image name for docker-cli packager
    StrictShell  odahuflowctl --verbose pack create -f ${RES_DIR}/arguments/image_name_cli.yaml
    ${res}=  StrictShell  odahuflowctl pack get --id ${PACK_ID} -o 'jsonpath=$[0].status.results[0].value'

    Should contain  ${res.stdout}  simple-model:1.2

Validate arguments for the docker-rest packager
    ${res}=  Shell  odahuflowctl --verbose pack create -f ${RES_DIR}/arguments/invalid_argument_rest.yaml
    Should not be equal  ${res.rc}  ${0}

Validate arguments for the docker-cli packager
    ${res}=  Shell  odahuflowctl --verbose pack create -f ${RES_DIR}/arguments/invalid_argument_cli.yaml
    Should not be equal  ${res.rc}  ${0}

Validate a target for the docker-rest packager
    ${res}=  Shell  odahuflowctl --verbose pack create -f ${RES_DIR}/arguments/invalid_target_rest.yaml
    Should not be equal  ${res.rc}  ${0}

Validate a target for the docker-cli packager
    ${res}=  Shell  odahuflowctl --verbose pack create -f ${RES_DIR}/arguments/invalid_target_cli.yaml
    Should not be equal  ${res.rc}  ${0}