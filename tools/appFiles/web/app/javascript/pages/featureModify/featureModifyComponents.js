import {
  React,
  CenteredPaperGrid,
  BasePageComponent,
  BackPage,
  TextField,
  RaisedButton,
  RadioButton,
  RadioButtonGroup,
  SelectField,
  MenuItem,
  List,
  ListItem,
  Grid,
  Row,
  Col,
  Editor,
  EditorState,
  RichUtils,
  ContentState,
  convertFromHTML
} from "../../globals/forms";

//todo finish up required field *
//todo mass delete and delete inline.
//

class FeatureModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    console.error(this.state)
    this.createComponentEvents();
  }

  createComponentEvents() {

    this.handleKeyChange = (event) => {
      this.setComponentState({Feature: {
        Key: event.target.value.toUpperCase(),
        Errors: {Key: ""}
      }});
    };


    this.handleNameChange = (event) => {
      this.setComponentState({Feature: {
        Name: event.target.value,
        Errors: {Name: ""}
      }});
    };


    this.handleDescriptionChange = (event) => {
      this.setComponentState({Feature: {
        Description: event.target.value,
        Errors: {Description: ""}
      }});
    };


    this.handleFeatureGroupIdChange = (event, index, value) => {
      this.setComponentState({Feature: {
        FeatureGroupId: value,
        Errors: {FeatureGroupId: ""}
      }});
    };


    this.handleBootstrapMetaChange = (event) => {
      this.setComponentState({Feature: {
        BootstrapMeta: event.target.value,
        Errors: {BootstrapMeta: ""}
      }});
    };



    this.save = () => {
        window.api.post({action: "UpdateFeatureDetails", state: this.state, controller:"features"});
    };

    this.createFeature = () => {
        window.api.post({action: "CreateFeature", state: this.state, controller:"features"});
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
                  floatingLabelText={"* " + window.pageContent.FeatureModifyKey}
                  hintText={"* " + window.pageContent.FeatureModifyKey}
                  fullWidth={true}
                  onChange={this.handleKeyChange}
                  errorText={this.globs.translate(this.state.Feature.Errors.Key)}
                  value={this.state.Feature.Key}
                />
                <br />
                <TextField
                  floatingLabelText={"* " + window.pageContent.FeatureModifyName}
                  hintText={"* " + window.pageContent.FeatureModifyName}
                  fullWidth={true}
                  onChange={this.handleNameChange}
                  errorText={this.globs.translate(this.state.Feature.Errors.Name)}
                  defaultValue={this.state.Feature.Name}
                />
                <br />
                <TextField
                  floatingLabelText={"* " + window.pageContent.FeatureModifyDescription}
                  hintText={"* " + window.pageContent.FeatureModifyDescription}
                  fullWidth={true}
                  onChange={this.handleDescriptionChange}
                  errorText={this.globs.translate(this.state.Feature.Errors.Description)}
                  defaultValue={this.state.Feature.Description}
                />
                <br />
                <SelectField
                  floatingLabelText={"* " + window.pageContent.FeatureModifyFeatureGroupId}
                  hintText={"* " + window.pageContent.FeatureModifyFeatureGroupId}
                  value={this.state.Feature.FeatureGroupId}
                  onChange={this.handleFeatureGroupIdChange}
                  style={{width: 400}}
                  errorText={this.globs.translate(this.state.Feature.Errors.FeatureGroupId)}>
                  {this.globs.map(this.state.FeatureGroups, (v) => <MenuItem key={v.Id} value={v.Id} primaryText={v.Name}/>)}
                </SelectField>
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

export default BackPage(FeatureModify);
