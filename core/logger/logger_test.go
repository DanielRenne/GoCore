package logger_test

import (
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core/logger"
)

// The simplest use of a TimeTrack caller is to simply call it with a start time and a message.
func ExampleTimeTrack() {
	/*
		import (
			"log"
			"time"

			"github.com/DanielRenne/GoCore/core/logger"
		)
	*/
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered ", r)
			return
		}
		log.Println(logger.TimeTrack(start, "some long running function"))
	}()
	time.Sleep(time.Second * 5)
}

// Heres how you can run a go routine and log if its still running every N seconds (set through serverSettings.WebConfig.Application.LogGopherInterval)
func ExampleGoRoutineLogger() {
	/*
		import (
			"log"
			"time"

			"github.com/DanielRenne/GoCore/core/logger"
		)
	*/
	go logger.GoRoutineLogger(func() {
		for {
			log.Println("test")
			time.Sleep(time.Second * 1)
		}
	}, "some long running function")
	/*
		Output: stdout

	*/
}
