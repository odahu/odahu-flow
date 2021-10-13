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
Library             odahuflow.robot.libraries.batch.BatchUtils  ${CLOUD_TYPE}  ${TEST_BUCKET}  ${CLUSTER_NAME}
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge
Force Tags          api  batch
Test Timeout        15 minutes


*** Test Cases ***
Create Batch Service
    [Tags]                      batch
    [Documentation]             create batch service
    Call API  service post  ${RES_DIR}/inferenceservice.yaml

Create Batch Job
    [Tags]                      batch
    [Documentation]             launch batch job
    ${job_id}                   Call API  job post  ${RES_DIR}/inferencejob.yaml
    ${exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  job  entity=${job_id}  exp_result=${exp_result}
    Status State Should Be      ${result}  succeeded
    ${result}                   check batch job response  ${RES_DIR}/inferencejob.yaml  ${RES_DIR}/output/response0.json
    Should Be True              ${result}

Create Batch Service Packed Model
    [Tags]                      batch
    [Documentation]             create batch service
    Call API  service post  ${RES_DIR}/inferenceservice-packed.yaml

Create Batch Job Packed
    [Tags]                      batch
    [Documentation]             launch batch job
    ${job_id}                   Call API  job post  ${RES_DIR}/inferencejob-packed.yaml
    ${exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  job  entity=${job_id}  exp_result=${exp_result}
    Status State Should Be      ${result}  succeeded
    ${result}                   check batch job response  ${RES_DIR}/inferencejob-packed.yaml  ${RES_DIR}/output/response0.json
    Should Be True              ${result}


Create Batch Service Embedded Model
    [Tags]                      batch
    [Documentation]             create batch service
    Call API  service post  ${RES_DIR}/inferenceservice-embedded.yaml

Create Batch Job Embedded
    [Tags]                      batch
    [Documentation]             launch batch job
    ${job_id}                   Call API  job post  ${RES_DIR}/inferencejob-embedded.yaml
    ${exp_result}               create list  succeeded  failed
    ${result}                   Wait until command finishes and returns result  job  entity=${job_id}  exp_result=${exp_result}
    Status State Should Be      ${result}  succeeded
    ${result}                   check batch job response  ${RES_DIR}/inferencejob-embedded.yaml  ${RES_DIR}/output/response0.json
    Should Be True              ${result}
