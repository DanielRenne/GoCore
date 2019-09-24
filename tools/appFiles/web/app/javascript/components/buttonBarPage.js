import React, {Component} from 'react';
import ButtonBar from './buttonBar';
import Page from '../pages/theme';
import BaseComponent from './base';
import GoCoreParentNode from './goCoreParentNode';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from '../pages/colors';


class ButtonBarPage extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.pageRef;
  }

  renderTitle() {
    if (this.props.page.props.selectedTab == "AccountList" && window.appState.AccountTypeShort == "deal") {
      return window.appContent.CustomersTitle;
    } else if (this.props.page.props.selectedTab == "AccountModify" && window.appState.AccountTypeShort == "deal"){
      return window.appContent.ModifyCustomersTitle;
    } else if (this.props.page.props.selectedTab == "AccountAdd" && window.appState.AccountTypeShort == "deal"){
      return window.appContent.TitlesAddANewCustomer;
    } else {
      return window.appContent[this.props.page.props.selectedTab]
    }
  }

  render() {
    try {
      this.logRender();

      var buttonBarPageTitle = this.renderTitle();

      return (
          <GoCoreParentNode>
            {(this.props.title) && $(window).height() > 750 ? <h1 className="buttonbar-el" style={{margin: '0em 0 .5em'}}><MuiThemeProvider muiTheme={MuiThemes.opposite}>{this.props.icon}</MuiThemeProvider> {this.props.title}</h1>: null}
            <MuiThemeProvider muiTheme={MuiThemes.default}>
               <ButtonBar className="buttonbar-el" {...this.props.page.props}/>
            </MuiThemeProvider>
            <div className={"margin-top " + ((window.location.hash.indexOf("List") == -1 && $(window).width() < 750) ? "margin-left-right": "")}>
              <MuiThemeProvider muiTheme={MuiThemes.default}>
                 {this.props.page}
               </MuiThemeProvider>
            </div>
          </GoCoreParentNode>
      );
    } catch(e) {
      return this.globs.ComponentError("ButtonBarPage", e.message, e);
    }
  }
}

ButtonBarPage.propTypes = {
  page: React.PropTypes.element,
  title: React.PropTypes.string,
  showMenuPageTitle: React.PropTypes.bool,
  icon: React.PropTypes.element
};

ButtonBarPage.defaultProps = {
    showMenuPageTitle: true
};

export default ButtonBarPage;
