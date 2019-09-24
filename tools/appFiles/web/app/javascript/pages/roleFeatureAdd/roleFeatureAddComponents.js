import {
  React,
  CenteredPaperGrid,
  BackPage,
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
import RoleFeatureModify from "../roleFeatureModify/roleFeatureModifyComponents";

class RoleFeatureAdd extends RoleFeatureModify {
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

    var roles = this.globs.map(this.state.Roles, (r) => {
      return <MenuItem
        key={r.Id}
        value={r.Id}
        primaryText={r.Name}
        />;
    });

    var features = this.globs.map(this.state.Features, (r) => {
      return <MenuItem
        key={r.Id}
        value={r.Id}
        primaryText={r.Description}
        />;
    });
    
    return (
        <CenteredPaperGrid>
                <SelectField
                  floatingLabelText={"* " + window.pageContent.RoleFeatureAddRoleId}
                  hintText={"* " + window.pageContent.RoleFeatureAddRoleId}
                  fullWidth={true}
                  onChange={this.handleRoleIdSelectChange}
                  errorText={this.globs.translate(this.state.RoleFeature.Errors.RoleId)}
                  value={this.state.RoleFeature.RoleId}
                >
                  {roles}
                </SelectField>
                <br />
                <SelectField
                  floatingLabelText={"* " + window.pageContent.RoleFeatureAddFeatureId}
                  hintText={"* " + window.pageContent.RoleFeatureAddFeatureId}
                  fullWidth={true}
                  onChange={this.handleFeatureIdSelectChange}
                  errorText={this.globs.translate(this.state.RoleFeature.Errors.FeatureId)}
                  value={this.state.RoleFeature.FeatureId}
                >
                  {features}
                </SelectField>
            <br />
            <br />
            <br />
            <RaisedButton
                label={window.pageContent.CreateRoleFeature}
                onTouchTap={this.createRoleFeature}
                secondary={true}
                id="save"
            />
        </CenteredPaperGrid>
    );
  }
}

export default BackPage(RoleFeatureAdd);
