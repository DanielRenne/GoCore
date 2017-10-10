import React from 'react';
import ReactDom from 'react-dom';
import -CAPCAMEL-AddModify from './-CAMEL-ModifyComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "-PAGE_RENDER_MODE-";

window.Load_-CAMEL-Modify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.Titles-CAPCAMEL-Modify;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<-CAPCAMEL-AddModify/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/-CAMEL-List", Active:""}, {Title:window.appContent.ModifyRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.Titles-CAPCAMEL-Modify} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <-CAPCAMEL-AddModify/>
      </Page>
    );
  }
};
