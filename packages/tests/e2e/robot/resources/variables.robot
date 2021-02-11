*** Variables ***
${ODAHUFLOW_NAMESPACE}                odahu-flow
${ODAHUFLOW_DEPLOYMENT_NAMESPACE}     odahu-flow-deployment
@{TEST_MODELS}                        Digit-Recognition  Test-Summation  Sklearn-Income
${MODEL_WITH_PROPS}                   Test-Summation
${MODEL_WITH_PROPS_ENDPOINT}          sum_and_pow
${MODEL_WITH_PROPS_PROP}              number.pow_of_ten
${TEST_MODEL_RESULT}                  42.0
${FEEDBACK_TAG}                       tag1
${TEST_MODEL_ARG_COPIES}              ${3000}
${TEST_MODEL_ARG_STR}                 test-model-invocation-string__
${FEEDBACK_LOCATION_MODELS_META_LOG}  model_log/request_response
${FEEDBACK_LOCATION_MODELS_RESP_LOG}  model_log/response_body
${FEEDBACK_LOCATION_MODELS_FEEDBACK}  model_log/feedback
${FEEDBACK_PARTITIONING_PATTERN}      year=%Y/month=%m/day=%d/%Y%m%d%H
${TEST_VCS}                           odahuflow
${ODAHUFLOW_ENTITIES_DIR}             ${CURDIR}/entities
${NODE_TAINT_KEY}                     dedicated
${NODE_TAINT_VALUE}                   jenkins-slave
${VCS_CONNECTION}                     odahu-flow-examples
${MP_SIMPLE_MODEL}                    simple-model
${MP_FAIL_MODEL}                      fail
${MP_COUNTER_MODEL}                   counter
${MP_FEEDBACK_MODEL}                  feedback
${TOOLCHAIN_INTEGRATION}              mlflow
${PI_REST}                            docker-rest
${PI_CLI}                             docker-cli
${CONN_SECRET_MASK}                   *****

# ---------------------------------  Error Templates  ---------------------------------
${400 BadRequest Template}          WrongHttpStatusCode: Got error from server: {} (status: 400)
${403 Forbidden Template}           WrongHttpStatusCode: Got error from server: {} (status: 403)
${404 NotFound Template}            WrongHttpStatusCode: Got error from server: entity "{}" is not found (status: 404)
${404 Model NotFoundTemplate}       Wrong status code returned: 404. Data: "". URL: "{}"
${409 Conflict Template}            EntityAlreadyExists: Got error from server: entity "{}" already exists (status: 409)
${Model WrongStatusCode Template}   Wrong status code returned: {status code}. Data: "{data}". URL: "{url}"
${APIConnectionException}           APIConnectionException: Can not reach {base url}
${IncorrectRefreshToken}            IncorrectAuthorizationToken: Refresh token is not correct.\nPlease login again
${IncorrectCredentials}             IncorrectAuthorizationToken: Credentials are not correct.\nPlease provide correct temporary token or disable non interactive mode
${IncorrectTemporaryToken}          IncorrectAuthorizationToken: Credentials are missed.\nPlease provide correct temporary token or disable non interactive mode

# ---------------------------------  Validation checks  ---------------------------------
${FailedConn}       Validation of connection is failed:
${FailedTI}         Validation of toolchain integration is failed:
${FailedPI}         Validation of packaging integration is failed:
${FailedTrain}      Validation of model training is failed:
${FailedPack}       Validation of model packaging is failed:
${FailedDeploy}     Validation of model deployment is failed:

${invalid_id}       ID is not valid
${empty_id}         empty "ID"

# ---------------------------------  connections  ---------------------------------
${empty_uri}                    empty uri
${s3_empty_keyID_keySecret}     s3 type requires that keyID and keySecret parameters must be non-empty
${gcs_empty_keySecret}          gcs type requires that keySecret parameter must be non-empty
${azureblob_req_keySecret}      azureblob type requires that keySecret parameter contains HTTP endpoint with SAS Token
${ecr_empty_keyID_keySecret}    ecr type requires that keyID and keySecret parameters must be non-empty
# ---------------------------------  toolchain  ---------------------------------
${TI_empty_entrypoint}          empty entrypoint
${TI_empty_defaultImage}        empty defaultImage
# ---------------------------------  packager  ---------------------------------
${PI_empty_entrypoint}          empty entrypoint
${PI_empty_defaultImage}        empty defaultImage
# ---------------------------------  training  ---------------------------------
${empty_model_name}             empty model.name
${empty_model_version}          empty model.version
${empty_VCS}                    empty vcsName
${empty_toolchain}              empty toolchain parameter
# ---------------------------------  packaging  ---------------------------------
${empty_artifactName}           empty artifactName
${empty_integrationName}        empty integrationName
# ---------------------------------  deployment  ---------------------------------
${max_smaller_min_replicas}     maximum number of replicas parameter must not be less than minimum number of replicas parameter
${empty_image}                  empty image parameter
${empty_predictor}              empty predictor parameter
${positive_livenessProbe}       livenessProbeInitialDelay must be non-negative integer
${positive_readinessProbe}      readinessProbeInitialDelay must be non-negative integer
${min_num_of_max_replicas}      maximum number of replicas parameter must not be less than 1
${min_num_of_min_replicas}      minimum number of replicas parameter must not be less than 0
