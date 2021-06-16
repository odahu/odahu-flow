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
${E2E_PACKAGING}        simple-model
${DEPLOYMENT_1_NAME}    bulk-test-md-1

*** Settings ***
Documentation       OdahuFlow's API operational check for bulk operations (with multiple resources)
Test Timeout        20 minutes
Resource            ../../resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Default Tags        cli  bulk
Suite Setup         Test Setup
Suite Teardown      Run keywords
...                 Cleanup resources
...                 AND  Remove File  ${LOCAL_CONFIG}

*** Keywords ***
Test Setup
    Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
    Login to the api and edge
    Cleanup resources

    # update image for model deployment
    ${res}=  StrictShell  odahuflowctl pack get --id ${E2E_PACKAGING} -o 'jsonpath=$[0].spec.image'
    StrictShell  yq w --inplace -d'6' ${RES_DIR}/correct.odahuflow.yaml 'spec.image' ${res.stdout}
    StrictShell  yq w --inplace -d'8' ${RES_DIR}/correct-v2.odahuflow.yaml 'spec.image' ${res.stdout}
    StrictShell  yq w --inplace -jP ${RES_DIR}/correct.odahuflow.json '[6].spec.image' ${res.stdout}

Remove ${entity_type} with id - "${name}"
    StrictShell  odahuflowctl --verbose ${entity_type} delete --id ${name} --ignore-not-found

Cleanup resources
    Remove conn with id - "${CONN_1_ID}"
    Remove conn with id - "${CONN_1_ID}"
    Remove conn with id - "${CONN_2_ID}"
    Remove ti with id - "${TI_1_ID}"
    Remove ti with id - "${TI_2_ID}"
    Remove pi with id - "${PI_1_ID}"
    Remove train with id - "${TRAINING_1_NAME}"
    Remove pack with id - "${PACKAGING_1_NAME}"
    Remove pack with id - "${PACKAGING_2_NAME}"
    Remove dep with id - "${DEPLOYMENT_1_NAME}"

Check ${entity_type} exists - "${name}"
    ${res}=  StrictShell  odahuflowctl --verbose ${entity_type} get --id ${name}
        Should contain     ${res.stderr}  ${name}

Check ${entity_type} doesn't exist - "${name}"
    ${res}=  Wait Until Keyword Succeeds  5m  2s  FailedShell  odahuflowctl --verbose ${entity_type} get --id ${name}
        Should contain     ${res.stderr}  "${name}" is not found

Apply bulk file and check counters
    [Arguments]  ${file}  ${created}  ${changed}  ${deleted}
    ${res}=  StrictShell  odahuflowctl --verbose bulk apply ${RES_DIR}/${file}
             Run Keyword If   ${created} > 0  Should contain   ${res.stdout}  created resources: ${created}
             Run Keyword If   ${changed} > 0  Should contain   ${res.stdout}  changed resources: ${changed}
             Run Keyword If   ${deleted} > 0  Should contain   ${res.stdout}  deleted resources: ${deleted}

Apply bulk file and check errors
    [Arguments]  ${file}  ${expected_error_message}
    ${res}=  FailedShell  odahuflowctl --verbose bulk apply ${RES_DIR}/${file}
             Should contain       ${res.stdout}  ${expected_error_message}

Apply bulk file and check parse errors
    [Arguments]  ${file}  ${expected_error_message}
    ${res}=  FailedShell  odahuflowctl --verbose bulk apply ${RES_DIR}/${file}
             Should contain       ${res.stderr}  ${expected_error_message}

Remove bulk file and check counters
    [Arguments]  ${file}  ${created}  ${changed}  ${deleted}
    ${res}=  StrictShell  odahuflowctl --verbose bulk delete ${RES_DIR}/${file}
             Run Keyword If   ${created} > 0  Should contain   ${res.stdout}  created resources: ${created}
             Run Keyword If   ${changed} > 0  Should contain   ${res.stdout}  changed resources: ${changed}
             Run Keyword If   ${deleted} > 0  Should contain   ${res.stdout}  removed resources: ${deleted}

Template. Apply good profile, check resources and remove on teardown
    [Arguments]  ${file}
    [Teardown]  Cleanup resources
    Check connection doesn't exist - "${CONN_1_ID}"
    Check connection doesn't exist - "${CONN_2_ID}"
    Check toolchain-integration doesn't exist - "${TI_1_ID}"
    Check packaging-integration doesn't exist - "${PI_1_ID}"
    Check training doesn't exist - "${TRAINING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_1_NAME}"
    Check deployment doesn't exist - "${DEPLOYMENT_1_NAME}"
    Apply bulk file and check counters     ${file}  7  0  0
    Check connection exists - "${CONN_1_ID}"
    Check connection exists - "${CONN_2_ID}"
    Check toolchain-integration exists - "${TI_1_ID}"
    Check packaging-integration exists - "${PI_1_ID}"
    Check training exists - "${TRAINING_1_NAME}"
    Check packaging exists - "${PACKAGING_1_NAME}"
    Check deployment exists - "${DEPLOYMENT_1_NAME}"
    Remove bulk file and check counters    ${file}  0  0  7
    Check connection doesn't exist - "${CONN_1_ID}"
    Check connection doesn't exist - "${CONN_2_ID}"
    Check toolchain-integration doesn't exist - "${TI_1_ID}"
    Check packaging-integration doesn't exist - "${PI_1_ID}"
    Check training doesn't exist - "${TRAINING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_1_NAME}"
    Check deployment doesn't exist - "${DEPLOYMENT_1_NAME}"

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
    Check connection doesn't exist - "${CONN_1_ID}"
    Check connection doesn't exist - "${CONN_2_ID}"
    Check toolchain-integration doesn't exist - "${TI_1_ID}"
    Check toolchain-integration doesn't exist - "${TI_2_ID}"
    Check packaging-integration doesn't exist - "${PI_1_ID}"
    Check training doesn't exist - "${TRAINING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_2_NAME}"
    Check deployment doesn't exist - "${DEPLOYMENT_1_NAME}"
    Apply bulk file and check counters     correct.odahuflow.yaml     7  0  0
    Check connection exists - "${CONN_1_ID}"
    Check connection exists - "${CONN_2_ID}"
    Check toolchain-integration exists - "${TI_1_ID}"
    Check toolchain-integration doesn't exist - "${TI_2_ID}"
    Check packaging-integration exists - "${PI_1_ID}"
    Check training exists - "${TRAINING_1_NAME}"
    Check packaging exists - "${PACKAGING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_2_NAME}"
    Check deployment exists - "${DEPLOYMENT_1_NAME}"
    Apply bulk file and check counters     correct.odahuflow.json     0  7  0
    Check connection exists - "${CONN_1_ID}"
    Check connection exists - "${CONN_2_ID}"
    Check toolchain-integration exists - "${TI_1_ID}"
    Check toolchain-integration doesn't exist - "${TI_2_ID}"
    Check packaging-integration exists - "${PI_1_ID}"
    Check training exists - "${TRAINING_1_NAME}"
    Check packaging exists - "${PACKAGING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_2_NAME}"
    Check deployment exists - "${DEPLOYMENT_1_NAME}"
    Apply bulk file and check counters     correct-v2.odahuflow.yaml  2  7  0
    Check connection exists - "${CONN_1_ID}"
    Check connection exists - "${CONN_2_ID}"
    Check toolchain-integration exists - "${TI_1_ID}"
    Check toolchain-integration exists - "${TI_2_ID}"
    Check packaging-integration exists - "${PI_1_ID}"
    Check training exists - "${TRAINING_1_NAME}"
    Check packaging exists - "${PACKAGING_1_NAME}"
    Check packaging exists - "${PACKAGING_2_NAME}"
    Check deployment exists - "${DEPLOYMENT_1_NAME}"
    Remove bulk file and check counters    correct-v2.odahuflow.yaml  0  0  9
    Check connection doesn't exist - "${CONN_1_ID}"
    Check connection doesn't exist - "${CONN_2_ID}"
    Check toolchain-integration doesn't exist - "${TI_1_ID}"
    Check toolchain-integration doesn't exist - "${TI_2_ID}"
    Check packaging-integration doesn't exist - "${PI_1_ID}"
    Check training doesn't exist - "${TRAINING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_1_NAME}"
    Check packaging doesn't exist - "${PACKAGING_2_NAME}"
    Check deployment doesn't exist - "${DEPLOYMENT_1_NAME}"

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
