*** Variables ***
${RES_DIR}              ${CURDIR}/resources
${LOCAL_CONFIG}         odahuflow/config_e2e_mlflow_wine
${WINE_ID}              test-e2e-wine

*** Settings ***
Documentation       Check wine model
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
Force Tags          e2e  wine  edi

*** Keywords ***
Cleanup resources
    StrictShell  odahuflowctl --verbose train delete --id ${WINE_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${WINE_ID} --ignore-not-found
    StrictShell  odahuflowctl --verbose dep delete --id ${WINE_ID} --ignore-not-found

*** Test Cases ***
Wine model
    StrictShell  odahuflowctl --verbose train create -f ${RES_DIR}/wine/training.odahuflow.yaml
    ${res}=  StrictShell  odahuflowctl train get --id ${WINE_ID} -o 'jsonpath=$[0].status.artifacts[0].artifactName'

    StrictShell  odahuflowctl --verbose pack create -f ${RES_DIR}/wine/packaging.odahuflow.yaml --artifact-name ${res.stdout}
    ${res}=  StrictShell  odahuflowctl pack get --id ${WINE_ID} -o 'jsonpath=$[0].status.results[0].value'

    StrictShell  odahuflowctl --verbose dep create -f ${RES_DIR}/wine/deployment.odahuflow.yaml --image ${res.stdout}

    Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl model info --mr ${WINE_ID}

    Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl model invoke --mr ${WINE_ID} --json-file ${RES_DIR}/wine/request.json

    ${res}=  Shell  odahuflowctl model invoke --mr ${WINE_ID} --json-file ${RES_DIR}/wine/request.json --jwt wrong-token
    should not be equal  ${res.rc}  0

