import React, {Component} from 'react';
import ActionButton from './actionButton'
import BaseComponent from './base';

export default function AddRecordPage(Component, Action, Tooltip){

  class AddRecordPageComponent extends Component {
    constructor(props, context) {
      super(props, context);
    }

    render(){

      try {
        let props={};
        if (typeof(Action) == "function") {
          props.onClick = Action;
        } else {
          props.action = Action;
        }

        return (
          <div>
            <Component/>
            <ActionButton {...props} tooltip={window.pageContent[Tooltip]} type="Add"/>
          </div>

        );
      } catch(e) {
        return this.globs.ComponentError("AddRecordPageComponent", e.message, e);
      }
    }
  }

  return AddRecordPageComponent;
}
