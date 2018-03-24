import React, {Component} from 'react';
import OpenIcon from 'material-ui/svg-icons/action/input';
import {grey500, blueGrey500} from 'material-ui/styles/colors';
import CommunicationVpnKey from 'material-ui/svg-icons/communication/vpn-key';
import FileFileDownload from 'material-ui/svg-icons/file/file-download';
import AddOrImportPage from '../../components/addOrImportComponent';
import BasePageComponent from '../../components/basePageComponent';
import WidgetList from '../../components/widgetList';
import {blueGrey400, blueGrey900} from 'material-ui/styles/colors';
import RaisedButton from 'material-ui/RaisedButton';
import {EditIcon, LoginAs} from "../../globals/icons";


class AccountList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.ImportUsers;

    this.handleSetAccount = (row) => {
      window.globals.handleAccountSwitch(row.Id);
    };

    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "accountModify"});
    };

    this.handleExport = (rows) => {
      window.api.download({action: "ExportCSV", state: {Accounts: rows}, controller: "accountList", fileName: "export.csv"});
    };

  }


  render() {
    this.logRender();
    var fields = [];

    fields.push({
                  func: (row, currentValue) => {
                    var currentAccount = grey500;

                   if (window.appState.AccountId == row.Id) {
                     currentAccount = blueGrey500;
                   }
                   this.row = row;
                    return <span>{currentValue}</span>
                  },
                  tooltip: window.pageContent.AccountName,
                  headerDisplay: window.pageContent.Name,
                  stateKey: "AccountName",
                  sortOn: "AccountName",
                  responsiveKeep: true
                });
    if (window.appState.AccountTypeShort == 'atl') {
      fields.push({
                    tooltip: "Related Account",
                    headerDisplay: "Related Account",
                    stateKey: "Joins.RelatedAccount.AccountName"
                  });
    }

    fields.push({
                  tooltip: window.pageContent.AccountTypeShort,
                  headerDisplay: window.pageContent.AccountTypeShort,
                  stateKey: "AccountTypeLong"
                });
    fields.push({
                  func: (row, currentValue) => {
                    return <span>{row.City + ", " + row.Region}</span>
                  },
                  tooltip: window.pageContent.Region,
                  headerDisplay: window.pageContent.Region,
                  sortOn: "City",
                  stateKey: "Region"
                });
    fields.push({
                  func: (row, currentValue) => {
                    return <a href={"tel:" + row.PrimaryPhone.Value}>{row.PrimaryPhone.Value}</a>
                  },
                  tooltip: window.pageContent.Phone,
                  headerDisplay: window.pageContent.Phone,
                  stateKey: "PrimaryPhone.Value"
                });
    fields.push({
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdated,
                  stateKey: "Views.UpdateFromNow",
                  sortOn: "UpdateDate",
                  tooltipKey: "Views.UpdateDate"
                });
    fields.push({
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdatedBy,
                  sortable: false,
                  stateKey: "Joins.LastUpdateUser.Views.FullName"
                });


    return (

        <div>
          <WidgetList
              {...this.globs.widgetListDefaults()}
              {...this.globs.widgetListButtonBarOffset()}
              name="accountList"
              listViewModel={this.state.WidgetList}
              controller="accountList"
              listTitle={this.globs.translate(this.state.WidgetList.ListTitle)}
              checkboxButtons={[
                {
                  func: (rows) => {
                    if (!this.globs.HasRole("ACCOUNT_EXPORT")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    label={window.appContent.Export}
                    labelColor={blueGrey900}
                    onTouchTap={this.handleExport}
                    icon={<FileFileDownload color={blueGrey500}/>}
                  />
                }
              ]}

              fields={fields}
              rowButtons={[
                {
                  func: (row) => {
                    if (!this.globs.HasRole("ACCOUNT_MODIFY")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.EditAccountTooltip}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                }
             ]}
             dataKey={this.state.WidgetList.DataKey}
             data={this.state[this.state.WidgetList.DataKey]}
             addRecordOnClick={this.globs.HasRole("ACCOUNT_ADD") ? () => this.globs.FloatingActionButtonClick(null, () => this.globs.clickCurrentAddOrImportActionButton(), "AddImport", "CONTROLLER_ACCOUNTADD"): null}
             addRecordOnClickToolTip={window.pageContent.AddAccountTooltip}
          />
        </div>
    );
  }
}

export default AddOrImportPage(AccountList, "CONTROLLER_ACCOUNTADD", "AddAccountTooltip", "account_import.csv", "accountList", "ImportCSV");
