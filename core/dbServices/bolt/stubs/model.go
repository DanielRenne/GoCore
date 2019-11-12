package boltStubs

var Model string

func init() {

	Model = `
package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
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
	DoesIdExist(interface{}) bool // no callers exist yet with this method, but this is stubbed out for future use if we ever wanted to get ById on the interface we could use this or something similar
}

type modelCollection interface {
	Rollback(transactionId string) error
}

type collection interface {
	Query() *Query
}

type BootstrapMeta struct {
	Version      int      ` + "`" + `json:"Version" bson:"Version"` + "`" + `
	Domain       string   ` + "`" + `json:"Domain" bson:"Domain"` + "`" + `
	ReleaseMode  string   ` + "`" + `json:"ReleaseMode" bson:"ReleaseMode"` + "`" + `
	ProductName  string   ` + "`" + `json:"ProductName" bson:"ProductName"` + "`" + `
	Domains      []string ` + "`" + `json:"Domains" bson:"Domains"` + "`" + `
	ProductNames []string ` + "`" + `json:"ProductNames" bson:"ProductNames"` + "`" + `
	DeleteRow    bool     ` + "`" + `json:"DeleteRow" bson:"DeleteRow"` + "`" + `
	AlwaysUpdate bool     ` + "`" + `json:"AlwaysUpdate" bson:"AlwaysUpdate"` + "`" + `
}

type BootstrapSync struct {
	sync.Mutex
	Items [][]byte
}

type tQueue struct {
	sync.RWMutex
	queue map[string]*transactionsToPersist
	ids map[string][]string
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

func init() {
	transactionQueue.ids = make(map[string][]string)
	transactionQueue.queue = make(map[string]*transactionsToPersist)
	go clearTransactionQueue()
}

func Q(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{k: v}
}

func QTs(k string, v time.Time) map[string]time.Time {
	return map[string]time.Time{k: v}
}

func RangeQ(k string, min interface{}, max interface{}) map[string]Range {
	var rge map[string]Range
	rge = make(map[string]Range)
	rge[k] = Range{
		Max: max,
		Min: min,
	}
	return rge
}

func MinQ(k string, min interface{}) map[string]Min {
	var rge map[string]Min
	rge = make(map[string]Min)
	rge[k] = Min{
		Min: min,
	}
	return rge
}

func MaxQ(k string, max interface{}) map[string]Max {
	var rge map[string]Max
	rge = make(map[string]Max)
	rge[k] = Max{
		Max: max,
	}
	return rge
}

//Every 12 hours check the transactionQueue and remove any outstanding stale transactions > 48 hours old
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
	collection := strings.Replace(reflect.TypeOf(x).String() + "s", "model.", "", -1)

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

	var syncedItems BootstrapSync
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
				jsonData, err := ioutil.ReadFile(path)
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
`
}
