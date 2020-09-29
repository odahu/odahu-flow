*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_deploy_route_model
${RES_DIR}                          ${CURDIR}/resources/deploy_route_model

${PACKAGING}                        simple-model
${DEPLOYMENT}                       wine-api-testing
${MODEL}                            ${DEPLOYMENT}
${REQUEST}                          SEPARATOR=
...                                 { "columns": [ "a", "b" ], "data": [ [ 1.0, 2.0 ] ] }
${REQUEST_RESPONSE}                 { "prediction": [ [ 42 ] ], "columns": [ "result" ] }
${WrongHttpStatusCode}              SEPARATOR=
...                                 WrongHttpStatusCode: Got error from server: entity "{entity name}" is not found (status: 404)
${WrongStatusCodeReturned}          SEPARATOR=
...                                 Wrong status code returned: 404. Data: . URL: {model url}

${DEPLOYMENT_NOT_EXIST}             deployment-api-not-exist
${MODEL_NOT_EXIST}                  ${DEPLOYMENT_NOT_EXIST}

*** Settings ***
Documentation       API of deployment, route and model
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
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

Format WrongHttpStatusCode
    [Arguments]       ${entity name}
    ${error output}   format string  ${WrongHttpStatusCode}  entity name=${entity name}
    [return]          ${error output}

Format WrongStatusCodeReturned
    [Arguments]       ${model url}
    ${error output}   format string  ${WrongStatusCodeReturned}  model url=${model url}
    [return]          ${error output}

*** Test Cases ***
Check deployment doesn't exist
    [Tags]                      deployment
    Command response list should not contain id  deployment  ${DEPLOYMENT}

Check route doesn't exist before deployment
    [Tags]                      route
    Command response list should not contain id  route  ${MODEL}

Create deployment
    [Tags]                      deployment
    ${image}                    Pick packaging_image  ${PACKAGING}
    Call API                    deployment post  ${RES_DIR}/valid/deployment.create.yaml  ${image}
    ${exp_result}               create List   Ready
    ${result}                   Wait until command finishes and returns result  deployment  entity=${DEPLOYMENT}  exp_result=${exp_result}
    Status State Should Be      ${result}  Ready

Update deployment
    [Tags]                      deployment
    ${image}                    Pick packaging_image  ${PACKAGING}
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
    ${model_url}                Get model Url  ${MODEL}
    ${result}                   Call API  model get  url=${model_url}
    should be equal             ${result['info']['description']}  This is a EDI server.

Invoke model
    [Tags]                        model
    ${model_url}                  Get model Url  ${MODEL}
    ${result}                     Call API  model post  url=${model_url}  json_input=${REQUEST}
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
    ${StatusCode}               Format WrongHttpStatusCode  ${DEPLOYMENT}
    Call API and get Error      ${StatusCode}  deployment get id  ${DEPLOYMENT}

#############################
#    NEGATIVE TEST CASES    #
#############################

#  DEPLOYMENT
#############
Try Create Deployment that already exists
    [Tags]                      deployment  negative
    [Setup]                     Cleanup resource  deployment  ${DEPLOYMENT}
    [Teardown]                  Cleanup resource  deployment  ${DEPLOYMENT}
    Call API                    deployment post  ${RES_DIR}/valid/deployment.update.json  packaging_image
    ${EntityAlreadyExists}      Format EntityAlreadyExists  ${DEPLOYMENT}
    Call API and get Error      ${EntityAlreadyExists}  deployment post  ${RES_DIR}/valid/deployment.create.yaml  packaging_image

Try Update not existing Deployment
    [Tags]                      deployment  negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DEPLOYMENT_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  deployment put  ${RES_DIR}/invalid/deployment.update.not_exist.json  packaging_image

Try Update deleted Deployment
    [Tags]                      deployment  negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DEPLOYMENT}
    Call API and get Error      ${WrongHttpStatusCode}  deployment put  ${RES_DIR}/valid/deployment.create.yaml  packaging_image

Try Get id not existing Deployment
    [Tags]                      deployment  negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DEPLOYMENT_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  deployment get id  ${DEPLOYMENT_NOT_EXIST}

Try Get id deleted Deployment
    [Tags]                      deployment  negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DEPLOYMENT}
    Call API and get Error      ${WrongHttpStatusCode}  deployment get id  ${DEPLOYMENT}

Try Delete not existing Deployment
    [Tags]                      deployment  negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DEPLOYMENT_NOT_EXIST}
    Call API and get Error      ${WrongHttpStatusCode}  deployment delete  ${DEPLOYMENT_NOT_EXIST}

Try Delete deleted Deployment
    [Tags]                      deployment  negative
    ${WrongHttpStatusCode}      Format WrongHttpStatusCode  ${DEPLOYMENT}
    Call API and get Error      ${WrongHttpStatusCode}  deployment delete  ${DEPLOYMENT}

#  MODEL
#############
Try Get info not existing Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT_NOT_EXIST}
    ${WrongStatusCodeReturned}  Format WrongStatusCodeReturned  ${model_url}/api/model/info
    Call API and get Error      ${WrongStatusCodeReturned}  model get  url=${model_url}

Try Get info deleted Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT}
    ${WrongStatusCodeReturned}  Format WrongStatusCodeReturned  ${model_url}/api/model/info
    Call API and get Error      ${WrongStatusCodeReturned}  model get  url=${model_url}

Try Invoke not existing and deleted Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT_NOT_EXIST}
    ${WrongStatusCodeReturned}  Format WrongStatusCodeReturned  ${model_url}/api/model/invoke
    Call API and get Error      ${WrongStatusCodeReturned}  model post  url=${model_url}  json_input=${REQUEST}

Try Invoke deleted Model
    [Tags]                      model  negative
    ${model_url}                Get model Url  ${DEPLOYMENT}
    ${WrongStatusCodeReturned}  Format WrongStatusCodeReturned  ${model_url}/api/model/invoke
    Call API and get Error      ${WrongStatusCodeReturned}  model post  url=${model_url}  json_input=${REQUEST}
