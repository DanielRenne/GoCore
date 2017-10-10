import React from 'react';
import ReactDom from 'react-dom';
import SideBarMenu from './sidebarMenu-component';
import Page from '../pages/theme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from '../pages/colors';

window.collapseSideBarMenu = function(){
  var bar = $(".site-menubar");

  if(bar.css("opacity") != 0){
    $("#btnSideBarMenuCollapse").click();
  }
};

window.expandSideBarMenu = function(){
  $("#btnSideBarMenuCollapse").click();
}

window.LoadSideBarMenu = function() {
    window.goCore.sideBarMenu = ReactDom.render(
        <MuiThemeProvider muiTheme={MuiThemes.default}>
            <SideBarMenu/>
        </MuiThemeProvider>
    , document.getElementById('GoCore-sidebarMenu'));
};
