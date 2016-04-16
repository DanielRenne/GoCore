package dbServices

import (
	"core/extensions"
	"core/serverSettings"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"strings"
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

type createObject struct {
	Tables  []tableDef `json:"tables"`
	Indexes []indexDef `json:"indexes"`
}

func RunDBCreate() {

	jsonData, err := ioutil.ReadFile("db/" + serverSettings.WebConfig.DbConnection.AppName + "/create.json")
	if err != nil {
		fmt.Println("Reading of create.json failed:  " + err.Error())
		return
	}

	var co createObject
	errUnmarshal := json.Unmarshal(jsonData, &co)
	if errUnmarshal != nil {
		fmt.Println("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
		return
	}
	if serverSettings.WebConfig.DbConnection.Driver == "sqlite3" {
		createSQLiteTables(co.Tables)
		createSQLiteIndexes(co.Indexes)
	}

	// fmt.Printf("%+v\n", co)
}

func createSQLiteTables(tables []tableDef) {

	for _, table := range tables {
		sqlStmt := "create table if not exists "
		sqlStmt += table.Name + " "
		sqlStmt += "("

		for _, field := range table.Fields {
			sqlStmt += field.Name + " "
			sqlStmt += field.FieldType + " "
			if field.AllowNull == false {
				sqlStmt += "not null "
			}
			if field.Primary == true {
				sqlStmt += "primary key "
			}
			if field.IsUnique == true {
				sqlStmt += "unique "
			}
			if field.Check != "" {
				sqlStmt += "check(" + field.Check + ") "
			}
			if field.Collate != "" {
				sqlStmt += "collate " + field.Collate + " "
			}
			if field.Default != "" {
				fieldChar := getFieldCharacter(field.FieldType)
				sqlStmt += "Default " + fieldChar + field.Default + fieldChar + " "
			}

			sqlStmt += ","
		}

		sqlStmt = extensions.TrimSuffix(sqlStmt, ",")

		sqlStmt += ");"

		fmt.Println(sqlStmt)

		_, errDBExec := DB.Exec(sqlStmt)
		if errDBExec != nil {
			fmt.Println("Creation of table \"" + table.Name + "\" failed:  " + errDBExec.Error())
			continue
		}

		fmt.Println("Creation of table \"" + table.Name + "\" successful.")
	}
}

func createSQLiteIndexes(indexes []indexDef) {
	for _, index := range indexes {
		sqlStmt := "create "
		if index.IsUnique == true {
			sqlStmt += "unique "
		}

		sqlStmt += "index if not exists "
		sqlStmt += index.Name + " on "
		sqlStmt += index.TableName + " ("

		for _, field := range index.Fields {
			sqlStmt += field
			sqlStmt += ","
		}

		sqlStmt = extensions.TrimSuffix(sqlStmt, ",")

		sqlStmt += ");"

		_, errDBExec := DB.Exec(sqlStmt)
		if errDBExec != nil {
			fmt.Println("Creation of index \"" + index.Name + "\" failed:  " + errDBExec.Error())
			continue
		}

		fmt.Println("Creation of index \"" + index.Name + "\" successful.")
	}
}

func getFieldCharacter(value string) string {
	if strings.Contains(value, "CHAR") || strings.Contains(value, "TEXT") || strings.Contains(value, "CLOB") {
		return "'"
	}

	return ""
}
