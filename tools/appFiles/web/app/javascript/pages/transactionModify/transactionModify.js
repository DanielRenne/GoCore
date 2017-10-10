import React from 'react';
import ReactDom from 'react-dom';
import TransactionAddModify from './transactionModifyComponents';
import Page from '../theme';

const settings_page = false;

window.Load_transactionModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesModifyTransaction;
  if (settings_page) {
    window.Load_settings();
  } else {
    window.global.functions.render(
      <Page>
          <TransactionAddModify/>
      </Page>
    );
  }
};
