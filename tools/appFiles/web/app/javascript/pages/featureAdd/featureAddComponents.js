import {
  React,
  CenteredPaperGrid,
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
import FeatureModify from "../featureModify/featureModifyComponents";

class FeatureAdd extends FeatureModify {
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
                floatingLabelText={"* " + window.pageContent.FeatureAddKey}
                hintText={"* " + window.pageContent.FeatureAddKey}
                fullWidth={true}
                onChange={this.handleKeyChange}
                errorText={this.globs.translate(this.state.Feature.Errors.Key)}
                value={this.state.Feature.Key}
                defaultValue={this.state.Feature.Key}
            />
            <br />
            <TextField
                floatingLabelText={"* " + window.pageContent.FeatureAddName}
                hintText={"* " + window.pageContent.FeatureAddName}
                fullWidth={true}
                onChange={this.handleNameChange}
                errorText={this.globs.translate(this.state.Feature.Errors.Name)}
                defaultValue={this.state.Feature.Name}
            />
            <br />
            <TextField
                floatingLabelText={"* " + window.pageContent.FeatureAddDescription}
                hintText={"* " + window.pageContent.FeatureAddDescription}
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
            label={window.pageContent.CreateFeature}
            onTouchTap={this.createFeature}
            secondary={true}
            id="save"
        />
        </CenteredPaperGrid>
    );
  }
}

export default BackPage(FeatureAdd);
