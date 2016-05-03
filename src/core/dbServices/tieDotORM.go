package dbServices

import (
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
)

func createTieDotORM(collections []NOSQLCollection, packageName string) {

	for _, collection := range collections {
		val := generateTieDotORM(collection.Schema)
		os.Mkdir("src/"+packageName+"/orm", 0777)
		writeTieDotORMCollection(val, "src/"+packageName+"/orm/"+collection.Name+".go")
		color.Green("Created NOSQL Collection " + collection.Name + " successfully.")
	}
}

func writeTieDotORMCollection(value string, path string) {

	err := ioutil.WriteFile(path, []byte(value), 0777)
	if err != nil {
		color.Red("Error creating ORM for TieDot Collection:  " + err.Error())
	}
}

func generateTieDotORM(schema NOSQLSchema) string {

	val := "package orm\n\n"
	val += "import(\n"
	val += "\t\"core/dbServices\"\n"
	val += "\t\"fmt\"\n"
	val += ")\n\n"

	val += genTieDotSchema(schema)

	return val
}

func genTieDotSchema(schema NOSQLSchema) string {

	val := ""
	schemasToCreate := []NOSQLSchema{}

	val += "type " + strings.Title(schema.Name) + " struct{\n"
	for _, field := range schema.Fields {
		if field.Type == "object" && field.Type == "objectArray" {
			schemasToCreate = append(schemasToCreate, field.Schema)
		}

		val += "\n\t" + strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t" + genTieDotFieldType(field) + "\t\t`json:\"" + field.Name + "\"`"
	}

	val += "\n}\n"
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
		return field.Schema.Name
	case "intArray":
		return "[]int"
	case "floatArray":
		return "[]float64"
	case "stringArray":
		return "[]string"
	case "boolArray":
		return "[]bool"
	case "objectArray":
		return "[]" + field.Schema.Name
	}
	return ""
}
