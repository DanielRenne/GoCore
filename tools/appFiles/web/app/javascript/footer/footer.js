import React from 'react';
import ReactDom from 'react-dom';
import Footer from './footer-component';
import Page from '../pages/theme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import MuiThemes from '../pages/colors';


window.LoadSiteFooter = function() {
    window.goCore.footer = ReactDom.render(
        <MuiThemeProvider muiTheme={MuiThemes.default}>
            <Footer/>
        </MuiThemeProvider>
    , document.getElementById('GoCore-footer'));
};