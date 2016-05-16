package dbServices

import (
	"core/extensions"
	"strings"
)

const basePath = "/api"

func genSchemaWebAPI(collection NOSQLCollection, schema NOSQLSchema, dbPackageName string, driver string, versionDir string) string {

	val := extensions.GenPackageImport("webAPI", []string{dbPackageName, "core/ginServer", "github.com/gin-gonic/gin", "core/extensions", "strings", "io/ioutil", "encoding/json"})
	val += "func init(){\n\n"
	//val += "\tginServer.AddRouterGroup(\"" + versionDir + "\", \"/" + extensions.MakeFirstLowerCase(schema.Name) + "\", \"GET\", get" + strings.Title(schema.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/single" + strings.Title(collection.Name) + "\", \"GET\", getSingle" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/search" + strings.Title(collection.Name) + "\", \"GET\", getSearch" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/sort" + strings.Title(collection.Name) + "\", \"GET\", getSort" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/range" + strings.Title(collection.Name) + "\", \"GET\", getRange" + strings.Title(collection.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/" + extensions.MakeFirstLowerCase(collection.Name) + "\", \"GET\", get" + strings.Title(collection.Name) + ")\n"

	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/" + extensions.MakeFirstLowerCase(schema.Name) + "\", \"POST\", post" + strings.Title(schema.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/" + extensions.MakeFirstLowerCase(schema.Name) + "\", \"PUT\", put" + strings.Title(schema.Name) + ")\n"
	val += "\tginServer.AddRouterGroup(\"" + basePath + "/" + versionDir + "\", \"/" + extensions.MakeFirstLowerCase(schema.Name) + "\", \"DELETE\", delete" + strings.Title(schema.Name) + ")\n"

	val += "}\n\n"

	val += genNOSQLSingleCollectionGET(collection)
	val += genNOSQLSearchCollectionGET(collection)
	val += genNOSQLSortCollectionGET(collection)
	val += genNOSQLRangeCollectionGET(collection)
	val += genNOSQLCollectionGET(collection)

	val += genNOSQLSchemaPost(schema)
	val += genNOSQLSchemaPut(schema)
	val += genNOSQLSchemaDelete(schema)

	//Add Swagger Paths
	addNOSQLSwaggerCollectionGet("/"+extensions.MakeFirstLowerCase(collection.Name), collection)
	addNOSQLSwaggerSearchCollectionGet("/search"+strings.Title(collection.Name), collection)
	addNOSQLSwaggerSingleCollectionGet("/single"+strings.Title(collection.Name), collection)
	addNOSQLSwaggerSortCollectionGet("/sort"+strings.Title(collection.Name), collection)
	addNOSQLSwaggerRangeCollectionGet("/range"+strings.Title(collection.Name), collection)

	addNOSQLSwaggerSchemaPostBody("/"+extensions.MakeFirstLowerCase(schema.Name), schema)
	addNOSQLSwaggerSchemaPutBody("/"+extensions.MakeFirstLowerCase(schema.Name), schema)
	addNOSQLSwaggerSchemaDeleteBody("/"+extensions.MakeFirstLowerCase(schema.Name), schema)

	addNOSQLSwaggerSchemaDefinition(schema)

	return val
}

func genNOSQLCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func get" + name + "(c *gin.Context){\n"

	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray, _ := items.AllAdvanced(limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray, _ := items.All()\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

//Adds to the Swagger.json root the Definitions
func addNOSQLSwaggerSchemaDefinition(schema NOSQLSchema) {
	var def Swagger2Schema
	requiredProperties := []string{}

	def.Type = "object"

	if len(schema.Fields) > 0 {
		def.Properties = make(map[string]Swagger2Schema)
	}

	for _, field := range schema.Fields {

		if field.Required == true {
			requiredProperties = append(requiredProperties, field.Name)
		}

		fieldSwaggerSchema := getNOSQLSwaggerSchemaFieldDefinition(field)

		def.Properties[extensions.MakeFirstLowerCase(field.Name)] = fieldSwaggerSchema
	}

	def.Required = requiredProperties
	AddSwaggerDefinition(extensions.MakeFirstLowerCase(schema.Name), def)
}

func getNOSQLSwaggerSchemaFieldDefinition(field NOSQLSchemaField) Swagger2Schema {

	fieldSwaggerSchema := Swagger2Schema{}
	if field.Type == "object" {
		fieldSwaggerSchema.Ref = "#/definitions/" + extensions.MakeFirstLowerCase(field.Schema.Name)
		addNOSQLSwaggerSchemaDefinition(field.Schema)
	} else {
		fieldSwaggerSchema.Type = getSwaggerType(field.Type)
		fieldSwaggerSchema.Format = getSwaggerFormat(field.Type)

		if fieldSwaggerSchema.Type == "array" {

			if field.Type == "objectArray" {
				item := Swagger2Item{
					Ref: "#/definitions/" + extensions.MakeFirstLowerCase(field.Schema.Name),
				}

				fieldSwaggerSchema.Items = &item

				addNOSQLSwaggerSchemaDefinition(field.Schema)
			} else {
				item := Swagger2Item{
					Type: getSwaggerArrayType(field.Type),
				}
				fieldSwaggerSchema.Items = &item
			}

		}
	}
	return fieldSwaggerSchema
}

func addNOSQLSwaggerSchemaPostBody(path string, schema NOSQLSchema) {

	apiPath := getSwaggerPOSTPath()

	apiPath.POST.Tags = append(apiPath.POST.Tags, strings.Title(schema.Name))
	apiPath.POST.Summary = "Save a new " + strings.Title(schema.Name) + ""
	apiPath.POST.Description = "POST a new " + strings.Title(schema.Name) + " as JSON format via http body request."
	apiPath.POST.Produces = []string{"application/json"}

	body := getSwaggerBodyParameter(strings.Title(schema.Name)+" that needs to be added.", true, "#/definitions/"+extensions.MakeFirstLowerCase(schema.Name))
	apiPath.POST.Parameters = []Swagger2Parameter{body}

	updateSwaggerOperationResponseRef(apiPath.POST, "200", "#/definitions/"+extensions.MakeFirstLowerCase(schema.Name), "successful operation")
	updateSwaggerOperationResponseRef(apiPath.POST, "406", "#/definitions/errorResponse", "Faild to parse JSON")
	updateSwaggerOperationResponseRef(apiPath.POST, "500", "#/definitions/errorResponse", "Faild to save object to database")

	AddSwaggerPOSTPath(path, apiPath)
	AddSwaggerTag(strings.Title(schema.Name), strings.Title(schema.Name)+" object", "", "")
}

func addNOSQLSwaggerSchemaPutBody(path string, schema NOSQLSchema) {

	apiPath := getSwaggerPUTPath()

	apiPath.PUT.Tags = append(apiPath.PUT.Tags, strings.Title(schema.Name))
	apiPath.PUT.Summary = "Save an existing " + strings.Title(schema.Name) + ""
	apiPath.PUT.Description = "PUT an existing " + strings.Title(schema.Name) + " as JSON format via http body request."
	apiPath.PUT.Produces = []string{"application/json"}

	body := getSwaggerBodyParameter(strings.Title(schema.Name)+" that needs to be updated.", true, "#/definitions/"+extensions.MakeFirstLowerCase(schema.Name))
	apiPath.PUT.Parameters = []Swagger2Parameter{body}

	updateSwaggerOperationResponseRef(apiPath.PUT, "200", "#/definitions/"+extensions.MakeFirstLowerCase(schema.Name), "successful operation")
	// updateSwaggerOperationResponseRef(apiPath.PUT, "404", "#/definitions/errorResponse", strings.Title(schema.Name)+" not found")
	updateSwaggerOperationResponseRef(apiPath.PUT, "406", "#/definitions/errorResponse", "Faild to parse JSON")
	updateSwaggerOperationResponseRef(apiPath.PUT, "500", "#/definitions/errorResponse", "Faild to save object to database")

	AddSwaggerPUTPath(path, apiPath)

}

func addNOSQLSwaggerSchemaDeleteBody(path string, schema NOSQLSchema) {

	apiPath := getSwaggerDELETEPath()

	apiPath.DELETE.Tags = append(apiPath.DELETE.Tags, strings.Title(schema.Name))
	apiPath.DELETE.Summary = "Remove an existing " + strings.Title(schema.Name) + ""
	apiPath.DELETE.Description = "DELETE an existing " + strings.Title(schema.Name) + " as JSON format via http body request."
	apiPath.DELETE.Produces = []string{"application/json"}

	body := getSwaggerBodyParameter(strings.Title(schema.Name)+" that needs to be updated.", true, "#/definitions/"+extensions.MakeFirstLowerCase(schema.Name))
	apiPath.DELETE.Parameters = []Swagger2Parameter{body}

	updateSwaggerOperationResponseRef(apiPath.DELETE, "200", "#/definitions/"+extensions.MakeFirstLowerCase(schema.Name), "successful operation")
	// updateSwaggerOperationResponseRef(apiPath.PUT, "404", "#/definitions/errorResponse", strings.Title(schema.Name)+" not found")
	updateSwaggerOperationResponseRef(apiPath.DELETE, "406", "#/definitions/errorResponse", "Faild to parse JSON")
	updateSwaggerOperationResponseRef(apiPath.DELETE, "500", "#/definitions/errorResponse", "Faild to delete object from database")

	AddSwaggerDELETEPath(path, apiPath)

}

func addNOSQLSwaggerCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, strings.Title(collection.Name))
	apiPath.GET.Summary = "Gets All " + strings.Title(collection.Name) + ""
	apiPath.GET.Description = "Can be filtered by limit and skip."
	apiPath.GET.Produces = []string{"application/json"}

	limit := getSwaggerParameter("limit", "query", "Limit the number of records returned.", false, "integer")
	skip := getSwaggerParameter("skip", "query", "Skip an amount of records from the collection returned from the database query.", false, "integer")

	apiPath.GET.Parameters = []Swagger2Parameter{limit, skip}

	updateSwaggerOperationResponse(apiPath.GET, "array", "#/definitions/"+extensions.MakeFirstLowerCase(collection.Schema.Name))

	AddSwaggerGETPath(path, apiPath)
	AddSwaggerTag(strings.Title(collection.Name), "All "+strings.Title(collection.Name), "", "")
}

func addNOSQLSwaggerSearchCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, "Search "+strings.Title(collection.Name))
	apiPath.GET.Summary = "Searches " + strings.Title(collection.Name)
	apiPath.GET.Description = "Is searched by the field and value.  Can be filtered by limit and skip."
	apiPath.GET.Produces = []string{"application/json"}

	field := getSwaggerParameter("field", "query", "Field to search for "+collection.Schema.Name+" on.", true, "string")
	value := getSwaggerParameter("value", "query", "Value to search for "+collection.Schema.Name+" on.", true, "string")
	limit := getSwaggerParameter("limit", "query", "Limit the number of records returned.", false, "integer")
	skip := getSwaggerParameter("skip", "query", "Skip an amount of records from the collection returned from the database query.", false, "integer")

	apiPath.GET.Parameters = []Swagger2Parameter{field, value, limit, skip}

	updateSwaggerOperationResponse(apiPath.GET, "array", "#/definitions/"+extensions.MakeFirstLowerCase(collection.Schema.Name))

	AddSwaggerGETPath(path, apiPath)
	AddSwaggerTag("Search "+strings.Title(collection.Name), "Search for "+strings.Title(collection.Name), "", "")
}

func addNOSQLSwaggerSingleCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, "Single "+strings.Title(collection.Schema.Name))
	apiPath.GET.Summary = "Get a " + strings.Title(collection.Schema.Name)
	apiPath.GET.Description = "Returns a single " + collection.Schema.Name + " searched by field and value."
	apiPath.GET.Produces = []string{"application/json"}

	field := getSwaggerParameter("field", "query", "Field to search for "+collection.Schema.Name+" on.", true, "string")
	value := getSwaggerParameter("value", "query", "Value to search for "+collection.Schema.Name+" on.", true, "string")

	apiPath.GET.Parameters = []Swagger2Parameter{field, value}

	updateSwaggerOperationResponseRef(apiPath.GET, "200", "#/definitions/"+extensions.MakeFirstLowerCase(collection.Schema.Name), "successful operation")

	AddSwaggerGETPath(path, apiPath)
	AddSwaggerTag("Single "+strings.Title(collection.Schema.Name), "A Single "+strings.Title(collection.Schema.Name), "", "")
}

func addNOSQLSwaggerSortCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, "Sort "+strings.Title(collection.Name))
	apiPath.GET.Summary = "Sort " + strings.Title(collection.Name)
	apiPath.GET.Description = "Sort Collection by Field.  Can be filtered by limit and skip."
	apiPath.GET.Produces = []string{"application/json"}

	field := getSwaggerParameter("field", "query", "Field to sort "+collection.Name+" on.", true, "string")
	limit := getSwaggerParameter("limit", "query", "Limit the number of records returned.", false, "integer")
	skip := getSwaggerParameter("skip", "query", "Skip an amount of records from the collection returned from the database query.", false, "integer")

	apiPath.GET.Parameters = []Swagger2Parameter{field, limit, skip}

	updateSwaggerOperationResponse(apiPath.GET, "array", "#/definitions/"+extensions.MakeFirstLowerCase(collection.Schema.Name))

	AddSwaggerGETPath(path, apiPath)
	AddSwaggerTag("Sort "+strings.Title(collection.Name), "Sorted "+strings.Title(collection.Name), "", "")
}

func addNOSQLSwaggerRangeCollectionGet(path string, collection NOSQLCollection) {

	apiPath := getSwaggerGETPath()

	apiPath.GET.Tags = append(apiPath.GET.Tags, "Range "+strings.Title(collection.Name))
	apiPath.GET.Summary = "Range " + strings.Title(collection.Name)
	apiPath.GET.Description = "Is searched by the field and value.  Can be filtered by limit and skip."
	apiPath.GET.Produces = []string{"application/json"}

	field := getSwaggerParameter("field", "query", "Field to range records for "+collection.Schema.Name+" on.", true, "string")
	min := getSwaggerParameter("min", "query", "Minimum value to range on "+collection.Schema.Name+" on.", true, "string")
	max := getSwaggerParameter("max", "query", "Maximum value to range on "+collection.Schema.Name+" on.", true, "string")
	limit := getSwaggerParameter("limit", "query", "Limit the number of records returned.", false, "integer")
	skip := getSwaggerParameter("skip", "query", "Skip an amount of records from the collection returned from the database query.", false, "integer")

	apiPath.GET.Parameters = []Swagger2Parameter{field, min, max, limit, skip}

	updateSwaggerOperationResponse(apiPath.GET, "array", "#/definitions/"+extensions.MakeFirstLowerCase(collection.Schema.Name))

	AddSwaggerGETPath(path, apiPath)
	AddSwaggerTag("Range "+strings.Title(collection.Name), "Range of "+strings.Title(collection.Name), "", "")
}

func genNOSQLSearchCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getSearch" + name + "(c *gin.Context){\n"

	val += "\tfield := strings.Title(c.DefaultQuery(\"field\",\"\"))\n"
	val += "\tvalue := c.DefaultQuery(\"value\",\"\")\n"
	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray, _ := items.SearchAdvanced(field, value, limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray, _  := items.Search(field, value)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genNOSQLSortCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getSort" + name + "(c *gin.Context){\n"

	val += "\tfield := strings.Title(c.DefaultQuery(\"field\",\"\"))\n"
	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray, _ := items.AllByIndexAdvanced(field, limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray, _ := items.AllByIndex(field)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genNOSQLRangeCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getRange" + name + "(c *gin.Context){\n"

	val += "\tfield := strings.Title(c.DefaultQuery(\"field\",\"\"))\n"
	val += "\tlimit := extensions.StringToInt(c.DefaultQuery(\"limit\",\"\"))\n"
	val += "\tskip := extensions.StringToInt(c.DefaultQuery(\"skip\",\"\"))\n"
	val += "\tmin := c.DefaultQuery(\"min\",\"\")\n"
	val += "\tmax := c.DefaultQuery(\"max\",\"\")\n"
	val += "\titems := model." + name + "{}\n"
	val += "\tif limit != 0 || skip != 0{\n"
	val += "\titemsArray, _ := items.RangeAdvanced(min, max, field, limit, skip)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"
	val += "\t\treturn"
	val += "\t}\n"
	val += "\titemsArray, _ := items.Range(min, max, field)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genNOSQLSingleCollectionGET(collection NOSQLCollection) string {

	name := strings.Title(collection.Name)

	val := ""
	val += "func getSingle" + name + "(c *gin.Context){\n"

	val += "\tfield := strings.Title(c.DefaultQuery(\"field\",\"\"))\n"
	val += "\tvalue := c.DefaultQuery(\"value\",\"\")\n"
	val += "\titems := model." + name + "{}\n"
	val += "\titemsArray, _ := items.Single(field, value)\n"
	val += "\tginServer.RespondJSON(itemsArray, c)\n"

	val += "}\n\n"

	return val
}

func genNOSQLSchemaPost(schema NOSQLSchema) string {

	name := strings.Title(schema.Name)

	val := ""
	val += "func post" + name + "(c *gin.Context){\n"
	val += "\tbody := c.Request.Body\n"
	val += "\tx, _ := ioutil.ReadAll(body)\n"

	val += "\tvar obj model." + name + "\n"
	val += "\terrMarshal := json.Unmarshal(x, &obj)\n"

	val += "\tif errMarshal != nil{\n"
	val += "\t\tc.Data(406, gin.MIMEHTML, ginServer.RespondError(errMarshal.Error()))\n"
	val += "\treturn\n"
	val += "\t}\n"
	val += "\terrSave := obj.Save()\n"
	val += "\tif errSave != nil{\n"
	val += "\t\tc.Data(500, gin.MIMEHTML, ginServer.RespondError(errSave.Error()))\n"
	val += "\treturn\n"
	val += "\t}\n"
	val += "\tginServer.RespondJSON(obj, c)\n"

	val += "}\n\n"

	return val
}

func genNOSQLSchemaPut(schema NOSQLSchema) string {

	name := strings.Title(schema.Name)

	val := ""
	val += "func put" + name + "(c *gin.Context){\n"
	val += "\tbody := c.Request.Body\n"
	val += "\tx, _ := ioutil.ReadAll(body)\n"

	val += "\tvar obj model." + name + "\n"
	val += "\terrMarshal := json.Unmarshal(x, &obj)\n"

	val += "\tif errMarshal != nil{\n"
	val += "\t\tc.Data(406, gin.MIMEHTML, ginServer.RespondError(errMarshal.Error()))\n"
	val += "\treturn\n"
	val += "\t}\n"

	val += "\terrSave := obj.Save()\n"
	val += "\tif errSave != nil{\n"
	val += "\t\tc.Data(500, gin.MIMEHTML, ginServer.RespondError(errSave.Error()))\n"
	val += "\treturn\n"
	val += "\t}\n"
	val += "\tginServer.RespondJSON(obj, c)\n"

	val += "}\n\n"

	return val
}

func genNOSQLSchemaDelete(schema NOSQLSchema) string {

	name := strings.Title(schema.Name)

	val := ""
	val += "func delete" + name + "(c *gin.Context){\n"
	val += "\tbody := c.Request.Body\n"
	val += "\tx, _ := ioutil.ReadAll(body)\n"

	val += "\tvar obj model." + name + "\n"
	val += "\terrMarshal := json.Unmarshal(x, &obj)\n"

	val += "\tif errMarshal != nil{\n"
	val += "\t\tc.Data(406, gin.MIMEHTML, ginServer.RespondError(errMarshal.Error()))\n"
	val += "\treturn\n"
	val += "\t}\n"

	val += "\terrDelete := obj.Delete()\n"
	val += "\tif errDelete != nil{\n"
	val += "\t\tc.Data(500, gin.MIMEHTML, ginServer.RespondError(errDelete.Error()))\n"
	val += "\treturn\n"
	val += "\t}\n"
	val += "\tginServer.RespondJSON(obj, c)\n"

	val += "}\n\n"

	return val
}
