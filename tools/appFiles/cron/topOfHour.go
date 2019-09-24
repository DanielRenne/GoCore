package cron

import (
	"time"

	"github.com/DanielRenne/GoCore/core"
)

func ClearDebugMemory(x time.Time) {
	// Just in case some job is running and session_functions.Dump but no requests are clearing it and memory is growing
	core.TransactionLogMutex.Lock()
	core.TransactionLog = ""
	core.TransactionLogMutex.Unlock()
}

func DeleteImageHistory(x time.Time) {
}
