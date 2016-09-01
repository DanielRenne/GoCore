package dbServices

import (
	"encoding/json"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	// "fmt"
	"encoding/base64"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type fieldValidation struct {
	Required  bool   `json:"required"`
	Type      string `json:"type"`
	Min       string `json:"min"`
	Max       string `json:"max"`
	Length    string `json:"length"`
	LengthMax string `json:"lengthMax"`
	LengthMin string `json:"lengthMin"`
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
	OmitEmpty    bool             `json:"omitEmpty"`
	DefaultValue string           `json:"defaultValue"`
	Required     bool             `json:"required"`
	Schema       NOSQLSchema      `json:"schema"`
	Validation   *fieldValidation `json:"validate, omitempty"`
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

	var collections []NOSQLCollection

	var wg sync.WaitGroup

	schemasCreated := make([]NOSQLSchema, 0)

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
					collections = append(collections, col)
				}

				createNoSQLModel(schemaDB.Collections, serverSettings.WebConfig.DbConnection.Driver, versionDir, &schemasCreated)
			}()
		}

		return e
	})
	if err != nil {
		color.Red("Walk of path failed:  " + err.Error())
	}

	wg.Wait()

	finalizeModelFile(versionDir, collections)
}

func createNoSQLModel(collections []NOSQLCollection, driver string, versionDir string, schemasCreated *[]NOSQLSchema) {

	//Clean the Model and API Directory
	extensions.RemoveDirectory(serverSettings.APP_LOCATION + "/models/" + versionDir + "/model")
	extensions.RemoveDirectory(serverSettings.APP_LOCATION + "/webAPIs/" + versionDir + "/webAPI")

	//Create a NOSQLBucket Model
	bucket := generateNoSQLModelBucket(driver)
	os.Mkdir(serverSettings.APP_LOCATION+"/models/", 0777)
	os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir, 0777)
	os.Mkdir(serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/", 0777)

	writeNOSQLModelBucket(bucket, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/bucket.go")

	//Copy Stub Files
	if driver == DATABASE_DRIVER_MONGODB {
		copyNoSQLStub(serverSettings.GOCORE_PATH+"/core/dbServices/mongo/stubs/transaction.go", serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/transaction.go")

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
		val := generateNoSQLModel(collection.Schema, collection, driver, schemasCreated)
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

func finalizeModelFile(versionDir string, collections []NOSQLCollection) {

	sort.Sort(SchemaNameSorter(collections))

	modelToWrite += "func resolveEntity(key string) modelEntity{\n\n"
	modelToWrite += "switch key{\n"

	for _, collection := range collections {
		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "\":\n"
		modelToWrite += " return &" + strings.Title(collection.Schema.Name) + "{}\n"

		modelToWrite += "case \"" + strings.Title(collection.Schema.Name) + "HistoryRecord\":\n"
		modelToWrite += " return &" + strings.Title(collection.Schema.Name) + "HistoryRecord{}\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"

	modelToWrite += "\n"
	modelToWrite += "func resolveHistoryCollection(key string) modelCollection{\n\n"
	modelToWrite += "switch key{\n"

	for _, collection := range collections {
		modelToWrite += "case \"" + strings.Title(collection.Name) + "History\":\n"
		modelToWrite += " return &" + strings.Title(collection.Name) + "History{}\n"
	}

	modelToWrite += "}\n"
	modelToWrite += "return nil\n"
	modelToWrite += "}\n\n"

	writeNoSQLStub(modelToWrite, serverSettings.APP_LOCATION+"/models/"+versionDir+"/model/model.go")
}

func generateNoSQLModel(schema NOSQLSchema, collection NOSQLCollection, driver string, schemasCreated *[]NOSQLSchema) string {

	val := ""

	timeImport := ""
	if checkSchemaForDateTime(schema) {
		timeImport = "time"
	}

	switch driver {
	case DATABASE_DRIVER_BOLTDB:
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "github.com/asdine/storm", timeImport})
	case DATABASE_DRIVER_MONGODB:
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "gopkg.in/mgo.v2", "gopkg.in/mgo.v2/bson", "log", timeImport, "errors", "encoding/base64", "reflect"})
		// val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "gopkg.in/mgo.v2/bson", "log", "time"})
	}

	val += genNoSQLCollection(collection, driver)
	val += genNoSQLSchema(schema, driver, schemasCreated, 0)
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

func genNoSQLCollection(collection NOSQLCollection, driver string) string {
	val := ""
	val += "type " + strings.Title(collection.Name) + " struct{}\n\n"

	if driver == DATABASE_DRIVER_MONGODB {

		val += "var mongo" + strings.Title(collection.Name) + "Collection *mgo.Collection\n\n"

		val += "func init(){\n"
		val += "go func() {\n\n"
		val += "for{\n"
		val += "if dbServices.MongoDB != nil {\n"
		val += "log.Println(\"Building Indexes for MongoDB collection " + collection.Name + ":\")\n"
		val += "mongo" + strings.Title(collection.Name) + "Collection = dbServices.MongoDB.C(\"" + collection.Name + "\")\n"
		val += "ci := mgo.CollectionInfo{ForceIdIndex: true}\n"
		val += "mongo" + strings.Title(collection.Name) + "Collection.Create(&ci)\n"
		val += "var obj " + strings.Title(collection.Name) + "\n"
		val += "obj.Index()\n"
		val += "obj.Bootstrap()\n"
		val += "return\n"
		val += "}\n"
		val += "<- dbServices.WaitForDatabase()\n"
		val += "}\n"
		val += "}()\n"
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
func genNoSQLSchema(schema NOSQLSchema, driver string, schemasCreated *[]NOSQLSchema, seed int) string {

	val := ""
	schemasToCreate := []NOSQLSchema{}

	if hasGeneratedModelSchema(schema.Name, schemasCreated) == true {
		log.Println("Skipping duplicate schema:  " + schema.Name)
		return ""
	}

	*schemasCreated = append(*schemasCreated, schema)

	val += "type " + strings.Title(schema.Name) + " struct{\n"

	for _, field := range schema.Fields {

		if field.Type == "object" || field.Type == "objectArray" {
			schemasToCreate = append(schemasToCreate, field.Schema)
		}

		additionalTags := genNoSQLAdditionalTags(field, driver)
		omitEmpty := ""
		if field.OmitEmpty {
			omitEmpty = ",omitempty"
		}

		val += "\n\t" + strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t" + genNoSQLFieldType(schema, field, driver) + "\t\t`json:\"" + strings.Title(field.Name) + omitEmpty + "\"" + additionalTags + "`"
	}

	//Add Validation
	if seed == 0 {
		val += "\n"
		val += "Errors struct{\n"

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

		val += "} `json:\"Errors\" bson:\"-\"`\n"
	}

	val += "\n}\n\n"

	for _, schemaToCreate := range schemasToCreate {
		val += genNoSQLSchema(schemaToCreate, driver, schemasCreated, 1)
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

// Checks to see if the schema has been created for a model or not.
func hasGeneratedModelSchema(name string, schemasCreated *[]NOSQLSchema) bool {
	for _, schema := range *schemasCreated {
		if schema.Name == name {
			return true
		}
	}
	return false
}

func genNoSQLAdditionalTags(field NOSQLSchemaField, driver string) string {

	validationTags := ""

	if field.Validation != nil {
		validationTags = "validate:\""

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

		tags := " bson:\"" + strings.Title(field.Name) + "\"" + " " + validationTags
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

func genNoSQLFieldType(schema NOSQLSchema, field NOSQLSchemaField, driver string) string {

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
		return strings.Title(field.Schema.Name)
	case "intArray":
		return "[]int"
	case "float64Array":
		return "[]float64"
	case "stringArray":
		return "[]string"
	case "boolArray":
		return "[]bool"
	case "objectArray":
		return "[]" + strings.Title(field.Schema.Name)
	case "selfArray":
		return "[]" + strings.Title(schema.Name)
	case "self":
		return strings.Title(schema.Name)
	}

	return field.Type
}

func genNoSQLRuntime(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""
	val += genNoSQLSchemaSingle(collection, schema, driver)
	val += genNoSQLSchemaSearch(collection, schema, driver)
	val += genNoSQLSchemaAll(collection, schema, driver)
	val += genNoSQLSchemaAllByIndex(collection, schema, driver)
	val += genNoSQLSchemaRange(collection, schema, driver)
	val += genNoSQLSchemaIndex(collection, schema, driver)
	val += genNoSQLBootstrap(collection, schema, driver)
	val += genNoSQLSchemaRunTransaction(collection, schema, driver)
	val += genNoSQLSchemaNew(collection, schema)
	val += genNoSQLSchemaSave(collection, schema, driver)
	val += genNoSQLSchemaSaveByTran(collection, schema, driver)
	val += genNoSQLValidate(collection, schema, driver)
	val += genNoSQLSchemaDelete(collection, schema, driver)
	val += genNoSQLSchemaDeleteWithTran(collection, schema, driver)
	val += genNoSQLUnMarshal(collection, schema, driver)
	val += genNoSQLSchemaJSONRuntime(schema)

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
		val += "return errors.New(\"Collection " + collection.Name + " not initialized\")\n"
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
		val += "}\n"

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

		val += "//Get the Original Record if it is a Update\n"
		val += "if isUpdate {\n\n"

		val += "	histRecord.Type = TRANSACTION_CHANGETYPE_UPDATE\n"
		val += "	eTransactionNew.changeType = TRANSACTION_CHANGETYPE_UPDATE\n\n"

		val += "	var col " + strings.Title(collection.Name) + "\n"
		val += "	original, err := col.Single(\"id\", self.Id.Hex())\n\n"

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
	val += "return validateFields(" + strings.Title(schema.Name) + "{}, self, reflect.ValueOf(self).Elem())"
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

		val += "if field == \"id\"{\n"
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
		val += "if field == \"id\"{\n"
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
		val += "if field == \"id\"{\n"
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

	val += "func (obj *" + strings.Title(collection.Name) + ") Index() error {\n"
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

	val += "func (obj *" + strings.Title(collection.Name) + ") Bootstrap() error {\n"

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
		val += "for _, doc := range v{\n\n"

		val += "err = doc.Save()\n"
		val += "	if err != nil {\n"
		val += "		log.Println(\"Failed to bootstrap data for " + strings.Title(schema.Name) + ":  \" + doc.Id.Hex() + \"  \" + err.Error())\n"
		val += "	}\n"
		val += "	log.Println(\"Successfully bootstraped " + strings.Title(schema.Name) + ":  \" + doc.Id.Hex())\n\n"
		val += "	log.Printf(\"%+v\\n\", doc)\n\n"
		val += "}\n\n"

	}
	val += "return nil"
	val += "}\n\n"
	return val
}

func genNoSQLSchemaRunTransaction(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") RunTransaction(objects []" + strings.Title(schema.Name) + ") error {\n\n"
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

	val += "func (obj *" + strings.Title(collection.Name) + ") New() *" + strings.Title(schema.Name) + " {\n"
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
		val += "return mongo" + strings.Title(collection.Name) + "Collection.Remove(self)"
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

		val += "var eTransactionNew entityTransaction\n"
		val += "eTransactionNew.changeType = TRANSACTION_CHANGETYPE_DELETE\n"
		val += "eTransactionNew.entity = self\n\n"

		val += "var eTransactionOriginal entityTransaction\n"
		val += "eTransactionOriginal.changeType = TRANSACTION_CHANGETYPE_DELETE\n"
		val += "eTransactionOriginal.entity = &histRecord\n\n"

		val += "var col " + strings.Title(collection.Name) + "\n"
		val += "original, err := col.Single(\"id\", self.Id.Hex())\n\n"

		val += "originalJson, err := original.JSONString()\n\n"

		val += "if err != nil {\n"
		val += "	return err\n"
		val += "}\n\n"

		val += "originalBase64 := getBase64(originalJson)\n"
		val += "histRecord.Data = originalBase64\n\n"

		val += "tPersist.originalItems = append(tPersist.originalItems, eTransactionOriginal)\n"
		val += "tPersist.newItems = append(tPersist.newItems, eTransactionNew)\n"

		val += "return nil"
	}
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
