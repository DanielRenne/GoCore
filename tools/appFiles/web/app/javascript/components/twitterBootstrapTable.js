import React, {Component} from 'react';
import BaseComponent from './base';

// This class chunks out all your output into twitter bootstrap col-md-XX column sizes

class TwitterBootstrapTable extends BaseComponent {
  constructor(props, context) {
    super(props, context);
  }

  render() {
    try {
      this.logRender();
      if (this.props.data.length > 0) {
        var ret = [];
        var colSize = (12 / this.props.columnsPerRow);
        if (colSize % 1 != 0) {
           console.error("Incorrect columnsPerRow.  Needs to divide by 12");
        }
        if (this.props.data.length == 1) {
          try {
            //blow up on anything that is not a card
            this.props.data[0].props.header
            this.props.data[0].props.figureImage
          } catch(e) {
            colSize = 12;
          }
        }
        var colMdClass = "col-md-" + colSize;
        $.each(this.props.data.chunk(this.props.columnsPerRow), (k, o) => {
          var cols = [];
          $.each(o, (kr, r) => {
            cols.push((<div key={"card-col-" + k + kr} className={colMdClass}>{r}</div>));
          });

          ret.push((<div key={"card" + k} className="row" style={this.props.rowStyle}>{cols}</div>));
        });
        return (
          <div onClick={this.props.onClick}>
            {ret}
          </div>
        );
      } else {
        return null;
      }
    } catch(e) {
      return this.globs.ComponentError("TwitterBootstrapTable", e.message, e);
    }
  }
}

TwitterBootstrapTable.propTypes = {
  data: React.PropTypes.array.isRequired,
  rowStyle: React.PropTypes.object,
  columnsPerRow: React.PropTypes.number // must be 6,4,3,2
};

TwitterBootstrapTable.defaultProps = {
  columnsPerRow: 4,
  rowStyle: {}
};


export default TwitterBootstrapTable;


