*** Variables ***
${RES_DIR}                          ${CURDIR}/resources/training_packaging_deploy
${TRAIN_MLFLOW_DEFAULT}             wine-mlflow-default
${TRAIN_MLFLOW_NOT_DEFAULT}         wine-mlflow-not-default
${TRAIN_MLFLOW-GPU_NOT_DEFAULT}     reuters-classifier-mlflow-gpu-not-default
${TRAINING_ARTIFACT_NAME}           wine-mlflow-not-default-1.0.zip
${PACKAGING}                        wine-api-testing
${DEPLOYMENT}                       wine-api-testing
${MODEL}                            ${DEPLOYMENT}
${MODEL_URL}                        ${EDGE_URL}/model/${MODEL}
${REQUEST}                          SEPARATOR=
...                                 { "columns": [ "fixed acidity", "volatile acidity", "citric acid",
...                                 "residual sugar", "chlorides", "free sulfur dioxide", "total sulfur dioxide", "density",
...                                 "pH", "sulphates", "alcohol" ], "data": [ [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ] ] }
${REQUEST_RESPONSE}                 {'prediction': [4.675934265196686], 'columns': ['quality']}

${WrongHttpStatusCode}              WrongHttpStatusCode: Got error from server: modeldeployments.odahuflow.odahu.org
...                                 "wine-api-testing" not found (status: 404)

*** Settings ***
Documentation       API of training, packaging, deployment, route and model
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
Force Tags          api  sdk
Test Timeout        60 minutes

*** Test Cases ***
Check model trainings do not exist
    [Tags]                      training
    [Documentation]             should not contain training that has not been run
    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}
    ...                                          ${TRAIN_MLFLOW_NOT_DEFAULT}  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Create Model Training, mlflow toolchain, default
    [Tags]                      training
    [Documentation]             create model training and check that one exists
    ${result}                   Call API  training post  ${RES_DIR}/valid/training.mlflow.default.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded

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
    ${result}                   Wait until command finishes and returns result  training  40  30s  entity=${TRAIN_MLFLOW-GPU_NOT_DEFAULT}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded

    Command response list should not contain id  training  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
                                Call API  training delete  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}
    Command response list should not contain id  training  ${TRAIN_MLFLOW-GPU_NOT_DEFAULT}

Get Model Training by id
    [Tags]                      training
    ${result}                   Call API  training get id  ${TRAIN_MLFLOW_NOT_DEFAULT}
    ID should be equal          ${result}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Update Model Training, mlflow toolchain, default
    [Tags]                      training
    ${result}                   Call API  training put  ${RES_DIR}/valid/training.mlflow.default.update.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  training  entity=${TRAIN_MLFLOW_DEFAULT}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded
    CreatedAt and UpdatedAt times should not be equal  ${result}

Get Logs of training
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

Check deployment doesn't exist
    [Tags]                      deployment
    Command response list should not contain id  deployment  ${DEPLOYMENT}

Check route doesn't exist before deployment
    [Tags]                      route
    Command response list should not contain id  route  ${MODEL}

Create deployment
    [Tags]                      deployment
    ${image}                    Pick packaging image  ${PACKAGING}
    Call API                    deployment post  ${RES_DIR}/valid/deployment.create.yaml  ${image}
    ${exp_result}               create List   Ready
    ${result}                   Wait until command finishes and returns result  deployment  entity=${DEPLOYMENT}  exp_result=${exp_result}
    Status State Should Be      ${result}  Ready

Update deployment
    [Tags]                      deployment
    ${image}                    Pick packaging image  ${PACKAGING}
    Call API                    deployment put  ${RES_DIR}/valid/deployment.update.json  ${image}
    ${check_changes}            Call API  deployment get id  ${DEPLOYMENT}
    should be equal             ${check_changes.spec.role_name}  test_updated
    ${exp_result}               create List   Ready
    ${result}                   Wait until command finishes and returns result  deployment  entity=${DEPLOYMENT}  exp_result=${exp_result}
    Status State Should Be      ${result}  Ready
    CreatedAt and UpdatedAt times should not be equal  ${result}

Check by id that deployment exists
    [Tags]                      deployment
    ${result}                   Call API  deployment get id  ${DEPLOYMENT}
    ID should be equal          ${result}  ${DEPLOYMENT}

Check that list of routes contains
    [Tags]                      route
    Command response list should contain id  route  ${MODEL}

Check by id that route exists
    [Tags]                      route
    ${result}                   Call API  route get id  ${MODEL}
    ID should be equal          ${result}  ${MODEL}

Check existance of model route by id
    [Tags]                      route
    ${result}                   Call API  route get id  ${MODEL}
    ID should be equal          ${result}  ${MODEL}

Get info about model
    [Tags]                      model
    ${result}                   Call API  model get  url=${MODEL_URL}
    should be equal             ${result['info']['description']}  This is a EDI server.

Invoke model
    [Tags]                        model
    ${result}                     Call API  model post  url=${MODEL_URL}  json_input=${REQUEST}
    ${expected response}          evaluate  ${REQUEST_RESPONSE}
    dictionaries should be equal  ${result}  ${expected response}

Delete Model Trainings and Check that Model Trainging do not exist
    [Tags]                      training
    [Documentation]             delete model trainings
    Command response list should contain id  training  ${TRAIN_MLFLOW_DEFAULT}  ${TRAIN_MLFLOW_NOT_DEFAULT}

    Call API                    training delete  ${TRAIN_MLFLOW_DEFAULT}
    Call API                    training delete  ${TRAIN_MLFLOW_NOT_DEFAULT}

    Command response list should not contain id  training  ${TRAIN_MLFLOW_DEFAULT}  ${TRAIN_MLFLOW_NOT_DEFAULT}

Delete Model Packaging and Check that Model Packaging does not exist
    [Tags]                      packaging
    Command response list should contain id  packaging  ${PACKAGING}
    Call API                    packaging delete  ${PACKAGING}
    Command response list should not contain id  packaging  ${PACKAGING}

Delete Model Deployment and Check that Model Deployment does not exist
    [Tags]                      deployment
    [Documentation]             check that after deletion of deployment the model and route are also deleted
    Command response list should contain id  deployment  ${DEPLOYMENT}
    Call API                    deployment delete  ${DEPLOYMENT}
    Wait until delete finished  deployment  entity=${DEPLOYMENT}
    Command response list should not contain id  deployment  ${DEPLOYMENT}
    Command response list should not contain id  route  ${MODEL}
    Call API and get Error      ${WrongHttpStatusCode}  deployment get id  ${DEPLOYMENT}
