import React from 'react';
import ReactDom from 'react-dom';
import FileObjectList from './fileObjectListComponents';
import Page from '../theme';
import PlainPage from '../../components/plainPage'
import FontIcon from 'material-ui/FontIcon';

const render_mode = "plain_page";

window.Load_fileObjectList = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFileObjectList;
  if (render_mode == "settings") {
    window.Load_settings();
  } else if (render_mode == "plain_page") {
      var page = (<FileObjectList/>);
      var breadCrumbs = [];
      window.global.functions.render(
        <PlainPage page={page} title={appContent.TitlesFileObjectList} icon={
          <FontIcon className="large material-icons">business</FontIcon>
        } breadCrumbs={breadCrumbs}/>
      );
  } else {
    window.global.functions.render(
      <Page>
          <FileObjectList/>
      </Page>
    );
  }
};
