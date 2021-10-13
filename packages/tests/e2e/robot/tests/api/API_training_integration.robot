*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_training_integration
${RES_DIR}              ${CURDIR}/resources/training_integration
${MLFLOW}               mlflow-api-testing
${MLFLOW_NOT_EXIST}     mlflow-api-not-exist

*** Settings ***
Documentation       API of training integrations
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.TrainingIntegration
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup All Resources
Suite Teardown      Run Keywords
...                 Cleanup All Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  training-integration
Test Timeout        5 minutes

*** Keywords ***
Cleanup All Resources
    Cleanup resource  training-integration  ${MLFLOW}
    Cleanup resource  training-integration  ${MLFLOW_NOT_EXIST}

*** Test Cases ***
Get list of training integrations
    [Documentation]  check that training integrations that would be created do not exist now
    Command response list should not contain id  training integration  ${MLFLOW}

Create mlflow training integration
    Call API                    training integration post  ${RES_DIR}/valid/mlflow_create.yaml
    ${check}                    Call API  training integration get id  ${MLFLOW}
    Default Docker image should be equal  ${check}  created
    Default Entrypoint should be equal  ${check}  created

Update mlflow training integration
    sleep                       1s
    Call API                    training integration put  ${RES_DIR}/valid/mlflow_update.json
    ${check}                    Call API  training integration get id  ${MLFLOW}
    Default Docker image should be equal  ${check}  updated
    Default Entrypoint should be equal  ${check}  updated

Get updated list of training integrations
    Command response list should contain id  training integration  ${MLFLOW}

Get mlflow training integrations by id
    ${result}                   Call API  training integration get id  ${MLFLOW}
    ID should be equal          ${result}  ${MLFLOW}

Delete mlflow training integration
    ${result}                   Call API  training integration delete  ${MLFLOW}
    should be equal             ${result.get('message')}  TrainingIntegration ${MLFLOW} was deleted

Check that training integrations do not exist
    Command response list should not contain id  training integration  ${MLFLOW}

#############################
#    NEGATIVE TEST CASES    #
#############################
Try Create Training integration that already exists
    [Tags]                      negative
    [Setup]                     Cleanup resource  training-integration  ${MLFLOW}
    [Teardown]                  Cleanup resource  training-integration  ${MLFLOW}
    Call API                    training integration post  ${RES_DIR}/valid/mlflow_update.json
    ${EntityAlreadyExists}      format string  ${409 Conflict Template}  ${MLFLOW}
    Call API and get Error      ${EntityAlreadyExists}  training integration post  ${RES_DIR}/valid/mlflow_update.json

Try Update not existing Training integration
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${404NotFound}  training integration put  ${RES_DIR}/invalid/mlflow_update_not_exist.yaml

Try Update deleted Training integration
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW}
    Call API and get Error      ${404NotFound}  training integration put  ${RES_DIR}/valid/mlflow_create.yaml

Try Get id not existing Training integration
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${404NotFound}  training integration get id  ${MLFLOW_NOT_EXIST}

Try Get id deleted Training integration
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW}
    Call API and get Error      ${404NotFound}  training integration get id  ${MLFLOW}

Try Delete not existing Training integration
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW_NOT_EXIST}
    Call API and get Error      ${404NotFound}  training integration delete  ${MLFLOW_NOT_EXIST}

Try Delete deleted Training integration
    [Tags]                      negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${MLFLOW}
    Call API and get Error      ${404NotFound}  training integration delete  ${MLFLOW}
