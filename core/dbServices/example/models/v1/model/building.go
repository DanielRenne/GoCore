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

var Buildings modelBuildings

type modelBuildings struct{}

var collectionBuildingsMutex *sync.RWMutex

type BuildingJoinItems struct {
	Count int         `json:"Count"`
	Items *[]Building `json:"Items"`
}

var GoCoreBuildingsHasBootStrapped atomicTypes.AtomicBool

var mongoBuildingsCollection *mgo.Collection

func init() {
	store.RegisterStore(Buildings)
	collectionBuildingsMutex = &sync.RWMutex{}
}

func (self *Building) GetId() string {
	return self.Id.Hex()
}

type Building struct {
	Id            bson.ObjectId  `json:"Id" bson:"_id,omitempty"`
	Name          string         `json:"Name" bson:"Name" validate:"true,,,,,,"`
	ImageCustom   string         `json:"ImageCustom" bson:"ImageCustom"`
	ImageFileName string         `json:"ImageFileName" bson:"ImageFileName"`
	SiteId        string         `json:"SiteId" dbIndex:"index" bson:"SiteId" validate:"true,,,,,,"`
	AccountId     string         `json:"AccountId" dbIndex:"index" bson:"AccountId" validate:"true,,,,,,"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id            string `json:"Id"`
		Name          string `json:"Name"`
		ImageCustom   string `json:"ImageCustom"`
		ImageFileName string `json:"ImageFileName"`
		SiteId        string `json:"SiteId"`
		AccountId     string `json:"AccountId"`
	} `json:"Errors" bson:"-"`

	Views struct {
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
	} `json:"Views" bson:"-"`

	Joins struct {
		FileObjects    *FileObject     `json:"FileObjects,omitempty" join:"FileObjects,FileObject,ImageCustom,false,"`
		Floors         *FloorJoinItems `json:"Floors,omitempty" join:"Floors,Floor,Id,true,BuildingId"`
		Site           *Site           `json:"Site,omitempty" join:"Sites,Site,SiteId,false,"`
		Account        *Account        `json:"Account,omitempty" join:"Accounts,Account,AccountId,false,"`
		LastUpdateUser *User           `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
	} `json:"Joins" bson:"-"`
}

func (obj modelBuildings) SetCollection(mdb *mgo.Database) {
	collectionBuildingsMutex.Lock()
	mongoBuildingsCollection = mdb.C("Buildings")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoBuildingsCollection.Create(&ci)
	collectionBuildingsMutex.Unlock()
}

func (obj modelBuildings) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj Building
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj *Building) DoesIdExist(objectID interface{}) bool {
	var retObj Building
	row := modelBuildings{}
	q := row.Query()
	err := q.ById(objectID, &retObj)
	if err == nil {
		return true
	} else {
		return false
	}
}
func (obj modelBuildings) NewByReflection() (value reflect.Value) {
	retObj := Building{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelBuildings) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []Building
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

func (obj modelBuildings) CountByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (count int, err error) {
	var retObj []Building
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

func (obj modelBuildings) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionBuildingsMutex.RLock()
		collection := mongoBuildingsCollection
		bootstrapped := GoCoreBuildingsHasBootStrapped.Get()
		collectionBuildingsMutex.RUnlock()

		if collection != nil && bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 1000)
		if elapseMs%10000 == 0 {
			log.Println("Buildings has not bootstrapped and has yet to get a collection pointer")
		}
	}
	collectionBuildingsMutex.RLock()
	collection := mongoBuildingsCollection
	collectionBuildingsMutex.RUnlock()
	query.collection = collection
	query.entityName = "Building"
	return query
}
func (obj modelBuildings) RemoveAll() {
	var elapseMs int
	collection := mongoBuildingsCollection
	for {
		bootstrapped := GoCoreBuildingsHasBootStrapped.Get()

		if collection != nil && bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 1000)
		if elapseMs%10000 == 0 {
			log.Println("Buildings has not bootstrapped and has yet to get a collection pointer")
		}
	}
	collection.RemoveAll(bson.M{})
	return
}
func (obj modelBuildings) Index() error {
	log.Println("Building Indexes for MongoDB collection Buildings:")
	for key, value := range dbServices.GetDBIndexes(Building{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}

		collectionBuildingsMutex.RLock()
		collection := mongoBuildingsCollection
		collectionBuildingsMutex.RUnlock()
		err := collection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for Building." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for Building." + key)
		}
	}
	return nil
}

func (obj modelBuildings) BootStrapComplete() {
	GoCoreBuildingsHasBootStrapped.Set(true)
}
func (obj modelBuildings) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Buildings Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	collectionBuildingsMutex.RLock()
	query.collection = mongoBuildingsCollection
	collectionBuildingsMutex.RUnlock()
	var rows []Building
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwogIHsKICAgICJJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGUiLAogICAgIk5hbWUiOiAiQm9vdHN0cmFwcGVkIFdlYlNvY2tldCBCdWlsZGluZyIsCiAgICAiSW1hZ2VDdXN0b20iOiAiMTIzIiwKICAgICJJbWFnZUZpbGVOYW1lIjogIiIsCiAgICAiU2l0ZUlkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0ZCIsCiAgICAiQWNjb3VudElkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0YyIKICB9LAogIHsKICAgICJJZCI6ICI2MzNlMjE0MTJhMWI0OWY0MzFlZTZmNGYiLAogICAgIk5hbWUiOiAiQm9vdHN0cmFwcGVkIENoYW5uZWxzIGFuZCBHb1JvdXRpbmUgQnVpbGRpbmciLAogICAgIkltYWdlQ3VzdG9tIjogIjEyMzQ1IiwKICAgICJJbWFnZUZpbGVOYW1lIjogIiIsCiAgICAiU2l0ZUlkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0ZCIsCiAgICAiQWNjb3VudElkIjogIjYzM2UyMTQxMmExYjQ5ZjQzMWVlNmY0YyIKICB9Cl0K"

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Buildings")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Buildings due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("buildings", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Buildings: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Buildings: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []Building
	for _, file := range files {
		var fileBootstrap []Building
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Buildings", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Buildings: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Buildings: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Buildings", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Buildings")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Buildings", len(v))

	for _, doc := range v {
		var original Building
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Buildings" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Buildings:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for Buildings:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Buildings skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Buildings")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("buildings", "Buildings")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Buildings")
			if actualCount != originalCount {
				logger.Message("Buildings counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelBuildings) New() *Building {
	return &Building{}
}

func (obj *Building) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *Building) Save() error {
	if !AllowWrites {
		return nil
	}
	collectionBuildingsMutex.RLock()
	collection := mongoBuildingsCollection
	collectionBuildingsMutex.RUnlock()
	t := time.Now()
	objectId := self.Id
	if self.Id == "" {
		objectId = bson.NewObjectId()
		self.CreateDate = t
	}
	self.UpdateDate = t
	changeInfo, err := collection.UpsertId(objectId, &self)
	if err != nil {
		log.Println("Failed to upsertId for Building:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	dbServices.CollectionCache{}.Remove("Buildings", self.Id.Hex())
	if store.OnChangeRecord != nil && len(store.OnRecordUpdate) > 0 {
		if store.OnRecordUpdate[0] == "*" || utils.InArray("Buildings", store.OnRecordUpdate) {
			value := reflect.ValueOf(&self)
			store.OnChangeRecord("Buildings", self.Id.Hex(), value.Interface())
		}
	}
	pubsub.Publish("Buildings.Save", self)
	return nil
}

func (self *Building) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *Building) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *Building) CreateWithTran(t *Transaction, forceCreate bool) error {

	transactionQueue.Lock()
	defer func() {
		transactionQueue.Unlock()
	}()

	// collectionBuildingsMutex.RLock()
	// collection := mongoBuildingsCollection
	// collectionBuildingsMutex.RUnlock()
	// if collection == nil {
	// 	initBuildings()
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

	t.Collections = append(t.Collections, "BuildingsHistory")
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
	dbServices.CollectionCache{}.Remove("Buildings", self.Id.Hex())
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
	var histRecord BuildingHistoryRecord
	histRecord.TId = t.Id.Hex()
	histRecord.Data = newBson
	histRecord.Type = TRANSACTION_CHANGETYPE_INSERT

	histRecord.ObjId = self.Id.Hex()
	histRecord.CreateDate = time.Now()
	//Get the Original Record if it is a Update
	if isUpdate {

		_, ok := transactionQueue.queue[t.Id.Hex()].newItems["Building_"+self.Id.Hex()]
		if ok {
			transactionQueue.queue[t.Id.Hex()].newItems["Building_"+self.Id.Hex()] = eTransactionNew
		}
		histRecord.Type = TRANSACTION_CHANGETYPE_UPDATE
		eTransactionNew.changeType = TRANSACTION_CHANGETYPE_UPDATE
		var original Building
		err := Buildings.Query().ById(self.Id, &original)
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
	transactionQueue.queue[t.Id.Hex()].newItems["Building_"+self.Id.Hex()] = eTransactionNew
	transactionQueue.queue[t.Id.Hex()].originalItems["Building_"+self.Id.Hex()] = eTransactionOriginal
	return nil
}
func (self *Building) ValidateAndClean() error {

	return validateFields(Building{}, self, reflect.ValueOf(self).Elem())
}

func (self *Building) Reflect() []Field {

	return Reflect(Building{})
}

func (self *Building) Delete() error {
	dbServices.CollectionCache{}.Remove("Buildings", self.Id.Hex())
	err := mongoBuildingsCollection.RemoveId(self.Id)
	if err == nil {
		pubsub.Publish("Buildings.Delete", self)
	}
	return err
}

func (self *Building) DeleteWithTran(t *Transaction) error {
	transactionQueue.Lock()
	defer func() {
		transactionQueue.Unlock()
	}()
	if self.Id.Hex() == "" {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_RECORD_NOT_EXISTS)
	}

	dbServices.CollectionCache{}.Remove("Buildings", self.Id.Hex())
	_, ok := transactionQueue.queue[t.Id.Hex()]
	if ok == false {
		return errors.New(dbServices.ERROR_CODE_TRANSACTION_NOT_PRESENT)
	}

	var histRecord BuildingHistoryRecord
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

	var original Building
	err := Buildings.Query().ById(self.Id, &original)

	if err != nil {
		return err
	}

	originalJson, err := original.JSONString()

	if err != nil {
		return err
	}

	originalBase64 := getBase64(originalJson)
	histRecord.Data = originalBase64

	t.Collections = append(t.Collections, "BuildingsHistory")
	if len(transactionQueue.queue[t.Id.Hex()].originalItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].originalItems = make(map[string]entityTransaction, 0)
	}
	if len(transactionQueue.queue[t.Id.Hex()].newItems) == 0 {
		transactionQueue.queue[t.Id.Hex()].newItems = make(map[string]entityTransaction, 0)
	}
	transactionQueue.queue[t.Id.Hex()].newItems["Buildings_"+self.Id.Hex()] = eTransactionNew

	transactionQueue.queue[t.Id.Hex()].originalItems["Buildings_"+self.Id.Hex()] = eTransactionOriginal

	transactionQueue.ids[t.Id.Hex()] = append(transactionQueue.ids[t.Id.Hex()], eTransactionNew.entity.GetId())
	return nil
}

func (self *Building) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *Building) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Building) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Building) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *Building) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *Building) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *Building) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelBuildings) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

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
	case "ImageCustom":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "ImageFileName":
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

func (obj modelBuildings) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

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
	case "ImageCustom":
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
	case "ImageFileName":
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
