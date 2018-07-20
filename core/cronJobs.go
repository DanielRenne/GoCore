package core

import (
	"log"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/fileCache"
)

//CronJobs provides a cron job engine for golang function callbacks to schedule events to execute at specific times of the day and week.
var CronJobs cronJobs

//RecurringType defines a type of cron job.
type RecurringType int

var recurringJobs recurringJobsSync

const (
	CRON_TOP_OF_MINUTE RecurringType = iota
	CRON_TOP_OF_HOUR
	CRON_TOP_OF_DAY
	CRON_TOP_OF_30_SECONDS
	CRON_TOP_OF_SECOND
)

type cronJobs struct {
}

type onDemandJobsSync struct {
	sync.RWMutex
	items []OnDemandEvent
}

type recurringJobsSync struct {
	sync.RWMutex
	items []recurringEvent
}

//CronJob entity provides details of the cron job to be executed.
type CronJob struct {
}

//CronJobEvent is used as the callback function for the event.
type OnDemandEvent func(id string, eventTime time.Time, context interface{})

//CronEvent is a callback function called by the cron job engine.
type RecurringEvent func(eventDate time.Time)
type OneTimeEvent func() bool

type recurringEvent struct {
	Type  RecurringType
	Event RecurringEvent
}

//Starts the cron job engine.
func (jobs *cronJobs) ShouldRunEveryTenMinutes(x time.Time) (run bool) {
	if x.Minute() == 0 || x.Minute() == 10 || x.Minute() == 20 || x.Minute() == 30 || x.Minute() == 40 || x.Minute() == 50 {
		run = true
	}
	return
}

func (jobs *cronJobs) ShouldRunEveryFiveMinutes(x time.Time) (run bool) {
	if x.Minute() == 0 || x.Minute() == 5 || x.Minute() == 10 || x.Minute() == 15 || x.Minute() == 20 || x.Minute() == 25 || x.Minute() == 30 || x.Minute() == 35 || x.Minute() == 40 || x.Minute() == 45 || x.Minute() == 50 || x.Minute() == 55 {
		run = true
	}
	return
}

func (jobs *cronJobs) ShouldRunEvery15Minutes(x time.Time) (run bool) {
	if x.Minute() == 0 || x.Minute() == 15 || x.Minute() == 30 || x.Minute() == 45 {
		run = true
	}
	return
}

func (jobs *cronJobs) Start() {
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
					go callRecurringEvents(CRON_TOP_OF_MINUTE, tm)
					callTopMinute = false
				}

				if callTop30Seconds {
					go callRecurringEvents(CRON_TOP_OF_30_SECONDS, tm)
					callTop30Seconds = false
				}
			}

			if sec == 30 { //Top of the Minute && Top of 30 Seconds
				if callTop30Seconds {
					go callRecurringEvents(CRON_TOP_OF_30_SECONDS, tm)
					callTop30Seconds = false
				}
			}

			if sec == 0 && min == 0 { //Top of the Hour
				if callTopHour {
					go callRecurringEvents(CRON_TOP_OF_HOUR, tm)
					callTopHour = false
				}
			}
			if sec == 0 && min == 0 && hour == 0 { //Top of the Day
				if callTopDay {
					go callRecurringEvents(CRON_TOP_OF_DAY, tm)
					callTopDay = false
				}
			}
			if sec == 1 {
				callTopMinute = true
				callTopHour = true
				callTopDay = true
				callTop30Seconds = true
			}
			if sec == 31 {
				callTop30Seconds = true
			}

			if previousSec != sec {
				previousSec = sec
				go callRecurringEvents(CRON_TOP_OF_SECOND, tm)
			}
		}
	}()

}

//Register provides a method to register for a callback that is called at the start of the cron job engine and 5 seconds before each day occures.
func (jobs *cronJobs) RegisterRecurring(t RecurringType, callback RecurringEvent) {
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

func (jobs *cronJobs) ExecuteOneTimeJob(jobName string, callback OneTimeEvent) {
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

func processRecurringTick(tm time.Time) {

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
