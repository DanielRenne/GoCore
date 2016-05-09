package dbServices

import (
	"core/extensions"
)

func genSchemaWebAPI(schema NOSQLSchema, dbPackageName string, driver string) string {

	val := extensions.GenPackageImport("webAPI", []string{dbPackageName, "core/ginServer", "encoding/json"})
	val += "func init(){\n\n"
	val += "}\n"
	return val
}
