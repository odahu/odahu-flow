*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_training_not-default
${RES_DIR}                          ${CURDIR}/resources/training_packaging
${TRAIN_MLFLOW_NOT_DEFAULT}         wine-mlflow-not-default

*** Settings ***
Documentation       API for not default training
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup All Resources
Suite Teardown      Run Keywords
...                 Cleanup All Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  training
Test Timeout        60 minutes

*** Keywords ***
Cleanup All Resources
    Cleanup resource  training  ${TRAIN_MLFLOW_NOT_DEFAULT}

*** Test Cases ***
Check model trainings do not exist
    [Documentation]             should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}

Create Model Training, mlflow toolchain, not default
    Call API  training post  ${RES_DIR}/valid/training.mlflow.not_default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_NOT_DEFAULT}  exp_result=@{exp_result}
    Get Logs                    training  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Status State Should Be      ${result}  succeeded

Get Model Training by id
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ID should be equal          ${result}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Delete Model Trainings and Check that Model Training do not exist
    Command response list should contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Command response list should not contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}
