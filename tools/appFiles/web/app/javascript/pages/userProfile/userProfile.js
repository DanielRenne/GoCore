import React from 'react';
import ReactDom from 'react-dom';
import UserProfile from './userProfileComponents';
import Page from '../theme';

// window.Load_userProfile = function() {
//   console.log(window.pageState);
//   window.global.functions.render(
//     <Page>
//         <UserProfile/>
//     </Page>
//   );
// };

window.Load_userProfile = function() {
    document.title = window.global.functions.productTitle() + window.appContent.TitlesUserProfile;
    window.Load_settings();
};