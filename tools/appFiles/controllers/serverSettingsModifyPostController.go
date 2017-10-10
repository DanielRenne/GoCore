package controllers

import (
	"fmt"
	"time"

	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel/posts"
)

func (self *ServerSettingsModifyController) UpdateServerSettings(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.ServerSettingsModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}
	respond(PARAM_REDIRECT_NONE, constants.SERVER_SETTING_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)

}


func (self *ServerSettingsModifyController) UpdateGatewayTimeSettings(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.ServerSettingsModifyViewModel
	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	changedTimeZone := false
	tzSetting, err := queries.ServerSettings.ById(constants.SERVER_SETTING_TIMEZONE)
	if err == nil && tzSetting.Value != vm.TimeZone.Value {
		vm.TimeZone.SaveWithTran(t)
		changedTimeZone = true
	}

	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to Change Time Zone", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	if changedTimeZone {
		tme := time.Now()
		br.Schedules.LoadDay(tme)
	}

	format := "2006-01-02"
	formatTime := "03:04:05 PM MST"
	//
	// session_functions.Log("Time Details", fmt.Sprintf("%+v", vm.DateToSet)+","+fmt.Sprintf("%+v", vm.TimeToSet))
	if vm.DateToSet.Unix() > 0 && vm.TimeToSet.Unix() > 0 { //Set both the Date and Time
		t := time.Now()
		vm.TimeToSet.Add(time.Second * time.Duration(t.Second()))
		br.SetDateAndTime(fmt.Sprintf("%+v", vm.DateToSet.Format(format)), fmt.Sprintf("%+v", vm.TimeToSet.Format(formatTime)))
		tme := time.Now()
		br.Schedules.LoadDay(tme)
	}
	if vm.DateToSet.Unix() > 0 && vm.TimeToSet.Unix() < 0 { //Set only the Date

		br.SetDate(fmt.Sprintf("%+v", vm.DateToSet.Format(format)))
	}
	if vm.TimeToSet.Unix() > 0 { //Set only the Time
		t := time.Now()
		vm.TimeToSet.Add(time.Second * time.Duration(t.Second()))
		br.SetTime(fmt.Sprintf("%+v", vm.TimeToSet.Format(formatTime)))
		tme := time.Now()
		br.Schedules.LoadDay(tme)
	}

	respond(PARAM_REDIRECT_RERENDER, constants.SERVER_SETTING_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}

func (self *ServerSettingsModifyController) EnableNTPServer(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.ServerSettingsModifyViewModel
	vm.Parse(state)

	err := br.EnableNTPServer()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, "Failed to Enable NTP Service", PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	respond(PARAM_REDIRECT_RERENDER, constants.SERVER_SETTING_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)
}


func (self *ServerSettingsModifyController) UpdateLockoutSettings(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {

	var vm viewModel.ServerSettingsModifyViewModel

	vm.Parse(state)

	t, err := session_functions.StartTransaction(context())

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_RETRIEVE, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, viewModel.EmptyViewModel{})
		return
	}

	valid := constants.BITWISE_TRUE

	if vm.LockoutSettings.Lockout.Value == "" {
		valid &= constants.BITWISE_FALSE
		vm.LockoutSettings.Lockout.Errors.Value = constants.ERROR_REQUIRED_FIELD
	}

	if valid != constants.BITWISE_TRUE {
		respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_ERROR, nil, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	vm.LockoutSettings.Lockout.SaveWithTran(t)

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.SERVER_SETTING_SAVE_FAIL, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	err = t.Commit()

	if err != nil {
		respond(PARAM_REDIRECT_NONE, constants.ERRORS_TRANSACTION_FAILED_TO_COMMIT, PARAM_SNACKBAR_TYPE_ERROR, err, PARAM_TRANSACTION_ID_NONE, vm)
		return
	}

	respond(PARAM_REDIRECT_NONE, constants.SERVER_SETTING_SAVE_SUCCESS, PARAM_SNACKBAR_TYPE_SUCCESS, nil, t.Id.Hex(), vm)
}


func (self *ServerSettingsModifyController) RestartMachine(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.ServerSettingsModifyViewModel
	vm.Parse(state)

	var vm2 viewModel.StatusUpdateViewModel
	vm2.Mode = "dbdump"
	vm2.Message2 = "RestartDone"
	vm2.Message = queries.AppContent.GetTranslation(session_functions.PassContext(context()), "RestartDone")
	app.PublishWebSocketJSON("FirmwareStatus", vm2)

	respond(PARAM_REDIRECT_NONE, "Restarting Machine ...", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)

	go logger.GoRoutineLogger(func() {
		time.Sleep(time.Second * 7)
		err := br.Reboot()
		if err != nil {
			session_functions.Log("Error", "Failed to Reboot Machine:  "+err.Error())
		} else {
			session_functions.Log("Reboot", "Reboot Successful")
		}

	}, "serverSettingsModifyPostController.go->Restart")
}

func (self *ServerSettingsModifyController) ShutdownMachine(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm viewModel.ServerSettingsModifyViewModel
	vm.Parse(state)

	respond(PARAM_REDIRECT_NONE, "Shutting Down Machine...", PARAM_SNACKBAR_TYPE_SUCCESS, nil, "", vm)

	go logger.GoRoutineLogger(func() {
		time.Sleep(time.Millisecond * 1000)
		err := br.Shutdown()
		if err != nil {
			session_functions.Log("Error", "Failed to Shutdown Machine:  "+err.Error())
		} else {
			session_functions.Log("Shutdown", "Shutting Down Succesful")
		}

	}, "serverSettingsModifyPostController.go->ShutdownMachine")
}


func (self *ServerSettingsModifyController) GetTimeZones(context session_functions.RequestContext, state string, respond session_functions.ServerResponse) {
	var vm posts.TimeZoneVM
	vm.TimeZones = model.TimeZoneLocations
	respond(PARAM_REDIRECT_NONE, PARAM_SNACKBAR_MESSAGE_NONE, PARAM_SNACKBAR_TYPE_SUCCESS, nil, PARAM_TRANSACTION_ID_NONE, vm)
}
