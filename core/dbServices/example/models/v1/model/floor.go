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

var Floors modelFloors

type modelFloors struct{}

var collectionFloorsMutex *sync.RWMutex

type FloorJoinItems struct {
	Count int      `json:"Count"`
	Items *[]Floor `json:"Items"`
}

var GoCoreFloorsHasBootStrapped atomicTypes.AtomicBool

var mongoFloorsCollection *mgo.Collection

func init() {
	store.RegisterStore(Floors)
	collectionFloorsMutex = &sync.RWMutex{}
}

func (self *Floor) GetId() string {
	return self.Id.Hex()
}

type Floor struct {
	Id            bson.ObjectId  `json:"Id" bson:"_id,omitempty"`
	Name          string         `json:"Name" bson:"Name" validate:"true,,,,,,"`
	SiteId        string         `json:"SiteId" dbIndex:"index" bson:"SiteId" validate:"true,,,,,,"`
	BuildingId    string         `json:"BuildingId" dbIndex:"index" bson:"BuildingId" validate:"true,,,,,,"`
	AccountId     string         `json:"AccountId" dbIndex:"index" bson:"AccountId" validate:"true,,,,,,"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id         string `json:"Id"`
		Name       string `json:"Name"`
		SiteId     string `json:"SiteId"`
		BuildingId string `json:"BuildingId"`
		AccountId  string `json:"AccountId"`
	} `json:"Errors" bson:"-"`

	Views struct {
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
	} `json:"Views" bson:"-"`

	Joins struct {
		Site           *Site     `json:"Site,omitempty" join:"Sites,Site,SiteId,false,"`
		Building       *Building `json:"Building,omitempty" join:"Buildings,Building,BuildingId,false,"`
		Account        *Account  `json:"Account,omitempty" join:"Accounts,Account,AccountId,false,"`
		LastUpdateUser *User     `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
	} `json:"Joins" bson:"-"`
}

func (obj modelFloors) SetCollection(mdb *mgo.Database) {
	collectionFloorsMutex.Lock()
	mongoFloorsCollection = mdb.C("Floors")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoFloorsCollection.Create(&ci)
	collectionFloorsMutex.Unlock()
}

func (obj modelFloors) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj Floor
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj *Floor) DoesIdExist(objectID interface{}) bool {
	var retObj Floor
	row := modelFloors{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}
func (obj modelFloors) NewByReflection() (value reflect.Value) {
	retObj := Floor{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelFloors) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []Floor
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

func (obj modelFloors) CountByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (count int, err error) {
	var retObj []Floor
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

func (obj modelFloors) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionFloorsMutex.RLock()
		collection := mongoFloorsCollection
		bootstrapped := GoCoreFloorsHasBootStrapped.Get()
		collectionFloorsMutex.RUnlock()

		if collection != nil && bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 1000)
		if elapseMs%10000 == 0 {
			log.Println("Floors has not bootstrapped and has yet to get a collection pointer")
		}
	}
	collectionFloorsMutex.RLock()
	collection := mongoFloorsCollection
	collectionFloorsMutex.RUnlock()
	query.collection = collection
	query.entityName = "Floor"
	return query
}
func (obj modelFloors) RemoveAll() {
	var elapseMs int
	collection := mongoFloorsCollection
	for {
		bootstrapped := GoCoreFloorsHasBootStrapped.Get()

		if collection != nil && bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 1000)
		if elapseMs%10000 == 0 {
			log.Println("Floors has not bootstrapped and has yet to get a collection pointer")
		}
	}
	collection.RemoveAll(bson.M{})
	return
}
func (obj modelFloors) Index() error {
	log.Println("Building Indexes for MongoDB collection Floors:")
	for key, value := range dbServices.GetDBIndexes(Floor{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}

		collectionFloorsMutex.RLock()
		collection := mongoFloorsCollection
		collectionFloorsMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for Floor." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for Floor." + key)
		}
	}
	return nil
}

func (obj modelFloors) BootStrapComplete() {
	GoCoreFloorsHasBootStrapped.Set(true)
}
func (obj modelFloors) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Floors Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	collectionFloorsMutex.RLock()
	query.collection = mongoFloorsCollection
	collectionFloorsMutex.RUnlock()
	var rows []Floor
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwogIHsKICAgICJJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNTAiLAogICAgIk5hbWUiOiAiQm9vdHN0cmFwcGVkIEZsb29yIDEiLAogICAgIlNpdGVJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGQiLAogICAgIkJ1aWxkaW5nSWQiOiAiNjMzZTIxNDEyYTFiNDlmNDMxZWU2ZjRlIiwKICAgICJBY2NvdW50SWQiOiAiNjMzZTIxNDEyYTFiNDlmNDMxZWU2ZjRjIgogIH0sCiAgewogICAgIklkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY1MSIsCiAgICAiTmFtZSI6ICJCb290c3RyYXBwZWQgRmxvb3IgMiIsCiAgICAiU2l0ZUlkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0ZCIsCiAgICAiQnVpbGRpbmdJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGUiLAogICAgIkFjY291bnRJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGMiCiAgfSwKICB7CiAgICAiSWQiOiAiNjMzZTIxNDEyYTFiNDlmNDMxZWU2ZjUyIiwKICAgICJOYW1lIjogIkJvb3RzdHJhcHBlZCBGbG9vciAxIiwKICAgICJTaXRlSWQiOiAiNjMzZTIxNDEyYTFiNDlmNDMxZWU2ZjRkIiwKICAgICJCdWlsZGluZ0lkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0ZiIsCiAgICAiQWNjb3VudElkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0YyIKICB9LAogIHsKICAgICJJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNTMiLAogICAgIk5hbWUiOiAiQm9vdHN0cmFwcGVkIEZsb29yIDIiLAogICAgIlNpdGVJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGQiLAogICAgIkJ1aWxkaW5nSWQiOiAiNjMzZTIxNDEyYTFiNDlmNDMxZWU2ZjRmIiwKICAgICJBY2NvdW50SWQiOiAiNjMzZTIxNDEyYTFiNDlmNDMxZWU2ZjRjIgogIH0KXQo="

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Floors")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Floors due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("floors", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Floors: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Floors: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []Floor
	for _, file := range files {
		var fileBootstrap []Floor
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Floors", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Floors: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Floors: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Floors", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Floors")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Floors", len(v))

	for _, doc := range v {
		var original Floor
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Floors" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Floors:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for Floors:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Floors skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Floors")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("floors", "Floors")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Floors")
			if actualCount != originalCount {
				logger.Message("Floors counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelFloors) New() *Floor {
	return &Floor{}
}

func (obj *Floor) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *Floor) Save() error {
	if !AllowWrites {
		return nil
	}
	collectionFloorsMutex.RLock()
	collection := mongoFloorsCollection
	collectionFloorsMutex.RUnlock()
	t := time.Now()
	objectId := self.Id
	if self.Id == "" {
		objectId = bson.NewObjectId()
		self.CreateDate = t
	}
	self.UpdateDate = t
	changeInfo, err := collection.UpsertId(objectId, &self)
	if err != nil {
		log.Println("Failed to upsertId for Floor:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	dbServices.CollectionCache{}.Remove("Floors", self.Id.Hex())
	if store.OnChangeRecord != nil && len(store.OnRecordUpdate) > 0 {
		if store.OnRecordUpdate[0] == "*" || utils.InArray("Floors", store.OnRecordUpdate) {
			value := reflect.ValueOf(&self)
			store.OnChangeRecord("Floors", self.Id.Hex(), value.Interface())
		}
	}
	pubsub.Publish("Floors.Save", self)
	return nil
}

func (self *Floor) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *Floor) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *Floor) CreateWithTran(t *Transaction, forceCreate bool) error {

	transactionQueue.Lock()
	defer func() {
		transactionQueue.Unlock()
	}()

	// collectionFloorsMutex.RLock()
	// collection := mongoFloorsCollection
	// collectionFloorsMutex.RUnlock()
	// if collection == nil {
	// 	initFloors()
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

	t.Collections = append(t.Collections, "FloorsHistory")
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
	dbServices.CollectionCache{}.Remove("Floors", self.Id.Hex())
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
	var histRecord FloorHistoryRecord
	histRecord.TId = t.Id.Hex()
	histRecord.Data = newBson
	histRecord.Type = TRANSACTION_CHANGETYPE_INSERT

	histRecord.ObjId = self.Id.Hex()
	histRecord.CreateDate = time.Now()
	//Get the Original Record if it is a Update
	if isUpdate {

		_, ok := transactionQueue.queue[t.Id.Hex()].newItems["Floor_"+self.Id.Hex()]
		if ok {
			transactionQueue.queue[t.Id.Hex()].newItems["Floor_"+self.Id.Hex()] = eTransactionNew
		}
		histRecord.Type = TRANSACTION_CHANGETYPE_UPDATE
		eTransactionNew.changeType = TRANSACTION_CHANGETYPE_UPDATE
		var original Floor
		err := Floors.Query().ById(self.Id, &original)
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
	transactionQueue.queue[t.Id.Hex()].newItems["Floor_"+self.Id.Hex()] = eTransactionNew
	transactionQueue.queue[t.Id.Hex()].originalItems["Floor_"+self.Id.Hex()] = eTransactionOriginal
	return nil
}
func (self *Floor) ValidateAndClean() error {

	return validateFields(Floor{}, self, reflect.ValueOf(self).Elem())
}

func (self *Floor) Reflect() []Field {

	return Reflect(Floor{})
}

func (self *Floor) Delete() error {
	dbServices.CollectionCache{}.Remove("Floors", self.Id.Hex())
	err := mongoFloorsCollection.RemoveId(self.Id)
	if err == nil {
		pubsub.Publish("Floors.Delete", self)
	}
	return err
}

func (self *Floor) DeleteWithTran(t *Transaction) error {
	transactionQueue.Lock()
	defer func() {
		transactionQueue.Unlock()
	}()
	if self.Id.Hex() == "" {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_RECORD_NOT_EXISTS)
	}

	dbServices.CollectionCache{}.Remove("Floors", self.Id.Hex())
	_, ok := transactionQueue.queue[t.Id.Hex()]
	if ok == false {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_NOT_PRESENT)
	}

	var histRecord FloorHistoryRecord
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

	var original Floor
	err := Floors.Query().ById(self.Id, &original)

	if err != nil {
		return err
	}

	originalJson, err := original.JSONString()

	if err != nil {
		return err
	}

	originalBase64 := getBase64(originalJson)
	histRecord.Data = originalBase64

	t.Collections = append(t.Collections, "FloorsHistory")
	if len(transactionQueue.queue[t.Id.Hex()].originalItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].originalItems = make(map[string]entityTransaction, 0)
	}
	if len(transactionQueue.queue[t.Id.Hex()].newItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].newItems = make(map[string]entityTransaction, 0)
	}
	transactionQueue.queue[t.Id.Hex()].newItems["Floors_"+self.Id.Hex()] = eTransactionNew

	transactionQueue.queue[t.Id.Hex()].originalItems["Floors_"+self.Id.Hex()] = eTransactionOriginal

	transactionQueue.ids[t.Id.Hex()] = append(transactionQueue.ids[t.Id.Hex()], eTransactionNew.entity.GetId())
	return nil
}

func (self *Floor) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *Floor) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Floor) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Floor) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *Floor) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *Floor) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *Floor) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelFloors) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Id":
		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Name":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "SiteId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "BuildingId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountId":
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

func (obj modelFloors) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "SiteId":
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
	case "BuildingId":
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
	case "AccountId":
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
	case "Name":
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
	}
	return
}
