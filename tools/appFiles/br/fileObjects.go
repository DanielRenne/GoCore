package br

import (
	"fmt"

	"bytes"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"path/filepath"
	"strings"
)

type fileObjectsBr struct{}

func (self fileObjectsBr) Resize(content []byte, fileName string, width uint, height uint) (resizeContent []byte, err error) {
	vm := viewModel.FileObjectResizeViewModel{
		ImageResize: content,
	}
	extension := strings.ToLower(filepath.Ext(fileName))
	extensionsIgnored := utils.Array(".svg")
	extensionsSupported := utils.Array(".jpg", ".jpeg", ".png", ".gif", ".tif", ".tiff", ".bmp")
	reader := strings.NewReader(string(vm.ImageResize))
	m, _, err := image.Decode(reader)
	if !utils.InArray(extension, extensionsIgnored) {
		if width > 0 && height > 0 && utils.InArray(extension, extensionsSupported) {
			if err != nil {
				return
			}
			resized := resize.Resize(width, height, m, resize.NearestNeighbor)
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, resized, nil)
			resizeContent = buf.Bytes()
		}
	}

	return
}

func (self fileObjectsBr) Create(context session_functions.RequestContext, vm *viewModel.FileObjectModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}
	acct, actErr := queries.Common.GetAccount(context)
	if actErr == nil {
		vm.FileObject.AccountId = acct.Id.Hex()
	}
	err = vm.FileObject.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->fileObjectsBr->Create", fmt.Sprintf("%+v", vm.FileObject.Errors))
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_FILEOBJECT_SAVE) + core.Debug.HandleError(err)
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->fileObjectsBr->Create", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_FILEOBJECT_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.FILEOBJECT_SAVE_SUCCESS

	return

}

func (self fileObjectsBr) Update(context session_functions.RequestContext, vm *viewModel.FileObjectModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}

	err = vm.FileObject.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->fileObjectsBr->Update", fmt.Sprintf("%+v", vm.FileObject.Errors))
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_FILEOBJECT_SAVE) + core.Debug.HandleError(err)
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->fileObjectsBr->Update", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_FILEOBJECT_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.FILEOBJECT_SAVE_SUCCESS

	return

}

func (self fileObjectsBr) Delete(context session_functions.RequestContext, vm *viewModel.FileObjectModifyViewModel, t *model.Transaction) (message string, err error) {
	err = vm.FileObject.DeleteWithTran(t)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		message = constants.ERRORS_FILEOBJECT_DELETE
		return
	}

	message = constants.FILEOBJECT_DELETE_SUCCESSFUL
	return
}

func (self fileObjectsBr) Copy(context session_functions.RequestContext, copyRowVm *viewModel.FileObjectModifyViewModel, t *model.Transaction) (message string, err error) {

	err = model.FileObjects.Query().ById(copyRowVm.FileObject.Id, &copyRowVm.FileObject)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		message = constants.ERRORS_FILEOBJECT_COPY
		return
	}

	copyRowVm.FileObject.Id = ""

	message, err = self.Create(context, copyRowVm, t)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		message = constants.ERRORS_FILEOBJECT_COPY
		return
	}

	message = constants.FILEOBJECT_COPY_SUCCESSFUL
	return
}

func (self fileObjectsBr) WrapViewModel(FileObjectInstance model.FileObject) viewModel.FileObjectModifyViewModel {
	vm := viewModel.FileObjectModifyViewModel{}
	vm.FileObject = FileObjectInstance
	return vm
}

func (self fileObjectsBr) Validate(vm *viewModel.FileObjectModifyViewModel) (message string, err error) {
	return
}

func (self fileObjectsBr) GetVmDefaults(vm *viewModel.FileObjectModifyViewModel) (message string, err error) {
	return
}
