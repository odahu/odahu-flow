*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_deploy_route_model
${RES_DIR}                          ${CURDIR}/resources/deploy_route_model

${PACKAGING}                        simple-model
${DEPLOYMENT}                       wine-api-testing
${MODEL}                            ${DEPLOYMENT}
${REQUEST}                          SEPARATOR=
...                                 { "columns": [ "a", "b" ], "data": [ [ 1.0, 2.0 ] ] }
${REQUEST_RESPONSE}                 { "prediction": [ [ 42 ] ], "columns": [ "result" ] }
${DEPLOYMENT_NOT_EXIST}             deployment-api-not-exist
${MODEL_NOT_EXIST}                  ${DEPLOYMENT_NOT_EXIST}

*** Settings ***
Documentation       API of deployment, route and model
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.Login
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Library             odahuflow.robot.libraries.sdk_wrapper.ModelDeployment
Library             odahuflow.robot.libraries.sdk_wrapper.ModelRoute
Library             odahuflow.robot.libraries.sdk_wrapper.Model
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
    Cleanup resource  deployment  ${DEPLOYMENT}
    Cleanup resource  deployment  ${DEPLOYMENT_NOT_EXIST}

Get model Url
    [Arguments]       ${model_id}
    ${model_url}      set variable  ${EDGE_URL}/model/${model_id}
    [return]          ${model_url}

*** Test Cases ***
Check deployment doesn't exist
    [Tags]                      deployment
    Command response list should not contain id  deployment  ${DEPLOYMENT}

Create deployment
    [Tags]                      deployment
    ${image}                    Pick packaging image  ${PACKAGING}
    Call API                    deployment post  ${RES_DIR}/valid/deployment.create.yaml  ${image}
    ${exp_result}               create List   Ready
    ${result}                   Wait until command finishes and returns result  deployment  entity=${DEPLOYMENT}  exp_result=${exp_result}
    Check model started         ${DEPLOYMENT}
    Status State Should Be      ${result}  Ready

Update deployment
    [Tags]                      deployment
    ${image}                    Pick packaging image  ${PACKAGING}
    Call API                    deployment put  ${RES_DIR}/valid/deployment.update.json  ${image}
    ${check_changes}            Call API  deployment get id  ${DEPLOYMENT}
    should be equal             ${check_changes.spec.role_name}  test_updated
    ${exp_result}               create list   Ready
    ${result}                   Wait until command finishes and returns result  deployment  entity=${DEPLOYMENT}  exp_result=${exp_result}
    Check model started  ${DEPLOYMENT}
    Status State Should Be      ${result}  Ready
    CreatedAt and UpdatedAt times should not be equal  ${result}

Check by id that deployment exists
    [Tags]                      deployment
    ${result}                   Call API  deployment get id  ${DEPLOYMENT}
    ID should be equal          ${result}  ${DEPLOYMENT}

Check that list of routes contains
    [Tags]                      route
    ${result}                   Call API  deployment get default route  ${MODEL}
    Command response list should contain id  route  ${result.id}

Check by id that route exists
    [Tags]                      route
    ${default_route}            Call API  deployment get default route  ${MODEL}
    ${result}                   Call API  route get id  ${default_route.id}
    ID should be equal          ${result}  ${default_route.id}


Get info about model
    [Tags]                      model
    ${model_url}                Get model Url  ${MODEL}
    ${result}                   Call API  model get  url=${model_url}  token=${AUTH_TOKEN}
    should be equal             ${result['info']['description']}  This is a EDI server.

Invoke model
    [Tags]                        model
    ${model_url}                  Get model Url  ${MODEL}
    ${result}                     Call API  model post  url=${model_url}  token=${AUTH_TOKEN}  json_input=${REQUEST}
    ${expected response}          evaluate  ${REQUEST_RESPONSE}
    dictionaries should be equal  ${result}  ${expected response}

Delete Model Deployment and Check that Model Deployment does not exist
    [Tags]                      deployment
    [Documentation]             check that after deletion of deployment the model and route are also deleted
    Command response list should contain id  deployment  ${DEPLOYMENT}
    Call API                    deployment delete  ${DEPLOYMENT}
    Wait until delete finished  deployment  entity=${DEPLOYMENT}
    Command response list should not contain id  deployment  ${DEPLOYMENT}
    Command response list should not contain id  route  ${MODEL}
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT}
    Call API and get Error      ${404NotFound}  deployment get id  ${DEPLOYMENT}



#############################
#    NEGATIVE TEST CASES    #
#############################

#  DEPLOYMENT
#############
Try Create Deployment that already exists
    [Tags]                      deployment  negative
    [Setup]                     Cleanup resource  deployment  ${DEPLOYMENT}
    [Teardown]                  Cleanup resource  deployment  ${DEPLOYMENT}
    Call API                    deployment post  ${RES_DIR}/valid/deployment.update.json  ${PACKAGE_IMAGE_STUB}
    ${EntityAlreadyExists}      format string  ${409 Conflict Template}  ${DEPLOYMENT}
    Call API and get Error      ${EntityAlreadyExists}  deployment post  ${RES_DIR}/valid/deployment.create.yaml  ${PACKAGE_IMAGE_STUB}

Try Update not existing Deployment
    [Tags]                      deployment  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  deployment put  ${RES_DIR}/invalid/deployment.update.not_exist.json  ${PACKAGE_IMAGE_STUB}

Try Update deleted Deployment
    [Tags]                      deployment  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT}
    Call API and get Error      ${404NotFound}  deployment put  ${RES_DIR}/valid/deployment.create.yaml  ${PACKAGE_IMAGE_STUB}

Try Get id not existing Deployment
    [Tags]                      deployment  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  deployment get id  ${DEPLOYMENT_NOT_EXIST}

Try Get id deleted Deployment
    [Tags]                      deployment  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT}
    Call API and get Error      ${404NotFound}  deployment get id  ${DEPLOYMENT}

Try Delete not existing Deployment
    [Tags]                      deployment  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT_NOT_EXIST}
    Call API and get Error      ${404NotFound}  deployment delete  ${DEPLOYMENT_NOT_EXIST}

Try Delete deleted Deployment
    [Tags]                      deployment  negative
    ${404NotFound}              format string  ${404 NotFound Template}  ${DEPLOYMENT}
    Call API and get Error      ${404NotFound}  deployment delete  ${DEPLOYMENT}

#  MODEL
#############
Try Get info not existing Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT_NOT_EXIST}
    ${404ModelNotFound}         format string  ${404 Model NotFoundTemplate}  ${model_url}/api/model/info
    Call API and get Error      ${404ModelNotFound}  model get  url=${model_url}  token=${AUTH_TOKEN}

Try Get info deleted Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT}
    ${404ModelNotFound}         format string  ${404 Model NotFoundTemplate}  ${model_url}/api/model/info
    Call API and get Error      ${404ModelNotFound}  model get  url=${model_url}  token=${AUTH_TOKEN}

Try Invoke not existing and deleted Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT_NOT_EXIST}
    ${404ModelNotFound}         format string  ${404 Model NotFoundTemplate}  ${model_url}/api/model/invoke
    Call API and get Error      ${404ModelNotFound}  model post  url=${model_url}  token=${AUTH_TOKEN}  json_input=${REQUEST}

Try Invoke deleted Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT}
    ${404ModelNotFound}         format string  ${404 Model NotFoundTemplate}  ${model_url}/api/model/invoke
    Call API and get Error      ${404ModelNotFound}  model post  url=${model_url}  token=${AUTH_TOKEN}  json_input=${REQUEST}
