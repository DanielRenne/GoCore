import React, {Component} from 'react';
import BasePageComponent from '../../components/basePageComponent';
import WidgetList from '../../components/widgetList';
import InviteIcon from 'material-ui/svg-icons/communication/mail-outline';
import RaisedButton from 'material-ui/RaisedButton';
import {DeleteIcon, EditIcon, ExportIcon} from '../../globals/icons'
import {red500, blueGrey400, blueGrey900, deepOrange500} from 'material-ui/styles/colors';
import IconButton from 'material-ui/IconButton';
import {RevokeUserIcon} from '../../icons/icons';
import {
    ConfirmPopup
} from "../../globals/forms";
import AddOrImportPage from '../../components/addOrImportComponent';

class UserList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);

    this.confirmDeleteRowRef;

    this.handleEditUser = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Joins.User.Id}, controller: "userModify"});
    };

    this.handleDeleteUserInline = (row) => {
      this.confirmDeleteRowRef.handleOpen(row);
    };

    this.handleDeleteUser = (row) => {
      window.api.post({action:"Revoke", state:row, controller:"userList"});
    }
  }

  // Button ideas
  // ----------------------------------------------------------
  // DELETE, EXPORT, COPY, UPDATE (MASS) always
  // CHANGE ROLES (only if all same roles)
  // BOTTOM ACTION BUTTON HOVER.... Import (show large text box to paste emails)
  // Individual row for each row, EDIT, DELETE, RESEND CONFIRMATION


  // Add spinner logic for long running queries to overlay data
  // resend confirmation button?  based on state of the record?  I thought dan was using another table for signups.
  // ensure created updated date automatically adding.  Add in created by and added by columns but ensure that buttons all render on all browser sizes.
  // Revoke access instead of delete if you are a client admin instead of deleting a parent company user
  // Fix checkbox bugs on widget list
  // --  when you check all funny things happen with whats checked
  // -- on paginate forward deselect
  // Add company name to user list inside of client account view?  This way I know what people within my company or the dealer and i can individually revoke certain people out of the dealers access
  // Date formatting on user list - we just need the date, not date time
  //  -- we need to ensure how we are handling dates with timezone offsets
  // Add search functionality (need tons of data for this on testing)


  /*
  * FUTURE AND OTHER WIDGET LISTS
  *
                {
                  func: (rows) => {
                    return true
                  },
                    button: <RaisedButton
                    label="Export"
                    labelColor={blueGrey900}
                    onTouchTap={this.handleInvite2}
                    icon={<ExportIcon/>}
                  />
                },
                {
                  func: (rows) => {
                    return true
                  },
                    button: <RaisedButton
                    label="Copy"
                    labelColor={blueGrey900}
                    onTouchTap={this.handleInvite2}
                    icon={<DeleteIcon/>}
                  />
                },
  * */


  render() {
    this.logRender();
    return (
        <div>
          <WidgetList
              {...this.globs.widgetListDefaults()}
              {...this.globs.widgetListButtonBarOffset()}
              name="userList"
              listViewModel={this.state.WidgetList}
              controller="userList"
              listTitle={window.pageContent.UserListShowingAllUsers}
              showCheckboxes={false}
              fields={[
                {
                  func: (row, currentValue) => {
                    if (currentValue != "") {
                      return <a title={window.appContent.ListEmailThisPerson} href={"mailto:" + row.Joins.User.Email}><InviteIcon style={{width: 16, height: 16, marginRight: 5}} />{currentValue}</a>
                    }
                    return "";
                  },
                  tooltip: window.pageContent.UserListLastNameFirstName,
                  headerDisplay: window.appContent.GlobalListsName,
                  stateKey: "Joins.User.Views.FullName",
                  sortOn: "Last",
                  sortable: false,
                  responsiveKeep: true
                },
                {
                  tooltip: window.pageContent.UserListRole,
                  headerDisplay: window.pageContent.UserListRole,
                  stateKey: "Joins.Role.Name",
                  sortable: false,
                  responsiveKeep: true
                },
                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdated,
                  sortable: false,
                  stateKey: "Joins.User.Views.UpdateFromNow",
                  sortOn: "UpdateDate",
                  tooltipKey: "Joins.User.Views.UpdateDate"
                },
                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdatedBy,
                  sortable: false,
                  stateKey: "Joins.User.Joins.LastUpdateUser.Views.FullName"
                },
              ]}
              rowButtons={[
                {
                  func: (row) => {
                    if (!this.globs.HasRole("USER_MODIFY") || row.Errors.Id == "true") {
                      return false
                    }

                    return true;
                  },
                  button:
                    <RaisedButton
                      title={window.appContent.GlobalButtonsEditThisRecord}
                      onTouchTap={this.handleEditUser}
                      icon={<EditIcon color={blueGrey400}/>}
                    />
                },

                {
                  func: (row) => {
                    if (!this.globs.HasRole("USER_REVOKE_FROM_ACCOUNT")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    style={{height: 36, minHeight: 36, marginTop: 6}}
                    title={window.pageContent.UserListRevokeUser}
                    onTouchTap={this.handleDeleteUserInline}
                    icon={<RevokeUserIcon color={deepOrange500} width={20} height={20} position="absolute" top={10} left={10}/>}
                  />
                }
             ]}
             data={this.state[this.state.WidgetList.DataKey]}
             dataKey={this.state.WidgetList.DataKey}
          />
          <ConfirmPopup
              onSubmit={this.handleDeleteUser}
              areYouSureMsg={window.pageContent.AreYouSure}
              ref={(component) => this.confirmDeleteRowRef = component}
              />
        </div>
    );
  }
}


export default AddOrImportPage(UserList, "CONTROLLER_USERADD", "AddUserTooltip", "user_import.csv", "userList", "ImportCSV", (callback) => {
  return (
          window.global.functions.HasRole("ACCOUNT_INVITE") ?
          <div className="col-sm-4" style={{cursor: 'pointer', marginRight: 26, marginLeft: 52}} onClick={() => callback()}>
            <div className="counter counter-lg counter-inverse bg-red-600 vertical-align height-150">
              <div className="vertical-align-middle">
                <div className="counter-icon margin-bottom-5" style={{marginTop: 18
}}><i className="icon md-account" aria-hidden="true"/></div>
                <span className="counter-number">
                  {window.appContent.InviteUser}
                  <br/>
                  {window.appState.AccountName}
                  <div style={{marginTop: 50}}></div>
                </span>
              </div>
            </div>
          </div>: null)
}, false, true, 2);