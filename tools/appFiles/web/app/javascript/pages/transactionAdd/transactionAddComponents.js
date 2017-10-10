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
import TransactionModify from '../transactionModify/transactionModifyComponents';

class TransactionAdd extends TransactionModify {
  constructor(props, context) {
    super(props, context);

    this.createTransaction = () => {
        window.api.post({action: "CreateTransaction", state: this.state, controller:"transactions"});
    };
  }

  render() {
    this.logRender();
    return (
      <CenteredPaperGrid>
			<TextField
				floatingLabelText={window.pageContent.TransactionAddUserId}
				hintText={window.pageContent.TransactionAddUserId}
				fullWidth={true}
				onChange={this.handleUserIdChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.UserId)}
				defaultValue={this.state.Transaction.userId}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddLastUpdate}
				hintText={window.pageContent.TransactionAddLastUpdate}
				fullWidth={true}
				onChange={this.handleLastUpdateChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.LastUpdate)}
				defaultValue={this.state.Transaction.lastUpdate}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddCompleteDate}
				hintText={window.pageContent.TransactionAddCompleteDate}
				fullWidth={true}
				onChange={this.handleCompleteDateChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.CompleteDate)}
				defaultValue={this.state.Transaction.completeDate}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddRollbackDate}
				hintText={window.pageContent.TransactionAddRollbackDate}
				fullWidth={true}
				onChange={this.handleRollbackDateChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.RollbackDate)}
				defaultValue={this.state.Transaction.rollbackDate}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddCommitted}
				hintText={window.pageContent.TransactionAddCommitted}
				fullWidth={true}
				onChange={this.handleCommittedChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.Committed)}
				defaultValue={this.state.Transaction.committed}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddError}
				hintText={window.pageContent.TransactionAddError}
				fullWidth={true}
				onChange={this.handleErrorChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.Error)}
				defaultValue={this.state.Transaction.error}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddCollections}
				hintText={window.pageContent.TransactionAddCollections}
				fullWidth={true}
				onChange={this.handleCollectionsChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.Collections)}
				defaultValue={this.state.Transaction.collections}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddDetails}
				hintText={window.pageContent.TransactionAddDetails}
				fullWidth={true}
				onChange={this.handleDetailsChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.Details)}
				defaultValue={this.state.Transaction.details}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddRolledBack}
				hintText={window.pageContent.TransactionAddRolledBack}
				fullWidth={true}
				onChange={this.handleRolledBackChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.RolledBack)}
				defaultValue={this.state.Transaction.rolledBack}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddRolledBackBy}
				hintText={window.pageContent.TransactionAddRolledBackBy}
				fullWidth={true}
				onChange={this.handleRolledBackByChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.RolledBackBy)}
				defaultValue={this.state.Transaction.rolledBackBy}
			/>
			<br />			<TextField
				floatingLabelText={window.pageContent.TransactionAddRollbackReason}
				hintText={window.pageContent.TransactionAddRollbackReason}
				fullWidth={true}
				onChange={this.handleRollbackReasonChange}
				errorText={this.globs.translate(this.state.Transaction.Errors.RollbackReason)}
				defaultValue={this.state.Transaction.rollbackReason}
			/>
			<br />
        <br />
        <br />
        <RaisedButton
            label={window.pageContent.CreateTransaction}
            onTouchTap={this.createTransaction}
            secondary={true}
            id="save"
        />
        </CenteredPaperGrid>
    );
  }
}

export default BackPage(TransactionAdd);
