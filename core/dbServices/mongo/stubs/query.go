package model

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

const (
	ERROR_INVALID_ID_VALUE = "Invalid Id value for query."
)

type Range struct {
	Min interface{}
	Max interface{}
}

type Query struct {
	q          *mgo.Query
	m          bson.M
	limit      int
	skip       int
	sort       []string
	collection *mgo.Collection
	e          error
}

func (self *Query) ById(val interface{}, x interface{}) error {

	objId, err := self.getIdHex(val)
	if err != nil {
		return err
	}

	self.q = self.collection.FindId(objId)
	return self.q.One(x)

}

func (self *Query) Filter(criteria map[string]interface{}) *Query {

	val, hasId := criteria["Id"]
	if hasId {

		objId, err := self.getIdHex(val)
		if err != nil {
			self.e = err
			return self
		}

		self.q = self.collection.FindId(objId)

		return self
	} else {

		if self.m == nil {
			self.m = make(bson.M)
		}

		for key, val := range criteria {
			if key != "" {
				self.m[key] = val
			}
		}

	}
	return self
}

func (self *Query) In(field string, values ...string) *Query {

	if field == "Id" {
		var ids []bson.ObjectId

		for _, val := range values {

			objId, err := self.getIdHex(val)
			if err != nil {
				self.e = err
				continue
			}

			ids = append(ids, objId)
		}

		self.q = self.collection.Find(bson.M{"_id": bson.M{"$in": ids}})

	} else {

		self.q = self.collection.Find(bson.M{field: bson.M{"$in": values}})
	}
	return self

}

func (self *Query) Range(criteria map[string]Range) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	for key, val := range criteria {
		if key != "" {
			self.m[key] = bson.M{"$gte": val.Min, "$lte": val.Max}
		}
	}

	return self
}

func (self *Query) Sort(fields ...string) *Query {
	self.sort = fields
	return self
}

func (self *Query) Limit(val int) *Query {
	self.limit = val
	return self
}

func (self *Query) Skip(val int) *Query {
	self.skip = val
	return self
}

func (self *Query) All(x interface{}) error {

	if self.e != nil {
		return self.e
	}

	q := self.generateQuery()
	return q.All(x)
}

func (self *Query) One(x interface{}) error {

	if self.e != nil {
		return self.e
	}

	q := self.generateQuery()
	return q.One(x)
}

func (self *Query) Count(x interface{}) (int, error) {

	if self.e != nil {
		return 0, self.e
	}

	q := self.generateQuery()
	return q.Count()
}

func (self *Query) Distinct(key string, x interface{}) error {

	if self.e != nil {
		return self.e
	}

	q := self.generateQuery()
	return q.Distinct(key, x)
}

func (self *Query) generateQuery() *mgo.Query {

	q := self.collection.Find(bson.M{})

	if self.q != nil {
		q = self.q
	}

	if self.m != nil {
		q = self.collection.Find(self.m)
	}

	if self.limit > 0 {
		q = q.Limit(self.limit)
	}

	if self.skip > 0 {
		q = q.Skip(self.skip)
	}

	if len(self.sort) > 0 {
		q = q.Sort(self.sort...)
	}

	return q
}

func (self *Query) getIdHex(val interface{}) (bson.ObjectId, error) {

	myIdType := reflect.TypeOf(val)
	myIdInstance := reflect.ValueOf(val)
	myIdInstanceHex := myIdInstance.MethodByName("Hex")

	if myIdType.Name() != "ObjectId" && myIdType.Kind() == reflect.String && val.(string) != "" {
		return bson.ObjectIdHex(val.(string)), nil
	} else if myIdType.Name() == "ObjectId" && myIdType.Kind() == reflect.String && val.(string) != "" {
		return bson.ObjectIdHex(myIdInstanceHex.Call([]reflect.Value{})[0].String()), nil
	} else {
		return bson.NewObjectId(), errors.New(ERROR_INVALID_ID_VALUE)
	}
}
