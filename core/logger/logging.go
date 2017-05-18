package logger

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/fatih/color"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Color int

var VerboseBornAndDeadGophers bool
var TotalSystemGoRoutines int32
var RunningGophers []string
var GopherTimeRunning map[string]time.Time
var gopherMutex sync.RWMutex

const (
	RED     = 1
	GREEN   = 2
	YELLOW  = 3
	BLUE    = 4
	MAGENTA = 5
	CYAN    = 6
	WHITE   = 7
)

func init() {
	rand.Seed(time.Now().UnixNano())
	RunningGophers = utils.Array()
	GopherTimeRunning = make(map[string]time.Time, 0)
	go func() {
		time.Sleep(time.Minute * 1)
		for {
			serverSettings.WebConfigMutex.RLock()
			if serverSettings.WebConfig.Application.LogGophers {
				ViewRunningGophers()
			}
			serverSettings.WebConfigMutex.RUnlock()
			time.Sleep(time.Minute * 1)
		}
	}()
}

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
		TimeTrack(time.Now(), time.Now().String()+" "+id+" finished ["+routineDesc+"] died ;.-(")
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
      ` + extensions.IntToString(len(RunningGophers)) + ` Gophers workin up in here!
`)
		for i, gopher := range RunningGophers {
			val, ok := GopherTimeRunning[gopher]
			var timeRunning string
			if ok {
				timeRunning = " (" + time.Since(val).String() + " elapsed)"
			}
			log.Println("#" + extensions.IntToString(i) + ":" + gopher + timeRunning)
		}
	}
	gopherMutex.RUnlock()
}

func GoRoutineLogger(fn func(), routineDesc string) {
	if serverSettings.WebConfig.Application.LogGophers {
		id := getGopherGender() + utils.RandStringRunes(5)
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
	fn()
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	//if elapsed.Seconds() > 1 {
	//}
	log.Printf("<Timing>%s took %s</Timing>", name, elapsed)
}

func TimeTrackQuery(start time.Time, name string, collection *mgo.Collection, m bson.M, q *mgo.Query) {
	elapsed := time.Since(start)
	log.Println("<Timing>")
	log.Println()
	log.Printf("%#v", collection)
	log.Println()
	log.Println()
	log.Printf("%#v", m)
	log.Println()
	log.Println()
	log.Printf("%#v", q)
	log.Println()
	log.Println()

	log.Printf("%s took %s", name, elapsed)
	log.Println()
	log.Println("</Timing>")
}
