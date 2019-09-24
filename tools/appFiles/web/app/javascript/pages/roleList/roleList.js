import React from 'react';
import ReactDom from 'react-dom';
import RoleList from './roleListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "settings";

window.Load_roleList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleList;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<RoleList/>);
      var breadCrumbs = [];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesRoleList} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <RoleList/>
      </Page>
    );
  }
};
