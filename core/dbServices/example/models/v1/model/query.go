package model

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/DanielRenne/GoCore/core/serverSettings"
	stacktrace "github.com/go-errors/errors"

	"fmt"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/logger"
	dateformatter "github.com/altipla-consulting/i18n-dateformatter"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

type Min struct {
	Min interface{}
}

type Max struct {
	Max interface{}
}

type join struct {
	collectionName       string
	joinFieldRefName     string
	joinFieldName        string
	joinSchemaName       string
	joinSpecified        string
	joinType             string
	isMany               bool
	joinForeignFieldName string
	whiteListedFields    []string
	blackListedFields    []string
}

type joinType struct {
	Type string
}

type view struct {
	fieldName string
	ref       string
	format    string
}

type Query struct {
	q              *mgo.Query
	m              bson.M
	o              []bson.M
	ao             map[string][]map[string][]bson.M
	stopLog        bool
	limit          int
	skip           int
	sort           []string
	collection     *mgo.Collection
	entityName     string
	e              error
	joins          map[string]joinType
	format         DataFormat
	renderViews    bool
	whiteListed    []QueryFieldFilter
	blackListed    []QueryFieldFilter
	filteredFields bool
}

type QueryIterator struct {
	q    *Query
	iter *mgo.Iter
}

type QueryFieldFilter struct {
	CollectionName string
	Fields         []string
}

type DataFormat struct {
	Language      string `json:"Language"`
	DateFormat    string `json:"DateFormat"`
	LocalTimeZone string `json:"LocalTimeZone"`
}

func (obj *DataFormat) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (self *DataFormat) Parse(data string) (err error) {
	err = json.Unmarshal([]byte(data), &self)
	return
}

func Criteria(key string, value interface{}) map[string]interface{} {
	criteria := make(map[string]interface{}, 1)
	criteria[key] = value
	return criteria
}

func (self *QueryIterator) Next(x interface{}) (gotRecord bool, err error) {

	// This callback is used for if the Ethernet port is unplugged
	callback := func() (err error) {
		gotRecord = self.iter.Next(x)
		if gotRecord == false {
			err = self.iter.Close()
			return
		}
		err = self.q.processJoinsAndViews(x)
		return
	}

	gotRecord = self.iter.Next(x)
	if gotRecord == true {
		err = self.q.processJoinsAndViews(x)
		return
	} else {
		err = self.iter.Close()
		if err != nil {
			err = self.q.handleQueryError(err, callback)
		} else {
			return
		}

	}
	return
}

func (self *Query) ById(objectId interface{}, modelInstance interface{}) error {

	objId, err := self.getIdHex(objectId)
	if err != nil {
		return err
	}

	if !serverSettings.WebConfig.Application.LogQueries && serverSettings.WebConfig.Application.LogQueryTimes {
		defer func() {
			log.Println(logger.TimeTrackQuery(time.Now(), "q.ById("+objId.Hex()+")", self.collection, self.m, self.q))
		}()
	} else if serverSettings.WebConfig.Application.LogQueries {
		defer func() {
			log.Println(logger.TimeTrack(time.Now(), "q.ById("+objId.Hex()+")"))
		}()
	}

	if !self.stopLog && serverSettings.WebConfig.Application.LogQueries {
		self.LogQuery("q.ById(" + objId.Hex() + ")")
	}
	q := self.collection.FindId(objId)

	cacheData := true
	foundCache := false

	if self.filteredFields {
		cacheData = false
	} else {
		foundCache = dbServices.CollectionCache{}.Fetch(self.collection.Name, objId.Hex(), modelInstance)
	}

	if foundCache == false {
		err = q.One(modelInstance)
	}

	if err != nil {
		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.One(modelInstance)
			if err != nil {
				return err
			}
			if cacheData {

				if modelInstance != nil && reflect.TypeOf(modelInstance).Name() == self.entityName {
					dbServices.CollectionCache{}.Store(self.collection.Name, objId.Hex(), modelInstance)
				}
			}

			return self.processJoinsAndViews(modelInstance)
		}

		return self.handleQueryError(err, callback)
	}

	if foundCache == false && cacheData {
		if modelInstance != nil && reflect.TypeOf(modelInstance).Name() == self.entityName {
			dbServices.CollectionCache{}.Store(self.collection.Name, objId.Hex(), modelInstance)
		}
	}
	return self.processJoinsAndViews(modelInstance)
}

func (self *Query) GetCollection() *mgo.Collection {
	return self.collection
}

func (self *Query) Join(criteria string) *Query {

	if self.joins == nil {
		self.joins = make(map[string]joinType)
	}

	_, ok := self.joins[JOIN_ALL]
	if ok {
		return self
	}

	self.joins[criteria] = joinType{Type: "Inner"}

	return self
}

func (self *Query) LeftJoin(criteria string) *Query {

	if self.joins == nil {
		self.joins = make(map[string]joinType)
	}

	_, ok := self.joins[JOIN_ALL]
	if ok {
		return self
	}

	self.joins[criteria] = joinType{Type: "Left"}

	return self
}

func (self *Query) ToggleLogFlag(toggle bool) *Query {
	self.stopLog = toggle
	return self
}

func (self *Query) Or(criteria map[string]interface{}) *Query {

	val, hasId := criteria["Id"]
	if hasId {

		objId, err := self.getIdHex(val)
		if err != nil {
			self.e = err
			return self
		}

		if self.o == nil {
			self.o = make([]bson.M, 0)
		}

		self.o[0]["_id"] = objId
		return self
	} else {

		if self.o == nil {
			self.o = make([]bson.M, 0)
		}

		for key, val := range criteria {
			if key != "" {
				self.o = append(self.o, Q(key, self.CheckForObjectId(val)))
			}
		}

	}
	return self
}

func (self *Query) GetAndOr() map[string][]map[string][]bson.M {
	return self.ao
}

func (self *Query) InitAndOr() *Query {
	self.ao = nil
	self.ao = map[string][]map[string][]bson.M{}
	self.AddAndOr()
	return self
}

func (self *Query) AddAndOr() *Query {
	self.ao["$and"] = append(self.ao["$and"], map[string][]bson.M{
		"$or": make([]bson.M, 0),
	})
	if len(self.ao["$and"]) > 1 {
		// first index is reserved for search criteria so we always will fill in a stub conditionally when no search is passed.
		// Other results we should put in something just in case a filter is not passed so we dont have to litter the _id exists hack across all the code

		self.OrFilter(len(self.ao["$and"])-1, Q("_id", Q("$exists", true)))
	}

	return self
}

func (self *Query) AddAnd() *Query {
	self.ao["$and"] = append(self.ao["$and"], map[string][]bson.M{
		"$and": make([]bson.M, 0),
	})
	if len(self.ao["$and"]) > 1 {
		// first index is reserved for search criteria so we always will fill in a stub conditionally when no search is passed.
		// Other results we should put in something just in case a filter is not passed so we dont have to litter the _id exists hack across all the code

		self.AndFilter(len(self.ao["$and"])-1, Q("_id", Q("$exists", true)))
	}

	return self
}

func (self *Query) AddBlankAndOr() *Query {
	self.ao["$and"] = append(self.ao["$and"], map[string][]bson.M{
		"$or": make([]bson.M, 0),
	})
	return self
}

func (self *Query) AndFilter(index int, criteria map[string]interface{}) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$and"] = append(self.ao["$and"][index]["$and"], Q(key, val))
	}
	return self
}

func (self *Query) OrFilter(index int, criteria map[string]interface{}) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$or"] = append(self.ao["$and"][index]["$or"], Q(key, val))
	}
	return self
}

func (self *Query) AndRange(index int, criteria map[string]Range) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$and"] = append(self.ao["$and"][index]["$and"], bson.M{"$gte": val.Min, "$lte": val.Max})
	}

	return self
}

func (self *Query) OrRange(index int, criteria map[string]Range) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$or"] = append(self.ao["$and"][index]["$or"], bson.M{"$gte": val.Min, "$lte": val.Max})
	}

	return self
}

func (self *Query) AndLessThanEqualTo(index int, criteria map[string]Min) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$and"] = append(self.ao["$and"][index]["$and"], bson.M{"$lte": val.Min})
	}
	return self
}

func (self *Query) OrLessThanEqualTo(index int, criteria map[string]Min) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$or"] = append(self.ao["$and"][index]["$or"], bson.M{"$lte": val.Min})
	}
	return self
}

func (self *Query) AndLessThan(index int, criteria map[string]Min) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$and"] = append(self.ao["$and"][index]["$and"], bson.M{"$lt": val.Min})
	}
	return self
}

func (self *Query) OrLessThan(index int, criteria map[string]Min) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$or"] = append(self.ao["$and"][index]["$or"], bson.M{"$lt": val.Min})
	}
	return self
}

func (self *Query) AndGreaterThanEqualTo(index int, criteria map[string]Max) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$and"] = append(self.ao["$and"][index]["$and"], bson.M{"$gte": val.Max})
	}
	return self
}

func (self *Query) OrGreaterThanEqualTo(index int, criteria map[string]Max) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$or"] = append(self.ao["$and"][index]["$or"], bson.M{"$gte": val.Max})
	}
	return self
}

func (self *Query) OrGreaterThan(index int, criteria map[string]Max) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$or"] = append(self.ao["$and"][index]["$or"], bson.M{"$gt": val.Max})
	}
	return self
}

func (self *Query) AndGreaterThan(index int, criteria map[string]Max) *Query {
	for key, val := range criteria {
		if key == "Id" {
			key = "_id"
		}
		self.ao["$and"][index]["$and"] = append(self.ao["$and"][index]["$and"], bson.M{"$gt": val.Max})
	}
	return self
}

func (self *Query) AndExclude(index int, criteria map[string]interface{}) *Query {
	self.andOrInNot(index, criteria, "$nin", "$and")
	return self
}

func (self *Query) OrExclude(index int, criteria map[string]interface{}) *Query {
	self.andOrInNot(index, criteria, "$nin", "$or")
	return self
}

func (self *Query) AndIn(index int, criteria map[string]interface{}) *Query {
	self.andOrInNot(index, criteria, "$in", "$and")
	return self

}

func (self *Query) OrIn(index int, criteria map[string]interface{}) *Query {
	self.andOrInNot(index, criteria, "$in", "$or")
	return self

}

func (self *Query) andOrInNot(index int, criteria map[string]interface{}, queryType string, andorType string) {
	for key, value := range criteria {

		if key == "Id" {
			var ids []bson.ObjectId

			k := reflect.TypeOf(value).Kind()
			if k == reflect.Slice || k == reflect.Array {
				values := reflect.ValueOf(value)

				for i := 0; i < values.Len(); i++ {
					val := values.Index(i).Interface()
					objId, err := self.getIdHex(val)
					// if no errors on Id parse, append, else drop it
					if err == nil {
						ids = append(ids, objId)
					}

				}
			} else {
				objId, err := self.getIdHex(value)
				if err != nil {
					self.e = err
					continue
				}

				ids = append(ids, objId)
			}
			self.ao["$and"][index][andorType] = append(self.ao["$and"][index][andorType], Q("_id", bson.M{queryType: ids}))
		} else {
			var valuesToQuery []interface{}

			k := reflect.TypeOf(value).Kind()
			if k == reflect.Slice || k == reflect.Array {
				values := reflect.ValueOf(value)
				for i := 0; i < values.Len(); i++ {
					val := values.Index(i).Interface()
					valuesToQuery = append(valuesToQuery, self.CheckForObjectId(val))
				}
			} else {
				valuesToQuery = append(valuesToQuery, self.CheckForObjectId(value))
			}
			self.ao["$and"][index][andorType] = append(self.ao["$and"][index][andorType], Q(key, bson.M{queryType: valuesToQuery}))
		}
	}
}

func (self *Query) Whitelist(collection string, fields []string) *Query {
	var qff QueryFieldFilter
	qff.CollectionName = collection
	qff.Fields = fields

	idFound := false
	for _, val := range fields {
		if val == "_id" {
			idFound = true
		}
	}

	if idFound == false {
		qff.Fields = append(qff.Fields, "_id")
	}

	self.whiteListed = append(self.whiteListed, qff)
	return self
}

func (self *Query) Blacklist(collection string, fields []string) *Query {
	var qff QueryFieldFilter
	qff.CollectionName = collection
	qff.Fields = fields

	self.blackListed = append(self.blackListed, qff)
	return self
}

func (self *Query) Where(field string, val interface{}) *Query {
	return self.Filter(Q(field, val))
}

func (self *Query) Exists(criteria map[string]interface{}) *Query {

	_, hasId := criteria["Id"]
	if hasId {
		if self.m == nil {
			self.m = make(bson.M)
		}

		self.m["_id"] = Q("$exists", true)
		return self
	} else {

		if self.m == nil {
			self.m = make(bson.M)
		}

		for key, _ := range criteria {
			if key != "" {
				self.m[key] = Q("$exists", true)
			}
		}
	}
	return self
}

func (self *Query) NotExists(criteria map[string]interface{}) *Query {

	_, hasId := criteria["Id"]
	if hasId {
		if self.m == nil {
			self.m = make(bson.M)
		}

		self.m["_id"] = Q("$exists", false)
		return self
	} else {

		if self.m == nil {
			self.m = make(bson.M)
		}

		for key, _ := range criteria {
			if key != "" {
				self.m[key] = Q("$exists", false)
			}
		}
	}
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

		if self.m == nil {
			self.m = make(bson.M)
		}

		self.m["_id"] = objId
		return self
	} else {

		if self.m == nil {
			self.m = make(bson.M)
		}

		for key, val := range criteria {
			if key != "" {
				self.m[key] = self.CheckForObjectId(val)
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

func (self *Query) RenderViews(format DataFormat) *Query {

	self.renderViews = true
	self.format = format

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

					// if no errors on Id parse, append, else drop it
					if err == nil {
						ids = append(ids, objId)
					}
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
					valuesToQuery = append(valuesToQuery, self.CheckForObjectId(val))
				}

			} else {
				valuesToQuery = append(valuesToQuery, self.CheckForObjectId(value))
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

func (self *Query) LessThanEqualTo(criteria map[string]Min) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	for key, val := range criteria {
		if key != "" {
			self.m[key] = bson.M{"$lte": val.Min}
		}
	}

	return self
}

func (self *Query) LessThan(criteria map[string]Min) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	for key, val := range criteria {
		if key != "" {
			self.m[key] = bson.M{"$lt": val.Min}
		}
	}

	return self
}

func (self *Query) GreaterThanEqualTo(criteria map[string]Max) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	for key, val := range criteria {
		if key != "" {
			self.m[key] = bson.M{"$gte": val.Max}
		}
	}

	return self
}

func (self *Query) GreaterThan(criteria map[string]Max) *Query {

	if self.m == nil {
		self.m = make(bson.M)
	}

	for key, val := range criteria {
		if key != "" {
			self.m[key] = bson.M{"$gt": val.Max}
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
	if !serverSettings.WebConfig.Application.LogQueries && serverSettings.WebConfig.Application.LogQueryTimes {
		defer func() {
			log.Println(logger.TimeTrackQuery(time.Now(), "q.All()", self.collection, self.m, self.q))
		}()
	} else if serverSettings.WebConfig.Application.LogQueries {
		defer func() {
			log.Println(logger.TimeTrack(time.Now(), "q.All()"))
		}()
	}
	if self.e != nil {
		return self.e
	}

	q := self.GenerateQuery()

	if !self.stopLog && serverSettings.WebConfig.Application.LogQueries {
		self.LogQuery("q.All()")
	}

	err := q.All(x)

	if err != nil {
		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.All(x)
			if err != nil {
				return err
			}
			return self.processJoinsAndViews(x)
		}

		return self.handleQueryError(err, callback)
	}
	return self.processJoinsAndViews(x)
}

func (self *Query) One(x interface{}) error {
	if !serverSettings.WebConfig.Application.LogQueries && serverSettings.WebConfig.Application.LogQueryTimes {
		defer func() {
			log.Println(logger.TimeTrackQuery(time.Now(), "q.One()", self.collection, self.m, self.q))
		}()
	} else if serverSettings.WebConfig.Application.LogQueries {
		defer func() {
			log.Println(logger.TimeTrack(time.Now(), "q.One()"))
		}()
	}
	if self.e != nil {
		return self.e
	}

	q := self.GenerateQuery()

	if !self.stopLog && serverSettings.WebConfig.Application.LogQueries {
		self.LogQuery("q.One()")
	}

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
			return self.processJoinsAndViews(x)
		}

		return self.handleQueryError(err, callback)
	}
	cnt, _ := q.Count()
	if cnt != 1 {
		return errors.New("Did not return exactly one row.  Returned " + extensions.IntToString(cnt))
	}

	return self.processJoinsAndViews(x)
}

func (self *Query) GetOrCreate(x interface{}, t *Transaction) (err error) {

	count, err := self.Count(x)

	if err != nil {
		return
	}

	if count == 1 {
		err = self.One(x)
		return
	} else if count == 0 {

		valToCall := reflect.ValueOf(x)
		val := reflect.ValueOf(x).Elem()

		for key, value := range self.m {
			if key == "Id" {
				continue
			}
			fieldVal := val.FieldByName(key)
			if fieldVal.CanSet() {
				fieldVal.Set(reflect.ValueOf(value))
			}
		}

		method := valToCall.MethodByName("SaveWithTran")
		in := []reflect.Value{}
		in = append(in, reflect.ValueOf(t))
		values := method.Call(in)
		if values[0].Interface() == nil {
			err = nil
			return
		}
		err = values[0].Interface().(error)

		if err != nil {
			return
		}

		err = self.processJoinsAndViews(x)
		log.Printf("%+v\n\n", x)
		return

	} else {
		err = errors.New("More than one record exists for GetOrCreate.")
	}
	return
}

func (self *Query) TotalRows() int {
	if !serverSettings.WebConfig.Application.LogQueries && serverSettings.WebConfig.Application.LogQueryTimes {
		defer func() {
			log.Println(logger.TimeTrackQuery(time.Now(), "q.TotalRows()", self.collection, self.m, self.q))
		}()
	} else if serverSettings.WebConfig.Application.LogQueries {
		defer func() {
			log.Println(logger.TimeTrack(time.Now(), "q.TotalRows()"))
		}()
	}

	q := self.GenerateQuery()

	if !self.stopLog && serverSettings.WebConfig.Application.LogQueries {
		self.LogQuery("q.Length()")
	}
	count, _ := q.Count()

	return count
}

func (self *Query) Count(x interface{}) (int, error) {
	if !serverSettings.WebConfig.Application.LogQueries && serverSettings.WebConfig.Application.LogQueryTimes {
		defer func() {
			log.Println(logger.TimeTrackQuery(time.Now(), "q.Count()", self.collection, self.m, self.q))
		}()
	} else if serverSettings.WebConfig.Application.LogQueries {
		defer func() {
			log.Println(logger.TimeTrack(time.Now(), "q.Count()"))
		}()
	}
	if self.e != nil {
		return 0, self.e
	}

	q := self.GenerateQuery()

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
	if !serverSettings.WebConfig.Application.LogQueries && serverSettings.WebConfig.Application.LogQueryTimes {
		defer func() {
			log.Println(logger.TimeTrackQuery(time.Now(), "q.Distinct()", self.collection, self.m, self.q))
		}()
	} else if serverSettings.WebConfig.Application.LogQueries {
		defer func() {
			log.Println(logger.TimeTrack(time.Now(), "q.Distinct()"))
		}()
	}

	if self.e != nil {
		return self.e
	}

	q := self.GenerateQuery()

	if !self.stopLog && serverSettings.WebConfig.Application.LogQueries {
		self.LogQuery("q.Distinct()")
	}

	err := q.Distinct(key, x)

	if err != nil {

		// This callback is used for if the Ethernet port is unplugged
		callback := func() error {
			err = q.Distinct(key, x)
			if err != nil {
				return err
			}
			return self.processJoinsAndViews(x)

		}

		return self.handleQueryError(err, callback)
	}

	return self.processJoinsAndViews(x)
}

func (self *Query) processJoinsAndViews(x interface{}) (err error) {
	err = self.processJoins(x)
	if err != nil {
		return
	}

	if self.renderViews {
		err = self.processViews(x)
		if err != nil {
			return
		}
	}

	return
}

func (self *Query) processJoins(x interface{}) (err error) {

	if len(self.joins) > 0 {

		//Check if x is a single struct or an Array
		_, isArray := valueType(x)

		if isArray {
			source := reflect.ValueOf(x).Elem()

			var joins []join
			if source.Len() > 0 {
				joins, err = self.getJoins(source.Index(0))
			}

			if len(joins) == 0 {
				return
			}

			//Advanced way is to get the count and chunk with an IN query.  For now we will iterate.

			for i := 0; i < source.Len(); i++ {
				s := source.Index(i)
				for _, j := range joins {
					id := reflect.ValueOf(self.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()
					joinsField := s.FieldByName("Joins")
					setField := joinsField.FieldByName(j.joinFieldName)

					errJoin := joinField(j, id, setField, j.joinSpecified, self, false, 10)
					if errJoin != nil {
						if j.joinType == "Inner" {
							fields := self.printStruct(s)
							err = errors.New("Failed to Inner Join on " + j.collectionName + " with id = " + id + "\n\nFailed to Join Object: " + s.Type().Name() + " \n" + fields)
							return
						}
					}
				}
			}

		} else {
			source := reflect.ValueOf(x).Elem()

			var joins []join
			joins, err = self.getJoins(source)

			if len(joins) == 0 {
				return
			}

			s := source
			for _, j := range joins {
				id := reflect.ValueOf(self.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()
				joinsField := s.FieldByName("Joins")
				setField := joinsField.FieldByName(j.joinFieldName)

				errJoin := joinField(j, id, setField, j.joinSpecified, self, false, 10)
				if errJoin != nil {
					if j.joinType == "Inner" {
						fields := self.printStruct(source)
						err = errors.New("Failed to Inner Join on " + j.collectionName + " with id = " + id + "\n\tFailed to Join Object:\n" + fields)
						return
					}
				}
			}
		}
		return nil
	}
	return nil
}

func (self *Query) getJoins(x reflect.Value) (joins []join, err error) {

	joinsField := x.FieldByName("Joins")

	if joinsField.Kind() != reflect.Struct {
		err = errors.New("Could not resolve a field due to it not being a struct: " + fmt.Sprintf("%+v", x))
		return
	}

	allJoin, ok := self.joins[JOIN_ALL]
	var hasJoins bool
	if ok {
		for i := 0; i < joinsField.NumField(); i++ {

			typeField := joinsField.Type().Field(i)
			name := typeField.Name

			fmt.Println("getJoins Name")
			fmt.Println(name)
			tagValue := typeField.Tag.Get("join")
			splitValue := strings.Split(tagValue, ",")
			var j join
			j.collectionName = splitValue[0]
			j.joinSchemaName = splitValue[1]
			j.joinFieldRefName = splitValue[2]
			j.isMany = extensions.StringToBool(splitValue[3])
			j.joinForeignFieldName = splitValue[4]
			j.joinFieldName = name
			j.joinSpecified = JOIN_ALL
			j.joinType = allJoin.Type

			//Add WhiteList Fields
			for i := range self.whiteListed {
				wl := &self.whiteListed[i]
				if wl.CollectionName == j.collectionName {
					j.whiteListedFields = wl.Fields
				}
			}

			//Add Blacklist Fields
			for i := range self.blackListed {
				bl := &self.blackListed[i]
				if bl.CollectionName == j.collectionName {
					j.blackListedFields = bl.Fields
				}
			}

			joins = append(joins, j)
		}
	} else {
		for key, val := range self.joins {

			fields := strings.Split(key, ".")
			fieldName := fields[0]

			typeField, ok := joinsField.Type().FieldByName(fieldName)
			if serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println("getJoins fieldName")
				fmt.Println(fieldName)
			}
			if ok == false {
				continue
			}
			hasJoins = true
			tagValue := typeField.Tag.Get("join")
			splitValue := strings.Split(tagValue, ",")
			var j join
			j.collectionName = splitValue[0]
			j.joinSchemaName = splitValue[1]
			j.joinFieldRefName = splitValue[2]
			j.isMany = extensions.StringToBool(splitValue[3])
			j.joinForeignFieldName = splitValue[4]
			j.joinFieldName = fieldName
			j.joinSpecified = strings.Replace(key, fieldName+".", "", 1)
			j.joinType = val.Type

			//Add WhiteList Fields
			for i := range self.whiteListed {
				wl := &self.whiteListed[i]
				if wl.CollectionName == j.collectionName {
					j.whiteListedFields = wl.Fields
				}
			}

			//Add Blacklist Fields
			for i := range self.blackListed {
				bl := &self.blackListed[i]
				if bl.CollectionName == j.collectionName {
					j.blackListedFields = bl.Fields
				}
			}

			joins = append(joins, j)
		}
	}
	if !hasJoins {
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("Could not resolve a field  (getJoins query): " + " on " + x.Type().String() + " object")
		}
		//err = errors.New("Could not resolve a field  (getJoins query): " +  " on " + x.Type().String() + " object")
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

func (self *Query) processViews(x interface{}) (err error) {

	//Check if x is a single struct or an Array
	_, isArray := valueType(x)

	if isArray {

		source := reflect.ValueOf(x).Elem()

		var views []view
		if source.Len() > 0 {
			views = self.getViews(source.Index(0))
		}

		if len(views) == 0 {
			return
		}

		var wg sync.WaitGroup
		for i := 0; i < source.Len(); i++ {

			wg.Add(1)

			s := source.Index(i)
			go func(s reflect.Value) {
				defer func() {
					if r := recover(); r != nil {
						log.Println("Panic at query.go->processViews:  " + fmt.Sprintf("%+v", r))
						return
					}
				}()

				defer wg.Done()
				var wgSetViews sync.WaitGroup
				for _, v := range views { //Update and format the view fields that have ref
					wgSetViews.Add(1)
					defer func() {

						if r := recover(); r != nil {
							log.Println("Panic at query.go->processViews:  " + fmt.Sprintf("%+v", r))
							return
						}
					}()

					go func(v view, s reflect.Value) {
						defer wgSetViews.Done()
						self.setViewValue(v, s)
					}(v, s)
				}
				wgSetViews.Wait()
			}(s)
			wg.Wait()
		}
	} else {

		source := reflect.ValueOf(x)
		if source.Kind() == reflect.Ptr {
			source = source.Elem()
		}

		var views []view
		views = self.getViews(source)

		if len(views) == 0 {
			return
		}

		var wgSetViews sync.WaitGroup
		for _, v := range views { //Update and format the view fields that have ref
			wgSetViews.Add(1)
			go func(v view, s reflect.Value) {

				defer func() {
					if r := recover(); r != nil {
						log.Println("Panic at query.go->processViews:  " + fmt.Sprintf("%+v", r))
						return
					}
				}()

				defer wgSetViews.Done()
				self.setViewValue(v, s)
			}(v, source)
		}
		wgSetViews.Wait()

	}

	return
}

func (self *Query) setViewValue(v view, source reflect.Value) {
	viewValue := dbServices.GetStructReflectionValue(v.ref, source)

	//Check for Bools first
	if viewValue == "true" || viewValue == "false" {
		substitute := ""
		switch v.format {
		case "YesNo":
			substitute = "No"
			if viewValue == "true" {
				substitute = "Yes"
			}
			dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), substitute)

		case "yesno":
			substitute = "no"
			if viewValue == "true" {
				substitute = "yes"
			}
			dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), substitute)
		case "EnabledDisabled":
			substitute = "Disabled"
			if viewValue == "true" {
				substitute = "Enabled"
			}
			dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), substitute)
		case "enableddisabled":
			substitute = "disabled"
			if viewValue == "true" {
				substitute = "enabled"
			}
			dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), substitute)
		case "TrueFalse":
			substitute = "False"
			if viewValue == "true" {
				substitute = "True"
			}
			dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), substitute)
		case "":
			dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), viewValue)
		}
		return
	}

	locale := "en"
	if self.format.Language != "" {
		locale = self.format.Language
	}

	dateFormat := "mm/dd/yyyy"
	if self.format.DateFormat != "" {
		dateFormat = self.substituteDateFormat(self.format.DateFormat)
	} else {
		dateFormat = self.substituteDateFormat(dateFormat)
	}

	timeZone := "US/Eastern"
	if self.format.LocalTimeZone != "" {
		timeZone = self.format.LocalTimeZone
	}

	location, err := time.LoadLocation(timeZone)

	if err != nil {
		log.Println("Failed to Load time.LoadLocation at query.go->SetViewValue:  " + err.Error())
		return
	}

	switch v.format {
	case "DateNumeric":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), dateformatter.Format(t, locale, dateFormat))
	case "DateLong":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), dateformatter.Format(t, locale, "Monday, January 01, 2006"))
	case "DateShort":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), dateformatter.Format(t, locale, "January 01, 2006"))
	case "DateMonthYearShort":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), dateformatter.Format(t, locale, "Jan 2006"))
	case "Time":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), t.Format("03:04:05 PM"))
	case "DateTime":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), dateformatter.Format(t, locale, dateFormat)+" "+t.Format("03:04:05 PM"))
	case "TimeMilitary":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), t.Format("15:04:05"))
	case "DateTimeMilitary":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), dateformatter.Format(t, locale, dateFormat)+" "+t.Format("15:04:05"))
	case "TimeFromNow":
		i, _ := strconv.ParseInt(viewValue, 10, 64)
		t := time.Unix(i, 0).In(location)
		diff := time.Now().Sub(t)
		diffAdd := t.Sub(time.Now())
		self.processTimeFromNow(v.fieldName, source.FieldByName("Views"), diff, diffAdd)
	case "":
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), viewValue)
	}

	if strings.Contains(v.format, "Concatenate:") {
		concatenateData := v.format[12:]
		args := extensions.ExtractArgsWithinBrackets(concatenateData)
		for _, arg := range args {
			substitution := dbServices.GetStructReflectionValue(strings.Title(arg), source)
			concatenateData = strings.Replace(concatenateData, "{"+arg+"}", substitution, 1)
		}
		dbServices.SetFieldValue(v.fieldName, source.FieldByName("Views"), viewValue+concatenateData)
	}

}

func (self *Query) processTimeFromNow(key string, field reflect.Value, diff time.Duration, diffAdd time.Duration) {

	diffToTake := diff

	if diff.Seconds() < 0 {
		diffToTake = diffAdd
	}

	if diffToTake.Minutes() < 1 { //Seconds
		label := "Second"
		if diffToTake.Seconds() > 1 {
			label = "Seconds"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Seconds(), 0)+" "+label) //Translate label
	} else if diffToTake.Hours() < 1 { //Minutes
		label := "Minute"
		if diffToTake.Minutes() > 1 {
			label = "Minutes"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Minutes(), 0)+" "+label) //Translate label
	} else if diffToTake.Hours() < 24 { //Hours
		label := "Hour"
		if diffToTake.Hours() > 1 {
			label = "Hours"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Hours(), 0)+" "+label) //Translate label
	} else if diffToTake.Hours() < 24*7 { //Days
		label := "Day"
		if diffToTake.Hours() > 24*2 {
			label = "Days"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Hours()/24, 0)+" "+label) //Translate label
	} else if diffToTake.Hours() < 24*30 { // Weeks
		label := "Week"
		if diffToTake.Hours() > 24*7+24*3.5 {
			label = "Weeks"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Hours()/(24*7), 0)+" "+label) //Translate label
	} else if diffToTake.Hours() < 24*365 { // Months
		label := "Month"
		if diffToTake.Hours() > 24*30+24*15 {
			label = "Months"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Hours()/(24*30), 0)+" "+label) //Translate label
	} else { // Years
		label := "Year"
		if diffToTake.Hours() > 24*365+24*182.5 {
			label = "Years"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diffToTake.Hours()/(24*365), 0)+" "+label) //Translate label
	}
}

func (self *Query) getViews(x reflect.Value) (views []view) {

	viewsField := x.FieldByName("Views")

	if viewsField.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < viewsField.NumField(); i++ {

		typeField := viewsField.Type().Field(i)
		name := typeField.Name
		tagValue := typeField.Tag.Get("ref")

		if tagValue == "" {
			continue
		}

		tagValues := strings.Split(tagValue, "~")

		var v view
		v.fieldName = name
		v.ref = strings.Title(tagValues[0])
		if len(tagValues) > 1 {
			v.format = tagValues[1]
		}

		views = append(views, v)
	}

	return
}

func (self *Query) substituteDateFormat(dateFormat string) string {
	dateFormat = strings.Replace(dateFormat, "dd", "02", -1)
	dateFormat = strings.Replace(dateFormat, "d", "2", -1)

	dateFormat = strings.Replace(dateFormat, "mmmm", "January", -1)
	dateFormat = strings.Replace(dateFormat, "mmm", "Jan", -1)
	dateFormat = strings.Replace(dateFormat, "mm", "01", -1)
	dateFormat = strings.Replace(dateFormat, "m", "1", -1)

	dateFormat = strings.Replace(dateFormat, "yyyy", "2006", -1)
	dateFormat = strings.Replace(dateFormat, "yy", "06", -1)

	return dateFormat
}

func (self *Query) LogQuery(functionName string) {
	if serverSettings.WebConfig.Application.LogQueryStackTraces {
		caller := stacktrace.Errorf("GoCore caller:")
		core.Debug.Dump("Desc-> Called Function query.go#"+functionName, "Desc->Caller for Query:", caller.ErrorStack(), core.Debug.GetDump("Desc->Collection", self.collection, "Desc->Limit", self.limit, "Desc->Skip", self.skip, "Desc->Sort", self.sort, "Desc->mgo Query", self.q, "Desc->Queryset", self.m, "Desc->Joins", self.joins))
	} else {
		core.Debug.Dump("Desc-> Called Function query.go#"+functionName, core.Debug.GetDump("Desc->Collection", self.collection, "Desc->Limit", self.limit, "Desc->Skip", self.skip, "Desc->Sort", self.sort, "Desc->mgo Query", self.q, "Desc->Queryset", self.m, "Desc->Joins", self.joins))
	}
}

func (self *Query) handleQueryError(err error, callback queryError) error {

	if self.isDBConnectionError(err) {

		for i := 0; i < 2; i++ {

			log.Println("Attempting to Refresh Mongo Session")
			dbServices.DBMutex.Lock()
			dbServices.MongoSession.Refresh()
			dbServices.DBMutex.Unlock()

			err = callback()
			if !self.isDBConnectionError(err) {
				return err
			}

			time.Sleep(200 * time.Millisecond)
		}
	}

	if serverSettings.WebConfig.Application.LogQueries {
		core.Debug.Dump("Desc->You got a mongo error!!!! ", err.Error())
		self.LogQuery("handleQueryError")
	}

	return err

}

func (self *Query) isDBConnectionError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "Closed explicitly") ||
		strings.Contains(err.Error(), "read: operation timed out") ||
		strings.Contains(err.Error(), "i/o timeout") ||
		strings.Contains(err.Error(), "read: connection reset by peer") ||
		strings.Contains(err.Error(), "EOF") {
		return true
	}
	return false
}

func (self *Query) Iter() (qi *QueryIterator) {
	qi = new(QueryIterator)
	qi.q = self
	qi.iter = self.GenerateQuery().Iter()
	return
}

func (self *Query) GenerateQuery() *mgo.Query {

	q := self.collection.Find(bson.M{})

	if self.q != nil {
		q = self.q
	}

	if self.ao != nil {
		core.Debug.Dump(self.ao)
		q = self.collection.Find(self.ao)
	} else if self.m != nil {
		q = self.collection.Find(self.m)
	} else if self.o != nil {
		q = self.collection.Find(bson.M{"$or": self.o})
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

	//Add WhiteList Fields
	for i := range self.whiteListed {
		wl := &self.whiteListed[i]
		if wl.CollectionName == self.entityName {
			q = q.Select(whiteList(wl.Fields))
			self.filteredFields = true
		}
	}

	//Add Blacklist Fields
	for i := range self.blackListed {
		bl := &self.blackListed[i]
		if bl.CollectionName == self.entityName {

			var whiteListFields []string

			obj := ResolveEntity(self.entityName)
			reflectedFields := obj.Reflect()

			for i := 0; i < len(reflectedFields); i++ {
				if !reflectedFields[i].IsView {

					addField := true
					for j := range bl.Fields {
						blField := bl.Fields[j]
						if blField == reflectedFields[i].Name {
							addField = false
							break
						}
					}
					if addField {
						whiteListFields = append(whiteListFields, reflectedFields[i].Name)
						self.filteredFields = true
					}
				}
			}

			q = q.Select(whiteList(whiteListFields))
			self.filteredFields = true
		}
	}

	return q
}

func whiteList(q []string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

func (self *Query) getIdHex(val interface{}) (bson.ObjectId, error) {

	myIdType := reflect.TypeOf(val)
	myIdInstance := reflect.ValueOf(val)
	myIdInstanceHex := myIdInstance.MethodByName("Hex")

	if myIdType.Name() != "ObjectId" && (myIdType.Kind() == reflect.String && bson.IsObjectIdHex(val.(string)) == false) {
		return bson.NewObjectId(), errors.New(ERROR_INVALID_ID_VALUE)
	}

	if myIdType.Name() != "ObjectId" && myIdType.Kind() == reflect.String && val.(string) != "" {
		return bson.ObjectIdHex(val.(string)), nil
	} else if myIdType.Name() == "ObjectId" && myIdType.Kind() == reflect.String {
		return bson.ObjectIdHex(myIdInstanceHex.Call([]reflect.Value{})[0].String()), nil
	} else {
		return bson.NewObjectId(), errors.New(ERROR_INVALID_ID_VALUE)
	}
}

func (self *Query) CheckForObjectId(val interface{}) interface{} {

	myIdType := reflect.TypeOf(val)
	myIdInstance := reflect.ValueOf(val)
	myIdInstanceHex := myIdInstance.MethodByName("Hex")

	if myIdType.Name() == "ObjectId" && myIdType.Kind() == reflect.String {
		return myIdInstanceHex.Call([]reflect.Value{})[0].String()
	} else {
		return val
	}
}

func (self *Query) printStruct(s reflect.Value) string {
	fields := ""
	for i := 0; i < s.NumField(); i++ {
		valField := s.Field(i)
		typeField := s.Type().Field(i)
		name := typeField.Name
		if valField.Kind() == reflect.Ptr {
			// fields +=  name + ":*\n"
		} else if valField.Kind() == reflect.Struct {
			// fields += name + ":{}\n"
		} else if valField.Kind() == reflect.Array || valField.Kind() == reflect.Slice {
			// fields += name + ":{}\n"
		} else if valField.Kind() == reflect.String {
			if name == "Id" {
				myIdInstanceHex := valField.MethodByName("Hex")
				fields += "\t" + name + ":  " + myIdInstanceHex.Call([]reflect.Value{})[0].String() + "\n"
			} else {
				fields += "\t" + name + ":  " + valField.String() + "\n"
			}
		} else if valField.Kind() == reflect.Int || valField.Kind() == reflect.Int8 || valField.Kind() == reflect.Int16 || valField.Kind() == reflect.Int32 || valField.Kind() == reflect.Int64 {
			fields += "\t" + name + ":  " + strconv.FormatInt(valField.Int(), 10) + "\n"
		} else if valField.Kind() == reflect.Bool {
			fields += "\t" + name + ":  " + extensions.BoolToString(valField.Bool()) + "\n"
		} else if valField.Kind() == reflect.Float32 {
			fields += "\t" + name + ":  " + strconv.FormatFloat(valField.Float(), 'E', -1, 32) + "\n"
		} else if valField.Kind() == reflect.Float64 {
			fields += "\t" + name + ":  " + strconv.FormatFloat(valField.Float(), 'E', -1, 64) + "\n"
		}
	}
	return fields
}
