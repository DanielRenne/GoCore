package model

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var HistCollection modelHistCollection

type modelHistCollection struct{}

var mongoHistCollectionCollection *mgo.Collection

func init() {
	go func() {

		for {
			if (dbServices.MongoDB != nil && !serverSettings.WebConfig.HasDbAuth) || (serverSettings.WebConfig.HasDbAuth && dbServices.MongoDBAuth != nil) {
				initHistCollection()
				return
			}
			time.Sleep(time.Millisecond * 5)
		}
	}()
}

func initHistCollection() {
	log.Println("Building Indexes for MongoDB collection HistCollection:")
	fmt.Sprint(serverSettings.WebConfig.HasDbAuth)
	//CollectionVariable
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoHistCollectionCollection.Create(&ci)
	HistCollection.Index()
}

type HistEntity struct {
	Id         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId        string        `json:"tId" dbIndex:"index" bson:"tId"`
	ObjId      string        `json:"objId" dbIndex:"index" bson:"objId"`
	Data       string        `json:"data" bson:"data"`
	Type       int           `json:"type" bson:"type"`
	CreateDate time.Time     `json:"createDate" dbIndex:"index" bson:"createDate"`
}

func (obj modelHistCollection) Query() *Query {
	var query Query
	for {
		if mongoHistCollectionCollection != nil {
			break
		}
		time.Sleep(time.Millisecond * 2)
	}
	query.collection = mongoHistCollectionCollection
	return &query
}

func (self *modelHistCollection) Rollback(transactionId string) error {
	var rows []HistEntity
	err := self.Query().Filter(Criteria("tId", transactionId)).All(&rows)

	if err != nil {
		return err
	}

	for _, entity := range rows {

		originalEntityData, err := decodeBase64(entity.Data)

		if err != nil {
			return err
		}

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
	for key, value := range dbServices.GetDBIndexes(HistEntity{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}

		err := mongoHistCollectionCollection.EnsureIndex(index)
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

func (self *HistEntity) Save() error {
	if mongoHistCollectionCollection == nil {
		initHistCollection()
	}
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	changeInfo, err := mongoHistCollectionCollection.UpsertId(objectId, &self)
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
	return mongoHistCollectionCollection.Remove(self)
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

func (self *HistEntity) Unmarshal(data string) error {

	err := json.Unmarshal([]byte(data), &self)
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
