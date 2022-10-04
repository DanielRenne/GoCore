// Package logger provides a simple logging package for GoCore.
// It can also log running goRoutines, track time of execution easily and tail files
// Note: To view running gopher logs, you must set serverSettings.WebConfig.Application.LogGophers to true so it will print out the gopher logs at set interval of serverSettings.WebConfig.Application.LogGopherInterval (set this as well)
package logger

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/fatih/color"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Color is a color type
type Color int

// VerboseBornAndDeadGophers is a flag to turn on and off the verbose logging of gophers.
var VerboseBornAndDeadGophers bool

// TotalSystemGoRoutines is a counter of all the go routines running in the system.
var TotalSystemGoRoutines int32

// RunningGophers is a list of all the gophers currently running.
var RunningGophers []string

// GopherTimeRunning is a map of all the gophers currently running and the time they started.
var GopherTimeRunning map[string]time.Time

var gopherMutex sync.RWMutex

const (
	// RED is a color constant
	RED = 1
	// GREEN is a color constant
	GREEN = 2
	// YELLOW is a color constant
	YELLOW = 3
	// BLUE is a color constant
	BLUE = 4
	// MAGENTA is a color constant
	MAGENTA = 5
	// CYAN is a color constant
	CYAN = 6
	// WHITE is a color constant
	WHITE = 7
)

func init() {
	rand.Seed(time.Now().UnixNano())
	RunningGophers = utils.Array()
	GopherTimeRunning = make(map[string]time.Time, 0)
	go func() {
		time.Sleep(time.Second * 15)
		for {
			serverSettings.WebConfigMutex.RLock()
			if serverSettings.WebConfig.Application.LogGophers {
				ViewRunningGophers()
			}
			sleepTime := serverSettings.WebConfig.Application.LogGopherInterval
			serverSettings.WebConfigMutex.RUnlock()
			time.Sleep(time.Second * time.Duration(sleepTime))
		}
	}()
}

//Log is a wrapper for the standard log package.  Pass in unlimited number of parameters.
func Log(dataValues ...interface{}) {
	for _, value := range dataValues {
		Message(fmt.Sprintf("%+v", value), CYAN)
	}
}

//Message takes in a string and a color and prints it to the console.
func Message(message string, c Color) {
	switch c {
	case 1:
		color.Red(message)
	case 2:
		color.Green(message)
	case 3:
		color.Yellow(message)
	case 4:
		color.Blue(message)
	case 5:
		color.Magenta(message)
	case 6:
		color.Cyan(message)
	case 7:
		color.White(message)
	}
}

func deferGoRoutine(routineDesc string, goRoutineIdStarted int32, id string) {
	if VerboseBornAndDeadGophers {
		log.Println(time.Now().String() + " " + id + " finished [" + routineDesc + "] died ;.-(")
	}
	atomic.AddInt32(&TotalSystemGoRoutines, -1)
	gopherMutex.Lock()
	RunningGophers = utils.ArrayRemove(RunningGophers, id+"-> ("+routineDesc+")")
	gopherMutex.Unlock()

}

func getGopherGender() string {
	if rand.Intn(2) == 1 {
		return "Mrs."
	} else {
		return "Mr."
	}
}

//ViewRunningGophers prints out all the gophers currently running in the system who have been wrapped in GoRoutineLogger or GoRoutineLoggerWithId
func ViewRunningGophers() {
	gopherMutex.RLock()
	if len(RunningGophers) > 0 {
		log.Println(`
           ,_---~~~~~----._
    _,,_,*^____      _____ -g--"*,
   / __/ /'     ^.  /      \ ^@q  f
  [  @f | @))    |  | @))   l  0 _/
   \ /   \~____ / __ \_____/    \
    |           _l__l_           I
    }          [______]           I
    |            | | |            |
    ]             ~ ~             |
    |                            |
     |                           |
at ` + time.Now().String() + " " + extensions.IntToString(len(RunningGophers)) + ` Gophers workin up in here!
`)
		for i, gopher := range RunningGophers {
			val, ok := GopherTimeRunning[gopher]
			var timeRunning string
			if ok {
				timeRunning = " (" + time.Since(val).String() + " elapsed)"
			}
			log.Println("#" + extensions.IntToString(i) + ":" + gopher + timeRunning)
		}
	} else {
		log.Println(time.Now().String() + " no gophers in memory yay!")
	}
	gopherMutex.RUnlock()
}

//GoRoutineLoggerWithId is a wrapper for go routines that will log the start and end of the go routine.  Pass in a function to be executed in the go routine.
func GoRoutineLoggerWithId(fn func(), routineDesc string, Id string) {
	if serverSettings.WebConfig.Application.LogGophers {
		id := getGopherGender()
		if Id == "" {
			id += utils.RandStringRunes(5)
		} else {
			id += Id
		}
		gopherMutex.Lock()
		descId := id + "-> (" + routineDesc + ")"
		GopherTimeRunning[descId] = time.Now()
		RunningGophers = append(RunningGophers, descId)
		gopherMutex.Unlock()
		atomic.AddInt32(&TotalSystemGoRoutines, 1)
		goRoutineIdStarted := atomic.LoadInt32(&TotalSystemGoRoutines)
		defer deferGoRoutine(routineDesc, goRoutineIdStarted, id)
		if VerboseBornAndDeadGophers {
			log.Println(time.Now().String() + " " + id + " is starting to [" + routineDesc + "]")
		}
	}
	if fn != nil {
		fn()
	}
}

//GoRoutineLogger is a wrapper for go routines that will log the start and end of the go routine.  Pass in a function to be executed in the go routine.
func GoRoutineLogger(fn func(), routineDesc string) {
	GoRoutineLoggerWithId(fn, routineDesc, "")
}

//TimeTrack is typically called in your defer function to log the time it took to execute a function.  But can be used anywhere.
func TimeTrack(start time.Time, name string) (log string) {
	elapsed := time.Since(start)
	//if elapsed.Seconds() > 1 {
	//}
	return fmt.Sprintf("%s took %s", name, elapsed)
}

//TimeTrackQuery is meant to be used in conjunction with the mgo package.  It will log the time it took to execute a query and the query itself.
func TimeTrackQuery(start time.Time, name string, collection *mgo.Collection, m bson.M, q *mgo.Query) (log string) {
	elapsed := time.Since(start)
	log += "\n\n"
	log += fmt.Sprintf("%#v", collection) + "\n\n"
	log += fmt.Sprintf("%#v", m) + "\n\n"
	log += fmt.Sprintf("%#v", q) + "\n\n"
	log += fmt.Sprintf("%s took %s", name, elapsed) + "\n\n"
	return
}

// Tail will return the last n bytes of a file
func Tail(path string, length int64) (data string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	buf := make([]byte, length)
	stat, err := os.Stat(path)
	if err != nil {
		return
	}
	start := stat.Size() - length
	if start > 0 {
		_, err = file.ReadAt(buf, start)
		if err == nil {
			data = string(buf)
		}
	} else {
		_, err := file.Read(buf)
		if err == nil {
			data = string(buf)
		}
	}
	return
}
