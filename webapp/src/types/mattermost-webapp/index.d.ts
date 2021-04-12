export interface PluginRegistry {
    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference

    registerRightHandSidebarComponent(component, title)

    registerChannelHeaderButtonAction(icon, action, dropdownText, tooltipText)
}
