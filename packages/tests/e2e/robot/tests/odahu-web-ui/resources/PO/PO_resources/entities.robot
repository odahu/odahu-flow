*** Settings ***
Documentation   common Page Objects for different entities
Library         SeleniumLibrary  timeout=10s
Library         Collections

*** Variables ***
${ENTITIES.PAGE_HEADING}        xpath://*[@id="root"]/div/main/div[2]/div/div/div[1]/h6
# ${ENTITIES.PAGE_HEADING.CLASS}  MuiTypography-root jss165 MuiTypography-h6
${ENTITIES.NEW_BUTTON}          xpath://*[@id="root"]/div/main/div[2]/div/div/div[1]/button[1]
${ENTITIES.REFRESH_BUTTON}      xpath://*[@id="root"]/div/main/div[2]/div/div/div[1]/button[2]

# table
${ENTITIES.TABLE}        xpath://*[@id="root"]/div/main/div[2]/div/div/div[2]/table
${ENTITIES.TABLE.HEAD}   xpath://*[@id="root"]/div/main/div[2]/div/div/div[2]/table/thead
${ENTITIES.TABLE.BODY}   xpath://*[@id="root"]/div/main/div[2]/div/div/div[2]/table/tbody

${ENTITIES.TABLE.ROWS_PER_PAGE}         css:.MuiTablePagination-selectRoot
${ENTITIES.TABLE.ROWS_PER_PAGE.LIST}    css:.MuiListItem-root[role="option"]
${ENTITIES.TABLE.FOOTER.PREVIOUS_PAGE}  css:[title="Previous page"]
${ENTITIES.TABLE.FOOTER.NEXT_PAGE}      css:[title="Next page"]

# tabs
${ENTITIES.TABS}

# ${ENTITIES.SUCCESS_MESSAGE}     Success

${LAST_COLUMN}      Updated at

*** Keywords ***
Click "+New" Button
    Keywords.Mouse over and click button  ${ENTITIES.NEW_BUTTON}

Validate page with entities
    [Arguments]  ${page_location}  ${page_heading}
    Keywords.Validate page  ${page_location}  ${page_heading}
    Keywords.Validate visible element and text  ${ENTITIES.PAGE_HEADING}  ${page_heading}

Validate Table Headers
    [Arguments]  @{headers}
    page should contain element  ${ENTITIES.TABLE.HEAD}
    FOR  ${header}  IN  @{headers}
        table header should contain  ${ENTITIES.TABLE.HEAD}  ${header}
    END

Iterate over row
    [Arguments]  @{headers}  ${row_element}
    ${col_amount} =  Get length  @{headers}
    @{row list} =  Create List
    &{row dict}  evaluate  {}
    FOR  ${col element}  IN RANGE  1  100
        ${cell} =  get table cell  ${ENTITIES.TABLE}  ${row_element}  ${col element}
        append to list  ${row List}  ${cell}
        set to dictionary  ${row dict}  ${headers}[${col element}]=${cell}
        log  ${cell}
    END
    log  ${row list}
    log  ${row dict}

Get Table Headers
    @{header list} =  Create List
    FOR  ${col element}  IN RANGE  1  100
        ${cell} =  get table cell  ${ENTITIES.TABLE}  1  ${col element}
        append to list  ${header list}  ${cell}
        EXIT For Loop If  '${cell}' == '${LAST_COLUMN}'
    END
    log  ${header list}

Validate that Entity in Table
    # 1st row, it's headers of table
    [Arguments]  @{headers}  ${entity}
    page should contain element  ${ENTITIES.TABLE.BODY}
    FOR  ${row_element}  IN RANGE  1  6
        log  ${row_element}
        Iterate over row  @{headers}  ${col_ammount}  ${row_element}
    END

Validate Table loaded
    # @{headers}: the list of headlines for columns
    # @{entities}: the list of created entities
    [Arguments]  @{headers}  # @{entities}=${EMPTY}
    Keywords.Validate visible element  ${ENTITIES.TABLE}
    Validate Table Headers  @{headers}

    # testing
    ${col_ammount} =  Get length    ${headers}

    Get Table Headers

    FOR  ${row_element}  IN RANGE  2  6
        log  ${row_element}
        Iterate over row  ${headers}  ${col_ammount}  ${row_element}
    END


### --  Table Related Keywords  -- ###
Change number of entities per page
    [Arguments]  ${rows_per_page}=25
    Keywords.Mouse over and click element  ${ENTITIES.TABLE.ROWS_PER_PAGE}
    Keywords.Select Item from the list  ${ENTITIES.TABLE.ROWS_PER_PAGE.LIST}  ${rows_per_page}