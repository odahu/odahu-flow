*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_training_packaging
${RES_DIR}                          ${CURDIR}/resources/training_packaging
${TRAIN_MLFLOW_DEFAULT}             wine-mlflow-default
${PACKAGING}                        wine-api-testing
${TRAINING_ARTIFACT_NAME}           wine-mlflow-not-default-1.0.zip

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
...                 Cleanup Resources
Suite Teardown      Run Keywords
...                 Cleanup Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk
Test Timeout        60 minutes

*** Keywords ***
Cleanup Resources
    [Documentation]  Deletes of created resources
    StrictShell  odahuflowctl --verbose train delete --id ${TRAIN_MLFLOW_DEFAULT} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${PACKAGING} --ignore-not-found

*** Test Cases ***
Check model trainings do not exist
    [Tags]                      training
    [Documentation]             should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Create Model Training, mlflow toolchain, default
    [Tags]                      training
    [Documentation]             create model training with default resources and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    @{exp_result}               create list  running  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    Get Logs                    training  ${TRAIN_MLFLOW_DEFAULT}
    Limits resources should be equal     ${result}  250m  ${NONE}  256Mi
    Requested resources should be equal  ${result}  125m  ${NONE}  128Mi

Update Model Training, mlflow toolchain, default
    [Tags]                      training
    ${result}                   Call API  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    Get Logs                    training  ${TRAIN_MLFLOW_DEFAULT}
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
    Get Logs                    packaging  ${PACKAGING}
    Status State Should Be      ${result}  succeeded

Update packaging
    [Tags]                      packaging
    ${result}                   Call API  packaging put  ${RES_DIR}/valid/packaging.update.yaml  ${TRAINING_ARTIFACT_NAME}
    ${result_pack}              Call API  packaging get id  ${PACKAGING}
    should be equal             ${result_pack.spec.integration_name}  docker-rest
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  entity=${PACKAGING}  exp_result=@{exp_result}
    Get Logs                    packaging  ${PACKAGING}
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
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_DEFAULT}
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}

Delete Model Packaging and Check that Model Packaging does not exist
    [Tags]                      packaging
    Command response list should contain id  packaging  ${PACKAGING}
    Call API                    packaging delete  ${PACKAGING}
    Command response list should not contain id  packaging  ${PACKAGING}
