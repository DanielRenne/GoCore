import React from 'react';
import ReactDom from 'react-dom';
import Roles from './rolesComponents';
import Page from '../theme';

window.Load_roles = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoles;
  window.global.functions.render(
    <Page>
        <Roles/>
    </Page>
  );
};
