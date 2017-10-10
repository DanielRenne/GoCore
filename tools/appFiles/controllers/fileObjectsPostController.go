package controllers

import (
	"bytes"
	"encoding/csv"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"gopkg.in/mgo.v2/bson"
)

func (self *FileObjectsController) UpdateFileObjectDetails(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	message, err := br.FileObjects.Update(context, &vm, t)
	if err != nil {
		respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(constants.PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FileObjectPostCommitHook("UpdateFileObjectDetails", context, vm.FileObject.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FileObjectPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(constants.PARAM_REDIRECT_BACK, constants.FILEOBJECT_SAVE_SUCCESS, constants.PARAM_SNACKBAR_TYPE_ERROR, nil, constants.PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *FileObjectListController) ExportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectListViewModel
	vm.Parse(state)
	var i viewModel.FileObjectImport
	schema := i.LoadSchema(context)
	b := &bytes.Buffer{}
	record := viewModel.GetCSVHeaderArray(context, schema)
	wr := csv.NewWriter(b)
	wr.Write(record)
	wr.Flush()

	for _, row := range vm.FileObjects {
		record := []string{row.Id.Hex(), row.Name, row.Content, row.Type, row.MD5}
		wr := csv.NewWriter(b)
		wr.Write(record)
		wr.Flush()
	}
	respond(PARAM_REDIRECT_NONE, "export.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", b.Bytes())
}

func (self *FileObjectListController) GetImportCSVTemplate(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var i viewModel.FileObjectImport
	schema := i.LoadSchema(context)
	output := viewModel.GetCSVTemplate(context, schema)
	respond(PARAM_REDIRECT_NONE, "import_FileObjects.csv", PARAM_SNACKBAR_TYPE_DOWNLOAD_CONTENT, nil, "", []byte(output))
}

func (self *FileObjectListController) ImportCSV(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var message string
	var err error
	var vm viewModel.FileObjectListViewModel
	vm.Parse(state)

	var i viewModel.FileObjectImport
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
		var vm viewModel.FileObjectModifyViewModel
		if row[i.Id.Idx] != "" {
			isUpdating = true
			vm.FileObject.Id = bson.ObjectIdHex(row[i.Id.Idx])
		} else {
			isUpdating = false
		}
		//vm.FileObject.XXXXXXX = row[i.XXXXXXXXX.Idx]
		vm.FileObject.Name = row[i.Name.Idx]
		vm.FileObject.Content = row[i.Content.Idx]
		vm.FileObject.Type = row[i.Type.Idx]
		vm.FileObject.MD5 = row[i.MD5.Idx]

		if isUpdating {
			message, err = br.FileObjects.Update(context, &vm, t)
		} else {
			message, err = br.FileObjects.Create(context, &vm, t)
		}

		if err != nil {
			errors = append(errors, message)
			invalidRows = append(invalidRows, row)
		}
		if isUpdating {
			actions[vm.FileObject.Id.Hex()] = 1
		} else {
			actions[vm.FileObject.Id.Hex()] = 2
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
			hook = "UpdateFileObject"
		} else {
			createdAffected += 1
			hook = "CreateFileObject"
		}
		hookSuccess := FileObjectPostCommitHook(hook, context, k)
		if !hookSuccess {
			err = queries.Transactions.Rollback(context, "FileObjectPostCommitHook failed", t.Id.Hex())
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

func (self *FileObjectsController) CreateFileObject(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	message, err := br.FileObjects.Create(context, &vm, t)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, message, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FileObjectPostCommitHook("CreateFileObject", context, vm.FileObject.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FileObjectPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(constants.PARAM_REDIRECT_BACK, constants.FILEOBJECT_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, constants.PARAM_TRANSACTION_ID_NONE, vm)

}

func (self *FileObjectsController) DeleteManyFileObjects(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectListViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	var FileObjectId string
	for i := 0; i < len(vm.DeletedFileObjects); i++ {
		var vmModify viewModel.FileObjectModifyViewModel
		vmModify.FileObject = vm.DeletedFileObjects[i]
		FileObjectId = vmModify.FileObject.Id.Hex()
		message, err := br.FileObjects.Delete(context, &vmModify, t)
		if err != nil {
			respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
	}
	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FileObjectPostCommitHook("DeleteManyFileObjects", context, FileObjectId)
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FileObjectPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_RERENDER, constants.FILEOBJECT_DELETE_SUCCESSFUL, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *FileObjectsController) DeleteFileObject(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	message, err := br.FileObjects.Delete(context, &vm, t)
	if err != nil {
		respond(constants.PARAM_REDIRECT_NONE, message, constants.PARAM_SNACKBAR_TYPE_ERROR, err, constants.PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FileObjectPostCommitHook("DeleteFileObject", context, vm.FileObject.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FileObjectPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_RERENDER, message, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *FileObjectsController) CopyFileObject(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	if session_functions.BlockByDeveloperModeOff(context()) {
		return
	}
	var vm viewModel.FileObjectModifyViewModel
	vm.Parse(state)
	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	message, err := br.FileObjects.Copy(context, &vm, t)
	if err != nil {
		respond(PARAM_REDIRECT_NONE, message, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	hookSuccess := FileObjectPostCommitHook("CopyFileObject", context, vm.FileObject.Id.Hex())
	if !hookSuccess {
		err = queries.Transactions.Rollback(context, "FileObjectPostCommitHook failed", t.Id.Hex())
		if err != nil {
			respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_POST_HOOK, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
			return
		}
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_RERENDER, message, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)

}
