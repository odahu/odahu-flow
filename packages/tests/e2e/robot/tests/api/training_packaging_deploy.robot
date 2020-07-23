*** Variables ***
${RES_DIR}                          ${CURDIR}/resources/training_packaging_deploy
${TRAIN_MLFLOW_DEFAULT}             wine-mlflow-default
${TRAIN_MLFLOW_NOT_DEFAULT}         wine-mlflow-not-default
${TRAIN_MLFLOW-GPU_NOT_DEFAULT}     reuters-classifier-mlflow-gpu-not-default

${PACKAGING}                        wine-api-testing
${TRAINING_ARTIFACT_NAME}           wine-mlflow-not-default-1.0.zip

${DEPLOYMENT}                       wine-api-testing
${MODEL}                            wine-api-testing

${REQUEST}                          SEPARATOR=
...                                 { "columns": [ "fixed acidity", "volatile acidity", "citric acid",
...                                 "residual sugar", "chlorides", "free sulfur dioxide", "total sulfur dioxide", "density",
...                                 "pH", "sulphates", "alcohol" ], "data": [ [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ] ] }

*** Settings ***
Documentation       API of training, packaging and deployment
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Library             odahuflow.robot.libraries.sdk_wrapper.ModelDeployment
Library             odahuflow.robot.libraries.sdk_wrapper.ModelRoute
Library             odahuflow.robot.libraries.sdk_wrapper.Model
Suite Setup         Run Keywords
...                 Login to the api and edge
Force Tags          api  sdk  testing

*** Test Cases ***
Get list of model training
    [Tags]                      training
    [Documentation]             should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Create Model Training, mlflow toolchain, default
    [Tags]                      training
    [Documentation]             create model training and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW_DEFAULT}
    should be equal             ${result.status.state}  succeeded

Create Model Training, mlflow toolchain, not default
    [Tags]                      training
    [Documentation]             create model training and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.not_default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW_NOT_DEFAULT}
    should be equal             ${result.status.state}  succeeded

Create Model Training, mlflow-gpu toolchain, not default
    [Tags]                      training
    [Documentation]             create model training with mlflow-gpu toolchain and not default values
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow-gpu.not_default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
    should be equal             ${result.status.state}  succeeded

Get Model Training by id
    [Tags]                      training
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ${result_id}                Log id  ${result}
    should be equal             ${result_id}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Update Model Training, mlflow toolchain, not default
    [Tags]                      training
    ${result}                   Call API  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml
    log                         ${result.spec.model.version}
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  result=${result}  exp_result=@{exp_result}
    Call API                    training get log  ${TRAIN_MLFLOW_DEFAULT}
    should be equal             ${result.status.state}  succeeded
    ${result_status}            Log Status  ${result}
    should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Get Logs of training
    [Tags]                      training  log
    ${result}                   Call API  training get log  ${TRAIN_MLFLOW_DEFAULT}
    should contain              ${result}  INFO

Get list of packagings
    [Tags]                      packaging
    Command response list should not contain id  packaging  ${PACKAGING}

Create packaging
    [Tags]                      packaging
    ${result_train}             Call API  training get id  ${TRAIN_MLFLOW_DEFAULT}
    ${artifactName}             Pick artifact name  ${result_train}
    ${result_pack}              Call API  packaging post  ${RES_DIR}/valid/packaging.yaml  ${artifactName}
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  result=${result_pack}  exp_result=@{exp_result}
    should be equal             ${result.status.state}  succeeded

Update packaging
    [Tags]                      packaging
    ${result_pack}              Call API  packaging put  ${RES_DIR}/valid/packaging.update.yaml  ${TRAINING_ARTIFACT_NAME}
    ${check_changes}            Call API  packaging get id  ${PACKAGING}
    should be equal             ${check_changes.spec.integration_name}  docker-rest
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  packaging  result=${result_pack}  exp_result=@{exp_result}
    should be equal             ${result.status.state}  succeeded
    ${result_status}            Log Status  ${result}
    should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Get packaging by id
    [Tags]                      packaging
    ${result}                   Call API  packaging get id  ${PACKAGING}
    should be equal             ${result.id}  ${PACKAGING}

Get logs of packaging
    [Tags]                      packaging  log
    ${result}                   Call API  packaging get log  ${PACKAGING}
    should contain              ${result}  INFO

Get list of deployments
    [Tags]                      deployment
    Command response list should not contain id  deployment  ${DEPLOYMENT}

Get list of routes
    [Tags]                      route
    Command response list should not contain id  route  ${DEPLOYMENT}

Create deployment
    [Tags]                      deployment
    ${result_pack}              Call API  packaging get id  ${PACKAGING}
    ${image}                    Pick packaging image  ${result_pack}
    ${result_deploy}            Call API  deployment post  ${RES_DIR}/valid/deployment.create.yaml  ${image}
    ${exp_result}               create List   Ready
    ${result}                   Wait until command finishes and returns result  deployment  result=${result_deploy}  exp_result=${exp_result}
    should be equal             ${result.status.state}  Ready

Update deployment
    [Tags]                      deployment
    ${result_pack}              Call API  packaging get id  ${PACKAGING}
    ${image}                    Pick packaging image  ${result_pack}
    ${result_deploy}            Call API  deployment put  ${RES_DIR}/valid/deployment.update.json  ${image}
    ${check_changes}            Call API  deployment get id  ${DEPLOYMENT}
    should be equal             ${check_changes.spec.role_name}  test_updated
    ${exp_result}               create List   Ready
    ${result}                   Wait until command finishes and returns result  deployment  result=${result_deploy}  exp_result=${exp_result}
    should be equal             ${result.status.state}  Ready
    ${result_status}            Log Status  ${result}
    should not be equal         ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Get deployment by id
    [Tags]                      deployment
    ${result}                   Call API  deployment get id  ${DEPLOYMENT}
    should be equal             ${result.id}  ${DEPLOYMENT}

Get model route
    [Tags]                      route
    ${result}                   Call API  route get id  ${DEPLOYMENT}


Get info about model
    [Tags]                      model
    ${result}                   Call API  model get  https://gke04.dev.odahu.org/model/wine-api-testing

Invoke model
    [Tags]                      model
    ${url}                      set variable  https://gke04.dev.odahu.org/model/wine-api-testing
    log                         ${REQUEST}
    ${REQUEST}                  convert to string  ${REQUEST}
    ${request_type}             evaluate  type(${REQUEST})
    ${result}                   Call API  model post  url=${url}  json_input=REQUEST

Check updated list of Model Trainings
    [Tags]                      training
    [Documentation]             check that new training are in the list
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Delete Model Trainings
    [Tags]                      training
    [Documentation]             delete model trainings
    Call API                    training delete  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_NOT_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Check that Model Trainging do not exist
    [Tags]                      training
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Check updated list of Model Packagings
    [Tags]                      packaging
    Command response list should contain id  packaging  ${PACKAGING}

Delete Model Packaging
    [Tags]                      packaging
    Call API                    packaging delete  ${PACKAGING}

Check that Model Packaging does not exist
    [Tags]                      packaging
    Command response list should not contain id  packaging  ${PACKAGING}

Check updated list of Model Deployments
    [Tags]                      deployment
    Command response list should not contain id  deployment  ${DEPLOYMENT}

Delete Model Deployment
    [Tags]                      deployment
    Call API                    deployment delete  ${PACKAGING}

Check that Model Deployment does not exist
    [Tags]                      deployment
    Command response list should not contain id  deployment  ${DEPLOYMENT}

Check updated list of Models
    [Tags]                      model
    Command response list should not contain id  model  ${MODEL}
