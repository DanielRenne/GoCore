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
    red50, red100, red200, red300, red400, red500, red600, red700, red800, red900, pink50, pink100, pink200, pink300, pink400, pink500, pink600, pink700, pink800, pink900, purple50, purple100, purple200, purple300, purple400, purple500, purple600, purple700, purple800, purple900,deepPurple50, deepPurple100, deepPurple200, deepPurple300, deepPurple400, deepPurple500, deepPurple600, deepPurple700, deepPurple800, deepPurple900, indigo50, indigo100, indigo200, indigo300, indigo400, indigo500, indigo600, indigo700, indigo800, indigo900, blue50, blue100, blue200, blue300, blue400, blue500, blue600, blue700, blue800, blue900, lightBlue50, lightBlue100, lightBlue200, lightBlue300, lightBlue400, lightBlue500, lightBlue600, lightBlue700, lightBlue800, lightBlue900, cyan50, cyan100, cyan200, cyan300, cyan400, cyan500, cyan600, cyan700, cyan800, cyan900,  teal50, teal100, teal200, teal300, teal400, teal500, teal600, teal700, teal800, teal900, green50, green100, green200, green300, green400, green500, green600, green700, green800, green900, lightGreen50, lightGreen100, lightGreen200, lightGreen300, lightGreen400, lightGreen500, lightGreen600, lightGreen700, lightGreen800, lightGreen900, lime50, lime100, lime200, lime300, lime400, lime500, lime600, lime700, lime800, lime900, yellow50, yellow100, yellow200, yellow300, yellow400, yellow500, yellow600, yellow700, yellow800, yellow900, amber50, amber100, amber200, amber300, amber400, amber500, amber600, amber700, amber800, amber900,orange50, orange100, orange200, orange300, orange400, orange500, orange600, orange700, orange800, orange900, deepOrange50, deepOrange100, deepOrange200, deepOrange300, deepOrange400, deepOrange500, deepOrange600, deepOrange700, deepOrange800, deepOrange900, brown50, brown100, brown200, brown300, brown400, brown500, brown600, brown700, brown800, brown900, grey50, grey100, grey200, grey300, grey400, grey500, grey600, grey700, grey800, grey900, blueGrey50, blueGrey100, blueGrey200, blueGrey300, blueGrey400, blueGrey500, blueGrey600, blueGrey700, blueGrey800, blueGrey900,
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


class -CAPCAMEL-List extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.confirmDeleteAllRowsRef;
    this.confirmDeleteRowRef;
    this.widgetList;
    this.createComponentEvents();
  }

  createComponentEvents() {
    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "-CAMEL-Modify"});
    };

    this.handleDelete-CAMEL-Inline = (row) => {
      this.confirmDeleteRowRef.handleOpen(row);
    };

    this.handleExport = (rows) => {
      window.api.download({action: "ExportCSV", state: {-CAPCAMELPLURAL-: rows}, controller: "-CAMEL-List", fileName: "export.csv"});
    };

    this.handleDelete = (row) => {
      window.api.post({action:"Delete-CAPCAMEL-", state: {-CAPCAMEL-: row}, controller: "-CAMELPLURAL-"});
    };

    this.delete-CAPCAMELPLURAL- = () => {
      window.api.post({action: "Delete-CAPCAMELPLURAL-", state: this.state, controller:"-CAMELPLURAL-"});
    };

    this.handleCopy = (row) => {
      window.api.post({action: "Copy-CAPCAMEL-", state: {-CAPCAMEL-: row}, controller:"-CAMELPLURAL-"});
    };

    this.handleDeleteAll = (rows) => {
      this.setComponentState({Deleted-CAPCAMELPLURAL-: rows}, () => {
        window.api.post({action: "DeleteMany-CAPCAMELPLURAL-", state: this.state, controller:"-CAMELPLURAL-", callback: (vm) => {
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
              ref={(c) => this.widgetList = c}
              -PAGE_RENDER_EXTRA-
              name="-CAMEL-List"
              listViewModel={this.state.WidgetList}
              controller="-CAMEL-List"
              listTitle={this.globs.translate(this.state.WidgetList.ListTitle)}
              checkboxButtons={[
                {
                  func: (rows) => {
                    if (!this.globs.HasRole("-CAPITALIZED-_DELETE")) {
                      return false
                    }
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
                    if (!this.globs.HasRole("-CAPITALIZED-_EXPORT")) {
                      return false
                    }
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
                  -LIST_FIELDS-
              ]}
              rowButtons={[
                {
                  func: (row) => {
                    if (!this.globs.HasRole("-CAPITALIZED-_MODIFY")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.Edit-CAPCAMEL-}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    if (!this.globs.HasRole("-CAPITALIZED-_COPY")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.Copy-CAPCAMEL-}
                    onTouchTap={this.handleCopy}
                    icon={<CopyIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    if (!this.globs.HasRole("-CAPITALIZED-_DELETE")) {
                      return false
                    }
                    return true
                  },
                  button: <RaisedButton
                    title={window.appContent.GlobalButtonsDeleteThisRecord}
                    onTouchTap={this.handleDelete-CAMEL-Inline}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
             ]}
             dataKey={this.state.WidgetList.DataKey}
             data={this.state[this.state.WidgetList.DataKey]}
             addRecordOnClick={this.globs.HasRole("-CAPITALIZED-_ADD") ? () => this.globs.FloatingActionButtonClick(null, () => this.globs.clickCurrentAddOrImportActionButton(), "AddImport", "CONTROLLER_-CAPITALIZED-ADD"): null}
             addRecordOnClickToolTip={window.pageContent.Add-CAPCAMEL-}
             -PAGE_RENDER_HEIGHT_WIDGETLIST-
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

export default AddOrImportPage(-CAPCAMEL-List, "CONTROLLER_-CAPITALIZED-ADD", "Add-CAPCAMEL-", "-CAMEL-_import.csv", "-CAMEL-List", "ImportCSV");
