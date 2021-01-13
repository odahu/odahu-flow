*** Variables ***
${RES_DIR}              ${CURDIR}/resources/wine-triton
${LOCAL_CONFIG}         odahuflow/e2e_mlflow_wine_triton
${EXAMPLE_ID}              test-e2e-wine-triton

*** Settings ***
Documentation       Check wine model
Test Timeout        120 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         Run Keywords
...                 Set Environment Variable  ODAHUFLOW_CONFIG  ${LOCAL_CONFIG}  AND
...                 Login to the api and edge  AND
...                 Cleanup example resources  ${WINE_ID}
Suite Teardown      Run Keywords
...                 Cleanup example resources  ${WINE_ID}  AND
...                 Remove file  ${LOCAL_CONFIG}
Force Tags          e2e  wine  cli  triton

*** Test Cases ***
Wine model
    Download file  mlflow/sklearn-triton/wine/odahuflow/training.odahuflow.yaml  ${RES_DIR}/training.odahuflow.yaml
    Download file  mlflow/sklearn-triton/wine/odahuflow/packaging.odahuflow.yaml  ${RES_DIR}/packaging.odahuflow.yaml

    StrictShell  odahuflowctl --verbose train create -f ${manifests_dir}/training.odahuflow.yaml --id ${EXAMPLE_ID}

    ${res}=  StrictShell  odahuflowctl train get --id ${EXAMPLE_ID}  -o 'jsonpath=$[0].status.artifacts[0].artifactName'

    StrictShell  odahuflowctl --verbose pack create -f ${manifests_dir}/packaging.odahuflow.yaml --artifact-name ${res.stdout} --id ${EXAMPLE_ID}
    ${res}=  StrictShell  odahuflowctl pack get --id ${EXAMPLE_ID}  -o 'jsonpath=$[0].status.results[0].value'
