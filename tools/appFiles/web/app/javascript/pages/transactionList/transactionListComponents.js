import {
    React,
    CenteredPaperGrid,
    BasePageComponent,
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
    Divider
} from "../../globals/forms";
import {blueGrey400, blueGrey900, red500} from "material-ui/styles/colors";
import {DeleteIcon, EditIcon, ExportIcon} from "../../globals/icons";


class TransactionList extends BasePageComponent {
  constructor(props, context) {
    super(props, context);

    //this.confirmDeleteAllRowsRef;

    this.handleEdit = (row) => {
      window.api.get({action:"Load", uriParams: {Id:row.id}, controller: "transactionModify"});
    };

    this.handleDelete = (row) => {
      alert("todo");
      //window.api.get({action:"Load", uriParams: {Id:row.Id}, controller: "TransactionDelete"});
    };

    this.deleteTransactions = () => {
        window.api.post({action: "DeleteTransactions", state: this.state, controller:"transactions"});
    };

    this.handleDeleteAll = (row) => {
      window.api.post({action: "DeleteTransactions", state: this.state, controller:"transactions", callback: (vm) => {
        if (window.appState.DialogOpen || window.appState.DialogGenericOpen) {
          this.confirmDeleteAllRowsRef.handleClose();
        }
      }});
    };

  }


  render() {
    this.logRender();
    return (

        <div>
          <WidgetList
              {...this.globs.widgetListDefaults()}
              {...this.globs.widgetListButtonBarOffset()}
              name="TransactionList"
              listViewModel={this.state.WidgetList}
              controller="transactionList"
              listTitle={this.globs.translate(this.state.WidgetList.ListTitle)}
              checkboxButtons={[
                {
                  func: (rows) => {
                    return true
                  },
                  button: <RaisedButton
                    label={window.appContent.WidgetListDeleteAll}
                    labelColor={blueGrey900}
                    onTouchTap={this.handleDeleteAll}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
              ]}

              fields={[
                {
                  tooltip: window.pageContent.TransactionListToolTipId,
                  headerDisplay: window.pageContent.TransactionListHeaderId,
                  stateKey: "id"
                },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipUserId,
                //   headerDisplay: window.pageContent.TransactionListHeaderUserId,
                //   stateKey: "userId"
                // },

                // {
                //   tooltip: window.appContent.GlobalListsDateOfCreation,
                //   headerDisplay: window.appContent.GlobalListsCreated,
                //   stateKey: "createDate"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipLastUpdate,
                //   headerDisplay: window.pageContent.TransactionListHeaderLastUpdate,
                //   stateKey: "lastUpdate"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipCompleteDate,
                //   headerDisplay: window.pageContent.TransactionListHeaderCompleteDate,
                //   stateKey: "completeDate"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipRollbackDate,
                //   headerDisplay: window.pageContent.TransactionListHeaderRollbackDate,
                //   stateKey: "rollbackDate"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipCommitted,
                //   headerDisplay: window.pageContent.TransactionListHeaderCommitted,
                //   stateKey: "committed"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipError,
                //   headerDisplay: window.pageContent.TransactionListHeaderError,
                //   stateKey: "error"
                // },
                {
                  tooltip: window.pageContent.TransactionListToolTipCollections,
                  headerDisplay: window.pageContent.TransactionListHeaderCollections,
                  stateKey: "collections"
                },
                {
                  tooltip: window.pageContent.TransactionListToolTipDetails,
                  headerDisplay: window.pageContent.TransactionListHeaderDetails,
                  stateKey: "details"
                },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipRolledBack,
                //   headerDisplay: window.pageContent.TransactionListHeaderRolledBack,
                //   stateKey: "rolledBack"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipRolledBackBy,
                //   headerDisplay: window.pageContent.TransactionListHeaderRolledBackBy,
                //   stateKey: "rolledBackBy"
                // },
                // {
                //   tooltip: window.pageContent.TransactionListToolTipRollbackReason,
                //   headerDisplay: window.pageContent.TransactionListHeaderRollbackReason,
                //   stateKey: "rollbackReason"
                // },

              ]}
              rowButtons={[
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.pageContent.EditTransaction}
                    onTouchTap={this.handleEdit}
                    icon={<EditIcon color={blueGrey400}/>}
                  />
                },
                {
                  func: (row) => {
                    return true
                  },
                  button: <RaisedButton
                    title={window.appContent.GlobalButtonsDeleteThisRecord}
                    onTouchTap={this.handleDelete}
                    icon={<DeleteIcon color={red500}/>}
                  />
                }
             ]}
             data={this.state[this.state.WidgetList.DataKey]}
             dataKey="Transactions"
          />
        </div>
    );
  }
}

export default AddRecordPage(TransactionList, "CONTROLLER_TRANSACTIONADD", "AddTransaction");
