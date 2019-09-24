import React from 'react';
import ReactDom from 'react-dom';
import BaseComponent from "./base";

class CenterAnything extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.state = {};
    this.state.DisplayHeight = window.innerHeight;
    this.state.DisplayWidth = window.innerWidth;
    this.resizeEventCenter = (e) => this.handleResizeCenter(e);
  }

  handleResizeCenter(e) {
    if (window.innerWidth != this.state.DisplayWidth || window.innerHeight != this.state.DisplayHeight) {
      this.setComponentState({
        DisplayHeight: window.innerHeight - window.HeaderHeight - window.FooterHeight,
        DisplayWidth: window.innerWidth,
      });
    }
  }

  componentDidMount() {
    window.addEventListener('resize', this.resizeEventCenter);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeEventCenter);
  }

  render() {
    return <div className="Aligner" style={{width: this.state.DisplayWidth, height: this.state.DisplayHeight}}>
      <div className="Aligner-item Aligner-item--top"></div>
      <div className="Align-center">
        {this.props.children}
      </div>
      <div className="Aligner-item Aligner-item--bottom"></div>
    </div>
  }
}

export default CenterAnything;
