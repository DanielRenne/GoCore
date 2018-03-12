import React from 'react';
import ReactDom from 'react-dom';
import FeatureGroups from './featureGroupsComponents';
import Page from '../theme';

window.Load_featureGroups = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatureGroups;
  window.global.functions.render(
    <Page>
        <FeatureGroups/>
    </Page>
  );
};
