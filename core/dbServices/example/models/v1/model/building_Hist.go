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

var BuildingsHistory modelBuildingsHistory

type modelBuildingsHistory struct{}

var mongoBuildingsHistoryCollection *mgo.Collection
var collectionBuildingsHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&BuildingsHistory)
	collectionBuildingsHistoryMutex = &sync.RWMutex{}
}

type BuildingHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelBuildingsHistory) SetCollection(mdb *mgo.Database) {
	collectionBuildingsHistoryMutex.Lock()
	mongoBuildingsHistoryCollection = mdb.C("BuildingsHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoBuildingsHistoryCollection.Create(&ci)
	collectionBuildingsHistoryMutex.Unlock()
}

func (obj modelBuildingsHistory) Query() *Query {
	var query Query
	for {
		collectionBuildingsHistoryMutex.RLock()
		collection := mongoBuildingsHistoryCollection
		collectionBuildingsHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionBuildingsHistoryMutex.RLock()
	collection := mongoBuildingsHistoryCollection
	collectionBuildingsHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelBuildingsHistory) Rollback(transactionId string) error {
	var rows []BuildingHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity Building
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
			dbServices.CollectionCache{}.Remove("Building", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelBuildingsHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection BuildingsHistory:")
	for key, value := range dbServices.GetDBIndexes(BuildingHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionBuildingsHistoryMutex.RLock()
		collection := mongoBuildingsHistoryCollection
		collectionBuildingsHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for BuildingHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for BuildingHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelBuildingsHistory) New() *BuildingHistoryRecord {
	return &BuildingHistoryRecord{}
}

func (obj *BuildingHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj BuildingHistoryRecord
	row := modelBuildingsHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *BuildingHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionBuildingsHistoryMutex.RLock()
	changeInfo, err := mongoBuildingsHistoryCollection.UpsertId(objectId, &self)
	collectionBuildingsHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for BuildingHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *BuildingHistoryRecord) Reflect() []Field {
	return Reflect(BuildingHistoryRecord{})
}

func (self *BuildingHistoryRecord) Delete() error {
	collectionBuildingsHistoryMutex.RLock()
	collection := mongoBuildingsHistoryCollection
	collectionBuildingsHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *BuildingHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *BuildingHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *BuildingHistoryRecord) GetType() int {
	return self.Type
}

func (self *BuildingHistoryRecord) GetData() string {
	return self.Data
}

func (self *BuildingHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *BuildingHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *BuildingHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *BuildingHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *BuildingHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *BuildingHistoryRecord) GetId() string {
	return self.Id.Hex()
}
