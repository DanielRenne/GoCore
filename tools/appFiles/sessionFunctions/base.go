package session_functions

import (
	"encoding/base64"
	"strings"

	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"log"
	"sync"

	"runtime/debug"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/app"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/goCoreAppTemplate/constants"
	"github.com/DanielRenne/goCoreAppTemplate/models/v1/model"
	"github.com/DanielRenne/goCoreAppTemplate/scheduleEngine"
	"github.com/DanielRenne/goCoreAppTemplate/settings"
	"github.com/DanielRenne/goCoreAppTemplate/viewModel/socketViews"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//const LOG_VERBOSE = ""
//Logger.STANDARD_LEVELS = [ 'all', 'trace', 'debug', 'info', 'warn', 'error', 'fatal' ];
//Logger.DEFAULT_LEVEL = 'info';

type logBuffer struct {
	sync.RWMutex
	buffers map[string]string
}

var logBuffers logBuffer

type ServerResponseStruct struct {
	CompletedSuccessfully bool
	Context               RequestContext
	Redirect              string
	GlobalMessage         string
	GlobalMessageType     string
	Trace                 error
	TransactionId         string
	ViewModel             interface{}
}

func init() {
	logBuffers = logBuffer{
		buffers: make(map[string]string, 0),
	}
}

func ServerResponseToStruct(completedSuccessfully bool, context RequestContext, redirect string, globalMessage string, globalMessageType string, trace error, transactionId string, viewModel interface{}) ServerResponseStruct {
	return ServerResponseStruct{
		CompletedSuccessfully: completedSuccessfully,
		Context:               context,
		Redirect:              redirect,
		GlobalMessage:         globalMessage,
		GlobalMessageType:     globalMessageType,
		Trace:                 trace,
		TransactionId:         transactionId,
		ViewModel:             viewModel,
	}
}

type ServerResponse func(string, string, string, error, string, interface{})

type RequestContext func() *gin.Context

func GetSessionUser(c *gin.Context) (user model.User, err error) {
	if c == nil || c.Keys == nil {
		err = errors.New("Gin Context nil trying to get Session User.")
		return
	}
	authUserIdObj, ok := c.Keys[constants.COOKIE_AUTH_USER_ID]
	if ok {
		authUserId, parsed := authUserIdObj.(string)
		if parsed {
			err = model.Users.Query().ById(authUserId, &user)
		} else {
			err = errors.New("Gin Context COOKIE_AUTH_USER_ID did not parse to string")
		}

	} else {
		if c.Request == nil {
			err = errors.New("Gin Context.Request is nil trying to get Session User.")
			return
		}
		err = model.Users.Query().ById(ginServer.GetSessionKey(c, constants.COOKIE_AUTH_USER_ID), &user)
	}

	return
}

func GetSessionAccount(c *gin.Context) (account model.Account, err error) {
	if c == nil || c.Keys == nil {
		err = errors.New("Gin Context nil trying to get Session Account.")
		return
	}
	authAcctIdObj, ok := c.Keys[constants.COOKIE_AUTH_ACCOUNT_ID]
	if ok {
		authAcctId, parsed := authAcctIdObj.(string)
		if parsed {
			err = model.Accounts.Query().ById(authAcctId, &account)
		} else {
			err = errors.New("Gin Context COOKIE_AUTH_ACCOUNT_ID did not parse to string")
		}

	} else {
		if c.Request == nil {
			err = errors.New("Gin Context.Request is nil trying to get Session Account.")
			return
		}

		err = model.Accounts.Query().ById(ginServer.GetSessionKey(c, constants.COOKIE_AUTH_ACCOUNT_ID), &account)
	}

	return
}

func GetSessionAuthToken(c *gin.Context) (token string) {
	authTokenObj, ok := c.Keys[constants.COOKIE_AUTH_TOKEN]
	if ok {
		authToken, parsed := authTokenObj.(string)
		if parsed {
			token = authToken
		} else {
			token = ""
		}

	} else {
		token = ginServer.GetSessionKey(c, constants.COOKIE_AUTH_TOKEN)
	}

	return
}

func GetSessionDateCreated(c *gin.Context) (created string) {
	authDateCreatedObj, ok := c.Keys[constants.COOKIE_DATE_CREATED]
	if ok {
		authDateCreated, parsed := authDateCreatedObj.(string)
		if parsed {
			created = authDateCreated
		} else {
			created = ""
		}

	} else {
		created = ginServer.GetSessionKey(c, constants.COOKIE_DATE_CREATED)
	}

	return
}

func PassContext(c *gin.Context) RequestContext {

	return func() *gin.Context {
		return c
	}
}

func GetRedirect(controller string, action string, uriParams map[string]string) string {
	params := "{"
	for key, value := range uriParams {
		params += "\"" + key + "\":\"" + strings.Replace(value, "\"", "\"\"", -1) + "\","
	}
	params = params[:len(params)-1]
	params += "}"
	params = base64.StdEncoding.EncodeToString([]byte(params))
	return "/#/" + controller + "?action=" + action + "&uriParams=" + params

}

func GetPartialRedirect(controller string, action string, uriParams map[string]string) string {
	params := "{"
	for key, value := range uriParams {
		params += "\"" + key + "\":\"" + strings.Replace(value, "\"", "\"\"", -1) + "\","
	}
	params = params[:len(params)-1]
	params += "}"
	params = base64.StdEncoding.EncodeToString([]byte(params))
	return controller + "?action=" + action + "&uriParams=" + params

}

func GetSessionRole(c *gin.Context) (role model.Role, err error) {
	authRoleIdObj, ok := c.Keys[constants.COOKIE_AUTH_ROLE_ID]
	if ok {
		authRoleId, parsed := authRoleIdObj.(string)
		if parsed {
			err = model.Roles.Query().ById(authRoleId, &role)
		} else {
			err = errors.New("Gin Context COOKIE_AUTH_ROLE_ID did not parse to string")
		}
	} else {
		err = model.Roles.Query().ById(ginServer.GetSessionKey(c, constants.COOKIE_AUTH_ROLE_ID), &role)
	}

	return
}

func GetSessionAccountRole(c *gin.Context) (accountRole model.AccountRole, err error) {
	authAcctRoleIdObj, ok := c.Keys[constants.COOKIE_AUTH_ACCOUNTROLE_ID]
	if ok {
		authAcctRoleId, parsed := authAcctRoleIdObj.(string)
		if parsed {
			err = model.AccountRoles.Query().ById(authAcctRoleId, &accountRole)
		} else {
			err = errors.New("Gin Context COOKIE_AUTH_ACCOUNTROLE_ID did not parse to string")
		}

	} else {
		err = model.AccountRoles.Query().ById(ginServer.GetSessionKey(c, constants.COOKIE_AUTH_ACCOUNTROLE_ID), &accountRole)
	}
	return
}

func CheckRoleAccess(c *gin.Context, featureId string) (result bool) {
	var accountRole model.AccountRole
	var err error

	authAcctRoleIdObj, ok := c.Keys[constants.COOKIE_AUTH_ACCOUNTROLE_ID]
	if ok {
		authAcctRoleId, parsed := authAcctRoleIdObj.(string)
		if parsed {
			err = model.AccountRoles.Query().ById(authAcctRoleId, &accountRole)
		} else {
			err = errors.New("Gin Context COOKIE_AUTH_ACCOUNTROLE_ID did not parse to string")
		}

	} else {
		err = model.AccountRoles.Query().ById(ginServer.GetSessionKey(c, constants.COOKIE_AUTH_ACCOUNTROLE_ID), &accountRole)
	}

	if err != nil {
		return
	}

	authUserIdObj, ok := c.Keys[constants.COOKIE_AUTH_USER_ID]
	if ok {
		authUserId, parsed := authUserIdObj.(string)
		if parsed {
			result = CheckRoleAccessByRole(authUserId, accountRole.RoleId, featureId)
		} else {
			err = errors.New("Gin Context COOKIE_AUTH_USER_ID did not parse to string")
		}

	} else {
		result = CheckRoleAccessByRole(ginServer.GetSessionKey(c, constants.COOKIE_AUTH_USER_ID), accountRole.RoleId, featureId)
	}

	return result
}

func CheckRoleAccessByRole(userId string, roleId string, featureId string) (result bool) {

	filter := make(map[string]interface{}, 2)
	filter[model.FIELD_ACCOUNTROLE_USERID] = userId
	c1, err := model.AccountRoles.Query().Filter(filter).Count()
	if c1 > 0 {
		return true
	}
	if err == nil {
		count, err := model.RoleFeatures.Query().Filter(model.Q(model.FIELD_ROLEFEATURE_ROLEID, roleId)).Filter(model.Q(model.FIELD_ROLEFEATURE_FEATUREID, featureId)).Count()
		if err != nil {
			return
		}

		if count >= 1 {
			// just in case there is a data issue
			if count != 1 {
				Dump("There are more than one RoleFeature rows when there should be one for " + featureId + "!!!!!!!!!!!!")
			}
			result = true
		}
	}

	return result
}

func BlockByRoleAccess(c *gin.Context, featureId string) (blocked bool) {
	if !CheckRoleAccess(c, featureId) {
		Dump(featureId + " BLOCKED!!")
		ginServer.RenderHTML(constants.HTTP_NOT_AUTHORIZED, c)
		c.AbortWithError(403, errors.New("Role Not Authorized"))
		blocked = true
		return
	}
	return
}

func GetProtocol(c *gin.Context) (schema string) {
	if c.Request.TLS == nil {
		schema = "http://"
	} else {
		schema = "https://"
	}
	return
}

func BlockByDeveloperModeOff(c *gin.Context) (blocked bool) {
	if !settings.AppSettings.DeveloperMode && settings.AppSettings.DemoMode == false {
		ginServer.RenderHTML(constants.HTTP_NOT_AUTHORIZED, c)
		c.AbortWithError(403, errors.New("Role Not Authorized"))
		blocked = true
		return
	}
	return
}

func BlockByAccountOwnerShip(c *gin.Context, rowAccount string) (blocked bool) {
	if rowAccount == "" {
		return
	}
	act, err := GetSessionAccount(c)
	if err == nil && rowAccount != act.Id.Hex() {
		ginServer.RenderHTML(constants.HTTP_NOT_AUTHORIZED, c)
		c.AbortWithError(403, errors.New("AccountId does not match this row"))
		blocked = true
		return
	}
	return
}

func CheckAccountOwnerShip(c *gin.Context, rowAccount string) (blocked bool) {
	if rowAccount == "" {
		return
	}
	act, err := GetSessionAccount(c)
	//_IR is for a special ID in equipment.
	if err == nil && rowAccount != act.Id.Hex() && rowAccount != act.Id.Hex()+"_IR" && !act.IsSystemAccount {
		blocked = true
		return
	}
	return
}

func StartTransaction(c *gin.Context) (t *model.Transaction, err error) {

	if c == nil {
		t, err = model.Transactions.New(constants.APP_CONSTANTS_USERS_ANONYMOUS_ID)
		return
	}

	user, err := GetSessionUser(c)

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	t, err = model.Transactions.New(user.Id.Hex())

	Log("StartTransaction", user.First+" "+user.Last+" Id: "+t.Id.Hex())
	if settings.AppSettings.DeveloperMode {
		core.Debug.Dump("Stack")
	}
	return
}

func StartTransactionWithUser(user model.User) (t *model.Transaction, err error) {
	t, err = model.Transactions.New(user.Id.Hex())
	return
}

func StoreDataFormat(c *gin.Context, language string, timeZone string, dateFormat string) {
	var df model.DataFormat
	df.Language = language
	df.DateFormat = dateFormat
	df.LocalTimeZone = timeZone
	dfJson, err := df.JSONString()

	_, ok := c.Keys[constants.COOKIE_DATA_FORMAT]

	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		if ok {
			c.Keys[constants.COOKIE_DATA_FORMAT] = "{\"Language\":\"en\",\"DateFormat\":\"mm/dd/yyyy\",\"TimeZone\":\"US/Eastern\"}"
		} else {
			ginServer.SetSessionKey(c, constants.COOKIE_DATA_FORMAT, "{\"Language\":\"en\",\"DateFormat\":\"mm/dd/yyyy\",\"TimeZone\":\"US/Eastern\"}")
		}

		return
	}
	if ok {
		c.Keys[constants.COOKIE_DATA_FORMAT] = dfJson
	} else {
		ginServer.SetSessionKey(c, constants.COOKIE_DATA_FORMAT, dfJson)
	}

}

func GetDataFormat(c *gin.Context) model.DataFormat {

	key := ""
	authDataFormatObj, ok := c.Keys[constants.COOKIE_DATA_FORMAT]
	if ok {
		authDataFormat, parsed := authDataFormatObj.(string)
		if parsed {
			key = authDataFormat
		} else {
			key = ""
		}
	} else {
		key = ginServer.GetSessionKey(c, constants.COOKIE_DATA_FORMAT)
	}

	if key == "" {
		return model.DataFormat{Language: "en", LocalTimeZone: "US/Eastern", DateFormat: "mm/dd/yyyy"}
	}
	var df model.DataFormat
	err := df.Parse(key)
	if err != nil {
		return model.DataFormat{Language: "en", LocalTimeZone: "US/Eastern", DateFormat: "mm/dd/yyyy"}
	}
	return df
}

// Logs out a hex dump and should be the most common thing you want to use in a TCP connection log
func LogWithoutQuotes(uniqueId string, desc string, message string) {
	logWithQuote(uniqueId, desc, message, false, false, false)
	return
}


func LogHex(uniqueId string, desc string, message string) {
	logWithQuote(uniqueId, desc, message, true, false, true)
	return
}

func Log(desc string, message string) {
	logWithQuote("", desc, message, true, false, false)
	return
}

func Dump(values ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("\n\nPanic Stack: " + string(debug.Stack()))
			log.Println("Panic Recovered at Dump():  ", r)
			return
		}
	}()

	var newValues []interface{}

	for _, message := range values {
		if strings.TrimSpace(fmt.Sprintf("%T", message)) == "string" && !extensions.IsPrintable(message.(string)) {
			msg := strings.Replace(message.(string), "\x00", "NULL", -1)
			dataToSend := make([]byte, uint32(len(msg)))
			binary.LittleEndian.PutUint32(dataToSend[0:], uint32(len(msg)))
			copy(dataToSend[0:], []byte(msg))
			output := hex.Dump(dataToSend[:])
			newValues = append(newValues, output)
		} else {
			newValues = append(newValues, message)
		}
	}

	logWithQuote("", "", core.Debug.GetDumpWithInfoAndTimeString(scheduleEngine.GetLocalTime(time.Now()).Format(time.RFC1123), newValues...), false, true, false)
	return
}

func logWithQuote(uniqueId string, desc string, message string, includeQuotes bool, skipBuildDesc bool, printDump bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("\n\nPanic Stack: " + string(debug.Stack()))
			log.Println("Panic Recovered at logWithQuote():  ", r)
			return
		}
	}()

	go logger.GoRoutineLogger(func() {

		defer func() {
			if r := recover(); r != nil {
				log.Println("Panic Recovered at logWithQuote():  ", r)
				return
			}
		}()

		var output string
		if !skipBuildDesc {
			if !extensions.IsPrintable(message) || printDump == true {
				message = strings.Replace(message, "\x00", "NULL", -1)
				dataToSend := make([]byte, uint32(len(message)))
				binary.LittleEndian.PutUint32(dataToSend[0:], uint32(len(message)))
				copy(dataToSend[0:], []byte(message))
				output = hex.Dump(dataToSend[:])
			} else {
				if message != "" {
					if len(message) > 2 {
						lastTwo := message[len(message)-2:]
						lastTwo = strings.Replace(lastTwo, "\r", "\\r", -1)
						lastTwo = strings.Replace(lastTwo, "\n", "\\n", -1)
						coreMessage := message[:len(message)-2]
						coreMessage = strings.Replace(coreMessage, "\n", "\n\t", -1)
						coreMessage = strings.Replace(coreMessage, "\r", "\r\t", -1)
						if includeQuotes {
							output = "\"" + coreMessage + lastTwo + "\""
						} else {
							output = coreMessage + lastTwo
						}

					} else {
						if includeQuotes {
							output = "\"" + message + "\""
						} else {
							output = message
						}
					}
				}
			}
			desc = "\n\n" + scheduleEngine.GetLocalTime(time.Now()).Format(time.RFC1123) + "\t\n\t" + fmt.Sprintf("#### %-30s\t####", desc) + "\n\t" + output
		} else {
			desc = message
		}

		//Ensure this is run first as we concatenate the entire log so there are no delays in the time a new entry is put into the mutex.
		if settings.AppSettings.DeveloperMode {
			WriteLog(uniqueId, desc)
		} else {
			logBuffers.Lock()
			_, ok := logBuffers.buffers[uniqueId]
			if !ok {
				logBuffers.buffers[uniqueId] = ""
			}
			logBuffers.buffers[uniqueId] += desc
			logBuffers.Unlock()
		}

		if uniqueId == "" {
			BroadcastLog("app", desc)
		} else {
			BroadcastLog(uniqueId, desc)
		}
	}, "logWithQuote")
	return
}

func Println(logInfo string) {
	logWithQuote("", "", logInfo+"\n", false, true, false)
	return
}

func Print(logInfo string) {
	logWithQuote("", "", logInfo, false, true, false)
	return
}

func FlushAllLogs() (err error) {
	logBuffers.Lock()
	for logId, contents := range logBuffers.buffers {
		err = WriteLog(logId, contents)
		logBuffers.buffers[logId] = ""
	}
	logBuffers.Unlock()
	return
}

func FlushLog(uniqueId string) (err error) {
	logBuffers.Lock()
	val, ok := logBuffers.buffers[uniqueId]
	if !ok {
		return
	}
	err = WriteLog(uniqueId, val)
	logBuffers.buffers[uniqueId] = ""
	logBuffers.Unlock()
	return
}

func WriteLog(uniqueId string, message string) (err error) {

	log := "src/github.com/DanielRenne/goCoreAppTemplate/log/plugins/" + uniqueId + ".log"
	if uniqueId == "" {
		log = "src/github.com/DanielRenne/goCoreAppTemplate/log/app.log"
	}

	var mode int
	if extensions.DoesFileNotExist(log) {
		mode = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	} else {
		mode = os.O_APPEND | os.O_WRONLY
	}
	f2, err := os.OpenFile(log, mode, 0777)
	if err == nil {
		defer f2.Close()
		w := bufio.NewWriter(f2)
		fmt.Fprint(w, message)
		err = w.Flush()
	}
	return
}

func BroadcastLog(id string, data string) {
	type LogData struct {
		Id   string `json:"Id"`
		Data string `json:"Data"`
	}
	var jsonResponse LogData
	jsonResponse.Id = id
	jsonResponse.Data = data

	var connections []*app.WebSocketConnection
	app.WebSocketConnections.RLock()
	for i := range app.WebSocketConnections.Connections {
		c := app.WebSocketConnections.Connections[i]
		connections = append(connections, c)
	}

	app.WebSocketConnections.RUnlock()

	for i := range connections {
		conn := connections[i]
		conn.RLock()
		value, ok := ParseSocketClientStatus(conn)
		conn.RUnlock()
		if ok == true && (value.Page == "logs" || value.Page == "\"logs\"") {
			app.ReplyToWebSocketPubSub(conn, "LogData", jsonResponse)
		}
	}
}

func BroadcastTime(date string, time string) {
	type timeData struct {
		Date string `json:"Date"`
		Time string `json:"Time"`
	}
	var jsonResponse timeData
	jsonResponse.Time = time
	jsonResponse.Date = date

	var connections []*app.WebSocketConnection
	app.WebSocketConnections.RLock()
	for i := range app.WebSocketConnections.Connections {
		c := app.WebSocketConnections.Connections[i]
		connections = append(connections, c)
	}

	app.WebSocketConnections.RUnlock()

	for i := range connections {
		conn := connections[i]
		conn.RLock()
		value, ok := ParseSocketClientStatus(conn)
		conn.RUnlock()
		if ok == true && (value.Page == "serverSettingsModify" || value.Page == "\"serverSettingsModify\"") {
			app.ReplyToWebSocketPubSub(conn, "Clock", jsonResponse)
		}
	}
}

func ParseSocketClientStatus(conn *app.WebSocketConnection) (value socketViews.ClientStatus, ok bool) {
	if conn.ContextType == "ClientStatus" {
		ok = true
		value.Parse(conn.ContextString)
		return
	}
	return
}
