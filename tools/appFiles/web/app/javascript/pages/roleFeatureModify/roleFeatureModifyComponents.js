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

class RoleFeatureModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.createComponentEvents();
  }

  createComponentEvents() {

    this.handleRoleIdChange = (event) => {
      this.setComponentState({RoleFeature: {
        RoleId: event.target.value,
        Errors: {RoleId: ""}
      }});
    };

    this.handleRoleIdSelectChange = (event, index, value) => {
      this.setComponentState({RoleFeature: {
        RoleId: value,
        Errors: {RoleId: ""}
      }});
    };

    this.handleFeatureIdSelectChange = (event, index, value) => {
      this.setComponentState({RoleFeature: {
        FeatureId: value,
        Errors: {FeatureId: ""}
      }});
    };

    this.handleFeatureIdChange = (event) => {
      this.setComponentState({RoleFeature: {
        FeatureId: event.target.value,
        Errors: {FeatureId: ""}
      }});
    };

    this.handleBootstrapMetaChange = (event) => {
      this.setComponentState({RoleFeature: {
        BootstrapMeta: event.target.value,
        Errors: {BootstrapMeta: ""}
      }});
    };



    this.save = () => {
        window.api.post({action: "UpdateRoleFeatureDetails", state: this.state, controller:"roleFeatures"});
    };

    this.createRoleFeature = () => {
        window.api.post({action: "CreateRoleFeature", state: this.state, controller:"roleFeatures"});
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
        <div>
            <CenteredPaperGrid>
                <SelectField
                  floatingLabelText={"* " + window.pageContent.RoleFeatureModifyRoleId}
                  hintText={"* " + window.pageContent.RoleFeatureModifyRoleId}
                  fullWidth={true}
                  onChange={this.handleRoleIdSelectChange}
                  errorText={this.globs.translate(this.state.RoleFeature.Errors.RoleId)}
                  value={this.state.RoleFeature.RoleId}
                >
                  {roles}
                </SelectField>
                <br />
                <SelectField
                  floatingLabelText={"* " + window.pageContent.RoleFeatureModifyFeatureId}
                  hintText={"* " + window.pageContent.RoleFeatureModifyFeatureId}
                  fullWidth={true}
                  onChange={this.handleFeatureIdSelectChange}
                  errorText={this.globs.translate(this.state.RoleFeature.Errors.FeatureId)}
                  value={this.state.RoleFeature.FeatureId}
                >
                  {features}
                </SelectField>
                <br />
                
                <TextField
                  floatingLabelText={window.pageContent.RoleFeatureModifyBootstrapMeta}
                  hintText={window.pageContent.RoleFeatureModifyBootstrapMeta}
                  fullWidth={true}
                  onChange={this.handleBootstrapMetaChange}
                  errorText={this.globs.translate(this.state.RoleFeature.Errors.BootstrapMeta)}
                  defaultValue={this.state.RoleFeature.BootstrapMeta}
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

export default BackPage(RoleFeatureModify);
