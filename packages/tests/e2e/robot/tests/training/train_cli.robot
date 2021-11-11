*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${SETUP_DIR}            ${CURDIR}/setup
${LOCAL_CONFIG}         odahuflow/config_training_training_cli
${TRAIN_ID}             test-algorithm-source
${TRAIN_STUFF_DIR}      ${CURDIR}/../../../../stuff


*** Settings ***
Documentation       Check training model via cli with various algorithm sources
Test Timeout        60 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Variables           ../../variables.py
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.odahu_k8s_reporter.OdahuKubeReporter
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
...                 AND  Login to the api and edge
...                 AND  Cleanup all resources
...                 AND  StrictShell  ${SETUP_DIR}/train_setup.sh
Suite Teardown      Run Keywords
...                 Cleanup all resources  AND
...                 Remove file  ${LOCAL_CONFIG}
Force Tags          training  algorithm-source  cli

*** Keywords ***
Cleanup all resources
    [Documentation]  cleanups resources created during whole test suite, hardcoded training IDs
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-vcs
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-object-storage

Cleanup resources
    [Arguments]  ${training id}
    StrictShell  odahuflowctl --verbose train delete --id ${training id} --ignore-not-found

Train valid model
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/valid/${training_file} --id ${training id}
    report training pods  ${training id}

*** Test Cases ***
Vaild downloading parameters
    [Documentation]  Verify valid algorithm sourcses
    [Template]  Train valid model
    ${TRAIN_ID}-vcs                   vcs.training.odahuflow.yaml
    ${TRAIN_ID}-object-storage        object_storage.training.odahuflow.yaml
