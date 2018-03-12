import React from 'react';
import ReactDom from 'react-dom';
import AccountAddModify from './accountModifyComponents';
import Page from '../theme';

// window.Load_accountModify = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <AccountAddModify/>
//     </Page>
//   );
// };

window.Load_accountModify = function() {
  if (window.appState.AccountTypeShort == "deal") {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesModifyCustomer;
  } else {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesModifyAccount;
  }

  window.Load_settings();
};
