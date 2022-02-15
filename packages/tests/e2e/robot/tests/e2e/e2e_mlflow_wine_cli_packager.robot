*** Variables ***
${RES_DIR}              ${CURDIR}/resources/wine
${LOCAL_CONFIG}         odahuflow/config_e2e_mlflow_wine_cli_packager
${WINE_ID}              test-e2e-wine-cli-packager
${PREDICT_FILE_NAME}    answer.txt

*** Settings ***
Documentation       Check wine model
Test Timeout        120 minutes
Variables           ../../load_variables_from_profiles.py    ${CLUSTER_PROFILE}
Resource            ../../resources/keywords.robot
Resource            ../../resources/variables.robot
Library             OperatingSystem
Library             Collections
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.model.Model
Library             odahuflow.robot.libraries.examples_loader.ExamplesLoader  https://raw.githubusercontent.com/odahu/odahu-examples  ${EXAMPLES_VERSION}
Suite Setup         E2E test setup  ${WINE_ID}
Suite Teardown      E2E test teardown  ${WINE_ID}
Force Tags          e2e  wine  cli  cli-packager

*** Keywords ***
Run example model
    [Arguments]  ${example_id}  ${manifests_dir}
    StrictShell  odahuflowctl --verbose train create -f ${manifests_dir}/training.odahuflow.yaml --id ${example_id}
    report training pods  ${example_id}

    ${artifact_name}=  StrictShell  odahuflowctl train get --id ${example_id} -o 'jsonpath=$[0].status.artifacts[0].artifactName'

    StrictShell  odahuflowctl --verbose pack create -f ${manifests_dir}/packaging.cli.odahuflow.yaml --artifact-name ${artifact_name.stdout} --id ${example_id}
    report packaging pods  ${example_id}
    ${packaged_image}=  StrictShell  odahuflowctl pack get --id ${example_id} -o 'jsonpath=$[0].status.results[0].value'

    # download docker image and make prediction request
    ${docker_command}=  catenate
    ...     docker run --rm --net host
    ...     --mount type=bind,source=$(pwd)/packages/tests/e2e/robot/tests/e2e/resources/wine,target=/volume
    ...     ${packaged_image.stdout} predict --output_file_name ${PREDICT_FILE_NAME} /volume/request.json /volume/
    ${prediction_request}=  Strictshell   ${docker_command}
    should be equal as strings   Prediction is successful. Result file: /volume/${PREDICT_FILE_NAME}  ${prediction_request.stdout}
    ${prediction_file_content}=  get file  packages/tests/e2e/robot/tests/e2e/resources/wine/${PREDICT_FILE_NAME}
    should be equal as strings  ${WINE_MODEL_RESULT}   ${prediction_file_content}


*** Test Cases ***
Wine model for docker-cli packager
    Download file  mlflow/sklearn/wine/odahuflow/training.odahuflow.yaml  ${RES_DIR}/training.odahuflow.yaml
    Download file  mlflow/sklearn/wine/odahuflow/packaging.cli.odahuflow.yaml  ${RES_DIR}/packaging.cli.odahuflow.yaml
    Download file  mlflow/sklearn/wine/odahuflow/deployment.odahuflow.yaml  ${RES_DIR}/deployment.odahuflow.yaml
    Download file  mlflow/sklearn/wine/odahuflow/request.json  ${RES_DIR}/request.json

    Run example model  ${WINE_ID}  ${RES_DIR}
