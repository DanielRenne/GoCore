import React from 'react';
import ReactDom from 'react-dom';
import Banner from './banner-component';
import Page from '../pages/theme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from '../pages/colors';


window.collapseNavbarAvatar = function(){
  var navbarAvatar = $(".navbar-avatar");
  if (navbarAvatar.attr( "aria-expanded") == "true") {
    navbarAvatar.click();
  }
};

window.collapseBannerComponents = function(){
  var searchBar = $(".navbar-search-overlap");

  if(searchBar.css("display") != "none"){
    $(".input-search-close").click();
  }

  window.collapseSideBarMenu();
};

window.LoadSiteBanner = function() {
    window.goCore.banner = ReactDom.render(
        <MuiThemeProvider muiTheme={MuiThemes.default}>
          <Banner/>
        </MuiThemeProvider>
    , document.getElementById('GoCore-banner'));
};
