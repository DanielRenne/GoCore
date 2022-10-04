package channels_test

import (
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core/atomicTypes"
	"github.com/DanielRenne/GoCore/core/channels"
	"github.com/DanielRenne/GoCore/core/extensions"
)

// ExampleSignal: Two examples of signal and wait. One timesout and the other does not.
func ExampleSignal() {
	/*
		import (
			"log"
			"time"

			"github.com/DanielRenne/GoCore/core/atomicTypes"
			"github.com/DanielRenne/GoCore/core/channels"
			"github.com/DanielRenne/GoCore/core/extensions"
		)

	*/
	timeout := atomicTypes.AtomicInt{}
	timeout.Set(1000)
	channelQueue := channels.Queue{
		TimeoutMilliseconds: &timeout,
	}

	type myCoolStruct struct {
		Hello       string
		WasReturned bool
	}
	c, any := channelQueue.Wait(myCoolStruct{})

	if any == false {
		go func() {
			log.Println("Sleep 500 ms")
			time.Sleep(500 * time.Millisecond)
			// from some other package or function signal the channel
			channelQueue.Signal(myCoolStruct{
				Hello:       "World",
				WasReturned: true,
			})
		}()
	}

	data := <-c
	log.Println("done, cast your interface to the proper type back to the caller.  Response data: " + data.(myCoolStruct).Hello + "  WasReturned: " + extensions.BoolToString(data.(myCoolStruct).WasReturned))

	timeout2 := atomicTypes.AtomicInt{}
	timeout2.Set(100)
	channelQueue2 := channels.Queue{
		TimeoutMilliseconds: &timeout2,
	}

	c, any = channelQueue2.Wait("initial data")

	if any == false {
		go func() {
			log.Println("Sleep 5 seconds for timeout")
			time.Sleep(5000 * time.Millisecond)
			// from some other package or function signal the channel
			channelQueue.Signal("return will never make it")
		}()
	}

	data2 := <-c
	log.Println("done, cast your interface to the proper type back to the caller.  Response data: " + data2.(string))
	/*
		Output:
			2022/10/04 16:29:57 Sleep 500 ms
			2022/10/04 16:29:57 done, cast your interface to the proper type back to the caller.  Response data: World  WasReturned: true
			2022/10/04 16:29:57 Sleep 5 seconds for timeout
			2022/10/04 16:29:57 done, cast your interface to the proper type back to the caller.  Response data: initial data
	*/
}
