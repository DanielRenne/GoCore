import React from 'react';
import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import BaseComponent from '../components/base';

class ConfirmPopup extends BaseComponent {
  constructor(props, context) {
    super(props, context);


    this.state = {
      title: window.appContent.ConfirmPopupConfirm,
      open: props.open,
      pageData: null,
      areYouSureMessage: props.areYouSureMsg,
      actionSubmitLabel: (props.actionSubmitLabel) ? props.actionSubmitLabel : window.appContent.ConfirmPopupSubmit,
      actionCancelLabel: (props.actionCancelLabel) ? props.actionCancelLabel : window.appContent.ConfirmPopupCancel
    };

    if (props.title) {
      this.state.title = props.title;
    }


    this.handleOpen = (customData) => {
      var updates = {};
      if (customData) {
        updates.pageData = customData;
      }

      if (this.state.areYouSureMessage) {
        if (Array.isArray(customData)) {
          updates.areYouSureMessage = window.global.functions._(this.state.areYouSureMessage, {total: customData.length});
        } else {
          updates.areYouSureMessage = window.global.functions._(this.state.areYouSureMessage);
        }
      }
      updates.open = true;
      this.setComponentState(updates);
    };

    this.handleComplete = () => {
      if (!this.props.onSubmitFullControl) {
        (this.props.onSubmit)(this.state.pageData);
        this.setComponentState({open: false});
      } else {
        let funcCall = () => {
          (this.props.onSubmit)(this.state.pageData);
        };
        funcCall();
      }
    };

    this.handleClose = () => {
      if (this.props.onClose) {
        this.props.onClose();
      }
      this.setComponentState({open: false});
    };

    this.handleCloseRules = (forceClose=true) => {
      if (this.props.autoClose || forceClose) {
        this.setComponentState({open: false});
        if (this.props.onClose) {
          this.props.onClose();
        }
      }
    };
  }

  componentWillReceiveProps(nextProps) {
    var changes = {};
    let isDifferent = false;
    if (nextProps.title != undefined && nextProps.title != this.state.title) {
      changes.title = nextProps.title;
      isDifferent = true;
    }
    if (nextProps.actionSubmitLabel != undefined && nextProps.actionSubmitLabel != this.state.actionSubmitLabel) {
      changes.actionSubmitLabel = nextProps.actionSubmitLabel;
      isDifferent = true;
    }
    if (nextProps.areYouSureMessage != undefined && nextProps.areYouSureMessage != this.state.areYouSureMessage) {
      changes.areYouSureMessage = nextProps.areYouSureMessage;
      isDifferent = true;
    }
    if (nextProps.actionCancelLabel != undefined && nextProps.actionCancelLabel != this.state.actionCancelLabel) {
      changes.actionCancelLabel = nextProps.actionCancelLabel;
      isDifferent = true;
    }
    if (nextProps.open != undefined && nextProps.open != this.state.open) {
      changes.open = nextProps.open;
      isDifferent = true;
    }
    if (isDifferent) {
      this.setComponentState(changes);
    }
  }

  render() {
    try {
      this.logRender();
      if (!this.state.open) {
        return null;
      }

      let actions = [];
      if (this.props.showActionButtons) {
        if (this.props.showActionCancel) {
          actions.push(
              <FlatButton
                  label={this.state.actionCancelLabel}
                  primary={true}
                  onTouchTap={() => {
                    this.handleClose(true);
                  }}
              />);
        }
        if (this.props.showActionSubmit) {
          actions.push(<FlatButton
            label={this.state.actionSubmitLabel}
            primary={true}
            onTouchTap={this.handleComplete}
          />)
        }
      }

      var areYouSureMessage = null;
      if (!this.props.popupHTML) {
        areYouSureMessage = (this.state.areYouSureMessage) ? this.state.areYouSureMessage: window.appContent.ConfirmPopupMessage
      }

      return (
        <div>
          <Dialog
            contentStyle={this.props.width ? {width:this.props.width} : {}}
            title={this.state.title}
            actions={actions}
            modal={false}
            secondary={true}
            open={this.state.open}
            onRequestClose={this.handleCloseRules}
            style={{zIndex:3000}}
          >
            {areYouSureMessage}
            {(this.props.popupHTML) ? this.props.popupHTML: null}
          </Dialog>
        </div>
      );
    } catch(e) {
      return this.globs.ComponentError(this.getClassName(), e.message);
    }
  }
}

ConfirmPopup.propTypes = {
  open: React.PropTypes.bool,
  width: React.PropTypes.number,
  title: React.PropTypes.string,
  onSubmitFullControl: React.PropTypes.bool, // dont let it close automatically and apply your own stuff
  onSubmit: React.PropTypes.func.isRequired,
  autoClose: React.PropTypes.bool.isRequired,
  showActionButtons: React.PropTypes.bool,
  showActionSubmit: React.PropTypes.bool,
  showActionCancel: React.PropTypes.bool,
  actionSubmitLabel: React.PropTypes.string,
  actionCancelLabel: React.PropTypes.string,
  onClose: React.PropTypes.func,
  areYouSureMsg: React.PropTypes.node,
  popupHTML: React.PropTypes.node,
};

ConfirmPopup.defaultProps = {
  onSubmitFullControl: false,
  showActionCancel: true,
  showActionSubmit: true,
  showActionButtons: true,
  autoClose: true
};

export default ConfirmPopup;
