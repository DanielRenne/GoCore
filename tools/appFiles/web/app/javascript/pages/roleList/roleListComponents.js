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
    AppBar,
    FileUpload,
    Slider,
    Editor,
    EditorState,
    RichUtils,
    DraftJsStyleMap,
    DraftJsGetBlockStyle,
    DraftJsStyleButton,
    DraftJsBlockStyleControls,
    DraftJsInlineStyleControls,
    ContentState,
    convertFromHTML,
    stateToHTML
} from "../../globals/forms";
import {DeleteIcon, CopyIcon, EditIcon, ExportIcon, DownloadIcon, InfoIcon} from "../../globals/icons";


class RoleList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.confirmDeleteAllRowsRef;
    this.confirmDeleteRowRef;
    this.createComponentEvents();
  }

  componentDidMount() {
    window.clickCurrentAddOrImportActionButton = () => {
      this.open();
    }
  }

  createComponentEvents() {
    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "roleModify"});
    };

    this.handleView = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id, ReadOnly:true}, controller: "roleModify"});
    };

    this.handleDeleteroleInline = (row) => {
      this.confirmDeleteRowRef.handleOpen(row);
    };

    this.handleExport = (rows) => {
      window.api.download({action: "ExportCSV", state: {Roles: rows}, controller: "roleList", fileName: "export.csv"});
    };

    this.handleDelete = (row) => {
      window.api.post({action:"DeleteRole", state: {Role: row}, controller: "roles"});
    };

    this.deleteRoles = () => {
        window.api.post({action: "DeleteRoles", state: this.state, controller:"roles"});
    };

    this.handleCopy = (row) => {
        window.api.post({action: "CopyRole", state: {Role: row}, controller:"roles"});
    };

    this.handleDeleteAll = (rows) => {
      this.setComponentState({DeletedRoles: rows}, () => {
        window.api.post({action: "DeleteManyRoles", state: this.state, controller:"roles", callback: (vm) => {
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
    if (window.appState.DeveloperMode) {
      this.createComponentEvents();
    }
    return true;
  }

  render() {
    this.logRender();

    var fields = [];
    fields.push({
                  tooltip: window.pageContent.RoleListToolTipName,
                  headerDisplay: window.pageContent.RoleListHeaderName,
                  sortable: true,
                  stateKey: "Name",
                  responsiveKeep: true
                });
    if (this.globs.IsMasterAccount()) {
      fields.push({
                    tooltip: "Account Type Visible",
                    headerDisplay: "Account Type Visible",
                    sortable: true,
                    stateKey: "AccountType"
                  });
    }
    fields.push({
                  tooltip: window.pageContent.EnabledFeatures,
                  headerDisplay: window.pageContent.EnabledFeatures,
                  sortable: true,
                  stateKey: "Joins.RoleFeatures.Count"
                });
    fields.push({
                  func: (row) => {
                    if (!row.CanDelete) {
                      return window.pageContent.System
                    } else {
                      return window.pageContent.CurrentAccount
                    }
                  },
                  tooltip: window.pageContent.Owner,
                  headerDisplay: window.pageContent.Owner,
                  sortable: true,
                  stateKey: "Joins.RoleFeatures.Count"
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
              name="roleList"
              listViewModel={this.state.WidgetList}
              controller="roleList"
              listTitle={this.globs.translate(this.state.WidgetList.ListTitle)}
              showCheckboxes={false}
              fields={fields}
              rowButtons={[
                {
                  func: (row) => {
                    if (!this.globs.HasRole("ROLE_MODIFY") || (!row.CanDelete && !this.globs.IsMasterAccount())) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.EditRole}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (rows) => {
                    if (!this.globs.HasRole("ROLE_VIEW")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.appContent.ViewDetails}
                    labelColor={blueGrey900}
                    onTouchTap={this.handleView}
                    icon={<InfoIcon color={blueGrey500}/>}
                  />
                },
                {
                  func: (row) => {
                    if (!this.globs.HasRole("ROLE_COPY")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.CopyRole}
                    onTouchTap={this.handleCopy}
                    icon={<CopyIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    if (!this.globs.HasRole("ROLE_DELETE") || !row.CanDelete) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.appContent.GlobalButtonsDeleteThisRecord}
                    onTouchTap={this.handleDeleteroleInline}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
             ]}
             dataKey="Roles"
             data={this.state[this.state.WidgetList.DataKey]}
             searchEnabled={false}
             addRecordOnClick={this.globs.HasRole("ROLE_ADD") ? () => this.globs.FloatingActionButtonClick(null, false, "AddImport", "CONTROLLER_ROLEADD"): null}
             addRecordOnClickToolTip={window.pageContent.AddRole}
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

export default AddRecordPage(RoleList, "CONTROLLER_ROLEADD", "AddRole");
