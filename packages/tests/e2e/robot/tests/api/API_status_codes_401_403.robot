*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_status_codes_401_403
${RES_DIR}              ${CURDIR}/resources

${invalid_token}        not-valid-token
${NOT_EXIST_ENTITY}     not-exist

${REQUEST}              SEPARATOR=
...                     { "columns": [ "a", "b" ], "data": [ [ 1.0, 2.0 ] ] }

*** Settings ***
Documentation       tests for API status codes 401, 403
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             String
Library             odahuflow.robot.libraries.sdk_wrapper.Configuration
Library             odahuflow.robot.libraries.sdk_wrapper.Connection
Library             odahuflow.robot.libraries.sdk_wrapper.Toolchain
Library             odahuflow.robot.libraries.sdk_wrapper.Packager
Library             odahuflow.robot.libraries.sdk_wrapper.ModelTraining
Library             odahuflow.robot.libraries.sdk_wrapper.ModelPackaging
Library             odahuflow.robot.libraries.sdk_wrapper.ModelDeployment
Library             odahuflow.robot.libraries.sdk_wrapper.ModelRoute
Library             odahuflow.robot.libraries.sdk_wrapper.Model
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Shell  odahuflowctl config set API_URL ${API_URL}  AND
...                 Shell  odahuflowctl config set MODEL_HOST ${EDGE_URL}
# Suite Teardown      Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  negative  test
Test Timeout        1 minute

*** Keywords ***
Try Call API - Unathorized
    [Arguments]  ${command}  @{options}  &{keyword arguments}
    Call API and get Error  ${IncorrectTemporaryToken}  ${command}  @{options}  &{keyword arguments}  token=${EMPTY}
    Call API and get Error  ${IncorrectCredentials}  ${command}  @{options}  &{keyword arguments}  token=${invalid_token}

Try Call API - Forbidden
    [Arguments]  ${command}  @{options}  &{keyword arguments}
    Log many   ${API_URL}  ${EDGE_URL}
    ${403 Forbidden}  format string  ${403 Forbidden Template}   None
    Call API and get Error  ${403 Forbidden}  ${command}  @{options}  &{keyword arguments}


*** Test Cases ***
# Status Code 401 - Unathorized
#     [Template]  Try Call API - Unathorized
#     # config
#     config get
#     # connection
#     connection get
#     connection get id  ${VCS_CONNECTION}
#     connection get id decrypted  ${VCS_CONNECTION}
#     connection post  ${RES_DIR}/connection/valid/docker_connection_create.json
#     connection put  ${RES_DIR}/connection/valid/git_connection_update.yaml
#     connection delete  ${NOT_EXIST_ENTITY}
#     # toolchains
#     toolchain get
#     toolchain get id  ${TOOLCHAIN_INTEGRATION}
#     toolchain post  ${RES_DIR}/toolchain/valid/mlflow_create.yaml
#     toolchain put  ${RES_DIR}/toolchain/valid/mlflow_update.json
#     toolchain delete  ${NOT_EXIST_ENTITY}
#     # packagers
#     packager get
#     packager get id  ${PI_REST}
#     packager post  ${RES_DIR}/packager/valid/docker_rest_create.json
#     packager put  ${RES_DIR}/packager/valid/docker_rest_update.yaml
#     packager delete  ${NOT_EXIST_ENTITY}
#     # training
#     training get
#     training get id  ${NOT_EXIST_ENTITY}
#     training get log  ${NOT_EXIST_ENTITY}
#     training post  ${RES_DIR}/training_packaging/valid/training.mlflow.default.yaml
#     training put  ${RES_DIR}/training_packaging/valid/training.mlflow.default.yaml
#     training delete  ${NOT_EXIST_ENTITY}
#     # packaging
#     packaging get
#     packaging get id  ${NOT_EXIST_ENTITY}
#     packaging get log  ${NOT_EXIST_ENTITY}
#     packaging post  ${RES_DIR}/training_packaging/valid/packaging.create.yaml
#     packaging put  ${RES_DIR}/training_packaging/valid/packaging.create.yaml
#     packaging delete  ${NOT_EXIST_ENTITY}
#     # deployment
#     deployment get
#     deployment get id  ${NOT_EXIST_ENTITY}
#     deployment post  ${RES_DIR}/deploy_route_model/valid/deployment.create.yaml
#     deployment put  ${RES_DIR}/deploy_route_model/valid/deployment.create.yaml
#     deployment delete  ${NOT_EXIST_ENTITY}
#     # route
#     route get
#     route get id  ${NOT_EXIST_ENTITY}
#     route post  ${RES_DIR}/deploy_route_model/valid/route.yaml
#     route put  ${RES_DIR}/deploy_route_model/valid/route.yaml
#     route delete  ${NOT_EXIST_ENTITY}
#     # model
#     model get   url=${EDGE_URL}/model/${NOT_EXIST_ENTITY}
#     model post  url=${EDGE_URL}/model/${NOT_EXIST_ENTITY}  json_input=${REQUEST}


Status Code 403 - Forbidden - Data Scientist
    [Template]  Try Call API - Forbidden
    [Setup]  Login to the api and edge  ${SA_DATA_SCIENTIST}
    # [Teardown]  Remove File  ${LOCAL_CONFIG}
    connection get
    connection get id  ${VCS_CONNECTION}
    connection get id decrypted  ${VCS_CONNECTION}
    # toolchains
    toolchain post  ${RES_DIR}/toolchain/valid/mlflow_create.yaml
    toolchain put  ${RES_DIR}/toolchain/valid/mlflow_update.json
    toolchain delete  ${NOT_EXIST_ENTITY}
    # packagers
    packager post  ${RES_DIR}/packager/valid/docker_rest_create.json
    packager put  ${RES_DIR}/packager/valid/docker_rest_update.yaml
    packager delete  ${NOT_EXIST_ENTITY}
    # route
    route post  ${RES_DIR}/deploy_route_model/valid/route.yaml
    route put  ${RES_DIR}/deploy_route_model/valid/route.yaml
    route delete  ${NOT_EXIST_ENTITY}

# Status Code 403 - Forbidden - Viewer
#     [Template]  Template Error Keyword
#     [Setup]  Login to the api and edge  ${SA_DATA_SCIENTIST}
#     [Teardown]  Remove File  ${LOCAL_CONFIG}
#
# Status Code 403 - Forbidden - Custom Role
#     [Template]  Template Error Keyword
#     [Setup]  Login to the api and edge  ${SA_DATA_SCIENTIST}
#     [Teardown]  Remove File  ${LOCAL_CONFIG}
#