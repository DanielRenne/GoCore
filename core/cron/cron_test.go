package cron_test

import (
	"github.com/DanielRenne/GoCore/core/cron"
	"github.com/DanielRenne/GoCore/core/fileCache" // In order to use crons, you must import fileCache and call Initialize() if you are not using app.Initialize() or app.InitializeLite() to run a webserver (which do this automatically for you)
	"github.com/DanielRenne/GoCore/core/zip"
)

// ExampleStart is an example of how to start the cron job engine (note you must initialize fileCache first if you are not using app.Initialize() or app.InitializeLite() to run a webserver (which do this automatically for you)
func ExampleStart() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start()
}

// ExampleExecuteOneTimeJob is an example of how to execute a one time event.
func ExampleExecuteOneTimeJob() {
	/*
		import (
			"github.com/DanielRenne/GoCore/core/cron"
			"github.com/DanielRenne/GoCore/core/fileCache" // In order to use crons, you must import fileCache and call Initialize() if you are not using app.Initialize() or app.InitializeLite() to run a webserver (which do this automatically for you)
			"github.com/DanielRenne/GoCore/core/zip"
		)
	*/
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start()
	cb := func() bool {
		err := zip.Unzip("test", "test", []string{})
		if err != nil {
			return false
		}
		return true
	}
	cron.CronJobs.ExecuteOneTimeJob("run me and i will save a key of this string and write the boolean of success or not to /usr/local/goCore/jobs/jobs.json for historical purposes", cb)
}
