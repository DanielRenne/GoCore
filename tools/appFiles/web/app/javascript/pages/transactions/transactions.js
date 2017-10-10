import React from 'react';
import ReactDom from 'react-dom';
import Transactions from './transactionsComponents';
import Page from '../theme';

window.Load_transactions = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesTransactions;
  window.global.functions.render(
    <Page>
        <Transactions/>
    </Page>
  );
};
