import React from 'react';
import ReactDom from 'react-dom';
import RoleFeatureList from './roleFeatureListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_roleFeatureList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesRoleFeatureList;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<RoleFeatureList/>);
      var breadCrumbs = [];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesRoleFeatureList} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <RoleFeatureList/>
      </Page>
    );
  }
};
