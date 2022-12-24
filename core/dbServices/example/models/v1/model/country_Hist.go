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

var CountriesHistory modelCountriesHistory

type modelCountriesHistory struct{}

var mongoCountriesHistoryCollection *mgo.Collection
var collectionCountriesHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&CountriesHistory)
	collectionCountriesHistoryMutex = &sync.RWMutex{}
}

type CountryHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelCountriesHistory) SetCollection(mdb *mgo.Database) {
	collectionCountriesHistoryMutex.Lock()
	mongoCountriesHistoryCollection = mdb.C("CountriesHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoCountriesHistoryCollection.Create(&ci)
	collectionCountriesHistoryMutex.Unlock()
}

func (obj modelCountriesHistory) Query() *Query {
	var query Query
	for {
		collectionCountriesHistoryMutex.RLock()
		collection := mongoCountriesHistoryCollection
		collectionCountriesHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionCountriesHistoryMutex.RLock()
	collection := mongoCountriesHistoryCollection
	collectionCountriesHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelCountriesHistory) Rollback(transactionId string) error {
	var rows []CountryHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity Country
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
			dbServices.CollectionCache{}.Remove("Country", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelCountriesHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection CountriesHistory:")
	for key, value := range dbServices.GetDBIndexes(CountryHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionCountriesHistoryMutex.RLock()
		collection := mongoCountriesHistoryCollection
		collectionCountriesHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for CountryHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for CountryHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelCountriesHistory) New() *CountryHistoryRecord {
	return &CountryHistoryRecord{}
}

func (obj *CountryHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj CountryHistoryRecord
	row := modelCountriesHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *CountryHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionCountriesHistoryMutex.RLock()
	changeInfo, err := mongoCountriesHistoryCollection.UpsertId(objectId, &self)
	collectionCountriesHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for CountryHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *CountryHistoryRecord) Reflect() []Field {
	return Reflect(CountryHistoryRecord{})
}

func (self *CountryHistoryRecord) Delete() error {
	collectionCountriesHistoryMutex.RLock()
	collection := mongoCountriesHistoryCollection
	collectionCountriesHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *CountryHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *CountryHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *CountryHistoryRecord) GetType() int {
	return self.Type
}

func (self *CountryHistoryRecord) GetData() string {
	return self.Data
}

func (self *CountryHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *CountryHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *CountryHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *CountryHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *CountryHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *CountryHistoryRecord) GetId() string {
	return self.Id.Hex()
}
