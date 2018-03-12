import React from 'react';
import ReactDom from 'react-dom';
import RoleAdd from './roleAddComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "settings";

window.Load_roleAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleAdd;

  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<RoleAdd/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/roleList", Active:""}, {Title:window.appContent.AddRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesRoleAdd} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <RoleAdd/>
      </Page>
    );
  }
};
