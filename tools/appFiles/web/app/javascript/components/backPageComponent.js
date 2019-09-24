import React, {Component} from 'react';
import ActionButton from './actionButton'
import BaseComponent from './base';

export default function BackPage(Component, ShowBack){

  class BackPageComponent extends Component {
    constructor(props, context) {
      super(props, context);
    }

    render(){
      try {
        return (
          <div>
            <Component/>
            <ActionButton visible={ShowBack} tooltip={window.appContent.backPageComponentBack} type="Back"/>
          </div>
        );
      } catch(e) {
        return this.globs.ComponentError("BackPageComponent", e.message, e);
      }
    }
  }

  return BackPageComponent;
}
