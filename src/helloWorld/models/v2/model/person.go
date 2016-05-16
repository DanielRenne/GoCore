package model

import (
	"core/dbServices"
	"encoding/json"
	"github.com/asdine/storm"
)

type Persons struct{}

type Person struct {
	Id      int64        `json:"id" storm:"id"`
	Worth   float64      `json:"worth"`
	First   string       `json:"first" storm:"index"`
	IsCool  bool         `json:"isCool"`
	Blob    []byte       `json:"blob"`
	Hand    HandDetails  `json:"hand"`
	Field7  []int        `json:"field7"`
	Field8  []float64    `json:"field8"`
	Field9  []string     `json:"field9"`
	Field10 []bool       `json:"field10"`
	Field12 []Field12Sub `json:"field12"`
}

type HandDetails struct {
	FingerCount int64 `json:"fingerCount"`
}

type Field12Sub struct {
	SubField int64 `json:"subField"`
}

func (obj *Persons) Single(field string, value interface{}) (retObj Person, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj *Persons) Search(field string, value interface{}) (retObj []Person, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []Person{}
	}
	return
}

func (obj *Persons) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []Person, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	return
}

func (obj *Persons) All() (retObj []Person, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []Person{}
	}
	return
}

func (obj *Persons) AllAdvanced(limit int, skip int) (retObj []Person, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	return
}

func (obj *Persons) AllByIndex(index string) (retObj []Person, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []Person{}
	}
	return
}

func (obj *Persons) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Person, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	return
}

func (obj *Persons) Range(min, max, field string) (retObj []Person, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []Person{}
	}
	return
}

func (obj *Persons) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Person, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Person{}
		}
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
	return dbServices.BoltDB.Remove(obj)
}

func (obj *Person) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Person) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}
