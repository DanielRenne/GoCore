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

//todo finish up required field *
//todo mass delete and delete inline.
//

class TransactionModify extends BasePageComponent {
  constructor(props, context) {
    super(props, context);


    // todo later if we can edit but i doubt you will be able to edit a tranny
    this.handleUserIdChange = (event) => {
    };


    this.handleLastUpdateChange = (event) => {
    };


    this.handleCompleteDateChange = (event) => {
    };


    this.handleRollbackDateChange = (event) => {
    };


    this.handleCommittedChange = (event) => {
    };


    this.handleErrorChange = (event) => {
    };


    this.handleCollectionsChange = (event) => {
    };


    this.handleDetailsChange = (event) => {
    };


    this.handleRolledBackChange = (event) => {
    };


    this.handleRolledBackByChange = (event) => {
    };


    this.handleRollbackReasonChange = (event) => {
    };



    this.save = () => {
        window.api.post({action: "UpdateTransactionDetails", state: this.state, controller:"transactions"});
    };
  }

  render() {
    this.logRender();
    return (
        <div>
            <CenteredPaperGrid>
                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyUserId}
                  hintText={window.pageContent.TransactionModifyUserId}
                  fullWidth={true}
                  onChange={this.handleUserIdChange}
                  defaultValue={this.state.Transaction.Joins.User.First + " " + this.state.Transaction.Joins.User.Last}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyLastUpdate}
                  hintText={window.pageContent.TransactionModifyLastUpdate}
                  fullWidth={true}
                  onChange={this.handleLastUpdateChange}
                  defaultValue={this.state.Transaction.lastUpdate}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyCompleteDate}
                  hintText={window.pageContent.TransactionModifyCompleteDate}
                  fullWidth={true}
                  onChange={this.handleCompleteDateChange}
                  defaultValue={this.state.Transaction.completeDate}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyRollbackDate}
                  hintText={window.pageContent.TransactionModifyRollbackDate}
                  fullWidth={true}
                  onChange={this.handleRollbackDateChange}
                  defaultValue={this.state.Transaction.rollbackDate}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyCommitted}
                  hintText={window.pageContent.TransactionModifyCommitted}
                  fullWidth={true}
                  onChange={this.handleCommittedChange}
                  defaultValue={this.state.Transaction.committed}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyError}
                  hintText={window.pageContent.TransactionModifyError}
                  fullWidth={true}
                  onChange={this.handleErrorChange}
                  defaultValue={this.state.Transaction.error}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyCollections}
                  hintText={window.pageContent.TransactionModifyCollections}
                  fullWidth={true}
                  onChange={this.handleCollectionsChange}
                  defaultValue={this.state.Transaction.collections}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyDetails}
                  hintText={window.pageContent.TransactionModifyDetails}
                  fullWidth={true}
                  onChange={this.handleDetailsChange}
                  defaultValue={this.state.Transaction.details}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyRolledBack}
                  hintText={window.pageContent.TransactionModifyRolledBack}
                  fullWidth={true}
                  onChange={this.handleRolledBackChange}
                  defaultValue={this.state.Transaction.rolledBack}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyRolledBackBy}
                  hintText={window.pageContent.TransactionModifyRolledBackBy}
                  fullWidth={true}
                  onChange={this.handleRolledBackByChange}
                  defaultValue={this.state.Transaction.rolledBackBy}
                />
                <br />                <TextField
                  floatingLabelText={window.pageContent.TransactionModifyRollbackReason}
                  hintText={window.pageContent.TransactionModifyRollbackReason}
                  fullWidth={true}
                  onChange={this.handleRollbackReasonChange}
                  defaultValue={this.state.Transaction.rollbackReason}
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

export default BackPage(TransactionModify);
