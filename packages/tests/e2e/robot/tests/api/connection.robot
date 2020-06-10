*** Variables ***
${RES_DIR}             ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_feedback
${MD_FEEDBACK_MODEL}   feedback-model
${TEST_MODEL_NAME}     feedback
${TEST_MODEL_VERSION}  5.5

*** Settings ***
Documentation       Feedback loop (fluentd) check
Resource            resources/keywords.robot
#  Variables           ../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge
Force Tags          testing


*** Test Cases ***
testing of API for connection
    Call API              connection get

testing strict API
    Strict Call API       connection get

testing strict API with arguments
    Call API              connection get id  docker-ci