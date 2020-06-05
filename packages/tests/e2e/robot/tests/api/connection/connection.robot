*** Variables ***
${RES_DIR}             ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_feedback
${MD_FEEDBACK_MODEL}   feedback-model
${TEST_MODEL_NAME}     feedback
${TEST_MODEL_VERSION}  5.5

*** Settings ***
Documentation       Feedback loop (fluentd) check
Resource            ../resources/keywords.robot
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             Collections
Library             odahuflow.robot.libraries.feedback.Feedback  ${CLOUD_TYPE}  ${FEEDBACK_BUCKET}  ${CLUSTER_NAME}
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup resources  AND
...                 Run API deploy from model packaging  ${MP_FEEDBACK_MODEL}  ${MD_FEEDBACK_MODEL}  ${RES_DIR}/simple-model.deployment.odahuflow.yaml  AND
...                 Check model started  ${MD_FEEDBACK_MODEL}
Suite Teardown      Run keywords  Cleanup resources  AND
...                 Remove File  ${LOCAL_CONFIG}
Force Tags          testing


*** Test Cases ***
testing of API for connection
    Call API               connection get
    Strict Call API       connection get