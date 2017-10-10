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

func (self *RoleFeaturesController) UpdateRoleFeatureDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.RoleFeatureModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := UpdateRoleFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		r = session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	hookSuccess := RoleFeaturePostCommitHook("UpdateRoleFeatureDetails", context, vm.RoleFeature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "RoleFeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *RoleFeatureListController) ExportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.RoleFeatureListViewModel
	vm.Parse(state)
	var i viewModel.RoleFeatureImport
	schema := i.LoadSchema(context)
	b := &bytes.Buffer{}
	record := viewModel.GetCSVHeaderArray(context, schema)
	wr := csv.NewWriter(b)
	wr.Write(record)
	wr.Flush()

	for _, row := range vm.RoleFeatures {
		record := []string{row.Id.Hex()}
		wr := csv.NewWriter(b)
		wr.Write(record)
		wr.Flush()
	}
	respond(PARAM_REDIRECT_NONE, "export.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", b.Bytes())
}

func (self *RoleFeatureListController) GetImportCSVTemplate(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var i viewModel.RoleFeatureImport
	schema := i.LoadSchema(context)
	output := viewModel.GetCSVTemplate(context, schema)
	respond(PARAM_REDIRECT_NONE, "import_RoleFeatures.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", []byte(output))
}

func (self *RoleFeatureListController) ImportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var r session_functions.ServerResponseStruct
	var vm viewModel.RoleFeatureListViewModel
	vm.Parse(state)

	var i viewModel.RoleFeatureImport
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
		var vm viewModel.RoleFeatureModifyViewModel
		if row[i.Id.Idx] != "" {
			isUpdating = true
			vm.RoleFeature.Id = bson.ObjectIdHex(row[i.Id.Idx])
		}
		//vm.RoleFeature.XXXXXXX = row[i.XXXXXXXXX.Idx]

		if isUpdating {
			r = UpdateRoleFeatureRow(context, &vm, t)
		} else {
			r = CreateRoleFeatureRow(context, &vm, t)
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

func (self *RoleFeaturesController) CreateRoleFeature(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.RoleFeatureModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := CreateRoleFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := RoleFeaturePostCommitHook("CreateRoleFeature", context, vm.RoleFeature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "RoleFeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *RoleFeaturesController) DeleteManyRoleFeatures(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.RoleFeatureListViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var RoleFeatureId string
	for i := 0; i < len(vm.DeletedRoleFeatures); i++ {
		var vmModify viewModel.RoleFeatureModifyViewModel
		vmModify.RoleFeature = vm.DeletedRoleFeatures[i]
		RoleFeatureId = vmModify.RoleFeature.Id.Hex()
		r := DeleteRoleFeatureRow(context, &vmModify, t)
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

	hookSuccess := RoleFeaturePostCommitHook("DeleteManyRoleFeatures", context, RoleFeatureId)
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "RoleFeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_RERENDER, constants.ROLEFEATURE_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *RoleFeaturesController) DeleteRoleFeature(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.RoleFeatureModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := DeleteRoleFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := RoleFeaturePostCommitHook("DeleteRoleFeature", context, vm.RoleFeature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "RoleFeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *RoleFeaturesController) CopyRoleFeature(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.RoleFeatureModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := CopyRoleFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := RoleFeaturePostCommitHook("CopyRoleFeature", context, vm.RoleFeature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "RoleFeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}
