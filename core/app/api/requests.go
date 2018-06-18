package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error *errorObj `json:"error,omitEmpty"`
}

type errorObj struct {
	Message    string `json:"Message"`
	Code       string `json:"code"`
	Stacktrace string `json:"stackTrace"`
}

type emptyResponse struct{}

func processGETAPI(c *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Stack: " + string(debug.Stack()))
			log.Println("Recover Error:  " + fmt.Sprintf("%+v", r))
			var e errorResponse
			e.Error.Message = "Recover Error:  " + fmt.Sprintf("%+v", r)
			c.JSON(http.StatusInternalServerError, e)
			return
		}
	}()

	var e errorResponse
	e.Error = new(errorObj)

	controller := c.Query("controller")
	action := c.Query("action")
	uriParams := c.Query("uriParams")

	if action == "" {
		action = "Root"
	}

	uriParamsData, err := base64.StdEncoding.DecodeString(uriParams)

	if err != nil {
		e.Error.Message = "Failed to decode uriParams:  " + err.Error()
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	processRequest(controller, action, uriParamsData, c)

}

func processPOSTAPI(c *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Stack: " + string(debug.Stack()))
			log.Println("Recover Error:  " + fmt.Sprintf("%+v", r))
			var e errorResponse
			e.Error.Message = "Recover Error:  " + fmt.Sprintf("%+v", r)
			c.JSON(http.StatusInternalServerError, e)
			return
		}
	}()

	controller := c.Query("controller")
	action := c.Query("action")

	body, _ := ginServer.GetRequestBody(c)

	processRequest(controller, action, body, c)
}

func processRequest(controller string, action string, data []byte, c *gin.Context) {

	var e errorResponse
	e.Error = new(errorObj)

	ctl := getController(controller)
	method := ctl.MethodByName(action)

	if !method.IsValid() {
		e.Error.Message = "Method " + action + " not available to call."
		c.JSON(http.StatusNotImplemented, e)
		return
	}

	methodType := method.Type()
	paramCnt := methodType.NumIn()
	paramType := methodType.In(0)
	genericType := reflect.TypeOf((*interface{})(nil))
	in := []reflect.Value{}

	if paramCnt == 0 {

		var tmp string
		in = append(in, reflect.ValueOf(tmp))

		value := method.Call(in)
		if len(value) > 0 {
			y := value[0].Interface()
			c.JSON(http.StatusOK, y)
		} else {
			c.JSON(http.StatusOK, emptyResponse{})
		}
		return
	}

	if len(data) == 0 {
		e.Error.Message = "No data posted.  Method expects a parameter of data."
		c.JSON(http.StatusBadRequest, e)
		return
	}

	if paramType == genericType || paramType.String() == "interface {}" {

		var x interface{}
		err := json.Unmarshal(data, &x)

		if err != nil {
			e.Error.Message = "Failed to unmarshal uriParams or post body data:  " + err.Error()
			c.JSON(http.StatusInternalServerError, e)
			return
		}

		if x != nil {
			in = append(in, reflect.ValueOf(x))
		}
	} else {

		raw := data

		if len(raw) > 0 {
			param := reflect.New(paramType)

			if paramType.String() == "string" {
				paramValue := string(data)
				in = append(in, reflect.ValueOf(paramValue))
			} else if paramType.String() == "int" {
				paramValue := extensions.IntToString(string(data))
				in = append(in, reflect.ValueOf(paramValue))
			} else {
				err := json.Unmarshal(raw, param.Interface())
				if err != nil {
					e.Error.Message = "Failed to unmarshal raw uriParamsData:  " + err.Error()
					log.Printf("%+v", e)
					c.JSON(http.StatusInternalServerError, e)
					return
				}

				in = append(in, param.Elem())
			}

		} else {
			var tmp string
			in = append(in, reflect.ValueOf(tmp))
		}
	}

	value := method.Call(in)
	if len(value) > 0 {
		y := value[0].Interface()
		c.JSON(http.StatusOK, y)
	} else {
		c.JSON(http.StatusOK, emptyResponse{})
	}
}
