** Variables ***
${RES_DIR}             ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_feedback
${MD_FEEDBACK_MODEL}   feedback-model
${TEST_MODEL_NAME}     feedback
${TEST_MODEL_VERSION}  5.5

*** Settings ***
Documentation       API check for configuration
Resource            ../../../resources/keywords.robot
Resource            ../../../resources/variables.robot
Variables           ../../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             Collections
Library             odahuflow.robot.libraries.feedback.Feedback  ${CLOUD_TYPE}  ${FEEDBACK_BUCKET}  ${CLUSTER_NAME}
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.sdk.clients.configuration.ConfigurationClient
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup resources  AND
...                 Run API deploy from model packaging  ${MP_FEEDBACK_MODEL}  ${MD_FEEDBACK_MODEL}  ${RES_DIR}/simple-model.deployment.odahuflow.yaml  AND
...                 Check model started  ${MD_FEEDBACK_MODEL}
Suite Teardown      Run keywords  Cleanup resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          deployment  api  cli  feedback

*** Variables ***
${REQUEST_ID_CHECK_RETRIES}         30
@{FORBIDDEN_HEADERS}  authorization  x-jwt  x-user  x-email

*** Keywords ***
Call API
    StrictShell
