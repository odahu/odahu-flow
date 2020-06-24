*** Settings ***
Documentation       API keywords
Library             Collections

*** Keywords ***
Call API
    [Arguments]                  ${keyword}  @{arguments}
    ${result}                    Run Keyword  ${keyword}  @{arguments}
    Log                          ${result}
    [Return]                     ${result}

Log id
    [Arguments]                  ${input}
    ${output}                    set variable  ${input.id}
    Log  ${output}
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

Status State Should Be
    [Arguments]                  ${result}  ${exp_state}
    Should Be Equal              ${result.status.state}  ${exp_state}

Status Reason Should Contain
    [Arguments]                  ${result}  ${exp_reason}
    should contain               ${result.status.reason}  ${exp_reason}

keySecret connection should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.key_secret}  ${exp_value}

keySecret connection should not be equal
    [Arguments]                  ${result}  ${exp_value}
    should not be equal          ${result.spec.key_secret}  ${exp_value}

Password connection should be equal
    [Arguments]                  ${result}  ${exp_value}
    should be equal              ${result.spec.password}  ${exp_value}

Password connection should not be equal
    [Arguments]                  ${result}  ${exp_value}
    should not be equal          ${result.spec.password}  ${exp_value}

Wait until command finishes and returns result
    [Arguments]    ${command}  ${cycles}=20  ${sleep_time}=30s  ${result}=  @{exp_result}=succeeded
    FOR     ${i}    IN RANGE   ${cycles}
        ${result}              Call API  ${command} get id  ${result.id}
        ${list_contain}        get match count  ${exp_result}  ${result.status.state}
        log   ${list_contain}
        exit for loop if       '${list_contain}' != '0'
        Sleep                  ${sleep_time}
    END
    [Return]  ${result}

Command response list should contain id
    [Arguments]                         ${command}  @{value}
    ${result}                           Call API  ${command} get
    ${list_length}                      get length  ${result}

    FOR     ${i}   IN  @{result}
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

Get Value From Dictionary
    [Arguments]             ${Dictionary Name}      ${Key}
    ${KeyIsPresent}=        Run Keyword And Return Status       Dictionary Should Contain Key       ${Dictionary Name}      ${Key}
    ${Value}=               Run Keyword If      ${KeyIsPresent}     Get From Dictionary             ${Dictionary Name}      ${Key}
    [Return]                ${Value}
