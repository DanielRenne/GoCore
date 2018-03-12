package controllers

import (
	"bytes"
	"encoding/csv"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/globalsign/mgo/bson"
)

func (self *FeaturesController) UpdateFeatureDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := UpdateFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		r = session_functions.ServerResponseToStruct(PARAM_FAILED, context, PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
	}

	hookSuccess := FeaturePostCommitHook("UpdateFeatureDetails", context, vm.Feature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *FeatureListController) ExportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureListViewModel
	vm.Parse(state)
	var i viewModel.FeatureImport
	schema := i.LoadSchema(context)
	b := &bytes.Buffer{}
	record := viewModel.GetCSVHeaderArray(context, schema)
	wr := csv.NewWriter(b)
	wr.Write(record)
	wr.Flush()

	for _, row := range vm.Features {
		record := []string{row.Id.Hex(), row.Key, row.Name, row.Description, row.FeatureGroupId}
		wr := csv.NewWriter(b)
		wr.Write(record)
		wr.Flush()
	}
	respond(PARAM_REDIRECT_NONE, "export.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", b.Bytes())
}

func (self *FeatureListController) GetImportCSVTemplate(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var i viewModel.FeatureImport
	schema := i.LoadSchema(context)
	output := viewModel.GetCSVTemplate(context, schema)
	respond(PARAM_REDIRECT_NONE, "import_Features.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", []byte(output))
}

func (self *FeatureListController) ImportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var r session_functions.ServerResponseStruct
	var vm viewModel.FeatureListViewModel
	vm.Parse(state)

	var i viewModel.FeatureImport
	rows, err := i.LoadSchemaAndParseFile(context, vm.FileUpload.Content)
	if err != nil {
		vm.FileUpload.Meta.CompleteFailure = true
		respond(PARAM_REDIRECT_NONE, "CSVError", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	actions := make(map[string]int, 0)
	errors, invalidRows, validRows := i.ValidateRows(context, rows)

	t, err := session_functions.StartTransaction(context())
	if err != nil {
		vm.FileUpload.Meta.CompleteFailure = true
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	var isUpdating bool
	for _, row := range validRows {
		var vm viewModel.FeatureModifyViewModel
		if row[i.Id.Idx] != "" {
			isUpdating = true
			vm.Feature.Id = bson.ObjectIdHex(row[i.Id.Idx])
		} else {
			isUpdating = false
		}
		//vm.Feature.XXXXXXX = row[i.XXXXXXXXX.Idx]
		vm.Feature.Key = row[i.Key.Idx]
		vm.Feature.Name = row[i.Name.Idx]
		vm.Feature.Description = row[i.Description.Idx]
		vm.Feature.FeatureGroupId = row[i.FeatureGroupId.Idx]

		if isUpdating {
			r = UpdateFeatureRow(context, &vm, t)
		} else {
			r = CreateFeatureRow(context, &vm, t)
		}

		if !r.CompletedSuccessfully {
			errors = append(errors, r.GlobalMessage)
			invalidRows = append(invalidRows, row)
		}

		if isUpdating {
			actions[vm.Feature.Id.Hex()] = 1
		} else {
			actions[vm.Feature.Id.Hex()] = 2
		}
	}
	err = t.Commit()
	if err != nil {
		vm.FileUpload.Meta.CompleteFailure = true
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	var updatedAffected int
	var createdAffected int
	for k, v := range actions {
		var hook string
		if v == 1 {
			updatedAffected += 1
			hook = "UpdateFeature"
		} else {
			createdAffected += 1
			hook = "CreateFeature"
		}
		hookSuccess := FeaturePostCommitHook(hook, context, k)
		if !hookSuccess {
			err = queries.Transactions.Rollback(context, "FeaturePostCommitHook failed", t.Id.Hex())
			if err != nil {
				respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
				return
			}
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
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
		if updatedAffected > 0 && createdAffected > 0 {
			key = "UpdatedXAndAddedXCSVRows"
		} else if updatedAffected > 0 && createdAffected == 0 {
			key = "UpdatedXCSVRows"
		} else {
			key = "CreatedXCSVRows"
		}
		if updatedAffected > 0 && createdAffected > 0 {
			replacements := queries.TagReplacements{
				Tag1: queries.Q("insert_count", extensions.IntToString(createdAffected)),
				Tag2: queries.Q("update_count", extensions.IntToString(updatedAffected)),
			}
			createdStr = queries.AppContent.GetTranslationWithReplacements(context, key, &replacements)
		} else {
			replacements := queries.TagReplacements{
				Tag1: queries.Q("row_count", extensions.IntToString(len(validRows))),
			}
			createdStr = queries.AppContent.GetTranslationWithReplacements(context, key, &replacements)
		}
	}

	vm.FileUpload.Meta.FileErrors = errors
	vm.FileUpload.Meta.RowsCommitted = len(validRows)
	vm.FileUpload.Meta.RowsCommittedInfo = createdStr
	vm.FileUpload.Meta.RowsSkipped = len(invalidRows)
	vm.FileUpload.Meta.RowsSkippedInfo = skippedStr

	respond(PARAM_REDIRECT_NONE, createdStr+skippedStr, snackbarSuccess, nil, PARAM_TRANSACTION_ID_NONE, vm)
}

func (self *FeaturesController) CreateFeature(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := CreateFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FeaturePostCommitHook("CreateFeature", context, vm.Feature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *FeaturesController) DeleteManyFeatures(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureListViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var FeatureId string
	for i := 0; i < len(vm.DeletedFeatures); i++ {
		var vmModify viewModel.FeatureModifyViewModel
		vmModify.Feature = vm.DeletedFeatures[i]
		FeatureId = vmModify.Feature.Id.Hex()
		r := DeleteFeatureRow(context, &vmModify, t)
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

	hookSuccess := FeaturePostCommitHook("DeleteManyFeatures", context, FeatureId)
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_RERENDER, constants.FEATURE_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *FeaturesController) DeleteFeature(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := DeleteFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FeaturePostCommitHook("DeleteFeature", context, vm.Feature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}

func (self *FeaturesController) CopyFeature(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FeatureModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	r := CopyFeatureRow(context, &vm, t)
	if !r.CompletedSuccessfully {
		respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FeaturePostCommitHook("CopyFeature", context, vm.Feature.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FeaturePostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(r.Redirect, r.GlobalMessage, r.GlobalMessageType, r.Trace, r.TransactionId, r.ViewModel)
}
