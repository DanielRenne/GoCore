package model

import (
	"core/dbServices"
	"encoding/json"
	"github.com/asdine/storm"
)

type Cars struct{}

type Car struct {
	Id    int64  `json:"id" storm:"id"`
	Color string `json:"color"`
}

func (obj *Cars) Single(field string, value interface{}) (retObj Car, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj *Cars) Search(field string, value interface{}) (retObj []Car, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []Car{}
	}
	return
}

func (obj *Cars) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []Car, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	return
}

func (obj *Cars) All() (retObj []Car, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []Car{}
	}
	return
}

func (obj *Cars) AllAdvanced(limit int, skip int) (retObj []Car, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	return
}

func (obj *Cars) AllByIndex(index string) (retObj []Car, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []Car{}
	}
	return
}

func (obj *Cars) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Car, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	return
}

func (obj *Cars) Range(min, max, field string) (retObj []Car, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []Car{}
	}
	return
}

func (obj *Cars) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Car, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Car{}
		}
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
	return dbServices.BoltDB.Remove(obj)
}

func (obj *Car) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Car) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}
