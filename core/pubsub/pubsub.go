// Package pubsub provides simple pub/sub functionality for people who register on a given topic, functions will be called back in a thread-safe manner
package pubsub

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

//SubscriptionCallback is the callback function for published messgages for use in your interfaces
type SubscriptionCallback func(key string, x interface{})
type subscriptionCallbacks struct {
	sync.RWMutex
	callbacks []SubscriptionCallback
}

// Appends an item to the concurrent slice
func (subscription *subscriptionCallbacks) append(callback SubscriptionCallback) {
	subscription.Lock()
	defer subscription.Unlock()

	subscription.callbacks = append(subscription.callbacks, callback)
}

//iter over the subscription callbacks
func (subscription *subscriptionCallbacks) iter() <-chan SubscriptionCallback {
	c := make(chan SubscriptionCallback)

	f := func() {
		subscription.Lock()
		defer subscription.Unlock()
		for index := range subscription.callbacks {
			c <- subscription.callbacks[index]
		}
		close(c)
	}
	go f()

	return c
}

var subscribers sync.Map

//Subscribe to a publisher message key and function to call
func Subscribe(key string, callback SubscriptionCallback) {
	subscriptionObj, ok := subscribers.Load(key)
	if ok {
		subscription := subscriptionObj.(*subscriptionCallbacks)
		subscription.append(callback)
	} else {
		scs := new(subscriptionCallbacks)
		scs.append(callback)
		subscribers.Store(key, scs)
	}
}

//Publish a message with a payload.
func Publish(key string, x interface{}) {
	go pub(key, x)
}

func pub(key string, x interface{}) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	subscriptionObj, ok := subscribers.Load(key)
	if ok {
		subscription := subscriptionObj.(*subscriptionCallbacks)
		for callback := range subscription.iter() {
			go func(keyLocal string, xLocal interface{}) {
				defer func() {
					if r := recover(); r != nil {
						log.Println("Panic Recovered in pub.  You have some bad code here: "+string(debug.Stack()), fmt.Sprintf("%+v", r))
						return
					}
				}()
				callback(keyLocal, xLocal)
			}(key, x)
		}
	}
}
