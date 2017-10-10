import React from 'react';
import ReactDom from 'react-dom';
import FileObjects from './fileObjectsComponents';
import Page from '../theme';

window.Load_fileObjects = function() {
  document.title = window.global.functions.productTitle() + window.appContent.TitlesFileObjects;
  window.global.functions.render(
    <Page>
        <FileObjects/>
    </Page>
  );
};
