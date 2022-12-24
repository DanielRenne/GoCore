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

var AccountsHistory modelAccountsHistory

type modelAccountsHistory struct{}

var mongoAccountsHistoryCollection *mgo.Collection
var collectionAccountsHistoryMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&AccountsHistory)
	collectionAccountsHistoryMutex = &sync.RWMutex{}
}

type AccountHistoryRecord struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelAccountsHistory) SetCollection(mdb *mgo.Database) {
	collectionAccountsHistoryMutex.Lock()
	mongoAccountsHistoryCollection = mdb.C("AccountsHistory")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoAccountsHistoryCollection.Create(&ci)
	collectionAccountsHistoryMutex.Unlock()
}

func (obj modelAccountsHistory) Query() *Query {
	var query Query
	for {
		collectionAccountsHistoryMutex.RLock()
		collection := mongoAccountsHistoryCollection
		collectionAccountsHistoryMutex.RUnlock()
		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	collectionAccountsHistoryMutex.RLock()
	collection := mongoAccountsHistoryCollection
	collectionAccountsHistoryMutex.RUnlock()

	query.collection = collection
	return &query
}

func (self *modelAccountsHistory) Rollback(transactionId string) error {
	var rows []AccountHistoryRecord
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData := []byte(entity.Data)

		var rollbackEntity Account
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
			dbServices.CollectionCache{}.Remove("Account", rollbackEntity.Id.Hex())
		}

	}
	return err
}

func (obj *modelAccountsHistory) Index() error {
	log.Println("Building Indexes for MongoDB collection AccountsHistory:")
	for key, value := range dbServices.GetDBIndexes(AccountHistoryRecord{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}
		collectionAccountsHistoryMutex.RLock()
		collection := mongoAccountsHistoryCollection
		collectionAccountsHistoryMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for AccountHistoryRecord." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for AccountHistoryRecord." + key)
		}
	}
	return nil
}

func (obj *modelAccountsHistory) New() *AccountHistoryRecord {
	return &AccountHistoryRecord{}
}

func (obj *AccountHistoryRecord) DoesIdExist(objectID interface{}) bool {
	var retObj AccountHistoryRecord
	row := modelAccountsHistory{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (self *AccountHistoryRecord) Save() error {
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	collectionAccountsHistoryMutex.RLock()
	changeInfo, err := mongoAccountsHistoryCollection.UpsertId(objectId, &self)
	collectionAccountsHistoryMutex.RUnlock()
	if err != nil {
		log.Println("Failed to upsertId for AccountHistoryRecord:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (obj *AccountHistoryRecord) Reflect() []Field {
	return Reflect(AccountHistoryRecord{})
}

func (self *AccountHistoryRecord) Delete() error {
	collectionAccountsHistoryMutex.RLock()
	collection := mongoAccountsHistoryCollection
	collectionAccountsHistoryMutex.RUnlock()
	return collection.Remove(self)
}

func (self *AccountHistoryRecord) SaveWithTran(t *Transaction) error {
	return nil
}

func (self *AccountHistoryRecord) JoinFields(s string, q *Query, x int) error {
	return nil
}

func (self *AccountHistoryRecord) GetType() int {
	return self.Type
}

func (self *AccountHistoryRecord) GetData() string {
	return self.Data
}

func (self *AccountHistoryRecord) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *AccountHistoryRecord) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *AccountHistoryRecord) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *AccountHistoryRecord) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *AccountHistoryRecord) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (self *AccountHistoryRecord) GetId() string {
	return self.Id.Hex()
}
