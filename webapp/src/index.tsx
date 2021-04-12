import { Store, Action } from 'redux';

import { GlobalState } from 'mattermost-redux/types/store';

import manifest from './manifest';

import React from 'react';
import { FormattedMessage } from 'react-intl';

// eslint-disable-next-line import/no-unresolved
import { PluginRegistry } from './types/mattermost-webapp';


import ChannelHeaderButtonIcon from './components/channel_header_button_icon'

import RHSView from './components/right_hand_sidebar';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        const { toggleRHSPlugin } = registry.registerRightHandSidebarComponent(
            RHSView,
            <FormattedMessage
                id='plugin.name'
                defaultMessage="Dice Roller"
            />
        )

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon />,
            () => store.dispatch(toggleRHSPlugin),
            <FormattedMessage
                id='plugin.name'
                defaultMessage='Dice Roller'
            />,
            'Throw some dice'
        )
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
