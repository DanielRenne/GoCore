import {
  React,
  CenteredPaperGrid,
  BasePageComponent,
  TextField,
  RaisedButton,
  RadioButton,
  RadioButtonGroup,
  Toggle,
  SelectField,
  MenuItem,
  List,
  ListItem,
  Grid,
  Row,
  Col,
  PhoneInput
} from "../../globals/forms";

class UserAddModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.createComponentEvents();
  }

  createComponentEvents() {
    //Edit

    this.handleEmailChange = (event) => {
      this.setComponentState({
        EmailChanged: true,
        User: {
          Email: event.target.value,
          Errors: {Email: ""}
        }
      });
    };

    this.handleFirstNameChange = (event) => {
      this.setComponentState({
        User: {
          First: event.target.value,
          Errors: {First: ""}
        }
      });
    };

    this.handleLastNameChange = (event) => {
      this.setComponentState({
        User: {
          Last: event.target.value,
          Errors: {Last: ""}
        }
      });
    };

    this.handlePhoneChange = (phoneState) => {
      this.setComponentState({User: {
        Phone: phoneState,
        Errors: {Phone: ""}
      }});
    };

    this.handleExtChange = (event) => {
      this.setComponentState({
        User: {
          Ext: event.target.value,
          Errors: {Ext: ""}
        }
      });
    };

    this.handleMobileChange = (phoneState) => {
      this.setComponentState({User: {
        Mobile: phoneState,
        Errors: {Mobile: ""}
      }});
    };

    this.handleJobTitleChange = (event) => {
      this.setComponentState({
        User: {
          JobTitle: event.target.value,
          Errors: {JobTitle: ""}
        }
      });
    };

    this.handleOfficeNameChange = (event) => {
      this.setComponentState({
        User: {
          OfficeName: event.target.value,
          Errors: {OfficeName: ""}
        }
      });
    };

    this.handleDeptChange = (event) => {
      this.setComponentState({
        User: {
          Dept: event.target.value,
          Errors: {Dept: ""}
        }
      });
    };

    this.handleBioChange = (event) => {
      this.setComponentState({
        User: {
          Bio: event.target.value,
          Errors: {Bio: ""}
        }
      });
    };

    this.handleSkypeIdChange = (event) => {
      this.setComponentState({
        User: {
          SkypeId: event.target.value,
          Errors: {SkypeId: ""}
        }
      });
    };

    this.handleProfileIconChange = (event) => {
      this.setComponentState({
        User: {
          PhotoIcon: event.target.value,
          Errors: {PhotoIcon: ""}
        }
      });
    };

    this.handleDefaultAccountChange = (event, index, value) => {
      this.setComponentState({
        User: {
          DefaultAccountId: value,
          Errors: {DefaultAccountId: ""}
        }
      });
    };

    this.handleLockedChange = (event, index, value) => {
      this.setComponentState({
        User: {
          Locked: value,
        }
      }, () => {
        this.setComponentState({
          User: {
            LoginAttempts: (value) ? 999: 0
          }
        });
      });
    };

    this.handleUserEnforcePasswordChangeChange = (e, value) => {
      this.setComponentState({
        User: {
          EnforcePasswordChange: value,
        }
      });
    };

    this.handleLanguageChange = (event, index, value) => {
      this.setComponentState({
        User: {
          Language: value,
          Errors: {Language: ""}
        }
      });
    };

    this.handleTimeZoneChange = (event, index, value) => {
      this.setComponentState({
        User: {
          TimeZone: value,
          Errors: {TimeZone: ""}
        }
      });
    };

    this.handleDateFormatChange = (event) => {
      this.setComponentState({
        User: {
          DateFormat: event.target.value,
          Errors: {DateFormat: ""}
        }
      });
    };

    this.handleChangeAcctRole = (event, index, value) => {
      this.setComponentState({
        AccountRole: {
          RoleId: value,
          Errors: {RoleId: ""}
        }
      });
    };


    this.save = () => {
        if (window.location.href.indexOf("userProfile") != - 1) {
          this.state.CurrentPage = "userProfile";
        } else {
          this.state.CurrentPage = "userModify";
        }
        window.api.post({action: "UpdateUserDetails", state: this.state, controller:"userModify"});
        
     };

    this.saveAccountRole = () => {
        this.state.CurrentPage = "userModify";
        window.api.post({action: "UpdateAccountRole", state: this.state, controller:"userModify"});
     };

    this.handleNewPasswordChange = (event) => {
      this.setComponentState({
        Password: event.target.value,
        PasswordErrors: ""
      });
    };

    this.handleNewConfirmPasswordChange = (event) => {
      this.setComponentState({
        ConfirmPassword: event.target.value,
        ConfirmPasswordErrors: ""
      });
    };

    this.handleSaveNewPassword = () => {
      window.api.post({action: "ChangeUserPassword", state: this.state, controller:"userModify"})
    };

    this.getTimeZones = () => {
      return this.globs.map(this.state.TimeZones, (tz) => <MenuItem
        key={tz.Location}
        value={tz.Location}
        primaryText={tz.Location + " (" + tz.Country + ")"}
      />)
    };

    this.getLocales = () => {
      return this.globs.map(this.state.Locales, (locale) => <MenuItem
        key={locale.Value}
        value={locale.Value}
        primaryText={locale.Language + " (" + locale.Value + ")"}
      />)
    };

    this.timeZoneAndLanguage = () => {
      return <span>
          <SelectField
            floatingLabelText={window.pageContent.UserAddEditLanguage}
            hintText={window.pageContent.UserAddEditLanguage}
            onChange={this.handleLanguageChange}
            style={{width: 400}}
            value={this.state.User.Language}
          >
          {this.getLocales()}
          </SelectField>
          <br />
          <SelectField
            floatingLabelText={window.pageContent.UserAddEditSetTimeZone}
            hintText={window.pageContent.UserAddEditSetTimeZone}
            onChange={this.handleTimeZoneChange}
            style={{width: 400}}
            value={this.state.User.TimeZone}
          >
            {this.getTimeZones()}
          </SelectField>
          <br />
      </span>
    }
  }

  componentWillReceiveProps(nextProps) {
    if (window.appState.DeveloperMode) {
      this.createComponentEvents();
    }
    return true;
  }

  render() {
    this.logRender();
    var lockedItemsObj = [{display:"Locked", value:true},{display:"Unlocked", value: false}];

    var lockedItems = this.globs.map(lockedItemsObj, (obj) => {
      return <MenuItem
        key={obj.value}
        value={obj.value}
        primaryText={obj.display}
        />;
    });

    var items = this.globs.map(this.state.Accounts, (ac) => {
      return <MenuItem
        key={ac.Id}
        value={ac.Id}
        primaryText={ac.AccountName + " - " + ac.AccountName + " - " + ac.City + ", " + ac.Region}
        />;
    });
    var roles = this.globs.map(this.state.Roles, (r) => {
      return <MenuItem
        key={r.Id}
        value={r.Id}
        primaryText={r.Name}
        />;
    });

    var acctRole = (this.globs.HasRole("USER_CHANGE_ROLE") || this.globs.IsMasterAccount()) ? (
      <div>
        <CenteredPaperGrid>
          <div>{window.pageContent.ChangeRole}</div>
          <SelectField
            floatingLabelText={"* " + window.pageContent.UserAddEditAccountRole}
            hintText={"* " + window.pageContent.UserAddEditAccountRole}
            onChange={this.handleChangeAcctRole}
            style={{width: 400}}
            value={this.state.AccountRole.RoleId}
          >
          {roles}
          </SelectField>
          <br />
          <br />
          <RaisedButton
              label={window.appContent.SaveChanges}
              onTouchTap={this.saveAccountRole}
              secondary={true}
          />
        </CenteredPaperGrid>
        <br />
      </div>
    ) : (<div></div>);

    if (this.state.User.TimeZone == "") {
      this.state.User.TimeZone = moment.tz.guess();
    }

    if (this.state.User.DateFormat == "") {
      this.state.User.DateFormat = "mm/dd/yyyy";
    }

    if (this.state.User.Language == "") {
      this.state.User.Language = this.state.UserLocale;
    }


    return (
      <div>
        <CenteredPaperGrid>
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditEmail}
            hintText={"* " + window.pageContent.UserAddEditEmail}
            defaultValue={this.state.User.Email}
            fullWidth={true}
            onChange={this.handleEmailChange}
            errorText={this.globs.translate(this.state.User.Errors.Email)}
          />
          <br />
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditFirstName}
            hintText={"* " + window.pageContent.UserAddEditFirstName}
            defaultValue={this.state.User.First}
            fullWidth={false}
            onChange={this.handleFirstNameChange}
            errorText={this.globs.translate(this.state.User.Errors.First)}
          />
          <br />
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditLastName}
            hintText={"* " + window.pageContent.UserAddEditLastName}
            defaultValue={this.state.User.Last}
            fullWidth={false}
            onChange={this.handleLastNameChange}
            errorText={this.globs.translate(this.state.User.Errors.Last)}
          />
          <br />
          <Grid style={{paddingLeft:0}}>
            <Row>
              <Col md={3}>
                <PhoneInput
                    InitialValue={this.state.User.Phone.Value}
                    ErrorText={this.globs.translate(this.state.User.Errors.Phone.Value)}
                    OnChange={this.handlePhoneChange}
                    Label={window.pageContent.UserAddEditPhone}
                />
              </Col>
              <Col md={3}>
                <TextField
                    floatingLabelText={window.pageContent.UserAddEditExtension}
                    hintText={window.pageContent.UserAddEditExtension}
                    defaultValue={this.state.User.Ext}
                    fullWidth={false}
                    onChange={this.handleExtChange}
                    style={{width: 100, 'marginLeft': 20}}
                  />
              </Col>
            </Row>
          </Grid>
          <br />
          <PhoneInput
              InitialValue={this.state.User.Mobile.Value}
              ErrorText={this.globs.translate(this.state.User.Errors.Mobile.Value)}
              OnChange={this.handleMobileChange}
              Label={window.pageContent.UserAddEditMobile}
          />
          <br />
          <TextField
            floatingLabelText={window.pageContent.UserAddEditJobTitle}
            hintText={window.pageContent.UserAddEditJobTitle}
            defaultValue={this.state.User.JobTitle}
            fullWidth={false}
            onChange={this.handleJobTitleChange}
          />
          <br />
          <TextField
            floatingLabelText={window.pageContent.UserAddEditOfficeName}
            hintText={window.pageContent.UserAddEditOfficeName}
            defaultValue={this.state.User.OfficeName}
            fullWidth={false}
            onChange={this.handleOfficeNameChange}
          />
          <br />
          <TextField
            floatingLabelText={window.pageContent.UserAddEditDepartment}
            hintText={window.pageContent.UserAddEditDepartment}
            defaultValue={this.state.User.Dept}
            fullWidth={false}
            onChange={this.handleDeptChange}
          />
          <br />
          <TextField
            floatingLabelText={window.pageContent.UserAddEditUserBio}
            hintText={window.pageContent.UserAddEditUserBio}
            defaultValue={this.state.User.Bio}
            fullWidth={true}
            onChange={this.handleBioChange}
            multiLine={true}
            rows={1}
          />
          <br />
          <TextField
            floatingLabelText={window.pageContent.UserAddEditSkypeId}
            hintText={window.pageContent.UserAddEditSkypeId}
            defaultValue={this.state.User.SkypeId}
            fullWidth={false}
            onChange={this.handleSkypeIdChange}
          />
          <br />
          <TextField
            floatingLabelText={window.pageContent.UserAddEditUserProfileIcon}
            hintText={window.pageContent.UserAddEditUserProfileIcon}
            defaultValue={this.state.User.PhotoIcon}
            fullWidth={false}
            onChange={this.handleProfileIconChange}
            style={{display: "none"}}
          />
          <br />
          {/*
            <SelectField
                        floatingLabelText={window.pageContent.UserAddEditSetDefaultAccount}
                        hintText={window.pageContent.UserAddEditSetDefaultAccount}
                        onChange={this.handleDefaultAccountChange}
                        style={{width: 400}}
                        value={this.state.User.DefaultAccountId}
                        hidden={true}
                      >
                      {items}
                      </SelectField>
            <br />
          */}
          {document.location.href.indexOf("userProfile") == -1 ?
            <SelectField
              floatingLabelText={window.pageContent.UserLockedStatus}
              hintText={window.pageContent.UserLockedStatus}
              onChange={this.handleLockedChange}
              style={{width: 400}}
              value={this.state.User.Locked}
            >
            {lockedItems}
            </SelectField> : null
          }
          {document.location.href.indexOf("userProfile") == -1 ?
              <span>
                <br />
                <br />
                <Toggle
                  label={window.pageContent.UpdatePasswordNextLogin}
                  defaultToggled={this.state.User.EnforcePasswordChange}
                  labelPosition="right"
                  onToggle={this.handleUserEnforcePasswordChangeChange}
                />
              </span>
              : null}
          <br />
          <br />
          <RaisedButton
              label={window.appContent.SaveChanges}
              onTouchTap={this.save}
              secondary={true}
          />
        </CenteredPaperGrid>

        <br />
        <CenteredPaperGrid>
          <div>{window.pageContent.LanguageFormatSettings}</div>
          {this.timeZoneAndLanguage()}

          <TextField
            floatingLabelText={window.pageContent.UserAddEditSetDateField}
            hintText={window.pageContent.UserAddEditSetDateField}
            defaultValue={this.state.User.DateFormat}
            fullWidth={false}
            onChange={this.handleDateFormatChange}
            errorText={this.globs.translate(this.state.DateFormatError)}
          />
          <br />
          <br />
          <RaisedButton
              label={window.appContent.SaveChanges}
              onTouchTap={this.save}
              secondary={true}
          />
        </CenteredPaperGrid>
        <br />

        {acctRole}

        {document.location.href.indexOf("userProfile") != -1 ?
          <CenteredPaperGrid>
            <div>{window.pageContent.ChangePassword}</div>
            <TextField
              floatingLabelText={window.pageContent.UserAddEditPassword}
              hintText={window.pageContent.UserAddEditPassword}
              fullWidth={false}
              type="password"
              onChange={this.handleNewPasswordChange}
              errorText={this.globs.translate(this.state.PasswordErrors)}
            />
            <br />
            <TextField
              floatingLabelText={window.pageContent.UserAddEditConfPassword}
              hintText={window.pageContent.UserAddEditConfPassword}
              fullWidth={false}
              type="password"
              onChange={this.handleNewConfirmPasswordChange}
              errorText={this.globs.translate(this.state.ConfirmPasswordErrors)}
            />
            <br />
            <br />
            <RaisedButton
                label={window.pageContent.SaveNewPassword}
                onTouchTap={this.handleSaveNewPassword}
                secondary={true}
            />
          </CenteredPaperGrid> : null}

      </div>
    );
  }
}

export default UserAddModify;
