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

//todo finish up required field *
//todo mass delete and delete inline.
//

class FileObjectModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.createComponentEvents();
  }

  createComponentEvents() {

    this.handleNameChange = (event) => {
      this.setComponentState({FileObject: {
        Name: event.target.value,
        Errors: {Name: ""}
      }});
    };


    this.handleContentChange = (event) => {
      this.setComponentState({FileObject: {
        Content: event.target.value,
        Errors: {Content: ""}
      }});
    };


    this.handleSizeChange = (event) => {
      this.setComponentState({FileObject: {
        Size: event.target.value,
        Errors: {Size: ""}
      }});
    };


    this.handleTypeChange = (event) => {
      this.setComponentState({FileObject: {
        Type: event.target.value,
        Errors: {Type: ""}
      }});
    };


    this.handleModifiedUnixChange = (event) => {
      this.setComponentState({FileObject: {
        ModifiedUnix: event.target.value,
        Errors: {ModifiedUnix: ""}
      }});
    };


    this.handleModifiedChange = (event) => {
      this.setComponentState({FileObject: {
        Modified: event.target.value,
        Errors: {Modified: ""}
      }});
    };


    this.handleMD5Change = (event) => {
      this.setComponentState({FileObject: {
        MD5: event.target.value,
        Errors: {MD5: ""}
      }});
    };



    this.save = () => {
      window.api.post({action: "UpdateFileObjectDetails", state: this.state, controller:"fileObjects"});
    };

    this.createFileObject = () => {
      window.api.post({action: "CreateFileObject", state: this.state, controller:"fileObjects"});
    };
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
            <CenteredPaperGrid>

                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifyName}
                  hintText={window.pageContent.FileObjectModifyName}
                  fullWidth={true}
                  onChange={this.handleNameChange}
                  errorText={this.globs.translate(this.state.FileObject.Errors.Name)}
                  value={this.state.FileObject.Name}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifyContent}
                  hintText={window.pageContent.FileObjectModifyContent}
                  fullWidth={true}
                  onChange={this.handleContentChange}
                  errorText={this.globs.translate(this.state.FileObject.Errors.Content)}
                  value={this.state.FileObject.Content}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifySize}
                  hintText={window.pageContent.FileObjectModifySize}
                  fullWidth={true}
                  onChange={this.handleSizeChange}
                  errorText={this.globs.translate(this.state.FileObject.Errors.Size)}
                  value={this.state.FileObject.Size}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifyType}
                  hintText={window.pageContent.FileObjectModifyType}
                  fullWidth={true}
                  onChange={this.handleTypeChange}
                  errorText={this.globs.translate(this.state.FileObject.Errors.Type)}
                  value={this.state.FileObject.Type}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifyModifiedUnix}
                  hintText={window.pageContent.FileObjectModifyModifiedUnix}
                  fullWidth={true}
                  onChange={this.handleModifiedUnixChange}
                  errorText={this.globs.translate(this.state.FileObject.Errors.ModifiedUnix)}
                  value={this.state.FileObject.ModifiedUnix}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifyModified}
                  hintText={window.pageContent.FileObjectModifyModified}
                  fullWidth={true}
                  onChange={this.handleModifiedChange}
                  errorText={this.globs.translate(this.state.FileObject.Errors.Modified)}
                  value={this.state.FileObject.Modified}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.FileObjectModifyMD5}
                  hintText={window.pageContent.FileObjectModifyMD5}
                  fullWidth={true}
                  onChange={this.handleMD5Change}
                  errorText={this.globs.translate(this.state.FileObject.Errors.MD5)}
                  value={this.state.FileObject.MD5}
                />
                <br />
                <br />
                <RaisedButton
                    label={window.appContent.SaveChanges}
                    onTouchTap={this.save}
                    secondary={true}
                />
            </CenteredPaperGrid>
        </div>
    );
  }
}

export default BackPage(FileObjectModify);
