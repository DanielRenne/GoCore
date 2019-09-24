import {React,
  BaseComponent,
  RaisedButton} from '../globals/forms';
import Dialog from 'material-ui/Dialog';

class PagePopup extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = {
      open: false
    }

    this.handleClose = () => {
      this.setComponentState({open: false});
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
          label={window.appContent.DeviceDetailsPopupClose}
          primary={true}
          onTouchTap={this.handleClose}
          style={{marginRight: 10}}
        />
      ];

      return (
        <div>
          <Dialog
            title={this.props.title}
            actions={actions}
            open={this.state.open}
            onRequestClose={this.handleClose}
          >
            {this.props.children}
          </Dialog>
        </div>
      );
    } catch(e) {
      return  <Dialog
            title={this.props.title}
            open={this.state.open}
            onRequestClose={this.handleClose}
          >
            {this.globs.ComponentError("PagePopup", e.message, e)}
          </Dialog>
    }
  }
}

PagePopup.propTypes = {
  open: React.PropTypes.bool,
  title: React.PropTypes.string
};

export default PagePopup;
