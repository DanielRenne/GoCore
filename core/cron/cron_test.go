package cron_test

import (
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/cron"
	"github.com/DanielRenne/GoCore/core/zip"
)

// You may have a situation that upon initializing, you have a user defined configuration to control the interval of the job.  This is how you would do that.
func ExampleRegisterRecurring() {
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

// ExampleRegisterFutureEvent is an example of how to execute a one time event.
func ExampleRegisterFutureEvent() {
	cb := func() {
		log.Println("asdfasdf")
	}
	t := time.Now().Add(time.Second * 5)
	cron.RegisterFutureEvent(t, cb)
}

// ExampleExecuteOneTimeJob is an example of how to execute a one time event.
func ExampleExecuteOneTimeJob() {
	cb := func() bool {
		err := zip.Unzip("test", "test", []string{})
		return err == nil
	}
	cron.ExecuteOneTimeJob("run me and i will save a key of this string and write the boolean of success or not to /usr/local/goCore/jobs/jobs.json for historical purposes", cb)
}

// execute job every 5 minutes
func ExampleShouldRunEveryFiveMinutes() {
	cb := func(currentTime time.Time) {
		// Here you wrap the current time with the helper function so that minutes 1, 2, 3, 4 are skipped, but minutes 0 and 5 are executed.
		if cron.ShouldRunEveryFiveMinutes(currentTime) {
			core.Dump("Job Executed on top of 5 minutes: " + currentTime.String())
		}
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
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
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
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
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
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
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
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
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
	cb := func(currentTime time.Time) {
		core.Dump("Job Executed on: " + currentTime.String())
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
