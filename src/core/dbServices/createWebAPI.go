package dbServices

import (
	"core/extensions"
	"strings"
)

const basePath = "/api"

func genSchemaWebAPI(collection NOSQLCollection, schema NOSQLSchema, dbPackageName string, driver string, versionDir string) string {

	val := extensions.GenPackageImport("webAPI", []string{dbPackageName, "core/ginServer", "github.com/gin-gonic/gin", "core/extensions"})
	val += "func init(){\n\n"
	//val += "\tginServer.AddRouterGroup(\"" + versionDir + "\", \"/" + strings.ToLower(schema.Name) + "\", \"GET\", get" + strings.Title(schema.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/single" + strings.Title(collection.Name) + "\", \"GET\", getSingle" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/search" + strings.Title(collection.Name) + "\", \"GET\", getSearch" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/sort" + strings.Title(collection.Name) + "\", \"GET\", getSort" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/range" + strings.Title(collection.Name) + "\", \"GET\", getRange" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/" + strings.ToLower(collection.Name) + "\", \"GET\", get" + strings.Title(collection.Name) + ")\n"
	val += "}\n\n"

	val += genSingleCollectionGET(collection)
	val += genSearchCollectionGET(collection)
	val += genSortCollectionGET(collection)
	val += genRangeCollectionGET(collection)
	val += genCollectionGET(collection)

	//Add Swagger Paths
	addNOSQLSwaggerCollectionGet("/"+strings.ToLower(collection.Name), collection)
	addNOSQLSwaggerSearchCollectionGet("/search"+strings.Title(collection.Name), collection)

	addNOSQLSwaggerSchemaDefinition(schema)

	return val
}

func genCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func get" + name + "(c *gin.Context){\n"

	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray := items.AllAdvanced(limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray := items.All()\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

//Adds to the Swagger.json root the Definitions
func addNOSQLSwaggerSchemaDefinition(schema NOSQLSchema) {
	var def Swagger2Schema
	requiredProperties := []string{}

	def.Type = "object"

	// if len(schema.Fields) > 0 {
	// 	def.Properties = make(map[string]Swagger2Schema)
	// }

	// for _, field := range schema.Fields {

	// 	if field.Required == true {
	// 		requiredProperties = append(requiredProperties, field.Name)
	// 	}

	// 	fieldSwaggerSchema := Swagger2Schema{}
	// 	if field.Type == "object" {
	// 		fieldSwaggerSchema.Ref = "#/definitions/" + strings.ToLower(field.Schema.Name)
	// 		addNOSQLSwaggerSchemaDefinition(field.Schema)
	// 	} else {
	// 		fieldSwaggerSchema.Type = getSwaggerType(field.Type)
	// 		fieldSwaggerSchema.Format = getSwaggerFormat(field.Type)
	// 	}

	// 	def.Properties[strings.ToLower(field.Name)] = fieldSwaggerSchema
	// }

	def.Required = requiredProperties
	AddSwaggerDefinition(strings.ToLower(schema.Name), def)
}

func addNOSQLSwaggerCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, strings.Title(collection.Name))
	apiPath.GET.Summary = "Gets All " + strings.Title(collection.Name) + ".  Can be filtered by limit and skip."
	apiPath.GET.Produces = []string{"application/json"}

	limit := getSwaggerParameter("limit", "query", "Limit the number of records returned.", false, "integer")
	skip := getSwaggerParameter("skip", "query", "Skip an amount of records from the collection returned from the database query.", false, "integer")

	apiPath.GET.Parameters = []Swagger2Parameter{limit, skip}

	updateSwaggerOperationResponse(apiPath.GET, "array", "#/definitions/"+strings.ToLower(collection.Schema.Name))

	AddSwaggerPath(path, apiPath)
	AddSwaggerTag(strings.Title(collection.Name), "A collection of "+strings.Title(collection.Name), "", "")
}

func addNOSQLSwaggerSearchCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, "Search "+strings.Title(collection.Name))
	apiPath.GET.Summary = "Searches " + strings.Title(collection.Name) + " by field and value.  Can be filtered by limit and skip."
	apiPath.GET.Produces = []string{"application/json"}

	field := getSwaggerParameter("field", "query", "Field to search for "+collection.Schema.Name+" on.", true, "string")
	value := getSwaggerParameter("value", "query", "Value to search for "+collection.Schema.Name+" on.", true, "string")
	limit := getSwaggerParameter("limit", "query", "Limit the number of records returned.", false, "integer")
	skip := getSwaggerParameter("skip", "query", "Skip an amount of records from the collection returned from the database query.", false, "integer")

	apiPath.GET.Parameters = []Swagger2Parameter{field, value, limit, skip}
	AddSwaggerPath(path, apiPath)
	AddSwaggerTag("Search "+strings.Title(collection.Name), "A collection of "+strings.Title(collection.Name), "", "")
}

func genSearchCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getSearch" + name + "(c *gin.Context){\n"

	val += "\tfield := c.DefaultQuery(\"field\",\"\")\n"
	val += "\tvalue := c.DefaultQuery(\"value\",\"\")\n"
	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray := items.SearchAdvanced(field, value, limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray := items.Search(field, value)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genSortCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getSort" + name + "(c *gin.Context){\n"

	val += "\tfield := c.DefaultQuery(\"field\",\"\")\n"
	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray := items.AllByIndexAdvanced(field, limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray := items.AllByIndex(field)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genRangeCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getRange" + name + "(c *gin.Context){\n"

	val += "\tfield := c.DefaultQuery(\"field\",\"\")\n"
	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\tmin := c.DefaultQuery(\"min\",\"\")\n"
	val += "\tmax := c.DefaultQuery(\"max\",\"\")\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray := items.RangeAdvanced(min, max, field, limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray := items.Range(min, max, field)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genSingleCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getSingle" + name + "(c *gin.Context){\n"

	val += "\tfield := c.DefaultQuery(\"field\",\"\")\n"
	val += "\tvalue := c.DefaultQuery(\"value\",\"\")\n"
	val += "\titems := model." + name + "{}\n"
	val += "\titemsArray := items.Single(field, value)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}
