/**
 * Created by Dan on 10/26/16.
 */
import {
  React,
  BaseComponent} from '../globals/forms';
import Drawer from 'material-ui/Drawer';



class ErrorNotification extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.onView = this.props.onView;

    this.state = {
      open: this.props.open,
      drawerWidth: this.props.width,
      title: this.props.title
    };

    this.close = () => {
      this.setComponentState({open: false});
    };

    this.handleViewError =() => {
      if (this.onView != undefined) {
        this.onView();
      }
    }

  }


  componentWillReceiveProps(nextProps) {
    if (this.state.open === true){
      return;
    }
    if (nextProps.open === true) {
      setTimeout(() => {
        this.setComponentState({open: false});
      },5000);
    }
      this.setComponentState({open:nextProps.open, title: nextProps.title});
  }

  render() {
    try {
      this.logRender();


      return (
        <Drawer width={this.state.drawerWidth} open={this.state.open} openSecondary={true} containerStyle={{zIndex:2000, top:71, height:103}} >
          <div className="alert alert-icon alert-danger alert-dismissible" role="alert" style={{height:103}}>
          <button type="button" className="close"  aria-label="Close" onClick={this.close}>
            <span>×</span>
          </button>
          <i className="icon md-notifications" aria-hidden="true"></i> {window.appContent.DialogError}
          <p className="margin-top-15">
            <button className="btn btn-danger" type="button" onClick={this.handleViewError}>{window.appContent.ViewError}</button>
          </p>

          </div>

        </Drawer>
      );
    } catch(e) {
      return <Drawer width={this.state.drawerWidth} open={this.state.open} openSecondary={true} containerStyle={{zIndex:2000, top:71, height:103}}>{this.globs.ComponentError("ErrorNotification", e.message, e)}</Drawer>
    }
  }
}

ErrorNotification.propTypes = {
    width:React.PropTypes.number,
    open: React.PropTypes.bool,
    onView: React.PropTypes.func
};

ErrorNotification.defaultProps = {
    width: 415,
    open:false
};

export default ErrorNotification;
