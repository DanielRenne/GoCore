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

//todo finish up required field *
//todo mass delete and delete inline.
//

class AppErrorModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);


    this.handleClientSideChange = (event) => {
      this.setComponentState({AppError: {
        ClientSide: event.target.value,
        Errors: {ClientSide: ""}
      }});
    };


    this.handleUrlChange = (event) => {
      this.setComponentState({AppError: {
        Url: event.target.value,
        Errors: {Url: ""}
      }});
    };


    this.handleMessageChange = (event) => {
      this.setComponentState({AppError: {
        Message: event.target.value,
        Errors: {Message: ""}
      }});
    };


    this.handleStackShownChange = (event) => {
      this.setComponentState({AppError: {
        StackShown: event.target.value,
        Errors: {StackShown: ""}
      }});
    };



    this.save = () => {
        window.api.post({action: "UpdateAppErrorDetails", state: this.state, controller:"appErrors"});
    };
  }

  componentWillReceiveProps(nextProps) {
    return true;
  }

  render() {
    this.logRender();
    return (
        <div>
            <CenteredPaperGrid>

                <TextField
                  floatingLabelText={window.pageContent.AppErrorModifyClientSide}
                  hintText={window.pageContent.AppErrorModifyClientSide}
                  fullWidth={true}
                  onChange={this.handleClientSideChange}
                  errorText={this.globs.translate(this.state.AppError.Errors.ClientSide)}
                  defaultValue={this.state.AppError.ClientSide}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.AppErrorModifyUrl}
                  hintText={window.pageContent.AppErrorModifyUrl}
                  fullWidth={true}
                  onChange={this.handleUrlChange}
                  errorText={this.globs.translate(this.state.AppError.Errors.Url)}
                  defaultValue={this.state.AppError.Url}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.AppErrorModifyMessage}
                  hintText={window.pageContent.AppErrorModifyMessage}
                  fullWidth={true}
                  onChange={this.handleMessageChange}
                  multiLine={true}
                  errorText={this.globs.translate(this.state.AppError.Errors.Message)}
                  defaultValue={this.state.AppError.Message}
                />
                <br />
                <TextField
                  floatingLabelText={window.pageContent.AppErrorModifyStackShown}
                  hintText={window.pageContent.AppErrorModifyStackShown}
                  fullWidth={true}
                  multiLine={true}
                  onChange={this.handleStackShownChange}
                  errorText={this.globs.translate(this.state.AppError.Errors.StackShown)}
                  defaultValue={this.state.AppError.StackShown}
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

export default BackPage(AppErrorModify);
