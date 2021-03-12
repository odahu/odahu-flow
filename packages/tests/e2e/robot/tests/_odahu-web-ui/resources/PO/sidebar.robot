*** Settings ***
Documentation   Page Object for "sidebar" part of UI
Library         SeleniumLibrary  timeout=10s
Library         Collections
Resource        ${RES_DIR}/keywords.robot

*** Variables ***
${SIDEBAR.SANDWICH_BUTTON}  xpath://*[@id="root"]/div/header/div/button

# link icon classes
${SIDEBAR.ACTIVE_LINK_ICON_CLASS}    MuiListItemIcon-root jss24 jss25
${SIDEBAR.INACTIVE_LINK_ICON_CLASS}  MuiListItemIcon-root jss24

# sidebar classes
${SIDEBAR.DRAWER.EXTENDED}          MuiDrawer-root MuiDrawer-docked jss20 jss21
${SIDEBAR.SIDEBAR.EXTENDED}         MuiPaper-root MuiDrawer-paper jss21 MuiDrawer-paperAnchorLeft MuiDrawer-paperAnchorDockedLeft MuiPaper-elevation0
${SIDEBAR.DRAWER.SHRINKED}          MuiDrawer-root MuiDrawer-docked jss20 jss22
${SIDEBAR.SIDEBAR.SHRINKED}         MuiPaper-root MuiDrawer-paper jss22 MuiDrawer-paperAnchorLeft MuiDrawer-paperAnchorDockedLeft MuiPaper-elevation0

${SIDEBAR.DRAWER.LOCATOR}           xpath://*[@id="root"]/div/div[2]
${SIDEBAR.SIDEBAR.LOCATOR}          xpath://*[@id="root"]/div/div[2]/div

# list of (icons) links to ODAHU Web UI pages
@{SIDEBAR.LINKS_LIST}           ${SIDEBAR.LINK.DASHBOARD.ICON}  ${SIDEBAR.LINK.TRAININGS.ICON}  ${SIDEBAR.LINK.PACKAGINGS.ICON}  ${SIDEBAR.LINK.DEPLOYMENTS.ICON}
...                             ${SIDEBAR.LINK.CONNECTIONS.ICON}  ${SIDEBAR.LINK.TOOLCHAINS.ICON}  ${SIDEBAR.LINK.PACKAGERS.ICON}

# icons for links to ODAHU Web UI pages
               # description path -> xpath://*[@id="root"]/div/div[2]/div/ul/a[1]/div/div[2]/span
${SIDEBAR.LINK.DASHBOARD.ICON}       xpath://*[@id="root"]/div/div[2]/div/ul/a[1]/div/div[1]
${SIDEBAR.LINK.TRAININGS.ICON}       xpath://*[@id="root"]/div/div[2]/div/ul/a[2]/div/div[1]
${SIDEBAR.LINK.PACKAGINGS.ICON}      xpath://*[@id="root"]/div/div[2]/div/ul/a[3]/div/div[1]
${SIDEBAR.LINK.DEPLOYMENTS.ICON}     xpath://*[@id="root"]/div/div[2]/div/ul/a[4]/div/div[1]
${SIDEBAR.LINK.CONNECTIONS.ICON}     xpath://*[@id="root"]/div/div[2]/div/ul/a[5]/div/div[1]
${SIDEBAR.LINK.TOOLCHAINS.ICON}      xpath://*[@id="root"]/div/div[2]/div/ul/a[6]/div/div[1]
${SIDEBAR.LINK.PACKAGERS.ICON}       xpath://*[@id="root"]/div/div[2]/div/ul/a[7]/div/div[1]

*** Keywords ***
Click "Sandwich Menu" button
    click button  ${SIDEBAR.SANDWICH_BUTTON}

Open ODAHU page
    [Arguments]  ${odahu-ui-page locator}
    wait until page contains element  ${odahu-ui-page locator}
    click element  ${odahu-ui-page locator}

Validate SideBar exists and visible
    element should be visible  ${SIDEBAR.SIDEBAR.LOCATOR}

Validate that active page icon changed the color
    [Arguments]  ${locator}
    Validate page contains the "Element" of the "Class"  ${locator}  ${SIDEBAR.ACTIVE_LINK_ICON_CLASS}

Validate that inactive page icon has the default color
    [Arguments]  ${locator}
    Validate page contains the "Element" of the "Class"  ${locator}  ${SIDEBAR.INACTIVE_LINK_ICON_CLASS}

Validate that the active page icon has one color and the others different
    [Arguments]  ${active page locator}
    @{icons list}=  create list  @{SIDEBAR.LINKS_LIST}
    remove values from list  ${icons list}  ${active page locator}
    Validate that active page icon changed the color  ${active page locator}
    FOR  ${inactive page icon}  IN  @{icons list}
        Validate that inactive page icon has the default color  ${inactive page icon}
    END

Validate that "SideBar" is extended
    Validate page contains the "Element" of the "Class"  ${SIDEBAR.DRAWER.LOCATOR}  ${SIDEBAR.DRAWER.EXTENDED}
    Validate page contains the "Element" of the "Class"  ${SIDEBAR.SIDEBAR.LOCATOR}  ${SIDEBAR.SIDEBAR.EXTENDED}

Validate that "SideBar" is shrinked
    Validate page contains the "Element" of the "Class"  ${SIDEBAR.DRAWER.LOCATOR}  ${SIDEBAR.DRAWER.SHRINKED}
    Validate page contains the "Element" of the "Class"  ${SIDEBAR.SIDEBAR.LOCATOR}  ${SIDEBAR.SIDEBAR.SHRINKED}

Validate that "ODAHU tab" links are visible on "SideBar"
    @{links list}=  create list  @{SIDEBAR.LINKS_LIST}
    FOR  ${link}  IN  @{links list}
        element should be visible  ${link}
    END
