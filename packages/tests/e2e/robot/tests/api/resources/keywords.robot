*** Settings ***
Documentation       API keywords
Library             Collections

*** Keywords ***
Call API
    [Arguments]                  ${keyword}  @{arguments}  &{named arguments}
    ${result}                    Run Keyword  ${keyword}  @{arguments}  &{named arguments}
    Log                          ${result}
    [Return]                     ${result}

Log id
    [Arguments]                  ${input}
    ${output}                    set variable  ${input.id}
    Log                          ${output}
    [Return]                     ${output}

Log Status
    [Arguments]                  ${input}
    ${output}                    set variable  ${input.status}
    Log  ${output}
    [Return]                     ${output}

Log State
    [Arguments]                  ${input}
    ${output}                    set variable  ${input.state}
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
    should be equal              ${result.spec.default_image} ${exp_value}

Default Entrypoint should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.entrypoint} ${exp_value}

CreatedAt and UpdatedAt times should not be equal
    [Arguments]                  ${result}
    ${result_status}             Log Status  ${result}
    should not be equal          ${result_status}.get('createdAt')  ${result_status}.get('updatedAt')

Wait until command finishes and returns result
    [Arguments]    ${command}    ${cycles}=20  ${sleep_time}=30s  ${result}=  @{exp_result}=succeeded
    FOR     ${i}    IN RANGE     ${cycles}
        ${result}                Call API  ${command} get id  ${result.id}
        ${list_contain}          get match count  ${exp_result}  ${result.status.state}
        log   ${list_contain}
        exit for loop if         '${list_contain}' != '0'
        Sleep                    ${sleep_time}
    END
    [Return]  ${result}

Command response list should contain id
    [Arguments]                         ${command}  @{value}
    ${result}                           Call API  ${command} get
    ${list_length}                      get length  ${result}

    FOR     ${i}  IN  @{result}
        exit for loop if                $i.id in $value
        ${list_length}                  evaluate  ${list_length} - 1
    END
    should not be equal as integers     ${list_length}  0

Command response list should not contain id
    [Arguments]                         ${command}  @{value}
    ${result}                           Call API  ${command} get
    ${list_length}                      get length  ${result}

    FOR     ${i}  IN  @{result}
        exit for loop if                $i.id in $value
        ${list_length}                  evaluate  ${list_length} - 1
    END
    should be equal as integers         ${list_length}  0

Pick artifact name
    [Arguments]                   ${training id}
    ${result}                     Call API  training get id  ${training id}
    ${artifact}                   set variable  ${result.status.artifacts[0]}
    ${artifact_name}              set variable  ${artifact.artifact_name}
    [Return]                      ${artifact_name}

Pick packaging image
    [Arguments]                   ${packaging id}
    ${result}                     Call API  packaging get id  ${packaging id}
    ${image}                      set variable  ${result.status.results[0]}
    ${image_value}                set variable  ${image.value}
    [Return]                      ${image_value}
