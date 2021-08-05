*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_training_packaging
${RES_DIR}                          ${CURDIR}/resources/training_packaging
${TRAIN_MLFLOW_DEFAULT}             wine-mlflow-default
${PACKAGING}                        wine-api-testing
${TRAINING_ARTIFACT_NAME}           wine-mlflow-default-updated-1.1.zip
${TRAIN_NOT_EXIST}                  train-api-not-exist
${PACKAGING_NOT_EXIST}              packaging-api-not-exist

*** Settings ***
Documentation       API of training and packaging
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup All Resources
Suite Teardown      Run Keywords
...                 Cleanup All Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk
Test Timeout        60 minutes

*** Keywords ***
Cleanup All Resources
    Cleanup resource  training  ${TRAIN_MLFLOW_DEFAULT}
    Cleanup resource  packaging  ${PACKAGING}
    Cleanup resource  training  ${TRAIN_NOT_EXIST}
    Cleanup resource  packaging  ${PACKAGING_NOT_EXIST}

*** Test Cases ***
Training's list doesn't contain not created training
    [Tags]                      training
    [Documentation]             check that the training to be created does not exist now
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Create Model Training, mlflow training, default
    [Tags]                      training
    [Documentation]             create model training with default resources and check that one exists
    Call API  training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    @{exp_result}               create list  running  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    wait until keyword succeeds  2m  2s    Get Logs  training  ${TRAIN_MLFLOW_DEFAULT}
    Limits resources should be equal     ${result}  250m  ${NONE}  256Mi
    Requested resources should be equal  ${result}  125m  ${NONE}  128Mi

Update Model Training, mlflow training, default
    [Tags]                      training
    Call API  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    wait until keyword succeeds  2m  2s    Get Logs  training  ${TRAIN_MLFLOW_DEFAULT}
    Status State Should Be      ${result}  succeeded
    CreatedAt and UpdatedAt times should not be equal  ${result}
    Limits resources should be equal     ${result}  3024m  ${NONE}  4024Mi
    Requested resources should be equal  ${result}  3024m  ${NONE}  3024Mi

Get short-term Logs of training
    [Tags]                      training  log
    ${result}                   Call API  training get log  ${TRAIN_MLFLOW_DEFAULT}
    should contain              ${result}  INFO

Get training by id
    [Tags]                      training
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_DEFAULT}
    ID should be equal          ${result}  ${TRAIN_MLFLOW_DEFAULT}

Get updated list of trainings
    [Tags]                      training
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Packaging's list doesn't contain not created packaging
    [Tags]                      packaging
    [Documentation]             check that the packaging to be created does not exist now
    Command response list should not contain id  packaging  ${PACKAGING}

Create packaging
    [Tags]                      packaging
    ${artifact_name}            Pick artifact name  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    packaging post  ${RES_DIR}/valid/packaging.create.yaml  ${artifact_name}
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  entity=${PACKAGING}  exp_result=@{exp_result}
    Get Logs                    packaging  ${PACKAGING}
    Status State Should Be      ${result}  succeeded

Update packaging
    [Tags]                      packaging
    Call API  packaging put     ${RES_DIR}/valid/packaging.update.yaml  ${TRAINING_ARTIFACT_NAME}
    ${result_pack}              Call API  packaging get id  ${PACKAGING}
    should be equal             ${result_pack.spec.integration_name}  docker-rest
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  entity=${PACKAGING}  exp_result=@{exp_result}
    Get Logs                    packaging  ${PACKAGING}
    Status State Should Be      ${result}  succeeded
    CreatedAt and UpdatedAt times should not be equal  ${result}

Get logs of packaging
    [Tags]                      packaging  log
    ${result}                   Call API  packaging get log  ${PACKAGING}
    should contain              ${result}  INFO

Get packaging by id
    [Tags]                      packaging
    ${result}                   Call API  packaging get id  ${PACKAGING}
    ID should be equal          ${result}  ${PACKAGING}

Get updated list of packagings
    [Tags]                      packaging
    Command response list should contain id  packaging  ${PACKAGING}

Delete Model Trainings and Check that Model Training do not exist
    [Tags]                      training
    [Documentation]             delete model trainings
    [Teardown]                  Cleanup resource  training  ${TRAIN_MLFLOW_DEFAULT}
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ${result}                   Call API  training delete  ${TRAIN_MLFLOW_DEFAULT}
    should be equal             ${result.get('message')}  Model training ${TRAIN_MLFLOW_DEFAULT} was deleted
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Delete Model Packaging and Check that Model Packaging does not exist
    [Tags]                      packaging
    [Teardown]                  Cleanup resource  packaging  ${PACKAGING}
    Command response list should contain id  packaging  ${PACKAGING}
    ${result}                   Call API  packaging delete  ${PACKAGING}
    should be equal             ${result.get('message')}  Model packaging ${PACKAGING} was deleted
    Command response list should not contain id  packaging  ${PACKAGING}

#############################
#    NEGATIVE TEST CASES    #
#############################

#  TRAINING
#############
Try Create Training that already exists
    [Tags]                      training  negative
    [Setup]                     Cleanup resource  training  ${TRAIN_MLFLOW_DEFAULT}
    [Teardown]                  Cleanup resource  training  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    ${EntityAlreadyExists}      format string  ${409 Conflict Template}  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${EntityAlreadyExists}  training post  ${RES_DIR}/valid/training.mlflow.default.update.yaml

Try Update not existing Training
    [Tags]                      training  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${TRAIN_NOT_EXIST}
    Call API and get Error      ${404NotFound}  training put  ${RES_DIR}/invalid/training.update.not_exist.json

Try Update deleted Training
    [Tags]                      training  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${404NotFound}  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml

Try Get id not existing Training
    [Tags]                      training  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${TRAIN_NOT_EXIST}
    Call API and get Error      ${404NotFound}  training get id  ${TRAIN_NOT_EXIST}

Try Get id deleted Training
    [Tags]                      training  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${404NotFound}  training get id  ${TRAIN_MLFLOW_DEFAULT}

Try Delete not existing Training
    [Tags]                      training  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${TRAIN_NOT_EXIST}
    Call API and get Error      ${404NotFound}  training delete  ${TRAIN_NOT_EXIST}

Try Delete deleted Training
    [Tags]                      training  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${404NotFound}  training delete  ${TRAIN_MLFLOW_DEFAULT}

#  PACKAGING
#############
Try Create Packaging that already exists
    [Tags]                      packaging  negative
    [Setup]                     Cleanup resource  packaging  ${PACKAGING}
    [Teardown]                  Cleanup resource  packaging  ${PACKAGING}
    Call API                    packaging post  ${RES_DIR}/valid/packaging.update.yaml  ${TRAINING_ARTIFACT_NAME}
    ${EntityAlreadyExists}      format string  ${409 Conflict Template}  ${PACKAGING}
    Call API and get Error      ${EntityAlreadyExists}  packaging post  ${RES_DIR}/valid/packaging.create.yaml  ${TRAINING_ARTIFACT_NAME}

Try Update not existing Packaging
    [Tags]                      packaging  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${PACKAGING_NOT_EXIST}
    Call API and get Error      ${404NotFound}  packaging put  ${RES_DIR}/invalid/packaging.update.not_exist.yaml  ${TRAINING_ARTIFACT_NAME}

Try Update deleted Packaging
    [Tags]                      packaging  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${PACKAGING}
    Call API and get Error      ${404NotFound}  packaging put  ${RES_DIR}/valid/packaging.create.yaml  ${TRAINING_ARTIFACT_NAME}

Try Get id not existing Packaging
    [Tags]                      packaging  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${PACKAGING_NOT_EXIST}
    Call API and get Error      ${404NotFound}  packaging get id  ${PACKAGING_NOT_EXIST}

Try Get id deleted Packaging
    [Tags]                      packaging  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${PACKAGING}
    Call API and get Error      ${404NotFound}  packaging get id  ${PACKAGING}

Try Delete not existing Packaging
    [Tags]                      packaging  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${PACKAGING_NOT_EXIST}
    Call API and get Error      ${404NotFound}  packaging delete  ${PACKAGING_NOT_EXIST}

Try Delete deleted Packaging
    [Tags]                      packaging  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${PACKAGING}
    Call API and get Error      ${404NotFound}  packaging delete  ${PACKAGING}
