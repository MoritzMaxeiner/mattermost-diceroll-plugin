import { id } from 'manifest';
import { connect } from 'react-redux';

import { GlobalState } from 'mattermost-redux/types/store';
import { getTheme } from 'mattermost-redux/selectors/entities/preferences';

import ChannelHeaderButtonIcon from './channel_header_button_icon';

const mapStateToProps = (state: GlobalState) => ({
	theme: getTheme(state),
})

export default connect(mapStateToProps)(ChannelHeaderButtonIcon);
