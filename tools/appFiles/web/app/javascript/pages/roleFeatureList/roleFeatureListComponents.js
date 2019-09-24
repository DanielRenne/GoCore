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
import {DeleteIcon, CopyIcon, EditIcon, ExportIcon, DownloadIcon} from "../../globals/icons";


class RoleFeatureList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.confirmDeleteAllRowsRef;
    this.confirmDeleteRowRef;
    this.widgetList;
    this.createComponentEvents();
  }

  createComponentEvents() {
    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "roleFeatureModify"});
    };

    this.handleDeleteroleFeatureInline = (row) => {
      this.confirmDeleteRowRef.handleOpen(row);
    };

    this.handleExport = (rows) => {
      window.api.download({action: "ExportCSV", state: {RoleFeatures: rows}, controller: "roleFeatureList", fileName: "export.csv"});
    };

    this.handleDelete = (row) => {
      window.api.post({action:"DeleteRoleFeature", state: {RoleFeature: row}, controller: "roleFeatures"});
    };

    this.deleteRoleFeatures = () => {
        window.api.post({action: "DeleteRoleFeatures", state: this.state, controller:"roleFeatures"});
    };

    this.handleCopy = () => {
        window.api.post({action: "CopyRoleFeature", state: this.state, controller:"roleFeatures"});
    };

    this.handleDeleteAll = (rows) => {
      this.setComponentState({DeletedRoleFeatures: rows}, () => {
        window.api.post({action: "DeleteManyRoleFeatures", state: this.state, controller:"roleFeatures", callback: (vm) => {
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
    return (

        <div>
          <WidgetList
              ref={(c) => this.widgetList = c}
              {...this.globs.widgetListDefaults()}
              name="roleFeatureList"
              listViewModel={this.state.WidgetList}
              controller="roleFeatureList"
              listTitle={this.globs.translate(this.state.WidgetList.ListTitle)}
              checkboxButtons={[
                {
                  func: (rows) => {
                    return true
                  },
                  button: <RaisedButton
                    label={window.appContent.WidgetListDeleteAll}
                    labelColor={blueGrey900}
                    onTouchTap={this.handleDeleteConfirmation}
                    icon={<DeleteIcon color={red500}/>}
                  />
                },
                {
                  func: (rows) => {
                    return true
                  },
                  button: <RaisedButton
                    label={window.appContent.Export}
                    labelColor={blueGrey900}
                    onTouchTap={this.handleExport}
                    icon={<DownloadIcon color={blueGrey500}/>}
                  />
                }
              ]}

              fields={[
                {
                  tooltip: "Role",
                  headerDisplay: "Role",
                  sortable: true,
                  stateKey: "Joins.Role.Name"
                },
                {
                  func: (row, currentValue) => {
                    if (currentValue != "") {
                      return <a href={"javascript:"} onClick={(e) => {
                        e.stopPropagation();
                        this.widgetList.buildSearchRequest({CustomCriteria: JSON.stringify({FeatureKey: row.Joins.Feature.Key})});
                      }}>{row.Joins.Feature.Key}</a>
                    }
                    return "";
                  },
                  tooltip: "Feature Key",
                  headerDisplay: "Feature Key",
                  sortable: true,
                  stateKey: "Joins.Feature.Key"
                },

                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdated,
                  stateKey: "Views.UpdateFromNow",
                  sortOn: "UpdateDate",
                  tooltipKey: "Views.UpdateDate"
                },
                {
                  tooltip: window.appContent.GlobalListsDateOfLastUpdate,
                  headerDisplay: window.appContent.GlobalListsUpdatedBy,
                  sortable: false,
                  stateKey: "Joins.LastUpdateUser.Views.FullName"
                },

              ]}
              rowButtons={[
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.EditRoleFeature}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.CopyRoleFeature}
                    onTouchTap={this.handleCopy}
                    icon={<CopyIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.appContent.GlobalButtonsDeleteThisRecord}
                    onTouchTap={this.handleDeleteroleFeatureInline}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
             ]}
             dataKey="RoleFeatures"
             data={this.state[this.state.WidgetList.DataKey]}
             addRecordOnClick={() => this.globs.FloatingActionButtonClick(null, () => this.globs.clickCurrentAddOrImportActionButton(), "AddImport", "CONTROLLER_ROLEFEATUREADD") }
             addRecordOnClickToolTip={window.pageContent.AddRoleFeature}
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

export default AddOrImportPage(RoleFeatureList, "CONTROLLER_ROLEFEATUREADD", "AddRoleFeature", "roleFeature_import.csv", "roleFeatureList", "ImportCSV");
