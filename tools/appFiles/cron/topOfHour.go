package cron

import (
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
)

func ClearDebugMemory(x time.Time) {
	// Just in case some job is running and session_functions.Dump but no requests are clearing it and memory is growing
	core.TransactionLogMutex.Lock()
	core.TransactionLog = ""
	core.TransactionLogMutex.Unlock()
}

func DeleteImageHistory(x time.Time) {
	var fo []model.FileObjectHistoryRecord
	err := model.FileObjectsHistory.Query().All(&fo)
	if err != nil {
		for _, f := range fo {
			f.Delete()
		}
	}
}
