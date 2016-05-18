package sqliteDBServices

import (
	"encoding/json"
	"fmt"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
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

func RunDBCreate() {

	jsonData, err := ioutil.ReadFile("db/schemas/create.json")
	if err != nil {
		fmt.Println("Reading of create.json failed:  " + err.Error())
		return
	}

	if serverSettings.WebConfig.DbConnection.Driver == "sqlite3" {
		var co createObject
		errUnmarshal := json.Unmarshal(jsonData, &co)
		if errUnmarshal != nil {
			color.Red("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
			return
		}
		if serverSettings.WebConfig.DbConnection.Driver == "sqlite3" {
			createSQLiteTables(co.Tables)
			createSQLiteIndexes(co.Indexes)
			createSQLiteForeignKeys(co.ForeignKeys, co.Tables)
		}

	}
}
