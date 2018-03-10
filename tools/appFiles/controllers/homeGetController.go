package controllers

import (
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"strings"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
)

func (self *HomeController) Root(context session_functions.RequestContext, uriParams map[string]string, respond session_functions.ServerResponse) {
	var vm viewModel.HomeViewModel
	vm.LoadDefaultState()

	if session_functions.GetSessionAuthToken(context()) != constants.COOKIE_AUTHED {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	user, err := session_functions.GetSessionUser(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_USER_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	role, err := session_functions.GetSessionRole(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ROLE_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	account, err := session_functions.GetSessionAccount(context())
	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_ACCOUNT_NOT_FOUND, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	session_functions.Dump(user, account, role, account)

	releaseNotes, err := extensions.ReadFile(serverSettings.APP_LOCATION + "/releaseNotes.txt")
	if err == nil {

		title := "-goCoreProductName "

		notes := string(releaseNotes)
		idx := strings.Index(notes, title)
		startOfVersion := notes[idx:]
		startIdx := strings.Index(startOfVersion, "\n")
		versionNotes := startOfVersion[startIdx:]
		idxPreviousVersion := strings.Index(versionNotes, title)
		var currentVersionNotes string
		if idxPreviousVersion != -1 {
			currentVersionNotes = versionNotes[:idxPreviousVersion]
		} else {
			currentVersionNotes = versionNotes
		}
		currentVersionNotes = strings.Replace(strings.TrimSpace(currentVersionNotes), "\t", "", -1)
		versionLineItems := strings.Split(currentVersionNotes, "\n")
		if strings.Index(currentVersionNotes, "[d]") != -1 {
			for _, line := range versionLineItems {
				if strings.Index(line, "[d]") != -1 {
					vm.ReleaseDescriptionLines = append(vm.ReleaseDescriptionLines, line[4:])
				}
			}
		} else {
			for _, line := range versionLineItems {
				if strings.Index(line, "[d]") == -1 {
					typeFeature := ""

					if strings.Index(line, "[*]") != -1 {
						typeFeature = "Bug Fix:"
					} else if strings.Index(line, "[+]") != -1 {
						typeFeature = "New Feature:"
					} else if strings.Index(line, "[-]") != -1 {
						typeFeature = "Removed Function:"
					}
					vm.ReleaseDescriptionLines = append(vm.ReleaseDescriptionLines, typeFeature+" "+line[4:])
				}
			}
		}
	}

	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
