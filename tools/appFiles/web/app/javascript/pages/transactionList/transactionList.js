import React from 'react';
import ReactDom from 'react-dom';
import TransactionList from './transactionListComponents';
import Page from '../theme';

const settings_page = false;

window.Load_transactionList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesTransactionList;
  if (settings_page) {
    window.Load_settings();
  } else {
    window.global.functions.render(
      <Page>
          <TransactionList/>
      </Page>
    );
  }
};
