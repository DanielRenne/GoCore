package dbServices

import (
	"core/extensions"
	"core/serverSettings"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
)

type fieldDef struct {
	Name      string `json:"name"`
	Primary   bool   `json:"primary"`
	AllowNull bool   `json:"allowNull"`
	FieldType string `json:"fieldType"`
	IsUnique  bool   `json:"isUnique"`
	Check  	  string `json:"check"`
	Collate   string `json:"collate"`	
	Default	  string `json:"default"`
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

var DB *sql.DB

func init() {
	fmt.Println("core dbServices initialized.")

	var err error
	DB, err = sql.Open(serverSettings.WebConfig.DbConnection.Driver, serverSettings.WebConfig.DbConnection.ConnectionString)

	if err != nil {
		fmt.Println("Open connection failed:" + err.Error())
		return
	}

	fmt.Println("Open Database Connections: " + string(DB.Stats().OpenConnections))

	// sqlStmt := `
	// create table foo (id integer not null primary key, name text);
	// delete from foo;
	// `
	// _, err = DB.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	// tx, err := DB.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// tx.Commit()

	// rows, err := DB.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	rows.Scan(&id, &name)
	// 	fmt.Println(id, name)
	// }

	// stmt, err = DB.Prepare("select name from foo where id = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// var name string
	// err = stmt.QueryRow("3").Scan(&name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(name)

	// _, err = DB.Exec("delete from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = DB.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err = DB.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	rows.Scan(&id, &name)
	// 	fmt.Println(id, name)
	// }

	runDBCreate()

}

func runDBCreate() {

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

func getFieldCharacter(value string) string{
	if strings.Contains(value, "CHAR") || strings.Contains(value, "TEXT") || strings.Contains(value, "CLOB"){
		return "'"
	}

	return ""
}
