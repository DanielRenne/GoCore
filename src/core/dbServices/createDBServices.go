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

type foreignKeyDef struct {
	Name     string `json:"name"`
	Field    string `json:"field"`
	FKTable  string `json:"fkTable"`
	FKField  string `json:"fkField"`
	IsDelete bool   `json:"isDelete"`
	IsUpdate bool   `json:"isUpdate"`
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
		createSQLiteForeignKeys(co.ForeignKeys, co.Tables)
	}

	// fmt.Printf("%+v\n", co)
}

func createSQLiteTables(tables []tableDef) {

	for _, table := range tables {

		sqlStmt := generateSQLiteTableCreate(table, []foreignKeyDef{})

		fmt.Println(sqlStmt)

		_, errDBExec := DB.Exec(sqlStmt)
		if errDBExec != nil {
			fmt.Println("Creation of table \"" + table.Name + "\" failed:  " + errDBExec.Error())
			continue
		}

		fmt.Println("Creation of table \"" + table.Name + "\" successful.")
	}
}

func generateSQLiteTableCreate(table tableDef, keys []foreignKeyDef) string {
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

	for _, fk := range keys {
		cascades := ""
		if fk.IsDelete == true {
			cascades += " ON DELETE CASCADE"
		}
		if fk.IsUpdate == true {
			cascades += " ON UPDATE CASCADE"
		}
		sqlStmt += ", FOREIGN KEY (" + fk.Field + ") references " + fk.FKTable + "(" + fk.FKField + ") " + cascades

	}

	sqlStmt += ");"
	fmt.Println(sqlStmt)

	return sqlStmt
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

func createSQLiteForeignKeys(foreignKeys []foreignKeyTableDef, tables []tableDef) {
	for _, fk := range foreignKeys {

		renameSQLiteTable(fk.Table, fk.Table+"_temp")

		//Create the new Table with FK's
		for _, table := range tables {
			if table.Name == fk.Table {
				sqlStmt := generateSQLiteTableCreate(table, fk.Keys)
				_, errDBExec := DB.Exec(sqlStmt)
				if errDBExec != nil {
					fmt.Println("Creation of table with Foreign Keys Failed:  Table  \"" + fk.Table + "_temp" + "\" " + errDBExec.Error())
					return
				}
				fmt.Println("Creation of table \"" + fk.Table + "\" successful.")
				break
			}
		}

		copySQLiteTable(fk.Table+"_temp", fk.Table)
		dropSQLiteTable(fk.Table + "_temp")
	}
}

func renameSQLiteTable(original string, newName string) {

	sqlStmt := "ALTER TABLE " + original + " RENAME TO " + newName + ";"

	fmt.Println("Renaming table from " + original + " to " + newName)

	_, errDBExec := DB.Exec(sqlStmt)
	if errDBExec != nil {
		fmt.Println("Rename of table \"" + original + "\" failed:  " + errDBExec.Error())
	}
	fmt.Println("Renamed " + original + " to " + newName)
}

func copySQLiteTable(from string, to string) {

	sqlStmt := "INSERT INTO " + to + " SELECT * FROM " + from + ";"

	fmt.Println("Copying Table Data from " + from + " to " + to)

	_, errDBExec := DB.Exec(sqlStmt)
	if errDBExec != nil {
		fmt.Println("Copy of table \"" + from + "\" to " + "\"" + to + "\" failed:  " + errDBExec.Error())
	}
	fmt.Println("Copied Data from " + from + " to " + to + " successfully.")
}

func dropSQLiteTable(tableName string) {

	sqlStmt := "DROP TABLE IF EXISTS " + tableName + ";"

	fmt.Println("Dropping table " + tableName)

	_, errDBExec := DB.Exec(sqlStmt)
	if errDBExec != nil {
		fmt.Println("Drop of Table \"" + tableName + "\" failed:  " + errDBExec.Error())
	}
	fmt.Println("Dropped Table " + tableName + " successfully.")
}

func getFieldCharacter(value string) string {
	if strings.Contains(value, "CHAR") || strings.Contains(value, "TEXT") || strings.Contains(value, "CLOB") {
		return "'"
	}

	return ""
}
