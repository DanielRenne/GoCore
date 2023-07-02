package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"

	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/gin-gonic/gin"
)

// ErrorResponse is the default error response object.
type ErrorResponse struct {
	Error *errorObj `json:"error"`
}

type errorObj struct {
	Message    string `json:"Message"`
	Code       string `json:"code"`
	Stacktrace string `json:"stackTrace"`
}

type emptyResponse struct{}

type socketAPIRequest struct {
	CallbackID int        `json:"callBackId"`
	Data       apiRequest `json:"data"`
}

type apiRequest struct {
	Action     string      `json:"action"`
	State      interface{} `json:"state"`
	Controller string      `json:"controller"`
}

type socketAPIResponse struct {
	CallbackId int         `json:"callBackId"`
	Data       interface{} `json:"data"`
}

func processGETAPI(c *gin.Context) {

	controller := c.Query("controller")
	action := c.Query("action")
	uriParams := c.Query("uriParams")

	log.Println("process GET API " + controller + " " + action)

	defer func() {
		if r := recover(); r != nil {

			log.Println("Failed to processGETAPI:  Controller:  " + controller + " Action: " + action + " " + fmt.Sprintf("%+v", r))
			log.Println("Panic Stack: " + string(debug.Stack()))

			var e ErrorResponse
			e.Error.Message = "Recover Error:  " + fmt.Sprintf("%+v", r)
			c.JSON(http.StatusInternalServerError, e)
			return
		}
	}()

	var e ErrorResponse
	e.Error = new(errorObj)

	if action == "" {
		action = "Root"
	}

	var uriParamsData []byte
	var err error

	if uriParams != "" {
		uriParamsData, err = base64.StdEncoding.DecodeString(uriParams)
	}

	if err != nil {
		e.Error.Message = "Failed to decode uriParams:  " + err.Error()
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	if serverSettings.WebConfig.Application.AllowCrossOriginRequests {
		c.Header("Access-Control-Allow-Origin", "*")
	}

	response := func(y interface{}, e ErrorResponse, httpStatus int) {
		processHTTPResponse(y, e, httpStatus, c)
	}

	processRequest(controller, action, uriParamsData, c, response)

}

func processPOSTAPI(c *gin.Context) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Stack: " + string(debug.Stack()))
			log.Println("Recover Error:  " + fmt.Sprintf("%+v", r))
			var e ErrorResponse
			e.Error.Message = "Recover Error:  " + fmt.Sprintf("%+v", r)
			c.JSON(http.StatusInternalServerError, e)
			return
		}
	}()

	controller := c.Query("controller")
	action := c.Query("action")

	body, _ := ginServer.GetRequestBody(c)

	// if serverSettings.WebConfig.Application.AllowCrossOriginRequests {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// }

	response := func(y interface{}, e ErrorResponse, httpStatus int) {
		processHTTPResponse(y, e, httpStatus, c)
	}

	processRequest(controller, action, body, c, response)
}

func processHTTPResponse(y interface{}, e ErrorResponse, httpStatus int, c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	if y == nil {
		if c.Writer.Status() == http.StatusFound || c.Writer.Status() == http.StatusMovedPermanently {
			return
		}
		c.JSON(httpStatus, e)
	} else {
		data, ok := y.([]byte)
		if ok {
			c.Writer.Header().Set("Content-Type", "application/octet-stream")
			c.Writer.Header().Set("Content-Length", strconv.Itoa((len(data))))
			c.Writer.Write(data)
		} else {
			c.JSON(httpStatus, y)
		}

	}
}

func processSocketAPI(c *gin.Context, data []byte, conn *app.WebSocketConnection) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Stack at requests.processSocketAPI: " + string(debug.Stack()))
			log.Println("Recover Error:  " + fmt.Sprintf("%+v", r))
			return
		}
	}()

	var request socketAPIRequest

	var e ErrorResponse
	e.Error = new(errorObj)

	var socketResponse socketAPIResponse

	errMarshal := json.Unmarshal(data, &request)
	if errMarshal != nil {

		e.Error.Message = "Failed to unmarshal socketAPIRequest:  " + errMarshal.Error()
		socketResponse.Data = e
		app.ReplyToWebSocketJSON(conn, socketResponse)
		return
	}

	socketResponse.CallbackId = request.CallbackID

	response := func(y interface{}, e ErrorResponse, httpStatus int) {
		if y == nil {
			socketResponse.Data = e
			app.ReplyToWebSocketJSON(conn, socketResponse)
		} else {
			socketResponse.Data = y
			app.ReplyToWebSocketJSON(conn, socketResponse)
		}
	}

	data, err := json.Marshal(request.Data.State)
	if err != nil {
		e.Error.Message = "Failed to Marshal socketAPIRequest.Data.State:  " + err.Error()
		socketResponse.Data = e
		app.ReplyToWebSocketJSON(conn, socketResponse)
		return
	}

	processRequest(request.Data.Controller, request.Data.Action, data, c, response)

}

// ProcessRequest will process a controller requeest.
func ProcessRequest(controller string, action string, data []byte, results func(y interface{}, e ErrorResponse, httpStatus int)) {
	processRequest(controller, action, data, nil, results)
}

func processRequest(controller string, action string, data []byte, c *gin.Context, results func(y interface{}, e ErrorResponse, httpStatus int)) {

	var e ErrorResponse
	e.Error = new(errorObj)

	ctl := getController(extensions.Title(controller))

	method := ctl.MethodByName(extensions.Title(action))

	if !method.IsValid() {
		e.Error.Message = "Method " + action + " not available to call."
		// c.JSON(http.StatusNotImplemented, e)

		results(nil, e, http.StatusNotImplemented)
		return
	}

	methodType := method.Type()
	paramCnt := methodType.NumIn()
	in := []reflect.Value{}

	if paramCnt == 0 {

		value := method.Call(in)
		if len(value) > 0 {
			y := value[0].Interface()
			results(y, e, http.StatusOK)
		} else {
			results(emptyResponse{}, e, http.StatusOK)
		}
		return
	}

	paramType := methodType.In(0)
	genericType := reflect.TypeOf((*interface{})(nil))

	if paramType == genericType || paramType.String() == "interface {}" {

		var x interface{}
		err := json.Unmarshal(data, &x)

		if err != nil {
			e.Error.Message = "Failed to unmarshal uriParams or post body data:  " + err.Error()
			results(nil, e, http.StatusInternalServerError)
			return
		}

		if x != nil {
			in = append(in, reflect.ValueOf(x))
		}
	} else {

		raw := data

		if len(raw) > 0 {
			param := reflect.New(paramType)

			err := json.Unmarshal(raw, param.Interface())
			if err != nil {
				e.Error.Message = "Failed to unmarshal raw uriParamsData:  " + err.Error()
				results(nil, e, http.StatusInternalServerError)
				return
			}

			in = append(in, param.Elem())
		} else {
			in = append(in, reflect.ValueOf(c))
		}
	}

	if paramCnt == 2 {
		in = append(in, reflect.ValueOf(c))
	}

	value := method.Call(in)
	if len(value) > 0 {
		y := value[0].Interface()
		results(y, e, http.StatusOK)
	} else {
		results(emptyResponse{}, e, http.StatusOK)
	}
}
