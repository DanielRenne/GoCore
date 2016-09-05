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

func (self *Query) ById(id interface{}, x interface{}) error {

	myIdType := reflect.TypeOf(id)
	myIdInstance := reflect.ValueOf(id)
	myIdInstanceHex := myIdInstance.MethodByName("Hex")
	if myIdType.Name() != "ObjectId" && myIdType.Kind() == reflect.String && id.(string) != "" {
		self.q = self.collection.FindId(bson.ObjectIdHex(id.(string)))
	} else if myIdType.Name() == "ObjectId" && myIdType.Kind() == reflect.String && id.(string) != "" {
		self.q = self.collection.FindId(bson.ObjectIdHex(myIdInstanceHex.Call([]reflect.Value{})[0].String()))
	} else {
		return errors.New(ERROR_INVALID_ID_VALUE)
	}

	return self.q.One(x)

}

func (self *Query) Filter(criteria map[string]interface{}) *Query {

	idVal, hasId := criteria["Id"]
	if hasId {
		myIdType := reflect.TypeOf(idVal)
		myIdInstance := reflect.ValueOf(idVal)
		myIdInstanceHex := myIdInstance.MethodByName("Hex")
		if myIdType.Name() != "ObjectId" && myIdType.Kind() == reflect.String {
			self.q = self.collection.FindId(bson.ObjectIdHex(idVal.(string)))
		} else if myIdType.Name() == "ObjectId" && myIdType.Kind() == reflect.String {
			self.q = self.collection.FindId(bson.ObjectIdHex(myIdInstanceHex.Call([]reflect.Value{})[0].String()))
		} else {
			self.e = errors.New(ERROR_INVALID_ID_VALUE)
		}
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
	q := self.generateQuery()
	return q.All(x)
}

func (self *Query) One(x interface{}) error {
	q := self.generateQuery()
	return q.One(x)
}

func (self *Query) Count(x interface{}) (int, error) {
	q := self.generateQuery()
	return q.Count()
}

func (self *Query) Distinct(key string, x interface{}) error {
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
