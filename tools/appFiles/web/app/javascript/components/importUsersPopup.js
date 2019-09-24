import Dialog from 'material-ui/Dialog';
import {blueGrey500, blueGrey900} from 'material-ui/styles/colors';
import CommunicationVpnKey from 'material-ui/svg-icons/communication/vpn-key';
import Subheader from 'material-ui/Subheader';
import Menu from 'material-ui/Menu';
import {Grid, Row, Col} from 'react-flexgrid-no-bs-conflict';
import Paper from 'material-ui/Paper';
import {React,
  CenteredPaperGrid,
  BaseComponent,
  WidgetList,
  AddRecordPage,
  BackPage,
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
  Divider} from '../globals/forms';

class ImportUsers extends BaseComponent {
  constructor(props, context) {
    super(props, context);
    //this.eventHandlerExample = (event) => {
      //window.api.post({action: "Warning", state: this.state});
    //};

    this.state = {
      open: false,
      ComponentAccounts: null,
      ImportAccountRoles: null,
      CheckedAccountRoles: {}
    };

    this.handleOpen = (obj) => {
      this.setComponentState(obj, () => {
          this.setSelectedAccounts();
      });
    };

    this.setSelectedAccounts = () => {
      window.api.post({action: "GetAllAccountUsers", state: this.state, controller:"accounts", callback: (vm) => {
        this.setComponentState({open: true,
          ImportAccountRoles: window.pageState.ImportAccountRoles});
      }});
    };

    this.handleClose = () => {
      this.setComponentState({open: false});
    };

    // AccountRole.Id SET TO data("account") (id of AccountRole that got changed)
    this.handleRoleIdChange = (id, event, index, value) =>  {
      // var changedAccountroleId = $(event.target).parents().find("[data-account]").first().data("account")
      // var changedAccountroleId = id
      var currentAccountRoles = this.state.ImportAccountRoles
      this.setComponentState({ImportAccountRoles: {}}, () => {
        this.mergeAndUpdateCheckBoxState(id, true, () => {
          this.setComponentState({ImportAccountRoles: currentAccountRoles.map((acc) => {
              if (id == acc.Id) {
                 acc.RoleId = value
              }
              return acc
            })
          })
        })
      })
    };

    // Target is bool for checked/unchecked
    this.handleCheckBoxChange = (id, event, target) => {
      // console.warn(id);
      this.mergeAndUpdateCheckBoxState(id, target)
    };

    // Target is bool for checked/unchecked
    this.mergeAndUpdateCheckBoxState = (id, target, callback) => {
      var currCheckedRoles = this.state.CheckedAccountRoles;
      currCheckedRoles[id] = target
      this.setComponentState({CheckedAccountRoles: currCheckedRoles}, callback)
    };

    this.handleInviteAll = () => {

      window.api.post({action: "ImportSelectedUsers", state: this.state, controller:"accounts", callback:(vm) => {
        this.setComponentState({open: false,
          ImportAccountRoles: {},
          CheckedAccountRoles: {}})
      }})
    };

    this.handleInternalInvite = () => {
      window.closeAndClearImportUsers = () => {
        this.setComponentState({open: false,
          ImportAccountRoles: {}})
        window.closeAndClearImportUsers = "" // last line of callback
      };

      Object.filter = (obj, predicate) =>
          Object.keys(obj)
                .filter( key => predicate(obj[key]) )
                .reduce( (res, key) => Object.assign(res, { [key]: obj[key] }), {} )
      // removes the roles that got checked and then unchecked - false values
      var filteredCheckedAccRoles = Object.filter(this.state.CheckedAccountRoles, val => val == true)

      var ImportRoles = []
      for (var property in filteredCheckedAccRoles) {
        if(filteredCheckedAccRoles.hasOwnProperty(property)) {
          var objectPos = this.state.ImportAccountRoles.map((x) => {return x.Id; }).indexOf(property)
          var objectFound = this.state.ImportAccountRoles[objectPos]
          ImportRoles.push(objectFound)
        }
      }

      this.setComponentState({ImportAccountRoles: {}}, () => {
        this.setComponentState({ImportAccountRoles: ImportRoles,
          CheckedAccountRoles: {}}, () => window.api.post({action: "ImportSelectedUsers", state: this.state, controller:"accounts", callback: (vm) => {
            this.setComponentState({open: false,
              ImportAccountRoles: {},
              CheckedAccountRoles: {}})
          }}))
      });
    };

  }

  render() {
    try {
      this.logRender();
      const actions = [
        <RaisedButton
          label={window.pageContent.ButtonCancel}
          primary={true}
          onTouchTap={this.handleClose}
          style={{marginRight: 10}}
        />,
        <RaisedButton
          label={window.pageContent.ButtonImport}
          secondary={true}
          onTouchTap={this.handleInternalInvite}
          style={{marginRight: 10}}
        />,
        <RaisedButton
          label={window.pageContent.ButtonImportAll}
          secondary={true}
          onTouchTap={this.handleInviteAll}
        />,
      ];

      var rows="";
      var attemptkey = 1;
      if(this.state.open == true && this.state.ImportAccountRoles != null && this.state.ImportAccountRoles.length > 0) {
        rows = this.state.ImportAccountRoles.map((a) => {
          var accounts = a
          var roles = (
            <SelectField
              value={accounts.RoleId}
              onChange={this.handleRoleIdChange.bind(this, accounts.Id)}
              floatingLabelText={"* " + window.pageContent.ImportUsersSelectRole}
              hintText={"* " + window.pageContent.ImportUsersSelectRole}
              key={"select-" + accounts.Id}
              id={accounts.Id}
            >
              {this.props.roles.map((role) =>{
                return (
                  <MenuItem value={role.Id} primaryText={role.Name} key={role.Id + "-" + role.Id} data-account={accounts.Id}/>
                );
              })}
            </SelectField>
            );
            return(
              <Row key={accounts.Id}>
                <Col md={6}>
                  <Checkbox key={"check-" + accounts.Id} checked={this.state.CheckedAccountRoles[accounts.Id]} onCheck={this.handleCheckBoxChange.bind(this, accounts.Id)} label={accounts.Joins.User.Views.FullName} style={{paddingTop:'28px'}} />
                </Col>
                <Col md={6}>
                  {roles}
                </Col>
                <Divider />
              </Row>
            );
        });
      }

      if(this.state.ComponentAccounts != null) {
      var selAccounts = this.state.ComponentAccounts.map((acc) => {
            return(
              acc.AccountName
            );
          }).join(", ")
      }

      return (
        <div>

          <Dialog
            title={window.pageContent.TitleImportUsers}
            actions={actions}
            modal={false}
            open={this.state.open}
            onRequestClose={this.handleClose}
          >
            {window.pageContent.ImportUsersSelectAllUsers + selAccounts}


            <h4>{window.pageContent.ImportUsersAvailUsers}</h4>
            <Paper style={{margin:'15px', maxHeight: 396, overflowY: 'scroll'}} zDepth={1}>
              <Grid>
                {rows}
              </Grid>
            </Paper>
          </Dialog>
        </div>
      );
    } catch(e) {
      return this.globs.ComponentError("ImportUsers", e.message, e);
    }
  }
}

export default ImportUsers;
