*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_training_packaging
${RES_DIR}                          ${CURDIR}/resources/training_packaging
${TRAIN_MLFLOW_DEFAULT}             wine-mlflow-default
${PACKAGING}                        wine-api-testing
${TRAINING_ARTIFACT_NAME}           wine-mlflow-not-default-1.0.zip
${TRAIN_NOT_EXIST}                  train-api-not-exist
${PACKAGING_NOT_EXIST}              packaging-api-not-exist

*** Settings ***
Documentation       API of training, packaging, deployment, route and model
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
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
Check model trainings do not exist
    [Tags]                      training
    [Documentation]             should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Create Model Training, mlflow training, default
    [Tags]                      training
    [Documentation]             create model training with default resources and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    @{exp_result}               create list  running  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    Limits resources should be equal     ${result}  250m  ${NONE}  256Mi
    Requested resources should be equal  ${result}  125m  ${NONE}  128Mi

Update Model Training, mlflow training, default
    [Tags]                      training
    ${result}                   Call API  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded
    CreatedAt and UpdatedAt times should not be equal  ${result}
    Limits resources should be equal     ${result}  3024m  ${NONE}  4024Mi
    Requested resources should be equal  ${result}  3024m  ${NONE}  3024Mi

Get short-term Logs of training
    [Tags]                      training  log
    ${result}                   Call API  training get log  ${TRAIN_MLFLOW_DEFAULT}
    should contain              ${result}  INFO

Packaging's list doesn't contain packaging
    [Tags]                      packaging
    Command response list should not contain id  packaging  ${PACKAGING}

Create packaging
    [Tags]                      packaging
    ${artifact_name}            Pick artifact name  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    packaging post  ${RES_DIR}/valid/packaging.create.yaml  ${artifact_name}
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  entity=${PACKAGING}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded

Update packaging
    [Tags]                      packaging
    ${result}                   Call API  packaging put  ${RES_DIR}/valid/packaging.update.yaml  ${TRAINING_ARTIFACT_NAME}
    ${result_pack}              Call API  packaging get id  ${PACKAGING}
    should be equal             ${result_pack.spec.integration_name}  docker-rest
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  entity=${PACKAGING}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded
    CreatedAt and UpdatedAt times should not be equal  ${result}

Get packaging by id
    [Tags]                      packaging
    ${result}                   Call API  packaging get id  ${PACKAGING}
    ID should be equal          ${result}  ${PACKAGING}

Get logs of packaging
    [Tags]                      packaging  log
    ${result}                   Call API  packaging get log  ${PACKAGING}
    should contain              ${result}  INFO

Delete Model Trainings and Check that Model Training do not exist
    [Tags]                      training
    [Documentation]             delete model trainings
    [Teardown]                  Cleanup resource  training  ${TRAIN_MLFLOW_DEFAULT}
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_DEFAULT}
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Delete Model Packaging and Check that Model Packaging does not exist
    [Tags]                      packaging
    [Teardown]                  Cleanup resource  packaging  ${PACKAGING}
    Command response list should contain id  packaging  ${PACKAGING}
    Call API                    packaging delete  ${PACKAGING}
    Command response list should not contain id  packaging  ${PACKAGING}

#############################
#    NEGATIVE TEST CASES    #
#############################

#  TRAINING
#############
Try Create Training that already exists
    [Tags]                      negative
    [Teardown]                  Cleanup resource  training  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    ${EntityAlreadyExists}      Format EntityAlreadyExists  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${EntityAlreadyExists}  training post  ${RES_DIR}/valid/training.mlflow.default.update.yaml

Try Update not existing and deleted Training
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${TRAIN_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  training put  ${RES_DIR}/invalid/training.update.not_exist.json
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${WrongHttpStatusCode}  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml

Try Get id not existing and deleted Training
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${TRAIN_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  training get id  ${TRAIN_NOT_EXIST}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${WrongHttpStatusCode}  training get id  ${TRAIN_MLFLOW_DEFAULT}

Try Delete not existing and deleted Training
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${TRAIN_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  training delete  ${TRAIN_NOT_EXIST}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${TRAIN_MLFLOW_DEFAULT}
    Call API and get Error      ${WrongHttpStatusCode}  training delete  ${TRAIN_MLFLOW_DEFAULT}

#  PACKAGING
#############
Try Create Packaging that already exists
    [Tags]                      negative
    [Teardown]                  Cleanup resource  packaging  ${PACKAGING}
    Call API                    packaging post  ${RES_DIR}/valid/packaging.update.yaml  ${TRAINING_ARTIFACT_NAME}
    ${EntityAlreadyExists}      Format EntityAlreadyExists  ${PACKAGING}
    Call API and get Error      ${EntityAlreadyExists}  packaging post  ${RES_DIR}/valid/packaging.create.yaml  ${TRAINING_ARTIFACT_NAME}

Try Update not existing and deleted Packaging
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${PACKAGING_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  packaging put  ${RES_DIR}/invalid/packaging.update.not_exist.yaml  ${TRAINING_ARTIFACT_NAME}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${PACKAGING}
    Call API and get Error      ${WrongHttpStatusCode}  packaging put  ${RES_DIR}/valid/packaging.create.yaml  ${TRAINING_ARTIFACT_NAME}

Try Get id not existing and deleted Packaging
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${PACKAGING_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  packaging get id  ${PACKAGING_NOT_EXIST}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${PACKAGING}
    Call API and get Error      ${WrongHttpStatusCode}  packaging get id  ${PACKAGING}

Try Delete not existing and deleted Packaging
    [Tags]                      negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${PACKAGING_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  packaging delete  ${PACKAGING_NOT_EXIST}
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${PACKAGING}
    Call API and get Error      ${WrongHttpStatusCode}  packaging delete  ${PACKAGING}
