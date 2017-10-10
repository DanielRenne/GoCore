import React from 'react';
import ReactDom from 'react-dom';
import AppErrorAdd from './appErrorAddComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_appErrorAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesAppErrorAdd;

  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<AppErrorAdd/>);
      var breadCrumbs = [{Title:"Todo.... add my links", Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesAppErrorAdd} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <AppErrorAdd/>
      </Page>
    );
  }
};
