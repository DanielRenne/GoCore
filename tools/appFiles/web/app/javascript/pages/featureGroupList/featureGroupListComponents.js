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


class FeatureGroupList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.confirmDeleteAllRowsRef;
    this.confirmDeleteRowRef;
    this.createComponentEvents();
  }

  createComponentEvents() {
    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "featureGroupModify"});
    };

    this.handleDeletefeatureGroupInline = (row) => {
      this.confirmDeleteRowRef.handleOpen(row);
    };

    this.handleExport = (rows) => {
      window.api.download({action: "ExportCSV", state: {FeatureGroups: rows}, controller: "featureGroupList", fileName: "export.csv"});
    };

    this.handleDelete = (row) => {
      window.api.post({action:"DeleteFeatureGroup", state: {FeatureGroup: row}, controller: "featureGroups"});
    };

    this.deleteFeatureGroups = () => {
        window.api.post({action: "DeleteFeatureGroups", state: this.state, controller:"featureGroups"});
    };

    this.handleCopy = () => {
        window.api.post({action: "CopyFeatureGroup", state: this.state, controller:"featureGroups"});
    };

    this.handleDeleteAll = (rows) => {
      this.setComponentState({DeletedFeatureGroups: rows}, () => {
        window.api.post({action: "DeleteManyFeatureGroups", state: this.state, controller:"featureGroups", callback: (vm) => {
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
              {...this.globs.widgetListDefaults()}

              name="featureGroupList"
              listViewModel={this.state.WidgetList}
              controller="featureGroupList"
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
                  tooltip: "Id",
                  headerDisplay: "Id",
                  sortable: true,
                  stateKey: "Id"
                },
                {
                  tooltip: window.pageContent.FeatureGroupListToolTipName,
                  headerDisplay: window.pageContent.FeatureGroupListHeaderName,
                  sortable: true,
                  stateKey: "Name"
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
                    title={window.pageContent.EditFeatureGroup}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.CopyFeatureGroup}
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
                    onTouchTap={this.handleDeletefeatureGroupInline}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
             ]}
             dataKey="FeatureGroups"
             data={this.state[this.state.WidgetList.DataKey]}
             addRecordOnClick={() => this.globs.FloatingActionButtonClick(null, () => this.globs.clickCurrentAddOrImportActionButton(), "AddImport", "CONTROLLER_FEATUREGROUPADD") }
             addRecordOnClickToolTip={window.pageContent.AddFeatureGroup}
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

export default AddOrImportPage(FeatureGroupList, "CONTROLLER_FEATUREGROUPADD", "AddFeatureGroup", "featureGroup_import.csv", "featureGroupList", "ImportCSV");
