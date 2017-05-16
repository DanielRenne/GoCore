package logger

import (
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/fatih/color"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync/atomic"
	"time"
)

type Color int

var TotalSystemGoRoutines int32

const (
	RED     = 1
	GREEN   = 2
	YELLOW  = 3
	BLUE    = 4
	MAGENTA = 5
	CYAN    = 6
	WHITE   = 7
)

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
	TimeTrack(time.Now(), "["+extensions.Int32ToString(goRoutineIdStarted)+"] total gopher(s)->"+id+"Mr."+id+" finished ["+routineDesc+"] died on "+time.Now().String())
	atomic.AddInt32(&TotalSystemGoRoutines, -1)
}

func GoRoutineLogger(fn func(), routineDesc string) {
	if serverSettings.WebConfig.Application.LogGophers {
		id := utils.RandStringRunes(5)
		atomic.AddInt32(&TotalSystemGoRoutines, 1)
		goRoutineIdStarted := atomic.LoadInt32(&TotalSystemGoRoutines)
		defer deferGoRoutine(routineDesc, goRoutineIdStarted, id)
		log.Println("[" + extensions.Int32ToString(goRoutineIdStarted) + "] total gopher(s)->Mr." + id + " is starting to [" + routineDesc + "] on " + time.Now().String())
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
