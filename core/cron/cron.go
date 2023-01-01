// Package cron contains cron jobs and other logging functions
package cron

import (
	"log"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/fileCache"
)

// RecurringType is a type of recurring event (int)
type RecurringType int

var recurringJobs recurringJobsSync
var recurringFutureJobs recurringFutureJobsSync

const (
	//CronTopOfMinute is a cron job type that is called at the top of every minute.
	CronTopOfMinute RecurringType = iota
	//CronTopOfHour is a cron job type that is called at the top of every hour.
	CronTopOfHour
	//CronTopOfDay is a cron job type that is called at the top of every day.
	CronTopOfDay
	//CronTopOf30Seconds is a cron job type that is called at the top of every 30 seconds.
	CronTopOf30Seconds
	//CronTopOfSecond is a cron job type that is called at the top of every second.
	CronTopOfSecond
)

type recurringJobsSync struct {
	sync.RWMutex
	items []recurringEvent
}
type recurringFutureJobsSync struct {
	sync.RWMutex
	items []futureEvent
}

// RecurringEvent is a callback function called by the cron job engine.
type RecurringEvent func(eventDate time.Time)
type FutureEvent func()

// OneTimeEvent is used to schedule a one time event to be executed
type OneTimeEvent func() bool

type recurringEvent struct {
	Type  RecurringType
	Event RecurringEvent
}
type futureEvent struct {
	Time  time.Time
	Event FutureEvent
}

func init() {
	fileCache.Initialize()
	start()
}

// ShouldRunEveryTenMinutes is a helper function that returns true if the time is the top of every 10 minutes. Note, please call cron.()
func ShouldRunEveryTenMinutes(x time.Time) (run bool) {
	if x.Minute() == 0 || x.Minute() == 10 || x.Minute() == 20 || x.Minute() == 30 || x.Minute() == 40 || x.Minute() == 50 {
		run = true
	}
	return
}

// ShouldRunEveryFiveMinutes is a helper function that returns true if the time is the top of every 5 minutes.
func ShouldRunEveryFiveMinutes(x time.Time) (run bool) {
	if x.Minute() == 0 || x.Minute() == 5 || x.Minute() == 10 || x.Minute() == 15 || x.Minute() == 20 || x.Minute() == 25 || x.Minute() == 30 || x.Minute() == 35 || x.Minute() == 40 || x.Minute() == 45 || x.Minute() == 50 || x.Minute() == 55 {
		run = true
	}
	return
}

// ShouldRunEvery15Minutes is a helper function that returns true if the time is the top of every 15 minutes.
func ShouldRunEvery15Minutes(x time.Time) (run bool) {
	if x.Minute() == 0 || x.Minute() == 15 || x.Minute() == 30 || x.Minute() == 45 {
		run = true
	}
	return
}

func start() {
	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {

		callTopMinute := true
		callTopHour := true
		callTopDay := true
		callTop30Seconds := true
		previousSec := 0

		for t := range ticker.C {
			tm := t
			hour, min, sec := t.Clock()
			if sec == 0 { //Top of the Minute && Top of 30 Seconds
				if callTopMinute {
					go callRecurringEvents(CronTopOfMinute, tm)
					callTopMinute = false
				}

				if callTop30Seconds {
					go callRecurringEvents(CronTopOf30Seconds, tm)
					callTop30Seconds = false
				}
			}

			if sec == 30 { //Top of the Minute && Top of 30 Seconds
				if callTop30Seconds {
					go callRecurringEvents(CronTopOf30Seconds, tm)
					callTop30Seconds = false
				}
			}

			if sec == 0 && min == 0 { //Top of the Hour
				if callTopHour {
					go callRecurringEvents(CronTopOfHour, tm)
					callTopHour = false
				}
			}
			if sec == 0 && min == 0 && hour == 0 { //Top of the Day
				if callTopDay {
					go callRecurringEvents(CronTopOfDay, tm)
					callTopDay = false
				}
			}
			if sec == 1 {
				callTopMinute = true
				callTopHour = true
				callTopDay = true
				callTop30Seconds = true
			}
			if sec == 30 {
				callTop30Seconds = true
			}

			if previousSec != sec {
				previousSec = sec
				go callRecurringFutureEvents(time.Now())
				go callRecurringEvents(CronTopOfSecond, tm)
			}
		}
	}()

}

// RegisterTopOf30SecondsJob executes top of 30 seconds job.
func RegisterTopOf30SecondsJob(callback RecurringEvent) {
	go RegisterRecurring(CronTopOf30Seconds, callback)
}

// RegisterTopOfMinuteJob executes top of 30 seconds job.
func RegisterTopOfMinuteJob(callback RecurringEvent) {
	go RegisterRecurring(CronTopOfMinute, callback)
}

// RegisterTopOfSecondJob executes top of 30 seconds job.
func RegisterTopOfSecondJob(callback RecurringEvent) {
	go RegisterRecurring(CronTopOfSecond, callback)
}

// RegisterTopOfDayJob executes top of 30 seconds job.
func RegisterTopOfDayJob(callback RecurringEvent) {
	go RegisterRecurring(CronTopOfDay, callback)
}

// RegisterTopOfHourJob executes top of 30 seconds job.
func RegisterTopOfHourJob(callback RecurringEvent) {
	go RegisterRecurring(CronTopOfHour, callback)
}

// RegisterRecurring provides a method to register for a callback that is called at the start of the cron job engine and 5 seconds before each day occures.
func RegisterRecurring(t RecurringType, callback RecurringEvent) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	recurringJobs.Lock()
	var re recurringEvent
	re.Event = callback
	re.Type = t
	recurringJobs.items = append(recurringJobs.items, re)
	recurringJobs.Unlock()
}

// RegisterFutureEvent provides a method to register for a callback that is called at at some time you specify in the future
func RegisterFutureEvent(t time.Time, callback FutureEvent) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	recurringFutureJobs.Lock()
	var re futureEvent
	re.Event = callback
	re.Time = t
	recurringFutureJobs.items = append(recurringFutureJobs.items, re)
	recurringFutureJobs.Unlock()
}

// ClearRecurringJobs clears all recurring jobs.
// You would do this in a case of some event reconfiguring possibly a dynamic event of RegisterRecurring() where an end user is in control of the execution of the cron jobs
// Then call the function which reads your static configuration or dynamic configuration to setup your cron jobs and there is no need to call Start again if its already been called
func ClearRecurringJobs() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	recurringJobs.Lock()
	var re []recurringEvent
	recurringJobs.items = re
	recurringJobs.Unlock()
}

// ClearFutureJobs clears all future jobs waiting to be executed
func ClearFutureJobs() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	recurringFutureJobs.Lock()
	var re []futureEvent
	recurringFutureJobs.items = re
	recurringFutureJobs.Unlock()
}

// ExecuteOneTimeJob executes a one time job.
func ExecuteOneTimeJob(jobName string, callback OneTimeEvent) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic recover ExecuteOneTimeJob", r)
			return
		}
	}()

	go func() {
		fileCache.Jobs.Lock()
		defer func() {
			fileCache.Jobs.Unlock()
			if r := recover(); r != nil {
				log.Println("panic recover ExecuteOneTimeJob", r)
				return
			}
		}()

		value, ok := fileCache.Jobs.Jobs[jobName]
		if !ok || !value {
			success := callback()
			fileCache.Jobs.Jobs[jobName] = success
			fileCache.WriteJobCacheFile()
		}
	}()

}

func callRecurringEvents(t RecurringType, tm time.Time) {
	recurringJobs.RLock()
	for _, item := range recurringJobs.items {
		i := item
		if i.Type == t {
			go func(e RecurringEvent) {
				e(tm)
			}(i.Event)
		}
	}
	recurringJobs.RUnlock()
}

func remove(s []futureEvent, i int) []futureEvent {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func callRecurringFutureEvents(tm time.Time) {
	recurringFutureJobs.Lock()
	for idx, item := range recurringFutureJobs.items {
		i := item
		if i.Time.Year() == tm.Year() && i.Time.Month() == tm.Month() && i.Time.Day() == tm.Day() && i.Time.Hour() == tm.Hour() && i.Time.Minute() == tm.Minute() && i.Time.Second() == tm.Second() {
			go func(e FutureEvent) {
				e()
			}(i.Event)
			// manage some memory after execution so these things dont just sit forever dormant on each call
			remove(recurringFutureJobs.items, idx)
		}
	}
	recurringFutureJobs.Unlock()
}
