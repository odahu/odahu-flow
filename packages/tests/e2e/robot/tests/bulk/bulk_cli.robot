*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}         odahuflow/config_bulk_cli
${CONN_1_ID}            bulk-test-conn-1
${CONN_2_ID}            bulk-test-conn-2
${TI_1_ID}              bulk-test-ti-1
${TI_2_ID}              bulk-test-ti-2
${PI_1_ID}              bulk-test-pi-1
${TRAINING_1_NAME}      bulk-test-mt-1
${PACKAGING_1_NAME}     bulk-test-mp-1
${PACKAGING_2_NAME}     bulk-test-mp-2

*** Settings ***
Documentation       OdahuFlow's API operational check for bulk operations (with multiple resources)
Test Timeout        20 minutes
Resource            ../../resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Default Tags        cli  bulk
Suite Setup         Run keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
...                 AND  Login to the api and edge
...                 AND  Cleanup resources
Suite Teardown      Run keywords
...                 Cleanup resources
...                 AND  Remove File  ${LOCAL_CONFIG}

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose conn delete --id ${CONN_1_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose conn delete --id ${CONN_2_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose ti delete --id ${TI_1_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose ti delete --id ${TI_2_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose pi delete --id ${PI_1_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose train delete --id ${TRAINING_1_NAME} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${PACKAGING_1_NAME} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${PACKAGING_2_NAME} --ignore-not-found

Check entity exists
    [Arguments]  ${entity_type}  ${name}
    ${res}=  StrictShell  odahuflowctl --verbose ${entity_type} get --id ${name}
        Should contain     ${res.stderr}  ${name}

Check entity doesn't exist
    [Arguments]  ${entity_type}  ${name}
    ${res}=  FailedShell  odahuflowctl --verbose ${entity_type} get --id ${name}
        Should contain     ${res.stderr}  "${name}" is not found

Apply bulk file and check counters
    [Arguments]  ${file}  ${created}  ${changed}  ${deleted}
    ${res}=  Shell  odahuflowctl --verbose bulk apply ${RES_DIR}/${file}
             Should be equal  ${res.rc}      ${0}
             Run Keyword If   ${created} > 0  Should contain   ${res.stdout}  created resources: ${created}
             Run Keyword If   ${changed} > 0  Should contain   ${res.stdout}  changed resources: ${changed}
             Run Keyword If   ${deleted} > 0  Should contain   ${res.stdout}  deleted resources: ${deleted}

Apply bulk file and check errors
    [Arguments]  ${file}  ${expected_error_message}
    ${res}=  Shell  odahuflowctl --verbose bulk apply ${RES_DIR}/${file}
             Should not be equal      ${res.rc}      ${0}
             Should contain       ${res.stdout}  ${expected_error_message}

Apply bulk file and check parse errors
    [Arguments]  ${file}  ${expected_error_message}
    ${res}=  Shell  odahuflowctl --verbose bulk apply ${RES_DIR}/${file}
             Should not be equal      ${res.rc}      ${0}
             Should contain       ${res.stderr}  ${expected_error_message}

Remove bulk file and check counters
    [Arguments]  ${file}  ${created}  ${changed}  ${deleted}
    ${res}=  Shell  odahuflowctl --verbose bulk delete ${RES_DIR}/${file}
             Should be equal  ${res.rc}      ${0}
             Run Keyword If   ${created} > 0  Should contain   ${res.stdout}  created resources: ${created}
             Run Keyword If   ${changed} > 0  Should contain   ${res.stdout}  changed resources: ${changed}
             Run Keyword If   ${deleted} > 0  Should contain   ${res.stdout}  removed resources: ${deleted}

Template. Apply good profile, check resources and remove on teardown
    [Arguments]  ${file}
    [Teardown]  Cleanup resources
    Check entity doesn't exist      connection  ${CONN_1_ID}
    Check entity doesn't exist      connection  ${CONN_2_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_1_ID}
    Check entity doesn't exist      packaging-integration  ${PI_1_ID}
    Check entity doesn't exist      training  ${TRAINING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_1_NAME}
    Apply bulk file and check counters     ${file}  6  0  0
    Check entity exists             connection  ${CONN_1_ID}
    Check entity exists             connection  ${CONN_2_ID}
    Check entity exists             toolchain-integration  ${TI_1_ID}
    Check entity exists             packaging-integration  ${PI_1_ID}
    Check entity exists             training  ${TRAINING_1_NAME}
    Check entity exists             packaging  ${PACKAGING_1_NAME}
    Remove bulk file and check counters    ${file}  0  0  6
    Check entity doesn't exist      connection  ${CONN_1_ID}
    Check entity doesn't exist      connection  ${CONN_2_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_1_ID}
    Check entity doesn't exist      packaging-integration  ${PI_1_ID}
    Check entity doesn't exist      training  ${TRAINING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_1_NAME}

*** Test Cases ***
Apply good profile, check resources and remove on teardown
    [Documentation]  Apply good profile, validate and remove entities on end
    [Teardown]  Cleanup resources
    [Template]  Template. Apply good profile, check resources and remove on teardown
    file=correct.odahuflow.yaml
    file=correct.odahuflow.json

Apply changes on a good profile, remove on teardown
    [Documentation]  Apply changes on a good profile, validate resources, remove on teardown
    [Teardown]  Cleanup resources
    Check entity doesn't exist      connection  ${CONN_1_ID}
    Check entity doesn't exist      connection  ${CONN_2_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_1_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_2_ID}
    Check entity doesn't exist      packaging-integration  ${PI_1_ID}
    Check entity doesn't exist      training  ${TRAINING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_2_NAME}
    Apply bulk file and check counters     correct.odahuflow.yaml     6  0  0
    Check entity exists             connection  ${CONN_1_ID}
    Check entity exists             connection  ${CONN_2_ID}
    Check entity exists             toolchain-integration  ${TI_1_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_2_ID}
    Check entity exists             packaging-integration  ${PI_1_ID}
    Check entity exists             training  ${TRAINING_1_NAME}
    Check entity exists             packaging  ${PACKAGING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_2_NAME}
    Apply bulk file and check counters     correct.odahuflow.json     0  6  0
    Check entity exists             connection  ${CONN_1_ID}
    Check entity exists             connection  ${CONN_2_ID}
    Check entity exists             toolchain-integration  ${TI_1_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_2_ID}
    Check entity exists             packaging-integration  ${PI_1_ID}
    Check entity exists             training  ${TRAINING_1_NAME}
    Check entity exists             packaging  ${PACKAGING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_2_NAME}
    Apply bulk file and check counters     correct-v2.odahuflow.yaml  2  6  0
    Check entity exists             connection  ${CONN_1_ID}
    Check entity exists             connection  ${CONN_2_ID}
    Check entity exists             toolchain-integration  ${TI_1_ID}
    Check entity exists             toolchain-integration  ${TI_2_ID}
    Check entity exists             packaging-integration  ${PI_1_ID}
    Check entity exists             training  ${TRAINING_1_NAME}
    Check entity exists             packaging  ${PACKAGING_1_NAME}
    Check entity exists             packaging  ${PACKAGING_2_NAME}
    Remove bulk file and check counters    correct-v2.odahuflow.yaml  0  0  8
    Check entity doesn't exist      connection  ${CONN_1_ID}
    Check entity doesn't exist      connection  ${CONN_2_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_1_ID}
    Check entity doesn't exist      toolchain-integration  ${TI_2_ID}
    Check entity doesn't exist      packaging-integration  ${PI_1_ID}
    Check entity doesn't exist      training  ${TRAINING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_1_NAME}
    Check entity doesn't exist      packaging  ${PACKAGING_2_NAME}

Try to apply profile with incorrect order
    [Documentation]  Try to apply profile with resources in incorrect order
    [Teardown]  Cleanup resources
    Apply bulk file and check errors  incorrect-order.odahuflow.yaml  not found

Try to apply profile with syntax error
    [Documentation]  Try to apply profile with resources in incorrect order
    [Teardown]  Cleanup resources
    Apply bulk file and check parse errors  corrupted.odahuflow.yaml  not valid JSON or YAML

Try to apply profile with wrong kind
    [Documentation]  Try to apply profile with resources with wong kind
    [Teardown]  Cleanup resources
    Apply bulk file and check parse errors  wrong_kind.odahuflow.yaml  Unknown kind of object
