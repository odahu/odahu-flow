*** Variables ***
${RES_DIR}             ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_invoke
${MD_SIMPLE_MODEL}     simple-model-invoke
${MD_SIMPLE_MODEL_1}   simple-model-multiver-1
${MD_SIMPLE_MODEL_2}   simple-model-multiver-2

*** Settings ***
Documentation       Tests invoke (predict) process through CLI for deployed ODAHU models
Test Timeout        20 minutes
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Force Tags          cli  deployment  invoke  testing
Suite Setup         Run Keywords  Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                               Login to the api and edge  AND
...                               Cleanup resources  AND
...                               Run API deploy from model packaging  ${MP_SIMPLE_MODEL}  ${MD_SIMPLE_MODEL}  ${RES_DIR}/simple-model.deployment.odahuflow.yaml  AND
...                               Check model started  ${MD_SIMPLE_MODEL}
Suite Teardown      Run keywords  Cleanup resources  AND
...                 Remove File  ${LOCAL_CONFIG}

*** Keywords ***
Validate invoke succeeded and result
    [Arguments]      ${res}
    Should be equal  ${res.rc}  ${0}
    Should contain   ${res.stdout}  42

Cleanup resources
    StrictShell  odahuflowctl --verbose dep delete --id ${MD_SIMPLE_MODEL} --ignore-not-found
    StrictShell  odahuflowctl --verbose dep delete --id ${MD_SIMPLE_MODEL_1} --ignore-not-found
    StrictShell  odahuflowctl --verbose dep delete --id ${MD_SIMPLE_MODEL_2} --ignore-not-found

Refresh security tokens
    [Documentation]  Refresh api and model tokens

    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"
             Should be equal  ${res.rc}  ${0}
    ${res}=  Shell  odahuflowctl config set MODEL_HOST ${EDGE_URL}
             Should be equal  ${res.rc}  ${0}

*** Test Cases ***
Invoke. Empty token
    [Documentation]  Trying to invoke model with empty token. Suceeds if config exists, fails if not
    [Teardown]  Login to the api and edge
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --token ""
             Validate invoke succeeded and result  ${res}

    Remove File  ${LOCAL_CONFIG}

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --token ""
             Should not be equal  ${res.rc}  ${0}

Invoke. Empty model service url
    [Documentation]  Fails if model service url is empty
    [Teardown]  Login to the api and edge
    [Setup]     Remove File  ${LOCAL_CONFIG}
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url "" --md "${MD_SIMPLE_MODEL}" --json-file ${RES_DIR}/simple-model.request.json
             Should not be equal  ${res.rc}      ${0}
             Should contain       ${res.stderr}  401

Invoke. Wrong token
    [Documentation]  Fails if token is wrong
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --base-url ${EDGE_URL} --token wrong
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  Credentials are not correct

Invoke. Pass parameters explicitly
    [Documentation]  Pass parameters explicitly
    ${JWT}=  Refresh security tokens
    # Ensure that next command will not use the config file
    Remove File  ${LOCAL_CONFIG}
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --base-url ${EDGE_URL} --token "${AUTH_TOKEN}"
             Validate invoke succeeded and result  ${res}

Invoke. Pass request body through file
    [Documentation]  Pass parameters through config file
    Refresh security tokens
    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json
             Validate invoke succeeded and result  ${res}

Invoke. Pass request body through command line
    [Documentation]  Model parameters as json
    Refresh security tokens
    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json '{"columns": ["a","b"],"data": [[1.0,2.0]]}'
             Validate invoke succeeded and result  ${res}

Invoke. Deploy 2 models with the same image and invoke
    Run API deploy from model packaging  ${MP_SIMPLE_MODEL}  ${MD_SIMPLE_MODEL_1}  ${RES_DIR}/simple-model-1.deployment.odahuflow.yaml
    Check model started  ${MD_SIMPLE_MODEL_1}

    Run API deploy from model packaging  ${MP_SIMPLE_MODEL}  ${MD_SIMPLE_MODEL_2}  ${RES_DIR}/simple-model-2.deployment.odahuflow.yaml
    Check model started  ${MD_SIMPLE_MODEL_1}
    Check model started  ${MD_SIMPLE_MODEL_2}

    ${resp}=        StrictShell  odahuflowctl --verbose dep get
                    Should contain              ${resp.stdout}      ${MD_SIMPLE_MODEL_1}
                    Should contain              ${resp.stdout}      ${MD_SIMPLE_MODEL_2}

    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${EDGE_URL} --md ${MD_SIMPLE_MODEL_1} --json-file ${RES_DIR}/simple-model.request.json
         Validate invoke succeeded and result  ${res}
    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${EDGE_URL} --md ${MD_SIMPLE_MODEL_2} --json-file ${RES_DIR}/simple-model.request.json
         Validate invoke succeeded and result  ${res}

    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${EDGE_URL} --url-prefix /model/${MD_SIMPLE_MODEL_1} --json-file ${RES_DIR}/simple-model.request.json --token ${AUTH_TOKEN}
         Validate invoke succeeded and result  ${res}
    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${EDGE_URL} --url-prefix /model/${MD_SIMPLE_MODEL_2} --json-file ${RES_DIR}/simple-model.request.json --token ${AUTH_TOKEN}
         Validate invoke succeeded and result  ${res}
