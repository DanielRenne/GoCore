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
import AppErrorModify from '../appErrorModify/appErrorModifyComponents';

class AppErrorAdd extends AppErrorModify {
  constructor(props, context) {
    super(props, context);

    this.createAppError = () => {
        window.api.post({action: "CreateAppError", state: this.state, controller:"appErrors"});
    };
  }

  componentWillReceiveProps(nextProps) {
    return true;
  }

  render() {
    this.logRender();
    return (
      <CenteredPaperGrid>

            <TextField
                floatingLabelText={window.pageContent.AppErrorAddClientSide}
                hintText={window.pageContent.AppErrorAddClientSide}
                fullWidth={true}
                onChange={this.handleClientSideChange}
                errorText={this.globs.translate(this.state.AppError.Errors.ClientSide)}
                defaultValue={this.state.AppError.ClientSide}
            />
            <br />
            <TextField
                floatingLabelText={window.pageContent.AppErrorAddUrl}
                hintText={window.pageContent.AppErrorAddUrl}
                fullWidth={true}
                onChange={this.handleUrlChange}
                errorText={this.globs.translate(this.state.AppError.Errors.Url)}
                defaultValue={this.state.AppError.Url}
            />
            <br />
            <TextField
                floatingLabelText={window.pageContent.AppErrorAddMessage}
                hintText={window.pageContent.AppErrorAddMessage}
                fullWidth={true}
                onChange={this.handleMessageChange}
                errorText={this.globs.translate(this.state.AppError.Errors.Message)}
                defaultValue={this.state.AppError.Message}
            />
            <br />
            <TextField
                floatingLabelText={window.pageContent.AppErrorAddStackShown}
                hintText={window.pageContent.AppErrorAddStackShown}
                fullWidth={true}
                onChange={this.handleStackShownChange}
                errorText={this.globs.translate(this.state.AppError.Errors.StackShown)}
                defaultValue={this.state.AppError.StackShown}
            />
            <br />
        <br />
        <br />
        <RaisedButton
            label={window.pageContent.CreateAppError}
            onTouchTap={this.createAppError}
            secondary={true}
            id="save"
        />
        </CenteredPaperGrid>
    );
  }
}

export default BackPage(AppErrorAdd);
