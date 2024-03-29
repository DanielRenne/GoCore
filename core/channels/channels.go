// Package channels provides a simple way to wait for a signal in a thread-safe manner
package channels

import (
	"log"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/atomicTypes"

	"github.com/DanielRenne/GoCore/core/extensions"
)

// Queue provides a factory to queue channels sequentially and pop / signal them one at a time in a daisy chain.
type Queue struct {
	// Defaults to 10000 ms if none passed
	TimeoutMilliseconds *atomicTypes.AtomicInt
	any                 atomicTypes.AtomicBool
	channels            sync.Map
}

// Signal stub function. Use Queue.Signal() instead
func Signal(x interface{}) {
	q := Queue{}
	q.Signal(x)
}

// Signal will only signal the first item in the queue.
func (q *Queue) Signal(x interface{}) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at channels.Signal():  ", r)
			return
		}
	}()

	anyMoreInRange := false
	signaled := false
	q.channels.Range(func(key interface{}, value interface{}) bool {
		if signaled {
			anyMoreInRange = true
			return false
		}

		signaled = true
		c := value.(chan interface{})
		go func() {
			q.channels.Delete(key)
		}()
		c <- x
		return true
	})
	q.any.Set(anyMoreInRange)
}

// Any will return true if there are any current channels waiting.
func (q *Queue) Any() (any bool) {
	any = q.any.Get()
	return
}

// Wait will return a channel for your function to wait on.
func (q *Queue) Wait(x interface{}) (c chan interface{}, any bool) {
	if q.TimeoutMilliseconds == nil {
		q.TimeoutMilliseconds = &atomicTypes.AtomicInt{}
		q.TimeoutMilliseconds.Set(10000)
	}
	timeoutMs := q.TimeoutMilliseconds.Get()
	any = q.any.Get()
	q.any.Set(true)
	c = make(chan interface{})
	randomValue := extensions.Random(0, 15)
	q.channels.Store(randomValue, c)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic Recovered at channels.Wait->Timeout():  ", r)
				return
			}
		}()
		time.Sleep(time.Millisecond * time.Duration(timeoutMs))
		_, ok := q.channels.Load(randomValue)
		if ok {
			q.channels.Delete(randomValue)
			c <- x
		}

	}()
	return
}
