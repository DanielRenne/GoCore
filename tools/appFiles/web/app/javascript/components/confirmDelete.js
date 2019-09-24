import React from 'react';
import Dialog from 'material-ui/Dialog';
import RaisedButton from 'material-ui/RaisedButton';
import FlatButton from 'material-ui/FlatButton';
import BaseComponent from './base';

class ConfirmDelete extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    //this.eventHandlerExample = (event) => {
      //window.api.post({action: "Warning", state: this.state});
    //};

    this.state = {
      open: false
    };

    this.handleOpen = () => {
      this.setComponentState({open: true});
    };

    this.handleClose = () => {
      this.setComponentState({open: false});
    };

  }

  render() {
    try {
      this.logRender();
      // if (!this.state.open) {
      //   return null;
      // }
      const actions = [
        <RaisedButton
          label={this.props.noCancel}
          primary={true}
          onTouchTap={this.handleClose}
          style={{marginRight: 10}}
        />,
        <RaisedButton
          label={this.props.yesDelete}
          secondary={true}
          onTouchTap={this.props.deleteFunction}
        />,
      ];

      return (
        <div>
          <RaisedButton
            label={this.props.buttonTriggerLabel}
            onTouchTap={this.handleOpen}
            secondary={true}
          />
            <Dialog
              title={this.props.dialogTitle}
              actions={actions}
              modal={false}
              open={this.state.open}
              onRequestClose={this.handleClose}
            >
            {this.props.dialogMessage}
          </Dialog>
        </div>
      );
    } catch(e) {
      return <Dialog
              title={this.props.dialogTitle}
              modal={false}
              open={this.state.open}
              onRequestClose={this.handleClose}
            >
            {this.globs.ComponentError("ConfirmDelete", e.message, e)}
          </Dialog>
    }
  }
}

ConfirmDelete.propTypes = {
  noCancel: React.PropTypes.string.isRequired,
  yesDelete: React.PropTypes.string.isRequired,
  buttonTriggerLabel: React.PropTypes.string.isRequired,
  dialogTitle: React.PropTypes.string.isRequired,
  dialogMessage: React.PropTypes.string.isRequired,
};

export default ConfirmDelete;
