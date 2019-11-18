*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}         odahuflow/config_e2e_mlflow_tensorflow
${TENSORFLOW_ID}        test-e2e-tensorflow

*** Settings ***
Documentation       Check tensorflow model
Test Timeout        60 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the edi and edge  AND
...                 Cleanup resources
Suite Teardown      Run Keywords
...                 Cleanup resources  AND
...                 Remove file  ${LOCAL_CONFIG}
Force Tags          e2e  tensorflow  edi

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose train delete --id ${TENSORFLOW_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${TENSORFLOW_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose dep delete --id ${TENSORFLOW_ID} --ignore-not-found

*** Test Cases ***
Tensorflow model
    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/tensorflow/training.odahuflow.yaml
    ${res}=  StrictShell  odahuflowctl train get --id ${TENSORFLOW_ID} -o 'jsonpath=$[0].status.artifacts[0].artifactName'

    StrictShell  odahuflowctl --verbose pack create -f ${RES_DIR}/tensorflow/packaging.odahuflow.yaml --artifact-name ${res.stdout}
    ${res}=  StrictShell  odahuflowctl pack get --id ${TENSORFLOW_ID} -o 'jsonpath=$[0].status.results[0].value'

    StrictShell  odahuflowctl --verbose dep create -f ${RES_DIR}/tensorflow/deployment.odahuflow.yaml --image ${res.stdout}

    Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl --verbose model info --mr ${TENSORFLOW_ID}

    Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl --verbose model invoke --mr ${TENSORFLOW_ID} --json-file ${RES_DIR}/tensorflow/request.json

    ${res}=  Shell  odahuflowctl model invoke --mr ${TENSORFLOW_ID} --json-file ${RES_DIR}/tensorflow/request.json --jwt wrong-token
    should not be equal  ${res.rc}  0
