import { id } from 'manifest';
import React from 'react';
import PropTypes from 'prop-types';

import "./style.scss"


export default class ChannelHeaderButtonIcon extends React.PureComponent {
	static propTypes = {
		theme: PropTypes.object.isRequired,
	}

	constructor(props) {
		super(props);
	}

	render() {
		console.log(this.props.theme.sidebarText)
		const mask = `url(/plugins/${id}/dice-d20-solid.svg) no-repeat center`
		return (
			<div className="test" style={{
				WebkitMask: mask, mask: mask
			}
			} />
		)
	}
}
