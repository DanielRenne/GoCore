package pubsub_test

import (
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/pubsub"
)

// TestPublish is a test struct for pubsub
type TestPublish struct {
	DeviceID string `json:"deviceId"`
	Value    bool   `json:"value"`
}

// ExamplePubSub ... simple pub sub examples
func ExamplePublish() {
	// Place me anywhere in your code
	pubsub.Subscribe("test", func(topic string, data interface{}) {
		// Do something with the data
		core.Debug.Dump("Dumping the data emojis wont work on windows because windows cmd line sucks, but on linux or mac it will kiss you 💋", topic, data)
		core.Debug.Dump("oh cool byte array dumps when non printables exist", "\x00 hello nulls!")
	})
	log.Print("Sleeping for 5 seconds")
	time.Sleep(time.Second * 5)
	core.Debug.Dump("Calling Publish")
	pubsub.Publish("test", TestPublish{DeviceID: "123", Value: true})
	log.Print("Sleeping for 2 seconds to allow the pubsub to run before it exits the program")
	time.Sleep(time.Second * 2)
	/* Output:
	2022/10/03 21:25:30 Sleeping for 5 seconds
	!!!!!!!!!!!!! DEBUG 2022-10-03 21:25:35.506013!!!!!!!!!!!!!


	#### string                                  [len:15]####
	Calling Publish

	!!!!!!!!!!!!! ENDDEBUG 2022-10-03 21:25:35.506013!!!!!!!!!!!!!
	2022/10/03 21:25:35 Sleeping for 2 seconds to allow the pubsub to run before it exits the program
	*/
}
