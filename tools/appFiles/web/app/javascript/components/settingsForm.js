/**
 * Created by Dan on 4/14/17.
 */
import {
  React,
  BaseComponent,
  Toggle,
  TextField,
  RaisedButton,
  SelectField,
  MenuItem
} from '../globals/forms';
import PaperExpander from '../components/paperExpander'

class SettingsForm extends BaseComponent {
  constructor(props, context) {
    super(props, context);

    this.state = {
      Settings:this.props.settings
    };

    this.handleSave = (groupName) => {
      if (!this.validate(groupName)) {
        return;
      }

      this.clearGroupErrors(groupName);
      this.setComponentState({Settings:this.state.Settings}, () => {
      if (this.props.onSave != undefined) {
        this.props.onSave(this.state.Settings, groupName);
      }});

    }

    this.validate = (groupName) => {

      var pass = true;

      for (var i =0; i < this.state.Settings.length; i++) {
        var s = this.state.Settings[i];
        if (s.GroupName == groupName && s.Validation != "") {
          if (s.Validation == "OmnistreamMulticast" && s.Value != "") {
            var min = globals.IPToNumber("225.0.0.1")
            var max = globals.IPToNumber("239.255.255.255")
            var val = globals.IPToNumber(s.Value);
            if (isNaN(val) || val < min || val > max) {
              s.Error = "Error:  Not a Valid Multicast Address";
              pass = false;
            }
          }
          if (s.Validation == "OmnistreamPort" && s.Value != "") {
            var val = parseInt(s.Value);
            if (val < 1000 || val > 12000) {
              s.Error = "Error:  Invalid Port Number (1000 to 12000)";
              pass = false;
            }
            if (val % 4 > 0) {
              s.Error = "Error:  Invalid Port Number (Must be multiple of 4)";
              pass = false;
            }
          }
          if (s.Validation.indexOf("Length") != -1 && s.Value != "") {
            var values = s.Validation.split(":");
            if (values[1] == ">") {
              if (s.Value.length < parseInt(values[2])) {
                s.Error = "Error:  Length must be greater than " + values[2];
                pass = false;
              }
            }
          }
          if (s.Validation.indexOf("Depends") != -1 && s.Value != "" && s.Value == "true") {
            var values = s.Validation.split(":");

            var sDepends = this.getSetting({Key:values[1], Index:s.Index});
            if (sDepends.Value == "") {
              s.Error = "Error:  " + s.Name + " depends on a valid " + sDepends.Name;
              pass = false;
            }
          }
        }
      }

      if (pass == false) {
        this.setComponentState({Settings:this.state.Settings});
      }

      return pass;
    }

    this.getSetting = (s) => {
      for(var i = 0; i < this.state.Settings.length; i++) {
        var setting = this.state.Settings[i];
        if (setting.Key == s.Key && setting.Index == s.Index) {
          return setting;
        }
      }
    }

    this.generateSubgroup = (cat) => {

      // return this.generateSetting(s);
      return (
        <div>
        {
          cat.SubGroups.map((sg) => {

            var groupedItems = [];
            for (var i = 0; i < cat.Items.length; i++) {
              var s = cat.Items[i];
              if (s.SubGroup == sg) {
                groupedItems.push(s);
              }
            }

            return (
              <div key={window.globals.guid()}>
                <span className="AlignerLeft" style={{marginLeft:5}}>
                <div style={{fontWeight:"bold", width:150}}>{sg}</div>
                <div>
                {
                  groupedItems.map((s, i) => {
                    return this.generateSetting(s, i);
                  })
                }
                </div>
                </span>
                <br/>
              </div>
            )

          })
        }
        </div>
      );

    }

    this.generateSetting = (s, i) => {
      if(!s.Hide) {
        if (s.Options.length != 0) {
          return (
            <span key={window.globals.guid()} className="AlignerRight" style={{marginTop:5, marginLeft:5, width:620}}>
              <div style={{minWidth: 250}}>{(window.pageContent[s.Key] == undefined) ? s.Name : window.pageContent[s.Key]}</div>
              <SelectField
                id={window.globals.guid()}
                value={s.Value}
                onChange={(event, index, value)=> {
                  s.Value = value;
                  this.setComponentState({Settings: this.state.Settings});
                }}
                key={window.globals.guid()}
                inputStyle={{marginTop:-7}}
                style={{marginLeft:10, width:250, height:35, marginTop:-5}}
              >

                {s.Options.map((opt) => {
                    return(<MenuItem
                            key={opt.Value}
                            value={opt.Value}
                            primaryText={opt.Display}
                            />);
                  })
              }
              </SelectField>
              <div style={{opacity:0.5,fontSize:10, marginLeft:5}}>{s.Hint}</div>
            </span>
          );

        } else {

          switch (s.Type) {
          case "bool":
            return (
                    <Toggle
                      key={window.globals.guid()}
                      style={{marginLeft:5, width:250, marginTop:5}}
                      label={(window.pageContent[s.Key] == undefined) ? s.Name : window.pageContent[s.Key]}
                      labelPosition="left"
                      disabled={(s.ReadOnly === true) ? true : false}
                      defaultToggled={(s.Value == 'true') ? true : false}
                      onToggle={(e,v,t) => {
                        s.Value = v.toString();
                      }}
                    />
                  );
            break;
          case "string":
              return (
                      <span key={window.globals.guid()} className="AlignerRight" style={{marginTop:5, marginLeft:5, width:620}}>
                      <div>{(window.pageContent[s.Key] == undefined) ? s.Name : window.pageContent[s.Key]}</div>
                      <TextField
                        id={window.globals.guid()}
                        floatingLabelStyle={{top:0, color:"black"}}
                        inputStyle={{marginTop:-7}}
                        key={window.globals.guid()}
                        style={{marginLeft:10, width:250, height:30, marginTop:-5}}
                        defaultValue={s.Value}
                        disabled={(s.ReadOnly === true) ? true : false}
                        onChange={(event) =>  {
                          s.Value = event.target.value;
                        }}
                      />
                      <div style={{opacity:0.5, fontSize:10, marginLeft:5}}>{s.Hint}</div>
                      </span>
                    );
              break;
            case "int":
                return (
                        <span key={window.globals.guid()} className="AlignerRight" style={{marginTop:5, marginLeft:5, width:620}}>
                        <div>{(window.pageContent[s.Key] == undefined) ? s.Name : window.pageContent[s.Key]}</div>
                        <TextField
                          id={window.globals.guid()}
                          floatingLabelStyle={{top:0, color:"black"}}
                          inputStyle={{marginTop:-7}}
                          key={window.globals.guid()}
                          style={{marginLeft:10, width:250, height:30, marginTop:-5}}
                          defaultValue={s.Value}
                          disabled={(s.ReadOnly === true) ? true : false}
                          onChange={(event) =>  {
                            s.Value = event.target.value;
                          }}
                        />
                        <div style={{opacity:0.5, fontSize:10, marginLeft:5}}>{s.Hint}</div>
                        </span>
                      );
                break;
            default:
              return null;
          }
        }
      }
    }
  }

  handleDeviceSettingsUpdate(data) {
    if (data.Id == this.props.settingsKey) {
      this.setComponentState({Settings:data.Settings});
    }
  }

  componentDidMount() {
    this.deviceSettingsUpdateCallbackId = window.api.registerSocketCallback((data) => {this.handleDeviceSettingsUpdate(data)}, "DeviceSettingsUpdate");
  }

  componentWillUnmount() {
    window.api.unRegisterSocketCallback(this.deviceSettingsUpdateCallbackId);
  }

  setGroupError(groupName, message) {

    for (var i =0; i < this.state.Settings.length; i++) {
      var s = this.state.Settings[i];
      if (s.GroupName == groupName) {
        s.ErrorMessage = message;
        this.setComponentState({Settings:this.state.Settings});
        return
      }
    }
  }

  getGroupError(groupName) {
    for (var i =0; i < this.state.Settings.length; i++) {
      var s = this.state.Settings[i];
      if (s.GroupName == groupName && s.Error && s.Error != "") {
        return s.Error;
      }
    }
    return "";
  }

  clearGroupErrors(groupName) {
    for (var i =0; i < this.state.Settings.length; i++) {
      var s = this.state.Settings[i];
      if (s.GroupName == groupName && s.Error != "") {
        s.Error = "";
      }
    }
  }

  render() {
    try {
      this.logRender();

      var categories = [];

      if (this.state == undefined || this.state.Settings == undefined || this.state.Settings == null) {
        return null;
      }

      for (var i =0; i < this.state.Settings.length; i++) {
        var s = this.state.Settings[i];
        var foundCat = false;
        for (var j = 0; j < categories.length; j++) {
          if (categories[j].Name == s.GroupName) {
            categories[j].Items.push(s);
            foundCat = true;
            break;
          }
        }
        if (!foundCat) {
          var items = [];
          items.push(s);
          var errorMessage = this.getGroupError(s.GroupName);
          categories.push({Name:s.GroupName, Items:items, ErrorMessage:errorMessage});
        }
      }

      for (var i = 0; i < categories.length; i++) {
        var c = categories[i];
        var subGroups = [];

        for (var j = 0; j < c.Items.length; j++) {
          var setting = c.Items[j];
          var foundSubGroup = false;
          for (var z = 0; z < subGroups.length; z++) {
            if (setting.SubGroup == subGroups[z]) {
              foundSubGroup = true;
            }
          }
          if (foundSubGroup == false) {
            subGroups.push(setting.SubGroup);
          }
        }
        c.SubGroups = subGroups;
      }

      return (
        <div style={{height:500}}>{
            categories.map((cat) => {
              return (
                <div key={cat.Name}>
                <PaperExpander title={cat.Name.toString()} style={{marginLeft:10,width:this.props.width - 40}}>

                    <div>
                      <div>
                      {
                        this.generateSubgroup(cat)
                      }
                      </div>
                      <span className="AlignerRight">
                        <RaisedButton
                          label={window.appContent.RoomModifyDeviceConfigureSave + " " + cat.Name}
                          onTouchTap={(e) => this.handleSave(cat.Name)}
                          secondary={true}
                          style={{marginBottom:20,marginLeft:20, marginTop:20, width:250}}
                        />
                        {
                          (cat.ErrorMessage && cat.ErrorMessage != "") ?
                          <div style={{color:"red", marginRight:20}}>{cat.ErrorMessage}</div>
                          : null
                        }
                      </span>
                    </div>


                </PaperExpander>
                <br/>
                </div>
              )
            })
          }
        </div>
      )
    } catch(e) {
      return this.globs.ComponentError(this.getClassName(), e.message);
    }
  }

}


SettingsForm.propTypes = {
  settings: React.PropTypes.array,
  width:React.PropTypes.number,
  onSave:React.PropTypes.func,
  settingsKey:React.PropTypes.string
};

export default SettingsForm;
