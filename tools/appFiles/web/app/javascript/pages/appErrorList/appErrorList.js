import React from 'react';
import ReactDom from 'react-dom';
import AppErrorList from './appErrorListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_appErrorList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesAppErrorList;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<AppErrorList/>);
      var breadCrumbs = [{Title:"Todo.... add my links", Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesAppErrorList} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <AppErrorList/>
      </Page>
    );
  }
};
