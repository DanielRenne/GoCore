package br

import "github.com/pkg/errors"

func ValidationError(msg string, errInfo error) (message string, err error) {
	err = errors.Wrap(errInfo, msg)
	message = msg
	return
}

var FileObjects fileObjectsBr
var Users usersBr
var Passwords passwordsBr
var AccountRoles accountRolesBr
var Server Server_Br
var Schedules schedulesBr
//NewBrVarsDontDeleteMe
