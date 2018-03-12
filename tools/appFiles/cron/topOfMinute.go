package cron

import (
	"errors"
	"os"
	"runtime/debug"
	"time"

	"fmt"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

func ClearLogs(x time.Time) {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("topOfMinute.go", "Panic Recovered at ClearLogs:  "+fmt.Sprintf("%+v", r))
			return
		}
	}()
	baseLogPath := os.Getenv("GOPATH") + "/" + serverSettings.APP_LOCATION + "/log/"

	size, err := extensions.GetFileSize(baseLogPath + "app.log")
	if size > 100*1024*1024 {
		err = os.Remove(baseLogPath + "app.log")
		if err != nil {
			err = errors.New("Failed to clearLogs:\n" + err.Error())
			return
		}
	}
	d, err := os.Open(baseLogPath + "plugins/")
	if err != nil {
		session_functions.Log("ClearLogs", err.Error())
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		session_functions.Log("ClearLogs", err.Error())
	}
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			if fi.Name() != ".gitkeep" {
				size, err := extensions.GetFileSize(os.Getenv("GOPATH") + "/" + serverSettings.APP_LOCATION + "/log/plugins/" + fi.Name())
				if size > 3*1024*1024 {

					err = os.Remove(os.Getenv("GOPATH") + "/" + serverSettings.APP_LOCATION + "/log/plugins/" + fi.Name())
					if err != nil {
						err = errors.New("Failed to clearLogs:\n" + err.Error())
					}
				}
			}
		}
	}
}
