*** Variables ***
${RES_DIR}                          ${CURDIR}/resources/training_packaging_deploy
${TRAIN_MLFLOW_DEFAULT}             wine-mlflow-default
${TRAIN_MLFLOW_NOT_DEFAULT}         wine-mlflow-not-default
${TRAIN_MLFLOW-GPU_NOT_DEFAULT}     reuters-classifier-mlflow-gpu-not-default

*** Settings ***
Documentation       API of training, packaging and deployment
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Library             odahuflow.robot.libraries.sdk_wrapper.ModelDeployment
Suite Setup         Run Keywords
...                 Login to the api and edge
Force Tags          api  training  packaging  deployment

*** Test Cases ***
Get list of model training
    [Documentation]  should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Create Model Training, mlflow toolchain, default
    [Documentation]  create model training and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    @{exp_result}               set variable  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW_DEFAULT}
    should be equal             ${result.status.state}  succeeded

Create Model Training, mlflow toolchain, not default
    [Documentation]  create model training and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.not_default.yaml
    @{exp_result}               set variable  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW_NOT_DEFAULT}
    should be equal             ${result.status.state}  succeeded

Create Model Training, mlflow-gpu toolchain, not default
    [Documentation]  create model training with mlflow-gpu toolchain and not default values
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow-gpu.not_default.yaml
    @{exp_result}               set variable  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
    should be equal             ${result.status.state}  succeeded

Get Model Training by id
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ${result_id}                Log id  ${result}
    should be equal             ${result_id}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Update Model Training, mlflow toolchain, not default
    ${result}                   Call API  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml
    log                         ${result.spec.model.version}
    @{exp_result}               set variable  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW_DEFAULT}
    should be equal             ${result.status.state}  succeeded
    ${result_status}            Log Status  ${result}
    should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

# Redo cannot stout this
Get Logs of training
    ${result}                   Call API  training get log  ${TRAIN_MLFLOW_NOT_DEFAULT}

Get updated list of model training
    [Documentation]  check that new training are in the list
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Delete Model Trainings
    [Documentation]  delete model trainings
    Call API                    training delete  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Check that Model Trainging do not exist
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
