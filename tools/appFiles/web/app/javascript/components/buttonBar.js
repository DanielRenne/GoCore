import React, {Component} from 'react';
import {BottomNavigation, BottomNavigationItem} from 'material-ui/BottomNavigation';
import Paper from 'material-ui/Paper';
import BaseComponent from './base';
import {Tabs, Tab} from 'material-ui/Tabs';

class ButtonBar extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    let selected = this.props.tabOrder[0];
    $.each(this.props.tabOrder, (k,v) => {
      if (v == this.props.selectedTab) {
        selected = v;
      }
    });
    this.state = {
      initialSelectedIndex: selected
    };
    this.select = (index) => {
      this.setComponentState({initialSelectedIndex: index})
    };
  }

  //componentDidUpdate() {
    // Hack because its tough to setState when passing children around after the fact
    //
    //$('.' + this.state.selectedIndex + ' .material-icons').css();
    //window.goCore.updateIconState({style: {color: MuiThemes.default.palette.accent1Color}});
    //console.log('asfasfdfsda',this.props.tabIcons[this.state.selectedIndex].setState({}));
  //}

  renderLabel(name) {
    if (name == "AccountList" && window.appState.AccountTypeShort == "deal") {
      return window.appContent.CustomersTitle;
    } else {
      return window.appContent[name]
    }
  }

  render() {
    try {
      this.logRender();
      let buttons = [];
      let i = 0;
      let b = 0;
      let nameToIndex = {};
      let otherTabSelected = {};
      $.each(this.props.tabOrder, (k, obj) => {
        let name = obj;
        var controller = this.props.tabControllers[i];
        var action = this.props.tabActions[i];
        var tabIsVisible = this.props.tabIsVisible[i];
        otherTabSelected[name] = this.props.otherTabSelected[i];
        var icon = this.props.tabIcons[name];
        var params = {};
        if (this.props.tabUriParams[name] != undefined) {
          params = this.props.tabUriParams[name];
        }
        if (tabIsVisible) {
          nameToIndex[name] = b;
          let button = <Tab key={obj} className={obj} label={this.renderLabel(name)} icon={icon} onActive={
            (event) => {
              this.select(obj);
              var parms = {action: action, uriParams: params};
              if (controller) {
                parms.controller = controller;
              }
              window.api.get(parms);
            }
          }/>;
          b++;
          buttons.push(button);
        }
        i++;
      });

      var tpl = null;
      if (buttons.length > 0) {
        tpl = (
            <Paper zDepth={2}>
              <Tabs className="buttonbar buttonbar-el"
                    initialSelectedIndex={(otherTabSelected[this.state.initialSelectedIndex] != "") ? nameToIndex[otherTabSelected[this.state.initialSelectedIndex]] : nameToIndex[this.state.initialSelectedIndex]}>
                {buttons}
              </Tabs>
            </Paper>
        );
      }
      return tpl;
    } catch (e) {
      return this.globs.ComponentError("ButtonBar", e.message, e);
    }
  }
}

ButtonBar.propTypes = {
  tabs: React.PropTypes.object.isRequired,
  tabActions: React.PropTypes.array,
  tabControllers: React.PropTypes.array,
  selectedTab: React.PropTypes.string.isRequired,
  tabOrder: React.PropTypes.array,
  tabIconsLarge: React.PropTypes.object,
  tabIcons: React.PropTypes.object
};

export default ButtonBar;
