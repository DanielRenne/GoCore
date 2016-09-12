package model

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ERROR_INVALID_ID_VALUE = "Invalid Id value for query."
)

type queryError func() error

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

func Criteria(key string, value interface{}) map[string]interface{} {
	criteria := make(map[string]interface{}, 1)
	criteria[key] = value
	return criteria
}

func (self *Query) ById(val interface{}, x interface{}) error {

	objId, err := self.getIdHex(val)
	if err != nil {
		return err
	}

	q := self.collection.FindId(objId)

	err = q.One(x)

	if err != nil {
		callback := func() error {
			return q.One(x)
		}

		return self.handleQueryError(err, callback)
	}

	return err

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

func (self *Query) Exclude(criteria map[string]interface{}) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	self.inNot(criteria, "$nin")

	return self

}

func (self *Query) In(criteria map[string]interface{}) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	self.inNot(criteria, "$in")

	return self

}

func (self *Query) inNot(criteria map[string]interface{}, queryType string) {
	for key, value := range criteria {

		if key == "Id" {
			var ids []bson.ObjectId

			k := reflect.TypeOf(value).Kind()
			if k == reflect.Slice || k == reflect.Array {
				values := reflect.ValueOf(value)

				for i := 0; i < values.Len(); i++ {
					val := values.Index(i).Interface()
					objId, err := self.getIdHex(val)
					if err != nil {
						self.e = err
						continue
					}

					ids = append(ids, objId)
				}
			} else {
				objId, err := self.getIdHex(value)
				if err != nil {
					self.e = err
					continue
				}

				ids = append(ids, objId)
			}

			self.m["_id"] = bson.M{queryType: ids}

		} else {
			var valuesToQuery []interface{}

			k := reflect.TypeOf(value).Kind()
			if k == reflect.Slice || k == reflect.Array {
				values := reflect.ValueOf(value)
				for i := 0; i < values.Len(); i++ {
					val := values.Index(i).Interface()
					valuesToQuery = append(valuesToQuery, val)
				}

			} else {
				valuesToQuery = append(valuesToQuery, value)
			}

			self.m[key] = bson.M{queryType: valuesToQuery}
		}

	}
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

	err := q.All(x)

	if err != nil {

		callback := func() error {
			return q.All(x)
		}

		return self.handleQueryError(err, callback)
	}
	return err
}

func (self *Query) One(x interface{}) error {

	if self.e != nil {
		return self.e
	}

	q := self.generateQuery()

	err := q.One(x)

	if err != nil {

		callback := func() error {
			return q.One(x)
		}

		return self.handleQueryError(err, callback)
	}
	return err
}

func (self *Query) Count(x interface{}) (int, error) {

	if self.e != nil {
		return 0, self.e
	}

	q := self.generateQuery()

	count, err := q.Count()

	if err != nil {

		callback := func() error {
			count, err = q.Count()
			return err
		}

		return count, self.handleQueryError(err, callback)
	}
	return count, err
}

func (self *Query) Distinct(key string, x interface{}) error {

	if self.e != nil {
		return self.e
	}

	q := self.generateQuery()

	err := q.Distinct(key, x)

	if err != nil {

		callback := func() error {
			return q.Distinct(key, x)
		}

		return self.handleQueryError(err, callback)
	}
	return err

}

func (self *Query) handleQueryError(err error, callback queryError) error {

	if self.isDBConnectionError(err) {

		for i := 0; i < 2; i++ {

			log.Println("Attempting to Refresh Mongo Session")
			dbServices.MongoSession.Refresh()

			err = callback()
			if !self.isDBConnectionError(err) {
				return err
			}

			time.Sleep(200 * time.Millisecond)
		}
	}

	return err

}

func (self *Query) isDBConnectionError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "Closed explicitly") || strings.Contains(err.Error(), "read: operation timed out") || strings.Contains(err.Error(), "read: connection reset by peer") {
		return true
	}
	return false
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
	} else if myIdType.Name() == "ObjectId" && myIdType.Kind() == reflect.String {
		return bson.ObjectIdHex(myIdInstanceHex.Call([]reflect.Value{})[0].String()), nil
	} else {
		return bson.NewObjectId(), errors.New(ERROR_INVALID_ID_VALUE)
	}
}
