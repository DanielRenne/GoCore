package br

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/pkg/errors"
)

type passwordsBr struct{}

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func (self passwordsBr) Create(context session_functions.RequestContext, vm *viewModel.PasswordModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}

	err = vm.Password.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->passwordsBr->Create", fmt.Sprintf("%+v", vm.Password.Errors))
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_PASSWORD_SAVE) + core.Debug.HandleError(err)
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->passwordsBr->Create", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_PASSWORD_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.PASSWORD_SAVE_SUCCESS
	return
}

func (self passwordsBr) Update(context session_functions.RequestContext, vm *viewModel.PasswordModifyViewModel, t *model.Transaction) (message string, err error) {

	message, err = self.Validate(vm)
	if err != nil {
		return
	}

	err = vm.Password.SaveWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->passwordsBr->Update", fmt.Sprintf("%+v", vm.Password.Errors))
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_PASSWORD_SAVE) + core.Debug.HandleError(err)
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->passwordsBr->Update", err.Error())
			message = queries.AppContent.GetTranslation(context, constants.ERRORS_PASSWORD_SAVE) + core.Debug.HandleError(err)
			return
		}
	}

	message = constants.PASSWORD_SAVE_SUCCESS
	return
}

func (self passwordsBr) Delete(context session_functions.RequestContext, vm *viewModel.PasswordModifyViewModel, t *model.Transaction) (message string, err error) {

	err = vm.Password.DeleteWithTran(t)

	if err != nil {
		if model.IsValidationError(err) {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->Validation_Error->passwordsBr->Delete", fmt.Sprintf("%+v", vm.Password.Errors))
			message = constants.ERRORS_PASSWORD_DELETE
			return
		} else {
			err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
			session_functions.Log("Error->passwordsBr->Delete", err.Error())
			message = constants.ERRORS_PASSWORD_DELETE
			return
		}
	}

	message = constants.PASSWORD_DELETE_SUCCESSFUL
	return
}

func (self passwordsBr) WrapViewModel(PasswordInstance model.Password) viewModel.PasswordModifyViewModel {
	vm := viewModel.PasswordModifyViewModel{}
	vm.Password = PasswordInstance
	return vm
}

func (self passwordsBr) Validate(vm *viewModel.PasswordModifyViewModel) (message string, err error) {
	message = ""
	return
}

func (self passwordsBr) NewPassword(length int) string {
	return rand_char(length, StdChars)
}

func rand_char(length int, chars []byte) string {
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			panic(err)
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword)
			}
		}
	}
}
