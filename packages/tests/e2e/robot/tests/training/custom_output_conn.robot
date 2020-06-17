*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${TRAIN_ID}  test-custom-output-connection-training
${PACK_ID}  test-custom-output-connection-pack
${LOCAL_CONFIG}         odahuflow/config_custom_output_con

*** Settings ***
Documentation       Check using custom output connection
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
Force Tags          training  packaging

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose train delete --id ${TRAIN_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${PACK_ID} --ignore-not-found

*** Test Cases ***
Use user defined connection for trained model binary storage
    [Documentation]    This test emulate user who use `OutputConnection` ModelTraining and ModelPackaging parameter
    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/custom-con/training.yaml
    ${res}=  StrictShell  odahuflowctl train get --id ${TRAIN_ID} -o 'jsonpath=$[0].status.artifacts[0].artifactName'

    StrictShell  odahuflowctl --verbose pack create -f ${RES_DIR}/custom-con/pack.yaml --artifact-name ${res.stdout}
    ${res}=  StrictShell  odahuflowctl pack get --id ${PACK_ID} -o 'jsonpath=$[0].status.results[0].value'