*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}         odahuflow/config_training_training_data
# This value locates in the odahuflow/tests/stuf/data/odahuflow.project.yaml file.
${RUN_ID}    training_data_test
${TRAIN_ID}  test-training-data

*** Settings ***
Documentation       Check downloading of a training data
Test Timeout        60 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Variables           ../../variables.py
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.odahu_k8s_reporter.OdahuKubeReporter
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup all resources
Suite Teardown      Run Keywords
...                 Cleanup all resources  AND
...                 Remove file  ${LOCAL_CONFIG}
Force Tags          training  training-data

*** Keywords ***
Cleanup all resources
    [Documentation]  cleanups resources created during whole test suite, hardcoded training IDs
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-dir-to-dir
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-remote-dir-to-dir
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-file-to-file
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-remote-file-to-file
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-file-to-dir
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-remote-file-to-dir
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-not-found-file
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-not-found-remote-file
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-not-valid-dir-path
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-not-valid-remote-dir
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-invalid-gppi

Cleanup resources
    [Arguments]  ${training id}
    StrictShell  odahuflowctl --verbose train delete --id ${training id} --ignore-not-found

Train valid model
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/valid/${training_file} --id ${training id}
    report training pods  ${training id}
    ${res}=  StrictShell  odahuflowctl train get --id ${training id} -o 'jsonpath=$[0].status.artifacts[0].runId'
    should be equal  ${RUN_ID}  ${res.stdout}

Train invalid model
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    ${res}=  Shell  odahuflowctl --verbose train create -f ${RES_DIR}/invalid/${training_file} --id ${training id}
    report training pods  ${training id}
    should not be equal  ${0}  ${res.rc}

Train model that create invalid GPPI artifact
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    ${res}=  Shell  odahuflowctl --verbose train create -f ${RES_DIR}/invalid/${training_file} --id ${training id}
    report training pods  ${training id}
    should not be equal  ${0}  ${res.rc}
    Should contain  ${res.stdout}  ${GPPI_VALIDATION_FAIL}


*** Test Cases ***
Vaild data downloading parameters
    [Documentation]  Verify various valid combination of connection uri, remote path and local path parameters
    [Template]  Train valid model
    ${TRAIN_ID}-dir-to-dir                  dir_to_dir.training.odahuflow.yaml
    ${TRAIN_ID}-remote-dir-to-dir           remote_dir_to_dir.training.odahuflow.yaml
    ${TRAIN_ID}-file-to-file                file_to_file.training.odahuflow.yaml
    ${TRAIN_ID}-remote-file-to-file         remote_file_to_file.training.odahuflow.yaml
    ${TRAIN_ID}-file-to-dir                 file_to_dir.training.odahuflow.yaml
    ${TRAIN_ID}-remote-file-to-dir          remote_file_to_dir.training.odahuflow.yaml

Invaild data downloading parameters
    [Documentation]  Verify various invalid combination of connection uri, remote path and local path parameters
    [Template]  Train invalid model
    ${TRAIN_ID}-not-found-file              not_found_file.training.odahuflow.yaml
    ${TRAIN_ID}-not-found-remote-file       not_found_remote_file.training.odahuflow.yaml
    ${TRAIN_ID}-not-valid-dir-path          not_valid_dir_path.training.odahuflow.yaml
    ${TRAIN_ID}-not-valid-remote-dir        not_valid_remote_dir.training.odahuflow.yaml

