*** Variables ***
${LOCAL_CONFIG}                     odahuflow/api_deploy_route_model
${RES_DIR}                          ${CURDIR}/resources/deploy_route_model

${PACKAGING}                        simple-model
${DEPLOYMENT}                       wine-api-testing
${MODEL}                            ${DEPLOYMENT}
${MODEL_URL}                        ${EDGE_URL}/model/${MODEL}
${REQUEST}                          SEPARATOR=
...                                 { "columns": [ "a", "b" ], "data": [ [ 1.0, 2.0 ] ] }
${REQUEST_RESPONSE}                 { "prediction": [ [ 42 ] ], "columns": [ "result" ] }
${WrongHttpStatusCode}              SEPARATOR=
...                                 WrongHttpStatusCode: Got error from server: entity "{entity name}" is not found (status: 404)

*** Settings ***
Documentation       API of training, packaging, deployment, route and model
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Library             odahuflow.robot.libraries.sdk_wrapper.ModelDeployment
Library             odahuflow.robot.libraries.sdk_wrapper.ModelRoute
Library             odahuflow.robot.libraries.sdk_wrapper.Model
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup Resources
Suite Teardown      Run Keywords
...                 Cleanup Resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk
Test Timeout        15 minutes

*** Keywords ***
Cleanup Resources
    [Documentation]  Deletes of created resources
    StrictShell  odahuflowctl --verbose dep delete --id ${DEPLOYMENT} --ignore-not-found

Format WrongHttpStatusCode
    [Arguments]                     ${entity name}
    ${error output}                 format string  ${WrongHttpStatusCode}  entity name=${entity name}
    [return]                        ${error output}

*** Test Cases ***
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
    Check model started  ${DEPLOYMENT}
    Status State Should Be      ${result}  Ready

Update deployment
    [Tags]                      deployment
    ${image}                    Pick packaging image  ${PACKAGING}
    Call API                    deployment put  ${RES_DIR}/valid/deployment.update.json  ${image}
    ${check_changes}            Call API  deployment get id  ${DEPLOYMENT}
    should be equal             ${check_changes.spec.role_name}  test_updated
    ${exp_result}               create List   Ready
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
    Command response list should contain id  route  ${MODEL}

Check by id that route exists
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
