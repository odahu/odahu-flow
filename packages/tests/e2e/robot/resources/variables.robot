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
${400 BadRequest Template}         WrongHttpStatusCode: Got error from server: {} (status: 400)
${401 Unathorized Template}
${403 Forbidden Template}
${404 NotFound Template}           WrongHttpStatusCode: Got error from server: entity "{}" is not found (status: 404)
${404 Model NotFoundTemplate}      Wrong status code returned: 404. Data: . URL: {}
${409 Conflict Template}           EntityAlreadyExists: Got error from server: entity "{}" already exists (status: 409)

${APIConnectionException}          APIConnectionException: Can not reach {base url}
${IncorrectToken}                  IncorrectAuthorizationToken: Refresh token is not correct.\nPlease login again

# ---------------------------------  Validation checks  ---------------------------------
${FailedConn}   Validation of connection is failed:
${FailedTI}     Validation of toolchain integration is failed:
${FailedPI}     Validation of packaging integration is failed:
${FailedTrain}  Validation of model training is failed:
${FailedPack}   Validation of model packaging is failed:


${invalid_id}   ID is not valid
# ---------------------------------  connections  ---------------------------------
@{connection types}            s3  gcs  azureblob  git  docker  ecr
${unknown type}                unknown type: . Supported types: [s3 gcs azureblob git docker ecr]
${empty_uri}                   the uri parameter is empty
${ecr_invalid_uri}             not valid uri for ecr type: docker-credential-ecr-login can only be used with Amazon Elastic Container Registry.
${s3_empty_keyID_keySecret}    s3 type requires that keyID and keySecret parameters must be non-empty
${gcs_empty_keySecret}         gcs type requires that keySecret parameter must be non-empty
${azureblob_req_keySecret}     azureblob type requires that keySecret parameter containsHTTP endpoint with SAS Token
${ecr_empty_keyID_keySecret}   ecr type requires that keyID and keySecret parameters must be non-empty
# ---------------------------------  toolchain  ---------------------------------
${TI_empty_entrypoint}         entrypoint must be no empty
${TI_empty_defaultImage}       defaultImage must be no empty
# ---------------------------------  packager  ---------------------------------
${PI_empty_entrypoint}         entrypoint must be nonempty
${PI_empty_defaultImage}       default image must be nonempty
# ---------------------------------  training  ---------------------------------
${empty_model_name}            model name must be non-empty
${empty_model_version}         model version must be non-empty
${empty_VCS}                   VCS name is empty
${empty_toolchain}             toolchain parameter is empty
# ---------------------------------  packaging  ---------------------------------
${empty_artifactName}         you should specify artifactName
${empty_integrationName}      integration name must be nonempty
# ---------------------------------  deployment  ---------------------------------
${max_smaller_min_replicas}     maximum number of replicas parameter must not be less than minimum number of replicas parameter
${empty_image}                  the image parameter is empty
${positive_livenessProbe}       liveness probe parameter must be positive number
${positive_readinessProbe}      readiness probe must be positive number
${min_num_of_max_replicas}      maximum number of replicas parameter must not be less than 1
${min_num_of_min_replicas}      minimum number of replicas parameter must not be less than 0

