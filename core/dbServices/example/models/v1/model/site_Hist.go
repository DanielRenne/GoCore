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

var SitesHistory modelSitesHistory

type modelSitesHistory struct{}

var mongoSitesHistoryCollection *mgo.Collection
var collectionSitesHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&SitesHistory)
	collectionSitesHistoryMutex = &sync.RWMutex{}
}

type SiteHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelSitesHistory) SetCollection(mdb *mgo.Database) {
	collectionSitesHistoryMutex.Lock()
	mongoSitesHistoryCollection = mdb.C("SitesHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoSitesHistoryCollection.Create(&ci)
	collectionSitesHistoryMutex.Unlock()
}

func (obj modelSitesHistory) Query() *Query {
	var query Query
	for {
		collectionSitesHistoryMutex.RLock()
		collection := mongoSitesHistoryCollection
		collectionSitesHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionSitesHistoryMutex.RLock()
	collection := mongoSitesHistoryCollection
	collectionSitesHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelSitesHistory) Rollback(transactionId string) error {
	var rows []SiteHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity Site
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
			dbServices.CollectionCache{}.Remove("Site", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelSitesHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection SitesHistory:")
	for key, value := range dbServices.GetDBIndexes(SiteHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionSitesHistoryMutex.RLock()
		collection := mongoSitesHistoryCollection
		collectionSitesHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for SiteHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for SiteHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelSitesHistory) New() *SiteHistoryRecord {
	return &SiteHistoryRecord{}
}

func (obj *SiteHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj SiteHistoryRecord
	row := modelSitesHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *SiteHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionSitesHistoryMutex.RLock()
	changeInfo, err := mongoSitesHistoryCollection.UpsertId(objectId, &self)
	collectionSitesHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for SiteHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *SiteHistoryRecord) Reflect() []Field {
	return Reflect(SiteHistoryRecord{})
}

func (self *SiteHistoryRecord) Delete() error {
	collectionSitesHistoryMutex.RLock()
	collection := mongoSitesHistoryCollection
	collectionSitesHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *SiteHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *SiteHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *SiteHistoryRecord) GetType() int {
	return self.Type
}

func (self *SiteHistoryRecord) GetData() string {
	return self.Data
}

func (self *SiteHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *SiteHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *SiteHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *SiteHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *SiteHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *SiteHistoryRecord) GetId() string {
	return self.Id.Hex()
}
