import React from 'react';
import ReactDom from 'react-dom';
import ServerSettingsModify from './serverSettingsModifyComponents';
import Page from '../theme';

window.Load_serverSettingsModify = function() {
  console.log(window.pageState);
  window.global.functions.render(
    <Page>
        <ServerSettingsModify/>
    </Page>
  );
};

// window.Load_serverSettingsModify = function() {
//   document.title = window.global.functions.productTitle() + window.appContent.TitlesServerSettingsModify;
//   // window.Load_settings();
// };
