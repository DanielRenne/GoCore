import {React, BaseComponent, RaisedButton} from "../globals/forms";
import Dialog from "material-ui/Dialog";

class InfoPopup extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = {
      open: false,
      message:""
    }

    this.open = () => {
      this.setComponentState({open: true});
    }

    this.setMessage = (message, open) => {
      this.setComponentState({message:message, open:open});
    }

    this.handleClose = () => {
      this.setComponentState({open: false}, () => {
        if (this.parent && this.props.parentStateKey) {
          let changes = {};
          changes[this.props.parentStateKey] = false;
          this.setParentState(changes, () => {
            if (this.props.hasOwnProperty("onClose")) {
              this.props.onClose();
            }
          });
        } else {
          if (this.props.hasOwnProperty("onClose")) {
            this.props.onClose();
          }
        }
      });
    }
  }

  componentWillReceiveProps(nextProps) {
      this.setComponentState({open:nextProps.open});
  }

  render() {
    try {
      this.logRender();

      const actions = [
        <RaisedButton
          label={window.appContent.ImportClose}
          primary={true}
          onTouchTap={this.handleClose}
          style={{marginRight: 10}}
        />
      ];

      return (
        <div>
          <Dialog
            title={window.pageContent.TitleImportUsers}
            actions={actions}
            open={this.state.open}
            onRequestClose={this.handleClose}
          >
            {(this.state.message != "") ? this.state.message : this.props.children}
          </Dialog>
        </div>
      );
    } catch(e) {
      return  <Dialog
            title={window.pageContent.TitleImportUsers}
            open={this.state.open}
            onRequestClose={this.handleClose}
          >
            {this.globs.ComponentError("InfoPopup", e.message, e)}
          </Dialog>
    }
  }
}

InfoPopup.propTypes = {
  parent: React.PropTypes.object,
  open: React.PropTypes.bool,
  onClose: React.PropTypes.func,
  parentStateKey: React.PropTypes.string,
};

export default InfoPopup;
