*** Variables ***
${RES_DIR}             ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_invoke
${MD_SIMPLE_MODEL}     simple-model-invoke

*** Settings ***
Documentation       OdahuFlow's API operational check for operations on ModelDeployment resources
Test Timeout        20 minutes
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Force Tags          cli  deployment  invoke
Suite Setup         Run Keywords  Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                               Login to the api and edge  AND
...                               Cleanup resources  AND
...                               Run API deploy from model packaging  ${MP_SIMPLE_MODEL}  ${MD_SIMPLE_MODEL}  ${RES_DIR}/simple-model.deployment.odahuflow.yaml  AND
...                               Check model started  ${MD_SIMPLE_MODEL}
Suite Teardown      Run keywords  Cleanup resources  AND
...                 Remove File  ${LOCAL_CONFIG}

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose dep delete --id ${MD_SIMPLE_MODEL} --ignore-not-found

Refresh security tokens
    [Documentation]  Refresh api and model tokens

    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"
             Should be equal  ${res.rc}  ${0}
    ${res}=  Shell  odahuflowctl config set MODEL_HOST ${EDGE_URL}
             Should be equal  ${res.rc}  ${0}

*** Test Cases ***
Invoke. No base url
    [Documentation]  Fails if base url is not specified
    [Teardown]  Login to the api and edge
    # Ensure that next command will not use the config file
    Remove File  ${LOCAL_CONFIG}
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  Base url is required

Invoke. Empty model service url
    [Documentation]  Fails if model service url is empty
    [Teardown]  Login to the api and edge
    [Setup]     Remove File  ${LOCAL_CONFIG}
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${EDGE_URL} --url-prefix /model/${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --jwt "some_token"
             Should not be equal  ${res.rc}      ${0}
             Should contain       ${res.stderr}  401

Invoke. Wrong jwt
    [Documentation]  Fails if jwt is wrong
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --base-url ${EDGE_URL} --jwt wrong
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  401

Invoke. Pass parameters explicitly
    [Documentation]  Pass parameters explicitly
    ${JWT}=  Refresh security tokens
    # Ensure that next command will not use the config file
    Remove File  ${LOCAL_CONFIG}
    StrictShell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"

    ${res}=  Shell  odahuflowctl --verbose model invoke --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json --base-url ${EDGE_URL} --jwt "${AUTH_TOKEN}"
             Should be equal  ${res.rc}  ${0}
             Should contain   ${res.stdout}  42

Invoke. Pass parameters through config file
    [Documentation]  Pass parameters through config file
    Refresh security tokens
    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${API_URL} --md ${MD_SIMPLE_MODEL} --json-file ${RES_DIR}/simple-model.request.json
             Should be equal  ${res.rc}  ${0}
             Should contain   ${res.stdout}  42

Invoke. Pass model parameters using json
    [Documentation]  Model parameters as json
    Refresh security tokens
    ${res}=  Shell  odahuflowctl --verbose model invoke --base-url ${API_URL} --md ${MD_SIMPLE_MODEL} --json '{"columns": ["a","b"],"data": [[1.0,2.0]]}'
             Should be equal  ${res.rc}  ${0}
             Should contain   ${res.stdout}  42
