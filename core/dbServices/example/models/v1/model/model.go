package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/asaskevich/govalidator"
	"github.com/fatih/camelcase"
	"github.com/globalsign/mgo/bson"
)

const (
	TRANSACTION_DATATYPE_ORIGINAL = 1
	TRANSACTION_DATATYPE_NEW      = 2

	TRANSACTION_CHANGETYPE_INSERT = 1
	TRANSACTION_CHANGETYPE_UPDATE = 2
	TRANSACTION_CHANGETYPE_DELETE = 3

	MGO_RECORD_NOT_FOUND = "not found"

	VALIDATION_ERROR                   = "ValidationError"
	VALIDATION_ERROR_REQUIRED          = "ValidationErrorRequiredFieldMissing"
	VALIDATION_ERROR_EMAIL             = "ValidationErrorInvalidEmail"
	VALIDATION_ERROR_SPECIFIC_REQUIRED = "ValidationFieldSpecificRequired"
	VALIDATION_ERROR_SPECIFIC_EMAIL    = "ValidationFieldSpecificEmailRequired"
)

type modelEntity interface {
	Save() error
	Delete() error
	SaveWithTran(*Transaction) error
	Reflect() []Field
	JoinFields(string, *Query, int) error
	GetId() string
	DoesIdExist(interface{}) bool
}

type modelCollection interface {
	Rollback(transactionId string) error
}

type collection interface {
	Query() *Query
}

// BootstrapMeta is a struct that is used to tell the bootstrap how or what you want to bootstrap
type BootstrapMeta struct {
	// Version is the version that you specifically want to bootstrap only to
	Version int `json:"Version" bson:"Version"`
	// Domain would be used to bootstrap only a specific domain
	Domain string `json:"Domain" bson:"Domain"`
	// ReleaseMode would be used to bootstrap only a specific release mode such as debug records
	ReleaseMode string `json:"ReleaseMode" bson:"ReleaseMode"`
	// ProductName would be used to bootstrap only a specific product name
	ProductName string `json:"ProductName" bson:"ProductName"`
	// Domains would be used to bootstrap only many domains
	Domains []string `json:"Domains" bson:"Domains"`
	// ProductNames would be used to bootstrap many product name
	ProductNames []string `json:"ProductNames" bson:"ProductNames"`
	// DeleteRow, if true, will delete the row from the bootstrap table after it is previousl bootstrapped
	DeleteRow bool `json:"DeleteRow" bson:"DeleteRow"`
	// AlwaysUpdate, if true, will always update the row from the bootstrap table after it is previously bootstrapped
	AlwaysUpdate bool `json:"AlwaysUpdate" bson:"AlwaysUpdate"`
}

type bootstrapSync struct {
	sync.Mutex
	Items [][]byte
}

type tQueue struct {
	sync.RWMutex
	queue map[string]*transactionsToPersist
	ids   map[string][]string
}

type transactionsToPersist struct {
	t             *Transaction
	newItems      map[string]entityTransaction
	originalItems map[string]entityTransaction
	startTime     time.Time
}

type entityTransaction struct {
	changeType int
	committed  bool
	entity     modelEntity
}

type Field struct {
	Name       string
	Label      string
	DataType   string
	IsView     bool
	Validation *dbServices.FieldValidation
}

var transactionQueue tQueue
var AllowWrites bool

func init() {
	AllowWrites = true
	transactionQueue.ids = make(map[string][]string)
	transactionQueue.queue = make(map[string]*transactionsToPersist)
	go clearTransactionQueue()
	go setupCollections()
}

func ConnectDB() {
	var ts []Transaction
	Transactions.Query().All(&ts)
}

func setupCollections() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at model.setupCollections(): " + fmt.Sprintf("%+v", r))
			return
		}
	}()

	for {

		mdb := dbServices.ReadMongoDB()
		if mdb != nil {

			//Make a connection to the Version Collection
			versionsCollection := mdb.C("Version")

			type Version struct {
				Id    bson.ObjectId `bson:"_id"`
				Value string        `bson:"Value"`
			}

			var version Version
			errFindVersion := versionsCollection.FindId(bson.ObjectIdHex("60942d9bab99e73ea651f967")).One(&version)

			// log.Printf("VERSION VERSION VERSION VERSION VERSION VERSION VERSION %+v", version)

			collectionNames := GetCollectionNames()

			for _, name := range collectionNames {
				for {
					collection, ok := store.GetCollection(name)
					if ok {
						log.Println("Initializing Collection " + name + " Binary Version:  " + store.Version + " DB Version: " + version.Value)
						collection.SetCollection(mdb)
						if errFindVersion != nil {
							collection.Index()
							collection.Bootstrap()
						} else if store.Version == "" {
							collection.Index()
							collection.Bootstrap()
						} else if version.Value != store.Version {
							collection.Index()
							collection.Bootstrap()
						} else {
							collection.BootStrapComplete()
						}
						break
					}
					time.Sleep(time.Millisecond * 1000)
				}
			}

			collectionHistoryNames := GetCollectionHistoryNames()

			for _, name := range collectionHistoryNames {
				for {
					collection, ok := store.GetCollectionHistory(name)
					if ok {
						// log.Println("Initializing" + name )
						collection.SetCollection(mdb)
						if store.Version == "" {
							collection.Index()
						} else if version.Value != store.Version {
							collection.Index()
						}
						break
					}
					time.Sleep(time.Millisecond * 1000)
				}
			}

			if store.Version == "" {
				UpdateAllRecordsToLatestSchema()
			} else if version.Value != store.Version {
				UpdateAllRecordsToLatestSchema()
			}

			if store.Version != "" {

				version.Id = bson.ObjectIdHex("60942d9bab99e73ea651f967")
				version.Value = store.Version
				versionsCollection.UpsertId(version.Id, &version)
				log.Printf("APPLIED VERSION %+v TO DATABASE", version)
			}

			break
		}

		time.Sleep(time.Millisecond * 5)
	}

}

// Q is a helper function to pass to things like Filter to filter a field value and key
func Q(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{k: v}
}

// QTs is a helper function to pass to things like Filter to filter a field value and a time.Time
func QTs(k string, v time.Time) map[string]time.Time {
	return map[string]time.Time{k: v}
}

// RangeQ is a helper function to pass to things like Filter to filter a field value and key
func RangeQ(k string, min interface{}, max interface{}) map[string]Range {
	var rge map[string]Range
	rge = make(map[string]Range)
	rge[k] = Range{
		Max: max,
		Min: min,
	}
	return rge
}

// MinQ is a helper function to pass to things like Filter to filter a field value and min value
func MinQ(k string, min interface{}) map[string]Min {
	var rge map[string]Min
	rge = make(map[string]Min)
	rge[k] = Min{
		Min: min,
	}
	return rge
}

// MaxQ is a helper function to pass to things like Filter to filter a field value and a max value
func MaxQ(k string, max interface{}) map[string]Max {
	var rge map[string]Max
	rge = make(map[string]Max)
	rge[k] = Max{
		Max: max,
	}
	return rge
}

// Every 12 hours check the transactionQueue and remove any outstanding stale transactions > 48 hours old
func clearTransactionQueue() {

	transactionQueue.Lock()

	for key, value := range transactionQueue.queue {

		if time.Since(value.startTime).Hours() > 48 {
			delete(transactionQueue.queue, key)
		}
	}

	transactionQueue.Unlock()

	time.Sleep(12 * time.Hour)
	clearTransactionQueue()
}

func getBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(data[:]), nil
}

func getNow() time.Time {
	return time.Now()
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == VALIDATION_ERROR || err.Error() == VALIDATION_ERROR_EMAIL {
		return true
	}
	return false
}

func validateFields(x interface{}, objectToUpdate interface{}, val reflect.Value) error {

	isError := false
	collection := strings.Replace(reflect.TypeOf(x).String()+"s", "model.", "", -1)

	obj := reflect.ValueOf(objectToUpdate)
	methodGetID := obj.MethodByName("GetId")
	inID := []reflect.Value{}
	objectID := methodGetID.Call(inID)

	for key, value := range dbServices.GetValidationTags(x) {

		fieldValue := dbServices.GetReflectionFieldValue(key, objectToUpdate)
		validations := strings.Split(value, ",")

		if validations[0] != "" {
			if err := validateRequired(fieldValue, validations[0]); err != nil {
				dbServices.SetFieldValue("Errors."+key, val, VALIDATION_ERROR_SPECIFIC_REQUIRED)
				if store.OnChange != nil {
					store.OnChange(collection, objectID[0].String(), "Errors."+key, VALIDATION_ERROR_SPECIFIC_REQUIRED, nil)
				}
				isError = true
			}
		}
		if validations[1] != "" {

			cleanup, err := validateType(fieldValue, validations[1])

			if err != nil {
				if err.Error() == VALIDATION_ERROR_EMAIL {
					dbServices.SetFieldValue("Errors."+key, val, VALIDATION_ERROR_SPECIFIC_EMAIL)
					if store.OnChange != nil {
						store.OnChange(collection, objectID[0].String(), "Errors."+key, VALIDATION_ERROR_SPECIFIC_EMAIL, nil)
					}
				}
				isError = true
			}

			if cleanup != "" {
				dbServices.SetFieldValue(key, val, cleanup)
			}

		}

	}

	if isError {
		return errors.New(VALIDATION_ERROR)
	}

	return nil
}

func validateRequired(value string, tagValue string) error {
	if tagValue == "true" {
		if value == "" {
			return errors.New(VALIDATION_ERROR_REQUIRED)
		}
		return nil
	}
	return nil
}

func validateType(value string, tagValue string) (string, error) {
	switch tagValue {
	case dbServices.VALIDATION_TYPE_EMAIL:
		return "", validateEmail(value)
	}
	return "", nil
}

func validateEmail(value string) error {
	if !govalidator.IsEmail(value) {
		return errors.New(VALIDATION_ERROR_EMAIL)
	}
	return nil
}

func getJoins(x reflect.Value, remainingRecursions string) (joins []join, err error) {
	if remainingRecursions == "" {
		return
	}

	fields := strings.Split(remainingRecursions, ".")
	fieldName := fields[0]

	joinsField := x.FieldByName("Joins")
	if joinsField.Kind() != reflect.Struct {
		return
	}

	if fieldName == JOIN_ALL {
		for i := 0; i < joinsField.NumField(); i++ {

			typeField := joinsField.Type().Field(i)
			name := typeField.Name
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
			joins = append(joins, j)
		}
	} else {
		typeField, ok := joinsField.Type().FieldByName(fieldName)
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println(fmt.Sprintf("%+v", fields), joinsField.Kind(), ok)
		}

		if ok == false {
			msg := "Could not resolve a field (getJoins model): " + fmt.Sprintf("%+v", remainingRecursions) + " on " + x.Type().String() + " object"
			if serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println(msg)
			}
			//new := errors.New(msg)
			//err = new
			return
		}
		name := typeField.Name
		tagValue := typeField.Tag.Get("join")
		splitValue := strings.Split(tagValue, ",")
		var j join
		j.collectionName = splitValue[0]
		j.joinSchemaName = splitValue[1]
		j.joinFieldRefName = splitValue[2]
		j.isMany = extensions.StringToBool(splitValue[3])
		j.joinForeignFieldName = splitValue[4]
		j.joinFieldName = name
		j.joinSpecified = strings.Replace(remainingRecursions, fieldName+".", "", 1)
		if strings.Contains(j.joinSpecified, "Count") && j.joinSpecified[:5] == "Count" {
			j.joinSpecified = "Count"
		}
		joins = append(joins, j)
	}
	return
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func Reflect(obj interface{}) []Field {
	var ret []Field
	val := reflect.ValueOf(obj)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if typeField.Name != "Errors" && typeField.Name != "Joins" && typeField.Name != "BootstrapMeta" {
			if typeField.Name == "Views" {
				for f := 0; f < val.FieldByName("Views").NumField(); f++ {
					field := Field{}
					field.IsView = true
					name := val.FieldByName("Views").Type().Field(f).Name
					namePart := camelcase.Split(name)
					for x := 0; x < len(namePart); x++ {
						if x > 0 {
							namePart[x] = strings.ToLower(namePart[x])
						}
					}
					field.Name = val.FieldByName("Views").Type().Field(f).Name
					field.Label = strings.Join(namePart[:], " ")
					field.DataType = val.FieldByName("Views").Type().Field(f).Type.Name()
					validate := val.FieldByName("Views").Type().Field(f).Tag.Get("validate")
					if validate != "" {
						field.Validation = &dbServices.FieldValidation{}
						parts := strings.Split(validate, ",")
						field.Validation.Required = extensions.StringToBool(parts[0])
						field.Validation.Type = parts[1]
						field.Validation.Min = parts[2]
						field.Validation.Max = parts[3]
						field.Validation.Length = parts[4]
						field.Validation.LengthMax = parts[5]
						field.Validation.LengthMin = parts[6]
					}
					ret = append(ret, field)
				}
			} else {
				field := Field{}
				validate := typeField.Tag.Get("validate")
				if validate != "" {
					field.Validation = &dbServices.FieldValidation{}
					parts := strings.Split(validate, ",")
					field.Validation.Required = extensions.StringToBool(parts[0])
					field.Validation.Type = parts[1]
					field.Validation.Min = parts[2]
					field.Validation.Max = parts[3]
					field.Validation.Length = parts[4]
					field.Validation.LengthMax = parts[5]
					field.Validation.LengthMin = parts[6]
				}
				name := typeField.Name
				namePart := camelcase.Split(name)
				for x := 0; x < len(namePart); x++ {
					if x > 0 {
						namePart[x] = strings.ToLower(namePart[x])
					}
				}
				field.Name = typeField.Name
				field.Label = strings.Join(namePart[:], " ")
				field.DataType = typeField.Type.Name()
				ret = append(ret, field)
			}
		}
	}
	return ret
}

func JoinEntity(collectionQ *Query, y interface{}, j join, id string, fieldToSet reflect.Value, remainingRecursions string, q *Query, endRecursion bool, recursionCount int) (err error) {

	defer func() {
		if r := recover(); r != nil {
			msg := "Panic Recovered at model.JoinEntity:  Failed to join " + j.joinSchemaName + " with id:" + id + "  Error:" + fmt.Sprintf("%+v", r)
			err = errors.New(msg)
			if serverSettings.WebConfig.Application.LogJoinQueries && err != nil {
				fmt.Println("err Recursion Line 399->" + fmt.Sprintf("%+v", err))
			}
			return
		}
	}()

	if serverSettings.WebConfig.Application.LogJoinQueries {
		fmt.Println("!!!!!!!!!!!!ddd!!!!!!!!")
		fmt.Printf("%+v", fieldToSet)
		fmt.Println("!!!!!!!!!!!!id!!!!!!!!")
		fmt.Printf("%+v", id)
	}

	//Add any whitelisting or blacklisting of fields
	if len(j.whiteListedFields) > 0 {
		collectionQ = collectionQ.Whitelist(collectionQ.entityName, j.whiteListedFields)
	} else if len(j.blackListedFields) > 0 {
		collectionQ = collectionQ.Blacklist(collectionQ.entityName, j.blackListedFields)
	}

	if IsZeroOfUnderlyingType(fieldToSet.Interface()) || j.isMany {
		if j.isMany && id != "" {
			if remainingRecursions == "Count" {
				cnt, err := collectionQ.ToggleLogFlag(true).Filter(Q(j.joinForeignFieldName, id)).Count(y)
				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("JoinEntity() Recursion Count Only err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
				}
				if serverSettings.WebConfig.Application.LogJoinQueries && err != nil {
					fmt.Println("err Recursion Line 413->" + fmt.Sprintf("%+v", err))
				}
				if err != nil {
					// err = errCnt
					return err
				}
				countField := fieldToSet.Elem().FieldByName("Count")
				countField.Set(reflect.ValueOf(cnt))
				return err
			}
			err = collectionQ.ToggleLogFlag(true).Filter(Q(j.joinForeignFieldName, id)).All(y)
			if serverSettings.WebConfig.Application.LogJoinQueries {
				collectionQ.LogQuery("JoinEntity({" + j.joinForeignFieldName + ": " + id + "}) Recursion Many err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
			}
		} else if id != "" {
			if j.joinForeignFieldName == "" {
				err = collectionQ.ToggleLogFlag(true).ById(id, y)
				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("JoinEntity() Recursion Single By Id (" + id + ") err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
				}
			} else {
				err = collectionQ.ToggleLogFlag(true).Filter(Q(j.joinForeignFieldName, id)).One(y)
				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("JoinEntity({" + j.joinForeignFieldName + ": " + id + "}) Recursion Single err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
				}
			}
		}

		if err == nil && id != "" {
			if endRecursion == false && recursionCount > 0 {
				recursionCount--

				in := []reflect.Value{}
				in = append(in, reflect.ValueOf(remainingRecursions))
				in = append(in, reflect.ValueOf(q))
				in = append(in, reflect.ValueOf(recursionCount))

				if j.isMany {

					myArray := reflect.ValueOf(y).Elem()
					for i := 0; i < myArray.Len(); i++ {
						s := myArray.Index(i)
						method := s.Addr().MethodByName("JoinFields")
						values := method.Call(in)
						if values[0].Interface() != nil {
							err = values[0].Interface().(error)
						}
					}
				} else {
					err = CallMethod(y, "JoinFields", in)
				}
			}
			if err != nil && serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println("err Recursion Line 465->" + fmt.Sprintf("%+v", err))
			}
			if err == nil {
				if j.isMany {

					itemsField := fieldToSet.Elem().FieldByName("Items")
					countField := fieldToSet.Elem().FieldByName("Count")
					itemsField.Set(reflect.ValueOf(y))
					countField.Set(reflect.ValueOf(reflect.ValueOf(y).Elem().Len()))
					//if serverSettings.WebConfig.Application.LogJoinQueries {
					//fmt.Println("!!!!!!!!!!!!reflected pointer for many!!!!!!!!")
					//fmt.Printf("%+v", itemsField)
					//fmt.Println("!!!!!!!!!!!!test interface!!!!!!!!")
					//fmt.Printf("%+v", itemsField.Interface())
					//}
				} else {
					fieldToSet.Set(reflect.ValueOf(y))
					//if serverSettings.WebConfig.Application.LogJoinQueries {
					//fmt.Println("!!!!!!!!!!!!reflected pointer for single row!!!!!!!!")
					//fmt.Printf("%+v", fieldToSet)
					//fmt.Println("!!!!!!!!!!!!test interface!!!!!!!!")
					//fmt.Printf("%+v", fieldToSet.Interface())
					//}

				}

				//if serverSettings.WebConfig.Application.LogJoinQueries {
				//	fmt.Println("!!!!!!!!!!!!reflected and set value to!!!!!!!!")
				//	fmt.Printf("%+v", reflect.ValueOf(y))
				//}

				if q.renderViews {
					err = q.processViews(y)
					if err != nil && serverSettings.WebConfig.Application.LogJoinQueries {
						collectionQ.LogQuery("err Recursion Line 479->" + fmt.Sprintf("%+v", err))
					}
					if err != nil {
						return
					}
				}

			}
		} else {
			if serverSettings.WebConfig.Application.LogJoinQueries && err != nil {
				fmt.Println("err Recursion Line 495->" + fmt.Sprintf("%+v", err))
			}
		}
	} else {
		if endRecursion == false && recursionCount > 0 {
			recursionCount--
			method := fieldToSet.MethodByName("JoinFields")
			in := []reflect.Value{}
			in = append(in, reflect.ValueOf(remainingRecursions))
			in = append(in, reflect.ValueOf(q))
			in = append(in, reflect.ValueOf(recursionCount))
			values := method.Call(in)
			if values[0].Interface() == nil {

				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("Recursion returning due to nil values[0] interface")
				}
				err = nil
				return
			}
			err = values[0].Interface().(error)
			if err != nil && serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println("err Recursion Line 503->" + fmt.Sprintf("%+v", err))
			}
		}
	}
	return
}

func CallMethod(i interface{}, methodName string, in []reflect.Value) (err error) {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		values := finalMethod.Call(in)
		if values[0].Interface() == nil {
			err = nil
			return
		}
		err = values[0].Interface().(error)
		return
	}

	// return or panic, method not found of either type
	return nil
}

func NewObjectId() string {
	return bson.NewObjectId().Hex()
}

func BootstrapMongoDump(directoryName string, collectionName string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at model.BootstrapMongoDump(): " + fmt.Sprintf("%+v", r))
			return
		}
	}()

	path := serverSettings.APP_LOCATION + "/db/bootstrap/" + directoryName + "/mongoDump"

	if extensions.DoesFileExist(path) == false {
		return
	}

	path = path + "/" + directoryName + "Dump.json"

	if extensions.DoesFileExist(path) == false {
		return
	}

	dbName := serverSettings.WebConfig.DbConnection.Database

	var commandPath string
	if runtime.GOOS == "linux" {
		commandPath = "/usr/bin/mongoimport"
	} else if runtime.GOOS == "darwin" {
		commandPath = "mongoimport"
	}

	if os.Getenv("MGO_COMMAND_ARGS") != "" {
		args := strings.Fields(commandPath + " " + os.Getenv("MGO_COMMAND_ARGS") + " --collection " + collectionName + " --file " + path + " --upsert")
		err = exec.Command(args[0], args[1:]...).Run()
		return
	}

	err = exec.Command(commandPath, "--db", dbName, "--collection", collectionName, "--file", path, "--upsert").Run()
	return

}

func BootstrapDirectory(directoryName string, collectionCount int) (files [][]byte, err error, directoryFound bool) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at model.BootstrapDirectory(): " + fmt.Sprintf("%+v", r))
			return
		}
	}()

	var syncedItems bootstrapSync
	var wg sync.WaitGroup
	path := serverSettings.APP_LOCATION + "/db/bootstrap/" + directoryName + "/dist"

	if extensions.DoesFileExist(path) == false {
		return
	}

	directoryFound = true

	err = fileCache.LoadCachedManifestFromKeyIntoMemory(directoryName)
	if err != nil {
		return
	}

	err = filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) (err error) {

		if errWalk != nil {
			err = errWalk
			return
		}

		var readFile bool
		if !f.IsDir() && !fileCache.DoesHashExistInManifestCache(directoryName, f.Name()) {
			fileCache.UpdateManifestMemoryCache(directoryName, f.Name(), extensions.Int64ToInt32(f.Size()))
			readFile = true
		}

		if !f.IsDir() && fileCache.DoesHashExistInManifestCache(directoryName, f.Name()) {
			var cachedSize int
			fileCache.ByteManifest.Lock()
			cachedSize = fileCache.ByteManifest.Cache[directoryName][f.Name()]
			fileCache.ByteManifest.Unlock()
			actualRetailPrice := extensions.Int64ToInt32(f.Size())
			if cachedSize != actualRetailPrice {
				fileCache.UpdateManifestMemoryCache(directoryName, f.Name(), extensions.Int64ToInt32(f.Size()))
				readFile = true
				log.Println(f.Name() + " is being read because of a difference in size (cached:" + extensions.IntToString(cachedSize) + " new bytes: " + extensions.IntToString(actualRetailPrice) + ")")
			}
		}

		if f.IsDir() || collectionCount == 0 {
			readFile = true
		}

		if readFile && filepath.Ext(f.Name()) == ".json" {
			wg.Add(1)

			go func() {
				defer wg.Done()
				jsonData, err := os.ReadFile(path)
				if err != nil {
					return
				}
				syncedItems.Lock()
				syncedItems.Items = append(syncedItems.Items, jsonData)
				syncedItems.Unlock()

			}()
		}
		return
	})

	if err == nil {
		fileCache.WriteManifestCacheFile(directoryName)
	}
	wg.Wait()
	files = syncedItems.Items

	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//GetCollectionNames returns a name of all collections

func GetCollectionNames() (names []string) {

	names = append(names, "Accounts")
	names = append(names, "Buildings")
	names = append(names, "Countries")
	names = append(names, "FileObjects")
	names = append(names, "Floors")
	names = append(names, "Sites")
	names = append(names, "Users")
	return

}

//GetCollectionHistoryNames returns a name of all collections history tables

func GetCollectionHistoryNames() (names []string) {

	names = append(names, "AccountsHistory")
	names = append(names, "BuildingsHistory")
	names = append(names, "CountriesHistory")
	names = append(names, "FileObjectsHistory")
	names = append(names, "FloorsHistory")
	names = append(names, "SitesHistory")
	names = append(names, "UsersHistory")
	names = append(names, "Transactions")
	return

}

// Each goCore application should probably call this once on server setup to iterate through all records in the system and re-save it so that new fields can be injected into the data and your javascript always will be able to access any record

func UpdateAllRecordsToLatestSchema() {

	var Account []Account
	Accounts.Query().All(&Account)
	for i := range Account {
		Account[i].Save()
	}
	var Building []Building
	Buildings.Query().All(&Building)
	for i := range Building {
		Building[i].Save()
	}
	var Country []Country
	Countries.Query().All(&Country)
	for i := range Country {
		Country[i].Save()
	}
	var FileObject []FileObject
	FileObjects.Query().All(&FileObject)
	for i := range FileObject {
		FileObject[i].Save()
	}
	var Floor []Floor
	Floors.Query().All(&Floor)
	for i := range Floor {
		Floor[i].Save()
	}
	var Site []Site
	Sites.Query().All(&Site)
	for i := range Site {
		Site[i].Save()
	}
	var User []User
	Users.Query().All(&User)
	for i := range User {
		User[i].Save()
	}
}

func ResolveEntity(key string) modelEntity {

	switch key {
	case "Account":
		return &Account{}
	case "AccountHistoryRecord":
		return &AccountHistoryRecord{}
	case "Building":
		return &Building{}
	case "BuildingHistoryRecord":
		return &BuildingHistoryRecord{}
	case "Country":
		return &Country{}
	case "CountryHistoryRecord":
		return &CountryHistoryRecord{}
	case "FileObject":
		return &FileObject{}
	case "FileObjectHistoryRecord":
		return &FileObjectHistoryRecord{}
	case "Floor":
		return &Floor{}
	case "FloorHistoryRecord":
		return &FloorHistoryRecord{}
	case "Site":
		return &Site{}
	case "SiteHistoryRecord":
		return &SiteHistoryRecord{}
	case "User":
		return &User{}
	case "UserHistoryRecord":
		return &UserHistoryRecord{}
	}
	return nil
}

func ResolveField(collectionName string, fieldName string) string {

	switch collectionName + fieldName {
	case "AccountId":
		return "int"
	case "AccountAccountName":
		return "string"
	case "AccountSecondaryPhone":
		return "object"
	case "AccountCreateDate":
		return "dateTime"
	case "AccountUpdateDate":
		return "dateTime"
	case "AccountLastUpdateId":
		return "string"
	case "SecondaryPhoneInfoValue":
		return "string"
	case "SecondaryPhoneInfoNumeric":
		return "string"
	case "SecondaryPhoneInfoDialCode":
		return "string"
	case "SecondaryPhoneInfoCountryISO":
		return "string"
	case "SecondaryPhoneInfoCreateDate":
		return "dateTime"
	case "SecondaryPhoneInfoUpdateDate":
		return "dateTime"
	case "SecondaryPhoneInfoLastUpdateId":
		return "string"
	case "CountryId":
		return "int"
	case "CountryIso":
		return "string"
	case "CountryName":
		return "string"
	case "CountryCreateDate":
		return "dateTime"
	case "CountryUpdateDate":
		return "dateTime"
	case "CountryLastUpdateId":
		return "string"
	case "FileObjectId":
		return "int"
	case "FileObjectName":
		return "string"
	case "FileObjectPath":
		return "string"
	case "FileObjectContent":
		return "string"
	case "FileObjectCreateDate":
		return "dateTime"
	case "FileObjectUpdateDate":
		return "dateTime"
	case "FileObjectLastUpdateId":
		return "string"
	case "BuildingId":
		return "int"
	case "BuildingName":
		return "string"
	case "BuildingImageCustom":
		return "string"
	case "BuildingImageFileName":
		return "string"
	case "BuildingSiteId":
		return "string"
	case "BuildingAccountId":
		return "string"
	case "BuildingCreateDate":
		return "dateTime"
	case "BuildingUpdateDate":
		return "dateTime"
	case "BuildingLastUpdateId":
		return "string"
	case "FloorId":
		return "int"
	case "FloorName":
		return "string"
	case "FloorSiteId":
		return "string"
	case "FloorBuildingId":
		return "string"
	case "FloorAccountId":
		return "string"
	case "FloorCreateDate":
		return "dateTime"
	case "FloorUpdateDate":
		return "dateTime"
	case "FloorLastUpdateId":
		return "string"
	case "SiteId":
		return "int"
	case "SiteName":
		return "string"
	case "SiteImageCustom":
		return "string"
	case "SiteAccountId":
		return "string"
	case "SiteCountryId":
		return "string"
	case "SiteCreateDate":
		return "dateTime"
	case "SiteUpdateDate":
		return "dateTime"
	case "SiteLastUpdateId":
		return "string"
	case "UserId":
		return "int"
	case "UserFirst":
		return "string"
	case "UserLast":
		return "string"
	case "UserSignupDate":
		return "dateTime"
	case "UserEmail":
		return "string"
	case "UserCreateDate":
		return "dateTime"
	case "UserUpdateDate":
		return "dateTime"
	case "UserLastUpdateId":
		return "string"
	}
	return ""
}

func ResolveCollection(key string) (collection, error) {

	if serverSettings.WebConfig.Application.LogJoinQueries {
		fmt.Println(key)
	}
	switch key {
	case "Accounts":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Accounts")
		}
		return &modelAccounts{}, nil
	case "Buildings":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Buildings")
		}
		return &modelBuildings{}, nil
	case "Countries":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Countries")
		}
		return &modelCountries{}, nil
	case "FileObjects":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! FileObjects")
		}
		return &modelFileObjects{}, nil
	case "Floors":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Floors")
		}
		return &modelFloors{}, nil
	case "Sites":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Sites")
		}
		return &modelSites{}, nil
	case "Users":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Users")
		}
		return &modelUsers{}, nil
	}
	return nil, errors.New("Failed to resolve collection:  " + key)
}

func ResolveHistoryCollection(key string) modelCollection {

	switch key {
	case "AccountsHistory":
		return &modelAccountsHistory{}
	case "BuildingsHistory":
		return &modelBuildingsHistory{}
	case "CountriesHistory":
		return &modelCountriesHistory{}
	case "FileObjectsHistory":
		return &modelFileObjectsHistory{}
	case "FloorsHistory":
		return &modelFloorsHistory{}
	case "SitesHistory":
		return &modelSitesHistory{}
	case "UsersHistory":
		return &modelUsersHistory{}
	}
	return nil
}

func joinField(j join, id string, fieldToSet reflect.Value, remainingRecursions string, q *Query, endRecursion bool, recursionCount int) (err error) {

	c, err2 := ResolveCollection(j.collectionName)
	if serverSettings.WebConfig.Application.LogJoinQueries {
		fmt.Println("joinFieldLogging")
		fmt.Println(fmt.Sprintf("%+v", j.collectionName))
		fmt.Println("c")
		fmt.Println(fmt.Sprintf("%+v", c))
		fmt.Println("err2")
		fmt.Println(fmt.Sprintf("%+v", err2))
	}
	if err2 != nil {
		err = errors.New("Failed to resolve collection:  " + j.collectionName)
		return
	}
	switch j.joinSchemaName {
	case "Account":
		var y Account
		if j.isMany {

			obj := fieldToSet.Interface().(*AccountJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Account
			var ji AccountJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Building":
		var y Building
		if j.isMany {

			obj := fieldToSet.Interface().(*BuildingJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Building
			var ji BuildingJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Country":
		var y Country
		if j.isMany {

			obj := fieldToSet.Interface().(*CountryJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Country
			var ji CountryJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "FileObject":
		var y FileObject
		if j.isMany {

			obj := fieldToSet.Interface().(*FileObjectJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []FileObject
			var ji FileObjectJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Floor":
		var y Floor
		if j.isMany {

			obj := fieldToSet.Interface().(*FloorJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Floor
			var ji FloorJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Site":
		var y Site
		if j.isMany {

			obj := fieldToSet.Interface().(*SiteJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Site
			var ji SiteJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "User":
		var y User
		if j.isMany {

			obj := fieldToSet.Interface().(*UserJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []User
			var ji UserJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	}
	err = errors.New("Failed to resolve schema :  " + j.joinSchemaName)
	return
}

const (
	FIELD_COUNTRY_ID                    = "Id"
	FIELD_COUNTRY_ISO                   = "Iso"
	FIELD_COUNTRY_NAME                  = "Name"
	FIELD_FILEOBJECT_ID                 = "Id"
	FIELD_FILEOBJECT_NAME               = "Name"
	FIELD_FILEOBJECT_PATH               = "Path"
	FIELD_FILEOBJECT_CONTENT            = "Content"
	FIELD_BUILDING_ID                   = "Id"
	FIELD_BUILDING_NAME                 = "Name"
	FIELD_BUILDING_IMAGECUSTOM          = "ImageCustom"
	FIELD_BUILDING_IMAGEFILENAME        = "ImageFileName"
	FIELD_BUILDING_SITEID               = "SiteId"
	FIELD_BUILDING_ACCOUNTID            = "AccountId"
	FIELD_FLOOR_ID                      = "Id"
	FIELD_FLOOR_NAME                    = "Name"
	FIELD_FLOOR_SITEID                  = "SiteId"
	FIELD_FLOOR_BUILDINGID              = "BuildingId"
	FIELD_FLOOR_ACCOUNTID               = "AccountId"
	FIELD_SITE_ID                       = "Id"
	FIELD_SITE_NAME                     = "Name"
	FIELD_SITE_IMAGECUSTOM              = "ImageCustom"
	FIELD_SITE_ACCOUNTID                = "AccountId"
	FIELD_SITE_COUNTRYID                = "CountryId"
	FIELD_USER_ID                       = "Id"
	FIELD_USER_FIRST                    = "First"
	FIELD_USER_LAST                     = "Last"
	FIELD_USER_SIGNUPDATE               = "SignupDate"
	FIELD_USER_EMAIL                    = "Email"
	FIELD_ACCOUNT_ID                    = "Id"
	FIELD_ACCOUNT_ACCOUNTNAME           = "AccountName"
	FIELD_ACCOUNT_SECONDARYPHONE        = "SecondaryPhone"
	FIELD_SECONDARYPHONEINFO_VALUE      = "Value"
	FIELD_SECONDARYPHONEINFO_NUMERIC    = "Numeric"
	FIELD_SECONDARYPHONEINFO_DIALCODE   = "DialCode"
	FIELD_SECONDARYPHONEINFO_COUNTRYISO = "CountryISO"
)
