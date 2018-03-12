import React from 'react';
import ReactDom from 'react-dom';
import RoleFeatureAdd from './roleFeatureAddComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_roleFeatureAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleFeatureAdd;

  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<RoleFeatureAdd/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/roleFeatureList", Active:""}, {Title:window.appContent.AddRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesRoleFeatureAdd} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <RoleFeatureAdd/>
      </Page>
    );
  }
};
