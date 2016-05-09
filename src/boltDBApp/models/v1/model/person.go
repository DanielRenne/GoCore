package model

import (
	"core/dbServices"
	"encoding/json"
	"github.com/asdine/storm"
)

type Persons struct{}

type Person struct {
	Id      int         `json:"Id" storm:"id"`
	Worth   float64     `json:"worth"`
	First   string      `json:"first" storm:"index"`
	IsCool  bool        `json:"isCool"`
	Blob    []byte      `json:"blob"`
	Hand    HandDetails `json:"hand"`
	Field7  []int       `json:"field7"`
	Field8  `json:"field8"`
	Field9  []string     `json:"field9"`
	Field10 []bool       `json:"field10"`
	Field12 []Field12Sub `json:"field12"`
}

type HandDetails struct {
	FingerCount int `json:"fingerCount"`
}

type Field12Sub struct {
	SubField int `json:"subField"`
}

func (obj *Persons) Single(field string, value string) (retObj Person) {
	dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj *Persons) Search(field string, value string) (retObj []Person) {
	dbServices.BoltDB.Find(field, value, &retObj)
	return
}

func (obj *Persons) SearchAdvanced(field string, value string, limit int, skip int) (retObj []Person) {
	if limit == 0 && skip == 0 {
		dbServices.BoltDB.Find(field, value, &retObj)
		return
	}
	if limit > 0 && skip > 0 {
		dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		return
	}
	if limit > 0 {
		dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		return
	}
	if skip > 0 {
		dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		return
	}
	return
}

func (obj *Persons) All() (retObj []Person) {
	dbServices.BoltDB.All(&retObj)
	return
}

func (obj *Persons) AllAdvanced(limit int, skip int) (retObj []Person) {
	if limit == 0 && skip == 0 {
		dbServices.BoltDB.All(&retObj)
		return
	}
	if limit > 0 && skip > 0 {
		dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		return
	}
	if limit > 0 {
		dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		return
	}
	if skip > 0 {
		dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		return
	}
	return
}

func (obj *Persons) AllByIndex(index string) (retObj []Person) {
	dbServices.BoltDB.AllByIndex(index, &retObj)
	return
}

func (obj *Persons) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Person) {
	if limit == 0 && skip == 0 {
		dbServices.BoltDB.AllByIndex(index, &retObj)
		return
	}
	if limit > 0 && skip > 0 {
		dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		return
	}
	if limit > 0 {
		dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		return
	}
	if skip > 0 {
		dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		return
	}
	return
}

func (obj *Persons) Range(min, max, field string) (retObj []Person) {
	dbServices.BoltDB.Range(field, min, max, &retObj)
	return
}

func (obj *Persons) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Person) {
	if limit == 0 && skip == 0 {
		dbServices.BoltDB.Range(field, min, max, &retObj)
		return
	}
	if limit > 0 && skip > 0 {
		dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		return
	}
	if limit > 0 {
		dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		return
	}
	if skip > 0 {
		dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		return
	}
	return
}

func (obj *Persons) Index() error {
	return dbServices.BoltDB.Init(&Person{})
}

func (obj *Persons) RunTransaction(objects []Person) error {

	tx, err := dbServices.BoltDB.Begin(true)

	for _, object := range objects {
		err = tx.Save(&object)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (obj *Persons) New() *Person {
	return &Person{}
}

func (obj *Person) Save() error {
	return dbServices.BoltDB.Save(obj)
}

func (obj *Person) Delete() error {
	return dbServices.BoltDB.Remove(&obj)
}

func (obj *Person) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Person) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}
