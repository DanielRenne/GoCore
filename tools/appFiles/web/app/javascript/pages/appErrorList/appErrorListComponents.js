import {
    React,
    CenteredPaperGrid,
    BasePageComponent,
    BaseComponent,
    WidgetList,
    AddRecordPage,
    BackPage,
    AddOrImportPage,
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
    ConfirmPopup,
    Grid,
    Row,
    Col,
    AutoComplete,
    PhoneInput,
    orange50, orange100, orange200, orange300, orange400, orange500, orange600, orange700, orange800, orange900, deepOrange50, deepOrange100, deepOrange200, deepOrange300, deepOrange400, deepOrange500, deepOrange600, deepOrange700, deepOrange800, deepOrange900,red50, red100, red200, red300, red400, red500, red600, red700, red800, red900, blueGrey50, blueGrey100, blueGrey200, blueGrey300, blueGrey400, blueGrey500, blueGrey600, blueGrey700, blueGrey800, blueGrey900, grey50, grey100, grey200, grey300, grey400, grey500, grey600, grey700, grey800, grey900, green50, green100, green200, green300, green400, green500, green600, green700, green800, green900, indigo50, indigo100, indigo200, indigo300, indigo400, indigo500, indigo600, indigo700, indigo800, indigo900,
    AppBar
} from "../../globals/forms";
import {DeleteIcon, EditIcon, ExportIcon} from "../../globals/icons";


class AppErrorList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);

    this.confirmDeleteAllRowsRef;
    this.confirmDeleteRowRef;

    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "appErrorModify"});
    };

    this.handleDeleteappErrorInline = (row) => {
      this.confirmDeleteRowRef.handleOpen(row);
    };

    this.handleDelete = (row) => {
      window.api.post({action:"DeleteAppError", state: {AppError: row}, controller: "appErrors"});
    };

    this.deleteAppErrors = () => {
        window.api.post({action: "DeleteAppErrors", state: this.state, controller:"appErrors"});
    };

    this.handleDeleteAll = (rows) => {
      this.setComponentState({DeletedAppErrors: rows}, () => {
        window.api.post({action: "DeleteManyAppErrors", state: this.state, controller:"appErrors", callback: (vm) => {
          if (window.appState.DialogOpen || window.appState.DialogGenericOpen) {
            this.confirmDeleteAllRowsRef.handleClose();
          }
        }});
      });
    };

    this.handleDeleteConfirmation = (rows) => {
      this.confirmDeleteAllRowsRef.handleOpen(rows);
    };

  }
  componentDidUpdate() {
    if (window.appState.DeveloperLogState) {
      console.log("componentDidUpdate", this.state);
    }
  }

  componentWillReceiveProps(nextProps) {
    return true;
  }

  render() {
    this.logRender();
    return (

        <div>
          <WidgetList
              {...this.globs.widgetListDefaults()}

              name="appErrorList"
              listViewModel={this.state.WidgetList}
              controller="appErrorList"
              listTitle={this.globs.translate(this.state.WidgetList.ListTitle)}
              checkboxButtons={[
                {
                  func: (rows) => {
                    return true
                  },
                  button: <RaisedButton
                    label={window.appContent.WidgetListDeleteAll}
                    labelColor={blueGrey900}
                    onTouchTap={this.handleDeleteAll}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
              ]}

              fields={[
                {
                  func: (row) => {
                    if (row.ClientSide) {
                      return "Yes"
                    } else {
                      return "No"
                    }
                  },
                  tooltip: window.pageContent.AppErrorListToolTipClientSide,
                  headerDisplay: window.pageContent.AppErrorListHeaderClientSide,
                  sortable: true,
                  stateKey: "ClientSide"
                },
                {
                  tooltip: window.pageContent.AppErrorListToolTipUrl,
                  headerDisplay: window.pageContent.AppErrorListHeaderUrl,
                  sortable: true,
                  stateKey: "Url"
                },
                {
                  func: (row) => {
                    if (row.Joins.hasOwnProperty("User")) {
                      return row.Joins.User.Views.FullName;
                    }
                    return "";
                  },
                  tooltip: "User",
                  headerDisplay: "User",
                  sortable: true,
                  stateKey: "ClientSide" // just a stub.  All logic in func
                },
                {
                  func: (row) => {
                    if (row.Joins.hasOwnProperty("Account")) {
                      var ptr = row.Joins.Account;
                      return ptr.AccountName + " (" + ptr.AccountTypeLong + ")";
                    }
                    return "";
                  },
                  tooltip: "Account",
                  headerDisplay: "Account",
                  sortable: true,
                  stateKey: "ClientSide" // just a stub.  All logic in func
                },
                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdated,
                  stateKey: "Views.UpdateFromNow",
                  sortOn: "UpdateDate",
                  tooltipKey: "Views.UpdateDate"
                }

              ]}
              rowButtons={[
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.EditAppError}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                } ,
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.appContent.GlobalButtonsDeleteThisRecord}
                    onTouchTap={this.handleDeleteappErrorInline}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
             ]}
             dataKey="AppErrors"
             data={this.state[this.state.WidgetList.DataKey]}
             addRecordOnClick={() => this.globs.FloatingActionButtonClick(null, () => this.globs.clickCurrentAddOrImportActionButton(), "AddImport", "CONTROLLER_APPERRORADD") }
             addRecordOnClickToolTip={window.pageContent.AddAppError}
             offsetHeightToList={92}
          />
          <ConfirmPopup
              onSubmit={this.handleDeleteAll}
              areYouSureMsg={window.pageContent.AreYouSure}
              ref={(component) => this.confirmDeleteAllRowsRef = component}/>

          <ConfirmPopup
              onSubmit={this.handleDelete}
              areYouSureMsg={window.pageContent.AreYouSureInline}
              ref={(component) => this.confirmDeleteRowRef = component}
              />

        </div>
    );
  }
}

export default AppErrorList;
