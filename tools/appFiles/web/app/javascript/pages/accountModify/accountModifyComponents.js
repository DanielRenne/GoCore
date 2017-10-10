import {
  React,
  CenteredPaperGrid,
  BasePageComponent,
  BackPage,
  ConfirmDelete,
  TextField,
  RaisedButton,
  IconButton,
  RadioButton,
  RadioButtonGroup,
  Toggle,
  SelectField,
  MenuItem,
  List,
  ListItem,
  Divider,
  AutoComplete,
  PhoneInput
} from "../../globals/forms";
import Paper from "material-ui/Paper";
import Subheader from "material-ui/Subheader";
import {AddIcon, DeleteIcon} from "../../globals/icons";

class AccountAddModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.confirmDeleteRef;
    this.phoneRef;
    this.secondaryPhoneRef;
    this.state.Account.States = this.state.States[this.state.Account.CountryId];
    this.state.Account.StateExist = this.state.Account.States === undefined ? false : true;

    this.handleAccountNameChange = (event) => {
      this.setComponentState({
        Account: {
          AccountName: event.target.value,
          Errors: {AccountName: ""}
        }
      });
    };

    this.handleAccountCompanyNameChange = (value) => {
      if (value == "" || (value.hasOwnProperty("target") && value.target.value == "")) {
        this.setComponentState({
          Account: {Errors: {IsEnterprise: "AccountAddEditIsEnterpriseErrorBlank"}}
        });
      } else if (!value.hasOwnProperty("target")) {
        this.setComponentState({
          AccountCompanyMapping: value,
          Account: {Errors: {IsEnterprise: ""}}
        });
      }
    };

    this.handleAddress1Change = (event) => {
      this.setComponentState({
        Account: {
          Address1: event.target.value,
          Errors: {Address1: ""}
        }
      });
    };

    this.handleAddress2Change = (event) => {
      this.setComponentState({
        Account: {
          Address2: event.target.value,
          Errors: {Address2: ""}
        }
      });
    };

    this.handleRegionChange = (event) => {
      this.setComponentState({
        Account: {
          Region: event.target.value,
          Errors: {Region: ""}
        }
      });
    };

    this.handleStateChange = (event, index, value) => {
      var stateId = "", state = "";
      if (this.state.Account.StateExist) {
        stateId = value
        var x = this.globs.map(this.state.Account.States, (v) =>
            v.Id === stateId ? state = v.Name : ""
        )
      } else {
        state = event.target.value
      }

      this.setComponentState({
        Account: {
          StateId: stateId,
          StateName: state,
          Errors: {StateId: "", StateName: ""}
        }
      });
    };

    this.handleCityChange = (event) => {
      this.setComponentState({
        Account: {
          City: event.target.value,
          Errors: {City: ""}
        }
      });
    };

    this.handlePostCodeChange = (event) => {
      this.setComponentState({
        Account: {
          PostCode: event.target.value,
          Errors: {PostCode: ""}
        }
      });
    };

    this.handleCountryIsoChange = (event, index, value) => {
      var states = this.state.States[value];
      this.setComponentState({
        Account: {
          CountryId: value,
          Errors: {CountryId: ""},
          States: states,
          StateExist: states === undefined ? false : true,
          StateName: "",
          StateId: ""
        }
      }, () => {
        this.updateSelectedISO(value, () => {
          this.phoneRef.getRef().setFlag(this.state.SelectedCountryIso, false);
          this.phoneRef.changeCountry(this.state.SelectedCountryIso);
          this.secondaryPhoneRef.getRef().setFlag(this.state.SelectedCountryIso, false);
          this.secondaryPhoneRef.changeCountry(this.state.SelectedCountryIso);
          this.setComponentState({
            Account: {
              PrimaryPhone: {},
              SecondaryPhone: {},
              Errors: {PrimaryPhone: "", SecondaryPhone: ""},
            }
          });
        });
      });
    };

    this.handlePrimaryPhoneChange = (phoneState) => {
      this.setComponentState({
        Account: {
          PrimaryPhone: phoneState,
          Errors: {PrimaryPhone: ""}
        }
      });
    };

    this.handleSecondaryPhoneChange = (phoneState) => {
      this.setComponentState({
        Account: {
          SecondaryPhone: phoneState,
          Errors: {SecondaryPhone: ""}
        }
      });
    };

    this.handleAccountEmailChange = (event) => {
      this.setComponentState({
        Account: {
          Email: event.target.value,
          Errors: {Email: ""}
        }
      });
    };

    this.handleEnterpriseChange = (event, value) => {
      this.setComponentState({
        Account: {
          IsEnterprise: value,
        }
      });
    };

    this.handlePrimaryAccountChange = (event, value) => {
      this.setComponentState({
        Account: {
          IsPrimaryLinkedAccount: value,
        }
      });
    };

    this.save = () => {
      window.api.post({action: "UpdateAccountDetails", state: this.state, controller: "accounts"});
    };

    this.transferAccount = () => {
      window.api.post({action: "TransferAccount", state: this.state, controller: "accounts"});
    };

    this.handleDeleteAccount = () => {
      this.finalizeState();
      window.api.post({
        action: "DeleteAccount", state: this.state, controller: "accounts", callback: (vm) => {
          if (window.appState.DialogOpen || window.appState.DialogGenericOpen) {
            this.confirmDeleteRef.handleClose();
          }
        }
      });
    };
  }


  updateSelectedISO(countryId, callback) {
    if (countryId != "") {
      var filtered = this.state.Countries.filter((c) => c.Id == countryId);
      if (filtered) {
        this.setComponentState({SelectedCountryIso: filtered[0].Iso}, callback);
      } else {
        this.setComponentState({SelectedCountryIso: "us"}, callback);
      }
    } else {
      this.setComponentState({SelectedCountryIso: "us"}, callback);
    }
  }

  finalizeState() {
    if (!this.state.Account.IsEnterprise) {
       this.setComponentState({AccountCompanyMapping: ""});
    }
  }

  componentDidMount() {
    this.updateSelectedISO(this.state.Account.CountryId);

    var updates = {};

    if (this.state.Companies && this.state.Companies.length > 0) {
      updates.Companies = this.state.Companies.map((o)=> o.CompanyName);
    } else {
      updates.Companies = [];
    }

    this.setComponentState(updates);
  }

  sharedCompanyRender() {
    return <div>
      <Toggle
        label={window.pageContent.AccountAddEditIsEnterpriseToggleLabel}
        defaultToggled={this.state.Account.IsEnterprise}
        labelPosition="right"
        onToggle={this.handleEnterpriseChange}
      />
      {(this.state.Account.IsEnterprise) ?     <AutoComplete
            floatingLabelText={"* " + window.pageContent.AccountAddEditIsEnterprise}
            hintText={"* " + window.pageContent.AccountAddEditIsEnterpriseHint}
            filter={AutoComplete.caseInsensitiveFilter}
            fullWidth={true}
            errorText={this.globs.translate(this.state.Account.Errors.IsEnterprise)}
            dataSource={this.state.Companies}
            searchText={this.state.AccountCompanyMapping}
            onBlur={this.handleAccountCompanyNameChange}
            onUpdateInput={this.handleAccountCompanyNameChange}
            onNewRequest={this.handleAccountCompanyNameChange}
          />
      : null}
    </div>
  }

  sharedCityAndPhoneRender() {
    return <div>
        <TextField
            floatingLabelText={"* " + window.pageContent.AccountAddEditCity}
            hintText={"* " + window.pageContent.AccountAddEditCity}
            defaultValue={this.state.Account.City}
            fullWidth={true}
            onChange={this.handleCityChange}
            errorText={this.globs.translate(this.state.Account.Errors.City)}
        />
        <br />
        <PhoneInput
            InitialValue={this.state.Account.PrimaryPhone.Value}
            ErrorText={this.globs.translate(this.state.Account.Errors.PrimaryPhone.Value)}
            OnChange={this.handlePrimaryPhoneChange}
            Label={"* " + window.pageContent.AccountAddEditPrimaryPhone}
            ref={(c) => this.phoneRef = c}
        />
        <PhoneInput
            InitialValue={this.state.Account.SecondaryPhone.Value}
            ErrorText={this.globs.translate(this.state.Account.Errors.SecondaryPhone.Value)}
            OnChange={this.handleSecondaryPhoneChange}
            Label={window.pageContent.AccountAddEditSecondaryPhone}
            ref={(c) => this.secondaryPhoneRef = c}
        />
    </div>
  }


  render() {
    this.logRender();

    var transferAccount = (

      <div>
      <span style={{fontWeight: "bold"}}>Change Dealer</span><br />
      <span>Your dealer stinks. Change dealer.</span><br /><br />
      <RaisedButton
          label={window.pageContent.TransferAccount}
          onTouchTap={this.transferAccount}
          secondary={true}
      />
      <br />
      <br />
      </div>
    );

    transferAccount = "";

    var primaryAccountName = window.pageContent.NoneOnAccounts;
    var primaryCityState = "";

    return (
        <div>
          <div onClick={this.closeDrawers}>
            <CenteredPaperGrid>

                <TextField
                    floatingLabelText={"* " + ((window.appState.AccountTypeShort == "deal") ? window.pageContent.AccountAddEditCustomerName : window.pageContent.AccountAddEditAccountName)}
                    hintText={"* " + window.pageContent.AccountAddEditAccountNameHint}
                    defaultValue={this.state.Account.AccountName}
                    fullWidth={true}
                    onChange={this.handleAccountNameChange}
                    errorText={this.globs.translate(this.state.Account.Errors.AccountName)}
                />
                <TextField
                    floatingLabelText={"* " + window.pageContent.AccountAddEditAddressLine1}
                    hintText={"* " + window.pageContent.AccountAddEditAddressLine1}
                    defaultValue={this.state.Account.Address1}
                    fullWidth={true}
                    onChange={this.handleAddress1Change}
                    errorText={this.globs.translate(this.state.Account.Errors.Address1)}
                />
                <br />
                <TextField
                    floatingLabelText={window.pageContent.AccountAddEditAddressLine2}
                    hintText={window.pageContent.AccountAddEditAddressLine2}
                    defaultValue={this.state.Account.Address2}
                    fullWidth={true}
                    onChange={this.handleAddress2Change}
                    // errorText={window.pageContent[this.state.Account.Errors.Address2]}
                />
                <br />
                <TextField
                    floatingLabelText={"* " + window.pageContent.AccountAddEditPostCode}
                    hintText={"* " + window.pageContent.AccountAddEditPostCode}
                    defaultValue={this.state.Account.PostCode}
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


                {/* <TextField
                    floatingLabelText={"* " + window.pageContent.AccountAddEditRegion}
                    hintText={"* " + window.pageContent.AccountAddEditRegion}
                    defaultValue={this.state.Account.Region}
                    onChange={this.handleRegionChange}
                    errorText={this.globs.translate(this.state.Account.Errors.Region)}
                /> */}
                <br />
                
                {this.sharedCityAndPhoneRender()}
                <TextField
                    floatingLabelText={"* " + window.pageContent.AccountAddEditEmail}
                    hintText={"* " + window.pageContent.AccountAddEditEmail}
                    defaultValue={this.state.Account.Email}
                    fullWidth={true}
                    onChange={this.handleAccountEmailChange}
                    errorText={this.globs.translate(this.state.Account.Errors.Email)}
                />
                <br />
                <TextField
                    floatingLabelText={"* " + window.pageContent.AccountTypeLong}
                    defaultValue={this.globs.translate(this.state.Account.AccountTypeLong)}
                    disabled={true}
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
            {(this.globs.HasRole("ACCOUNT_DELETE")) ? <div>
                <br />
                <CenteredPaperGrid>
                    <span style={{color: "red", fontWeight: "bold"}}>{window.pageContent.Danger}:</span><br /><br />
                    {transferAccount}
                    <Divider />
                    <br />
                    <span>{window.pageContent.DeleteAccountCertain}</span><br /><br />

                    <ConfirmDelete
                        deleteFunction={this.handleDeleteAccount}
                        noCancel={window.appContent.ConfirmDeleteNoCancel}
                        yesDelete={window.appContent.ConfirmDeleteYesDelete}
                        buttonTriggerLabel={window.pageContent.DeleteAccount}
                        dialogTitle={window.pageContent.DeleteAccount}
                        dialogMessage={window.pageContent.ConfirmDeleteMessage}
                        ref={(component) => this.confirmDeleteRef = component}
                        />
                </CenteredPaperGrid>
              </div>: null}
            </div>
        </div>
    );
  }
}

export default BackPage(AccountAddModify);
