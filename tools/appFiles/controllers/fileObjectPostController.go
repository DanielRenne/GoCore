package controllers

import (
	"encoding/base64"
	"image"
	"path/filepath"
	"strings"

	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/disintegration/imaging"
	"github.com/go-errors/errors"
)

var FileUpload FileUploadController

type FileUploadController struct {
}

func (self *FileUploadController) Save(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.FileObjectViewModel
	vm.Parse(state)
	extension := strings.ToLower(filepath.Ext(vm.FileObject.Name))
	extensionsIgnored := utils.Array(".svg")
	extensionsSupported := utils.Array(".jpg", ".jpeg", ".png", ".gif", ".tif", ".tiff", ".bmp")
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(vm.FileObject.Content))
	m, _, err := image.Decode(reader)
	if !utils.InArray(extension, extensionsIgnored) {
		if vm.Width > 0 && utils.InArray(extension, extensionsSupported) {
			if err != nil {
				respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
			}
			var dstImageFit *image.NRGBA
			if vm.Width > vm.Height {
				dstImageFit = imaging.Resize(m, vm.Width, 0, imaging.Lanczos)
			} else {
				dstImageFit = imaging.Resize(m, 0, vm.Height, imaging.Lanczos)
			}

			format := imaging.JPEG
			switch extension {
			case ".jpg":
			case ".jpeg":
				format = imaging.PNG
				break
			case ".png":
				format = imaging.PNG
				break
			case ".gif":
				format = imaging.GIF
				break
			case ".tif":
				format = imaging.TIFF
				break
			case ".tiff":
				format = imaging.TIFF
				break
			case ".bmp":
				format = imaging.BMP
				break
			}

			imaging.Encode(&vm, dstImageFit, format)
			vm.SaveResize()
		} else if vm.Width > 0 && err == nil {
			// This is a valid image file, but unsupported on the resize
			err = errors.New("Invalid extension for image resize.  Only " + strings.Join(extensionsSupported, ", "))
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
			return
		} else if vm.Width > 0 && err != nil {
			respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
			return
		}
	}

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	acct, err := session_functions.GetSessionAccount(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	vm.FileObject.AccountId = acct.Id.Hex()

	err = vm.FileObject.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	vm.FileObject.Content = ""

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
