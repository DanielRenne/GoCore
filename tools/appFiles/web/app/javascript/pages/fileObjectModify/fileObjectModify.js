import React from 'react';
import ReactDom from 'react-dom';
import FileObjectAddModify from './fileObjectModifyComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_fileObjectModify = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFileObjectModify;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FileObjectAddModify/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/fileObjectList", Active:""}, {Title:window.appContent.ModifyRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFileObjectModify} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FileObjectAddModify/>
      </Page>
    );
  }
};
