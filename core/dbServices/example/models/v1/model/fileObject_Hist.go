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

var FileObjectsHistory modelFileObjectsHistory

type modelFileObjectsHistory struct{}

var mongoFileObjectsHistoryCollection *mgo.Collection
var collectionFileObjectsHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&FileObjectsHistory)
	collectionFileObjectsHistoryMutex = &sync.RWMutex{}
}

type FileObjectHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelFileObjectsHistory) SetCollection(mdb *mgo.Database) {
	collectionFileObjectsHistoryMutex.Lock()
	mongoFileObjectsHistoryCollection = mdb.C("FileObjectsHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoFileObjectsHistoryCollection.Create(&ci)
	collectionFileObjectsHistoryMutex.Unlock()
}

func (obj modelFileObjectsHistory) Query() *Query {
	var query Query
	for {
		collectionFileObjectsHistoryMutex.RLock()
		collection := mongoFileObjectsHistoryCollection
		collectionFileObjectsHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionFileObjectsHistoryMutex.RLock()
	collection := mongoFileObjectsHistoryCollection
	collectionFileObjectsHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelFileObjectsHistory) Rollback(transactionId string) error {
	var rows []FileObjectHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity FileObject
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
			dbServices.CollectionCache{}.Remove("FileObject", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelFileObjectsHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection FileObjectsHistory:")
	for key, value := range dbServices.GetDBIndexes(FileObjectHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionFileObjectsHistoryMutex.RLock()
		collection := mongoFileObjectsHistoryCollection
		collectionFileObjectsHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for FileObjectHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for FileObjectHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelFileObjectsHistory) New() *FileObjectHistoryRecord {
	return &FileObjectHistoryRecord{}
}

func (obj *FileObjectHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj FileObjectHistoryRecord
	row := modelFileObjectsHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *FileObjectHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionFileObjectsHistoryMutex.RLock()
	changeInfo, err := mongoFileObjectsHistoryCollection.UpsertId(objectId, &self)
	collectionFileObjectsHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for FileObjectHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *FileObjectHistoryRecord) Reflect() []Field {
	return Reflect(FileObjectHistoryRecord{})
}

func (self *FileObjectHistoryRecord) Delete() error {
	collectionFileObjectsHistoryMutex.RLock()
	collection := mongoFileObjectsHistoryCollection
	collectionFileObjectsHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *FileObjectHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *FileObjectHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *FileObjectHistoryRecord) GetType() int {
	return self.Type
}

func (self *FileObjectHistoryRecord) GetData() string {
	return self.Data
}

func (self *FileObjectHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FileObjectHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *FileObjectHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *FileObjectHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *FileObjectHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *FileObjectHistoryRecord) GetId() string {
	return self.Id.Hex()
}
