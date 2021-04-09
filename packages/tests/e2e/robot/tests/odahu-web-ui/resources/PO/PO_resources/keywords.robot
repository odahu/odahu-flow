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
    [Arguments]  ${locator}
    wait until page contains element  ${locator}
    mouse over  ${locator}
    element should be visible  ${locator}

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
    Validate visible element    ${element locator}
    Element Should Have Class   ${element locator}  ${class}

Mouse over and click button
    [Arguments]  ${locator}
    Validate visible element  ${locator}
    click button  ${locator}

Mouse over and click element
    [Arguments]  ${locator}
    Validate visible element  ${locator}
    click element  ${locator}

Select Item from the list
    # ${list_items_locator}: css, xpath locators for <li> tags of needed list
    # ${item_value}:         value of list item that should be selected
    [Arguments]  ${list_items_locator}  ${item_value}
    ${get_items}    get webelements  ${list_items_locator}
    FOR  ${element}  IN  @{get_items}
        ${text}  get text  ${element}
        log  ${text}
        Exit For Loop IF    "${text}" == "${item_value}"
    END
    click element  ${element}
