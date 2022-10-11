package logger_test

import (
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
)

// The simplest use of a TimeTrack caller is to simply call it with a start time and a message.
func ExampleTimeTrack() {
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
	serverSettings.WebConfig.Application.LogGophers = true
	serverSettings.WebConfig.Application.LogGopherInterval = 5
	go logger.GoRoutineLogger(func() {
		for {
			log.Println("test")
			time.Sleep(time.Second * 100)
		}
	}, "some long running function")
	time.Sleep(time.Second * 20)
	log.Println("Done")
	/*
		Output:
		2022/10/05 00:20:12 test
		2022/10/05 00:20:22
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
		at 2022-10-05 00:20:22.830692815 +0000 UTC m=+10.010449585 1 Gophers workin up in here!

		2022/10/05 00:20:22 #0:Mr.ragjF-> (some long running function) (10.008959946s elapsed)
		2022/10/05 00:20:27
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
		at 2022-10-05 00:20:27.834213444 +0000 UTC m=+15.013970230 1 Gophers workin up in here!

		2022/10/05 00:20:27 #0:Mr.ragjF-> (some long running function) (15.012394763s elapsed)
		2022/10/05 00:20:32 Done

	*/
}
