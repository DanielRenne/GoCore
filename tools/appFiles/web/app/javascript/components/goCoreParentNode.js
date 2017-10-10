import React, {Component} from 'react';
import BaseComponent from './base';
import ReactDOM from 'react-dom';

class GoCoreParentNode extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.state = {
      fillerStyle: {}
    }
  }

  componentDidMount() {
    window.setTimeout(() => {
      // give time to get entire JS page loaded for height to get the best possible calculation for our filler div
      let div = ReactDOM.findDOMNode(this.mainNode);
      if (div != undefined) {
        let dimensions = div.getBoundingClientRect();
        if (window.ScrollHeight > 0) {
          let height = window.ScrollHeight - dimensions.height;
          if (height > 0) {
            this.setComponentState({fillerStyle:{height: height}})
          }
        }
      }
    }, 1000);
  }

  render() {
    this.logRender();
    return <span>
      <div ref={(c) => this.mainNode = c} className={this.props.addPadding ? "start-react-page": null} style={this.props.addPadding ? this.globs.basePageStyle(): {}}>
        {this.props.children}
      </div>
      <div style={this.state.fillerStyle}/>
    </span>;
  }
}

GoCoreParentNode.defaultProps = {
  addPadding: true
};

export default GoCoreParentNode;
