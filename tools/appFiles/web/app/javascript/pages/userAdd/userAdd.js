import React from 'react';
import ReactDom from 'react-dom';
import UserAdd from './userAddComponents';
import Page from '../theme';

// window.Load_userAdd = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <UserAdd/>
//     </Page>
//   );
// };

window.Load_userAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesUserAdd;
  window.Load_settings();
};
