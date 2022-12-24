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

var FloorsHistory modelFloorsHistory

type modelFloorsHistory struct{}

var mongoFloorsHistoryCollection *mgo.Collection
var collectionFloorsHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&FloorsHistory)
	collectionFloorsHistoryMutex = &sync.RWMutex{}
}

type FloorHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelFloorsHistory) SetCollection(mdb *mgo.Database) {
	collectionFloorsHistoryMutex.Lock()
	mongoFloorsHistoryCollection = mdb.C("FloorsHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoFloorsHistoryCollection.Create(&ci)
	collectionFloorsHistoryMutex.Unlock()
}

func (obj modelFloorsHistory) Query() *Query {
	var query Query
	for {
		collectionFloorsHistoryMutex.RLock()
		collection := mongoFloorsHistoryCollection
		collectionFloorsHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionFloorsHistoryMutex.RLock()
	collection := mongoFloorsHistoryCollection
	collectionFloorsHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelFloorsHistory) Rollback(transactionId string) error {
	var rows []FloorHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity Floor
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
			dbServices.CollectionCache{}.Remove("Floor", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelFloorsHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection FloorsHistory:")
	for key, value := range dbServices.GetDBIndexes(FloorHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionFloorsHistoryMutex.RLock()
		collection := mongoFloorsHistoryCollection
		collectionFloorsHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for FloorHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for FloorHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelFloorsHistory) New() *FloorHistoryRecord {
	return &FloorHistoryRecord{}
}

func (obj *FloorHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj FloorHistoryRecord
	row := modelFloorsHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *FloorHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionFloorsHistoryMutex.RLock()
	changeInfo, err := mongoFloorsHistoryCollection.UpsertId(objectId, &self)
	collectionFloorsHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for FloorHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *FloorHistoryRecord) Reflect() []Field {
	return Reflect(FloorHistoryRecord{})
}

func (self *FloorHistoryRecord) Delete() error {
	collectionFloorsHistoryMutex.RLock()
	collection := mongoFloorsHistoryCollection
	collectionFloorsHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *FloorHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *FloorHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *FloorHistoryRecord) GetType() int {
	return self.Type
}

func (self *FloorHistoryRecord) GetData() string {
	return self.Data
}

func (self *FloorHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FloorHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *FloorHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *FloorHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *FloorHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *FloorHistoryRecord) GetId() string {
	return self.Id.Hex()
}
