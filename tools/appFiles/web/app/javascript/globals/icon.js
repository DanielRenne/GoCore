import FontIcon from 'material-ui/FontIcon';
import Page from '../pages/theme';
import React, {Component} from 'react';
import BaseComponent from '../components/base';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from '../pages/colors';

class Icon extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    if (!this.props.initialStyle) {
      var style = {color: 'white'};
    } else {
      var style = this.props.initialStyle;
    }

    this.state = {
      icon_tag: '',
      style: style
    };

    // window.goCore.updateIconState = function(state) {
    //   $('.MainButtonBar .material-icons').css({color: 'rgba(0, 0, 0, 0.541176)'});
    //   this.setState(state);
    // }.bind(this);
  }

  render() {
      this.logRender();
      return (
        <MuiThemeProvider muiTheme={MuiThemes.default}>
          <FontIcon style={this.state.style} className={(this.props.large) ? 'large material-icons': 'material-icons'}>{this.props.iconTag}</FontIcon>
        </MuiThemeProvider>
      );
  }
}

Icon.propTypes = {
  initialStyle: React.PropTypes.object,
  iconTag: React.PropTypes.string,
  large: React.PropTypes.bool
};

Icon.defaultProps = {
    large: false,
};

export default Icon;
