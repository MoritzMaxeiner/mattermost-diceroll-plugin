import React from 'react';
import PropTypes from 'prop-types';

import { FormattedMessage } from 'react-intl';

import { Tabs, Tab } from 'react-bootstrap';

import { Panel, ListGroup, ListGroupItem } from 'react-bootstrap'

const { formatText, messageHtmlToComponent } = window.PostUtils;

import { Button, Form, FormGroup, FormControl, ControlLabel, Row, Col, Nav, NavItem } from 'react-bootstrap';
import { getTheme } from 'mattermost-redux/selectors/entities/preferences';

import "./style.scss"

export default class RHSView extends React.PureComponent {
    static propTypes = {
        theme: PropTypes.object.isRequired,
        user: PropTypes.object.isRequired,
        channel: PropTypes.object.isRequired,
        team: PropTypes.object.isRequired,
        executeCommand: PropTypes.func.isRequired,
    }

    constructor(props) {
        super(props);
        this.state = { numDice: 1, diceSystem: "", numDieSides: 2 }

        this.onRoll = this.onRoll.bind(this)
        this.onStateChange = this.onStateChange.bind(this)
    }

    onRoll(event) {
        event.preventDefault();

        const args = {
            channel_id: this.props.channel.id,
            team_id: this.props.team.id,
        }

        const numDice = this.state.numDice
        const typeDice = this.state.diceSystem != "" ? this.state.diceSystem : this.state.numDieSides

        this.props.executeCommand(`/roll ${numDice}d${typeDice}`, args);
    }

    onStateChange(event) {
        const target = event.target;
        const value = target.value;
        const name = target.name;

        this.setState({ [name]: value });
    }

    render() {
        const style = getStyle(this.props.theme)
        return (
            <div style={style.rhs}>
                <Tab.Container id="tabsDiceRoller" defaultActiveKey="generic">
                    <Row>
                        <Col sm={12}>
                            <Nav bsStyle="tabs">
                                <NavItem eventKey="generic">
                                    Generic
                                </NavItem>
                                <NavItem eventKey="about" className="pull-right">
                                    About
                                </NavItem>
                            </Nav>
                        </Col>
                        <Col sm={12}>
                            <Tab.Content animation>
                                <Tab.Pane eventKey="generic">
                                    <Form horizontal onSubmit={this.onRoll}>
                                        <FormGroup
                                            controlId="formNumDice"
                                        >
                                            <Col componentClass={ControlLabel} sm={4}><span className="pull-left">Number of Dice</span></Col>
                                            <Col sm={4} smOffset={4}>
                                                <FormControl type="number" name="numDice" value={this.state.numDice} min="1" onChange={this.onStateChange} />
                                            </Col>
                                        </FormGroup>
                                        <FormGroup controlId="formTypeDice">
                                            <Col componentClass={ControlLabel} sm={4}><span className="pull-left">Type of Dice</span></Col>
                                            <Col sm={4}>
                                                <FormControl componentClass="select" name="diceSystem" value={this.state.diceSystem} onChange={this.onStateChange}>
                                                    <option value="">N-Sided</option>
                                                    <option value="AE">Aetherium</option>
                                                </FormControl>
                                            </Col>
                                            <Col sm={4}>
                                                <FormControl style={this.state.diceSystem != "" ? { display: 'none' } : {}} type="number" name="numDieSides" value={this.state.numDieSides} min="2" onChange={this.onStateChange} />
                                            </Col>
                                        </FormGroup>
                                        <FormGroup controlId="formSubmit">
                                            <Col sm={2} smOffset={5}>
                                                <Button type="submit" className="btn-primary">Roll</Button>
                                            </Col>
                                        </FormGroup>
                                    </Form>
                                </Tab.Pane>
                                <Tab.Pane eventKey="about">
                                    <h3>Licensing Information</h3>
                                    <ListGroup>
                                        <ListGroupItem>{messageHtmlToComponent(formatText("[\"Dice D 20 Icon\"](https://fontawesome.com/icons/dice-d20/) by [Font Awesome](https://fontawesome.com/) is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)"))}</ListGroupItem>
                                    </ListGroup>
                                </Tab.Pane>
                            </Tab.Content>
                        </Col>
                    </Row>
                </Tab.Container>
            </div >
        );
    }
}

const getStyle = (theme) => ({
    button: {
        color: theme.buttonColor,
        backgroundColor: theme.buttonBg,
    },
    rhs: {
        padding: '10px',
    },
});

