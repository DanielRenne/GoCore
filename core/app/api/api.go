//api provides an API router and controller Registry to route controller api's.
package api

import (
	"log"
	"net/http"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

//private local variables
var registry sync.Map

//public variables

/*APICallback provides the routing to controller methods.
Implementation example-----------
ginServer.Router.GET("/apiGET", appAPI.APICallback)
ginServer.Router.POST("/api", appAPI.APICallback)
---------------------------------
*/
func APICallback(c *gin.Context) {

	log.Println(c.Request.Method)

	switch c.Request.Method {
	case http.MethodGet:
		processGETAPI(c)
	case http.MethodPost:
		processPOSTAPI(c)
	}
}

//RegisterController registers a controller object to be registered by the name of the object.
func RegisterController(controller interface{}) {
	registry.Store(getType(controller), reflect.ValueOf(controller))
}

//RegisterControllerByKey registers a controller object to be registered by a custom key.
func RegisterControllerByKey(key string, controller interface{}) {
	registry.Store(key, reflect.ValueOf(controller))
}

func getController(key string) reflect.Value {

	controller, ok := getControllerRegistry(key)
	if ok {
		return controller
	}
	return reflect.ValueOf(nil)
}

func getControllerRegistry(key string) (controller reflect.Value, ok bool) {

	obj, ok := registry.Load(key)
	if ok {
		controller = obj.(reflect.Value)
	}
	return
}

func getType(myvar interface{}) string {

	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}
