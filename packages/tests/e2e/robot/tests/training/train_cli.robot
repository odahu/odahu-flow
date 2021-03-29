*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}         odahuflow/config_training_training_data
${TRAIN_ID}             test-algorithm-source
${TRAIN_STUFF_DIR}      ../../../../stuff


*** Settings ***
Documentation       Check training model via cli with various algorithm sources
Test Timeout        20 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Variables           ../../variables.py
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.odahu_k8s_reporter.OdahuKubeReporter
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup all resources
Suite Teardown      Run Keywords
...                 Cleanup all resources  AND
...                 Remove file  ${LOCAL_CONFIG}
Force Tags          training  algorithm-source cli

*** Keywords ***
Cleanup all resources
    [Documentation]  cleanups resources created during whole test suite, hardcoded training IDs
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-vcs
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-object-storage
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-non-exist-vcs
    StrictShell  odahuflowctl --verbose train delete --ignore-not-found --id ${TRAIN_ID}-non-exist-storage

Cleanup resources
    [Arguments]  ${training id}
    StrictShell  odahuflowctl --verbose train delete --id ${training id} --ignore-not-found

Train valid model from VCS
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    ${res}=  StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/valid/${training_file} --id ${training id}
    report training pods  ${training id}
    should be equal ${res.rc}  ${0}

Train invalid model from VCS
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    ${res}=  StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/invalid/${training_file} --id ${training id}
    report training pods  ${training id}
    should not be equal ${res.rc}  ${0}

Train valid model from object storage
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    Download file  mlflow/sklearn/wine/MLproject  ${RES_DIR}/algorithm_source/MLproject
    StrictShell  ${TRAIN_STUFF_DIR}/training_stuff.sh bucket-copy "${RES_DIR}/algorithm_source" "/algorithm_source"
    ${res}=  StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/valid/${training_file} --id ${training id}
    report training pods  ${training id}
    should be equal ${res.rc}  ${0}

Train invalid model from object storage
    [Arguments]  ${training id}  ${training_file}
    [Teardown]  Cleanup resources  ${training id}
    Download file  mlflow/sklearn/wine/MLproject  ${RES_DIR}/algorithm_source/MLproject
    StrictShell  ${TRAIN_STUFF_DIR}/training_stuff.sh bucket-copy "${RES_DIR}/algorithm_source" "/algorithm_source"
    ${res}=  StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/invalid/${training_file} --id ${training id}
    report training pods  ${training id}
    should not be equal ${res.rc}  ${0}

*** Test Cases ***
Vaild VCS downloading parameters
    [Documentation]  Verify valid VCS sourcses
    [Template]  Train valid model from VCS
    ${TRAIN_ID}-vcs                   vcs.training.odahuflow.yaml

Invaild VCS downloading parameters
    [Documentation]  Verify invalid VCS sourcses
    [Template]  Train invalid model from VCS
    ${TRAIN_ID}-non-exist-vcs         non_exist_vcs.training.odahuflow.yaml

Vaild object storage downloading parameters
    [Documentation]  Verify valid object storage sourcses
    [Template]  Train valid model from object storage
    ${TRAIN_ID}-object-storage        object_storage.training.odahuflow.yaml

Invaild object storage downloading parameters
    [Documentation]  Verify invalid object storage sourcses
    [Template]  Train invalid model from object storage
    ${TRAIN_ID}-non-exist-storage     non_exist_object_storage.training.odahuflow.yaml