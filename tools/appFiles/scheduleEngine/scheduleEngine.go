package scheduleEngine

import (
	"fmt"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"log"
)

const parseFormat = "2-January-2006 03:04:05 PM"

type ScheduleCallback func(x interface{})
type ScheduleDay func(t time.Time)
type LocationCallback func() (loc *time.Location)

type scheduleItem struct {
	Time          time.Time
	CallbackParam interface{}
	Callback      ScheduleCallback
}

type ScheduleItemSlice []scheduleItem

func (p ScheduleItemSlice) Len() int {
	return len(p)
}

func (p ScheduleItemSlice) Less(i, j int) bool {
	return p[i].Time.Before(p[j].Time)
}

func (p ScheduleItemSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type scheduleSync struct {
	sync.RWMutex
	schedules         ScheduleItemSlice
	localTimeZoneTime time.Time
	local             *time.Location
}

var schedulesSync scheduleSync
var scheduleDay ScheduleDay
var locationCallback LocationCallback

func Clear() {
	schedulesSync.Lock()
	schedulesSync.schedules = nil
	schedulesSync.Unlock()
}

func AddItem(t time.Time, x interface{}, callback ScheduleCallback) {
	schedulesSync.Lock()
	var s scheduleItem
	s.Time = t
	s.Callback = callback
	s.CallbackParam = x

	schedulesSync.schedules = append(schedulesSync.schedules, s)
	schedulesSync.Unlock()
}

func Sort() {
	schedulesSync.Lock()
	sort.Sort(schedulesSync.schedules)
	schedulesSync.Unlock()
}

func SetScheduleDay(callback ScheduleDay) {
	scheduleDay = callback
}

func SetLocationCallback(callback LocationCallback) {
	locationCallback = callback
}

func GetSchedules() (schedules ScheduleItemSlice) {
	schedulesSync.RLock()
	schedules = schedulesSync.schedules
	schedulesSync.RUnlock()
	return
}

func GetLocalTime(t time.Time) (localTime time.Time) {
	localTime = t.In(locationCallback())
	return
}

func Trigger(t time.Time) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("\n\nPanic Stack: " + string(debug.Stack()))
			log.Println("Panic->scheduleEngine->Trigger", "Panic Recovered:  "+fmt.Sprintf("%+v", r))
			return
		}
	}()

	localTime := t.In(locationCallback())

	if localTime.Hour() == 0 && localTime.Minute() == 0 && localTime.Second() == 0 {
		dateString := ""
		dateString += extensions.IntToString(t.Day()) + "-" + t.Month().String() + "-" + extensions.IntToString(t.Year()) + " 00:00:00 AM"

		if scheduleDay != nil {
			scheduleDay(t)
		}
	}
	schedulesSync.Lock()
	for i := 0; i < len(schedulesSync.schedules); i++ {
		s := schedulesSync.schedules[i]

		if s.Time.Unix() > t.Unix() {
			break
		}
		if s.Time.Unix() < t.Unix() {
			schedulesSync.schedules = append(schedulesSync.schedules[:i], schedulesSync.schedules[i+1:]...)
			i--
		}
		if s.Time.Unix() == t.Unix() {
			go s.Callback(s.CallbackParam)
			schedulesSync.schedules = append(schedulesSync.schedules[:i], schedulesSync.schedules[i+1:]...)
			i--
		}
	}
	schedulesSync.Unlock()

}
