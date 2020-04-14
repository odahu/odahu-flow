*** Settings ***
Documentation  Checking of airflow settings
Resource            ../../resources/keywords.robot
Force Tags  airflow

*** Variables ***
${TEST_DAG_RUN_IDS}  health_check,airflow-wine-from-yamls

*** Test Cases ***
Airflow DAG
    [Documentation]  Checking Ariflow health status by DAG
    StrictShell  ${CURDIR}/resources/tools_test_airflow.sh --dags ${TEST_DAG_RUN_IDS}