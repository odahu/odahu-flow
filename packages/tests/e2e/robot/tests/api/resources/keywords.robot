*** Settings ***
Documentation       API keywords
Library             Collections

*** Variables ***
# Error Templates
${400 BadRequest Template}         WrongHttpStatusCode: Got error from server: {} (status: 400)
${401 Unathorized Template}
${403 Forbidden Template}
${404 NotFound Template}           WrongHttpStatusCode: Got error from server: entity "{}" is not found (status: 404)
${404 Model NotFoundTemplate}      Wrong status code returned: 404. Data: . URL: {}
${409 Conflict Template}           EntityAlreadyExists: Got error from server: entity "{}" already exists (status: 409)

${APIConnectionException}          APIConnectionException: Can not reach {base url}
${IncorrectToken}                  IncorrectAuthorizationToken: Refresh token is not correct.\nPlease login again

# Validation checks
${FailedConn}   Validation of connection is failed:

${invalid_id}   ID is not valid
# connections
${empty_uri}   the uri parameter is empty
${s3_empty_keyID_keySecret}  s3 type requires that keyID and keySecret parameters must be non-empty
${gcs_empty_keySecret}       gcs type requires that keySecret parameter must be non-empty
${azureblob_req_keySecret}   azureblob type requires that keySecret parameter containsHTTP endpoint with SAS Token

*** Keywords ***
Call API
    [Arguments]                  ${keyword}  @{arguments}  &{named arguments}
    ${result}                    run keyword  ${keyword}  @{arguments}  &{named arguments}
    Log                          ${result}
    [Return]                     ${result}

Call API and get Error
    [Arguments]                  ${expected_error}  ${keyword}  @{arguments}  &{named arguments}
    ${result}                    run keyword and expect error  ${expected_error}  ${keyword}  @{arguments}  &{named arguments}
    Log many                     ${result}
    [Return]                     ${result}

Call API and continue on Failure
    [Arguments]                  ${keyword}  @{arguments}  &{named arguments}
    ${result}                    run keyword and continue on failure  ${keyword}  @{arguments}  &{named arguments}
    Log                          ${result}
    [Return]                     ${result}

Get Logs
    [Arguments]                  ${entity type}  ${entity id}
    Call API                     ${entity type} get log  ${entity id}

Log Status
    [Arguments]                  ${input}
    ${output}                    set variable  ${input.status}
    Log  ${output}
    [Return]                     ${output}

Id should be equal
    [Arguments]                  ${result}  ${exp_id}
    should be equal              ${result.id}  ${exp_id}

Status State Should Be
    [Arguments]                  ${result}  ${exp_state}
    should be equal              ${result.status.state}  ${exp_state}

Status Reason Should Contain
    [Arguments]                  ${result}  ${exp_reason}
    should contain               ${result.status.reason}  ${exp_reason}

KeySecret connection should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.key_secret}  ${exp_value}

KeySecret connection should not be equal
    [Arguments]                  ${result}  ${exp_value}
    should not be equal          ${result.spec.key_secret}  ${exp_value}

Password connection should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.password}  ${exp_value}

Password connection should not be equal
    [Arguments]                  ${result}  ${exp_value}
    should not be equal          ${result.spec.password}  ${exp_value}

Default Docker image should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.default_image}  ${exp_value}

Default Entrypoint should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.entrypoint}  ${exp_value}

Requested resources should be equal
    [Arguments]                  ${result}  ${CPU}  ${GPU}  ${MEMORY}
    should be equal              ${result.spec.resources.requests.cpu}        ${CPU}
    should be equal              ${result.spec.resources.requests.gpu}        ${GPU}
    should be equal              ${result.spec.resources.requests.memory}     ${MEMORY}

Limits resources should be equal
    [Arguments]                 ${result}  ${CPU}  ${GPU}  ${MEMORY}
    should be equal              ${result.spec.resources.limits.cpu}        ${CPU}
    should be equal              ${result.spec.resources.limits.gpu}        ${GPU}
    should be equal              ${result.spec.resources.limits.memory}     ${MEMORY}

CreatedAt and UpdatedAt times should not be equal
    [Arguments]                  ${result}
    ${result_status}             Log Status  ${result}
    should not be equal          ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Wait until command finishes and returns result
    [Arguments]    ${command}    ${cycles}=120  ${sleep_time}=30s  ${entity}=  @{exp_result}=succeeded
    FOR     ${i}    IN RANGE     ${cycles}
        ${result}                Call API  ${command} get id  ${entity}
        ${result_state}          evaluate  str('${result.status.state}' or '')
        ${list_contain}          count values in list  ${exp_result}  ${result_state}
        exit for loop if         '${list_contain}' != '0'
        Sleep                    ${sleep_time}
    END
    [Return]  ${result}

Wait until delete finished
    [Arguments]    ${command}    ${cycles}=60  ${sleep_time}=30s  ${entity}=  @{exp_result}=deleting
    FOR     ${i}    IN RANGE     ${cycles}
        ${check}                 Check command response list contains id  ${command}  ${entity}

        exit for loop if         ${check} == ${FALSE}
        sleep                    ${sleep_time}
    END

Check command response list contains id
    [Arguments]                         ${command}  @{value}
    ${result}                           Call API  ${command} get
    ${list_length}                      get length  ${result}

    FOR     ${i}  IN  @{result}
        exit for loop if                $i.id in $value
        ${list_length}                  evaluate  ${list_length} - 1
    END

    ${result}                           set variable if  '${list_length}' != '0'  ${TRUE}  ${FALSE}
    [Return]                            ${result}

Command response list should contain id
    [Arguments]                         ${command}  @{value}
    ${response list}                    Call API  ${command} get
    ${value_length}                     get length  ${value}

    FOR     ${i}  IN  @{response list}
        continue for loop if            $i.id not in $value
        ${value_length}                 evaluate  ${value_length} - 1
    END
    should be equal as integers         ${value_length}  0

Command response list should not contain id
    [Arguments]                         ${command}  @{value}
    ${response list}                    Call API  ${command} get
    ${value_length}                     get length  ${value}

    FOR     ${i}  IN  @{response list}
        ${value_length}                 set variable if  $i.id in $value  0  ${value_length}
        exit for loop if                $i.id in $value
    END
    should not be equal as integers     ${value_length}  0

Pick artifact name
    [Arguments]                   ${training id}
    ${result}                     Call API  training get id  ${training id}
    ${artifact}                   get variable value  ${result.status.artifacts[0]}  ${EMPTY}
    ${artifact_name}              get variable value  ${artifact.artifact_name}  ${EMPTY}
    [Return]                      ${artifact_name}

Pick packaging image
    [Arguments]                   ${packaging id}
    ${result}                     Call API  packaging get id  ${packaging id}
    ${image}                      get variable value  ${result.status.results[0]}  ${EMPTY}
    ${image_value}                get variable value  ${image.value}  ${EMPTY}
    [Return]                      ${image_value}
