import React, {Component} from 'react';
import BasePageComponent from '../../components/basePageComponent';
class Accounts extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
        <div>DO NOT USE THIS UNTIL WE BUILD OUT THE LANDING PAGE FOR DRILLING INTO SITES{appContent.helloWorld} {this.__(pageContent.welcomeExample, {page:pageContent.objectName})}</div>
    );
  }
}

export default Accounts;
