import React from 'react';
import ReactDom from 'react-dom';
import -CAPCAMEL- from './-CAMEL-Components';
import Page from '../theme';

window.Load_-CAMEL- = function() {
  document.title = window.global.functions.productTitle() + window.appContent.Titles-CAPCAMEL-;
  window.global.functions.render(
    <Page>
        <-CAPCAMEL-/>
    </Page>
  );
};
