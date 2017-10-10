import React from 'react';
import ReactDom from 'react-dom';
import FeatureAdd from './featureAddComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_featureAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatureAdd;

  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FeatureAdd/>);
      var breadCrumbs = [{Title:"Back to list", Link:"/#/featureList", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFeatureAdd} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FeatureAdd/>
      </Page>
    );
  }
};
