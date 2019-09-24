/**
 * Created by Dan on 10/26/16.
 */
import {React, BaseComponent} from "../globals/forms";
import Drawer from "material-ui/Drawer";


class InfoNotification extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.onView = this.props.onView;

    this.state = {
      open: this.props.open,
      drawerWidth: this.props.width
    };

    this.close = () => {
      this.setComponentState({open: false}, () => this.postClosed());
    };

    this.handleViewError =() => {
      if (this.onView != undefined) {
        this.onView();
      }
    }
  }

  postClosed() {
    window.api.post({action: "NotificationClosed", state: {UniqueId: this.props.uniqueId, Type: this.props.type}, controller:"notifications"});
  }

  componentWillReceiveProps(nextProps) {
    if (this.state.open === true){
      return;
    }

    if (nextProps.open === true) {
      this.clearTimeout = setTimeout(this.close, 180000);
    }

    window.CurrentNotificationUniqueId = nextProps.uniqueId;
    window.CurrentNotificationType = nextProps.type;

    this.setComponentState({open:nextProps.open});
  }

  render() {
    try {
      this.logRender();
      if (!this.state.open) {
        return null
      }
      let line1Html = (this.props.IPAddress.length > 0) ? <span>{this.__(window.appContent.InfoMessageFromIP, {ipaddr: this.props.IPAddress})}</span>: this.props.User.hasOwnProperty("First") ? <span>{this.__(window.appContent.InfoMessageFromUser, {f_nm: this.props.User.First, l_nm: this.props.User.Last})}</span> : null;
      let line1 = (this.props.IPAddress.length > 0) ? this.__(window.appContent.InfoMessageFromIP, {ipaddr: this.props.IPAddress}): this.props.User.hasOwnProperty("First") ? this.__(window.appContent.InfoMessageFromUser, {f_nm: this.props.User.First, l_nm: this.props.User.Last}) : null;
      let line2 = window.appContent[this.props.title];
      let lineCount1 = Math.ceil(line1.length / 60);
      let lineCount2 = Math.ceil(line2.length / 60);
      let heightAdjustments = lineCount1 + lineCount2 - 2;
      let height = 103;
      if (heightAdjustments > 0) {
        height = height + (25*heightAdjustments);
      }
      session_functions.Dump(lineCount1, lineCount2);

      return (
        <Drawer width={this.state.drawerWidth} open={this.state.open} openSecondary={true} containerStyle={{zIndex:2000, top:71, height:height+27,backgroundColor:"#c8e6c9"}} >
          <div className="alert alert-icon alert-success alert-dismissible" role="alert" style={{height:height}}>
          <button type="button" className="close"  aria-label="Close" onClick={this.close}>
            <span>×</span>
          </button>
          <i className="icon md-notifications" aria-hidden="true"></i>
          <div ref={(c) => this.div1 = c} style={{color:"black"}}>{line1Html}</div>
          <div ref={(c) => this.div2 = c} style={{color:"black", marginTop:15}}>{line2}</div>
          <p ref={(c) => this.div3 = c} className="margin-top-15">
            <a href={this.props.link} onClick={() => {this.setComponentState({open:false}, () => {
              clearTimeout(this.clearTimeout)
            })}}>{this.props.linkTitle}</a>
          </p>

          </div>

        </Drawer>
      );
    } catch(e) {
      return <Drawer width={this.state.drawerWidth} open={this.state.open} openSecondary={true} containerStyle={{zIndex:2000, top:71, height:height+27,backgroundColor:"#c8e6c9"}} >{this.globs.ComponentError("InfoNotification", e.message, e)}</Drawer>
    }
  }
}

InfoNotification.propTypes = {
    width: React.PropTypes.number,
    title: React.PropTypes.string,
    link: React.PropTypes.string,
    linkTitle: React.PropTypes.string,
    uniqueId: React.PropTypes.string,
    type: React.PropTypes.string,
    open: React.PropTypes.bool
};

InfoNotification.defaultProps = {
    width: 500,
    open:false
};

export default InfoNotification;
