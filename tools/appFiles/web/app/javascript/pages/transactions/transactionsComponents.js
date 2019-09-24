import React, {Component} from 'react';
import BasePageComponent from '../../components/basePageComponent';

class Transactions extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    this.logRender();
    return (
        <div>{appContent.helloWorld} {this.__(pageContent.welcomeExample, {page:pageContent.objectName})}</div>
    );
  }
}

export default Transactions;
