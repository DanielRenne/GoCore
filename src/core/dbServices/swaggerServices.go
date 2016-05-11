package dbServices

import (
	"encoding/json"
	"fmt"
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
	Type             string          `json:"type"`
	Format           string          `json:"format,omitempty"`
	Items            *[]Swagger2Item `json:"items, omitempty"`
	CollectionFormat string          `json:"collectionFormat,omitempty"`
	Maximum          float64         `json:"maximum,omitempty"`
	ExclusiveMaximum bool            `json:"exclusiveMaximum"`
	Minimum          float64         `json:"minimum,omitempty"`
	ExclusiveMinimum bool            `json:"exclusiveMinimum"`
	MaxLength        int             `json:"maxLength,omitempty"`
	MinLength        int             `json:"minLength,omitempty"`
	Pattern          string          `json:"pattern,omitempty"`
	MaxItems         int             `json:"maxItems,omitempty"`
	MinItems         int             `json:"minItems,omitempty"`
	UniqueItems      bool            `json:"uniqueItems,omitempty"`
	MultipleOf       float64         `json:"multipleOf,omitempty"`
}

type Swagger2Header struct {
	Description      string          `json:"description"`
	Type             string          `json:"type"`
	Format           string          `json:"format,omitempty"`
	Items            *[]Swagger2Item `json:"items, omitempty"`
	CollectionFormat string          `json:"collectionFormat,omitempty"`
	Maximum          float64         `json:"maximum,omitempty"`
	ExclusiveMaximum bool            `json:"exclusiveMaximum"`
	Minimum          float64         `json:"minimum,omitempty"`
	ExclusiveMinimum bool            `json:"exclusiveMinimum"`
	MaxLength        int             `json:"maxLength,omitempty"`
	MinLength        int             `json:"minLength,omitempty"`
	Pattern          string          `json:"pattern,omitempty"`
	MaxItems         int             `json:"maxItems,omitempty"`
	MinItems         int             `json:"minItems,omitempty"`
	UniqueItems      bool            `json:"uniqueItems,omitempty"`
	MultipleOf       float64         `json:"multipleOf,omitempty"`
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
	Ref  string `json:"$ref,omitempty"`
	Type string `json:"type,omitempty"`
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
	Schema           *Swagger2Schema `json:"schema, omitempty"`
	Type             string          `json:"type,omitempty"`
	Format           string          `json:"format,omitempty"`
	AllowEmptyValue  bool            `json:allowEmptyValue,omitempty"`
	Items            *[]Swagger2Item `json:"items, omitempty"`
	CollectionFormat string          `json:"collectionFormat,omitempty"`
	Maximum          float64         `json:"maximum,omitempty"`
	ExclusiveMaximum bool            `json:"exclusiveMaximum"`
	Minimum          float64         `json:"minimum,omitempty"`
	ExclusiveMinimum bool            `json:"exclusiveMinimum"`
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
	Parameters  *[]Swagger2Parameter        `json:"parameters,omitempty"`
	Schemes     []string                    `json:"schemes,omitempty"`
	Deprecated  bool                        `json:"deprecated,omitempty"`
	Security    []map[string][]string       `json:"security,omitempty"`
}

type Swagger2Path struct {
	Ref        string               `json:"$ref,omitempty"`
	GET        *Swagger2Operation   `json:"get,omitempty"`
	PUT        *Swagger2Operation   `json:"put,omitempty"`
	POST       *Swagger2Operation   `json:"post,omitempty"`
	DELETE     *Swagger2Operation   `json:"delete,omitempty"`
	OPTIONS    *Swagger2Operation   `json:"options,omitempty"`
	HEAD       *Swagger2Operation   `json:"head,omitempty"`
	PATCH      *Swagger2Operation   `json:"patch,omitempty"`
	Parameters *[]Swagger2Parameter `json:"parameters,omitempty"`
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
	Definitions         map[string]interface{}            `json:"definitions,omitempty"`
	Parameters          map[string]Swagger2Parameter      `json:"parameters,omitempty"`
	Responses           map[string]Swagger2Response       `json:"responses,omitempty"`
	SecurityDefinitions map[string]Swagger2SecurityScheme `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string             `json:"security,omitempty"`
	ExternalDocs        *Swagger2ExternalDoc              `json:"externalDocs,omitempty"`
}

var SwaggerDefinition Swagger2

func init() {

	jsonData, err := ioutil.ReadFile(SWAGGER_SCHEMA_PATH + "/" + CURRENT_SWAGGER_VERSION + "/swagger.json")
	if err != nil {
		color.Red("Reading of swagger.json failed:  " + err.Error())
		return
	}

	errUnmarshal := json.Unmarshal(jsonData, &SwaggerDefinition)
	if errUnmarshal != nil {
		color.Red("Parsing / Unmarshaling of create.json failed:  " + errUnmarshal.Error())
		return
	}
	color.Green(SwaggerDefinition.Swagger)
}

func AddSwaggerPath(path string, swaggerPath Swagger2Path) error {
	SwaggerDefinition.Lock()
	SwaggerDefinition.Paths[path] = swaggerPath
	SwaggerDefinition.Unlock()
	fmt.Println(GetSwaggerDefinitionJSONString())
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

func getSwaggerGETPath() Swagger2Path {

	var apiPath Swagger2Path
	var op Swagger2Operation
	op.Responses = make(map[string]Swagger2Response)

	var response200 Swagger2Response
	response200.Description = "Successful operation"

	// response200Schema := Swagger2Schema{Type: "array"}
	// response200.Schema = &response200Schema

	op.Responses["200"] = response200
	apiPath.GET = &op

	return apiPath
}
