import React from 'react';
import ReactDom from 'react-dom';
import RoleFeatureAddModify from './roleFeatureModifyComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_roleFeatureModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleFeatureModify;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<RoleFeatureAddModify/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/roleFeatureList", Active:""}, {Title:window.appContent.ModifyRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesRoleFeatureModify} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <RoleFeatureAddModify/>
      </Page>
    );
  }
};
