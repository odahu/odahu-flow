*** Settings ***
Documentation  Checking of airflow settings
Test Timeout        60 minutes
Resource            ../../resources/keywords.robot
Library             odahuflow.robot.libraries.odahu_k8s_reporter.OdahuKubeReporter
Force Tags  airflow

*** Variables ***
${TEST_DAG_RUN_IDS}  health_check,airflow-wine-from-yamls

*** Test Cases ***
Airflow DAG
    [Documentation]  Checking Ariflow health status by DAG
    StrictShell  ${CURDIR}/resources/tools_test_airflow.sh --dags ${TEST_DAG_RUN_IDS}
    report training pods  airflow-wine-from-yamls
    report packaging pods  airflow-wine-from-yamls
    report model deployment pods  airflow-wine-from-yamls