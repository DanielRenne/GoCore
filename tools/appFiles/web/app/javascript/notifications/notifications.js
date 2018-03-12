import React from 'react';
import ReactDom from 'react-dom';
import Nofications from './notifications-component';
import Page from '../pages/theme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from '../pages/colors';

window.LoadSiteNotifications = function() {
    window.goCore.notifications = ReactDom.render(
        <MuiThemeProvider muiTheme={MuiThemes.default}>
            <Nofications/>
        </MuiThemeProvider>
    , document.getElementById('GoCore-notifications'));
};
