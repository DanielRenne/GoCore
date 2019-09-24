import React from 'react';
import ReactDom from 'react-dom';
import AccountList from '../accountList/accountListComponents';
import UserList from '../userList/userListComponents';
import UserAddModify from '../userModify/userModifyComponents';
import UserModify from '../userModify/userModifyPage';
import UserAdd from '../userAdd/userAddComponents';
//SettingsJSImports
import RoleList from '../roleList/roleListComponents';
import RoleAddModify from '../roleModify/roleModifyComponents';
import RoleAdd from '../roleAdd/roleAddComponents';
import UserProfile from '../userProfile/userProfileComponents';
import AccountAdd from '../accountAdd/accountAddComponents';
import AccountAddModify from '../accountModify/accountModifyComponents';
import ServerSettings from '../serverSettingsModify/serverSettingsModifyComponents'
import ButtonBarPage from '../../components/buttonBarPage'
import FontIcon from 'material-ui/FontIcon';
import {serverSettingsIconLarge, userListIconLarge, accountsIconLarge, userProfileIconLarge, notificationsIconLarge, serverSettingsIcon, userListIcon, accountsIcon, userProfileIcon, notificationsIcon, businessIcon, businessIconLarge, blockIcon, blockIconLarge} from '../../globals/icons'
// Note all pageContent used in any rendered component for all pages here needs to stay in settings/en/US.json


window.Load_settings = function() {
  var state = window.pageState;
  var page = null;
  var props = {};
  props.tabs = state.SettingsBar.ButtonBar.Config.VisibleTabs;
  props.selectedTab = state.SettingsBar.ButtonBar.Config.CurrentTab;
  props.tabOrder = state.SettingsBar.ButtonBar.Config.TabOrder;
  props.tabActions = state.SettingsBar.ButtonBar.Config.TabActions;
  props.tabControllers = state.SettingsBar.ButtonBar.Config.TabControllers;
  props.tabIsVisible = state.SettingsBar.ButtonBar.Config.TabIsVisible;
  props.otherTabSelected = state.SettingsBar.ButtonBar.Config.OtherTabSelected;
  props.tabIconsLarge = {};
  props.tabIcons = {};
  props.tabUriParams = {};
  props.tabOrder.forEach((value) => {
      let icon = '';
      let iconsm = '';
      let params = {};
      switch(value) {
          //SettingsIconSwitch
          case "RoleList":
          case "RoleAdd":
          case "RoleModify":
            icon = blockIconLarge;
            iconsm = blockIcon;
            break;

          case "AccountList":
          case "AccountAdd":
          case "AccountModify":
            icon = accountsIconLarge;
            iconsm = accountsIcon;
            break;
          case "UserList":
          case "UserAdd":
            icon = userListIconLarge;
            iconsm = userListIcon;
            break;
          case "ModifyServerSettings":
            icon = serverSettingsIconLarge;
            iconsm = serverSettingsIcon;
            break;
          case "UserModify":
          case "UserProfile":
            icon = userProfileIconLarge;
            iconsm = userProfileIcon;
            break;
      }
    props.tabIconsLarge[value] = icon;
    props.tabIcons[value] = iconsm;
    props.tabUriParams[value] = params;
  });

  if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.AccountList) {
    page = (<AccountList {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.AccountModify) {
    page = (<AccountAddModify {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.AccountAdd) {
    page = (<AccountAdd {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.RoleList) {
    page = (<RoleList {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.RoleModify) {
    page = (<RoleAddModify {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.RoleAdd) {
    page = (<RoleAdd {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.UserList) {
    page = (<UserList {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.UserAdd) {
    page = (<UserAdd {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.UserProfile) {
    page = (<UserProfile {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.UserModify) {
    page = (<UserModify {...props}/>);
  } else if (state.SettingsBar.ButtonBar.Config.CurrentTab == state.SettingsBar.Constants.ModifyServerSettings) {
    page = (<ServerSettings {...props}/>);
  }

  window.global.functions.render(
    <ButtonBarPage page={page} title={appContent.SrcSvgSettings} icon={
      <FontIcon className="large material-icons">settings</FontIcon>
    }/>
  );
};
