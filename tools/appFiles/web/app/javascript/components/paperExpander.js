import React from 'react';
import Paper from 'material-ui/Paper';
import Subheader from 'material-ui/Subheader';
import BaseComponent from './base';
import IconButton from 'material-ui/IconButton';
import {UpArrow, DownArrow} from '../globals/icons';
import {blueGrey500} from 'material-ui/styles/colors';

class PaperExpander extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = {
      expanded: this.props.expanded,
      style:this.props.style
    };

    this.handleToggle = () => {
      this.setComponentState({expanded:(this.state.expanded) ? false : true});
    }

    this.openPaper = () => {
      this.setComponentState({expanded:true});
    }
  }

  componentWillReceiveProps(nextProps) {
      this.setComponentState({style:nextProps.style, expanded: nextProps.expanded});
  }


  render() {
    try {
      this.logRender();
      var paperHeight = 10000;
      var expandedIcon = (<DownArrow  color={blueGrey500}/>);
      if (!this.state.expanded) {
        paperHeight = 50;
        expandedIcon = (<UpArrow  color={blueGrey500}/>);
      }
      var styleComposition = this.state.style;
      styleComposition.maxHeight = paperHeight;
      return (
          <Paper style={styleComposition}  zDepth={2} >
            <span className="AlignerRight">
              <Subheader style={{color:"black"}} className="paperExpanderSubHeader" >{this.props.title}</Subheader>
              <span className="AlignerLeft" style={{whiteSpace:"nowrap"}}>
                {this.props.button}
                <IconButton style={{padding:0, height:40,marginTop:5}}  onClick={(e) => this.handleToggle()}>{expandedIcon}</IconButton>
              </span>
            </span>
            <div>
              {this.state.expanded ? this.props.children: null}
            </div>
          </Paper>
      );
    } catch(e) {
      return this.globs.ComponentError("PaperExpander", e.message, e);
    }
  }
}



PaperExpander.propTypes = {
  expanded:React.PropTypes.bool,
  title:React.PropTypes.string,
  style:React.PropTypes.object,
  button:React.PropTypes.object
};

PaperExpander.defaultProps = {
    expanded: true,
    style:{}
};

export default PaperExpander;
