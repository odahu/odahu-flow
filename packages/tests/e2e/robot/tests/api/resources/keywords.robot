*** Settings ***
Documentation       API keywords
Library             Collections

*** Keywords ***
Call API
    [Arguments]                  ${keyword}  @{arguments}
    ${result}=                   Run Keyword  ${keyword}  @{arguments}
    Log                          ${result}
    [Return]                     ${result}

Log id
    [Arguments]                  ${input}
    ${output}=                   Log  ${input.id}
    [Return]                     ${output}

Log Status
    [Arguments]                  ${input}
    ${output}=                   Log  ${input.status}
    [Return]                     ${output}

Log State
    [Arguments]                  ${input}
    Log Status                   ${input.state}

Status State Should Be
    [Arguments]                  ${result}  ${exp_state}
    Should Be Equal              ${result.status.state}  ${exp_state}

Status Reason Should Contain
    [Arguments]                  ${result}  ${exp_reason}
    should contain               ${result.status.reason}  ${exp_reason}

Wait until command finishes and returns result
    [Arguments]    ${command}  ${cycles}=20  ${sleep_time}=30s  ${result}=  ${exp_result}='succeeded'
    FOR     ${i}    IN RANGE   ${cycles}
        ${result}=             Call API  ${command} get id  ${result.id}
        exit for loop if       '${exp_result}' == '${result.status.state}'
        Sleep                  ${sleep_time}
    END
    [Return]  ${result}
