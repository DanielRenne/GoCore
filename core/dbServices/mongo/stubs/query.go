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

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/extensions"
	dateformatter "github.com/altipla-consulting/i18n-dateformatter"
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
	joinType         string
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
	q           *mgo.Query
	m           bson.M
	limit       int
	skip        int
	sort        []string
	collection  *mgo.Collection
	e           error
	joins       map[string]joinType
	format      DataFormat
	renderViews bool
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
			return self.processJoinsAndViews(x)
		}

		return self.handleQueryError(err, callback)
	}

	return self.processJoinsAndViews(x)

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
				self.m[key] = self.checkForObjectId(val)
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
					valuesToQuery = append(valuesToQuery, self.checkForObjectId(val))
				}

			} else {
				valuesToQuery = append(valuesToQuery, self.checkForObjectId(value))
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
			return self.processJoinsAndViews(x)
		}

		return self.handleQueryError(err, callback)
	}
	return self.processJoinsAndViews(x)
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
				joins = self.getJoins(source.Index(0))
			}

			if len(joins) == 0 {
				return
			}

			//Advanced way is to get the count and chunk with an IN query.  For now we will iterate.

			for i := 0; i < source.Len(); i++ {
				s := source.Index(i)
				for _, j := range joins {
					id := reflect.ValueOf(s.FieldByName(j.joinFieldRefName).Interface()).String()
					joinsField := s.FieldByName("Joins")
					setField := joinsField.FieldByName(j.joinFieldName)

					errJoin := joinField(j.joinSchemaName, j.collectionName, id, setField, j.joinSpecified, self, false, 10)
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
			joins = self.getJoins(source)

			if len(joins) == 0 {
				return
			}

			s := source
			for _, j := range joins {
				id := reflect.ValueOf(s.FieldByName(j.joinFieldRefName).Interface()).String()
				joinsField := s.FieldByName("Joins")
				setField := joinsField.FieldByName(j.joinFieldName)

				errJoin := joinField(j.joinSchemaName, j.collectionName, id, setField, j.joinSpecified, self, false, 10)
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

func (self *Query) getJoins(x reflect.Value) (joins []join) {

	joinsField := x.FieldByName("Joins")

	if joinsField.Kind() != reflect.Struct {
		return
	}

	allJoin, ok := self.joins[JOIN_ALL]

	if ok {
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
			j.joinType = allJoin.Type
			joins = append(joins, j)
		}
	} else {
		for key, val := range self.joins {

			fields := strings.Split(key, ".")
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
			j.joinSpecified = strings.Replace(key, fieldName+".", "", 1)
			j.joinType = val.Type
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
				defer wg.Done()
				var wgSetViews sync.WaitGroup
				for _, v := range views { //Update and format the view fields that have ref
					wgSetViews.Add(1)
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

	location, _ := time.LoadLocation(timeZone)

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
		self.processTimeFromNow(v.fieldName, source.FieldByName("Views"), diff)
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

func (self *Query) processTimeFromNow(key string, field reflect.Value, diff time.Duration) {
	if diff.Minutes() < 1 { //Seconds
		label := "Second"
		if diff.Seconds() > 1 {
			label = "Seconds"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Seconds(), 0)+" "+label) //Translate label
	} else if diff.Hours() < 1 { //Minutes
		label := "Minute"
		if diff.Minutes() > 1 {
			label = "Minutes"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Minutes(), 0)+" "+label) //Translate label
	} else if diff.Hours() < 24 { //Hours
		label := "Hour"
		if diff.Hours() > 1 {
			label = "Hours"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Hours(), 0)+" "+label) //Translate label
	} else if diff.Hours() < 24*7 { //Days
		label := "Day"
		if diff.Hours() > 24*2 {
			label = "Days"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Hours()/24, 0)+" "+label) //Translate label
	} else if diff.Hours() < 24*30 { // Weeks
		label := "Week"
		if diff.Hours() > 24*7+24*3.5 {
			label = "Weeks"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Hours()/(24*7), 0)+" "+label) //Translate label
	} else if diff.Hours() < 24*365 { // Months
		label := "Month"
		if diff.Hours() > 24*30+24*15 {
			label = "Months"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Hours()/(24*30), 0)+" "+label) //Translate label
	} else { // Years
		label := "Year"
		if diff.Hours() > 24*365+24*182.5 {
			label = "Years"
		}
		dbServices.SetFieldValue(key, field, extensions.FloatToString(diff.Hours()/(24*365), 0)+" "+label) //Translate label
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

func (self *Query) checkForObjectId(val interface{}) interface{} {

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
