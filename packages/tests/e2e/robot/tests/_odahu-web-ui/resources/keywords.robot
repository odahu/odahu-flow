*** Keywords ***
Validate page
    [Arguments]  ${page_location}  ${page_title}  ${page_heading}
    wait until location is  ${page_location}
    title should be         ${page_title}
    wait until page contains  ${page_heading}

Validate visible element and text
    [Arguments]  ${element}  ${element_text}
    page should contain element  ${element}
    element should be visible  ${element}
    element text should be  ${element}  ${element_text}
