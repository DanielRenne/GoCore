package cron

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/DanielRenne/goCoreAppTemplate/br"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

func BroadcastTime(x time.Time) {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("Panic->cron->topOfSecond->BroadcastTime", "Panic Recovered:  "+fmt.Sprintf("%+v", r))
			return
		}
	}()
	loc, err := time.LoadLocation(br.Schedules.GetTimeZone())
	if err != nil {
		session_functions.Log("Failed to LoadLocation", br.Schedules.GetTimeZone()+":"+err.Error())
		return
	}
	t := x.In(loc)
	session_functions.BroadcastTime(fmt.Sprintf("%+v", t.Format("1-2-2006")), fmt.Sprintf("%+v", t.Format("03:04:05 PM")))
}

func FlushLogs(x time.Time) {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("Panic->cron->topOfSecond->FlushLogs", "Panic Recovered:  "+fmt.Sprintf("%+v", r))
			return
		}
	}()

	if x.Second() == 0 || x.Second() == 10 || x.Second() == 20 || x.Second() == 30 || x.Second() == 40 || x.Second() == 50 {
		//After all logs are cleared append the new ones that are in the buffer
		session_functions.FlushAllLogs()
	}
}
