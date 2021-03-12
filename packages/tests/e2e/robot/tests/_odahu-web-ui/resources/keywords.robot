*** Keywords ***
Validate ODAHU Web UI page loaded
    [Arguments]  @{elements}
    FOR  ${element}  IN  @{elements}
        wait until page contains element  ${element}
        element should be visible  ${element}
    END

Validate page
    [Arguments]  ${page_location}  ${page_heading}
    wait until location is  ${page_location}
    wait until page contains  ${page_heading}

Validate visible element and text
    [Arguments]  ${element}  ${element_text}
    page should contain element  ${element}
    element should be visible  ${element}
    element text should be  ${element}  ${element_text}

Element should have class
    [Arguments]  ${locator}  ${expected class}
    ${class}=           Get Element Attribute  ${locator}  class
    should be equal     ${class}  ${expected class}

Element should not have class
    [Arguments]  ${locator}  ${expected class}
    ${class}=            Get Element Attribute  ${locator}  class
    should be not equal  ${class}  ${expected class}

Validate page contains the "Element" of the "Class"
    [Arguments]  ${element locator}  ${class}
    element should be visible  ${element locator}
    Element Should Have Class  ${element locator}  ${class}