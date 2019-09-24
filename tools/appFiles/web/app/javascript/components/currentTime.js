import React, {Component} from 'react';
import BaseComponent from './base';

class CurrentTime extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = {
      CurrentDate:"",
      CurrentTime:""
    }

    this.handleClockUpdate = (data) => {
      this.setComponentState({CurrentTime:data.Time, CurrentDate:data.Date});
    };
  }

  componentDidMount() {
    this.clockUpdateCallbackId = window.api.registerSocketCallback((data) => {this.handleClockUpdate(data)}, "Clock");
  }

  componentWillUnmount() {
    window.api.unRegisterSocketCallback(this.clockUpdateCallbackId);
  }

  render() {
    try {
       return (
         <div>
           <span className="AlignerLeft">
           <div>{"Current Date:"}</div>
           <div style={{marginLeft:10, fontWeight:"bold"}}>{this.state.CurrentDate}</div>
           <div style={{marginLeft:10}}>{"Current Time:"}</div>
           <div style={{marginLeft:10, fontWeight:"bold"}}>{this.state.CurrentTime}</div>
           </span>
         </div>
       )
    } catch(e) {
      return this.globs.ComponentError("CurrentTime", e.message, e);
    }
  }
}

export default CurrentTime;
