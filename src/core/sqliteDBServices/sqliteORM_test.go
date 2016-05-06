package sqliteDBServices

import (
	// "io/ioutil"
	"fmt"
	"testing"
)

func TestGenerateSQLiteORM(t *testing.T) {

	// ormFile, err := ioutil.ReadFile("core/dbServices/testFiles/sqlite/testORM.go")

	// if err != nil {
	// 	t.Error("Reading of testFiles/sqlite/testORM.go failed:  " + err.Error())
	// 	return
	// }

	fields := []fieldDef{}
	fields = append(fields, fieldDef{
		Name:      "CustomerId",
		Primary:   true,
		AllowNull: false,
		FieldType: "INTEGER",
		IsUnique:  true,
		Check:     "",
		Collate:   "",
		Default:   "",
	})

	table := tableDef{Name: "Customer", Fields: fields}

	compareResult := generateSQLiteORM(table, "testORM")

	fmt.Println(compareResult)

	if getTestORM() != compareResult {
		t.Error("Failed generateSQLiteORM at sqliteORM.go:  Failed to match getTestORM().")
	}
}

func getTestORM() string {
	val := "package testORM\n\n"
	val += "import(\n"
	val += "\t\"database/sql\"\n"
	val += "\t\"fmt\"\n"
	val += "\t_ \"github.com/mattn/go-sqlite3\"\n"
	val += ")\n"
	return val
}
