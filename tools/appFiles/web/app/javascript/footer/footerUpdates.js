/**
 * Created by Dan on 4/14/17.
 */
import {
  React,
  BaseComponent
} from '../globals/forms';
import {LinkIcon, ReportProblemIcon} from "../globals/icons";
import {deepOrange300,red300} from 'material-ui/styles/colors';



class FooterUpdates extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    this.timeOutId = 0;
    this.state = {};
  }

  componentDidMount() {
    this.backupDetailsId = window.api.registerSocketCallback((data) => {this.handleBackupDetails(data)}, "BackupDetails");
  }

  componentWillUnmount() {
    window.api.unRegisterSocketCallback(this.backupDetailsId);
  }


  handleBackupDetails(data) {
    this.setComponentState({Details:data.Details});
    clearTimeout(this.timeOutId);
    this.timeOutId = setTimeout(() => {
      this.setComponentState({Details:""});
    }, 10000)
  }


  render() {
    this.logRender();

    var icon = (<LinkIcon color="white" style={{marginTop:0, marginLeft:20, height:16}}/>);
    if (this.props.title.indexOf("Primary") == -1 && this.props.title.indexOf("Standby") == -1) {
      icon = (<div style={{marginTop:0, marginLeft:10}}></div>);
    }
    if (this.props.title.indexOf("Inactive") != -1) {
      icon = (<ReportProblemIcon color={red300} style={{marginTop:0, marginLeft:20, height:16}}/>);
    }

    var backupTitle = (<span className="AlignerRight" style={{marginBottom:4}}>

                        <h5 className="Aligner-item" style={{display:"block",  whiteSpace: "nowrap", marginLeft:20, marginBottom:0, marginTop:0, color:this.props.color}}>
                          {window.appContent[this.props.title]}
                        </h5>
                        {icon}
                      </span>);

    var backupDetails = (<span className="AlignerRight" style={{marginBottom:4, display:"-webkit-inline-box"}}>
                            <div style={{marginTop:0, marginLeft:10}}></div>
                            <h6 style={{display:"block",  whiteSpace: "nowrap", marginTop:0, marginRight:5,marginBottom:0, color:"white"}}>{this.state.Details}
                            </h6>
                          </span>);

    return (
      <div className="AlignerRight" style={{display:"-webkit-inline-box", width:300}}>
        {backupTitle}
        {backupDetails}
      </div>
    );
  }

}

FooterUpdates.propTypes = {
  title:React.PropTypes.string,
  color:React.PropTypes.string
};

export default FooterUpdates;
