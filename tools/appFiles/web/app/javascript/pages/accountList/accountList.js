import React from 'react';
import ReactDom from 'react-dom';
import AccountList from './accountListComponents';
import Page from '../theme';

// window.Load_accountList = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <AccountList/>
//     </Page>
//   );
// };

window.Load_accountList = function() {
  if (window.appState.AccountTypeShort == "deal") {
      document.title = window.global.functions.productTitle() + window.appContent.TitlesCustomerList;
  } else {
      document.title = window.global.functions.productTitle() + window.appContent.TitlesAccountList;
  }

  window.Load_settings();
};
