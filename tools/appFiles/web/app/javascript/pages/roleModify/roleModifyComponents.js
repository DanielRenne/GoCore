import {
  React,
  CenteredPaperGrid,
  BasePageComponent,
  BackPage,
  TextField,
  RaisedButton,
  RadioButton,
  RadioButtonGroup,
  Toggle,
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

class RoleModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);
    this.createComponentEvents();
  }

  createComponentEvents() {

    this.handleRoleToggle = (e, key, value) => {
      var dict = {};
      dict[key] = value;
      this.setComponentState({FeaturesEnabled: dict});
    };

    this.handleNameChange = (event) => {
      this.setComponentState({Role: {
        Name: event.target.value,
        Errors: {Name: ""}
      }});
    };


    this.handleCanDeleteChange = (event) => {
      this.setComponentState({Role: {
        CanDelete: event.target.value,
        Errors: {CanDelete: ""}
      }});
    };


    this.handleShortNameChange = (event) => {
      this.setComponentState({Role: {
        ShortName: event.target.value,
        Errors: {ShortName: ""}
      }});
    };

    this.buildRoleChecks = () => {
      return this.globs.map(this.state.FeatureGroups, (group) => {
              return <div><CenteredPaperGrid><h3>{group.Name}</h3>{this.globs.map(group.Joins.Features.Items, (feature) => {
                return <span>
                  {this.uriParams.hasOwnProperty("ReadOnly") ? <span>{feature.Name + ": " + ((this.state.FeaturesEnabled[feature.Id]) ? window.appContent.Yes : window.appContent.No)}</span> :
                  <Toggle
                    key={feature.Id}
                    label={feature.Name}
                    defaultToggled={this.state.FeaturesEnabled[feature.Id]}
                    labelPosition="right"
                    onToggle={(e, value) => this.handleRoleToggle(e, feature.Id, value)}
                  />}
                  {group.Joins.Features.Items[group.Joins.Features.Items.length - 1] == feature ? null : <span>
                  <br />
                  <br /></span>}
                </span>
              })}</CenteredPaperGrid><br/><br/></div>;
            });
    };

    this.save = () => {
      var submit = () => {
        window.api.post({action: "UpdateRoleDetails", state: this.state, controller:"roles"});
      };

      if (!this.state.CanDelete && this.globs.IsMasterAccount() && !window.appState.DeveloperMode) {
        if (window.prompt("WARNING: Changing roles in production is extremely frowned upon (should only be done in emergency situations) and can effect all users in the system as well as affect our build of this software.  Please type UNDERSTAND to submit this change or cancel to stop") == "UNDERSTAND") {
          submit();
        }
      } else {
        submit();
      }
    };

    this.createRole = () => {
      window.api.post({action: "CreateRole", state: this.state, controller:"roles"});
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
                  disabled={this.uriParams.hasOwnProperty("ReadOnly")}
                  floatingLabelText={"* " + window.pageContent.RoleModifyName}
                  hintText={"* " + window.pageContent.RoleModifyName}
                  fullWidth={true}
                  onChange={this.handleNameChange}
                  errorText={this.globs.translate(this.state.Role.Errors.Name)}
                  value={this.state.Role.Name}
                />
            </CenteredPaperGrid>
            <br/>
            <br/>
            {this.buildRoleChecks()}
            {this.uriParams.hasOwnProperty("ReadOnly") ? null : <CenteredPaperGrid>
                <RaisedButton
                    label={window.appContent.SaveChanges}
                    onTouchTap={this.save}
                    secondary={true}
                />
            </CenteredPaperGrid>}
        </div>
    );
  }
}

export default BackPage(RoleModify);
