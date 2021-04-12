import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { executeCommand } from 'mattermost-redux/actions/integrations';

import RHSView from './rhs_view';

const mapStateToProps = (state) => {
    const currentUserId = state.entities.users.currentUserId;
    const currentChannelId = state.entities.channels.currentChannelId;
    const currentTeamId = state.entities.teams.currentTeamId;

    return {
        user: state.entities.users.profiles[currentUserId],
        channel: state.entities.channels.channels[currentChannelId],
        team: state.entities.teams.teams[currentTeamId],
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    executeCommand,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(RHSView);
