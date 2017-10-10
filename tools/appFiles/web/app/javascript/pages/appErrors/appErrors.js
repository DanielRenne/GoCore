import React from 'react';
import ReactDom from 'react-dom';
import AppErrors from './appErrorsComponents';
import Page from '../theme';

window.Load_appErrors = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesAppErrors;
  window.global.functions.render(
    <Page>
        <AppErrors/>
    </Page>
  );
};
