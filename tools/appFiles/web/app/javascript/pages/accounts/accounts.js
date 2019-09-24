import React from 'react';
import ReactDom from 'react-dom';
import Accounts from './accountsComponents';
import Page from '../theme';

window.Load_accounts = function() {
  window.global.functions.render(
    <Page>
        <Accounts/>
    </Page>
  );
};
