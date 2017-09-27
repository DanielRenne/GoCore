package app

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	randMath "math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/ginServer"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketRemoval func(conn *WebSocketConnection)

type WebSocketConnection struct {
	sync.RWMutex
	Id            string
	Connection    *websocket.Conn
	Req           *http.Request
	Context       interface{}
	ContextString string
	ContextType   string
}

type WebSocketConnectionCollection struct {
	sync.RWMutex
	Connections []*WebSocketConnection
}

type WebSocketCallbackSync struct {
	sync.RWMutex
	callbacks []WebSocketCallback
}

type WebSocketPubSubPayload struct {
	Key     string      `json:"Key"`
	Content interface{} `json:"Content"`
}

type WebSocketCallback func(conn *WebSocketConnection, c *gin.Context, messageType int, data []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var WebSocketConnections WebSocketConnectionCollection
var WebSocketCallbacks WebSocketCallbackSync
var WebSocketRemovalCallback WebSocketRemoval

func Initialize(path string, config string, cookieDomain string) (err error) {
	err = serverSettings.Initialize(path, config)
	if err != nil {
		return
	}

	serverSettings.WebConfigMutex.RLock()
	inRelease := serverSettings.WebConfig.Application.ReleaseMode == "release"
	serverSettings.WebConfigMutex.RUnlock()

	if inRelease {
		ginServer.Initialize(gin.ReleaseMode, cookieDomain)
	} else {
		ginServer.Initialize(gin.DebugMode, cookieDomain)
	}
	fileCache.Initialize()

	err = dbServices.Initialize()
	if err != nil {
		return
	}
	return
}

func InitializeLite() (err error) {
	ginServer.InitializeLite(gin.ReleaseMode)
	fileCache.Initialize()
	return
}

func RunLite(port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at RunLite():  ", r)
			time.Sleep(time.Millisecond * 3000)
			RunLite(port)
			return
		}
	}()

	ginServer.Router.GET("/ws", func(c *gin.Context) {
		webSocketHandler(c.Writer, c.Request, c)
	})

	log.Println("GoCore Application Started")

	ginServer.Router.Run(":" + strconv.Itoa(port))

}

func Run() {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at Run():  ", r)
			time.Sleep(time.Millisecond * 3000)
			Run()
			return
		}
	}()

	if serverSettings.WebConfig.Application.WebServiceOnly == false {

		loadHTMLTemplates()

		ginServer.Router.Static("/web", serverSettings.APP_LOCATION+"/web")

		ginServer.Router.GET("/ws", func(c *gin.Context) {
			webSocketHandler(c.Writer, c.Request, c)
		})
	}

	initializeStaticRoutes()

	go ginServer.Router.RunTLS(":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort), serverSettings.APP_LOCATION+"/keys/cert.pem", serverSettings.APP_LOCATION+"/keys/key.pem")

	log.Println("GoCore Application Started")

	ginServer.Router.Run(":" + strconv.Itoa(serverSettings.WebConfig.Application.HttpPort))

	// go ginServer.Router.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort))
	// })

}

func webSocketHandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {

	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at webSocketHandler():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			webSocketHandler(w, r, c)
			return
		}
	}()
	//log.Println("Web Socket Connection")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Failed to upgrade http connection to websocket:  " + err.Error())
		return
	}

	//Start the Reader, listen for Close Message, and Add to the Connection Array.

	wsConn := new(WebSocketConnection)
	wsConn.Connection = conn
	wsConn.Req = r
	uuid, err := newUUID()
	if err == nil {
		wsConn.Id = uuid
	} else {
		uuid = randomString(20)
		wsConn.Id = uuid
	}

	// log.Println("Upgrading Websocket")
	// log.Printf("%+v\n", r)

	//Reader
	go logger.GoRoutineLogger(func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err == nil {
				go logger.GoRoutineLogger(func() {
					WebSocketCallbacks.RLock()
					for _, callback := range WebSocketCallbacks.callbacks {
						if callback != nil {
							callback(wsConn, c, messageType, p)
						}
					}
					WebSocketCallbacks.RUnlock()
				}, "GoCore/app.go->webSocketHandler[Callback calls]")
			} else {
				deleteWebSocket(wsConn)
				return
			}
		}
	}, "GoCore/app.go->webSocketHandler[Reader]")

	WebSocketConnections.Lock()
	WebSocketConnections.Connections = append(WebSocketConnections.Connections, wsConn)
	WebSocketConnections.Unlock()
}

func loadHTMLTemplates() {

	if serverSettings.WebConfig.Application.HtmlTemplates.Enabled {

		levels := "/*"
		dirLevel := ""

		switch serverSettings.WebConfig.Application.HtmlTemplates.DirectoryLevels {
		case 0:
			levels = "/*"
			dirLevel = ""
		case 1:
			levels = "/**/*"
			dirLevel = "root/"
		case 2:
			levels = "/**/**/*"
			dirLevel = "root/root/"
		}

		ginServer.Router.LoadHTMLGlob(serverSettings.APP_LOCATION + "/web/" + serverSettings.WebConfig.Application.HtmlTemplates.Directory + levels)

		ginServer.Router.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, dirLevel+"index.tmpl", gin.H{})
		})
	} else {

		if serverSettings.WebConfig.Application.DisableRootIndex {
			return
		}

		ginServer.Router.GET("", func(c *gin.Context) {
			if serverSettings.WebConfig.Application.RootIndexPath == "" {
				ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/index.htm", c)
			} else {
				ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/"+serverSettings.WebConfig.Application.RootIndexPath, c)
			}
		})
	}
}

func initializeStaticRoutes() {

	ginServer.Router.GET("/swagger", func(c *gin.Context) {
		// c.Redirect(http.StatusMovedPermanently, "https://"+serverSettings.WebConfig.Application.Domain+":"+strconv.Itoa(serverSettings.WebConfig.Application.HttpsPort)+"/web/swagger/dist/index.html")

		ginServer.ReadHTMLFile(serverSettings.APP_LOCATION+"/web/swagger/dist/index.html", c)
	})
}

func RegisterWebSocketDataCallback(callback WebSocketCallback) {
	WebSocketCallbacks.Lock()
	WebSocketCallbacks.callbacks = append(WebSocketCallbacks.callbacks, callback)
	WebSocketCallbacks.Unlock()
}

func ReplyToWebSocket(conn *WebSocketConnection, data []byte) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at ReplyToWebSocket():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			ReplyToWebSocket(conn, data)
			return
		}

		//clear out locks when this is done
		WebSocketConnections.RUnlock()
	}()
	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		if ws.Id == conn.Id {
			go logger.GoRoutineLogger(func() {

				ws.Lock()
				ws.Connection.WriteMessage(websocket.BinaryMessage, data)
				ws.Unlock()
			}, "GoCore/app.go->ReplyToWebSocket[WriteMessage]")
			return
		}
	}
}

func ReplyToWebSocketJSON(conn *WebSocketConnection, v interface{}) {

	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at ReplyToWebSocketJSON():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			ReplyToWebSocketJSON(conn, v)
			return
		}
		WebSocketConnections.RUnlock()
	}()
	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		if ws.Id == conn.Id {
			go logger.GoRoutineLogger(func() {
				ws.Lock()
				ws.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
				ws.Connection.WriteJSON(v)
				ws.Unlock()
			}, "GoCore/app.go->ReplyToWebSocketJSON[WriteJSON]")
			return
		}
	}
}

func ReplyToWebSocketPubSub(conn *WebSocketConnection, key string, v interface{}) {

	WebSocketConnections.RLock()
	defer WebSocketConnections.RUnlock()

	var payload WebSocketPubSubPayload
	payload.Key = key
	payload.Content = v

	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		if ws.Id == conn.Id {
			ws.Lock()
			ws.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			ws.Connection.WriteJSON(payload)
			ws.Unlock()
			return
		}
	}
}

func BroadcastWebSocketData(data []byte) {

	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at WebSocketConnections():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			BroadcastWebSocketData(data)
			return
		}
	}()

	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		go logger.GoRoutineLogger(func() {

			ws.Lock()
			ws.Connection.WriteMessage(websocket.BinaryMessage, data)
			ws.Unlock()
		}, "GoCore/app.go->BroadcastWebSocketData[WriteMessage]")
	}
	WebSocketConnections.RUnlock()
}

func BroadcastWebSocketJSON(v interface{}) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at BroadcastWebSocketJSON():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			BroadcastWebSocketJSON(v)
			return
		}
	}()
	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		go logger.GoRoutineLogger(func() {
			ws.Lock()
			ws.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			ws.Connection.WriteJSON(v)
			ws.Unlock()
		}, "GoCore/app.go->BroadcastWebSocketData[WriteJSON]")
	}
	WebSocketConnections.RUnlock()
}

func PublishWebSocketJSON(key string, v interface{}) {
	defer func() {
		if recover := recover(); recover != nil {
			log.Println("Panic Recovered at PublishWebSocketJSON():  ", recover)
			time.Sleep(time.Millisecond * 3000)
			PublishWebSocketJSON(key, v)
			return
		}
	}()
	var payload WebSocketPubSubPayload
	payload.Key = key
	payload.Content = v

	//Serialize and Deserialize to prevent Race Conditions from caller.
	data, _ := json.Marshal(payload)
	json.Unmarshal(data, &payload)

	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		go logger.GoRoutineLogger(func() {
			ws.Lock()
			ws.Connection.SetWriteDeadline(time.Now().Add(time.Duration(10000) * time.Millisecond))
			ws.Connection.WriteJSON(payload)
			ws.Unlock()
		}, "GoCore/app.go->WriteJSON")
	}
	WebSocketConnections.RUnlock()
}

func deleteWebSocket(c *WebSocketConnection) {
	WebSocketConnections.Lock()

	for i := 0; i < len(WebSocketConnections.Connections); i++ {
		wsConn := WebSocketConnections.Connections[i]
		if wsConn.Id == c.Id {
			log.Println("Deleting Web Socket from client:  " + wsConn.Connection.RemoteAddr().String())
			WebSocketConnections.Connections = removeWebSocket(WebSocketConnections.Connections, i)
			if WebSocketRemovalCallback != nil {
				go func(c *WebSocketConnection) {
					defer func() {
						if recover := recover(); recover != nil {
							log.Println("Panic Recovered at deleteWebSocket():  ", recover)
							return
						}
					}()
					WebSocketRemovalCallback(c)
				}(wsConn)
			}
			i--
		}
	}

	WebSocketConnections.Unlock()
}

func removeWebSocket(s []*WebSocketConnection, i int) []*WebSocketConnection {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func randomString(strlen int) string {
	randMath.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[randMath.Intn(len(chars))]
	}
	return string(result)
}
