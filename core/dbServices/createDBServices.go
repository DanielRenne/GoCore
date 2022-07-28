package dbServices

import (
	"encoding/json"
	"fmt"
	"log"

	boltStubs "github.com/DanielRenne/GoCore/core/dbServices/bolt/stubs"
	commonStubs "github.com/DanielRenne/GoCore/core/dbServices/common/stubs"
	mongoStubs "github.com/DanielRenne/GoCore/core/dbServices/mongo/stubs"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"

	// "fmt"
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/davidrenne/heredoc"
	"github.com/fatih/color"
)

/*
Bolt todo:

Indexes are panicing
Transactions within bold using save with tran
WhiteList/Blacklist seems unsupported completely with storm
And/Or - in query
Ors - in query
Iter - in query
Reflection of In's
Not in... not sure if this is supported

*/

type FieldValidation struct {
	Required  bool   `json:"required"`
	Type      string `json:"type"`
	Min       string `json:"min"`
	Max       string `json:"max"`
	Length    string `json:"length"`
	LengthMax string `json:"lengthMax"`
	LengthMin string `json:"lengthMin"`
}

type join struct {
	CollectionName   string `json:"collectionName"`
	SchemaName       string `json:"schemaName"`
	FieldName        string `json:"fieldName"`
	ForeignFieldName string `json:"foreignFieldName"`
	IsMany           bool   `json:"isMany"`
}

type fieldDef struct {
	Name      string `json:"name"`
	Primary   bool   `json:"primary"`
	AllowNull bool   `json:"allowNull"`
	FieldType string `json:"fieldType"`
	IsUnique  bool   `json:"isUnique"`
	Check     string `json:"check"`
	Collate   string `json:"collate"`
	Default   string `json:"default"`
}

type tableDef struct {
	Name   string     `json:"name"`
	Fields []fieldDef `json:"fields"`
}

type indexDef struct {
	Name      string   `json:"name"`
	TableName string   `json:"tableName"`
	IsUnique  bool     `json:"isUnique"`
	Fields    []string `json:"fields"`
}

type foreignKeyDef struct {
	Name     string   `json:"name"`
	Fields   []string `json:"fields"`
	FKTable  string   `json:"fkTable"`
	FKFields []string `json:"fkFields"`
	OnDelete bool     `json:"onDelete"`
	OnUpdate bool     `json:"onUpdate"`
}

type foreignKeyTableDef struct {
	Table string          `json:"table"`
	Keys  []foreignKeyDef `json:"keys"`
}

type createObject struct {
	Tables      []tableDef           `json:"tables"`
	Indexes     []indexDef           `json:"indexes"`
	ForeignKeys []foreignKeyTableDef `json:"foreignKeys"`
}

type NOSQLSchemaField struct {
	Name         string           `json:"name"`
	Type         string           `json:"type"`
	Index        string           `json:"index"`
	View         bool             `json:"view"`
	Ref          string           `json:"ref"`
	Format       string           `json:"format"`
	OmitEmpty    bool             `json:"omitEmpty"`
	DefaultValue string           `json:"defaultValue"`
	Required     bool             `json:"required"`
	Schema       NOSQLSchema      `json:"schema"`
	Validation   *FieldValidation `json:"validate, omitempty"`
	Join         join             `json:"join"`
	NoPersist    bool             `json:"noPersist"`
}

type NOSQLSchema struct {
	Name   string             `json:"name"`
	Fields []NOSQLSchemaField `json"fields"`
}

type NOSQLCollection struct {
	Name       string      `json:"name"`
	ClearTable bool        `json:"clearTable"`
	Schema     NOSQLSchema `json:"schema"`
	FieldTypes map[string]FieldType
}

type FieldType struct {
	IsArray bool
	Value   string
}

type NOSQLSchemaDB struct {
	Collections []NOSQLCollection `json:"collections"`
}

type fieldType struct {
	Name string
	Type string
}

type entityList struct {
	Constants     []string
	FieldNames    []fieldType
	JoinConstants []string
}

type collectionsSet struct {
	sync.RWMutex
	Collections []NOSQLCollection
	Entities    map[string]entityList
}

type schemasCreatedSync struct {
	sync.RWMutex
	schemasCreated map[string]NOSQLSchema
}

var allCollections collectionsSet

var modelToWrite string

// AxisSorter sorts planets by axis.
type SchemaNameSorter []NOSQLCollection

func (a SchemaNameSorter) Len() int           { return len(a) }
func (a SchemaNameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SchemaNameSorter) Less(i, j int) bool { return a[i].Schema.Name < a[j].Schema.Name }

// This array holds a list of the Schema's created for a model version.
// It is used to NOT duplicate structs for the model.
// When we process each version we clear the array out first, then add and check against it.

func RunDBCreate() {

	// jsonData, err := ioutil.ReadFile("db/" + serverSettings.WebConfig.DbConnection.AppName + "/create.json")
	// if err != nil {
	// 	fmt.Println("Reading of create.json failed:  " + err.Error())
	// 	return
	// }

	// if serverSettings.WebConfig.DbConnection.Driver == "sqlite3" {
	// 	var co createObject
	// 	errUnmarshal := json.Unmarshal(jsonData, &co)
	// 	if errUnmarshal != nil {
	// 		color.Red("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
	// 		return
	// 	}
	// 	if serverSettings.WebConfig.DbConnection.Driver == "sqlite3" {
	// 		createSQLiteTables(co.Tables)
	// 		createSQLiteIndexes(co.Indexes)
	// 		createSQLiteForeignKeys(co.ForeignKeys, co.Tables)
	// 	}

	// } else {
	// var schemaDB NOSQLSchemaDB
	// errUnmarshal := json.Unmarshal(jsonData, &schemaDB)
	// if errUnmarshal != nil {
	// 	color.Red("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
	// 	return
	// }

	// createNoSQLModel(schemaDB.Collections, serverSettings.WebConfig.DbConnection.AppName, serverSettings.WebConfig.DbConnection.Driver)
	// }

	versionDir := "v1"
	modelPath := serverSettings.APP_LOCATION + "/models/" + versionDir + "/model"
	// delete model files
	fmt.Printf("RunDBCreate->Remove Model %v \n", modelPath)
	os.RemoveAll(modelPath)

	walkNoSQLSchema()
}

func walkNoSQLSchema() {
	allCollections.Lock()
	allCollections.Entities = make(map[string]entityList)
	allCollections.Unlock()
	basePath := serverSettings.APP_LOCATION + "/db/schemas"

	fileNames, errReadDir := ioutil.ReadDir(basePath)
	if errReadDir != nil {
		color.Red("Reading of " + basePath + " failed:  " + errReadDir.Error())
		return
	}

	versionNumber := ""
	for _, file := range fileNames {

		if file.IsDir() == true {

			version := extensions.Version{}
			version.Init(file.Name())
			versionDir := "v" + version.MajorString
			versionNumber = version.Value
			walkNoSQLVersion(basePath+"/"+file.Name(), versionDir)

			goFileNames, errGoReadDir := ioutil.ReadDir(serverSettings.APP_LOCATION + "/db/goFiles/v" + version.MajorString)

			if errGoReadDir == nil {
				for _, file := range goFileNames {
					err := extensions.CopyFile(serverSettings.APP_LOCATION+"/db/goFiles/v"+version.MajorString+"/"+file.Name(), serverSettings.APP_LOCATION+"/models/v"+version.MajorString+"/model/"+file.Name())
					if err != nil {
						color.Red(err.Error())
					}
				}
			} else {
				color.Red(errGoReadDir.Error())
			}

			//Create Swagger Definition With the latest Version being equal to swagger.json, all others swagger_1.0.0.json etc...
			//writeSwaggerConfiguration("/api/"+versionDir, version.Value)
		}
	}

	//Make a copy of the latest Swagger Version Definition for Swagger UI to default to swagger.json.  We will keep the latest version as a 2nd copy.
	extensions.CopyFile(serverSettings.SWAGGER_UI_PATH+"/swagger."+versionNumber+".json", serverSettings.SWAGGER_UI_PATH+"/swagger.json")

}

func walkNoSQLVersion(path string, versionDir string) {

	initializeModelFile()

	var scs schemasCreatedSync

	scs.schemasCreated = make(map[string]NOSQLSchema, 0)

	err := filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) error {

		if errWalk != nil {
			return errWalk
		}

		var e error

		if filepath.Ext(f.Name()) == ".json" {
			jsonData, err := ioutil.ReadFile(path)
			if err != nil {
				color.Red("Reading in walkNoSQLVersion of " + path + " failed:  " + err.Error())
				e = err
			}

			var schemaDB NOSQLSchemaDB
			errUnmarshal := json.Unmarshal(jsonData, &schemaDB)
			if errUnmarshal != nil {
				color.Red("Parsing / Unmarshaling of " + path + " failed:  " + errUnmarshal.Error())
				e = errUnmarshal
			}

			for _, col := range schemaDB.Collections {
				allCollections.Lock()
				allCollections.Collections = append(allCollections.Collections, col)
				allCollections.Unlock()
			}

			createNoSQLModel(schemaDB.Collections, serverSettings.WebConfig.DbConnection.Driver, versionDir, &scs)
		}

		return e
	})
	if err != nil {
		color.Red("Walk of path failed:  " + err.Error())
	}

	finalizeModelFile(versionDir)
}

func createNoSQLModel(collections []NOSQLCollection, driver string, versionDir string, scs *schemasCreatedSync) {

	//Clean the Model and API Directory
	// extensions.RemoveDirectory(serverSettings.APP_LOCATION + "/models/" + versionDir + "/model")
	extensions.RemoveDirectory(serverSettings.APP_LOCATION + "/webAPIs/" + versionDir + "/webAPI")

	//Create a NOSQLBucket Model
	// bucket := generateNoSQLModelBucket(driver)
	os.Mkdir(serverSettings.APP_LOCATION+"/models/", 0777)
	os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir, 0777)
	os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/", 0777)

	// writeNOSQLModelBucket(bucket, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/bucket.go")

	//Copy Stub Files
	if driver == DATABASE_DRIVER_MONGODB {
		writeNoSQLStub(mongoStubs.Query, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/query.go")
	} else if driver == DATABASE_DRIVER_BOLTDB {
		writeNoSQLStub(boltStubs.Query, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/query.go")
	}
	writeNoSQLStub(commonStubs.TimeZone, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/timeZone.go")
	writeNoSQLStub(commonStubs.TimeZoneLocations, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/timeZoneLocations.go")
	writeNoSQLStub(commonStubs.Locales, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/locales.go")

	var histTemplate []byte
	if driver == DATABASE_DRIVER_MONGODB {
		var err error
		histTemplate = []byte(mongoStubs.HistTemplate)

		if err != nil {
			color.Red("Error reading histTemplate.go:  " + err.Error())
			return
		}

	}

	var transactionTemplate []byte

	if driver == DATABASE_DRIVER_MONGODB {
		transactionTemplate = []byte(mongoStubs.Transaction)
	} else if driver == DATABASE_DRIVER_BOLTDB {
		transactionTemplate = []byte(boltStubs.Transaction)
	}
	transactionModified := string(transactionTemplate[:])

	if serverSettings.WebConfig.DbConnection.TransactionSizeMax > 0 {
		transactionModified = strings.Replace(transactionModified, "ci := mgo.CollectionInfo{ForceIdIndex: true}", "ci := mgo.CollectionInfo{ForceIdIndex: true, Capped:true, MaxBytes:"+extensions.IntToString(serverSettings.WebConfig.DbConnection.TransactionSizeMax)+"}\n", -1)
	}

	writeNoSQLStub(transactionModified, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/transaction.go")

	//Create the Collection Models
	for _, collection := range collections {
		val := generateNoSQLModel(collection.Schema, collection, driver, scs)
		os.Mkdir(serverSettings.APP_LOCATION+"/models/", 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir, 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/", 0777)
		writeNoSQLModelCollection(val, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+".go", collection)

		if string(histTemplate) != "" {

			//Create the Transaction History Table for the Collection
			histName := strings.Title(collection.Name) + "History"
			histModified := strings.Replace(string(histTemplate[:]), "HistCollection", histName, -1)
			histModified = strings.Replace(histModified, "//CollectionVariable", heredoc.Docf(`
			collection%sMutex.Lock()
			mongo%sCollection = dbServices.MongoDB.C("%s")
			collection%sMutex.Unlock()
			`, histName, histName, histName, histName), -1)
			histModified = strings.Replace(histModified, "HistEntity", strings.Title(collection.Schema.Name)+"HistoryRecord", -1)
			histModified = strings.Replace(histModified, "OriginalEntity", strings.Title(collection.Schema.Name), -1)
			if serverSettings.WebConfig.DbConnection.AuditHistorySizeMax > 0 {
				histModified = strings.Replace(histModified, "ci := mgo.CollectionInfo{ForceIdIndex: true}", "ci := mgo.CollectionInfo{ForceIdIndex: true, Capped:true, MaxBytes:"+extensions.IntToString(serverSettings.WebConfig.DbConnection.AuditHistorySizeMax)+"}\n", -1)
			}

			writeNoSQLStub(histModified, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+"_Hist.go")

		}
		os.Mkdir(serverSettings.APP_LOCATION+"/webAPIs/", 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/webAPIs/"+versionDir, 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/webAPIs/"+versionDir+"/webAPI/", 0777)

		cWebAPI := genSchemaWebAPI(collection, collection.Schema, strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+"/models/"+versionDir+"/model", driver, versionDir)
		writeNoSQLWebAPI(cWebAPI, serverSettings.APP_LOCATION+"/webAPIs/"+versionDir+"/webAPI/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+".go", collection)
	}

}

func initializeModelFile() {
	if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_MONGODB {
		modelToWrite = mongoStubs.Model
	} else {
		modelToWrite = boltStubs.Model
	}

	modelToWrite += "\n"

}

func finalizeModelFile(versionDir string) {
	allCollections.RLock()
	sort.Sort(SchemaNameSorter(allCollections.Collections))
	allCollections.RUnlock()

	modelToWrite += "//GetCollectionNames returns a name of all collections\n\n"
	modelToWrite += "func GetCollectionNames() (names []string) {\n\n"
	for _, collection := range allCollections.Collections {
		modelToWrite += "names = append(names,\"" + collection.Name + "\")\n"
	}
	modelToWrite += "return\n\n"
	modelToWrite += "}\n\n"

	modelToWrite += "//GetCollectionHistoryNames returns a name of all collections history tables\n\n"
	modelToWrite += "func GetCollectionHistoryNames() (names []string) {\n\n"
	for _, collection := range allCollections.Collections {
		modelToWrite += "names = append(names,\"" + collection.Name + "History\")\n"
	}
	modelToWrite += "names = append(names,\"Transactions\")\n"
	modelToWrite += "return\n\n"
	modelToWrite += "}\n\n"

	modelToWrite += "// Each goCore application should probably call this once on server setup to iterate through all records in the system and re-save it so that new fields can be injected into the data and your javascript always will be able to access any record\n\n"
	modelToWrite += "func UpdateAllRecordsToLatestSchema() {\n\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "var " + collection.Schema.Name + " []" + strings.Title(collection.Schema.Name) + "\n"
		modelToWrite += strings.Title(collection.Name) + ".Query().All(& " + collection.Schema.Name + ")\n"
		modelToWrite += "for i := range " + collection.Schema.Name + " {\n"
		modelToWrite += collection.Schema.Name + "[i].Save()\n"
		modelToWrite += "}\n"
	}

	modelToWrite += "}\n\n"

	modelToWrite += "func ResolveEntity(key string) modelEntity{\n\n"
	modelToWrite += "switch key{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "\":\n"
		modelToWrite += " return &" + strings.Title(collection.Schema.Name) + "{}\n"

		if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_MONGODB {
			modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "HistoryRecord\":\n"
			modelToWrite += " return &" + strings.Title(collection.Schema.Name) + "HistoryRecord{}\n"
		}
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"

	modelToWrite += "func ResolveField(collectionName string, fieldName string) string{\n\n"
	modelToWrite += "switch collectionName + fieldName {\n"

	for key, value := range allCollections.Entities {

		var fieldNameLookup sync.Map

		for _, constValue := range value.FieldNames {

			_, found := fieldNameLookup.Load(key + constValue.Name)
			if found {
				continue
			} else {
				fieldNameLookup.Store(key+constValue.Name, nil)
			}

			modelToWrite += "case \"" + key + constValue.Name + "\":\n"
			modelToWrite += "return \"" + constValue.Type + "\"\n"
		}
		modelToWrite += "case \"" + key + "CreateDate\":\n"
		modelToWrite += "return \"dateTime\"\n"
		modelToWrite += "case \"" + key + "UpdateDate\":\n"
		modelToWrite += "return \"dateTime\"\n"
		modelToWrite += "case \"" + key + "LastUpdateId\":\n"
		modelToWrite += "return \"string\"\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return \"\"\n"
	modelToWrite += "}\n\n"

	modelToWrite += "\n"
	modelToWrite += "func ResolveCollection(key string) (collection, error){\n\n"
	modelToWrite += " if serverSettings.WebConfig.Application.LogJoinQueries {\n"
	modelToWrite += "fmt.Println(key)\n"
	modelToWrite += " }\n"

	modelToWrite += "switch key{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Name) + "\":\n"
		modelToWrite += " if serverSettings.WebConfig.Application.LogJoinQueries {\n"
		modelToWrite += " fmt.Println(\"in case!! " + strings.Title(collection.Name) + "\")\n"
		modelToWrite += " }\n"
		modelToWrite += " return &model" + strings.Title(collection.Name) + "{}, nil\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil, errors.New(\"Failed to resolve collection:  \" + key)\n"
	modelToWrite += "}\n\n"

	modelToWrite += "\n"
	modelToWrite += "func ResolveHistoryCollection(key string) modelCollection{\n\n"

	if serverSettings.WebConfig.DbConnection.Driver == DATABASE_DRIVER_MONGODB {
		modelToWrite += "switch key{\n"

		for _, collection := range allCollections.Collections {
			modelToWrite += "case \"" + strings.Title(collection.Name) + "History\":\n"
			modelToWrite += " return &model" + strings.Title(collection.Name) + "History{}\n"
		}
		modelToWrite += "}\n"
	}
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"
	modelToWrite += "func joinField(j join, id string, fieldToSet reflect.Value, remainingRecursions string, q *Query, endRecursion bool, recursionCount int) (err error) {\n\n"
	modelToWrite += "c, err2 := ResolveCollection(j.collectionName)\n"
	modelToWrite += " if serverSettings.WebConfig.Application.LogJoinQueries {\n"
	modelToWrite += "fmt.Println(\"joinFieldLogging\")\n"
	modelToWrite += "fmt.Println(fmt.Sprintf(\"%+v\", j.collectionName))\n"
	modelToWrite += "fmt.Println(\"c\")\n"
	modelToWrite += "fmt.Println(fmt.Sprintf(\"%+v\", c))\n"
	modelToWrite += "fmt.Println(\"err2\")\n"
	modelToWrite += "fmt.Println(fmt.Sprintf(\"%+v\", err2))\n"
	modelToWrite += "}\n"

	modelToWrite += "if err2 != nil {\n"
	modelToWrite += "	err = errors.New(\"Failed to resolve collection:  \" + j.collectionName)\n"
	modelToWrite += "	return\n"
	modelToWrite += "}\n"

	modelToWrite += "switch j.joinSchemaName{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "\":\n"
		modelToWrite += "var y " + strings.Title(collection.Schema.Name) + "\n"
		modelToWrite += "if j.isMany {\n\n"

		modelToWrite += "obj := fieldToSet.Interface().(*" + strings.Title(collection.Schema.Name) + "JoinItems)\n"
		modelToWrite += "if obj != nil {\n"
		modelToWrite += "for i, _ := range *obj.Items {\n"
		modelToWrite += "item := &(*obj.Items)[i]\n"
		modelToWrite += "item.JoinFields(remainingRecursions, q, recursionCount)\n"
		modelToWrite += "}\n"
		modelToWrite += "return\n"
		modelToWrite += "}\n\n"

		modelToWrite += "var z []" + strings.Title(collection.Schema.Name) + "\n"
		modelToWrite += "var ji " + strings.Title(collection.Schema.Name) + "JoinItems\n"
		modelToWrite += "fieldToSet.Set(reflect.ValueOf(&ji))\n"
		modelToWrite += "JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)\n"
		modelToWrite += "}else{\n"
		modelToWrite += "JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)\n"
		modelToWrite += "}\n"
		modelToWrite += "return\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "err = errors.New(\"Failed to resolve schema :  \" + j.joinSchemaName)\n"
	modelToWrite += "return \n"
	modelToWrite += "}\n\n"

	modelToWrite += "const (\n"

	for key, value := range allCollections.Entities {

		for _, constValue := range value.Constants {
			modelToWrite += "FIELD_" + strings.ToUpper(key) + "_" + strings.ToUpper(constValue) + " = \"" + constValue + "\"\n"
		}
	}

	modelToWrite += ")\n\n"

	writeNoSQLStub(modelToWrite, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/model.go")
}

func generateNoSQLModel(schema NOSQLSchema, collection NOSQLCollection, driver string, scs *schemasCreatedSync) string {

	val := ""
	collection.FieldTypes = make(map[string]FieldType)

	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/logger", "github.com/DanielRenne/GoCore/core/pubsub", "github.com/DanielRenne/GoCore/core/serverSettings", "github.com/DanielRenne/GoCore/core/dbServices", "github.com/globalsign/mgo/bson", "encoding/json", "errors", "time", "github.com/asdine/storm", "reflect", "sync", "log", "encoding/base64", "github.com/DanielRenne/GoCore/core/utils", "fmt", "github.com/DanielRenne/GoCore/core/fileCache", "github.com/DanielRenne/GoCore/core", "encoding/hex", "github.com/DanielRenne/GoCore/core/store", "crypto/md5"})
	case DATABASE_DRIVER_MONGODB:
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "github.com/DanielRenne/GoCore/core/pubsub", "github.com/DanielRenne/GoCore/core/serverSettings", "encoding/json", "github.com/globalsign/mgo", "github.com/globalsign/mgo/bson", "log", "time", "errors", "encoding/base64", "reflect", "github.com/DanielRenne/GoCore/core/utils", "fmt", "github.com/DanielRenne/GoCore/core/logger", "github.com/DanielRenne/GoCore/core", "github.com/DanielRenne/GoCore/core/fileCache", "github.com/DanielRenne/GoCore/core/store", "github.com/DanielRenne/GoCore/core/atomicTypes", "crypto/md5", "encoding/hex", "sync"})
		// val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "gopkg.in/mgo.v2/bson", "log", "time"})
	}

	val += genNoSQLCollection(collection, schema, driver)
	val += genNoSQLSchema(&collection, schema, driver, scs, 0)
	val += genNoSQLRuntime(collection, schema, driver)

	return val
}

func generateNoSQLModelBucket(driver string) string {
	val := ""
	switch driver {
	case DATABASE_DRIVER_BOLTDB, DATABASE_DRIVER_MONGODB:
		// val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices"})
		val += extensions.GenPackageImport("model", []string{""})
	}

	val += genNoSQLBucketCore(driver)
	return val
}

func writeNoSQLWebAPI(value string, path string, collection NOSQLCollection) {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		color.Red("Error creating Web API for Collection " + collection.Name + ":  " + err.Error())
		return
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		color.Red("Failed to gofmt on file " + path + ":  " + err.Error())
		return
	}
	color.Green("Created NOSQL Web API Collection " + collection.Name + " successfully.")
}

func copyNoSQLStub(source string, dest string) {

	err := extensions.CopyFile(source, dest)
	if err != nil {
		color.Red("Failed to copy stub file from " + source + " to " + dest + ":  " + err.Error())
		return
	}

	color.Green("Successfully Copied Stub file to " + dest + ".")
}

func writeNoSQLStub(value string, path string) {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		color.Red("Failed to write stub file to " + path + ":  " + err.Error())
		return
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		color.Red("Failed to gofmt on file " + path + ":  " + err.Error())
		return
	}

	color.Green("Successfully Wrote Stub file to " + path + ".")
}

func writeNoSQLModelCollection(value string, path string, collection NOSQLCollection) {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		color.Red("Error creating Model for Collection " + collection.Name + ":  " + err.Error())
		return
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		color.Red("Failed to gofmt on file " + path + ":  " + err.Error())
		return
	}
	color.Green("Created NOSQL Model Collection " + collection.Name + " successfully.")
}

func writeNOSQLModelBucket(value string, path string) {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		color.Red("Error creating Model for Bucket:  " + err.Error())
		return
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		color.Red("Failed to gofmt on file " + path + ":  " + err.Error())
		return
	}

	color.Green("Created NOSQL Bucket successfully.")
}

func genNoSQLCollection(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "var " + strings.Title(collection.Name) + " model" + strings.Title(collection.Name) + "\n\n"
	val += "type model" + strings.Title(collection.Name) + " struct{}\n\n"
	val += "var collection" + strings.Title(collection.Name) + "Mutex *sync.RWMutex\n\n"
	val += "type " + strings.Title(schema.Name) + "JoinItems struct{\n"
	val += " Count int `json:\"Count\"`\n"
	val += " Items *[]" + strings.Title(schema.Name) + " `json:\"Items\"`\n"
	val += "}\n\n"

	val += "var GoCore" + strings.Title(collection.Name) + "HasBootStrapped atomicTypes.AtomicBool\n\n"
	if driver == DATABASE_DRIVER_MONGODB {

		val += "var mongo" + strings.Title(collection.Name) + "Collection *mgo.Collection\n"
		val += "func init(){\n"
		val += "store.RegisterStore(" + strings.Title(collection.Name) + ")\n"
		val += "collection" + strings.Title(collection.Name) + "Mutex = &sync.RWMutex{}\n"
		// val += "go func() {\n\n"
		// val += "for{\n"
		// val += "mdb := dbServices.ReadMongoDB()\n"
		// val += "if mdb != nil {\n"
		// val += "init" + strings.Title(collection.Name) + "()\n"
		// val += "return\n"
		// val += "}\n"
		// val += "time.Sleep(time.Millisecond * 5)\n"
		// val += "}\n"
		// val += "}()\n"
		val += "}\n\n"

		// val += "func init" + strings.Title(collection.Name) + "(){\n"
		// val += "log.Println(\"Building Indexes for MongoDB collection " + collection.Name + ":\")\n"
		// val += "mdb := dbServices.ReadMongoDB()\n"

		// val += "collection" + strings.Title(collection.Name) + "Mutex.Lock()\n"
		// val += "mongo" + strings.Title(collection.Name) + "Collection = mdb.C(\"" + collection.Name + "\")\n"
		// val += "collection" + strings.Title(collection.Name) + "Mutex.Unlock()\n"

		// val += "ci := mgo.CollectionInfo{ForceIdIndex: true}\n"
		// val += "collection" + strings.Title(collection.Name) + "Mutex.RLock()\n"
		// val += "mongo" + strings.Title(collection.Name) + "Collection.Create(&ci)\n"
		// val += "collection" + strings.Title(collection.Name) + "Mutex.RUnlock()\n"
		// val += strings.Title(collection.Name) + ".Index()\n"
		// val += strings.Title(collection.Name) + ".Bootstrap()\n"
		// val += "}\n\n"
	} else if driver == DATABASE_DRIVER_BOLTDB {
		val += "func init(){\n"
		val += "collection" + strings.Title(collection.Name) + "Mutex = &sync.RWMutex{}\n\n"
		val += "//" + strings.Title(collection.Name) + ".Index()\n"
		val += "go func(){ time.Sleep(time.Second * 5)\n" + strings.Title(collection.Name) + ".Bootstrap()}()\n"
		val += "store.RegisterStore(" + strings.Title(collection.Name) + ")\n"
		val += "}\n\n"
	}

	val += "func (self *" + strings.Title(schema.Name) + ") GetId() string { \n"
	val += "return self.Id.Hex()\n"
	val += "}\n\n"

	return val
}

func checkSchemaForDateTime(schema NOSQLSchema) bool {

	for _, field := range schema.Fields {

		if field.Type == "dateTime" {
			return true
		}

		if field.Type == "object" || field.Type == "objectArray" {
			objContainsDateTime := checkSchemaForDateTime(field.Schema)
			if objContainsDateTime {
				return true
			}
		}
	}

	return false
}

// Recursive function that will recurse type object and objectArray's in order to create structs for the model.
func genNoSQLSchema(collection *NOSQLCollection, schema NOSQLSchema, driver string, scs *schemasCreatedSync, seed int) string {

	allCollections.Lock()
	var el entityList
	allCollections.Entities[strings.Title(schema.Name)] = el
	allCollections.Unlock()

	schemasToCreate := []NOSQLSchema{}
	val := ""
	hasViews := false
	hasJoins := false
	prefix := ""

	if seed > 0 {
		prefix = strings.Title(collection.Name)
	}

	scs.Lock()
	_, ok := scs.schemasCreated[prefix+schema.Name]
	if ok {
		color.Cyan("%v", "Duplicate Entity:\n")
		log.Println("Skipping duplicate schema:  " + prefix + schema.Name)
		scs.Unlock()
		return ""
	}

	scs.schemasCreated[prefix+schema.Name] = schema
	scs.Unlock()

	val += "type " + prefix + strings.Title(schema.Name) + " struct{\n"

	for _, field := range schema.Fields {

		fieldPrefix := ""
		if field.View {
			hasViews = true
			continue
		}
		if field.Type == "join" {
			hasJoins = true
			continue
		}

		if field.Type == "object" || field.Type == "objectArray" {
			schemasToCreate = append(schemasToCreate, field.Schema)
			fieldPrefix = strings.Title(collection.Name)
		}

		additionalTags := genNoSQLAdditionalTags(field, driver)
		omitEmpty := ""
		if field.OmitEmpty {
			omitEmpty = ",omitempty"
		}

		allCollections.Lock()
		entityConsts := allCollections.Entities[strings.Title(schema.Name)]
		insert := true
		for _, con := range entityConsts.Constants {
			if strings.Title(field.Name) == con {
				insert = false
			}
		}
		if insert {
			entityConsts.Constants = append(entityConsts.Constants, strings.Title(field.Name))
		}
		entityConsts.FieldNames = append(entityConsts.FieldNames, fieldType{
			Name: strings.Title(field.Name),
			Type: field.Type,
		})
		allCollections.Entities[strings.Title(schema.Name)] = entityConsts
		allCollections.Unlock()

		ft := FieldType{}
		fieldStringType := genNoSQLFieldType(fieldPrefix, schema, field, driver)
		if strings.Contains(fieldStringType, "[]") && !strings.Contains(fieldStringType, "[]byte") {
			ft.IsArray = true
		}
		ft.Value = fieldStringType

		jsonField := strings.Title(field.Name)
		if field.NoPersist == true && driver == DATABASE_DRIVER_BOLTDB {
			jsonField = "-"
		}
		collection.FieldTypes[strings.Replace(strings.Title(field.Name), " ", "_", -1)] = ft
		val += "\n\t" + strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t" + fieldStringType + "\t\t`json:\"" + jsonField + omitEmpty + "\"" + additionalTags + "`"
	}

	if seed == 0 {
		val += "\n\t CreateDate time.Time `json:\"CreateDate\" bson:\"CreateDate\"`"
		val += "\n\t UpdateDate time.Time `json:\"UpdateDate\" bson:\"UpdateDate\"`"
		val += "\n\t LastUpdateId string `json:\"LastUpdateId\" bson:\"LastUpdateId\"`"
	}

	//Add Validation
	if seed == 0 {
		val += "\n"
		val += "BootstrapMeta          *BootstrapMeta             `json:\"BootstrapMeta\" bson:\"-\"`\n\n"
		val += "\n"
		val += "Errors struct{\n"

		for _, field := range schema.Fields {

			if field.View || field.Type == "join" {
				continue
			}

			if field.Type == "object" {
				val += strings.Title(field.Name) + " struct{\n"
				val += genNoSQLValidationRecusion(field.Schema)
				val += "} `json:\"" + strings.Title(field.Name) + "\"`\n"
				continue
			}
			if field.Type == "objectArray" {
				val += strings.Title(field.Name) + " []struct{\n"
				val += genNoSQLValidationRecusion(field.Schema)
				val += "} `json:\"" + strings.Title(field.Name) + "\"`\n"
				continue
			}

			val += strings.Title(field.Name) + " string `json:\"" + strings.Title(field.Name) + "\"`\n"
		}

		val += "} `json:\"Errors\" bson:\"-\"`\n\n"
	}

	//Add Views
	if hasViews {
		val += "\n"
		val += "Views struct{\n"

		for _, field := range schema.Fields {
			if field.View {

				viewTags := ""
				viewTagSpace := ""

				if field.Ref != "" {
					viewTagSpace = " "
					viewTags += "ref:\""
					viewTags += strings.Title(field.Ref)
					if field.Format != "" {
						viewTags += "~" + field.Format
					}
					viewTags += "\""
				}

				val += strings.Title(field.Name) + " " + genNoSQLFieldType("", schema, field, driver) + " `json:\"" + strings.Title(field.Name) + "\"" + viewTagSpace + viewTags + "`\n"
			}
		}

		val += "} `json:\"Views\" bson:\"-\"`\n\n"
	}

	//Add Joins
	if hasJoins {
		val += "\n"
		val += "Joins struct{\n"
		for _, field := range schema.Fields {
			if field.Type == "join" {
				schemaName := field.Join.SchemaName

				if field.Join.IsMany {
					schemaName = field.Join.SchemaName + "JoinItems"
				}
				val += strings.Title(field.Name) + " *" + schemaName + " `json:\"" + strings.Title(field.Name) + ",omitempty\" join:\"" + strings.Title(field.Join.CollectionName) + "," +
					strings.Title(field.Join.SchemaName) + "," +
					strings.Title(field.Join.FieldName) + "," +
					extensions.BoolToString(field.Join.IsMany) + "," +
					strings.Title(field.Join.ForeignFieldName) +
					"\"`\n"

				allCollections.Lock()
				entityConsts := allCollections.Entities[strings.Title(schema.Name)]
				entityConsts.JoinConstants = append(entityConsts.JoinConstants, strings.Title(field.Name))
				allCollections.Entities[strings.Title(schema.Name)] = entityConsts
				allCollections.Unlock()

			}
		}

		val += "} `json:\"Joins\" bson:\"-\"`\n\n"
	}

	// val += "\n\t sync.RWMutex `bson:\"-\"`\n"

	val += "\n}\n\n"

	for _, schemaToCreate := range schemasToCreate {
		val += genNoSQLSchema(collection, schemaToCreate, driver, scs, 1)
	}

	return val
}

//Recursive Function to Generate Validation Schema
func genNoSQLValidationRecusion(schema NOSQLSchema) string {
	val := ""
	for _, field := range schema.Fields {
		if field.Type == "object" {
			val += strings.Title(field.Name) + " struct{\n"
			val += genNoSQLValidationRecusion(field.Schema)
			val += "} `json:\"" + strings.Title(field.Name) + "\"`\n"
			continue
		}
		if field.Type == "objectArray" {
			val += strings.Title(field.Name) + " []struct{\n"
			val += genNoSQLValidationRecusion(field.Schema)
			val += "} `json:\"" + strings.Title(field.Name) + "\"`\n"
			continue
		}

		val += strings.Title(field.Name) + " string `json:\"" + strings.Title(field.Name) + "\"`\n"
	}
	return val
}

func genNoSQLAdditionalTags(field NOSQLSchemaField, driver string) string {

	validationTags := ""
	validationTagsGap := ""

	if field.Validation != nil {
		validationTags = "validate:\""
		validationTagsGap = " "

		validationTags += extensions.BoolToString(field.Validation.Required) + ","
		validationTags += field.Validation.Type + ","
		validationTags += field.Validation.Min + ","
		validationTags += field.Validation.Max + ","
		validationTags += field.Validation.Length + ","
		validationTags += field.Validation.LengthMax + ","
		validationTags += field.Validation.LengthMin + "\""
	}

	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		tags := validationTagsGap + validationTags
		switch field.Index {
		case "":
			return tags
		case "primary":
			return " storm:\"id\""
		case "index":
			return " storm:\"index\"" + tags
		case "unique":
			return " storm:\"unique\"" + tags
		}
	case DATABASE_DRIVER_MONGODB:

		tags := " bson:\"" + strings.Title(field.Name) + "\"" + validationTagsGap + validationTags
		if field.NoPersist == true {
			tags = " bson:\"-\"" + validationTagsGap + validationTags
		}
		switch field.Index {
		case "":
			return tags
		case "primary":
			return " bson:\"_id,omitempty\""
		case "index":
			return " dbIndex:\"index\"" + tags
		case "unique":
			return " dbIndex:\"unique\"" + tags
		}
	}
	return ""
}

func genNoSQLFieldType(prefix string, schema NOSQLSchema, field NOSQLSchemaField, driver string) string {

	if field.Index == "primary" {
		return "bson.ObjectId"
	}

	if field.Index == "primary" {
		if field.Type == "string" {
			return "string"
		} else {
			return "uint64"
		}
	}

	switch field.Type {
	case "dateTime":
		return "time.Time"
	case "interface":
		return "interface{}"
	case "interfaceArray":
		return "[]interface{}"
	case "byteArray":
		return "[]byte"
	case "object":
		return prefix + strings.Title(field.Schema.Name)
	case "intArray":
		return "[]int"
	case "float64":
		return "float64"
	case "float64Array":
		return "[]float64"
	case "stringArray":
		return "[]string"
	case "boolArray":
		return "[]bool"
	case "objectArray":
		return "[]" + prefix + strings.Title(field.Schema.Name)
	case "selfArray":
		return "[]" + prefix + strings.Title(schema.Name)
	case "self":
		return prefix + strings.Title(schema.Name)
	}

	return field.Type
}

func genNoSQLRuntime(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	if driver == DATABASE_DRIVER_BOLTDB {
		val += genNoSQLSchemaSingle(collection, schema, driver)
		val += genNoSQLSchemaSearch(collection, schema, driver)
		val += genNoSQLSchemaAll(collection, schema, driver)
		val += genNoSQLSchemaAllByIndex(collection, schema, driver)
		val += genNoSQLSchemaRange(collection, schema, driver)
	}

	val += genSetCollection(collection)

	val += genById(collection, schema, driver)
	val += genDoesIdExist(collection, schema, driver)
	val += genNewByReflection(collection, schema, driver)
	val += genByFilter(collection, schema, driver)
	val += genCountByFilter(collection, schema, driver)
	val += genNOSQLQuery(collection, schema, driver)
	val += genNOSQLRemoveAll(collection, schema, driver)
	val += genNoSQLSchemaIndex(collection, schema, driver)
	val += genNoSQLBootstrap(collection, schema, driver)
	// val += genNoSQLSchemaRunTransaction(collection, schema, driver)
	val += genNoSQLSchemaNew(collection, schema)
	val += genNewId(collection, schema, driver)
	val += genNoSQLSchemaSave(collection, schema, driver)
	val += genNoSQLSchemaSaveByTran(collection, schema, driver)
	val += genNoSQLValidate(collection, schema, driver)
	val += genNoSQLReflect(collection, schema, driver)
	val += genNoSQLSchemaDelete(collection, schema, driver)
	val += genNoSQLSchemaDeleteWithTran(collection, schema, driver)
	val += genNoSQLSchemaJoinFields(collection, schema, driver)
	val += genNoSQLUnMarshal(collection, schema, driver)
	val += genNoSQLSchemaJSONRuntime(schema)
	val += genParseInterface(collection, schema, driver)
	val += genNoSQLSchemaReflectByFieldName(collection)
	val += genNoSQLSchemaReflectBaseTypeByFieldName(collection)

	return val
}

func genNOSQLRemoveAll(collection NOSQLCollection, schema NOSQLSchema, driver string) string {

	val := ""
	//heredocs are concatenated like this because percent is a special thing and the modulus it will think its a tag

	if driver == DATABASE_DRIVER_MONGODB {
		val = heredoc.Docf(`
			func (obj model%s) RemoveAll() {
				var elapseMs int
				collection := mongo%sCollection
				for {
					bootstrapped := GoCore%sHasBootStrapped.Get()

					if collection != nil && bootstrapped {
						break
					}
					elapseMs = elapseMs + 2
					time.Sleep(time.Millisecond * 2)
					if elapseMs `, strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name)) + "%" + heredoc.Docf(`10000 == 0 {
						log.Println("%s has not bootstrapped and has yet to get a collection pointer")
					}
				}
				collection.RemoveAll(bson.M{})
				return
			}
		`, strings.Title(collection.Name))
	} else if driver == DATABASE_DRIVER_BOLTDB {
		val = `
			func (obj model` + strings.Title(collection.Name) + `) RemoveAll() {
				x, errDelete := obj.All()
				if errDelete == nil {
					for _, row := range x {
						row.Delete()
					}
				}
				return
			}
		`
	}
	return val
}

func genNOSQLQuery(collection NOSQLCollection, schema NOSQLSchema, driver string) string {

	val := ""
	//heredocs are concatenated like this because percent is a special thing and the modulus it will think its a tag

	if driver == DATABASE_DRIVER_MONGODB {
		val = heredoc.Docf(`
			func (obj model%s) Query() *Query {
				query := new(Query)
				var elapseMs int
				for {
					collection%sMutex.RLock()
					collection := mongo%sCollection
					bootstrapped := GoCore%sHasBootStrapped.Get()
					collection%sMutex.RUnlock()

					if collection != nil && bootstrapped {
						break
					}
					elapseMs = elapseMs + 2
					time.Sleep(time.Millisecond * 2)
					if elapseMs `, strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name)) + "%" + heredoc.Docf(`10000 == 0 {
						log.Println("%s has not bootstrapped and has yet to get a collection pointer")
					}
				}
				collection%sMutex.RLock()
				collection := mongo%sCollection
				collection%sMutex.RUnlock()
				query.collection = collection
				query.entityName = "%s"
				return query
			}
		`, strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(schema.Name))
	} else if driver == DATABASE_DRIVER_BOLTDB {
		val = heredoc.Docf(`
			func (obj model%s) Query() *Query {
				query := new(Query)
				var elapseMs int
				for {

					bootstrapped := GoCore%sHasBootStrapped.Get()
					if bootstrapped {
						break
					}
					elapseMs = elapseMs + 2
					time.Sleep(time.Millisecond * 2)
					if elapseMs `, strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name)) + "%" + heredoc.Docf(`10000 == 0 {
						log.Println("%s has not bootstrapped and has yet to get a collection pointer")
					}
				}
				query.collectionName = "%s"
				query.entityName = "%s"
				return query
			}
		`, strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(schema.Name))
	}
	return val
}

func genNoSQLSchemaJSONRuntime(schema NOSQLSchema) string {
	val := ""

	val += "func (obj *" + strings.Title(schema.Name) + ") JSONString() (string, error) {\n"
	val += "bytes, err := json.Marshal(obj)\n"
	val += "return string(bytes), err\n"
	val += "}\n\n"

	val += "func (obj *" + strings.Title(schema.Name) + ") JSONBytes() ([]byte, error) {\n"
	val += "	return json.Marshal(obj)\n"
	val += "}\n\n"

	val += "func (obj *" + strings.Title(schema.Name) + ") BSONString() (string, error) {\n"
	val += "bytes, err := bson.Marshal(obj)\n"
	val += "return string(bytes), err\n"
	val += "}\n\n"

	val += "func (obj *" + strings.Title(schema.Name) + ") BSONBytes() (in []byte, err error) {\n"
	val += "err = bson.Unmarshal(in, obj)\n"
	val += "return\n"
	val += "}\n\n"
	return val
}

func genNoSQLSchemaReflectByFieldName(collection NOSQLCollection) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error){\n\n"
	val += "switch fieldName{\n"

	for key, value := range collection.FieldTypes {

		valueType := strings.Replace(value.Value, "[]", "", -1)
		marshalJSON := true
		if valueType == "string" ||
			valueType == "bool" ||
			valueType == "int" ||
			valueType == "float64" ||
			valueType == "time.Time" ||
			valueType == "bson.ObjectId" {
			marshalJSON = false
		}

		val += "\tcase \"" + key + "\":\n"

		if valueType == "interface{}" {
			val += "\tvalue = reflect.ValueOf(x)\n"
			val += "return\n"
		} else {
			if value.IsArray {
				val += "xArray, ok := x.([]interface{})\n\n"

				val += "if ok {"
				val += "arrayToSet := make(" + value.Value + ", len(xArray))\n"
				val += "for i := range xArray {\n"
				val += "	inf := xArray[i]\n"
				if marshalJSON {
					val += "data, _ := json.Marshal(inf)\n"
					val += "var obj " + valueType + "\n"
					val += "err = json.Unmarshal(data, &obj)\n"
					val += "if err != nil {\n"
					val += "return\n"
					val += "}\n"
					val += "	arrayToSet[i] = obj\n"
				} else {
					val += "var ok bool\n"
					val += "	arrayToSet[i], ok = inf.(" + valueType + ")\n"
					val += "if !ok {\n"
					val += "err = errors.New(\"Failed to typecast interface.\")\n"
					val += "return\n"
					val += "}\n"
				}

				val += "}\n\n"

				val += "value = reflect.ValueOf(arrayToSet)\n"
				val += "}else {\n"
				val += "data, _ := json.Marshal(x)\n"
				val += "var obj []" + valueType + "\n"
				val += "err = json.Unmarshal(data, &obj)\n"
				val += "if err != nil {\n"
				val += "return\n"
				val += "}\n"
				val += "\tvalue = reflect.ValueOf(obj)\n"
				val += "}\n\n"

			} else {
				if marshalJSON {
					val += "data, _ := json.Marshal(x)\n"
					val += "var obj " + valueType + "\n"
					val += "err = json.Unmarshal(data, &obj)\n"
					val += "if err != nil {\n"
					val += "return\n"
					val += "}\n"
				} else {
					val += "\tobj, ok := x.(" + value.Value + ")\n"
					val += "if !ok {\n"
					val += "err = errors.New(\"Failed to typecast interface.\")\n"
					val += "return\n"
					val += "}\n"
				}

				val += "\tvalue = reflect.ValueOf(obj)\n"
				val += "return\n"
			}
		}
	}

	val += "}\n"
	val += "return\n"
	val += "}\n\n"

	return val
}

func genNoSQLSchemaReflectBaseTypeByFieldName(collection NOSQLCollection) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error){\n\n"
	val += "switch fieldName{\n"

	for key, value := range collection.FieldTypes {

		valueType := strings.Replace(value.Value, "[]", "", -1)

		marshalJSON := true
		if valueType == "string" ||
			valueType == "bool" ||
			valueType == "int" ||
			valueType == "float64" ||
			valueType == "time.Time" ||
			valueType == "bson.ObjectId" {
			marshalJSON = false
		}

		val += "\tcase \"" + key + "\":\n"

		if valueType == "interface{}" {
			val += "\tvalue = reflect.ValueOf(x)\n"
			val += "return\n"
		} else {
			if marshalJSON {
				val += "if x == nil{\n"
				val += "obj := " + valueType + "{}\n"
				val += "\tvalue = reflect.ValueOf(&obj)\n"
				val += "return\n"
				val += "}\n\n"

				val += "data, _ := json.Marshal(x)\n"
				val += "var obj " + valueType + "\n"
				val += "err = json.Unmarshal(data, &obj)\n"
				val += "if err != nil {\n"
				val += "return\n"
				val += "}\n"
			} else {
				val += "if x == nil{\n"
				val += "var obj " + valueType + "\n"
				val += "\tvalue = reflect.ValueOf(obj)\n"
				val += "return\n"
				val += "}\n\n"

				val += "\tobj, ok := x.(" + valueType + ")\n"
				val += "if !ok {\n"
				val += "err = errors.New(\"Failed to typecast interface.\")\n"
				val += "return\n"
				val += "}\n"
			}
			val += "\tvalue = reflect.ValueOf(obj)\n"
			val += "return\n"
		}
	}

	val += "}\n"
	val += "return\n"
	val += "}\n\n"

	return val
}

func genNoSQLSchemaSave(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") Save() error {\n"
	val += "if !AllowWrites {\n"
	val += "	return nil\n"
	val += "}\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "t := time.Now()\n"
		val += "if self.Id == \"\" {\n"
		val += "self.Id = bson.NewObjectId()\n"
		val += "self.CreateDate = t\n"
		val += "}\n"
		val += "self.UpdateDate = t \n"
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "err := dbServices.BoltDB.Save(self)\n"
		val += "if err == nil{\n"
		val += "pubsub.Publish(\"" + strings.Title(collection.Name) + ".Save\", self)\n"
		val += "}\n"
		val += "return nil\n"
	case DATABASE_DRIVER_MONGODB:
		val += "collection" + strings.Title(collection.Name) + "Mutex.RLock()\n"
		val += "collection := mongo" + strings.Title(collection.Name) + "Collection\n"
		val += "collection" + strings.Title(collection.Name) + "Mutex.RUnlock()\n"
		// val += "if collection == nil {\n"
		// val += "init" + strings.Title(collection.Name) + "()\n"
		// val += "}\n"
		val += "t := time.Now()\n"
		val += "objectId := self.Id\n"
		val += "if self.Id == \"\"{\n"
		val += "objectId = bson.NewObjectId()\n"
		val += "self.CreateDate = t\n"
		val += "}\n"
		val += "self.UpdateDate = t\n"
		val += "changeInfo, err := collection.UpsertId(objectId, &self)\n"
		val += "if err != nil {\n"
		val += "log.Println(\"Failed to upsertId for " + strings.Title(schema.Name) + ":  \" + err.Error())\n"
		val += "return err\n"
		val += "}\n"
		val += "if changeInfo.UpsertedId != nil {\n"
		val += "self.Id = changeInfo.UpsertedId.(bson.ObjectId)\n"
		val += "}\n"
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "if store.OnChangeRecord != nil && len(store.OnRecordUpdate) > 0 {\n"
		val += "if store.OnRecordUpdate[0] == \"*\" || utils.InArray(\"" + strings.Title(collection.Name) + "\", store.OnRecordUpdate) {\n"
		val += "  value := reflect.ValueOf(&self)\n"
		val += "  store.OnChangeRecord(\"" + strings.Title(collection.Name) + "\", self.Id.Hex(), value.Interface())\n"
		val += "}\n"
		val += "}\n"
		val += "pubsub.Publish(\"" + strings.Title(collection.Name) + ".Save\", self)\n"
		val += "return nil\n"
	}
	val += "}\n\n"
	return val
}

func genSetCollection(collection NOSQLCollection) string {
	return `func (obj model` + strings.Title(collection.Name) + `) SetCollection(mdb *mgo.Database) {
		collection` + strings.Title(collection.Name) + `Mutex.Lock()
		mongo` + strings.Title(collection.Name) + `Collection = mdb.C("` + strings.Title(collection.Name) + `")
		ci := mgo.CollectionInfo{ForceIdIndex: true}
		mongo` + strings.Title(collection.Name) + `Collection.Create(&ci)
		collection` + strings.Title(collection.Name) + `Mutex.Unlock()
	}

`
}

func genById(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	return `func (obj model` + strings.Title(collection.Name) + `) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
		var retObj ` + strings.Title(schema.Name) + `
		q := obj.Query()
		for i := range joins {
			joinValue := joins[i]
			q = q.Join(joinValue)
		}
		err = q.ById(objectID, &retObj)
		value = reflect.ValueOf(&retObj)
		return
	}
`
}

func genDoesIdExist(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	return `func (obj *` + strings.Title(schema.Name) + `) DoesIdExist(objectID interface{}) bool {
		var retObj ` + strings.Title(schema.Name) + `
		row := model` + strings.Title(collection.Name) + `{}
		q := row.Query()
		err := q.ById(objectID, &retObj)
		if err == nil {
			return true
		} else {
			return false
		}
	}
`
}

func genNewByReflection(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (obj model" + strings.Title(collection.Name) + ") NewByReflection() (value reflect.Value) {\n"
	val += "retObj := " + strings.Title(schema.Name) + "{}\n"
	val += "value = reflect.ValueOf(&retObj)\n"
	val += "return\n"
	val += "}\n\n"
	return val
}

func genParseInterface(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (obj *" + strings.Title(schema.Name) + ") ParseInterface(x interface{}) (err error) {\n"
	val += "data, err := json.Marshal(x)\n"
	val += "if err != nil {\n"
	val += "	return\n"
	val += "}\n"
	val += "err = json.Unmarshal(data, obj)\n"
	val += "return\n"
	val += "}\n"
	return val
}

func genNewId(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (obj *" + strings.Title(schema.Name) + ") NewId() {\n"
	val += "obj.Id = bson.NewObjectId()\n"
	val += "}\n\n"
	return val
}

func genCountByFilter(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (obj model" + strings.Title(collection.Name) + ") CountByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (count int, err error) {\n"
	val += "var retObj []" + strings.Title(schema.Name) + "\n"
	val += "q := obj.Query().Filter(filter)\n"
	val += "if len(inFilter) > 0 {\n"
	val += "	q = q.In(inFilter)\n"
	val += "}\n"
	val += "if len(excludeFilter) > 0 {\n"
	val += "	q = q.Exclude(excludeFilter)\n"
	val += "}\n"
	val += "// joins really make no sense here but just copy paste coding here\n"
	val += "for i := range joins {\n"
	val += "joinValue := joins[i]\n"
	val += "q = q.Join(joinValue)\n"
	val += "}\n"
	val += "cnt, errCount := q.Count(&retObj)\n"
	val += "return cnt, errCount\n"
	val += "}\n\n"
	return val
}

func genByFilter(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (obj model" + strings.Title(collection.Name) + ") ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {\n"
	val += "var retObj []" + strings.Title(schema.Name) + "\n"
	val += "q := obj.Query().Filter(filter)\n"
	val += "if len(inFilter) > 0 {\n"
	val += "	q = q.In(inFilter)\n"
	val += "}\n"
	val += "if len(excludeFilter) > 0 {\n"
	val += "	q = q.Exclude(excludeFilter)\n"
	val += "}\n"
	val += "for i := range joins {\n"
	val += "joinValue := joins[i]\n"
	val += "q = q.Join(joinValue)\n"
	val += "}\n"
	val += "err = q.All(&retObj)\n"
	val += "value = reflect.ValueOf(&retObj)\n"
	val += "return\n"
	val += "}\n\n"
	return val
}

func genNoSQLSchemaSaveByTran(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") SaveWithTran(t *Transaction) error {\n\n"
	val += "return self.CreateWithTran(t, false)\n"
	val += "}\n"
	val += "func (self *" + strings.Title(schema.Name) + ") ForceCreateWithTran(t *Transaction) error {\n\n"
	val += "return self.CreateWithTran(t, true)\n"
	val += "}\n"
	val += "func (self *" + strings.Title(schema.Name) + ") CreateWithTran(t *Transaction, forceCreate bool) error {\n\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "return self.Save()\n"
		val += "}\n"
	case DATABASE_DRIVER_MONGODB:
		val += `
		transactionQueue.Lock()
		defer func() {
			transactionQueue.Unlock()
		}()


		// collection` + strings.Title(collection.Name) + `Mutex.RLock()
		// collection := mongo` + strings.Title(collection.Name) + `Collection
		// collection` + strings.Title(collection.Name) + `Mutex.RUnlock()
		// if collection == nil {
		// 	init` + strings.Title(collection.Name) + `()
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

		t.Collections = append(t.Collections, "` + strings.Title(collection.Name) + `History")
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
		dbServices.CollectionCache{}.Remove("` + strings.Title(collection.Name) + `",self.Id.Hex())
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
		var histRecord ` + strings.Title(schema.Name) + `HistoryRecord
		histRecord.TId = t.Id.Hex()
		histRecord.Data = newBson
		histRecord.Type = TRANSACTION_CHANGETYPE_INSERT

		histRecord.ObjId = self.Id.Hex()
		histRecord.CreateDate = time.Now()
		//Get the Original Record if it is a Update
		if isUpdate {

			_, ok := transactionQueue.queue[t.Id.Hex()].newItems["` + strings.Title(schema.Name) + `_" + self.Id.Hex()]
			if ok {
				transactionQueue.queue[t.Id.Hex()].newItems["` + strings.Title(schema.Name) + `_" + self.Id.Hex()] = eTransactionNew
			}
			histRecord.Type = TRANSACTION_CHANGETYPE_UPDATE
			eTransactionNew.changeType = TRANSACTION_CHANGETYPE_UPDATE
			var original ` + strings.Title(schema.Name) + `
			err := ` + strings.Title(collection.Name) + `.Query().ById(self.Id, &original)
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
		transactionQueue.queue[t.Id.Hex()].newItems["` + strings.Title(schema.Name) + `_" + self.Id.Hex()] = eTransactionNew
		transactionQueue.queue[t.Id.Hex()].originalItems["` + strings.Title(schema.Name) + `_" + self.Id.Hex()] = eTransactionOriginal
		return nil
	}
`
	}

	return val
}

func genNoSQLValidate(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") ValidateAndClean() error {\n\n"
	val += "return validateFields(" + strings.Title(schema.Name) + "{}, self, reflect.ValueOf(self).Elem())\n"
	val += "}\n\n"
	return val
}

func genNoSQLReflect(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") Reflect() []Field {\n\n"
	val += "return Reflect(" + strings.Title(schema.Name) + "{})\n"
	val += "}\n\n"
	return val
}

func genNoSQLUnMarshal(collection NOSQLCollection, schema NOSQLSchema, driver string) string {

	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") Unmarshal(data []byte) error {\n\n"

	val += "err := bson.Unmarshal(data, &self)\n"

	val += "if err != nil {\n"
	val += "	return err\n"
	val += "}\n"
	val += "return nil\n"
	val += "}\n\n"
	return val
}

func genNoSQLSchemaArrayCheck(schema NOSQLSchema) string {
	val := "if len(retObj) == 0 {\n"
	val += "retObj = []" + strings.Title(schema.Name) + "{}\n"
	val += "}\n"
	return val
}

func genNoSQLSchemaSingle(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (self model" + strings.Title(collection.Name) + ") Single(field string, value interface{}) (retObj " + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "e = dbServices.BoltDB.One(field, value, &retObj)\n"
		val += "return\n"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaSearch(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") Search(field string, value interface{}) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "e = dbServices.BoltDB.Find(field, value, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	}
	val += "}\n\n"

	val += "func (obj model" + strings.Title(collection.Name) + ") SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "if limit == 0 && skip == 0{\n"
		val += "	e = dbServices.BoltDB.Find(field, value, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "	e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "	e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "	e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "return\n"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaAll(collection NOSQLCollection, schema NOSQLSchema, driver string) string {

	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") All() (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		val += "e = dbServices.BoltDB.All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	}
	val += "}\n\n"

	val += "func (obj model" + strings.Title(collection.Name) + ") AllAdvanced(limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "if limit == 0 && skip == 0{\n"
		val += "	e = dbServices.BoltDB.All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "	e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "	e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "	e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "return\n"
	}
	val += "}\n\n"

	return val
}

func genNoSQLSchemaAllByIndex(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") AllByIndex(index string) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "e = dbServices.BoltDB.AllByIndex(index, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"

	}
	val += "}\n\n"

	val += "func (obj model" + strings.Title(collection.Name) + ") AllByIndexAdvanced(index string, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "if limit == 0 && skip == 0{\n"
		val += "	e = dbServices.BoltDB.AllByIndex(index, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "	e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "	e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "	e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "return\n"
	}
	val += "}\n\n"

	return val
}

func genNoSQLSchemaRange(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") Range(min, max, field string) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "e = dbServices.BoltDB.Range(field, min, max, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	}
	val += "}\n\n"

	val += "func (obj model" + strings.Title(collection.Name) + ") RangeAdvanced(min, max, field string, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "if limit == 0 && skip == 0{\n"
		val += "	e = dbServices.BoltDB.Range(field, min, max, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "	e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "	e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "	e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "return\n"
	case DATABASE_DRIVER_MONGODB:

		val += "var query *mgo.Query\n"
		val += "m := make(bson.M)\n"
		val += "m[field] = bson.M{\"$gte\": min, \"$lte\": max}\n"
		val += "query = mongo" + strings.Title(collection.Name) + "Collection.Find(m)\n"

		val += "if limit == 0 && skip == 0{\n"
		val += "e = query.All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "e = query.Limit(limit).Skip(skip).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "e = query.Limit(limit).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "e = query.Skip(skip).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "return\n"
	}
	val += "}\n\n"

	return val
}

func genNoSQLSchemaIndex(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") Index() error {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "return dbServices.BoltDB.Init(&" + strings.Title(schema.Name) + "{})\n"
	case DATABASE_DRIVER_MONGODB:
		val += "log.Println(\"Building Indexes for MongoDB collection " + collection.Name + ":\")\n"
		val += "for key, value := range dbServices.GetDBIndexes(" + strings.Title(schema.Name) + "{}) {\n"
		val += "index := mgo.Index{\n"
		val += "Key:        []string{key},\n"
		val += "Unique:     false,\n"
		val += "Background: true,\n"
		val += "}\n\n"
		val += "if value == \"unique\" {\n"
		val += "index.Unique = true\n"
		val += "}\n"
		val += "\n"

		val += "collection" + strings.Title(collection.Name) + "Mutex.RLock()\n"
		val += "collection := mongo" + strings.Title(collection.Name) + "Collection\n"
		val += "collection" + strings.Title(collection.Name) + "Mutex.RUnlock()\n"

		val += "err := collection.EnsureIndex(index)\n"
		val += "if err != nil {\n"
		val += "log.Println(\"Failed to create index for " + strings.Title(schema.Name) + ".\" + key + \":  \" + err.Error())\n"
		val += "}else{\n"
		val += "log.Println(\"Successfully created index for " + strings.Title(schema.Name) + ".\" + key)\n"
		val += "}\n"
		val += "}\n"

		val += "return nil"
	}
	val += "}\n\n"
	return val
}

//Generates the Bootstrap Data for the application.
func genNoSQLBootstrap(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") BootStrapComplete() {\n"
	val += "	GoCore" + strings.Title(collection.Name) + "HasBootStrapped.Set(true)\n"
	val += "}\n"

	val += "func (obj model" + strings.Title(collection.Name) + ") Bootstrap() error {\n"

	val += "start := time.Now()\n"
	val += "defer func() {\n"
	val += "log.Println(logger.TimeTrack(start, \"Bootstraping of " + strings.Title(collection.Name) + " Took\"))\n"
	val += "}()\n"
	val += "if serverSettings.WebConfig.Application.BootstrapData == false {\n"
	val += "	obj.BootStrapComplete()\n"
	val += "	return nil\n"
	val += "}\n\n"

	val += "var isError bool\n"
	val += "var query Query\n"

	if driver == DATABASE_DRIVER_MONGODB {
		val += heredoc.Docf(`
			collection%sMutex.RLock()
			query.collection = mongo%sCollection
			collection%sMutex.RUnlock()
			`, strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name))
	}

	val += "var rows []" + strings.Title(schema.Name) + "\n"
	val += "cnt, errCount := query.Count(&rows)\n"
	val += "if errCount != nil{\n"
	val += "cnt = 1\n"
	val += "}\n\n"

	if collection.ClearTable && driver == DATABASE_DRIVER_MONGODB {
		val += "query.collection.DropCollection()\n"
		val += "cnt = 0\n"
	} else if collection.ClearTable && driver == DATABASE_DRIVER_BOLTDB {
		val += "x, errDelete := obj.All()\n"
		val += "if errDelete == nil {\n"
		val += "for _, row := range x {\n"
		val += "row.Delete()\n"
		val += "}\n"
		val += "}\n"
		val += "cnt = 0\n"
	}

	//First check if the path exists to bootstrap data
	path := serverSettings.APP_LOCATION + "/db/bootstrap/" + extensions.MakeFirstLowerCase(collection.Name) + "/" + extensions.MakeFirstLowerCase(collection.Name) + ".json"
	if extensions.DoesFileExist(path) {

		data, err := extensions.ReadFile(path)

		if err != nil {
			color.Red("Reading of " + path + " failed to create Bootstrap Data:  " + err.Error())
			val += "dataString :=\"\"\n\n"
		} else {
			val += "dataString :=\"" + base64.StdEncoding.EncodeToString(data[:]) + "\"\n\n"
		}

	} else {
		val += "dataString := \"\"\n\n"
	}

	val += heredoc.Docf(`
		var files [][]byte
		var err error
		var distDirectoryFound bool
		err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "%s")
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for %s due to caching issue: " + err.Error())
			return err
		}

		files, err, distDirectoryFound = BootstrapDirectory("%s", cnt)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for %s: " + err.Error())
			return err
		}

		if dataString != "" {
			data, err := base64.StdEncoding.DecodeString(dataString)
			if err != nil {
				obj.BootStrapComplete()
				log.Println("Failed to bootstrap data for %s: " + err.Error())
				return err
			}
			files = append(files, data)
		}

		var v []%s
		for _, file := range files{
			var fileBootstrap []%s
			hash := md5.Sum(file)
			hexString := hex.EncodeToString(hash[:])
			err = json.Unmarshal(file, &fileBootstrap)
			if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName + "%s", hexString) || cnt == 0 {
				if err != nil {

					logger.Message("Failed to bootstrap data for %s: " + err.Error(), logger.RED)
					utils.TalkDirtyToMe("Failed to bootstrap data for %s: " + err.Error())
					continue
				}

				fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName + "%s", hexString)

				for i, _ := range fileBootstrap {
					fb := fileBootstrap[i]
					v = append(v, fb)
				}
			}
		}
		fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "%s")

`, strings.Title(collection.Name), strings.Title(collection.Name), extensions.MakeFirstLowerCase(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(schema.Name), strings.Title(schema.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name))
	val += heredoc.Docf(`
		var actualCount int
		originalCount := len(v)
		log.Println("Total count of records attempting %s", len(v))

		for _, doc := range v {
			var original %s
			if doc.Id.Hex() == "" {
				doc.Id = bson.NewObjectId()
			}
			err = query.ById(doc.Id, &original)
			if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "%s" {
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
					err = doc.Delete()
					if err != nil {
						log.Println("Failed to delete data for %s:  " + doc.Id.Hex() + "  " + err.Error())
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
					if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.ProductNames) > 0 &&  !utils.InArray(serverSettings.WebConfig.Application.ProductName, doc.BootstrapMeta.ProductNames) {
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
							log.Println("Failed to bootstrap data for %s:  " + doc.Id.Hex() + "  " + err.Error())
							isError = true
						}
					} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
						log.Println("%s skipped a row for some reason on " + doc.Id.Hex() + " because of " +  core.Debug.GetDump(reason))
					}
				}
			} else {
				actualCount += 1
			}
		}
		if isError {
			log.Println("FAILED to bootstrap %s")
		} else {

			if distDirectoryFound == false {
				err = BootstrapMongoDump("%s", "%s")
			}
			if err == nil {
				log.Println("Successfully bootstrapped %s")
				if actualCount != originalCount {
					logger.Message("%s counts are different than original bootstrap and actual inserts, please inpect data." + core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
				}
			}
		}
`, strings.Title(collection.Name), strings.Title(schema.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), extensions.MakeFirstLowerCase(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name), strings.Title(collection.Name))
	val += "obj.BootStrapComplete()\n"
	val += "return nil\n"
	val += "}\n\n"
	return val
}

func genNoSQLSchemaRunTransaction(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") RunTransaction(objects []" + strings.Title(schema.Name) + ") error {\n\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "tx, err := dbServices.BoltDB.Begin(true)\n\n"

		val += "for _, object := range objects {\n"
		val += "	err = tx.Save(&object)\n"
		val += "	if err != nil {\n"
		val += "		tx.Rollback()\n"
		val += "		return err\n"
		val += "	}\n"
		val += "}\n\n"

		val += "tx.Commit()\n\n"

		val += "return nil\n"
	case DATABASE_DRIVER_MONGODB:
		val += "return nil\n"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaNew(collection NOSQLCollection, schema NOSQLSchema) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") New() *" + strings.Title(schema.Name) + " {\n"
	val += "return &" + strings.Title(schema.Name) + "{}\n"
	val += "}\n\n"
	return val
}

func genNoSQLSchemaSetKeyValue(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj model" + strings.Title(collection.Name) + ") SetKeyValue() error {\n"
	switch driver {
	case "boltDB":
		{
			val += "return dbServices.BoltDB.Init(&" + strings.Title(schema.Name) + "{})\n"
		}
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaDelete(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (self *" + strings.Title(schema.Name) + ") Delete() error {\n"
	switch driver {
	case "boltDB":
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "err := dbServices.BoltDB.Delete(\"" + strings.Title(schema.Name) + "\", self.Id.Hex())\n"
		val += "if err == nil{\n"
		val += "pubsub.Publish(\"" + strings.Title(collection.Name) + ".Delete\", self)\n"
		val += "}\n"
		val += "return err\n"
	case "mongoDB":
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "err := mongo" + strings.Title(collection.Name) + "Collection.RemoveId(self.Id)\n"
		val += "if err == nil{\n"
		val += "pubsub.Publish(\"" + strings.Title(collection.Name) + ".Delete\", self)\n"
		val += "}\n"
		val += "return err\n"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaDeleteWithTran(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (self *" + strings.Title(schema.Name) + ") DeleteWithTran(t *Transaction) error {\n"
	switch driver {
	case "boltDB":
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "return dbServices.BoltDB.Delete(\"" + strings.Title(collection.Name) + "\", self.Id.Hex())\n"
	case "mongoDB":
		val += "transactionQueue.Lock()\n"
		val += "defer func() {\n"
		val += "  transactionQueue.Unlock()\n"
		val += "}()\n"
		val += "if self.Id.Hex() == \"\" {\n"
		val += "return errors.New(dbServices.ERROR_CODE_TRANSACTION_RECORD_NOT_EXISTS)\n"
		val += "}\n\n"
		val += "dbServices.CollectionCache{}.Remove(\"" + strings.Title(collection.Name) + "\",self.Id.Hex())\n"
		val += "_, ok := transactionQueue.queue[t.Id.Hex()]\n"
		val += "if ok == false {\n"
		val += "	return errors.New(dbServices.ERROR_CODE_TRANSACTION_NOT_PRESENT)\n"
		val += "}\n\n"

		val += "var histRecord " + strings.Title(schema.Name) + "HistoryRecord\n"
		val += "histRecord.TId = t.Id.Hex()\n\n"
		val += "histRecord.Type = TRANSACTION_CHANGETYPE_DELETE\n"
		val += "histRecord.ObjId = self.Id.Hex()\n"
		val += "histRecord.CreateDate = time.Now()\n"

		val += "var eTransactionNew entityTransaction\n"
		val += "eTransactionNew.changeType = TRANSACTION_CHANGETYPE_DELETE\n"
		val += "eTransactionNew.entity = self\n\n"

		val += "var eTransactionOriginal entityTransaction\n"
		val += "eTransactionOriginal.changeType = TRANSACTION_CHANGETYPE_DELETE\n"
		val += "eTransactionOriginal.entity = &histRecord\n\n"

		val += "var original " + strings.Title(schema.Name) + "\n"
		val += "err := " + strings.Title(collection.Name) + ".Query().ById(self.Id, &original)\n\n"

		val += "if err != nil {\n"
		val += "	return err\n"
		val += "}\n\n"

		val += "originalJson, err := original.JSONString()\n\n"

		val += "if err != nil {\n"
		val += "	return err\n"
		val += "}\n\n"

		val += "originalBase64 := getBase64(originalJson)\n"
		val += "histRecord.Data = originalBase64\n\n"

		val += "t.Collections = append(t.Collections, \"" + strings.Title(collection.Name) + "History\")\n"

		val += "if len(transactionQueue.queue[t.Id.Hex()].originalItems) == 0 {\n"
		val += "	transactionQueue.queue[t.Id.Hex()].originalItems = make(map[string]entityTransaction, 0)\n"
		val += "}\n"
		val += "if len(transactionQueue.queue[t.Id.Hex()].newItems) == 0 {\n"
		val += "	transactionQueue.queue[t.Id.Hex()].newItems = make(map[string]entityTransaction, 0)\n"
		val += "}\n"
		val += "transactionQueue.queue[t.Id.Hex()].newItems[\"" + strings.Title(collection.Name) + "_\" + self.Id.Hex()] = eTransactionNew\n\n"
		val += "transactionQueue.queue[t.Id.Hex()].originalItems[\"" + strings.Title(collection.Name) + "_\" + self.Id.Hex()] = eTransactionOriginal\n\n"
		val += "transactionQueue.ids[t.Id.Hex()] = append(transactionQueue.ids[t.Id.Hex()], eTransactionNew.entity.GetId())\n"

		val += "return nil"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaJoinFields(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {\n\n"

	val += "source := reflect.ValueOf(self).Elem()\n\n"

	val += "var joins []join\n"
	val += "joins, err = getJoins(source, remainingRecursions)\n\n"

	val += "if len(joins) == 0 {\n"
	val += "	return\n"
	val += "}\n\n"

	val += "s := source\n"
	val += "for _, j := range joins {\n"
	val += "	id := reflect.ValueOf(q.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()\n"
	val += "	joinsField := s.FieldByName(\"Joins\")\n"
	val += "	setField := joinsField.FieldByName(j.joinFieldName)\n\n"
	val += " endRecursion := false\n"
	val += " if serverSettings.WebConfig.Application.LogJoinQueries {\n"
	val += " fmt.Print(\"Remaining Recursions\")\n"
	val += " fmt.Println(fmt.Sprintf(\"%+v\", remainingRecursions))\n"
	val += " fmt.Println(fmt.Sprintf(\"%+v\", j.collectionName))\n"
	val += "}\n"
	val += " if remainingRecursions == j.joinSpecified {\n"
	val += " endRecursion = true\n"
	val += "}\n"
	val += "	err = joinField(j, id, setField, j.joinSpecified, q, endRecursion, recursionCount)\n"
	val += "	if err != nil {\n"
	val += "		return\n"
	val += "	}\n"
	val += "}\n"
	val += "return\n"
	val += "}\n\n"

	return val
}

func genNoSQLBucketCore(driver string) string {
	val := ""

	val += "type Bucket struct{\n"
	val += "	Name string\n"
	val += "}\n\n"

	val += "func (obj *Bucket) SetKeyValue(key interface{}, value interface{}) error {\n"
	switch driver {
	case "boltDB":
		val += "return dbServices.BoltDB.Set(obj.Name, key, value)\n"
	case "mongoDB":
		val += "return nil\n"
	}
	val += "}\n\n"

	val += "func (obj *Bucket) GetKeyValue(key interface{}, value interface{}) error {\n"
	switch driver {
	case "boltDB":
		val += "return dbServices.BoltDB.Get(obj.Name, key, value)\n"
	case "mongoDB":
		val += "return nil\n"
	}
	val += "}\n\n"

	val += "func (obj *Bucket) DeleteKey(key interface{}) error {\n"
	switch driver {
	case "boltDB":
		val += "return dbServices.BoltDB.Delete(obj.Name, key)\n"
	case "mongoDB":
		val += "return nil\n"
	}
	val += "}\n\n"
	return val
}

func getNoSQLSchemaPrimaryKey(schema NOSQLSchema) string {

	for _, field := range schema.Fields {
		if field.Index == "primary" {
			return strings.Title(field.Name)
		}
	}

	return ""
}
