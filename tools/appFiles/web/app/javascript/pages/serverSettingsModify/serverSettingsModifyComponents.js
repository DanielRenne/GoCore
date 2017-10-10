import {
  React,
  CenteredPaperGrid,
  BasePageComponent,
  TextField,
  RaisedButton,
  Toggle,
  SelectField,
  MenuItem,
  green600,
  RadioButtonGroup,
  RadioButton,
  ConfirmPopup,
  IconButton
} from "../../globals/forms";
import {
  AddIcon,
  Ethernet,
  CPUIcon,
  Security,
  Email,
  Settings,
  License,
  ToolsIcon,
  DeleteIcon,
  List
} from "../../globals/icons";
import DatePicker from "material-ui/DatePicker";
import TimePicker from "material-ui/TimePicker";
import {DatabaseIcon} from "../../icons/icons";
import InfoPopup from "../../components/infoPopup";
import CurrentTime from "../../components/currentTime";
import XHRUploader from "react-xhr-uploader/dist-modules/index";
import CircularProgress from "material-ui/CircularProgress";
import {deepOrange500} from "material-ui/styles/colors";
import {Tabs, Tab} from "material-ui/Tabs";

class ServerSettingsModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);

    this.state.SelectedTab = "settings"

    let uriParams = this.globs.GetUriParams();

    if (uriParams.tab != "") {
      this.state.SelectedTab = uriParams.tab;
    }

    this.handleLoginAttemptsChange = (event, index, value) => {
      this.setComponentState({LockoutSettings: {
        Lockout: {
          Value: value,
          Errors: {Value: ""}
        },
       }});
    };


    this.handleTabChange = (value, cb) => {

      if (!value.hasOwnProperty("target")) {
        this.setComponentState({
          SelectedTab: value,
        });
      }
    };


    this.save = () => {
      //console.log(window.pageState);
      window.api.post({action: "UpdateServerSettings", state: this.state, controller:"serverSettingsModify"});
    };

    this.updateGatewaySettings = () => {
      //console.log(window.pageState);
      window.api.post({action: "UpdateGatewaySettings", state: this.state, controller:"serverSettingsModify"});
    };

    this.updateGatewayTimeSettings = () => {
        window.api.post({action: "UpdateGatewayTimeSettings", state: this.state, controller:"serverSettingsModify"});
    }

    this.enableNTPServer = () => {
        window.api.post({action: "EnableNTPServer", state: this.state, controller:"serverSettingsModify"});
    }

    this.updateLockoutSettings = () => {
      //console.log(window.pageState);
      window.api.post({action: "UpdateLockoutSettings", state: this.state, controller:"serverSettingsModify"});
    };

    this.getTimeZones = () => {
      return this.state.TimeZones.map((tz) => <MenuItem
        key={tz.Location}
        value={tz.Location}
        primaryText={tz.Location + " (" + tz.Country + ")"}
      />)
    };
  }

  render() {
    var lockoutItems = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"];
    var lockoutMap = lockoutItems.map(function(i) {
      if (i == "0") {
        return <MenuItem
            value={lockoutItems[i]}
            primaryText={"Never"}
            key={lockoutItems[i]}
          />;
      } else {
        return <MenuItem
            value={lockoutItems[i]}
            primaryText={lockoutItems[i]}
            key={lockoutItems[i]}
          />;
      }
    });

    var securtiyTab = this.state.SelectedTab != "security" ? null :(
      <div>
      <br/>
      <CenteredPaperGrid>
        <div>{window.pageContent.LockoutSettings}</div>
        <SelectField
          floatingLabelText={window.pageContent.LockoutLoginAttempts}
          hintText={window.pageContent.LockoutLoginAttempts}
          onChange={this.handleLoginAttemptsChange}
          errorText={this.globs.translate(this.state.LockoutSettings.Lockout.Errors.Value)}
          value={this.state.LockoutSettings.Lockout.Value}
        >
          {lockoutMap}
        </SelectField>
        <br />
        <RaisedButton
            label={window.pageContent.SaveSettings}
            onTouchTap={this.updateLockoutSettings}
            secondary={true}
        />
      </CenteredPaperGrid>

      </div>
    );

    var settingsTab = (
      <div>
      <CenteredPaperGrid>
        <div style={{fontWeight:"Bold"}}>{window.pageContent.TimeSettings}</div>
        <br/>
        <CurrentTime/>
        <SelectField
          floatingLabelText={window.pageContent.TimeZone}
          hintText={window.pageContent.TimeZone}
          onChange={(event, index, value) => {this.setComponentState({TimeZone:{Value:value}})}}
          style={{width: 400}}
          value={this.state.TimeZone.hasOwnProperty("Value") ? this.state.TimeZone.Value: "0"}
        >
          {this.getTimeZones()}
        </SelectField>
        <br/>
        <DatePicker
          hintText={window.pageContent.SetDate}
          onChange={(event, d) => {
            this.setComponentState({DateToSet:d})
          }}
          value = {this.state.DateToSet}
        />
        <br/>
        <TimePicker
          hintText={window.pageContent.SetTime}
          onChange={(event, t) => {
            this.setComponentState({TimeToSet:t})
          }}
          value = {this.state.TimeToSet}
        />
        <br/>
        <span className="AlignerRight">
        <RaisedButton
            label={window.pageContent.SaveSettings}
            onTouchTap={this.updateGatewayTimeSettings}
            secondary={true}
        />
        <RaisedButton
            label={window.pageContent.EnableNTP}
            onTouchTap={this.enableNTPServer}
            secondary={true}
        />
        </span>
      </CenteredPaperGrid>
      </div>
    );

    var categoryTabs = (
      <span>
        <input type="hidden" value={this.state.SelectedTab}/>
        <Tabs
          value={this.state.SelectedTab}
          onChange={this.handleTabChange}
          style={{width: "100%"}} >
            <Tab icon={<Settings color={"white"}/>} label={window.pageContent.Settings} value="settings">{settingsTab}</Tab>
            <Tab icon={<Security color={"white"}/>} label={window.pageContent.Security} value="security">{securtiyTab}</Tab>
        </Tabs>
      </span>
    );

    return (
        <div className="page">
          {categoryTabs}
        </div>
    );
  }
}

export default ServerSettingsModify;
