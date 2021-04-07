*** Settings ***
Documentation   keywords for Page Objects
Library         SeleniumLibrary  timeout=10s

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

Validate visible element
    [Arguments]  ${element}
    wait until page contains element  ${element}
    element should be visible  ${element}

Validate visible element and text
    [Arguments]  ${element}  ${element_text}
    Validate visible element  ${element}
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
    wait until page contains element  ${element locator}
    element should be visible  ${element locator}
    Element Should Have Class  ${element locator}  ${class}
