package dbServices

import (
	"core/extensions"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type tableSchema struct {
	CId          int
	Name         string
	Type         string
	NotNull      int
	DefaultValue sql.NullString
	PrimaryKey   int
}

func createSQLiteTables(tables []tableDef) {

	for _, table := range tables {

		doWeAlter := isAlterRequired(table, getSQLiteTableSchema(table))
		fmt.Printf("%+v\n", doWeAlter)
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

func isAlterRequired(table tableDef, existingSchema []tableSchema) bool {

	return false
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
			fieldChar := getSQLiteFieldCharacter(field.FieldType)
			sqlStmt += "Default " + fieldChar + field.Default + fieldChar + " "
		}

		sqlStmt += ","
	}

	sqlStmt = extensions.TrimSuffix(sqlStmt, ",")

	sqlStmt += genSQLiteForeignKeyCreate(keys)

	sqlStmt += ");"
	fmt.Println(sqlStmt)

	return sqlStmt
}

func genSQLiteForeignKeyCreate(keys []foreignKeyDef) string {

	sqlStmt := ""
	for _, fk := range keys {
		cascades := ""
		if fk.OnDelete == true {
			cascades += " ON DELETE CASCADE"
		}
		if fk.OnUpdate == true {
			cascades += " ON UPDATE CASCADE"
		}

		tblFields := ""
		for _, f := range fk.Fields {
			tblFields += f + ","
		}

		tblFields = extensions.TrimSuffix(tblFields, ",")

		fkFields := ""
		for _, f := range fk.FKFields {
			fkFields += f + ","
		}

		fkFields = extensions.TrimSuffix(fkFields, ",")

		sqlStmt += ", FOREIGN KEY (" + tblFields + ") references " + fk.FKTable + "(" + fkFields + ") " + cascades
	}

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

func getSQLiteFieldCharacter(value string) string {
	if strings.Contains(value, "CHAR") || strings.Contains(value, "TEXT") || strings.Contains(value, "CLOB") {
		return "'"
	}

	return ""
}

func getSQLiteTableSchema(table tableDef) []tableSchema {

	schemaRows := []tableSchema{}

	rows, err := DB.Query("PRAGMA table_info(" + table.Name + ");")
	if err != nil {
		fmt.Println("Prepare of DB Query Failed.  " + err.Error())
		return schemaRows
	}
	defer rows.Close()

	for rows.Next() {
		var schema tableSchema
		err = rows.Scan(&schema.CId, &schema.Name, &schema.Type, &schema.NotNull, &schema.DefaultValue, &schema.PrimaryKey)
		if err != nil {
			fmt.Println("Scan failed to load Table Schema for table \"" + table.Name + "\":  " + err.Error())
		}
		schemaRows = append(schemaRows, schema)
	}
	return schemaRows
}
