import React from 'react';
import ReactDom from 'react-dom';
import FeatureGroupList from './featureGroupListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_featureGroupList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatureGroupList;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FeatureGroupList/>);
      var breadCrumbs = [];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFeatureGroupList} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FeatureGroupList/>
      </Page>
    );
  }
};
