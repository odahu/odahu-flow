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
    [Documentation]             check that the training to be created does not exist now
    Command response list should not contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}

Create Model Training, mlflow toolchain, not default
    Call API  training post  ${RES_DIR}/valid/training.mlflow.not_default.yaml
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_NOT_DEFAULT}
    Get Logs                    training  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Status State Should Be      ${result}  succeeded

Get Model Training by id
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ID should be equal          ${result}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Delete Model Trainings and Check that Model Training do not exist
    Command response list should contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ${result}                   Call API  training delete  ${TRAIN_MLFLOW_NOT_DEFAULT}
    should be equal             ${result}  Model training ${TRAIN_MLFLOW_NOT_DEFAULT} was deleted
    Command response list should not contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}
