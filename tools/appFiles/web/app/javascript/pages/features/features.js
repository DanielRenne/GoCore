import React from 'react';
import ReactDom from 'react-dom';
import Features from './featuresComponents';
import Page from '../theme';

window.Load_features = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatures;
  window.global.functions.render(
    <Page>
        <Features/>
    </Page>
  );
};
