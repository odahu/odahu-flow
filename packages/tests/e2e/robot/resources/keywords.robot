*** Settings ***
Documentation       OdahuFlow robot resources
Resource            variables.robot
Variables           ../load_variables_from_profiles.py
Library             String
Library             OperatingSystem
Library             Collections
Library             DateTime
Library             odahuflow.robot.libraries.k8s.K8s  ${ODAHUFLOW_NAMESPACE}
Library             odahuflow.robot.libraries.utils.Utils
Library             odahuflow.robot.libraries.process.Process
Library             odahuflow.robot.libraries.odahu_k8s_reporter.OdahuKubeReporter

*** Keywords ***
Shell
    [Arguments]           ${command}
    ${result}=            Run Process without PIPE   ${command}    shell=True
    Log                   stdout = ${result.stdout}
    Log                   stderr = ${result.stderr}
    [Return]              ${result}

StrictShell
    [Arguments]           ${command}
    ${res}=   Shell  ${command}
              Should Be Equal  ${res.rc}  ${0}
    [Return]  ${res}

FailedShell
    [Arguments]           ${command}
    ${res}=   Shell  ${command}
              Should Not Be Equal  ${res.rc}  ${0}
    [Return]  ${res}

Build model
    [Arguments]  ${mt_name}  ${model_name}  ${model_version}  ${model_image_key_name}=\${TEST_MODEL_IMAGE}  ${entrypoint}=simple.py  ${kwargs}=&{EMPTY}

    ${args}=  prepare stub args  ${kwargs}
    StrictShell  odahuflowctl --verbose mt delete ${mt_name} --ignore-not-found
    StrictShell  odahuflowctl --verbose mt create ${mt_name} --toolchain-type python --vcs ${TEST_VCS} --workdir odahuflow/tests/e2e/models --entrypoint=${entrypoint} -a '--name ${model_name} --version ${model_version} ${args}'

    ${mt}=  get model training  ${mt_name}

    Set Suite Variable  ${model_image_key_name}  ${mt.trained_image}

    # --------- DEPLOY COMMAND SECTION -----------

Run API deploy from model packaging
    [Arguments]  ${mp_name}  ${md_name}  ${res_file}  ${role_name}=${EMPTY}

    ${res}=  StrictShell  odahuflowctl pack get --id ${mp_name} -o 'jsonpath=$[0].status.results[0].value'
    StrictShell  odahuflowctl --verbose dep create --id ${md_name} -f ${res_file} --image ${res.stdout}
    report model deployment pods  ${md_name}

Run API apply from model packaging
    [Arguments]  ${mp_name}  ${md_name}  ${res_file}  ${role_name}=${EMPTY}

    ${res}=  StrictShell  odahuflowctl pack get --id ${mp_name} -o 'jsonpath=$[0].status.results[0].value'
    StrictShell  odahuflowctl --verbose dep edit --id ${md_name} -f ${res_file} --image ${res.stdout}

Run API deploy from model packaging and check model started
    [Arguments]  ${mp_name}  ${md_name}  ${res_file}  ${role_name}=${EMPTY}
    Run API deploy from model packaging  ${mp_name}  ${md_name}  ${res_file}  ${role_name}

    Check model started  ${md_name}

    # --------- UNDEPLOY COMMAND SECTION -----------
Run API undeploy model and check
    [Arguments]           ${md_name}
    ${edi_state}=                Shell  odahuflowctl --verbose md delete ${md_name} --ignore-not-found
    Should Be Equal As Integers  ${edi_state.rc}        0
    ${edi_state} =               Shell  odahuflowctl --verbose md get
    Should Be Equal As Integers  ${edi_state.rc}        0
    Should not contain           ${edi_state.stdout}    ${md_name}

# --------- OTHER KEYWORDS SECTION -----------
Check model started
    [Documentation]  check if model run in container by http request
    [Arguments]           ${md_name}
    ${resp}=              Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl --verbose model info --md ${md_name}
    Log                   ${resp.stdout}

Verify model info from api
    [Arguments]      ${target_model}       ${model_name}        ${model_state}      ${model_replicas}
    Should Be Equal  ${target_model[0]}    ${model_name}        invalid model name
    Should Be Equal  ${target_model[1]}    ${model_state}       invalid model state
    Should Be Equal  ${target_model[2]}    ${model_replicas}    invalid model replicas
    # --------- TEMPLATE KEYWORDS SECTION -----------

Check if component domain has been secured
    [Arguments]     ${component}    ${enclave}
    [Documentation]  Check that a odahuflow component is secured by auth
    &{response} =    Run Keyword If   '${enclave}' == '${EMPTY}'    Get component auth page    ${HOST_PROTOCOL}://${component}.${HOST_BASE_DOMAIN}
    ...    ELSE      Get component auth page    ${HOST_PROTOCOL}://${component}-${enclave}.${HOST_BASE_DOMAIN}
    Log              Auth page for ${component} is ${response}
    Dictionary Should Contain Item    ${response}    response_code    200
    ${auth_page} =   Get From Dictionary   ${response}    response_text
    Should contain   ${auth_page}    Log in

Secured component domain should not be accessible by invalid credentials
    [Arguments]     ${component}    ${enclave}
    [Documentation]  Check that a secured odahuflow component does not provide access by invalid credentials
    &{creds} =       Create Dictionary 	login=admin   password=admin
    &{response} =    Run Keyword If   '${enclave}' == '${EMPTY}'    Post credentials to auth    ${HOST_PROTOCOL}://${component}    ${HOST_BASE_DOMAIN}    ${creds}
    ...    ELSE      Post credentials to auth    ${HOST_PROTOCOL}://${component}-${enclave}     ${HOST_BASE_DOMAIN}    ${creds}
    Log              Bad auth page for ${component} is ${response}
    Dictionary Should Contain Item    ${response}    response_code    200
    ${auth_page} =   Get From Dictionary   ${response}    response_text
    Should contain   ${auth_page}    Log in to Your Account
    Should contain   ${auth_page}    Invalid Email Address and password

Secured component domain should be accessible by valid credentials
    [Arguments]     ${component}    ${enclave}
    [Documentation]  Check that a secured odahuflow component does not provide access by invalid credentials
    &{creds} =       Create Dictionary    login=${STATIC_USER_EMAIL}    password=${STATIC_USER_PASS}
    &{response} =    Run Keyword If   '${enclave}' == '${EMPTY}'    Post credentials to auth    ${HOST_PROTOCOL}://${component}    ${HOST_BASE_DOMAIN}    ${creds}
    ...    ELSE      Post credentials to auth    ${HOST_PROTOCOL}://${component}-${enclave}     ${HOST_BASE_DOMAIN}    ${creds}
    Log              Bad auth page for ${component} is ${response}
    Dictionary Should Contain Item    ${response}    response_code    200
    ${auth_page} =   Get From Dictionary   ${response}    response_text
    Should contain   ${auth_page}    ${component}
    Should not contain   ${auth_page}    Invalid Email Address and password

Login to the api and edge
    [Arguments]  ${service_account}=${SA_ADMIN}
    [Documentation]  Login into API using odahuflowctl. ${service_account} should be object with .client_id and .client_secret attributes
    ${res}=  Shell  odahuflowctl --verbose login --url ${API_URL} --client_id "${service_account.client_id}" --client_secret "${service_account.client_secret}" --issuer "${ISSUER}"
    Should be equal  ${res.rc}  ${0}
    ${res}=  Shell  odahuflowctl config set MODEL_HOST ${EDGE_URL}
    Should be equal  ${res.rc}  ${0}

Cleanup example resources
    [Arguments]  ${example_id}
    StrictShell  odahuflowctl --verbose train delete --id ${example_id} --ignore-not-found
    StrictShell  odahuflowctl --verbose pack delete --id ${example_id} --ignore-not-found
    StrictShell  odahuflowctl --verbose dep delete --id ${example_id} --ignore-not-found

Run example model
    [Arguments]  ${example_id}  ${manifests_dir}
    StrictShell  odahuflowctl --verbose train create -f ${manifests_dir}/training.odahuflow.yaml --id ${example_id}
    report training pods  ${example_id}

    ${res}=  StrictShell  odahuflowctl train get --id ${example_id} -o 'jsonpath=$[0].status.artifacts[0].artifactName'

    StrictShell  odahuflowctl --verbose pack create -f ${manifests_dir}/packaging.odahuflow.yaml --artifact-name ${res.stdout} --id ${example_id}
    report packaging pods  ${example_id}
    ${res}=  StrictShell  odahuflowctl pack get --id ${example_id} -o 'jsonpath=$[0].status.results[0].value'

    StrictShell  odahuflowctl --verbose dep create -f ${manifests_dir}/deployment.odahuflow.yaml --image ${res.stdout} --id ${example_id}
    report model deployment pods  ${example_id}

    Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl model info --mr ${example_id}
    Wait Until Keyword Succeeds  1m  0 sec  StrictShell  odahuflowctl model invoke --mr ${example_id} --json-file ${manifests_dir}/request.json

    ${res}=  Shell  odahuflowctl model invoke --mr ${example_id} --json-file ${manifests_dir}/request.json --jwt wrong-token
    should not be equal  ${res.rc}  0