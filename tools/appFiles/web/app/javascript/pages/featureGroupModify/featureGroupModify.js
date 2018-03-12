import React from 'react';
import ReactDom from 'react-dom';
import FeatureGroupAddModify from './featureGroupModifyComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_featureGroupModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatureGroupModify;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FeatureGroupAddModify/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/featureGroupList", Active:""}, {Title:window.appContent.ModifyRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFeatureGroupModify} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FeatureGroupAddModify/>
      </Page>
    );
  }
};
