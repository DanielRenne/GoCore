package model

import (
	"core/dbServices"
)

type Bucket struct {
	Name string
}

func (obj *Bucket) SetKeyValue(key interface{}, value interface{}) error {
	return dbServices.BoltDB.Set(obj.Name, key, value)
}

func (obj *Bucket) GetKeyValue(key interface{}, value interface{}) error {
	return dbServices.BoltDB.Get(obj.Name, key, value)
}

func (obj *Bucket) DeleteKey(key interface{}) error {
	return dbServices.BoltDB.Delete(obj.Name, key)
}
