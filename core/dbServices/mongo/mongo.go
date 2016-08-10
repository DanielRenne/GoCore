package main

import (
	"fmt"
	_ "github.com/DanielRenne/GoCore/core/dbServices/mongo/acct"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var MongoSession *mgo.Session
var MongoDB *mgo.Database
var mongoCollection *mgo.Collection

type Account struct {
	Id      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Company string        `bson:"company" json:"company"`
}

func main() {
	openMongo()
}

func openMongo() {

	var err error
	MongoSession, err := mgo.Dial("10.0.0.8:27017") // open an connection -> Dial function
	if err != nil {                                 //  if you have a
		fmt.Println("Failed to create or open mongo Database:\n\t" + err.Error())
		return
	}

	MongoSession.SetMode(mgo.Monotonic, true) // Optional. Switch the session to a monotonic behavior.

	MongoDB = MongoSession.DB("studio")
	fmt.Printf("%+v\n", MongoDB)

	mongoCollection = MongoDB.C("accounts")

	var acct = Account{
		Company: "Test Inc.",
	}

	objectId := bson.NewObjectId()
	if acct.Id != "" {
		fmt.Println("dasfadsf")
		objectId = acct.Id
	}

	fmt.Printf("%+v\n", acct.Id)

	// acct.Id = bson.NewObjectId()

	changeInfo, err := mongoCollection.UpsertId(objectId, &acct)

	if err != nil { //  if you have a
		fmt.Println("Failed to upsertId for Account:\n\t" + err.Error())
		return
	}

	acct.Id = changeInfo.UpsertedId.(bson.ObjectId)

	fmt.Printf("%+v\n", changeInfo)
	fmt.Printf("%+v\n", acct)

	acct.Company = "Whatever"

	changeInfo, err = mongoCollection.UpsertId(acct.Id, &acct)

	if err != nil { //  if you have a
		fmt.Println("Failed to upsertId for Account:\n\t" + err.Error())
		return
	}

	fmt.Printf("%+v\n", changeInfo)
	fmt.Printf("%+v\n", acct)

	for {

	}

}
