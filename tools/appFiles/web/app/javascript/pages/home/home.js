import React from 'react';
import ReactDom from 'react-dom';
import Login from '../login/loginComponents';
import Home from './homeComponents';
import Page from '../theme';

window.Load_home = function() {

  if (!appState.loggedIn) {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesDashboard;
    window.global.functions.render(
        <Page>
            <Login/>
        </Page>
    );
    
  } else {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesLoginToTheSystem;
    window.global.functions.render(
        <Page>
            <Home/>
        </Page>
    );
  }
};

