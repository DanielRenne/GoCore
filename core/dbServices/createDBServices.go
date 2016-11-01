package dbServices

import (
	"encoding/json"
	"log"

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

	"github.com/fatih/color"
)

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
}

type NOSQLSchema struct {
	Name   string             `json:"name"`
	Fields []NOSQLSchemaField `json"fields"`
}

type NOSQLCollection struct {
	Name   string      `json:"name"`
	Schema NOSQLSchema `json:"schema"`
}

type NOSQLSchemaDB struct {
	Collections []NOSQLCollection `json:"collections"`
}

type entityList struct {
	Constants []string
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

	walkNoSQLSchema()

}

func walkNoSQLSchema() {

	allCollections.Entities = make(map[string]entityList)
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

			//Create Swagger Definition With the latest Version being equal to swagger.json, all others swagger_1.0.0.json etc...
			writeSwaggerConfiguration("/api/"+versionDir, version.Value)
		}
	}

	//Make a copy of the latest Swagger Version Definition for Swagger UI to default to swagger.json.  We will keep the latest version as a 2nd copy.
	extensions.CopyFile(serverSettings.SWAGGER_UI_PATH+"/swagger."+versionNumber+".json", serverSettings.SWAGGER_UI_PATH+"/swagger.json")

}

func walkNoSQLVersion(path string, versionDir string) {

	initializeModelFile()

	var wg sync.WaitGroup

	var scs schemasCreatedSync

	scs.schemasCreated = make(map[string]NOSQLSchema, 0)

	err := filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) error {

		if errWalk != nil {
			return errWalk
		}

		var e error

		if filepath.Ext(f.Name()) == ".json" {
			wg.Add(1)

			go func() {
				defer wg.Done()
				jsonData, err := ioutil.ReadFile(path)
				if err != nil {
					color.Red("Reading of " + path + " failed:  " + err.Error())
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
			}()
		}

		return e
	})
	if err != nil {
		color.Red("Walk of path failed:  " + err.Error())
	}

	wg.Wait()

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
		copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/transaction.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/transaction.go")
		copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/query.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/query.go")
		copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/timeZone.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/timeZone.go")
		copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/timeZoneLocations.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/timeZoneLocations.go")
		copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/locales.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/locales.go")
		////Support for Long Running Transactions Later Maybe
		//copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/transactionObjects.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/transactionObjects.go")
	}

	histTemplate, err := extensions.ReadFile(serverSettings.GOCORE_PATH + "/core/dbServices/mongo/stubs/histTemplate.go")

	if err != nil {
		color.Red("Error reading histTemplate.go:  " + err.Error())
		return
	}

	//Create the Collection Models
	for _, collection := range collections {
		val := generateNoSQLModel(collection.Schema, collection, driver, scs)
		os.Mkdir(serverSettings.APP_LOCATION+"/models/", 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir, 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/", 0777)
		writeNoSQLModelCollection(val, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+".go", collection)

		//Create the Transaction History Table for the Collection
		histModified := strings.Replace(string(histTemplate[:]), "HistCollection", strings.Title(collection.Name)+"History", -1)
		histModified = strings.Replace(histModified, "HistEntity", strings.Title(collection.Schema.Name)+"HistoryRecord", -1)
		histModified = strings.Replace(histModified, "OriginalEntity", strings.Title(collection.Schema.Name), -1)

		writeNoSQLStub(histModified, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+"_Hist.go")

		os.Mkdir(serverSettings.APP_LOCATION+"/webAPIs/", 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/webAPIs/"+versionDir, 0777)
		os.Mkdir(serverSettings.APP_LOCATION+"/webAPIs/"+versionDir+"/webAPI/", 0777)

		cWebAPI := genSchemaWebAPI(collection, collection.Schema, strings.Replace(serverSettings.APP_LOCATION, "src/", "", -1)+"/models/"+versionDir+"/model", driver, versionDir)
		writeNoSQLWebAPI(cWebAPI, serverSettings.APP_LOCATION+"/webAPIs/"+versionDir+"/webAPI/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+".go", collection)
	}

}

func initializeModelFile() {

	modelData, err := extensions.ReadFile(serverSettings.GOCORE_PATH + "/core/dbServices/mongo/stubs/model.go")

	if err != nil {
		color.Red("Failed to read and append model.go:  " + err.Error())
		return
	}

	modelToWrite = string(modelData[:])
	modelToWrite += "\n"

}

func finalizeModelFile(versionDir string) {

	sort.Sort(SchemaNameSorter(allCollections.Collections))

	modelToWrite += "func ResolveEntity(key string) modelEntity{\n\n"
	modelToWrite += "switch key{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "\":\n"
		modelToWrite += " return &" + strings.Title(collection.Schema.Name) + "{}\n"

		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "HistoryRecord\":\n"
		modelToWrite += " return &" + strings.Title(collection.Schema.Name) + "HistoryRecord{}\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"

	modelToWrite += "\n"
	modelToWrite += "func ResolveCollection(key string) collection{\n\n"
	modelToWrite += "switch key{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Name) + "\":\n"
		modelToWrite += " return &model" + strings.Title(collection.Name) + "{}\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"

	modelToWrite += "\n"
	modelToWrite += "func ResolveHistoryCollection(key string) modelCollection{\n\n"
	modelToWrite += "switch key{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Name) + "History\":\n"
		modelToWrite += " return &model" + strings.Title(collection.Name) + "History{}\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"

	modelToWrite += "func joinField(j join, id string, fieldToSet reflect.Value, remainingRecursions string, q *Query, endRecursion bool, recursionCount int) (err error) {\n\n"
	modelToWrite += "c := ResolveCollection(j.collectionName)\n"
	modelToWrite += "if c == nil {\n"
	modelToWrite += "	err = errors.New(\"Failed to resolve collection:  \" + j.collectionName)\n"
	modelToWrite += "	return\n"
	modelToWrite += "}\n"

	modelToWrite += "switch j.joinSchemaName{\n"

	for _, collection := range allCollections.Collections {
		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "\":\n"
		modelToWrite += "var y " + strings.Title(collection.Schema.Name) + "\n"
		modelToWrite += "if j.isMany {\n"
		modelToWrite += "var z []" + strings.Title(collection.Schema.Name) + "\n"
		modelToWrite += "var ji " + strings.Title(collection.Schema.Name) + "JoinItems\n"
		modelToWrite += "fieldToSet.Set(reflect.ValueOf(&ji))\n"
		modelToWrite += "JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)\n"
		modelToWrite += "}else{\n"
		modelToWrite += "JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)\n"
		modelToWrite += "}\n"
	}

	modelToWrite += "}\n"
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

	timeImport := ""
	if checkSchemaForDateTime(schema) {
		timeImport = "time"
	}

	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "github.com/asdine/storm", timeImport})
	case DATABASE_DRIVER_MONGODB:
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "github.com/DanielRenne/GoCore/core/serverSettings", "encoding/json", "gopkg.in/mgo.v2", "gopkg.in/mgo.v2/bson", "log", "time", "errors", "encoding/base64", "reflect"})
		// val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "gopkg.in/mgo.v2/bson", "log", "time"})
	}

	val += genNoSQLCollection(collection, schema, driver)
	val += genNoSQLSchema(collection.Name, schema, driver, scs, 0)
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
	val += "type " + strings.Title(schema.Name) + "JoinItems struct{\n"
	val += " Count int `json:\"Count\"`\n"
	val += " Items *[]" + strings.Title(schema.Name) + " `json:\"Items\"`\n"
	val += "}\n\n"

	if driver == DATABASE_DRIVER_MONGODB {

		val += "var mongo" + strings.Title(collection.Name) + "Collection *mgo.Collection\n\n"

		val += "func init(){\n"
		val += "go func() {\n\n"
		val += "for{\n"
		val += "if dbServices.MongoDB != nil {\n"
		val += "init" + strings.Title(collection.Name) + "()\n"
		val += "return\n"
		val += "}\n"
		val += "<- dbServices.WaitForDatabase()\n"
		val += "}\n"
		val += "}()\n"
		val += "}\n\n"

		val += "func init" + strings.Title(collection.Name) + "(){\n"
		val += "log.Println(\"Building Indexes for MongoDB collection " + collection.Name + ":\")\n"
		val += "mongo" + strings.Title(collection.Name) + "Collection = dbServices.MongoDB.C(\"" + collection.Name + "\")\n"
		val += "ci := mgo.CollectionInfo{ForceIdIndex: true}\n"
		val += "mongo" + strings.Title(collection.Name) + "Collection.Create(&ci)\n"
		val += strings.Title(collection.Name) + ".Index()\n"
		val += strings.Title(collection.Name) + ".Bootstrap()\n"
		val += "}\n\n"

	}

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
func genNoSQLSchema(collectionName string, schema NOSQLSchema, driver string, scs *schemasCreatedSync, seed int) string {

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
		prefix = strings.Title(collectionName)
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
			fieldPrefix = strings.Title(collectionName)
		}

		additionalTags := genNoSQLAdditionalTags(field, driver)
		omitEmpty := ""
		if field.OmitEmpty {
			omitEmpty = ",omitempty"
		}

		allCollections.Lock()
		entityConsts := allCollections.Entities[strings.Title(schema.Name)]
		entityConsts.Constants = append(entityConsts.Constants, strings.Title(field.Name))
		allCollections.Entities[strings.Title(schema.Name)] = entityConsts
		allCollections.Unlock()

		val += "\n\t" + strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t" + genNoSQLFieldType(fieldPrefix, schema, field, driver) + "\t\t`json:\"" + strings.Title(field.Name) + omitEmpty + "\"" + additionalTags + "`"
	}

	if seed == 0 {
		val += "\n\t CreateDate time.Time `json:\"CreateDate\" bson:\"CreateDate\"`"
		val += "\n\t UpdateDate time.Time `json:\"UpdateDate\" bson:\"UpdateDate\"`"
		val += "\n\t LastUpdateId string `json:\"LastUpdateId\" bson:\"LastUpdateId\"`"
	}

	//Add Validation
	if seed == 0 {
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
			}
		}

		val += "} `json:\"Joins\" bson:\"-\"`\n\n"
	}

	val += "\n}\n\n"

	for _, schemaToCreate := range schemasToCreate {
		val += genNoSQLSchema(collectionName, schemaToCreate, driver, scs, 1)
	}

	return val
}

//Recursive Function to Generate Validation Schema
func genNoSQLValidationRecusion(schema NOSQLSchema) string {
	val := ""
	for _, field := range schema.Fields {
		if field.Type == "object" {
			val += strings.Title(field.Name) + "struct{\n"
			val += genNoSQLValidationRecusion(field.Schema)
			val += "} `json:\"" + strings.Title(field.Name) + "\""
			continue
		}
		if field.Type == "objectArray" {
			val += strings.Title(field.Name) + "[]struct{\n"
			val += genNoSQLValidationRecusion(field.Schema)
			val += "} `json:\"" + strings.Title(field.Name) + "\""
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

		switch field.Index {
		case "":
			return ""
		case "primary":
			return " storm:\"id\""
		case "index":
			return " storm:\"index\""
		case "unique":
			return " storm:\"unique\""
		}
	case DATABASE_DRIVER_MONGODB:

		tags := " bson:\"" + strings.Title(field.Name) + "\"" + validationTagsGap + validationTags
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

	if driver == "mongoDB" && field.Index == "primary" {
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
	case "byteArray":
		return "[]byte"
	case "object":
		return prefix + strings.Title(field.Schema.Name)
	case "intArray":
		return "[]int"
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

	val += genNOSQLQuery(collection, schema, driver)
	val += genNoSQLSchemaIndex(collection, schema, driver)
	val += genNoSQLBootstrap(collection, schema, driver)
	val += genNoSQLSchemaRunTransaction(collection, schema, driver)
	val += genNoSQLSchemaNew(collection, schema)
	val += genNoSQLSchemaSave(collection, schema, driver)
	val += genNoSQLSchemaSaveByTran(collection, schema, driver)
	val += genNoSQLValidate(collection, schema, driver)
	val += genNoSQLReflect(collection, schema, driver)
	val += genNoSQLSchemaDelete(collection, schema, driver)
	val += genNoSQLSchemaDeleteWithTran(collection, schema, driver)
	val += genNoSQLSchemaJoinFields(collection, schema, driver)
	val += genNoSQLUnMarshal(collection, schema, driver)
	val += genNoSQLSchemaJSONRuntime(schema)

	return val
}

func genNOSQLQuery(collection NOSQLCollection, schema NOSQLSchema, driver string) string {

	val := ""
	val += "func (obj model" + strings.Title(collection.Name) + ") Query() *Query {\n"
	val += "	var query Query\n"
	val += "	query.collection = mongo" + strings.Title(collection.Name) + "Collection\n"
	val += "	return &query\n"
	val += "}\n"
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
	return val
}

func genNoSQLSchemaSave(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") Save() error {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "return dbServices.BoltDB.Save(self)\n"
	case DATABASE_DRIVER_MONGODB:
		val += "if mongo" + strings.Title(collection.Name) + "Collection == nil{\n"
		val += "init" + strings.Title(collection.Name) + "()\n"
		val += "}\n"
		val += "objectId := bson.NewObjectId()\n"
		val += "if self.Id != \"\"{\n"
		val += "objectId = self.Id\n"
		val += "}\n"
		val += "changeInfo, err := mongo" + strings.Title(collection.Name) + "Collection.UpsertId(objectId, &self)\n"
		val += "if err != nil {\n"
		val += "log.Println(\"Failed to upsertId for " + strings.Title(schema.Name) + ":  \" + err.Error())\n"
		val += "return err\n"
		val += "}\n"
		val += "if changeInfo.UpsertedId != nil {\n"
		val += "self.Id = changeInfo.UpsertedId.(bson.ObjectId)\n"
		val += "}\n"
		val += "return nil\n"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaSaveByTran(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (self *" + strings.Title(schema.Name) + ") SaveWithTran(t *Transaction) error {\n\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "return dbServices.BoltDB.Save(self)\n"
	case DATABASE_DRIVER_MONGODB:

		val += "//Validate the Model first.  If it fails then clean up the transaction in memory\n"
		val += "err := self.ValidateAndClean()\n"
		val += "if err != nil {\n"
		val += "transactionQueue.Lock()\n"
		val += "delete(transactionQueue.queue, t.Id.Hex())\n"
		val += "transactionQueue.Unlock()\n"
		val += "return err\n"
		val += "}\n"

		val += "transactionQueue.RLock()\n"
		val += "tPersist, ok := transactionQueue.queue[t.Id.Hex()]\n"
		val += "transactionQueue.RUnlock()\n\n"

		val += "if ok == false {\n"
		val += "	return errors.New(dbServices.ERROR_CODE_TRANSACTION_NOT_PRESENT)\n"
		val += "}\n\n"

		val += "t.Collections = append(t.Collections, \"" + strings.Title(collection.Name) + "History\")\n"

		val += "isUpdate := true\n"
		val += "if self.Id.Hex() == \"\"{\n"
		val += "	isUpdate = false\n"
		val += "	self.Id = bson.NewObjectId()\n"
		val += "  self.CreateDate = time.Now()\n"
		val += "}\n"

		val += "	self.UpdateDate = time.Now()\n"
		val += "	self.LastUpdateId = t.UserId\n"

		val += "newJson, err := self.JSONString()\n\n"

		val += "if err != nil {\n"
		val += "	return err\n"
		val += "}\n\n"

		val += "newBase64 := getBase64(newJson)\n"

		val += "var eTransactionNew entityTransaction\n"
		val += "eTransactionNew.changeType = TRANSACTION_CHANGETYPE_INSERT\n"
		val += "eTransactionNew.entity = self\n\n"

		val += "var histRecord " + strings.Title(schema.Name) + "HistoryRecord\n"
		val += "histRecord.TId = t.Id.Hex()\n"
		val += "histRecord.Data = newBase64\n"
		val += "histRecord.Type = TRANSACTION_CHANGETYPE_INSERT\n\n"
		val += "histRecord.ObjId = self.Id.Hex()\n"
		val += "histRecord.CreateDate = time.Now()\n"

		val += "//Get the Original Record if it is a Update\n"
		val += "if isUpdate {\n\n"

		val += "	histRecord.Type = TRANSACTION_CHANGETYPE_UPDATE\n"

		val += "	eTransactionNew.changeType = TRANSACTION_CHANGETYPE_UPDATE\n\n"

		val += "var original " + strings.Title(schema.Name) + "\n"
		val += "err := " + strings.Title(collection.Name) + ".Query().ById(self.Id, &original)\n\n"

		val += "if err != nil {\n"
		val += "	return err\n"
		val += "}\n\n"

		val += "	originalJson, err := original.JSONString()\n\n"

		val += "	if err != nil {\n"
		val += "		return err\n"
		val += "	}\n\n"

		val += "	originalBase64 := getBase64(originalJson)\n"
		val += "	histRecord.Data = originalBase64\n\n"

		val += "}\n\n"

		val += "var eTransactionOriginal entityTransaction\n"
		val += "eTransactionOriginal.entity = &histRecord\n\n"

		val += "tPersist.originalItems = append(tPersist.originalItems, eTransactionOriginal)\n"
		val += "tPersist.newItems = append(tPersist.newItems, eTransactionNew)\n\n"

		val += "return nil\n"

		val += "}\n\n"
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
	val += "func (self *" + strings.Title(schema.Name) + ") Unmarshal(data string) error {\n\n"

	val += "err := json.Unmarshal([]byte(data), &self)\n"

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

	val += "func (self *" + strings.Title(collection.Name) + ") Single(field string, value interface{}) (retObj " + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "e = dbServices.BoltDB.One(field, value, &retObj)\n"
		val += "return\n"
	case DATABASE_DRIVER_MONGODB:

		val += "if field == \"Id\"{\n"
		val += "query := mongo" + strings.Title(collection.Name) + "Collection.FindId(bson.ObjectIdHex(value.(string)))\n"
		val += "e = query.One(&retObj)\n"
		val += "return\n"
		val += "}\n"
		val += "m := make(bson.M)\n"
		val += "m[field] = value\n"
		val += "query := mongo" + strings.Title(collection.Name) + "Collection.Find(m)\n"
		val += "e = query.One(&retObj)\n"
		val += "return\n"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaSearch(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") Search(field string, value interface{}) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "e = dbServices.BoltDB.Find(field, value, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	case DATABASE_DRIVER_MONGODB:

		val += "var query *mgo.Query\n"
		val += "if field == \"Id\"{\n"
		val += "query = mongo" + strings.Title(collection.Name) + "Collection.FindId(bson.ObjectIdHex(value.(string)))\n"
		val += "}else{\n"
		val += "m := make(bson.M)\n"
		val += "m[field] = value\n"
		val += "query = mongo" + strings.Title(collection.Name) + "Collection.Find(m)\n"
		val += "}\n\n"
		val += "e = query.All(&retObj)\n"
		val += "return\n"
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
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
	case DATABASE_DRIVER_MONGODB:

		val += "var query *mgo.Query\n"
		val += "if field == \"Id\"{\n"
		val += "query = mongo" + strings.Title(collection.Name) + "Collection.FindId(bson.ObjectIdHex(value.(string)))\n"
		val += "}else{\n"
		val += "m := make(bson.M)\n"
		val += "m[field] = value\n"
		val += "query = mongo" + strings.Title(collection.Name) + "Collection.Find(m)\n"
		val += "}\n\n"

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

func genNoSQLSchemaAll(collection NOSQLCollection, schema NOSQLSchema, driver string) string {

	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") All() (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		val += "e = dbServices.BoltDB.All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	case "mongoDB":
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") AllAdvanced(limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
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
	case DATABASE_DRIVER_MONGODB:
		val += "if limit == 0 && skip == 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Limit(limit).Skip(skip).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Limit(limit).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Skip(skip).All(&retObj)\n"
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

	val += "func (obj *" + strings.Title(collection.Name) + ") AllByIndex(index string) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += "e = dbServices.BoltDB.AllByIndex(index, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"
	case DATABASE_DRIVER_MONGODB:
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Sort(index).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"

	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") AllByIndexAdvanced(index string, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
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

	case DATABASE_DRIVER_MONGODB:
		val += "if limit == 0 && skip == 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Sort(index).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0 && skip > 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Sort(index).Limit(limit).Skip(skip).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if limit > 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Sort(index).Limit(limit).All(&retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "	return\n"
		val += "}\n"
		val += "if skip > 0{\n"
		val += "e = mongo" + strings.Title(collection.Name) + "Collection.Find(bson.M{}).Sort(index).Skip(skip).All(&retObj)\n"
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

	val += "func (obj *" + strings.Title(collection.Name) + ") Range(min, max, field string) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case DATABASE_DRIVER_BOLTDB:

		val += "e = dbServices.BoltDB.Range(field, min, max, &retObj)\n"
		val += genNoSQLSchemaArrayCheck(schema)
		val += "return\n"

	case DATABASE_DRIVER_MONGODB:
		val += "var query *mgo.Query\n"

		val += "m := make(bson.M)\n"
		val += "m[field] = bson.M{\"$gte\": min, \"$lte\": max}\n"
		val += "query = mongo" + strings.Title(collection.Name) + "Collection.Find(m)\n"

		val += "e = query.All(&retObj)\n"
		val += "return\n"
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") RangeAdvanced(min, max, field string, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
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
		val += "err := mongo" + strings.Title(collection.Name) + "Collection.EnsureIndex(index)\n"
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

	val += "func (obj model" + strings.Title(collection.Name) + ") Bootstrap() error {\n"

	val += "if serverSettings.WebConfig.Application.BootstrapData == false {\n"
	val += "	return nil\n"
	val += "}\n"

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

	//Now parse the data into an array of the collection

	val += "if dataString == \"\"{\n"
	val += "	return nil\n"
	val += "}\n\n"

	val += "data, err := base64.StdEncoding.DecodeString(dataString)\n"
	val += "if err != nil{\n"
	val += "	log.Println(\"Failed to bootstrap data for " + collection.Name + ":  \" + err.Error())\n"
	val += "	return err\n"
	val += "}\n\n"

	val += "var v []" + strings.Title(schema.Name) + "\n"
	val += "err = json.Unmarshal(data, &v)\n"

	val += "if err != nil {\n"
	val += "	log.Println(\"Failed to bootstrap data for " + strings.Title(collection.Name) + ":  \" + err.Error())\n"
	val += "	return err\n"
	val += "}\n\n"

	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += ""
	case DATABASE_DRIVER_MONGODB:
		val += "var isError bool\n"
		val += "for _, doc := range v{\n\n"

		val += "err = doc.Save()\n"
		val += "	if err != nil {\n"
		val += "		log.Println(\"Failed to bootstrap data for " + strings.Title(schema.Name) + ":  \" + doc.Id.Hex() + \"  \" + err.Error())\n"
		val += "isError = true\n"
		val += "	}\n"

		// val += "	log.Printf(\"%+v\\n\", doc)\n\n"
		val += "}\n"
		val += "if isError{\n"
		val += "	log.Println(\"FAILED to bootstrap " + strings.Title(collection.Name) + "\")\n"
		val += "}else{\n"
		val += "	log.Println(\"Successfully bootstraped " + strings.Title(collection.Name) + "\")\n"
		val += "}\n\n"

	}
	val += "return nil"
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

	val += "func (obj *" + strings.Title(collection.Name) + ") SetKeyValue() error {\n"
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
		val += "return dbServices.BoltDB.Remove(self)\n"
	case "mongoDB":
		val += "return mongo" + strings.Title(collection.Name) + "Collection.RemoveId(self.Id)"
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaDeleteWithTran(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (self *" + strings.Title(schema.Name) + ") DeleteWithTran(t *Transaction) error {\n"
	switch driver {
	case "boltDB":
		val += "return dbServices.BoltDB.Remove(self)\n"
	case "mongoDB":
		val += "if self.Id.Hex() == \"\" {\n"
		val += "return errors.New(dbServices.ERROR_CODE_TRANSACTION_RECORD_NOT_EXISTS)\n"
		val += "}\n\n"

		val += "transactionQueue.RLock()\n"
		val += "tPersist, ok := transactionQueue.queue[t.Id.Hex()]\n"
		val += "transactionQueue.RUnlock()\n\n"

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

		val += "tPersist.originalItems = append(tPersist.originalItems, eTransactionOriginal)\n"
		val += "tPersist.newItems = append(tPersist.newItems, eTransactionNew)\n"

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
	val += "joins = getJoins(source, remainingRecursions)\n\n"

	val += "if len(joins) == 0 {\n"
	val += "	return\n"
	val += "}\n\n"

	val += "s := source\n"
	val += "for _, j := range joins {\n"
	val += "	id := reflect.ValueOf(q.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()\n"
	val += "	joinsField := s.FieldByName(\"Joins\")\n"
	val += "	setField := joinsField.FieldByName(j.joinFieldName)\n\n"
	val += " endRecursion := false\n"
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
