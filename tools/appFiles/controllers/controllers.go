package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/DanielRenne/goCoreAppTemplate/sessionFunctions"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel"

	"github.com/gin-gonic/gin"
)

var registry sync.Map

func init() {
	gob.Register(GinContextPayload{})
	gob.Register(url.URL{})
	gob.Register(url.Userinfo{})
	gob.Register(gin.Param{})
}

type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func getController(key string) reflect.Value {

	controller, ok := getControllerRegistry(key)
	if ok {
		return controller
	}

	switch key {
	case CONTROLLER_LOGIN:
		return reflect.ValueOf(&LoginController{})
	case CONTROLLER_HOME:
		return reflect.ValueOf(&HomeController{})
	case CONTROLLER_USERS:
		return reflect.ValueOf(&UsersController{})
	case CONTROLLER_SETTINGS:
		return reflect.ValueOf(&SettingsController{})
	case CONTROLLER_ACCOUNTS:
		return reflect.ValueOf(&AccountsController{})
	case CONTROLLER_ACCOUNTLIST:
		return reflect.ValueOf(&AccountListController{})
	case CONTROLLER_ACCOUNTMODIFY:
		return reflect.ValueOf(&AccountModifyController{})
	case CONTROLLER_USERMODIFY:
		return reflect.ValueOf(&UserModifyController{})
	case CONTROLLER_SERVERSETTINGSMODIFY:
		return reflect.ValueOf(&ServerSettingsModifyController{})
	case CONTROLLER_USERLIST:
		return reflect.ValueOf(&UserListController{})
	case CONTROLLER_ACCOUNTADD:
		return reflect.ValueOf(&AccountAddController{})
	case CONTROLLER_USERADD:
		return reflect.ValueOf(&UserAddController{})
	case CONTROLLER_USERPROFILE:
		return reflect.ValueOf(&UserProfileController{})
	case CONTROLLER_PASSWORDRESET:
		return reflect.ValueOf(&PasswordResetController{})
	case CONTROLLER_TRANSACTIONS:
		return reflect.ValueOf(&TransactionsController{})
	case CONTROLLER_TRANSACTIONMODIFY:
		return reflect.ValueOf(&TransactionModifyController{})
	case CONTROLLER_TRANSACTIONLIST:
		return reflect.ValueOf(&TransactionListController{})
	case CONTROLLER_TRANSACTIONADD:
		return reflect.ValueOf(&TransactionAddController{})
	case CONTROLLER_APPERRORS:
		return reflect.ValueOf(&AppErrorsController{})
	case CONTROLLER_APPERRORMODIFY:
		return reflect.ValueOf(&AppErrorModifyController{})
	case CONTROLLER_APPERRORLIST:
		return reflect.ValueOf(&AppErrorListController{})
	case CONTROLLER_APPERRORADD:
		return reflect.ValueOf(&AppErrorAddController{})
	case CONTROLLER_FILEUPLOAD:
		return reflect.ValueOf(&FileUploadController{})
	case CONTROLLER_FEATURES:
		return reflect.ValueOf(&FeaturesController{})
	case CONTROLLER_FEATUREMODIFY:
		return reflect.ValueOf(&FeatureModifyController{})
	case CONTROLLER_FEATURELIST:
		return reflect.ValueOf(&FeatureListController{})
	case CONTROLLER_FEATUREADD:
		return reflect.ValueOf(&FeatureAddController{})
	case CONTROLLER_ROLEFEATURES:
		return reflect.ValueOf(&RoleFeaturesController{})
	case CONTROLLER_ROLEFEATUREMODIFY:
		return reflect.ValueOf(&RoleFeatureModifyController{})
	case CONTROLLER_ROLEFEATURELIST:
		return reflect.ValueOf(&RoleFeatureListController{})
	case CONTROLLER_ROLEFEATUREADD:
		return reflect.ValueOf(&RoleFeatureAddController{})
	case CONTROLLER_FEATUREGROUPS:
		return reflect.ValueOf(&FeatureGroupsController{})
	case CONTROLLER_FEATUREGROUPMODIFY:
		return reflect.ValueOf(&FeatureGroupModifyController{})
	case CONTROLLER_FEATUREGROUPLIST:
		return reflect.ValueOf(&FeatureGroupListController{})
	case CONTROLLER_FEATUREGROUPADD:
		return reflect.ValueOf(&FeatureGroupAddController{})
	case CONTROLLER_ROLES:
		return reflect.ValueOf(&RolesController{})
	case CONTROLLER_ROLEMODIFY:
		return reflect.ValueOf(&RoleModifyController{})
	case CONTROLLER_ROLELIST:
		return reflect.ValueOf(&RoleListController{})
	case CONTROLLER_ROLEADD:
		return reflect.ValueOf(&RoleAddController{})
	case CONTROLLER_FILEOBJECTS:
		return reflect.ValueOf(&FileObjectsController{})
	case CONTROLLER_FILEOBJECTMODIFY:
		return reflect.ValueOf(&FileObjectModifyController{})
	case CONTROLLER_FILEOBJECTLIST:
		return reflect.ValueOf(&FileObjectListController{})
	case CONTROLLER_FILEOBJECTADD:
		return reflect.ValueOf(&FileObjectAddController{})
	case CONTROLLER_LOGS:
		return reflect.ValueOf(&LogsController{})
		//-DONT-REMOVE-NEW-CASE
	}
	return reflect.ValueOf(&HomeController{})

}

func renderStandardSideBar(context session_functions.RequestContext, sideBarMenu *viewModel.SideBarViewModel) (err error) {

	if sideBarMenu == nil {
		err = errors.New("SideBarMenu is a nil pointer.")
		return
	}

	user, err := session_functions.GetSessionUser(context())
	account, err := session_functions.GetSessionAccount(context())

	if err == nil {
		err = sideBarMenu.RenderApp(context(), account, user)
	} else {
		sideBarMenu.LoadDefaultState(context())
	}

	return
}

const (
	GIN_CONTEXT_METHOD_CLEAR_SESSION = "ClearSession"
	GIN_CONTEXT_METHOD_SAVE_SESSION  = "SaveSession"
	GIN_CONTEXT_METHOD_ABORT         = "AbortWithError"
	GIN_CONTEXT_METHOD_REDIRECT      = "Redirect"
	GIN_CONTEXT_METHOD_RENDERHTML    = "RenderHTML"
)

type GinContextPayload struct {
	Request struct {
		Header map[string][]string
		URL    *url.URL
	}
	Keys   map[string]interface{}
	Params []gin.Param
}

type GinCloudProxy struct {
	Method string
	Status int
}

func EncodeContextGOB(c *gin.Context) (data []byte, err error) {

	var payload GinContextPayload
	payload.Request.Header = c.Request.Header
	payload.Request.URL = c.Request.URL
	payload.Params = c.Params
	payload.Keys = make(map[string]interface{})

	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err = e.Encode(payload)
	data = b.Bytes()
	return
}

func DecodeContextGOB(data []byte, x *gin.Context) (err error) {
	var payload GinContextPayload
	// payload.Request.Header = make(map[string][]string)

	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	err = dec.Decode(&payload)
	if err != nil {
		return
	}

	x.Request.Header = payload.Request.Header
	x.Request.URL = payload.Request.URL
	x.Keys = payload.Keys
	x.Params = payload.Params

	return
}

func NewGinContext() (c *gin.Context) {
	var g gin.Context
	g.Request = &http.Request{}
	c = &g
	return
}

type meta struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Version string `json:"version"`
}

func (self *meta) Stringify() (value string, err error) {
	data, err := json.Marshal(self)
	value = string(data)
	return
}

func RegisterController(controller interface{}) {
	registry.Store(getType(controller), reflect.ValueOf(controller))
}

func RegisterControllerByKey(key string, controller interface{}) {
	registry.Store(key, reflect.ValueOf(controller))
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
