import React, {Component} from 'react';
import BasePageComponent from '../../components/basePageComponent';

class Roles extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
  }

  componentWillUpdate(nextProps, nextState) {
    return true;
  }

  render() {
    this.logRender();
    return (
        <div>{appContent.helloWorld} {this.__(pageContent.welcomeExample, {page:pageContent.objectName})}</div>
    );
  }
}

export default Roles;
