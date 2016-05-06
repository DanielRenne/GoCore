package sqliteDBServices

import (
	"reflect"
	"testing"
)

func TestIsSqliteFKChangeRequired(t *testing.T) {
	// t.Error("Testing Error")
	fkDef := foreignKeyDef{Name: "Location",
		Fields:   []string{"LocationId"},
		FKTable:  "Customer",
		FKFields: []string{"LocationId"},
		OnDelete: true,
		OnUpdate: true}

	combinedSchemas := []sqliteForeignKeyCombinedSchema{}
	combinedSchemas = append(combinedSchemas, sqliteForeignKeyCombinedSchema{
		Table:     "Customer",
		From:      "LocationId",
		To:        "LocationId",
		On_Update: "CASCADE",
		On_Delete: "CASCADE",
		Match:     "",
	})

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == true {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed to match foreignKeyDef to mocked Schema.")
	}

	fkDef.OnDelete = false

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == false {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed when OnDelete is Mismatched.")
	}

	fkDef.OnDelete = true

	fkDef.OnUpdate = false

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == false {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed when OnUpdate is Mismatched.")
	}

	fkDef.OnUpdate = true

	fkDef.Fields[0] = ""

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == false {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed when Table Field is Mismatched.")
	}

	fkDef.Fields[0] = "LocationId"

	fkDef.FKFields[0] = ""

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == false {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed when FK Field is Mismatched.")
	}

	fkDef.FKFields[0] = "LocationId"

	//Testing a multiField Foreign Key
	combinedSchemas[0].From += ",TestId"
	combinedSchemas[0].To += ",TestId"

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == false {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed when Multiple Foriegn Key fields are present and does not match.")
	}

	fkDef.Fields = append(fkDef.Fields, "TestId")
	fkDef.FKFields = append(fkDef.FKFields, "TestId")

	if isSqliteFKChangeRequired(fkDef, combinedSchemas) == true {
		t.Error("Failed isSqliteFKChangeRequired at sqliteDDL.go:  Failed when Multiple Foriegn Key fields are present and match.")
	}
}

func TestCombineSQLiteForeignKeys(t *testing.T) {

	FKSchemas := []sqliteForeignKeySchema{}
	FKSchemas = append(FKSchemas, sqliteForeignKeySchema{
		Id:        0,
		Seq:       0,
		Table:     "Customer",
		From:      "LocationId",
		To:        "LocationId",
		On_Update: "CASCADE",
		On_Delete: "CASCADE",
		Match:     "",
	})

	FKSchemas = append(FKSchemas, sqliteForeignKeySchema{
		Id:        0,
		Seq:       1,
		Table:     "Customer",
		From:      "Test",
		To:        "Test",
		On_Update: "CASCADE",
		On_Delete: "CASCADE",
		Match:     "",
	})

	validateCombined := []sqliteForeignKeyCombinedSchema{}
	validateCombined = append(validateCombined, sqliteForeignKeyCombinedSchema{
		Table:     "Customer",
		From:      "LocationId,Test",
		To:        "LocationId,Test",
		On_Update: "CASCADE",
		On_Delete: "CASCADE",
		Match:     "",
	})

	validateResult := combineSQLiteForeignKeys(FKSchemas)
	if reflect.DeepEqual(validateCombined, validateResult) == false {
		t.Error("Failed combineSQLiteForeignKeys at sqliteDDL.go.")
	}
}
