import React, {Component} from "react";
import Snackbar from "material-ui/Snackbar";
import {deepOrange300, red300, red500, green300, grey900, green900} from "material-ui/styles/colors";
import Dialog from "material-ui/Dialog";
import FlatButton from "material-ui/FlatButton";
import {ConfirmPopup, TextField} from "../globals/forms";
import BaseComponent from "../components/base";
import ErrorNotification from "../components/errorNotification";
import ActionBugReport from "material-ui/svg-icons/action/bug-report";
import FooterUpdates from "./footerUpdates";
let snackBarOrange = deepOrange300;
let snackBarRed = red300;
let snackBarGreen = green300;

class Footer extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.getTimeZones = () => {

      window.api.post({action:"GetTimeZones", controller: "serverSettingsModify", leaveStateAlone:true, callback:(data) => {
        data.TimeZones.splice(0, 0, {Location:"", Country:""});
        this.setComponentState({TimeZones:data.TimeZones});
      }});
    };

    this.getTimeZoneMenuItems = () => {
      if (this.state.TimeZones === undefined || this.state.TimeZones === null) {
        return null;
      }
      return this.state.TimeZones.map((tz) => <option
        key={tz.Location}
        value={tz.Location}
        >{(tz.Location == "") ? "" : tz.Location + " (" + tz.Country + ")"}</option>)
    };

    this.handlePasswordChange = (event) => {
      this.setComponentState({
        Password: event.target.value,
        PasswordErrors: ""
      });
    };

    this.handleConfirmPasswordChange = (event) => {
      this.setComponentState({
        ConfirmPassword: event.target.value,
        ConfirmPasswordErrors: ""
      });
    };

    this.handleAllDialogClose = (event) => {
      this.setComponentState({
        DialogOpen2: false,
        DialogGenericOpen: false,
        DialogGenericMessage: "",
        DialogMessage2: "",
        DialogTitle2: "",
        DialogTitleGeneric: "",
        DialogOpen:false
      });
      window.appState.DialogOpen2 = false;
      window.appState.DialogGenericOpen = false;
      window.appState.DialogGenericMessage = "";
      window.appState.DialogMessage2 = "";
      window.appState.DialogTitle2 = "";
      window.appState.DialogTitleGeneric = "";
      window.PopupServerMessage = "";
      window.PopupClientMessage = "";
    };

    this.handleDialogSendBug2 = (event) => {
      if (window.appState.DialogMessage2.indexOf("Client Side") != -1) {
        var subject = "App%20(Client%20Side)Error";
      } else {
        var subject = "App%20Error";
      }
      var navigatorInfo = "";
      try {
        navigatorInfo = encodeURIComponent(JSON.stringify({UA: window.navigator.userAgent, product: window.navigator.product, productSub: window.navigator.productSub, platform: window.navigator.platform, language: window.navigator.language, appVersion: window.navigator.appVersion, cookieEnabled: window.navigator.cookieEnabled}, "", 3));
      } catch (e) {}
      document.location.href = "mailto:unknown@unknown.com?subject=" + subject + "&body=Hello%20I%20am%20experiencing%20an%20issue%20using%20your%20software.%20%20%3CPlease%20describe%20what%20you%20were%20doing%20to%20recreate%20the%20issue%3E%20%0A%0A%0A(Please%20do%20not%20edit%20below%20these%20lines%20of%20equal%20signs%20so%20we%20can%20solve%20the%20problem%20for%20you%20quickly)%0A%0A%0A=========================================================%0A=========================================================%0A%0A%0A%0AAccount: " + window.appState.AccountName + " (id: " + window.appState.AccountId + ")%0A%0AError%20Message%20Shown:%0A%0A---------------------------------%0A%0A" + window.appState.DialogMessage2.split("\n").join("%0A") + "%0A%0AClient%20Information:%0A%0A" + navigatorInfo + "%0A%0AApp%20State:%0A%0A" + encodeURIComponent(JSON.stringify(window.appState, "", 3)) + "%0A%0APage%20State:%0A%0A" + encodeURIComponent(JSON.stringify(window.pageState, "", 3));

      window.PopupClientMessage = window.PopupClientMessage + "\n\nBug has been sent!"
      setTimeout(()=> {
        window.appState.PopupErrorSubmit2 = false;
        window.PopupServerMessage = "";
        window.PopupClientMessage = "";
        this.handleAllDialogClose();
        }, 2000
      )
    };

    this.handleDialogOpen = (event) =>  {
      this.setComponentState({DialogOpen: true});
    };

    this.handleDialogOpen2 = (event) =>  {
      this.setComponentState({DialogOpen2: true});
    };

    this.handleGenericDialogOpen = (event) =>  {
      this.setComponentState({DialogGenericOpen: true});
    };

    window.goCore.handleGenericDialogOpen = this.handleGenericDialogOpen;
    window.goCore.handleDialogOpen2 = this.handleDialogOpen2;
    window.goCore.handleDialogClose2 = this.handleAllDialogClose;
    window.goCore.handleDialogOpen = this.handleDialogOpen;
    window.goCore.handleDialogClose = this.handleAllDialogClose;

    this.handleActionTouchTap = (event) => {
      this.setComponentState({SnackbarOpen: false});
      window.appState.SnackbarOpen = false;
      if (this.state.ShowDialogSubmitBug2) {
        window.api.post({action: "RollBackFromSnackbar", state: this.state, controller:"transactions"});
      }
    };

    this.handleRequestClose = (event) => {
      this.setComponentState({SnackbarOpen: false});
      window.appState.SnackbarOpen = false;

      setTimeout(()=> {
        this.setComponentState({
          SnackbarMessage: "",
          SnackbarType: ""
        });
        window.appState.SnackbarMessage = "";
        window.appState.SnackbarType = "";
      }, 2000);
    };

    this.handleErrorView = () => {
      window.launcher.showDialog();
    };



    this.SubmitPassword = () => {
      var changePassword = (callback, userId=null) => {
        if (userId == null) {
          userId = window.appState.UserId;
        }
        window.api.post({action: "UserEnforcePasswordChange", state: {Password: this.state.Password, ConfirmPassword: this.state.ConfirmPassword, User: {Id: userId}}, controller:"userModify", leaveStateAlone: true, callback: (vm, message, trace, messagetype) => {
          if (vm.ConfirmPasswordErrors == "" && vm.PasswordErrors == "" && messagetype != "Error") {
            window.appState.UserEnforcePasswordChange = false;
            this.setComponentState({UserEnforcePasswordChange: false}, callback)
          } else {
            if (messagetype == "Error" && vm.ConfirmPasswordErrors == "" && vm.PasswordErrors == "") {
              vm.ConfirmPasswordErrors = "An unknown error occurred";
              vm.PasswordErrors = "An unknown error occurred";
            }
            this.setComponentState({ConfirmPasswordErrors: vm.ConfirmPasswordErrors, PasswordErrors: vm.PasswordErrors})
          }
        }});
      };
      changePassword();
    }

    this.state = window.appState;
    window.goCore.setFooterStateFromExternal = (state) => {
      window.appState = state;
      if (state.HTTPPort > 0) {
        window.launcher.ws.port = state.HTTPPort;
      }
      this.setComponentState(window.appState);
    };

    this.resizeEvent = (e) => this.handleResize(e);
  }

  handleResize(e) {
    var width = $(window).width();
    var height = $(window).height();
    if (width != this.state.windowWidth && height != this.state.windowHeight) {
      var changes = {windowWidth:  width, windowHeight: height};
      if (this.state.loggedIn) {
        if (width < 768) {
          $('body').css({paddingTop: 66});
        } else {
          $('body').css({paddingTop: 70});
        }
        $('.site-navbar').removeClass('hide');
      }
      this.setComponentState(changes);
    }
  }

  componentDidMount() {
    this.handleResize();
    window.addEventListener('resize', this.resizeEvent);
    this.getTimeZones();
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeEvent);
  }

  render() {
    this.logRender();
    var bkgColor = snackBarGreen;
    var action = "undo";

    var pass8Chars = red500;
    var passInteger = red500;
    var passLower = red500;
    var passUpper = red500;
    var passSpecial = red500;
    var passMatch = red500;

    if (this.state.SnackBarUndoTransactionId == "") {
      action = "";
    }
    if (this.state.SnackbarType == "Warning"){
      bkgColor = snackBarOrange;
    }
    if (this.state.SnackbarType == "Error"){
      bkgColor = snackBarRed;
      action = ""
    }
    var dialogButtons = [];

    if (this.state.PopupErrorSubmit2 && this.state.ShowDialogSubmitBug2) {
      dialogButtons.push(<FlatButton
              label={window.appContent.FooterSubmitBug}
              primary={false}
              icon={<ActionBugReport />}
              keyboardFocused={false}
              onTouchTap={this.handleDialogSendBug2}
            />);
    }

    dialogButtons.push(
            <FlatButton
              label={window.appContent["OK"]}
              primary={true}
              keyboardFocused={true}
              onTouchTap={this.handleAllDialogClose}
            />);
    return (
      <div className={(this.state.windowWidth > 768) ? "Footer": null}>
        {this.state.SnackbarOpen && this.state.SnackbarMessage != "" ? <Snackbar
            open={this.state.SnackbarOpen && this.state.SnackbarMessage != ""}
            message={this.state.SnackbarMessage}
            action={action}
            autoHideDuration={this.state.SnackbarAutoHideDuration}
            onActionTouchTap={this.handleActionTouchTap}
            onRequestClose={this.handleRequestClose}
            onClick={this.handleRequestClose}
            bodyStyle={{backgroundColor:bkgColor, marginBottom: 34, action:{color: 'blue'} }}
        />: null}

        {this.state.DialogGenericOpen ? <Dialog
          key="dialog-generic"
          style={{zIndex: 4000}}
          title={window.appContent.Alert}
          actions={<FlatButton
              label={window.appContent["OK"]}
              primary={true}
              keyboardFocused={true}
              onTouchTap={this.handleAllDialogClose}
            />}
          modal={false}
          open={this.state.DialogGenericOpen}
          onRequestClose={this.handleAllDialogClose}
          autoScrollBodyContent={true}
        >
          <div style={{marginTop: '1px'}}>
              <pre className="popup">
                {this.state.DialogGenericMessage}
              </pre>
          </div>
        </Dialog>: null}

        {this.state.DialogOpen2 ? <Dialog
          key="dialog-2"
          style={{zIndex: 4000}}
          title={this.state.DialogTitle2}
          actions={dialogButtons}
          modal={false}
          open={this.state.DialogOpen2}
          onRequestClose={this.handleAllDialogClose}
          autoScrollBodyContent={true}
        >
          <div style={{marginTop: '1px'}}>
            {typeof(this.state.DialogMessage2) == "string" ? <pre className="popup">
                {this.state.DialogMessage2}
            </pre>: this.state.DialogMessage2}
          </div>
        </Dialog>: null}
        {(window.appState.loggedIn && this.state.UserEnforcePasswordChange) ?
            <ConfirmPopup
              open={true}
              autoClose={false}
              showActionButtons={true}
              showActionSubmit={true}
              showActionCancel={false}
              onSubmit={this.SubmitPassword}
              title={window.appContent.ForcedToUpdatePassword}
              popupHTML={
                <div>
                  <div>
                    <span className ="AlignerLeft">
                      <span>
                        <TextField
                          floatingLabelText={"* " + window.appContent.UserAddEditPassword}
                          hintText={"* " + window.appContent.UserAddEditPassword}
                          fullWidth={false}
                          type="password"
                          onChange={this.handlePasswordChange}
                          errorText={this.globs.translate(this.state.PasswordErrors)}
                        />
                        <br />
                        <TextField
                          floatingLabelText={"* " + window.appContent.UserAddEditConfPassword}
                          hintText={"* " + window.appContent.UserAddEditConfPassword}
                          fullWidth={false}
                          type="password"
                          onChange={this.handleConfirmPasswordChange}
                          errorText={this.globs.translate(this.state.ConfirmPasswordErrors)}
                        />
                      </span>
                      <span style={{marginLeft: 100}}>
                        <label style={{color: grey900}}>Password criteria:</label><br />
                        <label style={{color: pass8Chars, marginLeft: 50}}>8 Characters</label><br />
                        <label style={{color: passInteger, marginLeft: 50}}>1 Integer</label><br />
                        <label style={{color: passLower, marginLeft: 50}}>1 Lowercase</label><br />
                        <label style={{color: passUpper, marginLeft: 50}}>1 Uppercase</label><br />
                        <label style={{color: passSpecial, marginLeft: 50}}>1 Special Character</label><br />
                        <label style={{color: passMatch, marginLeft: 50}}>{this.state.Password != "" ? this.state.Password === this.state.ConfirmPassword ? "Passwords Match!" : "Passwords do not match.":""}</label><br />
                      </span>
                    </span>
                  </div>
              </div>
            }
            />
            : null}
        {this.state.DialogOpen && this.state.ShowDialogSubmitBug2 ?
          <ErrorNotification
            open={this.state.DialogOpen}
            onView={this.handleErrorView}
          />: null}
        {(!window.appState.InRoomControl && this.state.windowWidth > 768) ? <span className="Footer-social">
          Copyright &#169;{this.state.CopyrightYear} XXXXXXXX  | Version: {(window.appState.DeveloperMode) ? <a href="javascript:" onClick={() => {
            this.globs.PopupWindow(window.atob(window.pageState.DeveloperLog))
        }}>{window.appState.displayVersion + " & View Logs"}</a>: window.appState.displayVersion}

        </span>: null}
      </div>
    );
  }
}

export default Footer;
