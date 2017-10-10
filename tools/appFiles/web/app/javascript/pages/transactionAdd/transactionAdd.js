import React from 'react';
import ReactDom from 'react-dom';
import TransactionAdd from './transactionAddComponents';
import Page from '../theme';
const settings_page = false;

window.Load_transactionAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesTransactionAdd;

  if (settings_page) {
    window.Load_settings();
  } else {
    window.global.functions.render(
      <Page>
          <TransactionAdd/>
      </Page>
    );
  }
};
