import React from 'react';
import ReactDom from 'react-dom';
import RoleFeatures from './roleFeaturesComponents';
import Page from '../theme';

window.Load_roleFeatures = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleFeatures;
  window.global.functions.render(
    <Page>
        <RoleFeatures/>
    </Page>
  );
};
