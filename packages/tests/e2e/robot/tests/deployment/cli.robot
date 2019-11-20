*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}        odahuflow/config_deployment_cli
${MD_SIMPLE_MODEL}     simple-model-cli

*** Settings ***
Documentation       OdahuFlow's API operational check for operations on ModelDeployment resources
Test Timeout        20 minutes
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Library             odahuflow.robot.libraries.utils.Utils
Library             Collections
Force Tags          deployment  api  cli
Suite Setup         Run Keywords  Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                               Login to the api and edge  AND
...                               Cleanup resources
Suite Teardown      Run keywords  Cleanup resources  AND
...                 Remove File  ${LOCAL_CONFIG}

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose dep delete --id ${MD_SIMPLE_MODEL} --ignore-not-found

*** Test Cases ***
Undeploy. Nonexistent model service
    [Documentation]  The undeploy command must fail if a model cannot be found by name
    ${res}=  Shell  odahuflowctl --verbose dep delete --id this-model-does-not-exsit
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  not found

Deploy. Zero timeout parameter
    [Documentation]  The deploy command must fail if timeout parameter is zero
    ${res}=  Shell  odahuflowctl --verbose dep create -f ${RES_DIR}/custom-resources.deployment.odahuflow.yaml --timeout=0
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  must be positive integer

Deploy. Negative timeout parameter
    [Documentation]  The deploy command must fail if it contains negative timeout parameter
    ${res}=  Shell  odahuflowctl --verbose dep create -f ${RES_DIR}/custom-resources.deployment.odahuflow.yaml --timeout=-500
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  must be positive integer

Missed the host parameter
    [Documentation]  The inspect command must fail if it does not contain an api host
    [Teardown]  Login to the api and edge
    [Setup]     Remove File  ${LOCAL_CONFIG}
    ${res}=  Shell  odahuflowctl --verbose dep --token "${AUTH_TOKEN}" get
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  Can not reach http://localhost:5000

Wrong token
    [Documentation]  The inspect command must fail if it does not contain a token
    [Teardown]  Login to the api and edge
    [Setup]     Remove File  ${LOCAL_CONFIG}
    ${res}=  Shell  odahuflowctl --verbose dep --url ${API_URL} --token wrong-token get
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  Credentials are not correct

Login. Override login values
    [Documentation]  Command line parameters must be overrided by config parameters
    [Teardown]  Login to the api and edge
    [Setup]     Remove File  ${LOCAL_CONFIG}
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --token "${AUTH_TOKEN}"
             Should be equal  ${res.rc}  ${0}

    ${res}=  Shell  odahuflowctl --verbose dep --url ${API_URL} --token wrong-token get
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  Credentials are not correct

Deploy fails when validation fails
    [Documentation]  Deploy fails when memory resource is incorect
    ${res}=  Shell  odahuflowctl --verbose dep create -f ${RES_DIR}/validation-fail.deployment.odahuflow.yaml
             Should not be equal  ${res.rc}  ${0}
             Should contain       ${res.stderr}  minimum number of replicas parameter must not be less than maximum number of replicas parameter
