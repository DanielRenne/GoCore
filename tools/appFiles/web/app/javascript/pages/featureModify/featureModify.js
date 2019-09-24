import React from 'react';
import ReactDom from 'react-dom';
import FeatureAddModify from './featureModifyComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_featureModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatureModify;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
    var page = (<FeatureAddModify/>);
    var breadCrumbs = [{Title:window.appContent.List, Link:"/#/featureList", Active:""}, {Title:window.appContent.ModifyRecord, Link:"", Active:""}];
    window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFeatureModify} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FeatureAddModify/>
      </Page>
    );
  }
};
