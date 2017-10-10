import React from 'react';
import ReactDom from 'react-dom';
import FileObjectAdd from './fileObjectAddComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_fileObjectAdd = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFileObjectAdd;

  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FileObjectAdd/>);
      var breadCrumbs = [{Title:window.appContent.List, Link:"/#/fileObjectList", Active:""}, {Title:window.appContent.AddRecord, Link:"", Active:""}];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFileObjectAdd} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FileObjectAdd/>
      </Page>
    );
  }
};
