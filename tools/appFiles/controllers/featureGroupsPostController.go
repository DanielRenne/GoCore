package controllers

import (
	"bytes"
	"encoding/csv"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"gopkg.in/mgo.v2/bson"
)

func (self *FeatureGroupsController) UpdateFeatureGroupDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureGroupModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := UpdateFeatureGroupRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		r = session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	hookSuccess := FeatureGroupPostCommitHook("UpdateFeatureGroupDetails", context, vm.FeatureGroup.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeatureGroupPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *FeatureGroupListController) ExportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureGroupListViewModel
	vm.Parse(state)
	var i viewModel.FeatureGroupImport
	schema := i.LoadSchema(context)
	b := &bytes.Buffer{}
	record := viewModel.GetCSVHeaderArray(context, schema)
	wr := csv.NewWriter(b)
	wr.Write(record)
	wr.Flush()

	for _, row := range vm.FeatureGroups {
		record := []string{row.Id.Hex(), row.Name}
		wr := csv.NewWriter(b)
		wr.Write(record)
		wr.Flush()
	}
	respond(PARAM_REDIRECT_NONE, "export.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", b.Bytes())
}

func (self *FeatureGroupListController) GetImportCSVTemplate(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var i viewModel.FeatureGroupImport
	schema := i.LoadSchema(context)
	output := viewModel.GetCSVTemplate(context, schema)
	respond(PARAM_REDIRECT_NONE, "import_FeatureGroups.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", []byte(output))
}

func (self *FeatureGroupListController) ImportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var r session_functions.ServerResponseStruct
	var vm viewModel.FeatureGroupListViewModel
	vm.Parse(state)

	var i viewModel.FeatureGroupImport
	rows, err := i.LoadSchemaAndParseFile(context, vm.FileUpload.Content)
	if err != nil {
		vm.FileUpload.Meta.CompleteFailure = true
		respond(PARAM_REDIRECT_NONE, "CSVError", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	errors, invalidRows, validRows := i.ValidateRows(context, rows)

	t, err := session_functions.StartTransaction(context())
	if err != nil {
		vm.FileUpload.Meta.CompleteFailure = true
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	var isUpdating bool
	for _, row := range validRows {
		var vm viewModel.FeatureGroupModifyViewModel
		if row[i.Id.Idx] != "" {
			isUpdating = true
			vm.FeatureGroup.Id = bson.ObjectIdHex(row[i.Id.Idx])
		}
		//vm.FeatureGroup.XXXXXXX = row[i.XXXXXXXXX.Idx]
		vm.FeatureGroup.Name = row[i.Name.Idx]

		if isUpdating {
			r = UpdateFeatureGroupRow(context, &vm, t)
		} else {
			r = CreateFeatureGroupRow(context, &vm, t)
		}

		if !r.CompletedSuccessfully {
			errors = append(errors, r.GlobalMessage)
			invalidRows = append(invalidRows, row)
		}
	}

	err = t.Commit()
	if err != nil {
		vm.FileUpload.Meta.CompleteFailure = true
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	// Clear out file contents to limit large payload size
	snackbarSuccess := SNACKBAR_TYPE_SUCCESS
	skippedStr := ""
	if len(invalidRows) > 0 {
		snackbarSuccess = SNACKBAR_TYPE_WARNING
		replacements := queries.TagReplacements{
			Tag1: queries.Q("row_count", extensions.IntToString(len(invalidRows))),
		}
		skippedStr = " " + queries.AppContent.GetTranslationWithReplacements(context, "SkippedXCSVRows", &replacements)
	}
	createdStr := ""
	if len(validRows) > 0 {
		var key string
		if isUpdating {
			key = "UpdatedXCSVRows"
		} else {
			key = "CreatedXCSVRows"
		}
		replacements := queries.TagReplacements{
			Tag1: queries.Q("row_count", extensions.IntToString(len(validRows))),
		}
		createdStr = queries.AppContent.GetTranslationWithReplacements(context, key, &replacements)
	}

	vm.FileUpload.Meta.FileErrors = errors
	vm.FileUpload.Meta.RowsCommitted = len(validRows)
	vm.FileUpload.Meta.RowsCommittedInfo = createdStr
	vm.FileUpload.Meta.RowsSkipped = len(invalidRows)
	vm.FileUpload.Meta.RowsSkippedInfo = skippedStr

	respond(PARAM_REDIRECT_NONE, createdStr+skippedStr, snackbarSuccess, nil, PARAM_TRANSACTION_ID_NONE, vm)
}

func (self *FeatureGroupsController) CreateFeatureGroup(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureGroupModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := CreateFeatureGroupRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FeatureGroupPostCommitHook("CreateFeatureGroup", context, vm.FeatureGroup.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeatureGroupPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *FeatureGroupsController) DeleteManyFeatureGroups(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureGroupListViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var FeatureGroupId string
	for i := 0; i < len(vm.DeletedFeatureGroups); i++ {
		var vmModify viewModel.FeatureGroupModifyViewModel
		vmModify.FeatureGroup = vm.DeletedFeatureGroups[i]
		FeatureGroupId = vmModify.FeatureGroup.Id.Hex()
		r := DeleteFeatureGroupRow(context, &vmModify, t)
		if !r.CompletedSuccessfully {
			respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
			return
		}
	}
	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FeatureGroupPostCommitHook("DeleteManyFeatureGroups", context, FeatureGroupId)
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeatureGroupPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_RERENDER, constants.FEATUREGROUP_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *FeatureGroupsController) DeleteFeatureGroup(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureGroupModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := DeleteFeatureGroupRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	hookSuccess := FeatureGroupPostCommitHook("DeleteFeatureGroup", context, vm.FeatureGroup.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeatureGroupPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *FeatureGroupsController) CopyFeatureGroup(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureGroupModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := CopyFeatureGroupRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FeatureGroupPostCommitHook("CopyFeatureGroup", context, vm.FeatureGroup.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeatureGroupPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}
