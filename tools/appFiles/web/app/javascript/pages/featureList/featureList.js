import React from 'react';
import ReactDom from 'react-dom';
import FeatureList from './featureListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_featureList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFeatureList;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FeatureList/>);
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFeatureList} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        }/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FeatureList/>
      </Page>
    );
  }
};
