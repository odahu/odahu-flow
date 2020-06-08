*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}         odahuflow/config_training_training_data
# This value locates in the odahuflow/tests/stuf/data/odahuflow.project.yaml file.
${RUN_ID}    training_data_test
${TRAIN_ID}  test-downloading-training-data

*** Settings ***
Documentation       Check downloading of a training data
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Variables           ../../variables.py
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
Force Tags          training  training-data  cli

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose train delete --id ${TRAIN_ID} --ignore-not-found

Train valid model
    [Arguments]  ${training_file}
    Cleanup resources

    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/valid/${training_file} --id ${TRAIN_ID}
    ${res}=  StrictShell  odahuflowctl train get --id ${TRAIN_ID} -o 'jsonpath=$[0].status.artifacts[0].runId'
    should be equal  ${RUN_ID}  ${res.stdout}

Train invalid model
    [Arguments]  ${training_file}
    Cleanup resources

    ${res}=  Shell  odahuflowctl --verbose train create -f ${RES_DIR}/invalid/${training_file} --id ${TRAIN_ID}
    should not be equal  ${0}  ${res.rc}

Train model that create invalid GPPI artifact
    [Arguments]  ${training_file}
    Cleanup resources

    ${res}=  Shell  odahuflowctl --verbose train create -f ${RES_DIR}/invalid/${training_file} --id ${TRAIN_ID}
    should not be equal  ${0}  ${res.rc}
    Should contain  ${res.stdout}  ${GPPI_VALIDATION_FAIL}


*** Test Cases ***
Vaild data downloading parameters
    [Documentation]  Verify various valid combination of connection uri, remote path and local path parameters
    [Template]  Train valid model
    dir_to_dir.training.odahuflow.yaml
    remote_dir_to_dir.training.odahuflow.yaml
    file_to_file.training.odahuflow.yaml
    remote_file_to_file.training.odahuflow.yaml
    file_to_dir.training.odahuflow.yaml
    remote_file_to_dir.training.odahuflow.yaml

Invaild data downloading parameters
    [Documentation]  Verify various invalid combination of connection uri, remote path and local path parameters
    [Template]  Train invalid model
    not_found_file.training.odahuflow.yaml
    not_found_remote_file.training.odahuflow.yaml
    not_valid_dir_path.training.odahuflow.yaml
    not_valid_remote_dir.training.odahuflow.yaml

Invalid GPPI output
    [Documentation]  Verify that validator step checks GPPI correctness and raise error if GPPI is not correct
    [Template]  Train model that create invalid GPPI artifact
    invalid_gppi_artifact.training.odahuflow.yaml
