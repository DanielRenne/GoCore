package model

import (
	"core/dbServices"
	"encoding/json"
	"github.com/asdine/storm"
)

type Cars struct{}

type Car struct {
	Id    int    `json:"Id" storm:"id"`
	Color string `json:"color"`
}

func (obj *Cars) Single(field string, value string) (retObj Car) {
	dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj *Cars) Search(field string, value string) (retObj []Car) {
	dbServices.BoltDB.Find(field, value, &retObj)
	return
}

func (obj *Cars) SearchAdvanced(field string, value string, limit int, skip int) (retObj []Car) {
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

func (obj *Cars) All() (retObj []Car) {
	dbServices.BoltDB.All(&retObj)
	return
}

func (obj *Cars) AllAdvanced(limit int, skip int) (retObj []Car) {
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

func (obj *Cars) AllByIndex(index string) (retObj []Car) {
	dbServices.BoltDB.AllByIndex(index, &retObj)
	return
}

func (obj *Cars) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Car) {
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

func (obj *Cars) Range(min, max, field string) (retObj []Car) {
	dbServices.BoltDB.Range(field, min, max, &retObj)
	return
}

func (obj *Cars) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Car) {
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

func (obj *Cars) Index() error {
	return dbServices.BoltDB.Init(&Car{})
}

func (obj *Cars) RunTransaction(objects []Car) error {

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

func (obj *Cars) New() *Car {
	return &Car{}
}

func (obj *Car) Save() error {
	return dbServices.BoltDB.Save(obj)
}

func (obj *Car) Delete() error {
	return dbServices.BoltDB.Remove(&obj)
}

func (obj *Car) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Car) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}
