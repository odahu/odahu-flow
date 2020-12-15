*** Variables ***
${ODAHUFLOW_NAMESPACE}                  odahu-flow
${ODAHUFLOW_DEPLOYMENT_NAMESPACE}       odahu-flow-deployment
@{TEST_MODELS}                          Digit-Recognition  Test-Summation  Sklearn-Income
${MODEL_WITH_PROPS}                     Test-Summation
${MODEL_WITH_PROPS_ENDPOINT}            sum_and_pow
${MODEL_WITH_PROPS_PROP}                number.pow_of_ten
${TEST_MODEL_RESULT}                    42.0
${FEEDBACK_TAG}                         tag1
${TEST_MODEL_ARG_COPIES}                ${3000}
${TEST_MODEL_ARG_STR}                   test-model-invocation-string__
${FEEDBACK_LOCATION_MODELS_META_LOG}    model_log/request_response
${FEEDBACK_LOCATION_MODELS_RESP_LOG}    model_log/response_body
${FEEDBACK_LOCATION_MODELS_FEEDBACK}    model_log/feedback
${FEEDBACK_PARTITIONING_PATTERN}        year=%Y/month=%m/day=%d/%Y%m%d%H
${TEST_VCS}                             odahuflow
${ODAHUFLOW_ENTITIES_DIR}               ${CURDIR}/entities
${NODE_TAINT_KEY}                       dedicated
${NODE_TAINT_VALUE}                     jenkins-slave
${MP_SIMPLE_MODEL}                      simple-model
${MP_FAIL_MODEL}                        fail
${MP_COUNTER_MODEL}                     counter
${MP_FEEDBACK_MODEL}                    feedback
${CONN_SECRET_MASK}                     *****

# Wine Model
${WINE_MODEL_RESULT}                    {"prediction": [6.3881577909662886], "columns": ["quality"]}

# Errors
${INVALID_URL_ERROR}                    Error: Can not reach
${INVALID_CREDENTIALS_ERROR}            Error: Credentials are not correct.
${MISSED_CREDENTIALS_ERROR}             Error: Credentials are missed.