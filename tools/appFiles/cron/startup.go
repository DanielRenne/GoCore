package cron

import (
	"time"
	"fmt"
	"runtime/debug"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

func startup() {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("startup.go", "Panic Recovered at startup:  "+fmt.Sprintf("%+v", r))
			time.Sleep(time.Millisecond * 3000)
			startup()
			return
		}
	}()

}