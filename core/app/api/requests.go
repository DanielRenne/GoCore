package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"

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
			c.JSON(http.StatusOK, e)
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

	ctl := getController(controller)
	method := ctl.MethodByName(action)

	if !method.IsValid() {
		e.Error.Message = "Method " + action + " not available to call."
		c.JSON(http.StatusOK, e)
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

	uriParamsData, err := base64.StdEncoding.DecodeString(uriParams)

	if err != nil {
		e.Error.Message = "Failed to decode uriParams:  " + err.Error()
		c.JSON(http.StatusOK, e)
		return
	}

	var x interface{}
	err = json.Unmarshal(uriParamsData, &x)

	if err != nil {
		e.Error.Message = "Failed to unmarshal uriParams:  " + err.Error()
		c.JSON(http.StatusOK, e)
		return
	}

	if paramType == genericType || paramType.String() == "interface {}" {
		if x != nil {
			in = append(in, reflect.ValueOf(x))
		}
	} else {

		raw := []byte(uriParamsData)

		if len(raw) > 0 {
			param := reflect.New(paramType)
			err := json.Unmarshal(raw, param.Interface())
			if err != nil {
				e.Error.Message = "Failed to unmarshal raw uriParamsData:  " + err.Error()
				c.JSON(http.StatusOK, e)
				return
			}

			in = append(in, param.Elem())
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

func processPOSTAPI(c *gin.Context) {

	controller := c.Query("controller")
	action := c.Query("action")
	uriParams := c.Query("uriParams")

	log.Printf("Controller:%s\nAction:%s\nURIParams:%s\n", controller, action, uriParams)

}
