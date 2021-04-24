*** Variables ***
${RES_DIR}              ${CURDIR}/resources/wine-triton
${LOCAL_CONFIG}         odahuflow/e2e_mlflow_wine_triton
${EXAMPLE_ID}              wine-triton-min

*** Settings ***
Documentation       Check wine model
Test Timeout        120 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Variables           ../../load_variables_from_config.py
Resource            ../../resources/keywords.robot
Library             Collections
Library             OperatingSystem
Library             RequestsLibrary
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}
...                 AND  Login to the api and edge
...                 AND  Cleanup example resources  ${EXAMPLE_ID}
Suite Teardown      Run Keywords
...                 Pass Execution If  '${NO_CLEAN_UP}' != '${False}'  Suite Teardown is not executed
...                 AND  Cleanup example resources  ${EXAMPLE_ID}
...                 AND  Remove file  ${LOCAL_CONFIG}
Force Tags          e2e  wine  cli  triton

*** Test Cases ***
Wine model
    Download file  mlflow/sklearn-triton/wine/odahuflow/training.odahuflow.yaml  ${RES_DIR}/training.odahuflow.yaml
    Download file  mlflow/sklearn-triton/wine/odahuflow/packaging.odahuflow.yaml  ${RES_DIR}/packaging.odahuflow.yaml
    Download file  mlflow/sklearn-triton/wine/odahuflow/deployment.odahuflow.yaml  ${RES_DIR}/deployment.odahuflow.yaml
    Download file  mlflow/sklearn-triton/wine/odahuflow/request.json  ${RES_DIR}/request.json

    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/training.odahuflow.yaml --id ${EXAMPLE_ID}
    ${res}=  StrictShell  odahuflowctl train get --id ${EXAMPLE_ID} -o 'jsonpath=$[0].status.artifacts[0].artifactName'
    ${model_name}=  StrictShell  odahuflowctl train get --id ${EXAMPLE_ID} -o 'jsonpath=$[0].spec.model.name'

    Log  ${model_name}

    StrictShell  odahuflowctl --verbose pack create -f ${RES_DIR}/packaging.odahuflow.yaml --artifact-name ${res.stdout} --id ${EXAMPLE_ID}
    ${res}=  StrictShell  odahuflowctl pack get --id ${EXAMPLE_ID} -o 'jsonpath=$[0].status.results[0].value'

    StrictShell  odahuflowctl --verbose dep create -f ${RES_DIR}/deployment.odahuflow.yaml --image ${res.stdout} --id ${EXAMPLE_ID}
    report model deployment pods  ${EXAMPLE_ID}

    ${deployment_path}=  Catenate  SEPARATOR=/  model  ${EXAMPLE_ID}
    ${model_path}=  catenate  SEPARATOR=/  ${deployment_path}  v2  models  wine
    ${readiness_path}=  Catenate  SEPARATOR=/  ${model_path}  ready
    ${infer_path}=  Catenate  SEPARATOR=/  ${model_path}  infer

    ${headers}=    Create Dictionary    Authorization=Bearer ${CONFIG}[API_TOKEN]
    Create Session  odahu  headers=${headers}  url=${CONFIG}[API_URL]

    Wait until keyword succeeds  30 min  10 sec  GET On Session  odahu  url=${model_path}
    Wait until keyword succeeds  1 min  10 sec  GET On Session  odahu  url=${readiness_path}

    ${body}=  Get file  ${RES_DIR}/request.json
    ${response}=  Post on session  odahu  url=${infer_path}  data=${body}
    Status Should Be  200  ${response}
