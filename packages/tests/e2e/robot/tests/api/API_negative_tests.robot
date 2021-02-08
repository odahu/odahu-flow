*** Variables ***
${LOCAL_CONFIG}         odahuflow/api_status_codes_400-403
${RES_DIR}              ${CURDIR}/resources
${invalid_token}        not-valid-token
${NOT_EXIST_ENTITY}     not-exist
${MODEL_DEPLOYMENT}     dep-status-code-400-403

${REQUEST}              SEPARATOR=
...                     { "columns": [ "a", "b" ], "data": [ [ 1.0, 2.0 ] ] }


*** Settings ***
Documentation       tests for API status codes 400, 403
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             String
Library             odahuflow.robot.libraries.sdk_wrapper.Login
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
...                 Login to the api and edge  AND
...                 reload config  AND
...                 Run API deploy from model packaging and check model started  simple-model  dep-status-code-400-403  ${RES_DIR}/deploy_route_model/valid/deployment.negative.403.yaml
Suite Teardown      Run Keywords
...                 Login to the api and edge  AND
...                 reload config  AND
...                 Run API undeploy model and check  ${MODEL_DEPLOYMENT}  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          api  sdk  negative  status-code-400-403
Test Timeout        1 minute

*** Keywords ***
Try Call API - Bad Request
    [Arguments]  ${format_string}  ${error}  ${command}  @{options}
    ${error}        format string   ${format_string}  ${error}
    Call API and get Error  ${error}  ${command}  @{options}

Try Call API - Forbidden
    [Arguments]  ${command}  @{options}  &{kwargs}
    Log many   ${API_URL}  ${EDGE_URL}
    ${403 Forbidden}  format string  ${403 Forbidden Template}  Forbidden
    Call API and get Error  ${403 Forbidden}  ${command}  @{options}  &{kwargs}

Try Call API - Forbidden.Model
    [Arguments]  ${command}  &{kwargs}
    ${403 Forbidden}  format string  ${Model WrongStatusCode Template}  status code=403  data=${EMPTY}  url=${kwargs.get("url")}${kwargs.pop("path_suffix")}
    Call API and get Error  ${403 Forbidden}  ${command}  &{kwargs}

*** Test Cases ***
Status Code 400 - Bad Request
    [Template]  Try Call API - Bad Request
    # connection
    ${400 BadRequest Template}  ${FailedConn} ${invalid_id}
    ...  connection put   ${RES_DIR}/connection/invalid/conn_invalid_id.yaml
    ${400 BadRequest Template}  ${FailedConn} ${s3_empty_keyID_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/s3_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${s3_empty_keyID_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/s3_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${gcs_empty_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/gcs_no_required_parameters
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${gcs_empty_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/gcs_no_required_parameters
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${azureblob_req_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/azureblob_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${azureblob_req_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/azureblob_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${invalid_id}; ${empty_id}; ${empty_uri}
     ...  connection post  ${RES_DIR}/connection/invalid/git_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${invalid_id}; ${empty_id}; ${empty_uri}
     ...  connection put   ${RES_DIR}/connection/invalid/git_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}
    ...  connection post  ${RES_DIR}/connection/invalid/docker_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}
    ...  connection put   ${RES_DIR}/connection/invalid/docker_no_required_parameters.yaml
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${ecr_empty_keyID_keySecret}
    ...  connection post  ${RES_DIR}/connection/invalid/ecr_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedConn} ${empty_uri}; ${ecr_empty_keyID_keySecret}
    ...  connection put   ${RES_DIR}/connection/invalid/ecr_no_required_parameters.json
    # toolchains
    ${400 BadRequest Template}  ${FailedTI} ${TI_empty_entrypoint}; ${TI_empty_defaultImage}
    ...  toolchain post  ${RES_DIR}/toolchain/invalid/toolchain_no_required_parameters.json
    ${400 BadRequest Template}  ${FailedTI} ${TI_empty_entrypoint}
    ...  toolchain put  ${RES_DIR}/toolchain/invalid/toolchain_no_required_parameters.yaml
    # packagers
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager post  ${RES_DIR}/packager/invalid/cli_no_required_params.yaml
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager put  ${RES_DIR}/packager/invalid/cli_no_required_params.yaml
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager post  ${RES_DIR}/packager/invalid/rest_no_required_params.json
    ${400 BadRequest Template}  ${FailedPI} ${PI_empty_entrypoint}; ${PI_empty_defaultImage}
    ...  packager put  ${RES_DIR}/packager/invalid/rest_no_required_params.json
    # model training
    ${400 BadRequest Template}  ${FailedTrain} ${empty_model_name}; ${empty_model_version}; ${empty_VCS}; ${empty_toolchain}
    ...  training post  ${RES_DIR}/training_packaging/invalid/training_no_required_params.yaml
    ${400 BadRequest Template}  ${FailedTrain} ${empty_model_name}; ${empty_model_version}; ${empty_VCS}; ${empty_toolchain}
    ...  training put  ${RES_DIR}/training_packaging/invalid/training_no_required_params.yaml
    # model packaging
    ${400 BadRequest Template}  ${FailedPack} ${empty_artifactName}; ${empty_integrationName}
    ...  packaging post  ${RES_DIR}/training_packaging/invalid/packaging_no_required_params.json
    ${400 BadRequest Template}  ${FailedPack} ${empty_artifactName}; ${empty_integrationName}
    ...  packaging put  ${RES_DIR}/training_packaging/invalid/packaging_no_required_params.json
    # model deployment
    ${400 BadRequest Template}  ${FailedDeploy} ${max_smaller_min_replicas}; ${invalid_id}; ${empty_id}; ${empty_image}; ${empty_predictor}
    ...  deployment post  ${RES_DIR}/deploy_route_model/invalid/deployment_empty_required_params.json
    ${400 BadRequest Template}  ${FailedDeploy} ${max_smaller_min_replicas}; ${invalid_id}; ${empty_id}; ${empty_image}; ${empty_predictor}
    ...  deployment put  ${RES_DIR}/deploy_route_model/invalid/deployment_empty_required_params.json
    ${400 BadRequest Template}  ${FailedDeploy} ${positive_livenessProbe}; ${positive_readinessProbe}; ${max_smaller_min_replicas}; ${min_num_of_max_replicas}; ${min_num_of_min_replicas}
    ...  deployment post  ${RES_DIR}/deploy_route_model/invalid/deployment_validation_checks.yaml
    ${400 BadRequest Template}  ${FailedDeploy} ${positive_livenessProbe}; ${positive_readinessProbe}; ${max_smaller_min_replicas}; ${min_num_of_max_replicas}; ${min_num_of_min_replicas}
    ...  deployment put  ${RES_DIR}/deploy_route_model/invalid/deployment_validation_checks.yaml

Status Code 403 - Forbidden - Data Scientist
    [Template]  Try Call API - Forbidden
    [Setup]     run keywords
    ...         Login to the api and edge  ${SA_DATA_SCIENTIST}  AND
    ...         reload config
    [Teardown]  Remove File  ${LOCAL_CONFIG}
    # connection
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

Status Code 403 - Forbidden - Viewer
    [Template]  Try Call API - Forbidden
    [Setup]     run keywords
    ...         Login to the api and edge  ${SA_VIEWER}  AND
    ...         reload config
    [Teardown]  Remove File  ${LOCAL_CONFIG}
    # connection
    connection get id decrypted  ${VCS_CONNECTION}
    connection post  ${RES_DIR}/connection/valid/docker_connection_create.json
    connection put  ${RES_DIR}/connection/valid/git_connection_update.yaml
    connection delete  ${NOT_EXIST_ENTITY}
    # toolchains
    toolchain post  ${RES_DIR}/toolchain/valid/mlflow_create.yaml
    toolchain put  ${RES_DIR}/toolchain/valid/mlflow_update.json
    toolchain delete  ${NOT_EXIST_ENTITY}
    # packaging
    packaging post  ${RES_DIR}/training_packaging/valid/packaging.create.yaml
    packaging put  ${RES_DIR}/training_packaging/valid/packaging.create.yaml
    packaging delete  ${NOT_EXIST_ENTITY}

Model. Status Code 403 - Forbidden - Viewer
    [Template]  Try Call API - Forbidden.Model
    [Setup]     run keywords
    ...         Login to the api and edge  ${SA_VIEWER}  AND
    ...         reload config
    [Teardown]  Remove File  ${LOCAL_CONFIG}
    # model

    model get   url=${EDGE_URL}/model/${MODEL_DEPLOYMENT}  path_suffix=/api/model/info
    model post  url=${EDGE_URL}/model/${MODEL_DEPLOYMENT}  path_suffix=/api/model/invoke  json_input=${REQUEST}

Status Code 403 - Forbidden - Custom Role
    [Template]  Try Call API - Forbidden
    [Setup]     run keywords
    ...         Login to the api and edge  ${SA_CUSTOM_USER}  AND
    ...         reload config
    [Teardown]  Remove File  ${LOCAL_CONFIG}
    # config
    config get
    # connection
    connection get
    connection get id  ${VCS_CONNECTION}
    connection get id decrypted  ${VCS_CONNECTION}
    connection post  ${RES_DIR}/connection/valid/docker_connection_create.json
    connection put  ${RES_DIR}/connection/valid/git_connection_update.yaml
    connection delete  ${NOT_EXIST_ENTITY}
    # packagers
    packager get
    packager get id  ${PI_REST}
    packager post  ${RES_DIR}/packager/valid/docker_rest_create.json
    packager put  ${RES_DIR}/packager/valid/docker_rest_update.yaml
    packager delete  ${NOT_EXIST_ENTITY}
    # training
    training get
    training get id  ${NOT_EXIST_ENTITY}
    training get log  ${NOT_EXIST_ENTITY}
    training post  ${RES_DIR}/training_packaging/valid/training.mlflow.default.yaml
    training put  ${RES_DIR}/training_packaging/valid/training.mlflow.default.yaml
    training delete  ${NOT_EXIST_ENTITY}
    # deployment
    deployment get
    deployment get id  ${NOT_EXIST_ENTITY}
    deployment post  ${RES_DIR}/deploy_route_model/valid/deployment.create.yaml
    deployment put  ${RES_DIR}/deploy_route_model/valid/deployment.create.yaml
    deployment delete  ${NOT_EXIST_ENTITY}

Model. Status Code 403 - Forbidden - Custom Role
    [Template]  Try Call API - Forbidden.Model
    [Setup]     run keywords
    ...         Login to the api and edge  ${SA_CUSTOM_USER}  AND
    ...         reload config
    [Teardown]  Remove File  ${LOCAL_CONFIG}
    # model
    model get   url=${EDGE_URL}/model/${MODEL_DEPLOYMENT}  path_suffix=/api/model/info
    model post  url=${EDGE_URL}/model/${MODEL_DEPLOYMENT}  path_suffix=/api/model/invoke  json_input=${REQUEST}
