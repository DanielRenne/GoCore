import React from 'react';
import ReactDom from 'react-dom';
import -CAPCAMEL-List from './-CAMEL-ListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "-PAGE_RENDER_MODE-";

window.Load_-CAMEL-List = function() {
  document.title = window.global.functions.productTitle() + window.appContent.Titles-CAPCAMEL-List;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<-CAPCAMEL-List/>);
      var breadCrumbs = [];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.Titles-CAPCAMEL-List} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <-CAPCAMEL-List/>
      </Page>
    );
  }
};
