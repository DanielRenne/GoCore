package cron

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/scheduleEngine"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

func Start() {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("main.go", "Panic Recovered at Start:  "+fmt.Sprintf("%+v", r))
			time.Sleep(time.Millisecond * 3000)
			Start()
			return
		}
	}()
	go core.CronJobs.RegisterRecurring(core.CRON_TOP_OF_30_SECONDS, ClearLogs)
	go core.CronJobs.RegisterRecurring(core.CRON_TOP_OF_HOUR, ClearDebugMemory)
	go core.CronJobs.RegisterRecurring(core.CRON_TOP_OF_HOUR, DeleteImageHistory)
	go core.CronJobs.RegisterRecurring(core.CRON_TOP_OF_SECOND, FlushLogs)
	go core.CronJobs.RegisterRecurring(core.CRON_TOP_OF_SECOND, scheduleEngine.Trigger)
	go core.CronJobs.RegisterRecurring(core.CRON_TOP_OF_SECOND, BroadcastTime)
	go startup()
}
