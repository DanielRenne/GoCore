import React from 'react';
import ReactDom from 'react-dom';
import RoleAddModify from './roleModifyComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "settings";

window.Load_roleModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleModify;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<RoleAddModify/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/roleList", Active:""}, {Title:window.appContent.ModifyRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesRoleModify} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <RoleAddModify/>
      </Page>
    );
  }
};
