package sqliteDBServices

import (
	"strings"
)

func createSQLiteORM(tables []tableDef, packageName string) {

}

func generateSQLiteORM(table tableDef, packageName string) string {

	val := "package " + packageName + "\n\n"
	val += "import(\n"
	val += "\t\"database/sql\"\n"
	val += "\t\"fmt\"\n"
	val += "\t_ \"github.com/mattn/go-sqlite3\"\n"
	val += ")\n\n"

	val += genSQLiteTable(table)

	return val
}

func genSQLiteTable(table tableDef) string {
	val := "type " + table.Name + " struct{\n"

	for _, field := range table.Fields {

		fieldType := getSQLiteFieldType(field.FieldType)

		val += strings.Replace(strings.Title(field.Name), " ", "_", -1) + "\t\t\t\t" + fieldType + "\n"
	}

	val += "}\n"
	return val
}

func getSQLiteFieldType(fieldType string) string {
	val := "int"
	switch fieldType {
	case "INTEGER":
		val = "int"
	case "INT":
		val = "int"
	case "TINYINT":
		val = "int"
	case "SMALLINT ":
		val = "int"
	case "MEDIUMINT":
		val = "int"
	case "BIGINT":
		val = "int"
	case "BOOLEAN":
		val = "bool"
	case "CHAR":
		val = "string"
	case "TEXT":
		val = "string"
	case "STRING":
		val = "string"
	case "VARCHAR":
		val = "string"
	case "DECIMAL":
		val = "float64"
	case "DOUBLE":
		val = "float64"
	case "REAL":
		val = "float64"
	case "BLOB":
		val = "[]byte"
	case "NONE":
		val = "[]byte"
	}

	return val
}
