package dbServices

import (
	"encoding/json"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	// "fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

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
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Index        string      `json:"index"`
	DefaultValue string      `json:"defaultValue"`
	Required     bool        `json:"required"`
	Schema       NOSQLSchema `json:"schema"`
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

	basePath := "db/" + serverSettings.WebConfig.DbConnection.AppName + "/schemas"

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
	extensions.CopyFile(serverSettings.SwaggerUIPath+"/swagger."+versionNumber+".json", serverSettings.SwaggerUIPath+"/swagger.json")

}

func walkNoSQLVersion(path string, versionDir string) {

	var wg sync.WaitGroup

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
					color.Red("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
					e = errUnmarshal
				}

				createNoSQLModel(schemaDB.Collections, serverSettings.WebConfig.DbConnection.AppName, serverSettings.WebConfig.DbConnection.Driver, versionDir)
			}()
		}

		return e
	})
	if err != nil {
		color.Red("Walk of path failed:  " + err.Error())
	}

	wg.Wait()
}

func createNoSQLModel(collections []NOSQLCollection, packageName string, driver string, versionDir string) {

	//Create a NOSQLBucket Model
	bucket := generateNoSQLModelBucket(driver)
	os.Mkdir("src/"+packageName+"/models/", 0777)
	os.Mkdir("src/"+packageName+"/models/"+versionDir, 0777)
	os.Mkdir("src/"+packageName+"/models/"+versionDir+"/model/", 0777)

	writeNOSQLModelBucket(bucket, "src/"+packageName+"/models/"+versionDir+"/model/bucket.go")

	for _, collection := range collections {
		val := generateNoSQLModel(collection.Schema, collection, driver)
		os.Mkdir("src/"+packageName+"/models/", 0777)
		os.Mkdir("src/"+packageName+"/models/"+versionDir, 0777)
		os.Mkdir("src/"+packageName+"/models/"+versionDir+"/model/", 0777)
		writeNoSQLModelCollection(val, "src/"+packageName+"/models/"+versionDir+"/model/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+".go", collection)

		os.Mkdir("src/"+packageName+"/webAPIs/", 0777)
		os.Mkdir("src/"+packageName+"/webAPIs/"+versionDir, 0777)
		os.Mkdir("src/"+packageName+"/webAPIs/"+versionDir+"/webAPI/", 0777)

		cWebAPI := genSchemaWebAPI(collection, collection.Schema, packageName+"/models/"+versionDir+"/model", driver, versionDir)
		writeNoSQLWebAPI(cWebAPI, "src/"+packageName+"/webAPIs/"+versionDir+"/webAPI/"+extensions.MakeFirstLowerCase(collection.Schema.Name)+".go", collection)
	}

}

func generateNoSQLModel(schema NOSQLSchema, collection NOSQLCollection, driver string) string {

	val := ""
	switch driver {
	case "boltDB":
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices", "encoding/json", "github.com/asdine/storm"})
	}

	val += genNoSQLCollection(collection)
	val += genNoSQLSchema(schema, driver)
	val += genNoSQLRuntime(collection, schema, driver)
	return val
}

func generateNoSQLModelBucket(driver string) string {
	val := ""
	switch driver {
	case "boltDB":
		val += extensions.GenPackageImport("model", []string{"github.com/DanielRenne/GoCore/core/dbServices"})
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

func genNoSQLCollection(collection NOSQLCollection) string {
	val := ""
	val += "type " + strings.Title(collection.Name) + " struct{}\n\n"
	return val
}

func genNoSQLSchema(schema NOSQLSchema, driver string) string {

	val := ""
	schemasToCreate := []NOSQLSchema{}

	val += "type " + strings.Title(schema.Name) + " struct{\n"

	for _, field := range schema.Fields {

		if field.Type == "object" || field.Type == "objectArray" {
			schemasToCreate = append(schemasToCreate, field.Schema)
		}

		additionalTags := genNoSQLAdditionalTags(field, driver)

		val += "\n\t" + strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t" + genNoSQLFieldType(field) + "\t\t`json:\"" + extensions.MakeFirstLowerCase(field.Name) + "\"" + additionalTags + "`"
	}

	val += "\n}\n\n"

	for _, schemaToCreate := range schemasToCreate {
		val += genNoSQLSchema(schemaToCreate, driver)
	}

	return val
}

func genNoSQLAdditionalTags(field NOSQLSchemaField, driver string) string {
	switch driver {
	case "boltDB":
		{
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
		}
	}
	return ""
}

func genNoSQLFieldType(field NOSQLSchemaField) string {

	switch field.Type {
	case "int":
		return "uint64"
	case "float64":
		return "float64"
	case "string":
		return "string"
	case "bool":
		return "bool"
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
	val += genNoSQLSchemaRunTransaction(collection, schema, driver)
	val += genNoSQLSchemaNew(collection, schema)
	val += genNoSQLSchemaSave(schema, driver)
	val += genNoSQLSchemaDelete(schema, driver)
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

func genNoSQLSchemaSave(schema NOSQLSchema, driver string) string {
	val := ""
	val += "func (obj *" + strings.Title(schema.Name) + ") Save() error {\n"
	switch driver {
	case "boltDB":
		{
			val += "return dbServices.BoltDB.Save(obj)\n"
		}
	}
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

	val += "func (obj *" + strings.Title(collection.Name) + ") Single(field string, value interface{}) (retObj " + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
			val += "e = dbServices.BoltDB.One(field, value, &retObj)\n"
			val += "return\n"
		}
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaSearch(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") Search(field string, value interface{}) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
			val += "e = dbServices.BoltDB.Find(field, value, &retObj)\n"
			val += genNoSQLSchemaArrayCheck(schema)
			val += "return\n"
		}
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
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
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaAll(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") All() (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
			val += "e = dbServices.BoltDB.All(&retObj)\n"
			val += genNoSQLSchemaArrayCheck(schema)
			val += "return\n"
		}
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") AllAdvanced(limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
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
	}
	val += "}\n\n"

	return val
}

func genNoSQLSchemaAllByIndex(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") AllByIndex(index string) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
			val += "e = dbServices.BoltDB.AllByIndex(index, &retObj)\n"
			val += genNoSQLSchemaArrayCheck(schema)
			val += "return\n"
		}
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") AllByIndexAdvanced(index string, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
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
	}
	val += "}\n\n"

	return val
}

func genNoSQLSchemaRange(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") Range(min, max, field string) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
			val += "e = dbServices.BoltDB.Range(field, min, max, &retObj)\n"
			val += genNoSQLSchemaArrayCheck(schema)
			val += "return\n"
		}
	}
	val += "}\n\n"

	val += "func (obj *" + strings.Title(collection.Name) + ") RangeAdvanced(min, max, field string, limit int, skip int) (retObj []" + strings.Title(schema.Name) + ",e error) {\n"
	switch driver {
	case "boltDB":
		{
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
		}
	}
	val += "}\n\n"

	return val
}

func genNoSQLSchemaIndex(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") Index() error {\n"
	switch driver {
	case "boltDB":
		{
			val += "return dbServices.BoltDB.Init(&" + strings.Title(schema.Name) + "{})\n"
		}
	}
	val += "}\n\n"
	return val
}

func genNoSQLSchemaRunTransaction(collection NOSQLCollection, schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(collection.Name) + ") RunTransaction(objects []" + strings.Title(schema.Name) + ") error {\n\n"
	switch driver {
	case "boltDB":
		{

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
		}
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

func genNoSQLSchemaDelete(schema NOSQLSchema, driver string) string {
	val := ""

	val += "func (obj *" + strings.Title(schema.Name) + ") Delete() error {\n"
	switch driver {
	case "boltDB":
		{
			val += "return dbServices.BoltDB.Remove(obj)\n"
		}
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
		{
			val += "return dbServices.BoltDB.Set(obj.Name, key, value)\n"
		}
	}
	val += "}\n\n"

	val += "func (obj *Bucket) GetKeyValue(key interface{}, value interface{}) error {\n"
	switch driver {
	case "boltDB":
		{
			val += "return dbServices.BoltDB.Get(obj.Name, key, value)\n"
		}
	}
	val += "}\n\n"

	val += "func (obj *Bucket) DeleteKey(key interface{}) error {\n"
	switch driver {
	case "boltDB":
		{
			val += "return dbServices.BoltDB.Delete(obj.Name, key)\n"
		}
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
