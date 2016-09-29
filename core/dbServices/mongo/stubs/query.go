package model

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/extensions"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ERROR_INVALID_ID_VALUE = "Invalid Id value for query."
	JOIN_ALL               = "All"
)

type queryError func() error

type Range struct {
	Min interface{}
	Max interface{}
}

type join struct {
	collectionName   string
	joinFieldRefName string
	joinFieldName    string
	joinSchemaName   string
	joinSpecified    string
}

type Query struct {
	q          *mgo.Query
	m          bson.M
	limit      int
	skip       int
	sort       []string
	collection *mgo.Collection
	e          error
	joins      []string
	allJoins   bool
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
		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.One(x)
			if err != nil {
				return err
			}
			return self.processJoins(x)
		}

		return self.handleQueryError(err, callback)
	}

	return self.processJoins(x)

}

func (self *Query) Join(criteria string) *Query {
	if self.allJoins {
		return self
	}
	if criteria == JOIN_ALL {
		self.allJoins = true
		return self
	}
	self.joins = append(self.joins, criteria)
	return self
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
		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.All(x)
			if err != nil {
				return err
			}
			return self.processJoins(x)
		}

		return self.handleQueryError(err, callback)
	}
	return self.processJoins(x)
}

func (self *Query) One(x interface{}) error {

	if self.e != nil {
		return self.e
	}

	q := self.generateQuery()

	err := q.One(x)

	if err != nil {

		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.One(x)
			if err != nil {
				return err
			}
			cnt, _ := q.Count()
			if cnt != 1 {
				return errors.New("Did not return exactly one row.  Returned " + extensions.IntToString(cnt))
			}
			return self.processJoins(x)
		}

		return self.handleQueryError(err, callback)
	}
	cnt, _ := q.Count()
	if cnt != 1 {
		return errors.New("Did not return exactly one row.  Returned " + extensions.IntToString(cnt))
	}

	return self.processJoins(x)
}

func (self *Query) Count() (int, error) {

	if self.e != nil {
		return 0, self.e
	}

	q := self.generateQuery()

	count, err := q.Count()

	if err != nil {

		// This callback is used for if the Ethernet port is unplugged
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

		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.Distinct(key, x)
			if err != nil {
				return err
			}
			return self.processJoins(x)

		}

		return self.handleQueryError(err, callback)
	}

	return self.processJoins(x)

}

func (self *Query) processJoins(x interface{}) (err error) {

	if self.allJoins || len(self.joins) > 0 {

		//Check if x is a single struct or an Array
		_, isArray := valueType(x)

		if isArray {
			source := reflect.ValueOf(x).Elem()

			var joins []join
			if source.Len() > 0 {
				joins = self.getJoins(source.Index(0))
			}

			if len(joins) == 0 {
				return
			}

			//Advanced way is to get the count and chunk.  For now we will iterate.
			for i := 0; i < source.Len(); i++ {
				s := source.Index(i)
				for _, j := range joins {
					id := reflect.ValueOf(s.FieldByName(j.joinFieldRefName).Interface()).String()
					joinsField := s.FieldByName("Joins")
					setField := joinsField.FieldByName(j.joinFieldName)

					err = joinField(j.joinSchemaName, j.collectionName, id, setField, j.joinSpecified)
					if err != nil {
						return
					}
				}
			}

		} else {
			source := reflect.ValueOf(x).Elem()

			var joins []join
			joins = self.getJoins(source)

			if len(joins) == 0 {
				return
			}

			s := source
			for _, j := range joins {
				id := reflect.ValueOf(s.FieldByName(j.joinFieldRefName).Interface()).String()
				joinsField := s.FieldByName("Joins")
				setField := joinsField.FieldByName(j.joinFieldName)

				err = joinField(j.joinSchemaName, j.collectionName, id, setField, j.joinSpecified)
				if err != nil {
					return
				}
			}
		}
		return nil
	}
	return nil
}

func (self *Query) getJoins(x reflect.Value) (joins []join) {

	joinsField := x.FieldByName("Joins")

	if joinsField.Kind() != reflect.Struct {
		return
	}

	if self.allJoins {
		for i := 0; i < joinsField.NumField(); i++ {

			typeField := joinsField.Type().Field(i)
			name := typeField.Name
			tagValue := typeField.Tag.Get("join")
			splitValue := strings.Split(tagValue, ",")
			var j join
			j.collectionName = splitValue[0]
			j.joinSchemaName = splitValue[1]
			j.joinFieldRefName = splitValue[2]
			j.joinFieldName = name
			j.joinSpecified = JOIN_ALL
			joins = append(joins, j)
		}
	} else {
		for _, name := range self.joins {

			fields := strings.Split(name, ".")
			fieldName := fields[0]

			typeField, ok := joinsField.Type().FieldByName(fieldName)
			if ok == false {
				continue
			}

			tagValue := typeField.Tag.Get("join")
			splitValue := strings.Split(tagValue, ",")
			var j join
			j.collectionName = splitValue[0]
			j.joinSchemaName = splitValue[1]
			j.joinFieldRefName = splitValue[2]
			j.joinFieldName = fieldName
			j.joinSpecified = strings.Replace(name, fieldName+".", "", 1)
			joins = append(joins, j)
		}
	}

	return
}

func valueType(m interface{}) (t reflect.Type, isArray bool) {
	t = reflect.Indirect(reflect.ValueOf(m)).Type()
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		isArray = true
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return

	}
	return
}

// func (self *Query) processJoinReflection

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
