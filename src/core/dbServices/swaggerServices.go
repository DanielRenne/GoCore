package dbServices

import (
	"core/serverSettings"
	"encoding/json"
	// "fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"sync"
)

const SWAGGER_SCHEMA_PATH = "swagger/schemas"
const CURRENT_SWAGGER_VERSION = "2.0"

type Swagger2License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Swagger2Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type Swagger2Item struct {
	Ref              string        `json:"$ref,omitempty"`
	Type             string        `json:"type,omitempty"`
	Format           string        `json:"format,omitempty"`
	Items            *Swagger2Item `json:"items,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
	Maximum          float64       `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          float64       `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        int           `json:"maxLength,omitempty"`
	MinLength        int           `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         int           `json:"maxItems,omitempty"`
	MinItems         int           `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       float64       `json:"multipleOf,omitempty"`
}

type Swagger2Header struct {
	Description      string        `json:"description"`
	Type             string        `json:"type"`
	Format           string        `json:"format,omitempty"`
	Items            *Swagger2Item `json:"items,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
	Maximum          float64       `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum"`
	Minimum          float64       `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum"`
	MaxLength        int           `json:"maxLength,omitempty"`
	MinLength        int           `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         int           `json:"maxItems,omitempty"`
	MinItems         int           `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       float64       `json:"multipleOf,omitempty"`
}

type Swagger2SecurityScheme struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name"`
	In               string            `json:"in"`
	Flow             string            `json:"flow"`
	AuthorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	Scopes           map[string]string `json:"scopes"`
}

type Swagger2Schema struct {
	Ref                  string                    `json:"$ref,omitempty"`
	Description          string                    `json:"description,omitempty"`
	Type                 string                    `json:"type,omitempty"`
	Format               string                    `json:"format,omitempty"`
	Items                *Swagger2Item             `json:"items,omitempty"`
	Default              string                    `json:"default,omitempty"`
	MultipleOf           float64                   `json:"multipleOf,omitempty"`
	Maximum              float64                   `json:"maximum,omitempty"`
	ExclusiveMaximum     bool                      `json:"exclusiveMaximum,omitempty"`
	Minimum              float64                   `json:"minimum,omitempty"`
	ExclusiveMinimum     bool                      `json:"exclusiveMinimum,omitempty"`
	MaxLength            int                       `json:"maxLength,omitempty"`
	MinLength            int                       `json:"minLength,omitempty"`
	Pattern              string                    `json:"pattern,omitempty"`
	MaxItems             int                       `json:"maxItems,omitempty"`
	MinItems             int                       `json:"minItems,omitempty"`
	UniqueItems          bool                      `json:"uniqueItems,omitempty"`
	MaxProperties        int                       `json:"maxProperties,omitempty"`
	MinProperties        int                       `json:"minProperties,omitempty"`
	Required             []string                  `json:"required,omitempty"`
	AdditionalProperties []string                  `json:"additionalProperties,omitempty"`
	Properties           map[string]Swagger2Schema `json:"properties,omitempty"`
}

type Swagger2Response struct {
	Description string                    `json:"description"`
	Schema      *Swagger2Schema           `json:"schema,omitempty"`
	Headers     map[string]Swagger2Header `json:"headers,omitempty"`
	Examples    string                    `json:"examples,omitempty"`
}

type Swagger2Parameter struct {
	In               string          `json:"in"`
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	Required         bool            `json:"required,omitempty"`
	Schema           *Swagger2Schema `json:"schema,omitempty"`
	Type             string          `json:"type,omitempty"`
	Format           string          `json:"format,omitempty"`
	AllowEmptyValue  bool            `json:"allowEmptyValue,omitempty"`
	Items            *Swagger2Item   `json:"items,omitempty"`
	CollectionFormat string          `json:"collectionFormat,omitempty"`
	Maximum          float64         `json:"maximum,omitempty"`
	ExclusiveMaximum bool            `json:"exclusiveMaximum,omitempty"`
	Minimum          float64         `json:"minimum,omitempty"`
	ExclusiveMinimum bool            `json:"exclusiveMinimum,omitempty"`
	MaxLength        int             `json:"maxLength,omitempty"`
	MinLength        int             `json:"minLength,omitempty"`
	Pattern          string          `json:"pattern,omitempty"`
	MaxItems         int             `json:"maxItems,omitempty"`
	MinItems         int             `json:"minItems,omitempty"`
	UniqueItems      bool            `json:"uniqueItems,omitempty"`
	MultipleOf       float64         `json:"multipleOf,omitempty"`
}

type Swagger2Operation struct {
	Responses   map[string]Swagger2Response `json:"responses"`
	Tags        []string                    `json:"tags,omitempty"`
	Summary     string                      `json:"summary,omitempty"`
	Description string                      `json:"description,omitempty"`
	OperationId string                      `json:"operationId,omitempty"`
	Consumes    []string                    `json:"consumes,omitempty"`
	Produces    []string                    `json:"produces,omitempty"`
	Parameters  []Swagger2Parameter         `json:"parameters,omitempty"`
	Schemes     []string                    `json:"schemes,omitempty"`
	Deprecated  bool                        `json:"deprecated,omitempty"`
	Security    []map[string][]string       `json:"security,omitempty"`
}

type Swagger2Path struct {
	Ref        string              `json:"$ref,omitempty"`
	GET        *Swagger2Operation  `json:"get,omitempty"`
	PUT        *Swagger2Operation  `json:"put,omitempty"`
	POST       *Swagger2Operation  `json:"post,omitempty"`
	DELETE     *Swagger2Operation  `json:"delete,omitempty"`
	OPTIONS    *Swagger2Operation  `json:"options,omitempty"`
	HEAD       *Swagger2Operation  `json:"head,omitempty"`
	PATCH      *Swagger2Operation  `json:"patch,omitempty"`
	Parameters []Swagger2Parameter `json:"parameters,omitempty"`
}

type Swagger2ExternalDoc struct {
	Description string `json:"description, omitempty"`
	URL         string `json:"url"`
}

type Swagger2Tag struct {
	Name         string               `json:"name"`
	Description  string               `json:"description,omitempty"`
	ExternalDocs *Swagger2ExternalDoc `json:"externalDocs,omitempty"`
}

type Swagger2Info struct {
	Title          string           `json:"title"`
	Description    string           `json:"description,omitempty"`
	TermsOfService string           `json:"termsOfService,omitempty"`
	Contact        *Swagger2Contact `json:"contact,omitempty"`
	License        *Swagger2License `json:"license,omitempty"`
	Version        string           `json:"version,omitempty"`
}

type Swagger2 struct {
	sync.RWMutex
	Swagger             string                            `json:"swagger"`
	Info                *Swagger2Info                     `json:"info"`
	Host                string                            `json:"host,omitempty"`
	BasePath            string                            `json:"basePath,omitempty"`
	Tags                []Swagger2Tag                     `json:"tags,omitempty"`
	Schemes             []string                          `json:"schemes,omitempty"`
	Paths               map[string]Swagger2Path           `json:"paths,omitempty"`
	Definitions         map[string]Swagger2Schema         `json:"definitions,omitempty"`
	Parameters          map[string]Swagger2Parameter      `json:"parameters,omitempty"`
	Responses           map[string]Swagger2Response       `json:"responses,omitempty"`
	SecurityDefinitions map[string]Swagger2SecurityScheme `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string             `json:"security,omitempty"`
	ExternalDocs        *Swagger2ExternalDoc              `json:"externalDocs,omitempty"`
}

var SwaggerDefinition Swagger2

func init() {
	LoadSwaggerTemplate()
}

func LoadSwaggerTemplate() {
	jsonData, err := ioutil.ReadFile(SWAGGER_SCHEMA_PATH + "/" + CURRENT_SWAGGER_VERSION + "/swagger.json")
	if err != nil {
		color.Red("Reading of swagger.json failed:  " + err.Error())
		return
	}

	var swagDef Swagger2
	errUnmarshal := json.Unmarshal(jsonData, &swagDef)
	if errUnmarshal != nil {
		color.Red("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
		return
	}
	SwaggerDefinition = swagDef
	color.Green(SwaggerDefinition.Swagger)
}

func AddSwaggerPath(path string, swaggerPath Swagger2Path) error {
	SwaggerDefinition.Lock()
	SwaggerDefinition.Paths[path] = swaggerPath
	SwaggerDefinition.Unlock()
	return nil
}

func AddSwaggerDefinition(name string, swaggerSchema Swagger2Schema) error {
	SwaggerDefinition.Lock()
	SwaggerDefinition.Definitions[name] = swaggerSchema
	SwaggerDefinition.Unlock()
	return nil
}

func AddSwaggerTag(name string, description string, docDescription string, docURL string) error {

	var st Swagger2Tag
	st.Name = name
	st.Description = description

	if docDescription != "" {
		ext := Swagger2ExternalDoc{Description: docDescription, URL: docURL}
		st.ExternalDocs = &ext
	}

	SwaggerDefinition.Tags = append(SwaggerDefinition.Tags, st)
	return nil
}

func GetSwaggerDefinitionJSONString() string {
	// bytes, err := json.MarshalIndent(SwaggerDefinition, "", "    ")
	bytes, err := json.Marshal(SwaggerDefinition)
	if err != nil {
		color.Red("Marshaling of SwaggerDefinition failed:  " + err.Error())
		return ""
	}
	return string(bytes)
}

func writeSwaggerConfiguration(verisonPath string, version string) {

	//Save Application Meta Data and Version information to Swagger
	contact := Swagger2Contact{
		Email: serverSettings.WebConfig.Application.Info.Contact.Email,
		Name:  serverSettings.WebConfig.Application.Info.Contact.Name,
		URL:   serverSettings.WebConfig.Application.Info.Contact.URL,
	}

	license := Swagger2License{
		Name: serverSettings.WebConfig.Application.Info.License.Name,
		URL:  serverSettings.WebConfig.Application.Info.License.URL,
	}

	info := Swagger2Info{
		Title:       serverSettings.WebConfig.Application.Info.Title,
		Description: serverSettings.WebConfig.Application.Info.Description,
		Contact:     &contact,
		License:     &license,
		Version:     version,
	}

	SwaggerDefinition.BasePath = verisonPath
	SwaggerDefinition.Host = serverSettings.WebConfig.Application.Domain
	SwaggerDefinition.Info = &info

	//Write out the swagger api Definition to the Application
	err := ioutil.WriteFile(serverSettings.SwaggerUIPath+"/swagger."+version+".json", []byte(GetSwaggerDefinitionJSONString()), 0777)
	if err != nil {
		color.Red("Error writing swagger.json:  " + err.Error())
		return
	}
	color.Green("Successfully created " + serverSettings.SwaggerUIPath + "/swagger." + version + ".json")

	LoadSwaggerTemplate()
}

func getSwaggerGETPath() Swagger2Path {

	var apiPath Swagger2Path
	var op Swagger2Operation
	op.Responses = make(map[string]Swagger2Response)

	var response200 Swagger2Response
	response200.Description = "Successful operation"

	op.Responses["200"] = response200
	apiPath.GET = &op

	return apiPath
}

func updateSwaggerOperationResponse(op *Swagger2Operation, responseType string, itemReference string) {

	response200 := op.Responses["200"]

	item := Swagger2Item{
		Ref: itemReference,
	}

	s := Swagger2Schema{
		Type:  responseType,
		Items: &item,
	}

	response200.Schema = &s
	op.Responses["200"] = response200
}

func updateSwaggerOperationResponseRef(op *Swagger2Operation, itemReference string) {

	response200 := op.Responses["200"]

	s := Swagger2Schema{
		Ref: itemReference,
	}

	response200.Schema = &s
	op.Responses["200"] = response200
}

func getSwaggerParameter(name string, in string, description string, required bool, valType string) (param Swagger2Parameter) {
	param.Name = name
	param.In = in
	param.Description = description
	param.Required = required
	param.Type = valType
	return
}

func getSwaggerType(value string) string {

	switch value {
	case "int":
		return "integer"
	case "float64":
		return "number"
	case "bool":
		return "boolean"
	case "byteArray":
		return "string"
	case "intArray":
		return "array"
	case "float64Array":
		return "array"
	case "stringArray":
		return "array"
	case "boolArray":
		return "array"
	case "objectArray":
		return "array"
	}
	return value
}

func getSwaggerArrayType(value string) string {

	switch value {
	case "intArray":
		return "integer"
	case "float64Array":
		return "number"
	case "stringArray":
		return "string"
	case "boolArray":
		return "boolean"
	}
	return value
}

func getSwaggerFormat(value string) string {

	switch value {
	case "int":
		return "int32"
	case "float64":
		return "float"
	case "bool":
		return ""
	case "string":
		return ""
	case "byteArray":
		return "binary"
	case "intArray":
		return ""
	case "float64Array":
		return ""
	case "stringArray":
		return ""
	case "boolArray":
		return ""
	case "objectArray":
		return ""
	}
	return value
}
