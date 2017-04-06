package logger

import (
	"github.com/fatih/color"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Color int

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

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	if elapsed.Seconds() > 1 {
	}
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
