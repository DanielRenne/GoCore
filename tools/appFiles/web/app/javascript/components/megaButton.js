import React from 'react';
import {Avatar, blueGrey100, IconButton} from "../globals/forms";
import MuiThemes from '../pages/colors'
import ActionAdd from 'material-ui/svg-icons/content/add';
import BaseComponent from './base';

class MegaIcon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  getIcon() {
    if (this.props.icon) {
      return this.props.icon;
    } else {
      return <ActionAdd color={this.props.iconColor}/>;
    }
  }


  render() {
    try {
      this.logRender();
      var outputButtonSize = true;
      if (this.props.size > 100) {
        var size1 = this.props.size - (this.props.size/8);
        var size2 = this.props.size - (this.props.size/4);
      } else if (this.props.size > 75) {
        var size1 = this.props.size - 15;
        var size2 = this.props.size - 40;
      } else if (this.props.size > 50 && this.props.size < 75) {
        var size1 = this.props.size/1.5;
        var size2 = this.props.size/2;
      } else if (this.props.size <= 50) {
        var size1 = this.props.size/1.5;
        var size2 = this.props.size/2;
        outputButtonSize = false;
      }
      var iconStyle= {color: this.props.iconColor, padding: 8};
      if (outputButtonSize) {
        iconStyle.width = size1;
        iconStyle.height = size1;
      }

      return (
          <div>
            <Avatar backgroundColor={this.props.avatarColor} style={{cursor: "pointer"}} size={this.props.size} onClick={this.props.onClick}>
              <IconButton tooltip={this.props.tooltip} disableTouchRipple={this.props.disableTouchRipple} style={iconStyle} iconStyle={{
                width: size2,
                height: size2,
                color: this.props.iconColor
              }}>
                {this.getIcon()}
              </IconButton>
            </Avatar>
          </div>
      );
    } catch(e) {
      return this.globs.ComponentError("MegaIcon", e.message, e);
    }
  }
}

MegaIcon.propTypes = {
  avatarColor: React.PropTypes.string,
  iconColor: React.PropTypes.string,
  tooltip: React.PropTypes.string,
  disableTouchRipple: React.PropTypes.bool,
  icon: React.PropTypes.node,
  onClick: React.PropTypes.func,
  size: React.PropTypes.number,
};

MegaIcon.defaultProps = {
  disableTouchRipple: false,
  size: 300,
  avatarColor: MuiThemes.default.palette.accent3Color,
  iconColor: blueGrey100,
};

export default MegaIcon;