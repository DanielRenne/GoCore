package model

import (
	"encoding/json"
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sync"
)

var UsersHistory modelUsersHistory

type modelUsersHistory struct{}

var mongoUsersHistoryCollection *mgo.Collection
var collectionUsersHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&UsersHistory)
	collectionUsersHistoryMutex = &sync.RWMutex{}
}

type UserHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelUsersHistory) SetCollection(mdb *mgo.Database) {
	collectionUsersHistoryMutex.Lock()
	mongoUsersHistoryCollection = mdb.C("UsersHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoUsersHistoryCollection.Create(&ci)
	collectionUsersHistoryMutex.Unlock()
}

func (obj modelUsersHistory) Query() *Query {
	var query Query
	for {
		collectionUsersHistoryMutex.RLock()
		collection := mongoUsersHistoryCollection
		collectionUsersHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionUsersHistoryMutex.RLock()
	collection := mongoUsersHistoryCollection
	collectionUsersHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelUsersHistory) Rollback(transactionId string) error {
	var rows []UserHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity User
		err = rollbackEntity.Unmarshal(originalEntityData)

		if err != nil {
			return err
		}

		rollbackType := entity.Type

		if rollbackType == TRANSACTION_CHANGETYPE_DELETE || rollbackType == TRANSACTION_CHANGETYPE_UPDATE {
			err = rollbackEntity.Save()
		} else {
			err = rollbackEntity.Delete()
		}

		if err == nil {
			dbServices.CollectionCache{}.Remove("User", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelUsersHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection UsersHistory:")
	for key, value := range dbServices.GetDBIndexes(UserHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionUsersHistoryMutex.RLock()
		collection := mongoUsersHistoryCollection
		collectionUsersHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for UserHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for UserHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelUsersHistory) New() *UserHistoryRecord {
	return &UserHistoryRecord{}
}

func (obj *UserHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj UserHistoryRecord
	row := modelUsersHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *UserHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionUsersHistoryMutex.RLock()
	changeInfo, err := mongoUsersHistoryCollection.UpsertId(objectId, &self)
	collectionUsersHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for UserHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *UserHistoryRecord) Reflect() []Field {
	return Reflect(UserHistoryRecord{})
}

func (self *UserHistoryRecord) Delete() error {
	collectionUsersHistoryMutex.RLock()
	collection := mongoUsersHistoryCollection
	collectionUsersHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *UserHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *UserHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *UserHistoryRecord) GetType() int {
	return self.Type
}

func (self *UserHistoryRecord) GetData() string {
	return self.Data
}

func (self *UserHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *UserHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *UserHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *UserHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *UserHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *UserHistoryRecord) GetId() string {
	return self.Id.Hex()
}
