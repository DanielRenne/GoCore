package model

import (
	"encoding/json"
	"errors"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type HistCollection struct{}

var mongoHistCollectionCollection *mgo.Collection

func init() {
	go func() {

		for {
			if dbServices.MongoDB != nil {
				log.Println("Building Indexes for MongoDB collection HistCollection:")
				mongoHistCollectionCollection = dbServices.MongoDB.C("HistCollection")
				ci := mgo.CollectionInfo{ForceIdIndex: true}
				mongoHistCollectionCollection.Create(&ci)
				var obj HistCollection
				obj.Index()
				return
			}
			<-dbServices.WaitForDatabase()
		}
	}()
}

type HistEntity struct {
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	TId  string        `json:"tId" dbIndex:"index" bson:"tId"`
	Data string        `json:"data" bson:"data"`
	Type int           `json:"type" bson:"type"`
}

func (self *HistCollection) Single(field string, value interface{}) (retObj HistEntity, e error) {
	if field == "id" {
		query := mongoHistCollectionCollection.FindId(bson.ObjectIdHex(value.(string)))
		e = query.One(&retObj)
		return
	}
	m := make(bson.M)
	m[field] = value
	query := mongoHistCollectionCollection.Find(m)
	e = query.One(&retObj)
	return
}

func (self *HistCollection) Rollback(transactionId string) error {
	rows, err := self.Search("tId", transactionId)

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

func (obj *HistCollection) Search(field string, value interface{}) (retObj []HistEntity, e error) {
	var query *mgo.Query
	if field == "id" {
		query = mongoHistCollectionCollection.FindId(bson.ObjectIdHex(value.(string)))
	} else {
		m := make(bson.M)
		m[field] = value
		query = mongoHistCollectionCollection.Find(m)
	}

	e = query.All(&retObj)
	return
}

func (obj *HistCollection) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []HistEntity, e error) {
	var query *mgo.Query
	if field == "id" {
		query = mongoHistCollectionCollection.FindId(bson.ObjectIdHex(value.(string)))
	} else {
		m := make(bson.M)
		m[field] = value
		query = mongoHistCollectionCollection.Find(m)
	}

	if limit == 0 && skip == 0 {
		e = query.All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = query.Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 {
		e = query.Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if skip > 0 {
		e = query.Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	return
}

func (obj *HistCollection) All() (retObj []HistEntity, e error) {
	e = mongoHistCollectionCollection.Find(bson.M{}).All(&retObj)
	if len(retObj) == 0 {
		retObj = []HistEntity{}
	}
	return
}

func (obj *HistCollection) AllAdvanced(limit int, skip int) (retObj []HistEntity, e error) {
	if limit == 0 && skip == 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if skip > 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	return
}

func (obj *HistCollection) AllByIndex(index string) (retObj []HistEntity, e error) {
	e = mongoHistCollectionCollection.Find(bson.M{}).Sort(index).All(&retObj)
	if len(retObj) == 0 {
		retObj = []HistEntity{}
	}
	return
}

func (obj *HistCollection) AllByIndexAdvanced(index string, limit int, skip int) (retObj []HistEntity, e error) {
	if limit == 0 && skip == 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Sort(index).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Sort(index).Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Sort(index).Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if skip > 0 {
		e = mongoHistCollectionCollection.Find(bson.M{}).Sort(index).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	return
}

func (obj *HistCollection) Range(min, max, field string) (retObj []HistEntity, e error) {
	var query *mgo.Query
	m := make(bson.M)
	m[field] = bson.M{"$gte": min, "$lte": max}
	query = mongoHistCollectionCollection.Find(m)
	e = query.All(&retObj)
	return
}

func (obj *HistCollection) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []HistEntity, e error) {
	var query *mgo.Query
	m := make(bson.M)
	m[field] = bson.M{"$gte": min, "$lte": max}
	query = mongoHistCollectionCollection.Find(m)
	if limit == 0 && skip == 0 {
		e = query.All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = query.Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if limit > 0 {
		e = query.Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	if skip > 0 {
		e = query.Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []HistEntity{}
		}
		return
	}
	return
}

func (obj *HistCollection) Index() error {
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

func (obj *HistCollection) New() *HistEntity {
	return &HistEntity{}
}

func (self *HistEntity) Save() error {
	if mongoHistCollectionCollection == nil {
		return errors.New("Collection HistCollection not initialized")
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

func (self *HistEntity) Delete() error {
	return mongoHistCollectionCollection.Remove(self)
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
