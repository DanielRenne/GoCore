package model

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/atomicTypes"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/pubsub"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"reflect"
	"sync"
	"time"
)

var Users modelUsers

type modelUsers struct{}

var collectionUsersMutex *sync.RWMutex

type UserJoinItems struct {
	Count int     `json:"Count"`
	Items *[]User `json:"Items"`
}

var GoCoreUsersHasBootStrapped atomicTypes.AtomicBool

var mongoUsersCollection *mgo.Collection

func init() {
	store.RegisterStore(Users)
	collectionUsersMutex = &sync.RWMutex{}
}

func (self *User) GetId() string {
	return self.Id.Hex()
}

type User struct {
	Id            bson.ObjectId  `json:"Id" bson:"_id,omitempty"`
	First         string         `json:"First" bson:"First" validate:"true,,,,,,"`
	Last          string         `json:"Last" bson:"Last" validate:"true,,,,,,"`
	SignupDate    time.Time      `json:"SignupDate" bson:"SignupDate"`
	Email         string         `json:"Email" dbIndex:"unique" bson:"Email" validate:"true,email,,,,,"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id         string `json:"Id"`
		First      string `json:"First"`
		Last       string `json:"Last"`
		SignupDate string `json:"SignupDate"`
		Email      string `json:"Email"`
	} `json:"Errors" bson:"-"`
}

func (obj modelUsers) SetCollection(mdb *mgo.Database) {
	collectionUsersMutex.Lock()
	mongoUsersCollection = mdb.C("Users")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoUsersCollection.Create(&ci)
	collectionUsersMutex.Unlock()
}

func (obj modelUsers) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj User
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj *User) DoesIdExist(objectID interface{}) bool {
	var retObj User
	row := modelUsers{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}
func (obj modelUsers) NewByReflection() (value reflect.Value) {
	retObj := User{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelUsers) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []User
	q := obj.Query().Filter(filter)
	if len(inFilter) > 0 {
		q = q.In(inFilter)
	}
	if len(excludeFilter) > 0 {
		q = q.Exclude(excludeFilter)
	}
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.All(&retObj)
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelUsers) CountByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (count int, err error) {
	var retObj []User
	q := obj.Query().Filter(filter)
	if len(inFilter) > 0 {
		q = q.In(inFilter)
	}
	if len(excludeFilter) > 0 {
		q = q.Exclude(excludeFilter)
	}
	// joins really make no sense here but just copy paste coding here
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	cnt, errCount := q.Count(&retObj)
	return cnt, errCount
}

func (obj modelUsers) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionUsersMutex.RLock()
		collection := mongoUsersCollection
		bootstrapped := GoCoreUsersHasBootStrapped.Get()
		collectionUsersMutex.RUnlock()

		if collection != nil && bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 1000)
		if elapseMs%10000 == 0 {
			log.Println("Users has not bootstrapped and has yet to get a collection pointer")
		}
	}
	collectionUsersMutex.RLock()
	collection := mongoUsersCollection
	collectionUsersMutex.RUnlock()
	query.collection = collection
	query.entityName = "User"
	return query
}
func (obj modelUsers) RemoveAll() {
	var elapseMs int
	collection := mongoUsersCollection
	for {
		bootstrapped := GoCoreUsersHasBootStrapped.Get()

		if collection != nil && bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 1000)
		if elapseMs%10000 == 0 {
			log.Println("Users has not bootstrapped and has yet to get a collection pointer")
		}
	}
	collection.RemoveAll(bson.M{})
	return
}
func (obj modelUsers) Index() error {
	log.Println("Building Indexes for MongoDB collection Users:")
	for key, value := range dbServices.GetDBIndexes(User{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}

		collectionUsersMutex.RLock()
		collection := mongoUsersCollection
		collectionUsersMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for User." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for User." + key)
		}
	}
	return nil
}

func (obj modelUsers) BootStrapComplete() {
	GoCoreUsersHasBootStrapped.Set(true)
}
func (obj modelUsers) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Users Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	collectionUsersMutex.RLock()
	query.collection = mongoUsersCollection
	collectionUsersMutex.RUnlock()
	var rows []User
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwogIHsKICAgICJJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGIiLAogICAgIkZpcnN0IjogIkJvb3RzdHJhcHBlZCBHbyIsCiAgICAiTGFzdCI6ICJCb290c3RyYXBwZWQgQ29yZSIsCiAgICAiRW1haWwiOiAiQm9vdHN0cmFwcGVkdGVzdGluZzEyMzU2Nzg5QHRlc3QuY29tIgogIH0KXQo="

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Users")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Users due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("users", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Users: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Users: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []User
	for _, file := range files {
		var fileBootstrap []User
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Users", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Users: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Users: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Users", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Users")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Users", len(v))

	for _, doc := range v {
		var original User
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Users" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Users:  " + doc.Id.Hex() + "  " + err.Error())
					isError = true
				}
			} else {
				valid := 0x01
				var reason map[string]bool
				reason = make(map[string]bool, 0)

				if doc.BootstrapMeta != nil && doc.BootstrapMeta.Version > 0 && doc.BootstrapMeta.Version <= serverSettings.WebConfig.Application.VersionNumeric {
					valid &= 0x00
					reason["Version Mismatch"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.Domain != "" && doc.BootstrapMeta.Domain != serverSettings.WebConfig.Application.ServerFQDN {
					valid &= 0x00
					reason["FQDN Mismatch With Domain"] = true
				}
				if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.Domains) > 0 && !utils.InArray(serverSettings.WebConfig.Application.ServerFQDN, doc.BootstrapMeta.Domains) {
					valid &= 0x00
					reason["FQDN Mismatch With Domains"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.ProductName != "" && doc.BootstrapMeta.ProductName != serverSettings.WebConfig.Application.ProductName {
					valid &= 0x00
					reason["ProductName does not Match"] = true
				}
				if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.ProductNames) > 0 && !utils.InArray(serverSettings.WebConfig.Application.ProductName, doc.BootstrapMeta.ProductNames) {
					valid &= 0x00
					reason["ProductNames does not Match Product"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.ReleaseMode != "" && doc.BootstrapMeta.ReleaseMode != serverSettings.WebConfig.Application.ReleaseMode {
					valid &= 0x00
					reason["ReleaseMode does not match"] = true
				}

				if valid == 0x01 {
					actualCount += 1
					err = doc.Save()
					if err != nil {
						log.Println("Failed to bootstrap data for Users:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Users skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Users")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("users", "Users")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Users")
			if actualCount != originalCount {
				logger.Message("Users counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelUsers) New() *User {
	return &User{}
}

func (obj *User) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *User) Save() error {
	if !AllowWrites {
		return nil
	}
	collectionUsersMutex.RLock()
	collection := mongoUsersCollection
	collectionUsersMutex.RUnlock()
	t := time.Now()
	objectId := self.Id
	if self.Id == "" {
		objectId = bson.NewObjectId()
		self.CreateDate = t
	}
	self.UpdateDate = t
	changeInfo, err := collection.UpsertId(objectId, &self)
	if err != nil {
		log.Println("Failed to upsertId for User:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	if store.OnChangeRecord != nil && len(store.OnRecordUpdate) > 0 {
		if store.OnRecordUpdate[0] == "*" || utils.InArray("Users", store.OnRecordUpdate) {
			value := reflect.ValueOf(&self)
			store.OnChangeRecord("Users", self.Id.Hex(), value.Interface())
		}
	}
	pubsub.Publish("Users.Save", self)
	return nil
}

func (self *User) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *User) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *User) CreateWithTran(t *Transaction, forceCreate bool) error {

	transactionQueue.Lock()
	defer func() {
		transactionQueue.Unlock()
	}()

	// collectionUsersMutex.RLock()
	// collection := mongoUsersCollection
	// collectionUsersMutex.RUnlock()
	// if collection == nil {
	// 	initUsers()
	// }
	//Validate the Model first.  If it fails then clean up the transaction in memory
	err := self.ValidateAndClean()
	if err != nil {
		delete(transactionQueue.queue, t.Id.Hex())
		return err
	}

	_, ok := transactionQueue.queue[t.Id.Hex()]
	if ok == false {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_NOT_PRESENT)
	}

	t.Collections = append(t.Collections, "UsersHistory")
	isUpdate := true
	if self.Id.Hex() == "" {
		isUpdate = false
		self.Id = bson.NewObjectId()
		self.CreateDate = time.Now()
	}
	if len(transactionQueue.queue[t.Id.Hex()].originalItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].originalItems = make(map[string]entityTransaction, 0)
	}
	if len(transactionQueue.queue[t.Id.Hex()].newItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].newItems = make(map[string]entityTransaction, 0)
	}
	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	if forceCreate {
		isUpdate = false
	}
	self.UpdateDate = time.Now()
	self.LastUpdateId = t.UserId
	newBson, err := self.BSONString()
	if err != nil {
		return err
	}

	var eTransactionNew entityTransaction
	eTransactionNew.changeType = TRANSACTION_CHANGETYPE_INSERT
	eTransactionNew.entity = self
	var histRecord UserHistoryRecord
	histRecord.TId = t.Id.Hex()
	histRecord.Data = newBson
	histRecord.Type = TRANSACTION_CHANGETYPE_INSERT

	histRecord.ObjId = self.Id.Hex()
	histRecord.CreateDate = time.Now()
	//Get the Original Record if it is a Update
	if isUpdate {

		_, ok := transactionQueue.queue[t.Id.Hex()].newItems["User_"+self.Id.Hex()]
		if ok {
			transactionQueue.queue[t.Id.Hex()].newItems["User_"+self.Id.Hex()] = eTransactionNew
		}
		histRecord.Type = TRANSACTION_CHANGETYPE_UPDATE
		eTransactionNew.changeType = TRANSACTION_CHANGETYPE_UPDATE
		var original User
		err := Users.Query().ById(self.Id, &original)
		if err == nil {
			// Found a match of an existing record, lets save history now on it
			originalBson, err := original.BSONString()
			if err != nil {
				return err
			}
			histRecord.Data = originalBson
		}
	}
	var eTransactionOriginal entityTransaction
	eTransactionOriginal.entity = &histRecord
	transactionQueue.ids[t.Id.Hex()] = append(transactionQueue.ids[t.Id.Hex()], eTransactionNew.entity.GetId())
	transactionQueue.queue[t.Id.Hex()].newItems["User_"+self.Id.Hex()] = eTransactionNew
	transactionQueue.queue[t.Id.Hex()].originalItems["User_"+self.Id.Hex()] = eTransactionOriginal
	return nil
}
func (self *User) ValidateAndClean() error {

	return validateFields(User{}, self, reflect.ValueOf(self).Elem())
}

func (self *User) Reflect() []Field {

	return Reflect(User{})
}

func (self *User) Delete() error {
	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	err := mongoUsersCollection.RemoveId(self.Id)
	if err == nil {
		pubsub.Publish("Users.Delete", self)
	}
	return err
}

func (self *User) DeleteWithTran(t *Transaction) error {
	transactionQueue.Lock()
	defer func() {
		transactionQueue.Unlock()
	}()
	if self.Id.Hex() == "" {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_RECORD_NOT_EXISTS)
	}

	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	_, ok := transactionQueue.queue[t.Id.Hex()]
	if ok == false {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_NOT_PRESENT)
	}

	var histRecord UserHistoryRecord
	histRecord.TId = t.Id.Hex()

	histRecord.Type = TRANSACTION_CHANGETYPE_DELETE
	histRecord.ObjId = self.Id.Hex()
	histRecord.CreateDate = time.Now()
	var eTransactionNew entityTransaction
	eTransactionNew.changeType = TRANSACTION_CHANGETYPE_DELETE
	eTransactionNew.entity = self

	var eTransactionOriginal entityTransaction
	eTransactionOriginal.changeType = TRANSACTION_CHANGETYPE_DELETE
	eTransactionOriginal.entity = &histRecord

	var original User
	err := Users.Query().ById(self.Id, &original)

	if err != nil {
		return err
	}

	originalJson, err := original.JSONString()

	if err != nil {
		return err
	}

	originalBase64 := getBase64(originalJson)
	histRecord.Data = originalBase64

	t.Collections = append(t.Collections, "UsersHistory")
	if len(transactionQueue.queue[t.Id.Hex()].originalItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].originalItems = make(map[string]entityTransaction, 0)
	}
	if len(transactionQueue.queue[t.Id.Hex()].newItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].newItems = make(map[string]entityTransaction, 0)
	}
	transactionQueue.queue[t.Id.Hex()].newItems["Users_"+self.Id.Hex()] = eTransactionNew

	transactionQueue.queue[t.Id.Hex()].originalItems["Users_"+self.Id.Hex()] = eTransactionOriginal

	transactionQueue.ids[t.Id.Hex()] = append(transactionQueue.ids[t.Id.Hex()], eTransactionNew.entity.GetId())
	return nil
}

func (self *User) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

	source := reflect.ValueOf(self).Elem()

	var joins []join
	joins, err = getJoins(source, remainingRecursions)

	if len(joins) == 0 {
		return
	}

	s := source
	for _, j := range joins {
		id := reflect.ValueOf(q.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()
		joinsField := s.FieldByName("Joins")
		setField := joinsField.FieldByName(j.joinFieldName)

		endRecursion := false
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Print("Remaining Recursions")
			fmt.Println(fmt.Sprintf("%+v", remainingRecursions))
			fmt.Println(fmt.Sprintf("%+v", j.collectionName))
		}
		if remainingRecursions == j.joinSpecified {
			endRecursion = true
		}
		err = joinField(j, id, setField, j.joinSpecified, q, endRecursion, recursionCount)
		if err != nil {
			return
		}
	}
	return
}

func (self *User) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *User) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *User) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *User) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *User) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *User) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelUsers) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "SignupDate":
		obj, ok := x.(time.Time)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Email":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Id":
		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "First":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Last":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	}
	return
}

func (obj modelUsers) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "First":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Last":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "SignupDate":
		if x == nil {
			var obj time.Time
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(time.Time)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Email":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Id":
		if x == nil {
			var obj bson.ObjectId
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	}
	return
}
