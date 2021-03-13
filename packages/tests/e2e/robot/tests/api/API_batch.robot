*** Variables ***
${LOCAL_CONFIG}                     odahuflow/batch
${RES_DIR}                          ${CURDIR}/resources/batch


*** Settings ***
Documentation       API of batch
Resource            ../../resources/keywords.robot
Resource            ./resources/keywords.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.sdk_wrapper.InferenceService
Library             odahuflow.robot.libraries.sdk_wrapper.InferenceJob
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge
Force Tags          api  batch
Test Timeout        15 minutes



*** Test cases ***
Create Batch Service
    [Tags]                      batch
    [Documentation]             create batch service
    ${result}                   Call API  service post  ${RES_DIR}/inferenceservice.yaml

Create Batch Job
    [Tags]                      batch
    [Documentation]             launch batch job
    ${job_id}                   Call API  job post  ${RES_DIR}/inferencejob.yaml
    @{exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  job  entity=${job_id}  exp_result=@{exp_result}
    Status State Should Be      ${result}  succeeded
