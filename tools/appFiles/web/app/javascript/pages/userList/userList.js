import React from 'react';
import ReactDom from 'react-dom';
import UserList from './userListComponents';
import Page from '../theme';

// window.Load_userListFull = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <UserList/>
//     </Page>
//   );
// };

window.Load_userList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesUserList;
  window.Load_settings();
};
