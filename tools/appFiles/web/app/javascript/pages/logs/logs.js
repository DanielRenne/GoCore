import React from 'react';
import ReactDom from 'react-dom';
import Logs from './logsComponents';
import Page from '../theme';

window.Load_logs = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesLogs;
  window.global.functions.render(
    <Page>
        <Logs/>
    </Page>
  );
};
