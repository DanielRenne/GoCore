import React from 'react';
import ReactDom from 'react-dom';
import -CAPCAMEL-Add from './-CAMEL-AddComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "-PAGE_RENDER_MODE-";

window.Load_-CAMEL-Add = function() {
  document.title = window.global.functions.productTitle() + window.appContent.Titles-CAPCAMEL-Add;

  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<-CAPCAMEL-Add/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/-CAMEL-List", Active:""}, {Title:window.appContent.AddRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.Titles-CAPCAMEL-Add} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <-CAPCAMEL-Add/>
      </Page>
    );
  }
};
