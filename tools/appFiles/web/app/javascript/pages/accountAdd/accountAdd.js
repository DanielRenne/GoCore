import React from 'react';
import ReactDom from 'react-dom';
import AccountAdd from './accountAddComponents';
import Page from '../theme';

// window.Load_accountAdd = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <AccountAdd/>
//     </Page>
//   );
// };

window.Load_accountAdd = function() {
  if (window.appState.AccountTypeShort == "deal") {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesAddANewCustomer;
  } else {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesAddANewAccount;
  }

  window.Load_settings();
};
