import UserAddModify from "../userModify/userModifyComponents";
import {
  React,
  CenteredPaperGrid,
  BackPage,
  TextField,
  RaisedButton,
  RadioButton,
  RadioButtonGroup,
  SelectField,
  MenuItem,
  List,
  ListItem
} from "../../globals/forms";

class UserAdd extends UserAddModify {
  constructor(props, context) {
    super(props, context);
    this.state.User.TimeZone = moment.tz.guess();
    this.state.User.DateFormat = "mm/dd/yyyy";
    this.state.Password = Math.random().toString(36).slice(-8);
    this.state.User.Language = this.state.UserLocale;


    this.handleAddEmailChange = (event) => {
      this.setComponentState({User: {
        Email: event.target.value,
        Errors: {Email: ""}
      }});
    };

    this.handleAddFirstChange = (event) => {
      this.setComponentState({User: {
        First: event.target.value,
        Errors: {First: ""}
      }});
    };

    this.handleAddLastChange = (event) => {
      this.setComponentState({User: {
        Last: event.target.value,
        Errors: {Last: ""}
      }});
    };

    this.handleAddPasswordChange = (event) => {
      this.setComponentState({
        Password: event.target.value,
        PasswordErrors: ""
      });
    };

    this.handleAddConfirmPasswordChange = (event) => {
      this.setComponentState({
        ConfirmPassword: event.target.value,
        ConfirmPasswordErrors: ""
      });
    };

    this.handleAddHiddenValues = (event) => {
      this.setComponentState({
        User: {
          DefaultAccountId: window.appState.AccountId,
        },
        AccountRole: {
          AccountId: window.appState.AccountId
        }
      });
    };

    this.handleRoleIdChange = (event, index, value) =>  {
      this.setComponentState({
        AccountRole: {
          RoleId: value
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

    this.saveNew = () => {
      this.handleAddHiddenValues();
      //console.log('window.api.post({action: "CreateNewUser", state: JSON.parse(\'' + JSON.stringify(this.state) + '\'), controller:"userModify"});');
      window.api.post({action: "CreateNewUser", state: this.state, controller:"userModify"});
    };
  }

  render() {
    this.logRender();
    var items = this.globs.map(this.state.Roles, (role) =>{
      return (
        <MenuItem value={role.Id} primaryText={role.Name} key={role.Id} />
      );
    });
    return (
        <CenteredPaperGrid>
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditEmail}
            hintText={"* " + window.pageContent.UserAddEditEmail}
            fullWidth={true}
            value={this.state.User.Email}
            onChange={this.handleAddEmailChange}
            errorText={this.globs.translate(this.state.User.Errors.Email)}
          />
          <br />
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditFirstName}
            hintText={"* " + window.pageContent.UserAddEditFirstName}
            fullWidth={false}
            value={this.state.User.First}
            onChange={this.handleAddFirstChange}
            errorText={this.globs.translate(this.state.User.Errors.First)}
          />
          <br />
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditLastName}
            hintText={"* " + window.pageContent.UserAddEditLastName}
            fullWidth={false}
            value={this.state.User.Last}
            onChange={this.handleAddLastChange}
            errorText={this.globs.translate(this.state.User.Errors.Last)}
          />
          <br />
          <TextField
            floatingLabelText={"* " + window.pageContent.UserAddEditPassword}
            hintText={"* " + window.pageContent.UserAddEditPassword}
            fullWidth={false}
            defaultValue={this.state.Password}
            onChange={this.handleAddPasswordChange}
            errorText={this.globs.translate(this.state.PasswordErrors)}
          />
          <br />
          <SelectField 
            value={this.state.AccountRole.RoleId}
            onChange={this.handleRoleIdChange}
            floatingLabelText={"* " + window.appContent.AccountInviteRoleType}
            hintText={"* " + window.appContent.AccountInviteRoleType}
            errorText={this.globs.translate(this.state.AccountRole.Errors.RoleId)}
          >
            {items}
          </SelectField>
          <br />
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
              label={window.appContent.CreateUser}
              onTouchTap={this.saveNew}
              secondary={true}
          />

        </CenteredPaperGrid>
    );
  }
}

export default BackPage(UserAdd);
