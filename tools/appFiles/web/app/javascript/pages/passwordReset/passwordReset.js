import React from 'react';
import ReactDom from 'react-dom';
import PasswordReset from './passwordResetComponents';
import Page from '../theme';

window.Load_passwordReset = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesResetPassword;
  window.global.functions.render(
    <Page>
        <PasswordReset/>
    </Page>
  );
};
