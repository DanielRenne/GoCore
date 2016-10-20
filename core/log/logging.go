package log

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	fmt.Print("<Timing>")
	fmt.Printf("%s took %s", name, elapsed)
	fmt.Println("</Timing>")
}

func TimeTrackQuery(start time.Time, name string, collection *mgo.Collection, m bson.M, q *mgo.Query) {
	elapsed := time.Since(start)
	fmt.Println("<Timing>")
	fmt.Println()
	fmt.Printf("%#v", collection)
	fmt.Println()
	fmt.Println()
	fmt.Printf("%#v", m)
	fmt.Println()
	fmt.Println()
	fmt.Printf("%#v", q)
	fmt.Println()
	fmt.Println()

	fmt.Printf("%s took %s", name, elapsed)
	fmt.Println()
	fmt.Println("</Timing>")
}
