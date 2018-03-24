import AccountAddModify from '../accountModify/accountModifyComponents';
import {
  React,
  CenteredPaperGrid,
  BasePageComponent,
  BaseComponent,
  WidgetList,
  AddRecordPage,
  BackPage,
  ConfirmDelete,
  AccountInvite,
  TextField,
  RaisedButton,
  FlatButton,
  IconButton,
  RadioButton,
  RadioButtonGroup,
  Toggle,
  FloatingActionButton,
  ContentAdd,
  Avatar,
  SelectField,
  MenuItem,
  List,
  ListItem,
  Checkbox,
  Divider,
  AutoComplete} from '../../globals/forms';

class AccountAdd extends AccountAddModify {
  constructor(props, context) {
    super(props, context);

    this.createAccount = () => {
      this.finalizeState();
      window.api.post({action: "CreateAccount", state: this.state, controller:"accounts"});
    };

    this.handleSendInvitationChange = (event, value) => {
      this.setComponentState({SendInvite:value});
    };
  }

  componentDidMount() {
    this.setComponentState({Account: {Email: window.appState.UserEmail}})
  }

  render() {
    this.logRender();

    return (
      <CenteredPaperGrid>

        <TextField
            floatingLabelText={"* " + ((window.appState.AccountTypeShort == "deal") ? window.pageContent.AccountAddEditCustomerName : window.pageContent.AccountAddEditAccountName)}
            hintText={"* " + window.pageContent.AccountAddEditAccountNameHint}
            fullWidth={true}
            onChange={this.handleAccountNameChange}
            errorText={this.globs.translate(this.state.Account.Errors.AccountName)}
        />
        <br />
        <TextField
            floatingLabelText={"* " + window.pageContent.AccountAddEditAddressLine1}
            hintText={"* " + window.pageContent.AccountAddEditAddressLine1}
            fullWidth={true}
            onChange={this.handleAddress1Change}
            errorText={this.globs.translate(this.state.Account.Errors.Address1)}
        />
        <br />
        <TextField
            floatingLabelText={window.pageContent.AccountAddEditAddressLine2}
            hintText={window.pageContent.AccountAddEditAddressLine2}
            fullWidth={true}
            onChange={this.handleAddress2Change}
            // errorText={window.pageContent[this.state.Account.Errors.Address2]}
        />
        <br />
        <TextField
            floatingLabelText={"* " + window.pageContent.AccountAddEditPostCode}
            hintText={"* " + window.pageContent.AccountAddEditPostCode}
            onChange={this.handlePostCodeChange}
            errorText={this.globs.translate(this.state.Account.Errors.PostCode)}
        />
        <br />
        <SelectField floatingLabelText={"* " + window.pageContent.AccountAddEditCountryIso}
            value={this.state.Account.CountryId}
            hintText={"* " + window.pageContent.AccountAddEditCountryIso}
            onChange={this.handleCountryIsoChange}
            style={{width: 400}}
            errorText={this.globs.translate(this.state.Account.Errors.CountryId)}>
          {this.globs.map(this.state.Countries, (v) => <MenuItem key={v.Id} value={v.Id} primaryText={v.Name}/>)}
        </SelectField>

        <SelectField floatingLabelText={"* " + window.pageContent.AccountAddEditState}
            value={this.state.Account.StateId}
            hintText={"* " + window.pageContent.AccountAddEditState}
            onChange={this.handleStateChange}
            style={this.state.Account.StateExist ? {width: 400}: {width: 400, display: 'none'}}
            errorText={this.globs.translate(this.state.Account.Errors.StateId)}>
          {this.globs.map(this.state.Account.States, (v) => <MenuItem key={v.Id} value={v.Id} primaryText={v.Name}/>)}
        </SelectField>
        <TextField
            floatingLabelText={"* " + window.pageContent.AccountAddEditState}
            hintText={"* " + window.pageContent.AccountAddEditState}
            fullWidth={true}
            onChange={this.handleStateChange}
            errorText={this.globs.translate(this.state.Account.Errors.StateName)}
            defaultValue={this.state.Account.StateName}
            style={!this.state.Account.StateExist ? {} : {display: 'none'}}
        />

        {/* <br />
        <TextField
            floatingLabelText={"* " + window.pageContent.AccountAddEditRegion}
            hintText={"* " + window.pageContent.AccountAddEditRegion}
            onChange={this.handleRegionChange}
            errorText={this.globs.translate(this.state.Account.Errors.Region)}
        /> */}
        {this.sharedCityAndPhoneRender()}
        <br />
        <span className="Aligner">
          <TextField
              floatingLabelText={"* " + window.pageContent.AccountAddEditEmail}
              hintText={"* " + window.pageContent.AccountAddEditEmail}
              fullWidth={true}
              onChange={this.handleAccountEmailChange}
              value={this.state.Account.Email}
              errorText={this.globs.translate(this.state.Account.Errors.Email)}
          />
          <Toggle
            style={{paddingTop:30}}
            label={window.pageContent.SendInvite}
            defaultToggled={false}
            labelPosition="right"
            onToggle={this.handleSendInvitationChange}/>

        </span>
        <br />
        <RaisedButton
            label={window.pageContent.CreateAccount}
            onTouchTap={this.createAccount}
            secondary={true}
            id="save"
        />
        </CenteredPaperGrid>
    );
  }
}

export default BackPage(AccountAdd);
