package app

import (
	"crypto/rand"
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
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	sync.Mutex
	Id         string
	Connection *websocket.Conn
}

type WebSocketConnectionCollection struct {
	sync.RWMutex
	Connections []*WebSocketConnection
}

type WebSocketCallbackSync struct {
	sync.RWMutex
	callbacks []WebSocketCallback
}

type WebSocketCallback func(conn *WebSocketConnection, c *gin.Context, messageType int, data []byte)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var WebSocketConnections WebSocketConnectionCollection
var WebSocketCallbacks WebSocketCallbackSync

func Initialize(path string) {
	serverSettings.Initialize(path)
	dbServices.Initialize()

	if serverSettings.WebConfig.Application.ReleaseMode == "release" {
		ginServer.Initialize(gin.ReleaseMode)
	} else {
		ginServer.Initialize(gin.DebugMode)
	}

	fileCache.Initialize()

}

func Run() {

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

	log.Println("Web Socket Connection")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	//Start the Reader, listen for Close Message, and Add to the Connection Array.

	wsConn := new(WebSocketConnection)
	wsConn.Connection = conn
	uuid, err := newUUID()
	if err == nil {
		wsConn.Id = uuid
	} else {
		uuid = randomString(20)
		wsConn.Id = uuid
	}

	//Reader
	go func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err == nil {
				go func() {
					WebSocketCallbacks.RLock()
					log.Println(string(p[:]))
					for _, callback := range WebSocketCallbacks.callbacks {
						callback(wsConn, c, messageType, p)
					}
					WebSocketCallbacks.RUnlock()
				}()
			} else {
				return
			}
		}
	}()

	//Close Message
	closeHandler := func(code int, text string) error {
		// log.Println("Closing Socket")
		// log.Printf("%+v\n", len(WebSocketConnections.Connections))
		deleteWebSocket(wsConn)
		// log.Printf("%+v\n", len(WebSocketConnections.Connections))
		return nil
	}

	conn.SetCloseHandler(closeHandler)

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

	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		if ws.Id == conn.Id {
			go func() {
				ws.Lock()
				ws.Connection.WriteMessage(websocket.BinaryMessage, data)
				ws.Unlock()
			}()
			return
		}
	}
	WebSocketConnections.RUnlock()
}

func ReplyToWebSocketJSON(conn *WebSocketConnection, v interface{}) {

	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		if ws.Id == conn.Id {
			go func() {
				ws.Lock()
				ws.Connection.WriteJSON(v)
				ws.Unlock()
			}()
			return
		}
	}
	WebSocketConnections.RUnlock()
}

func BroadcastWebSocketData(data []byte) {

	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		go func() {
			ws.Lock()
			ws.Connection.WriteMessage(websocket.BinaryMessage, data)
			ws.Unlock()
		}()
	}
	WebSocketConnections.RUnlock()
}

func BroadcastWebSocketJSON(v interface{}) {
	WebSocketConnections.RLock()
	for _, wsConn := range WebSocketConnections.Connections {
		ws := wsConn
		go func() {
			ws.Lock()
			ws.Connection.WriteJSON(v)
			ws.Unlock()
		}()
	}
	WebSocketConnections.RUnlock()
}

func deleteWebSocket(c *WebSocketConnection) {
	WebSocketConnections.Lock()
	for i, wsConn := range WebSocketConnections.Connections {
		if wsConn.Id == c.Id {
			log.Println("Removing Socket")
			WebSocketConnections.Connections = removeWebSocket(WebSocketConnections.Connections, i)
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
