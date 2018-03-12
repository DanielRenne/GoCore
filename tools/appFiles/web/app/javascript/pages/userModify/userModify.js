import React from 'react';
import ReactDom from 'react-dom';
import UserAddModify from './userModifyPage';
import Page from '../theme';

// window.Load_userModify = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <UserAddModify/>
//     </Page>
//   );
// };

window.Load_userModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesUserModify;
  window.Load_settings();
};
