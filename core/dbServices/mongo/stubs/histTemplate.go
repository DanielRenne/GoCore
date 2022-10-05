// Package mongoStubs - Internal Stub templates for mongo struct/ORM generation
package mongoStubs

var HistTemplate string

func init() {
	HistTemplate = `
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

var HistCollection modelHistCollection

type modelHistCollection struct{}

var mongoHistCollectionCollection *mgo.Collection
var collectionHistCollectionMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&HistCollection)
	collectionHistCollectionMutex = &sync.RWMutex{}
}

type HistEntity struct {
	Id         bson.ObjectId ` + "`" + `json:"id" bson:"_id,omitempty"` + "`" + `
	TId        string        ` + "`" + `json:"tId" dbIndex:"index" bson:"tId"` + "`" + `
	ObjId      string        ` + "`" + `json:"objId" dbIndex:"index" bson:"objId"` + "`" + `
	Data       string        ` + "`" + `json:"data" bson:"data"` + "`" + `
	Type       int           ` + "`" + `json:"type" bson:"type"` + "`" + `
	CreateDate time.Time     ` + "`" + `json:"createDate" dbIndex:"index" bson:"createDate"` + "`" + `
}

func (obj modelHistCollection) SetCollection(mdb *mgo.Database) {
	collectionHistCollectionMutex.Lock()
	mongoHistCollectionCollection = mdb.C("HistCollection")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoHistCollectionCollection.Create(&ci)
	collectionHistCollectionMutex.Unlock()
}

func (obj modelHistCollection) Query() *Query {
	var query Query
	for {
		collectionHistCollectionMutex.RLock()
		collection := mongoHistCollectionCollection
		collectionHistCollectionMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 2)
	}
	collectionHistCollectionMutex.RLock()
	collection := mongoHistCollectionCollection
	collectionHistCollectionMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelHistCollection) Rollback(transactionId string) error {
	var rows []HistEntity
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity OriginalEntity
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

	}
	return err
}

func (obj *modelHistCollection) Index() error {
	log.Println("Building Indexes for MongoDB collection HistCollection:")
	for key, value := range dbServices.GetDBIndexes(HistEntity{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionHistCollectionMutex.RLock()
		collection := mongoHistCollectionCollection
		collectionHistCollectionMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for HistEntity." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for HistEntity." + key)
		}
	}
	return nil
}

func (obj *modelHistCollection) New() *HistEntity {
	return &HistEntity{}
}

func (obj *HistEntity) DoesIdExist(objectID interface{}) bool {
	var retObj HistEntity
	row := modelHistCollection{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *HistEntity) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionHistCollectionMutex.RLock()
	changeInfo, err := mongoHistCollectionCollection.UpsertId(objectId, &self)
	collectionHistCollectionMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for HistEntity:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *HistEntity) Reflect() []Field {
	return Reflect(HistEntity{})
}

func (self *HistEntity) Delete() error {
	collectionHistCollectionMutex.RLock()
	collection := mongoHistCollectionCollection
	collectionHistCollectionMutex.RUnlock()
	return collection.Remove(self)
}

func (self *HistEntity) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *HistEntity) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *HistEntity) GetType() int {
	return self.Type
}

func (self *HistEntity) GetData() string {
	return self.Data
}

func (self *HistEntity) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *HistEntity) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *HistEntity) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *HistEntity) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *HistEntity) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *HistEntity) GetId() string {
	return self.Id.Hex()
}
`
}
