package br

// commented libraries are for schedules LoadDays
import (
	"fmt"
	//"runtime"
	"runtime/debug"
	"sync"
	"time"

	//"github.com/DanielRenne/GoCore/core"
	//"github.com/DanielRenne/GoCore/core/extensions"
	//"github.com/DanielRenne/goCoreAppTemplate/constants"
	//"github.com/DanielRenne/goCoreAppTemplate/queries"
	"github.com/DanielRenne/goCoreAppTemplate/scheduleEngine"
	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
)

type schedulesBr struct{}

type schedules struct {
	Schedules []schedule `json:"Schedules"`
}

type schedule struct {
	StartTime         string `json:"StartTime"`
	WeeklyValidations []int  `json:"WeeklyValidations"`
	MacroId           string `json:"MacroId"`
}

var timeZone string
var schedulesLock sync.RWMutex
var cachedLoc *time.Location

func init() {
	timeZone = "America/Los_Angeles"
	cachedLoc, _ = time.LoadLocation(timeZone)
	scheduleEngine.SetScheduleDay(Schedules.LoadDay)
	scheduleEngine.SetLocationCallback(Schedules.GetLocation)
}

func (self schedulesBr) UpdateLinuxToGMT() {
	if runtime.GOOS == "linux" {
		SetTimeZone("UTC")
	}
}

func (self schedulesBr) LoadDay(t time.Time) {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("Panic->br->schedules->LoadDay", "Panic Recovered:  "+fmt.Sprintf("%+v", r))
			return
		}
	}()

	//session_functions.Log("Loading Daily Schedule", t.String())
	//tUnix := t.Unix()
	//
	//self.UpdateLinuxToGMT()
	//
	//timeZoneSetting, _ := queries.ServerSettings.ById(constants.SERVER_SETTING_TIMEZONE)
	//if timeZoneSetting.Value != "" && timeZoneSetting.Value != "0" {
	//	self.SetTimeZone(timeZoneSetting.Value)
	//}
	//
	//loc, err := time.LoadLocation(self.GetTimeZone())
	//if err != nil {
	//	loc, _ = time.LoadLocation("America/Los_Angeles")
	//	self.SetTimeZone("America/Los_Angeles")
	//}
	//
	//GMTDate := t.In(loc)
	//
	//dateString := ""
	//dateString += extensions.IntToString(GMTDate.Day()) + "-" + GMTDate.Month().String() + "-" + extensions.IntToString(GMTDate.Year()) + " "
	//
	//roomSchedules, err := queries.RoomFeatures.GetSchedules()
	//if err != nil {
	//	session_functions.Log("Error->br->schedules->LoadDay", "Failed to get roomSchedules:  "+err.Error())
	//	return
	//}
	//
	//// var all schedules
	////
	//// err := extensions.ReadFileAndParse(serverSettings.APP_LOCATION+"/dist/schedules.json", &all)
	//// if err != nil {
	//// 	session_functions.Log("Error->br->schedules->LoadDay", err.Error())
	//// 	return
	//// }
	//
	//session_functions.Log("Clearing Schedule Engine", "")
	//scheduleEngine.Clear()
	//session_functions.Log("Clearing Schedule Engine Completed", "")
	//for i := range roomSchedules {
	//	item := roomSchedules[i]
	//	format := "2-January-2006 03:04:05 PM"
	//	tme, err := time.ParseInLocation(format, dateString+item.Value1, loc)
	//
	//	// session_functions.Log("Schedules", fmt.Sprintf("%+v", tme))
	//	// session_functions.Log("Schedules", fmt.Sprintf("%+v", t))
	//
	//	var err2 error
	//	var tme2 time.Time
	//	if item.Value2 != "0000-00-00" {
	//		tme2, err2 = time.ParseInLocation(format, item.Value2+" 12:00:00 AM", loc)
	//	}
	//	isInFuture := tme.Unix() >= tUnix
	//	isLessThanEndDate := (item.Value2 == "0000-00-00" || (item.Value2 != "0000-00-00" && tme2.Unix() <= t.Unix()))
	//	if err == nil && err2 == nil && isInFuture && isLessThanEndDate {
	//		weekDay := int(t.Weekday())
	//		if item.ScheduleRule == "weekly" {
	//			var isMatchedDay bool
	//			for j := range item.DayValidations {
	//				v := item.DayValidations[j]
	//				if v == weekDay {
	//					isMatchedDay = true
	//					// session_functions.Dump(item.MacroId, tme)
	//					session_functions.Log("Schedule Macro Added", item.MacroId)
	//					scheduleEngine.AddItem(tme, item.MacroId, self.TriggerSchedule)
	//				}
	//			}
	//			if !isMatchedDay {
	//				session_functions.Log("Schedule skipped due to weekday not matching on ", "Schedule Id: "+item.Id.Hex()+" (Macro "+item.MacroId+") v/DayValidations="+core.Debug.GetDump(item.DayValidations)+" weekDay="+extensions.IntToString(weekDay))
	//			}
	//		} else if item.ScheduleRule == "daily" {
	//			// session_functions.Dump(item.MacroId, tme)
	//			session_functions.Log("Schedule Macro Added", item.MacroId)
	//			scheduleEngine.AddItem(tme, item.MacroId, self.TriggerSchedule)
	//		}
	//	} else {
	//		session_functions.Log("Schedule skipped on ", "Schedule Id: "+item.Id.Hex()+" (Macro "+item.MacroId+")")
	//		session_functions.Log("Schedule skipped isLessThanEndDate", extensions.BoolToString(isLessThanEndDate))
	//		session_functions.Log("Schedule skipped isInFuture", extensions.BoolToString(isInFuture))
	//		if !isInFuture {
	//			session_functions.Log("Schedule skipped scheduled unix time", extensions.Int64ToString(tme.Unix()))
	//			session_functions.Log("Schedule skipped Current unix time", extensions.Int64ToString(tUnix))
	//		}
	//		session_functions.Log("Schedule skipped err", extensions.BoolToString(err == nil))
	//		session_functions.Log("Schedule skipped err2", extensions.BoolToString(err2 == nil))
	//	}
	//}
	//session_functions.Log("Sorting Schedule Engine", "")
	//scheduleEngine.Sort()
	//session_functions.Log("Sorting Schedule Engine Completed", "")
}

func (self schedulesBr) TriggerSchedule(x interface{}) {
	defer func() {
		if r := recover(); r != nil {
			session_functions.Print("\n\nPanic Stack: " + string(debug.Stack()))
			session_functions.Log("Panic->br->schedules->TriggerMacro", "Panic Recovered:  "+fmt.Sprintf("%+v", r))
			return
		}
	}()
	jobId, ok := x.(string)
	if ok {
		session_functions.Log("Triggering MacroId", jobId)
		//go controlEngine.RunMacro(macroId)
	}
}

func (self schedulesBr) GetTimeZone() (value string) {
	schedulesLock.RLock()
	value = timeZone
	if timeZone == "" {
		value = "America/Los_Angeles"
	}
	schedulesLock.RUnlock()
	return
}

func (self schedulesBr) SetTimeZone(value string) {
	schedulesLock.Lock()
	timeZone = value
	cachedLoc, _ = time.LoadLocation(timeZone)
	schedulesLock.Unlock()
	return
}

func (self schedulesBr) GetLocation() (loc *time.Location) {
	schedulesLock.RLock()
	loc = cachedLoc
	schedulesLock.RUnlock()
	return
}
