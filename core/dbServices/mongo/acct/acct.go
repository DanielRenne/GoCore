package acct

import (
	"github.com/DanielRenne/GoCore/core/dbServices"
	"log"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"time"
)

type Accounts struct {
}

func init() {
	go func() {

		for {
			if dbServices.MongoSession != nil {
				log.Println("Building Indexes for MongoDB collection Accounts:")
				return
			}

			time.Sleep(200 * time.Millisecond)

		}
	}()

}
