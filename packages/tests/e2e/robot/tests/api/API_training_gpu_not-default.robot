*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_training_gpu_not-default
${RES_DIR}                          ${CURDIR}/resources/training_packaging
${TRAIN_MLFLOW_NOT_DEFAULT}         wine-mlflow-not-default
${TRAIN_MLFLOW-GPU_NOT_DEFAULT}     reuters-classifier-mlflow-gpu-not-default

*** Settings ***
Documentation       API of training, packaging, deployment, route and model
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup Resources
Suite Teardown      Run Keywords
...                 Cleanup Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk
Test Timeout        60 minutes

*** Keywords ***
Cleanup Resources
    [Documentation]  Deletes of created resources
    StrictShell  odahuflowctl --verbose train delete --id ${TRAIN_MLFLOW_NOT_DEFAULT} --ignore-not-found
    StrictShell  odahuflowctl --verbose train delete --id ${TRAIN_MLFLOW-GPU_NOT_DEFAULT} --ignore-not-found

*** Test Cases ***
Check model trainings do not exist
    [Tags]                      training
    [Documentation]             should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Create Model Training, mlflow toolchain, not default
    [Tags]                      training
    [Documentation]             create model training and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.not_default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_NOT_DEFAULT}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded

Create and Delete Model Training, mlflow-gpu toolchain, not default
    [Tags]                      training
    [Documentation]             create model training with mlflow-gpu toolchain and not default values
    ...                         cluster with GPU node pools enabled
    Pass Execution If           not ${IS_GPU_ENABLED}  GPU node pools is not enabled on the cluster

    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow-gpu.not_default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW-GPU_NOT_DEFAULT}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded

    Command response list should contain id  training  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
    Command response list should not contain id  training  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Get Model Training by id
    [Tags]                      training
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ID should be equal          ${result}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Delete Model Trainings and Check that Model Training do not exist
    [Tags]                      training
    [Documentation]             delete model trainings
    Command response list should contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Command response list should not contain id  training  ${TRAIN_MLFLOW_NOT_DEFAULT}
