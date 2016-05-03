package dbServices

import (
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func createTieDotORM(collections []NOSQLCollection, packageName string) {

	//Create ORM
	ormVal := generateORMPackage()

	os.Mkdir("src/"+packageName+"/orm", 0777)
	writeTieDotModelCollection(ormVal, "src/"+packageName+"/orm/orm.go")
	color.Green("Created orm package successfully.")

	cmd := exec.Command("gofmt", "-w", "src/"+packageName+"/orm/orm.go")
	err := cmd.Start()
	if err != nil {
		color.Red("Failed to gofmt on file src/" + packageName + "/orm/orm.go:  " + err.Error())
	}

	//Create Model
	for _, collection := range collections {
		val := generateTieDotORM(collection.Schema, collection, packageName)
		os.Mkdir("src/"+packageName+"/model", 0777)
		writeTieDotModelCollection(val, "src/"+packageName+"/model/"+collection.Schema.Name+".go")
		color.Green("Created NOSQL Collection " + collection.Name + " successfully.")
	}

}

func generateORMPackage() string {
	val := genPackageImport("orm", []string{"core/dbServices", "github.com/HouzuoGuo/tiedot/db"})
	val += "func GetCollection(collectionName string) *db.Col{\n"
	val += "\treturn dbServices.TiedotDB.Use(collectionName)\n"
	val += "}\n\n"
	return val
}

func writeTieDotModelCollection(value string, path string) {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		color.Red("Error creating ORM for TieDot Collection:  " + err.Error())
	}

	cmd := exec.Command("gofmt", "-w", path)
	err = cmd.Start()
	if err != nil {
		color.Red("Failed to gofmt on file " + path + ":  " + err.Error())
	}
}

func generateTieDotORM(schema NOSQLSchema, collection NOSQLCollection, packageName string) string {

	// val := genPackageImport("model", []string{packageName + "/orm", "reflect", "encoding/json"})
	val := genPackageImport("model", []string{packageName + "/orm", "encoding/json"})

	val += genTieDotCollection(collection)
	val += genTieDotSchema(schema)
	val += genTieDotRuntime(collection, schema)
	return val
}

func genTieDotCollection(collection NOSQLCollection) string {
	val := ""
	val += "type " + strings.Title(collection.Name) + " struct{}\n\n"
	return val
}

func genTieDotSchema(schema NOSQLSchema) string {

	val := ""
	schemasToCreate := []NOSQLSchema{}

	val += "type " + strings.Title(schema.Name) + " struct{\n"
	val += "\n\tDocumentId\tint"
	for _, field := range schema.Fields {
		if field.Type == "object" || field.Type == "objectArray" {
			schemasToCreate = append(schemasToCreate, field.Schema)
		}

		val += "\n\t" + strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t" + genTieDotFieldType(field) + "\t\t`json:\"" + field.Name + "\"`"
	}

	val += "\n}\n\n"

	for _, schemaToCreate := range schemasToCreate {
		val += genTieDotSchema(schemaToCreate)
	}

	return val
}

func genTieDotFieldType(field NOSQLSchemaField) string {

	switch field.Type {
	case "int":
		return "int"
	case "float":
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
	case "floatArray":
		return "[]float64"
	case "stringArray":
		return "[]string"
	case "boolArray":
		return "[]bool"
	case "objectArray":
		return "[]" + strings.Title(field.Schema.Name)
	}
	return ""
}

func genTieDotRuntime(collection NOSQLCollection, schema NOSQLSchema) string {
	val := ""
	val += genTieDotCollectionInsert(collection, strings.Title(schema.Name))
	// val += genTieDotPersist(collection, schema)
	// val += genTieDotIsDirty(schema)
	return val
}

func genTieDotCollectionInsert(collection NOSQLCollection, schemaName string) string {
	val := ""
	val += "func (obj *" + strings.Replace(strings.Title(collection.Name), " ", "_", -1) + ") Insert(item " + schemaName + ") (int, error){\n"
	val += "itemBytes, _ := json.Marshal(item)\n"
	val += collection.Name + " := orm.GetCollection(\"" + collection.Name + "\")\n"
	val += "var docMap map[string]interface{}\n"
	val += "err := json.Unmarshal(itemBytes, &docMap)\n"
	val += "if err != nil {\n return 0, err\n}\n\n"
	val += "return " + collection.Name + ".Insert(docMap)"
	val += "}\n\n"
	return val
}

func genTieDotPersist(collection NOSQLCollection, schema NOSQLSchema) string {
	val := ""
	val += "func (obj *" + strings.Replace(strings.Title(schema.Name), " ", "_", -1) + ") Persist() (int, error){\n"
	val += "}\n\n"
	return val
}

func genTieDotIsDirty(schema NOSQLSchema) string {
	val := ""
	val += "func (obj *" + strings.Replace(strings.Title(schema.Name), " ", "_", -1) + ") IsDirty() (bool, error){\n"
	val += "}\n\n"
	return val
}

func genPackageImport(name string, imports []string) string {

	val := "package " + name + "\n\n"
	val += "import(\n"
	for _, imp := range imports {
		val += "\t\"" + imp + "\"\n"
	}
	val += ")\n\n"

	return val
}
