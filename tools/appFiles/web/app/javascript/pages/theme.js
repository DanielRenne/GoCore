import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from './colors';
import React, {Component} from 'react';
import BaseComponent from '../components/base';
import GoCoreParentNode from '../components/goCoreParentNode';

class Page extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
      this.logRender();
      return (
        <GoCoreParentNode>
          <MuiThemeProvider muiTheme={MuiThemes[this.props.muiThemeName]}>
            {this.props.children}
          </MuiThemeProvider>
        </GoCoreParentNode>
      );
  }
}

Page.propTypes = {
  muiThemeName: React.PropTypes.string
};

Page.defaultProps = {
    muiThemeName: "default"
};

export default Page;
