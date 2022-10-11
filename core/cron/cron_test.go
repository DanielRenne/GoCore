package cron_test

import (
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/cron"
	"github.com/DanielRenne/GoCore/core/fileCache" // In order to use crons, you must import fileCache and call Initialize() if you are not using app.Initialize() or app.InitializeLite() to run a webserver (which do this automatically for you)
	"github.com/DanielRenne/GoCore/core/zip"
)

// This is how to start the cron job engine (note you must initialize fileCache first if you are not using app.Initialize() or app.InitializeLite() to run a webserver (which do this automatically for you).  Only run Start once
func ExampleStart() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start()
}

// You may have a situation that upon initializing, you have a user defined configuration to control the interval of the job.  This is how you would do that.
func ExampleRegisterRecurring() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once

	// In this situation, you have a user defined configuration to control the interval of the job.  This is how you would do that:
	var runCronAt cron.RecurringType
	var someConfigValue string
	if someConfigValue == "top-of-hour" {
		runCronAt = cron.CronTopOfHour
	}
	if someConfigValue == "" {
		runCronAt = cron.CronTopOfSecond
	}
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on top of second (because it wasnt top-of-hour): " + currentTime.String())
		return
	}
	cron.RegisterRecurring(runCronAt, cb)
	c := make(chan string)
	go func() {
		time.Sleep(10 * time.Second)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)

}

// ExampleExecuteOneTimeJob is an example of how to execute a one time event.
func ExampleExecuteOneTimeJob() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func() bool {
		err := zip.Unzip("test", "test", []string{})
		if err != nil {
			return false
		}
		return true
	}
	cron.ExecuteOneTimeJob("run me and i will save a key of this string and write the boolean of success or not to /usr/local/goCore/jobs/jobs.json for historical purposes", cb)
}

// execute job every 5 minutes
func ExampleShouldRunEveryFiveMinutes() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func(currentTime time.Time) {
		// Here you wrap the current time with the helper function so that minutes 1, 2, 3, 4 are skipped, but minutes 0 and 5 are executed.
		if cron.ShouldRunEveryFiveMinutes(currentTime) {
			core.Dump("Job Executed on top of 5 minutes: " + currentTime.String())
		}
		return
	}
	cron.RegisterTopOfMinuteJob(cb)
	c := make(chan string)
	go func() {
		time.Sleep(11 * time.Minute)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)
}

// top of 30 seconds event
func ExampleRegisterTopOf30SecondsJob() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
		return
	}
	cron.RegisterTopOf30SecondsJob(cb)
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Minute)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)
}

// executes top of minute job.
func ExampleRegisterTopOfMinuteJob() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
		return
	}
	cron.RegisterTopOfMinuteJob(cb)
	c := make(chan string)
	go func() {
		time.Sleep(2 * time.Minute)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)
}

// executes top of day.
func ExampleRegisterTopOfDayJob() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
		return
	}
	cron.RegisterTopOfDayJob(cb)
	c := make(chan string)
	go func() {
		time.Sleep(24 * time.Hour)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)
}

// executes top of second job.
func ExampleRegisterTopOfSecondJob() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
		return
	}
	cron.RegisterTopOfSecondJob(cb)
	c := make(chan string)
	go func() {
		time.Sleep(10 * time.Second)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)
}

// executes top of the hour.
func ExampleRegisterTopOfHourJob() {
	// You can set the CACHE_STORAGE_PATH export to a folder on your system to use the goCore fileCache.
	// fileCache.CACHE_STORAGE_PATH = "/my/goCore/cachePath"
	fileCache.Initialize()
	cron.Start() // Please note, you do not call this on each cron job.  Just run it once
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
		return
	}
	cron.RegisterTopOfHourJob(cb)
	c := make(chan string)
	go func() {
		time.Sleep(1 * time.Hour)
		c <- "done"
	}()
	exiting := <-c
	log.Println(exiting)
}
