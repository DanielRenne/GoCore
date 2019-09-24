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
import FeatureGroupModify from '../featureGroupModify/featureGroupModifyComponents';

class FeatureGroupAdd extends FeatureGroupModify {
  constructor(props, context) {
    super(props, context);
    this.createComponentEvents();
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
      <CenteredPaperGrid>

            <TextField
                floatingLabelText={"* " + window.pageContent.FeatureGroupAddName}
                hintText={"* " + window.pageContent.FeatureGroupAddName}
                fullWidth={true}
                onChange={this.handleNameChange}
                errorText={this.globs.translate(this.state.FeatureGroup.Errors.Name)}
                defaultValue={this.state.FeatureGroup.Name}
            />
            <br />
            <br />
        {this.accountTypes()}

        <br />
        <br />
        <RaisedButton
            label={window.pageContent.CreateFeatureGroup}
            onTouchTap={this.createFeatureGroup}
            secondary={true}
            id="save"
        />
        </CenteredPaperGrid>
    );
  }
}

export default BackPage(FeatureGroupAdd);
