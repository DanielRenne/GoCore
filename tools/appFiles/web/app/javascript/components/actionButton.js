import React, {Component} from 'react';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import Tooltip from 'material-ui/internal/Tooltip';
import ContentAdd from 'material-ui/svg-icons/content/add';
import ContentEdit from 'material-ui/svg-icons/content/create';
import ContentDelete from 'material-ui/svg-icons/content/clear';
import ContentBack from 'material-ui/svg-icons/hardware/keyboard-backspace';
import {deepOrange500, blueGrey500} from 'material-ui/styles/colors';
import BaseComponent from './base';

class ActionButton extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.__ = window.global.functions._;
    this.state = window.pageState;
    this.state.showToolTip = false;

    this.handleHoverFloatButton = (id) =>  {
      this.setComponentState({showToolTip: true});
    };

    this.handleHoverOutFloatButton = (id) =>  {
      this.setComponentState({showToolTip: false});
    };
  }

  render() {
    try {
      this.logRender();
      var actionColor = deepOrange500;
      var buttonTooltip = "";
      var buttonTooltip = this.props.tooltip;

      var iconType = (<ContentAdd/>);
      if (this.props.type == "Edit"){
        iconType = (<ContentEdit/>);
      }

      if (this.props.type == "Delete"){
        iconType = (<ContentDelete/>);
      }

      if (this.props.type == "Back"){
        iconType = (<ContentBack/>);
        actionColor = blueGrey500;
        buttonTooltip = "Back";
      }

      var visible = "inherit";
      if (this.props.visible === false){
        visible = "none";
      }

      return (
        <div style={{
            margin: 0,
            top: 'auto',
            right: ($(window).width() < 500) ? 15 : 35,
            bottom: ($(window).width() < 500) ? 15 : 50,
            left: 'auto',
            position: 'fixed',
            display: visible,
            zIndex: 999
          }}>
          <FloatingActionButton backgroundColor={actionColor}
                                onMouseOver={this.handleHoverFloatButton}
                                onMouseOut={this.handleHoverOutFloatButton}
                                onTouchTap={ (e) => this.globs.FloatingActionButtonClick(e, this.props.onClick, this.props.type, this.props.action) }
          >
            {iconType}
          </FloatingActionButton>
          <Tooltip show={this.state.showToolTip}
                   label={buttonTooltip}
                   style={{right: 62, top:16, height:26, fontSize:'10px'}}
                   horizontalPosition="left"
                   verticalPosition="top"
                   touch={true}
          />
        </div>

      );
    } catch(e) {
      return this.globs.ComponentError(this.getClassName(), e.message);
    }
  }
}

ActionButton.propTypes= {
  onClick: React.PropTypes.func,
  tooltip: React.PropTypes.string,
  visible: React.PropTypes.bool
};

export default ActionButton;
